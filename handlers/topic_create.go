package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/xid"
	"github.com/scylladb/gocqlx/v2/qb"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-cass-pool/cassandra"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/graph/model"
)

func TopicCreate(ctx context.Context, courseID string, topic *model.TopicInput) (*model.Topic, error) {
	log.Info("TopicCreate called")
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	global.CassSession = session

	guid := xid.New()
	cassandraTopic := coursez.Topic{
		ID:          guid.String(),
		Name:        *topic.Name,
		Description: *topic.Description,
		Type:        *topic.Type,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		CourseID:    courseID,
		IsActive:    true,
	}
	if topic.ChapterID != nil {
		cassandraTopic.ChapterID = *topic.ChapterID
	}
	if topic.ModuleID != nil {
		cassandraTopic.ModuleID = *topic.ModuleID
	}
	if topic.Sequence != nil {
		cassandraTopic.Sequence = *topic.Sequence
	}
	if topic.CreatedBy != nil {
		cassandraTopic.CreatedBy = *topic.CreatedBy
	}
	if topic.UpdatedBy != nil {
		cassandraTopic.UpdatedBy = *topic.UpdatedBy
	}
	// set course in cassandra
	insertQuery := global.CassSession.Query(coursez.TopicTable.Insert()).BindStruct(cassandraTopic)
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
		CourseID:    &cassandraTopic.CourseID,
		ModuleID:    topic.ModuleID,
		Sequence:    topic.Sequence,
		CreatedBy:   topic.CreatedBy,
		UpdatedBy:   topic.UpdatedBy,
	}
	return &responseModel, nil
}

func TopicUpdate(ctx context.Context, topic *model.TopicInput) (*model.Topic, error) {
	log.Info("TopicUpdate called")
	if topic.ID == nil {
		return nil, fmt.Errorf("course id is required")
	}
	cassandraTopic := coursez.Topic{
		ID: *topic.ID,
	}
	if topic.CourseID != nil {
		cassandraTopic.CourseID = *topic.CourseID
	}
	if topic.ModuleID != nil {
		cassandraTopic.ModuleID = *topic.ModuleID
	}
	topics := []coursez.Topic{}
	getQuery := global.CassSession.Query(coursez.TopicTable.Get()).BindMap(qb.M{"id": cassandraTopic.ID})
	if err := getQuery.SelectRelease(&topics); err != nil {
		return nil, err
	}
	cassandraTopic = topics[0]
	updateCols := []string{}
	if topic.Description != nil {
		updateCols = append(updateCols, "description")
		cassandraTopic.Description = *topic.Description
	}
	if topic.Sequence != nil {
		updateCols = append(updateCols, "sequence")
		cassandraTopic.Sequence = *topic.Sequence
	}
	if topic.CreatedBy != nil {
		updateCols = append(updateCols, "created_by")
		cassandraTopic.CreatedBy = *topic.CreatedBy
	}
	if topic.Name != nil {
		updateCols = append(updateCols, "name")
		cassandraTopic.Name = *topic.Name
	}
	if topic.Type != nil {
		updateCols = append(updateCols, "type")
		cassandraTopic.Type = *topic.Type
	}
	updateCols = append(updateCols, "updated_at")
	cassandraTopic.UpdatedAt = time.Now().Unix()
	// set course in cassandra
	upStms, uNames := coursez.TopicTable.Update(updateCols...)
	updateQuery := global.CassSession.Query(upStms, uNames).BindStruct(&cassandraTopic)
	if err := updateQuery.ExecRelease(); err != nil {
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
		CourseID:    &cassandraTopic.CourseID,
		ModuleID:    topic.ModuleID,
		Sequence:    topic.Sequence,
		CreatedBy:   topic.CreatedBy,
	}
	return &responseModel, nil
}
