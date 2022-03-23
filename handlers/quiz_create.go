package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/rs/xid"
	"github.com/scylladb/gocqlx/v2/qb"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/lib/db/bucket"
	"github.com/zicops/zicops-course-creator/lib/googleprojectlib"
)

func CreateTopicQuiz(ctx context.Context, quiz *model.QuizInput) (*model.Quiz, error) {

	log.Info("CreateTopicQuiz called")
	guid := xid.New()
	cassandraQuiz := coursez.Quiz{
		ID:        guid.String(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		IsActive:  true,
	}
	if quiz.Name != nil {
		cassandraQuiz.Name = *quiz.Name
	}
	if quiz.Type != nil {
		cassandraQuiz.Type = *quiz.Type
	}
	if quiz.IsMandatory != nil {
		cassandraQuiz.IsMandatory = *quiz.IsMandatory
	}
	if quiz.TopicID != nil {
		cassandraQuiz.TopicID = *quiz.TopicID
	}
	if quiz.Category != nil {
		cassandraQuiz.Category = *quiz.Category
	}
	if quiz.StartTime != nil {
		cassandraQuiz.StartTime = *quiz.StartTime
	}
	if quiz.Sequence != nil {
		cassandraQuiz.Sequence = *quiz.Sequence
	}
	// set quiz in cassandra
	insertQuery := global.CassSession.Session.Query(coursez.QuizTable.Insert()).BindStruct(cassandraQuiz)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuiz.CreatedAt, 10)
	responseModel := model.Quiz{
		ID:          &cassandraQuiz.ID,
		Name:        quiz.Name,
		CreatedAt:   &created,
		UpdatedAt:   &created,
		Type:        quiz.Type,
		IsMandatory: quiz.IsMandatory,
		TopicID:     quiz.TopicID,
		StartTime:   quiz.StartTime,
		Sequence:    quiz.Sequence,
		Category:    quiz.Category,
	}
	return &responseModel, nil
}

func UpdateQuiz(ctx context.Context, quiz *model.QuizInput) (*model.Quiz, error) {
	log.Info("UpdateQuiz called")
	if quiz.ID == nil {
		return nil, fmt.Errorf("quiz id is required")
	}
	cassandraQuiz := coursez.Quiz{
		ID: *quiz.ID,
	}
	// set course in cassandra
	quizes := []coursez.Quiz{}
	getQuery := global.CassSession.Session.Query(coursez.QuizTable.Get()).BindMap(qb.M{"id": cassandraQuiz.ID})
	if err := getQuery.SelectRelease(&quizes); err != nil {
		return nil, err
	}
	if len(quizes) < 1 {
		return nil, fmt.Errorf("quiz not found")
	}
	cassandraQuiz = quizes[0]
	updateCols := []string{}
	if quiz.Name != nil {
		updateCols = append(updateCols, "name")
		cassandraQuiz.Name = *quiz.Name
	}
	if quiz.Type != nil {
		updateCols = append(updateCols, "type")
		cassandraQuiz.Type = *quiz.Type
	}
	if quiz.IsMandatory != nil {
		updateCols = append(updateCols, "ismandatory")
		cassandraQuiz.IsMandatory = *quiz.IsMandatory
	}
	if quiz.TopicID != nil {
		updateCols = append(updateCols, "topicid")
		cassandraQuiz.TopicID = *quiz.TopicID
	}
	if quiz.StartTime != nil {
		updateCols = append(updateCols, "starttime")
		cassandraQuiz.StartTime = *quiz.StartTime
	}
	if quiz.Sequence != nil {
		updateCols = append(updateCols, "sequence")
		cassandraQuiz.Sequence = *quiz.Sequence
	}
	if quiz.Category != nil {
		updateCols = append(updateCols, "category")
		cassandraQuiz.Category = *quiz.Category
	}
	updateCols = append(updateCols, "updated_at")
	cassandraQuiz.UpdatedAt = time.Now().Unix()
	// update quiz in cassandra
	upStms, uNames := coursez.QuizTable.Update(updateCols...)
	updateQuery := global.CassSession.Session.Query(upStms, uNames).BindStruct(&cassandraQuiz)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuiz.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuiz.UpdatedAt, 10)
	responseModel := model.Quiz{
		ID:          &cassandraQuiz.ID,
		Name:        &cassandraQuiz.Name,
		CreatedAt:   &created,
		UpdatedAt:   &updated,
		Type:        &cassandraQuiz.Type,
		IsMandatory: &cassandraQuiz.IsMandatory,
		TopicID:     &cassandraQuiz.TopicID,
		StartTime:   &cassandraQuiz.StartTime,
		Sequence:    &cassandraQuiz.Sequence,
		Category:    &cassandraQuiz.Category,
	}
	return &responseModel, nil

}

