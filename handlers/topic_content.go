package handlers

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/scylladb/gocqlx/v2/qb"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-course-creator/constants"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/lib/db/bucket"
	"github.com/zicops/zicops-course-creator/lib/googleprojectlib"
)

func TopicContentCreate(ctx context.Context, topicID string, courseID string, topicConent *model.TopicContentInput) (*model.TopicContent, error) {
	log.Info("TopicContentCreate called")
	guid := xid.New()
	cassandraTopicContent := coursez.TopicContent{
		ID:                 guid.String(),
		TopicId:            topicID,
		CourseId:           courseID,
		Language:           *topicConent.Language,
		CreatedAt:          time.Now().Unix(),
		UpdatedAt:          time.Now().Unix(),
		TopicContentBucket: "",
		Url:                "",
		SubtitleFile:       "",
		IsActive:           false,
	}
	if topicConent.StartTime != nil {
		cassandraTopicContent.StartTime = *topicConent.StartTime
	}
	if topicConent.Duration != nil {
		cassandraTopicContent.Duration = *topicConent.Duration
	}
	if topicConent.SkipIntroDuration != nil {
		cassandraTopicContent.SkipIntroDuration = *topicConent.SkipIntroDuration
	}
	if topicConent.NextShowTime != nil {
		cassandraTopicContent.NextShowtime = *topicConent.NextShowTime
	}
	if topicConent.FromEndTime != nil {
		cassandraTopicContent.FromEndTime = *topicConent.FromEndTime
	}
	if topicConent.Type != nil {
		cassandraTopicContent.Type = *topicConent.Type
	}
	if topicConent.IsDefault != nil {
		cassandraTopicContent.IsDefault = *topicConent.IsDefault
	}
	// set course in cassandra
	insertQuery := global.CassSession.Session.Query(coursez.TopicContentTable.Insert()).BindStruct(cassandraTopicContent)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraTopicContent.CreatedAt, 10)
	responseModel := model.TopicContent{
		ID:                &cassandraTopicContent.ID,
		Language:          topicConent.Language,
		StartTime:         topicConent.StartTime,
		CreatedAt:         &created,
		UpdatedAt:         &created,
		Duration:          topicConent.Duration,
		SkipIntroDuration: topicConent.SkipIntroDuration,
		NextShowTime:      topicConent.NextShowTime,
		FromEndTime:       topicConent.FromEndTime,
		TopicID:           &topicID,
		Type:              topicConent.Type,
		IsDefault:         topicConent.IsDefault,
	}
	return &responseModel, nil
}

func TopicExamCreate(ctx context.Context, topicID string, courseID string, exam *model.TopicExamInput) (*model.TopicExam, error) {
	log.Info("TopicExamCreate called")
	guid := xid.New()
	cassandraTopicContent := coursez.TopicExam{
		ID:        guid.String(),
		TopicId:   topicID,
		CourseId:  courseID,
		ExamId:    *exam.ExamID,
		Language:  *exam.Language,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		IsActive:  false,
	}
	// set course in cassandra
	insertQuery := global.CassSession.Session.Query(coursez.TopicExamTable.Insert()).BindStruct(cassandraTopicContent)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraTopicContent.CreatedAt, 10)
	responseModel := model.TopicExam{
		ID:        &cassandraTopicContent.ID,
		Language:  exam.Language,
		CourseID:  &courseID,
		CreatedAt: &created,
		UpdatedAt: &created,
		TopicID:   &topicID,
		ExamID:    exam.ExamID,
	}
	return &responseModel, nil
}

