package handlers

import (
	"bytes"
	"context"
	"io"
	"strconv"
	"time"

	"github.com/scylladb/gocqlx/v2/qb"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/lib/db/bucket"
	"github.com/zicops/zicops-course-creator/lib/googleprojectlib"
)

func TopicContentCreate(ctx context.Context, topicID string, topicConent *model.TopicContentInput) (*model.TopicContent, error) {
	log.Info("TopicContentCreate called")
	cassandraTopicContent := coursez.TopicContent{
		TopicId:            topicID,
		Language:           topicConent.Language,
		CreatedAt:          time.Now().Unix(),
		UpdatedAt:          time.Now().Unix(),
		StartTime:          *topicConent.StartTime,
		Duration:           *topicConent.Duration,
		SkipIntro:          *topicConent.SkipIntro,
		NextShowtime:       *topicConent.NextShowTime,
		FromEndTime:        *topicConent.FromEndTime,
		TopicContentBucket: "",
		Url:                "",
		IsDeleted:          false,
	}
	// set course in cassandra
	insertQuery := global.CassSession.Session.Query(coursez.TopicContentTable.Insert()).BindStruct(cassandraTopicContent)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraTopicContent.CreatedAt, 10)
	responseModel := model.TopicContent{
		Language:     topicConent.Language,
		StartTime:    topicConent.StartTime,
		CreatedAt:    &created,
		UpdatedAt:    &created,
		Duration:     topicConent.Duration,
		SkipIntro:    topicConent.SkipIntro,
		NextShowTime: topicConent.NextShowTime,
		FromEndTime:  topicConent.FromEndTime,
		TopicID:      topicID,
	}
	return &responseModel, nil
}

func UploadTopicVideo(ctx context.Context, file model.TopicVideo) (*bool, error) {
	log.Info("UploadTopicVideo called")
	isSuccess := false
	storageC := bucket.NewStorageHandler()
	gproject := googleprojectlib.GetGoogleProjectID()
	err := storageC.InitializeStorageClient(ctx, gproject)
	if err != nil {
		log.Errorf("Failed to upload video to course topic: %v", err.Error())
		return &isSuccess, nil
	}
	bucketPath := file.CourseID + "/" + file.TopicID + "/" + file.File.Filename
	writer, err := storageC.UploadToGCS(ctx, bucketPath)
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
	// update course image in cassandra
	updateQuery := global.CassSession.Session.Query(coursez.TopicContentTable.Update("topicContentBucket", "url")).BindMap(qb.M{"topicId": file.TopicID, "topicContentBucket": bucketPath, "url": getUrl})
	if err := updateQuery.ExecRelease(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}