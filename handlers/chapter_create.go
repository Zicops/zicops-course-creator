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

func ChapterCreate(ctx context.Context, courseID string, chapter *model.ChapterInput) (*model.Chapter, error) {
	log.Info("ChapterCreate called")
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	global.CassSession = session

	guid := xid.New()
	cassandraChapter := coursez.Chapter{
		ID:          guid.String(),
		Name:        *chapter.Name,
		Description: *chapter.Description,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		CourseID:    courseID,
	}
	if chapter.ModuleID != nil {
		cassandraChapter.ModuleID = *chapter.ModuleID
	}
	if chapter.Sequence != nil {
		cassandraChapter.Sequence = *chapter.Sequence
	}
	// set course in cassandra
	insertQuery := global.CassSession.Query(coursez.ChapterTable.Insert()).BindStruct(cassandraChapter)
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
		CourseID:    &cassandraChapter.CourseID,
		ModuleID:    chapter.ModuleID,
		Sequence:    chapter.Sequence,
	}
	return &responseModel, nil
}

func UpdateChapter(ctx context.Context, chapter *model.ChapterInput) (*model.Chapter, error) {
	log.Info("ChapterUpdate called")
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	global.CassSession = session

	if chapter.ID == nil {
		return nil, fmt.Errorf("chapter not found")
	}
	cassandraChapter := coursez.Chapter{
		ID: *chapter.ID,
	}
	chapters := []coursez.Chapter{}
	getQuery := global.CassSession.Query(coursez.ChapterTable.Get()).BindMap(qb.M{"id": cassandraChapter.ID})
	if err := getQuery.SelectRelease(&chapters); err != nil {
		return nil, err
	}
	if len(chapters) < 1 {
		return nil, fmt.Errorf("chapter not found")
	}
	cassandraChapter = chapters[0]
	updateCols := make([]string, 0)
	if *chapter.Name != "" {
		updateCols = append(updateCols, "name")
		cassandraChapter.Name = *chapter.Name
	}
	if *chapter.Description != "" {
		updateCols = append(updateCols, "description")
		cassandraChapter.Description = *chapter.Description
	}
	if chapter.Sequence != nil {
		updateCols = append(updateCols, "sequence")
		cassandraChapter.Sequence = *chapter.Sequence
	}
	if chapter.ModuleID != nil {
		updateCols = append(updateCols, "moduleid")
		cassandraChapter.ModuleID = *chapter.ModuleID
	}
	cassandraChapter.UpdatedAt = time.Now().Unix()
	updateCols = append(updateCols, "updated_at")
	upStms, uNames := coursez.ChapterTable.Update(updateCols...)
	updateQuery := global.CassSession.Query(upStms, uNames).BindStruct(&cassandraChapter)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraChapter.CreatedAt, 10)
	responseModel := model.Chapter{
		ID:          &cassandraChapter.ID,
		Name:        &cassandraChapter.Name,
		Description: &cassandraChapter.Description,
		CreatedAt:   &created,
		UpdatedAt:   &created,
		CourseID:    &cassandraChapter.CourseID,
		ModuleID:    &cassandraChapter.ModuleID,
		Sequence:    &cassandraChapter.Sequence,
	}
	return &responseModel, nil
}
