package handlers

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/scylladb/gocqlx/v2"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-cass-pool/cassandra"
	"github.com/zicops/zicops-course-creator/constants"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/helpers"
	"github.com/zicops/zicops-course-creator/lib/db/bucket"
	"github.com/zicops/zicops-course-creator/lib/googleprojectlib"
	deploy_static "github.com/zicops/zicops-course-creator/lib/static"
	"github.com/zicops/zicops-course-creator/lib/utils"
)

func TopicContentCreate(ctx context.Context, topicID string, courseID string, moduleID *string, topicConent *model.TopicContentInput) (*model.TopicContent, error) {
	log.Info("TopicContentCreate called")

	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspID := claims["lsp_id"].(string)
	if moduleID == nil {
		return nil, fmt.Errorf("moduleID is nil")
	}
	cassandraTopicContent := coursez.TopicContent{
		ID:                 uuid.New().String(),
		Language:           *topicConent.Language,
		TopicId:            topicID,
		CourseId:           courseID,
		ModuleId:           *moduleID,
		LspId:              lspID,
		CreatedAt:          time.Now().Unix(),
		UpdatedAt:          time.Now().Unix(),
		StartTime:          0,
		Duration:           0,
		SkipIntroDuration:  0,
		NextShowtime:       0,
		FromEndTime:        0,
		TopicContentBucket: "",
		Url:                "",
		SubtitleFile:       "",
		SubtitleFileBucket: "",
		Type:               "",
		IsActive:           true,
		IsDefault:          false,
	}
	contentCountByTopicID := GetTopicContentCountByTopicID(ctx, topicID, lspID, CassSession)
	if moduleID != nil && topicConent.Duration != nil && contentCountByTopicID == 0 {
		mod := GetModule(ctx, *moduleID, lspID, CassSession)
		if mod == nil {
			return nil, fmt.Errorf("module not found")
		}
		newDuration := *topicConent.Duration + mod.Duration
		queryStr := fmt.Sprintf("UPDATE coursez.module SET duration=%d WHERE id='%s' and lsp_id='%s' and is_active=true AND created_at=%d", newDuration, *moduleID, lspID, mod.CreatedAt)
		updateQ := CassSession.Query(queryStr, nil)
		if err := updateQ.ExecRelease(); err != nil {
			return nil, err
		}
	}
	if topicConent.Duration != nil && cassandraTopicContent.CourseId != "" && contentCountByTopicID == 0 {
		course := GetCourse(ctx, cassandraTopicContent.CourseId, lspID, CassSession)
		if course == nil {
			return nil, fmt.Errorf("course not found")
		}
		newDuration := course.Duration - cassandraTopicContent.Duration + *topicConent.Duration
		queryStr := fmt.Sprintf("UPDATE coursez.course SET duration=%d WHERE id='%s' and lsp_id='%s' and is_active=true and created_at=%d", newDuration, cassandraTopicContent.CourseId, lspID, course.CreatedAt)
		updateQ := CassSession.Query(queryStr, nil)
		if err := updateQ.ExecRelease(); err != nil {
			return nil, err
		}
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
	insertQuery := CassSession.Query(coursez.TopicContentTable.Insert()).BindStruct(cassandraTopicContent)
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
		CourseID:          &courseID,
	}
	return &responseModel, nil
}

func GetTopicContentCountByTopicID(ctx context.Context, topicID string, lspID string, CassSession *gocqlx.Session) int {
	queryStr := fmt.Sprintf("SELECT COUNT(*) FROM coursez.topic_content WHERE topicid='%s' and lsp_id='%s' and is_active=true ALLOW FILTERING", topicID, lspID)
	iter := CassSession.Query(queryStr, nil).Iter()
	var count int
	for iter.Scan(&count) {
	}
	return count
}

