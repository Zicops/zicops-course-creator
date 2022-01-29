package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/graph/model"
)

func CreateTopicQuiz(ctx context.Context, quiz *model.QuizInput) (*model.Quiz, error) {

	log.Info("CreateTopicQuiz called")
	guid := xid.New()
	cassandraQuiz := coursez.Quiz{
		ID:          guid.String(),
		Name:        quiz.Name,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		Type:        quiz.Type,
		IsMandatory: quiz.IsMandatory,
		TopicID:     quiz.TopicID,
		StartTime:   *quiz.StartTime,
		Sequence:    *quiz.Sequence,
		Category:    quiz.Category,
		IsDeleted:   false,
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
	// get quiz from cassandra
	getQuery := global.CassSession.Session.Query(coursez.QuizTable.Get()).BindStruct(cassandraQuiz)
	if err := getQuery.ExecRelease(); err != nil {
		return nil, err
	}
	if quiz.Name != "" {
		cassandraQuiz.Name = quiz.Name
	}
	if quiz.Type != "" {
		cassandraQuiz.Type = quiz.Type
	}
	cassandraQuiz.IsMandatory = quiz.IsMandatory
	cassandraQuiz.TopicID = quiz.TopicID
	cassandraQuiz.StartTime = *quiz.StartTime
	cassandraQuiz.Sequence = *quiz.Sequence
	cassandraQuiz.Category = quiz.Category
	cassandraQuiz.UpdatedAt = time.Now().Unix()
	// update quiz in cassandra
	updateQuery := global.CassSession.Session.Query(coursez.QuizTable.Update()).BindStruct(cassandraQuiz)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuiz.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuiz.UpdatedAt, 10)
	responseModel := model.Quiz{
		ID:          &cassandraQuiz.ID,
		Name:        cassandraQuiz.Name,
		CreatedAt:   &created,
		UpdatedAt:   &updated,
		Type:        cassandraQuiz.Type,
		IsMandatory: cassandraQuiz.IsMandatory,
		TopicID:     cassandraQuiz.TopicID,
		StartTime:   &cassandraQuiz.StartTime,
		Sequence:    &cassandraQuiz.Sequence,
		Category:    cassandraQuiz.Category,
	}
	return &responseModel, nil

}
