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

func ModuleCreate(ctx context.Context, courseID string, module *model.ModuleInput) (*model.Module, error) {
	log.Info("ModuleCreate called")
	guid := xid.New()
	cassandraModule := coursez.Module{
		ID:          guid.String(),
		Name:        module.Name,
		Description: module.Description,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		IsChapter:   module.IsChapter,
		CourseID:    courseID,
		Owner:       *module.Owner,
		Level:       *module.Level,
		Sequence:    *module.Sequence,
		SetGlobal:   *module.SetGlobal,
		IsDeleted:   false,
	}
	// set course in cassandra
	insertQuery := global.CassSession.Session.Query(coursez.ModuleTable.Insert()).BindStruct(cassandraModule)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraModule.CreatedAt, 10)
	responseModel := model.Module{
		ID:          &cassandraModule.ID,
		Name:        module.Name,
		Description: module.Description,
		CreatedAt:   &created,
		UpdatedAt:   &created,
		IsChapter:   module.IsChapter,
		CourseID:    cassandraModule.CourseID,
		Owner:       module.Owner,
		Level:       module.Level,
		Sequence:    module.Sequence,
		SetGlobal:   module.SetGlobal,
	}
	return &responseModel, nil
}