func TopicExamCreate(ctx context.Context, topicID string, courseID string, exam *model.TopicExamInput) (*model.TopicExam, error) {
	log.Info("TopicExamCreate called")

	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session
	sessionQbankz, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	CassSessionQBank := sessionQbankz
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspID := claims["lsp_id"].(string)
	cassExam := GetExam(ctx, *exam.ExamID, lspID, CassSessionQBank)
	course := GetCourse(ctx, courseID, lspID, CassSession)
	newDuration := course.Duration + cassExam.Duration
	queryStr := fmt.Sprintf("UPDATE coursez.course SET duration=%d WHERE id='%s' and lsp_id='%s' and is_active=true and created_at=%d", newDuration, courseID, lspID, course.CreatedAt)
	updateQ := CassSession.Query(queryStr, nil)
	if err := updateQ.ExecRelease(); err != nil {
		return nil, err
	}
	cassandraTopicContent := coursez.TopicExam{
		ID:        uuid.New().String(),
		TopicId:   topicID,
		CourseId:  courseID,
		ExamId:    *exam.ExamID,
		Language:  *exam.Language,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		IsActive:  true,
		LspId:     lspID,
	}
	// set course in cassandra
	insertQuery := CassSession.Query(coursez.TopicExamTable.Insert()).BindStruct(cassandraTopicContent)
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
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	if lspId == "" {
		return nil, fmt.Errorf("lsp_id is empty")
	}
	isSuccess := model.UploadResult{}
	if file.CourseID == nil || file.ContentID == nil {
		log.Errorf("Failed to upload video to course topic: %v", "courseID or contentId is nil")
		return &isSuccess, nil
	}
	bucketPath := *file.CourseID + "/" + *file.ContentID + "/" + base64.URLEncoding.EncodeToString([]byte(file.File.Filename))
	storageC := bucket.NewStorageHandler()
	gproject := googleprojectlib.GetGoogleProjectID()
	err = storageC.InitializeStorageClient(ctx, gproject, lspId)
	if err != nil {
		log.Errorf("Failed to upload video to course topic: %v", err.Error())
	}
	sendUploadRequestToUploaderQueue(ctx, *file.File, bucketPath, lspId)
	getUrl := storageC.GetSignedURLForObject(bucketPath)
	topicContent := GetTopicContent(ctx, *file.ContentID, lspId, CassSession)
	updateQuery := fmt.Sprintf("UPDATE coursez.topic_content SET topiccontentbucket='%s', url='%s' WHERE id='%s' AND lsp_id='%s' AND is_active=true and created_at=%d", bucketPath, getUrl, topicContent.ID, topicContent.LspId, topicContent.CreatedAt)
	updateQ := CassSession.Query(updateQuery, nil)
	if err := updateQ.ExecRelease(); err != nil {
		return nil, err
	}
	isSuccessRes := true
	isSuccess.Success = &isSuccessRes
	isSuccess.URL = &getUrl
	return &isSuccess, nil
}

func sendUploadRequestToUploaderQueue(ctx context.Context, file graphql.Upload, bucketPath string, lspId string) {
	// send message to uploader queue
	uploadRequest := utils.UploadRequest{
		BucketPath: bucketPath,
		File:       &file,
		LspId:      lspId,
	}

	// send message to uploader queue in utils
	utils.UploaderQueue <- &uploadRequest
}

