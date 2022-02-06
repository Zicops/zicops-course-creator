package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/graph/model"
)

func TopicCreate(ctx context.Context, courseID string, topic *model.TopicInput) (*model.Topic, error) {
	log.Info("TopicCreate called")
	guid := xid.New()
	cassandraTopic := coursez.Topic{
		ID:          guid.String(),
		Name:        topic.Name,
		Description: topic.Description,
		Type:        topic.Type,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		ChapterID:   *topic.ChapterID,
		CourseID:    courseID,
		ModuleID:    *topic.ModuleID,
		Sequence:    *topic.Sequence,
		CreatedBy:   *topic.CreatedBy,
		UpdatedBy:   *topic.UpdatedBy,
		IsActive:   true,
	}
	// set course in cassandra
	insertQuery := global.CassSession.Session.Query(coursez.TopicTable.Insert()).BindStruct(cassandraTopic)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraTopic.CreatedAt, 10)
	responseModel := model.Topic{
		ID:          &cassandraTopic.ID,
		Name:        topic.Name,
		Description: topic.Description,
		Type:        topic.Type,
		CreatedAt:   &created,
		UpdatedAt:   &created,
		ChapterID:   topic.ChapterID,
		CourseID:    cassandraTopic.CourseID,
		ModuleID:    topic.ModuleID,
		Sequence:    topic.Sequence,
		CreatedBy:   topic.CreatedBy,
	}
	return &responseModel, nil
}
