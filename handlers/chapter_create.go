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
		ModuleID:    *chapter.ModuleID,
		Sequence:    *chapter.Sequence,
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