func UploadTopicSubtitle(ctx context.Context, files []*model.TopicSubtitle) ([]*model.UploadResultSubtitles, error) {
	log.Info("UploadTopicSubtitle called")
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	isSuccess := []*model.UploadResultSubtitles{}
	lspId := claims["lsp_id"].(string)
	if lspId == "" {
		return nil, fmt.Errorf("lsp_id is empty")
	}
	for _, file := range files {
		isLocalSuccess := model.UploadResultSubtitles{}
		isLocal := false
		isLocalSuccess.Success = &isLocal
		isLocalSuccess.URL = nil
		storageC := bucket.NewStorageHandler()
		gproject := googleprojectlib.GetGoogleProjectID()
		err := storageC.InitializeStorageClient(ctx, gproject, lspId)
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
		bucketPath := mainBucket + base64.URLEncoding.EncodeToString([]byte(file.File.Filename))
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

func UpdateTopicContent(ctx context.Context, topicConent *model.TopicContentInput, moduleId *string) (*model.TopicContent, error) {
	log.Info("UpdateTopicContent called")
	contentID := topicConent.ContentID
	if *contentID == "" {
		return nil, fmt.Errorf("ContentID is required")
	}
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspID := claims["lsp_id"].(string)
	cassandraTopicContent := *GetTopicContent(ctx, *contentID, lspID, CassSession)
	updateCols := []string{}
	if topicConent.Duration != nil && *topicConent.Duration != cassandraTopicContent.Duration {
		updateCols = append(updateCols, "duration")
		cassandraTopicContent.Duration = *topicConent.Duration
		if moduleId != nil && topicConent.Duration != nil {
			mod := GetModule(ctx, *moduleId, lspID, CassSession)
			newDuration := mod.Duration - cassandraTopicContent.Duration + *topicConent.Duration
			queryStr := fmt.Sprintf("UPDATE coursez.module SET duration=%d WHERE id='%s'and lsp_id='%s' and is_active=true and created_at=%d", newDuration, *moduleId, lspID, mod.CreatedAt)
			updateQ := CassSession.Query(queryStr, nil)
			if err := updateQ.ExecRelease(); err != nil {
				return nil, err
			}
		}
		if cassandraTopicContent.CourseId != "" && topicConent.Duration != nil {
			course := GetCourse(ctx, cassandraTopicContent.CourseId, lspID, CassSession)
			newDuration := course.Duration - cassandraTopicContent.Duration + *topicConent.Duration
			queryStr := fmt.Sprintf("UPDATE coursez.course SET duration=%d WHERE id='%s' and lsp_id='%s' and is_active=true and created_at=%d", newDuration, cassandraTopicContent.CourseId, lspID, course.CreatedAt)
			updateQ := CassSession.Query(queryStr, nil)
			if err := updateQ.ExecRelease(); err != nil {
				return nil, err
			}
		}
	}
	if topicConent.StartTime != nil && *topicConent.StartTime != cassandraTopicContent.StartTime {
		updateCols = append(updateCols, "starttime")
		cassandraTopicContent.StartTime = *topicConent.StartTime
	}
	if topicConent.SkipIntroDuration != nil && *topicConent.SkipIntroDuration != cassandraTopicContent.SkipIntroDuration {
		updateCols = append(updateCols, "skipintroduration")
		cassandraTopicContent.SkipIntroDuration = *topicConent.SkipIntroDuration
	}
	if topicConent.Type != nil && *topicConent.Type != cassandraTopicContent.Type {
		updateCols = append(updateCols, "type")
		cassandraTopicContent.Type = *topicConent.Type
	}
	if topicConent.NextShowTime != nil && *topicConent.NextShowTime != cassandraTopicContent.NextShowtime {
		updateCols = append(updateCols, "nextshowtime")
		cassandraTopicContent.NextShowtime = *topicConent.NextShowTime
	}
	if topicConent.FromEndTime != nil && *topicConent.FromEndTime != cassandraTopicContent.FromEndTime {
		updateCols = append(updateCols, "fromendtime")
		cassandraTopicContent.FromEndTime = *topicConent.FromEndTime
	}
	if topicConent.Language != nil {
		updateCols = append(updateCols, "language")
		cassandraTopicContent.Language = *topicConent.Language
	}
	if topicConent.IsDefault != nil {
		updateCols = append(updateCols, "is_default")
		cassandraTopicContent.IsDefault = *topicConent.IsDefault
	}
	if len(updateCols) > 0 {
		updateCols = append(updateCols, "updated_at")
		cassandraTopicContent.UpdatedAt = time.Now().Unix()
		upStms, uNames := coursez.TopicContentTable.Update(updateCols...)
		updateQuery := CassSession.Query(upStms, uNames).BindStruct(&cassandraTopicContent)
		if err := updateQuery.ExecRelease(); err != nil {
			return nil, err
		}
	}
	created := strconv.FormatInt(cassandraTopicContent.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraTopicContent.UpdatedAt, 10)
	responseModel := model.TopicContent{
		ID:                &cassandraTopicContent.ID,
		Language:          topicConent.Language,
		StartTime:         topicConent.StartTime,
		CreatedAt:         &created,
		UpdatedAt:         &updated,
		Duration:          topicConent.Duration,
		SkipIntroDuration: topicConent.SkipIntroDuration,
		NextShowTime:      topicConent.NextShowTime,
		FromEndTime:       topicConent.FromEndTime,
		TopicID:           &cassandraTopicContent.TopicId,
		Type:              topicConent.Type,
		IsDefault:         topicConent.IsDefault,
		CourseID:          &cassandraTopicContent.CourseId,
	}
	return &responseModel, nil
}

func UpdateTopicExam(ctx context.Context, exam *model.TopicExamInput) (*model.TopicExam, error) {
	log.Info("UpdateTopicExam called")
	tExamId := exam.ID
	if *tExamId == "" {
		return nil, fmt.Errorf("TopicExamId is required")
	}
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspID := claims["lsp_id"].(string)
	cassandraTopicContent := *GetTopicExam(ctx, *tExamId, lspID, CassSession)
	updateCols := []string{}
	if exam.Language != nil && *exam.Language != cassandraTopicContent.Language {
		updateCols = append(updateCols, "language")
		cassandraTopicContent.Language = *exam.Language
	}
	if exam.TopicID != nil && *exam.TopicID != cassandraTopicContent.TopicId {
		updateCols = append(updateCols, "topicid")
		cassandraTopicContent.TopicId = *exam.TopicID
	}
	if exam.ExamID != nil && *exam.ExamID != cassandraTopicContent.ExamId {
		// get exam
		sessionQbankz, err := cassandra.GetCassSession("qbankz")
		if err != nil {
			return nil, err
		}
		CassSessionQBank := sessionQbankz
		cassExam := GetExam(ctx, cassandraTopicContent.ExamId, lspID, CassSessionQBank)
		cassExamNew := GetExam(ctx, *exam.ExamID, lspID, CassSessionQBank)
		// update course duration
		cassCourse := GetCourse(ctx, cassandraTopicContent.CourseId, lspID, CassSession)
		newDuration := cassCourse.Duration - cassExam.Duration + cassExamNew.Duration
		queryStr := fmt.Sprintf("UPDATE coursez.course SET duration=%d WHERE id='%s' and lsp_id='%s' and is_active=true and created_at=%d", newDuration, cassandraTopicContent.CourseId, lspID, cassCourse.CreatedAt)
		updateQ := CassSession.Query(queryStr, nil)
		if err := updateQ.ExecRelease(); err != nil {
			return nil, err
		}
		updateCols = append(updateCols, "examid")
		cassandraTopicContent.ExamId = *exam.ExamID
	}
	if len(updateCols) > 0 {
		updateCols = append(updateCols, "updated_at")
		cassandraTopicContent.UpdatedAt = time.Now().Unix()
		upStms, uNames := coursez.TopicExamTable.Update(updateCols...)
		updateQuery := CassSession.Query(upStms, uNames).BindStruct(&cassandraTopicContent)
		if err := updateQuery.ExecRelease(); err != nil {
			return nil, err
		}
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
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	if lspId == "" {
		return nil, fmt.Errorf("lsp_id is empty")
	}
	isSuccess := model.UploadResult{}
	bucketPath := ""
	getUrl := ""
	if (file.URL == nil || *file.URL == "") && file.File != nil {
		// check if file.File.FileName is .zip.gz
		if !strings.Contains(file.File.Filename, ".zip") {
			return nil, fmt.Errorf("file should be .zip")
		}
		storageC := bucket.NewStorageHandler()
		gproject := googleprojectlib.GetGoogleProjectID()
		err := storageC.InitializeStorageClient(ctx, gproject, "static-content-private")
		if err != nil {
			log.Errorf("Failed to upload static content to course topic: %v", err.Error())
			return &isSuccess, nil
		}
		err = deploy_static.Initialize(storageC)
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
		bucketPath = lspId + "/" + *file.CourseID + "/" + *file.ContentID + "/" + hashString + "/"
		w, err := storageC.UploadToGCS(ctx, bucketPath+base64.URLEncoding.EncodeToString([]byte(file.File.Filename)), map[string]string{})
		if err != nil {
			log.Errorf("Failed to upload static content to course topic: %v", err.Error())
			return &isSuccess, nil
		}
		// write file to bucket
		_, err = io.Copy(w, file.File.File)
		if err != nil {
			log.Errorf("Failed to upload static content to course topic: %v", err.Error())
			return &isSuccess, nil
		}
		err = w.Close()
		if err != nil {
			log.Errorf("Failed to upload static content to course topic: %v", err.Error())
			return &isSuccess, nil
		}
		currentType := strings.ToLower(strings.TrimSpace(file.Type.String()))
		dataStorage := deploy_static.StorageObjectData{}
		dataStorage.Bucket = "static-content-private"
		dataStorage.Name = bucketPath + base64.URLEncoding.EncodeToString([]byte(file.File.Filename))
		dataStorage.ContentID = *file.ContentID
		outputAmp, err := deploy_static.DeployStatic(ctx, dataStorage)
		if err != nil {
			log.Errorf("Failed to upload static content to course topic: %v", err.Error())
			return &isSuccess, nil
		}
		getUrl = outputAmp
		if currentType != "" {
			getUrl = getUrl + constants.StaticTypeMap[currentType]
		} else {
			return nil, fmt.Errorf("type is empty or not supported")
		}

	} else {
		getUrl = *file.URL
	}
	topicContent := GetTopicContent(ctx, *file.ContentID, lspId, CassSession)
	updateQuery := fmt.Sprintf("UPDATE coursez.topic_content SET topiccontentbucket='%s', url='%s' WHERE id='%s' AND lsp_id='%s' AND is_active=true and created_at=%d", bucketPath, getUrl, *file.ContentID, lspId, topicContent.CreatedAt)
	updateQ := CassSession.Query(updateQuery, nil)
	if err := updateQ.ExecRelease(); err != nil {
		return nil, err
	}
	isSuccessRes := true
	isSuccess.Success = &isSuccessRes
	isSuccess.URL = &getUrl
	return &isSuccess, nil
}

func UploadStaticZipHandler(c *gin.Context) (*model.UploadResult, error) {
	// unmarshal model.StaticContent
	var file model.StaticContent
	if err := c.ShouldBind(&file); err != nil {
		log.Errorf("Failed to bind request: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}
	return UploadTopicStaticContent(c, &file)
}
func GetTopicContent(ctx context.Context, courseID string, lspID string, session *gocqlx.Session) *coursez.TopicContent {
	chapters := []coursez.TopicContent{}
	getQueryStr := fmt.Sprintf("SELECT * FROM coursez.topic_content WHERE id='%s' and lsp_id='%s' and is_active=true", courseID, lspID)
	getQuery := session.Query(getQueryStr, nil)
	if err := getQuery.SelectRelease(&chapters); err != nil {
		return nil
	}
	return &chapters[0]
}

func GetTopicExam(ctx context.Context, courseID string, lspID string, session *gocqlx.Session) *coursez.TopicExam {
	chapters := []coursez.TopicExam{}
	getQueryStr := fmt.Sprintf("SELECT * FROM coursez.topic_exam WHERE id='%s' and lsp_id='%s' and is_active=true", courseID, lspID)
	getQuery := session.Query(getQueryStr, nil)
	if err := getQuery.SelectRelease(&chapters); err != nil {
		return nil
	}
	return &chapters[0]
}