func UploadQuizFile(ctx context.Context, courseID string, quiz model.QuizFile) (*model.UploadResult, error) {
	log.Info("UploadQuizFile called")
	isSuccess := model.UploadResult{}
	storageC := bucket.NewStorageHandler()
	gproject := googleprojectlib.GetGoogleProjectID()
	err := storageC.InitializeStorageClient(ctx, gproject)
	if err != nil {
		log.Errorf("Failed to upload video to course topic: %v", err.Error())
		return &isSuccess, nil
	}
	if courseID == "" || quiz.QuizID== nil {
		return nil, fmt.Errorf("course id and  quiz id is required")
	}
	bucketPath := courseID + "/" + *quiz.QuizID + "/" + quiz.File.Filename
	writer, err := storageC.UploadToGCS(ctx, bucketPath)
	if err != nil {
		log.Errorf("Failed to upload video to course topic: %v", err.Error())
		return &isSuccess, nil
	}
	defer writer.Close()
	fileBuffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(fileBuffer, quiz.File.File); err != nil {
		return &isSuccess, nil
	}
	currentBytes := fileBuffer.Bytes()
	_, err = io.Copy(writer, bytes.NewReader(currentBytes))
	if err != nil {
		return &isSuccess, err
	}
	getUrl := storageC.GetSignedURLForObject(bucketPath)
	cassandraQuizFile := coursez.QuizFile{
		QuizId:     *quiz.QuizID,
		BucketPath: bucketPath,
		Path:       getUrl,
		IsActive:   true,
	}
	if quiz.Type != nil {
		cassandraQuizFile.Type = *quiz.Type
	}
	if quiz.Name != nil {
		cassandraQuizFile.Name = *quiz.Name
	}
	// update course image in cassandra
	quizAdd := global.CassSession.Session.Query(coursez.QuizFileTable.Insert()).BindStruct(cassandraQuizFile)
	if err := quizAdd.ExecRelease(); err != nil {
		return &isSuccess, err
	}
	isSuccessRes := true
	isSuccess.Success = &isSuccessRes
	isSuccess.URL = &getUrl
	return &isSuccess, nil
}

func AddMCQQuiz(ctx context.Context, quiz *model.QuizMcq) (*bool, error) {
	log.Info("AddMCQQuiz called")
	if *quiz.QuizID == "" {
		return nil, fmt.Errorf("quiz id is required")
	}
	options := make([]string, 0)
	for _, option := range quiz.Options {
		options = append(options, *option)
	}
	cassandraQuiz := coursez.QuizMcq{
		QuizId:        *quiz.QuizID,
		Question:      *quiz.Question,
		Options:       options,
		CorrectOption: *quiz.CorrectOption,
		Explanation:   *quiz.Explanation,
		IsActive:      true,
	}
	// set quiz in cassandra
	insertQuery := global.CassSession.Session.Query(coursez.QuizMcqTable.Insert()).BindStruct(cassandraQuiz)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	isSuccess := true
	return &isSuccess, nil

}

func AddQuizDescriptive(ctx context.Context, quiz *model.QuizDescriptive) (*bool, error) {
	log.Info("AddQuizDescriptive called")
	if *quiz.QuizID == "" {
		return nil, fmt.Errorf("quiz id is required")
	}
	cassandraQuiz := coursez.QuizDescriptive{
		QuizId:        *quiz.QuizID,
		Question:      *quiz.Question,
		Explanation:   *quiz.Explanation,
		CorrectAnswer: *quiz.CorrectAnswer,
		IsActive:      true,
	}
	// set quiz in cassandra
	insertQuery := global.CassSession.Session.Query(coursez.QuizDescriptiveTable.Insert()).BindStruct(cassandraQuiz)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	isSuccess := true
	return &isSuccess, nil
}