func UploadTopicVideo(ctx context.Context, file model.TopicVideo) (*model.UploadResult, error) {
	log.Info("UploadTopicVideo called")
	isSuccess := model.UploadResult{}
	storageC := bucket.NewStorageHandler()
	gproject := googleprojectlib.GetGoogleProjectID()
	err := storageC.InitializeStorageClient(ctx, gproject)
	if err != nil {
		log.Errorf("Failed to upload video to course topic: %v", err.Error())
		return &isSuccess, nil
	}
	if file.CourseID == nil || file.ContentID == nil {
		log.Errorf("Failed to upload video to course topic: %v", "courseID or contentId is nil")
		return &isSuccess, nil
	}
	bucketPath := *file.CourseID + "/" + *file.ContentID + "/" + file.File.Filename
	writer, err := storageC.UploadToGCS(ctx, bucketPath, map[string]string{})
	if err != nil {
		log.Errorf("Failed to upload video to course topic: %v", err.Error())
		return &isSuccess, nil
	}
	defer writer.Close()
	fileBuffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(fileBuffer, file.File.File); err != nil {
		return &isSuccess, nil
	}
	currentBytes := fileBuffer.Bytes()
	_, err = io.Copy(writer, bytes.NewReader(currentBytes))
	if err != nil {
		return &isSuccess, err
	}
	getUrl := storageC.GetSignedURLForObject(bucketPath)
	where := qb.Eq("id")
	updateQB := qb.Update("coursez.topic_content").Set("topiccontentbucket").Set("url").Where(where)
	updateQuery := updateQB.Query(*global.CassSession.Session).BindMap(qb.M{"id": file.ContentID, "topiccontentbucket": bucketPath, "url": getUrl})
	if err := updateQuery.ExecRelease(); err != nil {
		return &isSuccess, err
	}
	isSuccessRes := true
	isSuccess.Success = &isSuccessRes
	isSuccess.URL = &getUrl
	return &isSuccess, nil
}

func UploadTopicSubtitle(ctx context.Context, files []*model.TopicSubtitle) ([]*model.UploadResultSubtitles, error) {
	log.Info("UploadTopicSubtitle called")
	isSuccess := []*model.UploadResultSubtitles{}
	for _, file := range files {
		isLocalSuccess := model.UploadResultSubtitles{}
		isLocal := false
		isLocalSuccess.Success = &isLocal
		isLocalSuccess.URL = nil
		storageC := bucket.NewStorageHandler()
		gproject := googleprojectlib.GetGoogleProjectID()
		err := storageC.InitializeStorageClient(ctx, gproject)
		if err != nil {
			log.Errorf("Failed to upload subtitle to course topic: %v", err.Error())
			isSuccess = append(isSuccess, &isLocalSuccess)
			continue
		}
		language := "en"
		if file.Language != nil {
			language = *file.Language
		}
		mainBucket := *file.CourseID + "/" + *file.TopicID + "/subtitles/"
		bucketPath := mainBucket + file.File.Filename
		writer, err := storageC.UploadToGCS(ctx, bucketPath, map[string]string{"language": language})
		if err != nil {
			log.Errorf("Failed to upload subtitle to course topic: %v", err.Error())
			isSuccess = append(isSuccess, &isLocalSuccess)
			continue
		}
		fileBuffer := bytes.NewBuffer(nil)
		if _, err := io.Copy(fileBuffer, file.File.File); err != nil {
			isSuccess = append(isSuccess, &isLocalSuccess)
			continue
		}
		currentBytes := fileBuffer.Bytes()
		_, err = io.Copy(writer, bytes.NewReader(currentBytes))
		if err != nil {
			isSuccess = append(isSuccess, &isLocalSuccess)
			continue
		}
		writer.Close()
		err = storageC.SetTags(ctx, bucketPath, map[string]string{"language": language})
		if err != nil {
			isSuccess = append(isSuccess, &isLocalSuccess)
			continue
		}
		getUrl := storageC.GetSignedURLForObject(bucketPath)
		isLocal = true
		isLocalSuccess.Success = &isLocal
		isLocalSuccess.URL = &getUrl
		isLocalSuccess.Language = file.Language
		isSuccess = append(isSuccess, &isLocalSuccess)
	}

	return isSuccess, nil
}

