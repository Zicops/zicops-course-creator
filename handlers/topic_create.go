package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/xid"
	"github.com/scylladb/gocqlx/v2"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-cass-pool/cassandra"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/helpers"
)

func TopicCreate(ctx context.Context, courseID string, topic *model.TopicInput) (*model.Topic, error) {
	log.Info("TopicCreate called")
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	email_creator := claims["email"].(string)
	lspId := claims["lsp_id"].(string)
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
		LspId:       lspId,
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

	cassandraTopic.CreatedBy = email_creator

	cassandraTopic.UpdatedBy = email_creator

	// set course in cassandra
	insertQuery := CassSession.Query(coursez.TopicTable.Insert()).BindStruct(cassandraTopic)
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
	CassSession, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	email_creator := claims["email"].(string)
	lspID := claims["lsp_id"].(string)
	cassandraTopic := coursez.Topic{
		ID: *topic.ID,
	}
	if topic.CourseID != nil {
		cassandraTopic.CourseID = *topic.CourseID
	}
	if topic.ModuleID != nil {
		cassandraTopic.ModuleID = *topic.ModuleID
	}
	cassandraTopic = *GetTopic(ctx, cassandraTopic.ID, lspID, CassSession)
	updateCols := []string{}
	if topic.Description != nil && cassandraTopic.Description != *topic.Description {
		updateCols = append(updateCols, "description")
		cassandraTopic.Description = *topic.Description
	}
	if topic.Sequence != nil && cassandraTopic.Sequence != *topic.Sequence {
		updateCols = append(updateCols, "sequence")
		cassandraTopic.Sequence = *topic.Sequence
	}
	if topic.Name != nil && cassandraTopic.Name != *topic.Name {
		updateCols = append(updateCols, "name")
		cassandraTopic.Name = *topic.Name
	}
	if topic.Type != nil && cassandraTopic.Type != *topic.Type {
		updateCols = append(updateCols, "type")
		cassandraTopic.Type = *topic.Type
	}
	if email_creator != "" && cassandraTopic.UpdatedBy != email_creator {
		updateCols = append(updateCols, "updated_by")
		cassandraTopic.UpdatedBy = email_creator
	}
	if len(updateCols) > 0 {
		updateCols = append(updateCols, "updated_at")
		cassandraTopic.UpdatedAt = time.Now().Unix()
		// set course in cassandra
		upStms, uNames := coursez.TopicTable.Update(updateCols...)
		updateQuery := CassSession.Query(upStms, uNames).BindStruct(&cassandraTopic)
		if err := updateQuery.ExecRelease(); err != nil {
			return nil, err
		}
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

func GetTopic(ctx context.Context, courseID string, lspID string, session *gocqlx.Session) *coursez.Topic {
	chapters := []coursez.Topic{}
	getQueryStr := fmt.Sprintf("SELECT * FROM coursez.topic WHERE id='%s' and lsp_id='%s' and is_active=true", courseID, lspID)
	getQuery := session.Query(getQueryStr, nil)
	if err := getQuery.SelectRelease(&chapters); err != nil {
		return nil
	}
	return &chapters[0]
}
