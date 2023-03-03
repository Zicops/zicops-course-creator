package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/scylladb/gocqlx/v2"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/helpers"
)

func ChapterCreate(ctx context.Context, courseID string, chapter *model.ChapterInput) (*model.Chapter, error) {
	log.Info("ChapterCreate called")
	session, err := global.CassPool.GetSession(ctx, "coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspID := claims["lsp_id"].(string)

	cassandraChapter := coursez.Chapter{
		ID:          uuid.New().String(),
		Name:        *chapter.Name,
		Description: *chapter.Description,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		CourseID:    courseID,
		LspId:       lspID,
		IsActive:    true,
	}
	if chapter.ModuleID != nil {
		cassandraChapter.ModuleID = *chapter.ModuleID
	}
	if chapter.Sequence != nil {
		cassandraChapter.Sequence = *chapter.Sequence
	}
	// set course in cassandra
	insertQuery := CassSession.Query(coursez.ChapterTable.Insert()).BindStruct(cassandraChapter)
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
	session, err := global.CassPool.GetSession(ctx, "coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	if chapter.ID == nil {
		return nil, fmt.Errorf("chapter not found")
	}
	cassandraChapter := coursez.Chapter{
		ID: *chapter.ID,
	}
	cassandraChapter = *GetChapter(ctx, cassandraChapter.ID, lspId, CassSession)
	updateCols := make([]string, 0)
	if *chapter.Name != "" && *chapter.Name != cassandraChapter.Name {
		updateCols = append(updateCols, "name")
		cassandraChapter.Name = *chapter.Name
	}
	if *chapter.Description != "" && *chapter.Description != cassandraChapter.Description {
		updateCols = append(updateCols, "description")
		cassandraChapter.Description = *chapter.Description
	}
	if chapter.Sequence != nil && *chapter.Sequence != cassandraChapter.Sequence {
		updateCols = append(updateCols, "sequence")
		cassandraChapter.Sequence = *chapter.Sequence
	}
	if chapter.ModuleID != nil && *chapter.ModuleID != cassandraChapter.ModuleID {
		updateCols = append(updateCols, "moduleid")
		cassandraChapter.ModuleID = *chapter.ModuleID
	}
	if len(updateCols) > 0 {
		cassandraChapter.UpdatedAt = time.Now().Unix()
		updateCols = append(updateCols, "updated_at")
		upStms, uNames := coursez.ChapterTable.Update(updateCols...)
		updateQuery := CassSession.Query(upStms, uNames).BindStruct(&cassandraChapter)
		if err := updateQuery.ExecRelease(); err != nil {
			return nil, err
		}
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

func GetChapter(ctx context.Context, courseID string, lspID string, session *gocqlx.Session) *coursez.Chapter {
	chapters := []coursez.Chapter{}
	getQueryStr := fmt.Sprintf("SELECT * FROM coursez.chapter WHERE id='%s' and lsp_id='%s' and is_active=true", courseID, lspID)
	getQuery := session.Query(getQueryStr, nil)
	if err := getQuery.SelectRelease(&chapters); err != nil {
		return nil
	}
	return &chapters[0]
}