func UpdateTopicContent(ctx context.Context, topicConent *model.TopicContentInput) (*model.TopicContent, error) {
	log.Info("UpdateTopicContent called")
	topicID := topicConent.ContentID
	if *topicID == "" {
		return nil, fmt.Errorf("ContentID is required")
	}
	cassandraTopicContent := coursez.TopicContent{
		ID: *topicID,
	}
	topicContents := []coursez.TopicContent{}
	getQuery := global.CassSession.Session.Query(coursez.TopicContentTable.Get()).BindMap(qb.M{"id": cassandraTopicContent.ID})
	if err := getQuery.SelectRelease(&topicContents); err != nil {
		return nil, err
	}
	if len(topicContents) < 1 {
		return nil, fmt.Errorf("quiz not found")
	}
	cassandraTopicContent = topicContents[0]
	updateCols := []string{}
	if topicConent.Duration != nil {
		updateCols = append(updateCols, "duration")
		cassandraTopicContent.Duration = *topicConent.Duration
	}
	if topicConent.StartTime != nil {
		updateCols = append(updateCols, "starttime")
		cassandraTopicContent.StartTime = *topicConent.StartTime
	}
	if topicConent.SkipIntroDuration != nil {
		updateCols = append(updateCols, "skipintroduration")
		cassandraTopicContent.SkipIntroDuration = *topicConent.SkipIntroDuration
	}
	if topicConent.Type != nil {
		updateCols = append(updateCols, "type")
		cassandraTopicContent.Type = *topicConent.Type
	}
	if topicConent.NextShowTime != nil {
		updateCols = append(updateCols, "nextshowtime")
		cassandraTopicContent.NextShowtime = *topicConent.NextShowTime
	}
	if topicConent.FromEndTime != nil {
		updateCols = append(updateCols, "fromendtime")
		cassandraTopicContent.FromEndTime = *topicConent.FromEndTime
	}
	updateCols = append(updateCols, "updated_at")
	cassandraTopicContent.UpdatedAt = time.Now().Unix()
	if topicConent.Language != nil {
		updateCols = append(updateCols, "language")
		cassandraTopicContent.Language = *topicConent.Language
	}
	if topicConent.IsDefault != nil {
		updateCols = append(updateCols, "is_default")
		cassandraTopicContent.IsDefault = *topicConent.IsDefault
	}
	// set course in cassandra

	upStms, uNames := coursez.TopicContentTable.Update(updateCols...)
	updateQuery := global.CassSession.Session.Query(upStms, uNames).BindStruct(&cassandraTopicContent)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraTopicContent.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraTopicContent.UpdatedAt, 10)
	responseModel := model.TopicContent{
		Language:          topicConent.Language,
		StartTime:         topicConent.StartTime,
		CreatedAt:         &created,
		UpdatedAt:         &updated,
		Duration:          topicConent.Duration,
		SkipIntroDuration: topicConent.SkipIntroDuration,
		NextShowTime:      topicConent.NextShowTime,
		FromEndTime:       topicConent.FromEndTime,
		TopicID:           topicID,
		Type:              topicConent.Type,
		IsDefault:         topicConent.IsDefault,
	}
	return &responseModel, nil
}

func UpdateTopicExam(ctx context.Context, exam *model.TopicExamInput) (*model.TopicExam, error) {
	log.Info("UpdateTopicExam called")
	tExamId := exam.ID
	if *tExamId == "" {
		return nil, fmt.Errorf("TopicExamId is required")
	}
	cassandraTopicContent := coursez.TopicExam{
		ID: *tExamId,
	}
	topicExams := []coursez.TopicExam{}
	getQuery := global.CassSession.Session.Query(coursez.TopicExamTable.Get()).BindMap(qb.M{"id": cassandraTopicContent.ID})
	if err := getQuery.SelectRelease(&topicExams); err != nil {
		return nil, err
	}
	if len(topicExams) < 1 {
		return nil, fmt.Errorf("quiz not found")
	}
	cassandraTopicContent = topicExams[0]
	updateCols := []string{}
	if exam.Language != nil {
		updateCols = append(updateCols, "language")
		cassandraTopicContent.Language = *exam.Language
	}
	if exam.TopicID != nil {
		updateCols = append(updateCols, "topicid")
		cassandraTopicContent.TopicId = *exam.TopicID
	}
	if exam.ExamID != nil {
		updateCols = append(updateCols, "examid")
		cassandraTopicContent.ExamId = *exam.ExamID
	}
	updateCols = append(updateCols, "updated_at")
	cassandraTopicContent.UpdatedAt = time.Now().Unix()
	upStms, uNames := coursez.TopicExamTable.Update(updateCols...)
	updateQuery := global.CassSession.Session.Query(upStms, uNames).BindStruct(&cassandraTopicContent)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraTopicContent.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraTopicContent.UpdatedAt, 10)
	responseModel := model.TopicExam{
		Language:  exam.Language,
		TopicID:   exam.TopicID,
		ExamID:    exam.ExamID,
		CreatedAt: &created,
		UpdatedAt: &updated,
		ID:        tExamId,
	}
	return &responseModel, nil
}

