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

func ChapterCreate(ctx context.Context, courseID string, chapter *model.ChapterInput) (*model.Chapter, error) {
	log.Info("ChapterCreate called")
	guid := xid.New()
	cassandraChapter := coursez.Chapter{
		ID:          guid.String(),
		Name:        chapter.Name,
		Description: chapter.Description,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		CourseID:    courseID,

	}
	if chapter.ModuleID!=nil {
		cassandraChapter.ModuleID = *chapter.ModuleID
	}
	if chapter.Sequence!=nil {
		cassandraChapter.Sequence = *chapter.Sequence
	}
	// set course in cassandra
	insertQuery := global.CassSession.Session.Query(coursez.ChapterTable.Insert()).BindStruct(cassandraChapter)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraChapter.CreatedAt, 10)
	responseModel := model.Chapter{
		ID:          &cassandraChapter.ID,
		Name:        chapter.Name,
		Description: chapter.Description,
		CreatedAt:   &created,
		UpdatedAt:   &created,
		CourseID:    cassandraChapter.CourseID,
		ModuleID:    chapter.ModuleID,
		Sequence:    chapter.Sequence,
	}
	return &responseModel, nil
}
func UpdateChapter(ctx context.Context, chapter *model.ChapterInput) (*model.Chapter, error) {
	log.Info("ChapterUpdate called")
	if chapter.ID == nil {
		return nil, fmt.Errorf("chapter not found")
	}
	cassandraChapter := coursez.Chapter{
		ID: *chapter.ID,
	}
	// set course in cassandra
	getQuery := global.CassSession.Session.Query(coursez.ChapterTable.Get()).BindStruct(cassandraChapter)
	if err := getQuery.ExecRelease(); err != nil {
		return nil, err
	}
	if chapter.Name != "" {
		cassandraChapter.Name = chapter.Name
	}
	if chapter.Description != "" {
		cassandraChapter.Description = chapter.Description
	}
	if chapter.Sequence != nil {
		cassandraChapter.Sequence = *chapter.Sequence
	}
	if chapter.ModuleID != nil {
		cassandraChapter.ModuleID = *chapter.ModuleID
	}
	cassandraChapter.UpdatedAt = time.Now().Unix()
	updateQuery := global.CassSession.Session.Query(coursez.ChapterTable.Update()).BindStruct(cassandraChapter)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraChapter.CreatedAt, 10)
	responseModel := model.Chapter{
		ID:          &cassandraChapter.ID,
		Name:        cassandraChapter.Name,
		Description: cassandraChapter.Description,
		CreatedAt:   &created,
		UpdatedAt:   &created,
		CourseID:    cassandraChapter.CourseID,
		ModuleID:    &cassandraChapter.ModuleID,
		Sequence:    &cassandraChapter.Sequence,
	}
	return &responseModel, nil
}