func UploadTopicStaticContent(ctx context.Context, file *model.StaticContent) (*model.UploadResult, error) {
	log.Info("UploadTopicStaticContent called")
	isSuccess := model.UploadResult{}
	bucketPath := ""
	getUrl := ""
	if (file.URL == nil || *file.URL == "") && file.File != nil {
		storageC := bucket.NewStorageHandler()
		gproject := googleprojectlib.GetGoogleProjectID()
		err := storageC.InitializeStorageClient(ctx, gproject)
		if err != nil {
			log.Errorf("Failed to upload static content to course topic: %v", err.Error())
			return &isSuccess, nil
		}
		baseDir := strings.TrimSuffix(file.File.Filename, filepath.Ext(file.File.Filename))
		baseDir = strings.Split(baseDir, ".")[0]
		// convert baseDir to cryptographic hash
		hash := sha256.New()
		hash.Write([]byte(baseDir))
		hashBytes := hash.Sum(nil)
		hashString := hex.EncodeToString(hashBytes)
		bucketPath = *file.CourseID + "/" + *file.ContentID + "/" + hashString
		zipPath := bucketPath + "/" + file.File.Filename
		writer, err := storageC.UploadToGCSPub(ctx, zipPath, map[string]string{})
		if err != nil {
			log.Errorf("Failed to upload static content to course topic: %v", err.Error())
			return &isSuccess, nil
		}
		defer writer.Close()
		fileBuffer := bytes.NewBuffer(nil)
		if _, err := io.Copy(fileBuffer, file.File.File); err != nil {
			return &isSuccess, nil
		}
		currentBytes := fileBuffer.Bytes()
		newReader := bytes.NewReader(currentBytes)
		_, err = io.Copy(writer, bytes.NewReader(currentBytes))
		if err != nil {
			return &isSuccess, err
		}
		b, err := ioutil.ReadAll(newReader)
		if err != nil {
			return nil, err
		}
		zr, err := zip.NewReader(newReader, int64(len(b)))
		if err != nil {
			return nil, err
		}
		buffer := make([]byte, 32*1024)
		for _, f := range zr.File {
			if f.FileInfo().IsDir() {
				continue
			}
			err := func() error {
				r, err := f.Open()
				if err != nil {
					return err
				}
				defer r.Close()

				filePath := filepath.Join(bucketPath, f.Name)
				w, err := storageC.UploadToGCSPub(ctx, filePath, map[string]string{})
				if err != nil {
					return err
				}
				_, err = io.CopyBuffer(w, r, buffer)
				if err != nil {
					return err
				}
				return w.Close()
			}()
			if err != nil {
				return nil, err
			}
		}
		currentType := strings.ToLower(strings.TrimSpace(file.Type.String()))
		urlPath := bucketPath
		if currentType != "" {
			urlPath = urlPath + "/" + constants.StaticTypeMap[currentType]
		} else {
			return nil, fmt.Errorf("type is empty or not supported")
		}
		getUrl = storageC.GetSignedURLForObjectPub(urlPath)
		bucketPath = urlPath

	} else {
		getUrl = *file.URL
	}

	where := qb.Eq("id")
	updateQB := qb.Update("coursez.topic_content").Set("topiccontentbucket").Set("url").Where(where)
	updateQuery := updateQB.Query(*global.CassSession.Session).BindMap(qb.M{"id": file.ContentID, "topiccontentbucket": bucketPath, "url": getUrl})
	if err := updateQuery.ExecRelease(); err != nil {
		return &isSuccess, err
	}
	isSuccessRes := true
	isSuccess.Success = &isSuccessRes
	isSuccess.URL = &getUrl
	return &isSuccess, nil
}
