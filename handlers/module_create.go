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
		Duration:    *module.Duration,
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
		Duration:    module.Duration,
		Level:       module.Level,
		Sequence:    module.Sequence,
		SetGlobal:   module.SetGlobal,
	}
	return &responseModel, nil
}

func UpdateModule(ctx context.Context, module *model.ModuleInput) (*model.Module, error) {

	log.Info("ModuleUpdate called")
	if module.ID == nil {
		return nil, fmt.Errorf("module id is required")
	}
	cassandraModule := coursez.Module{
		ID: *module.ID,
	}
	// set course in cassandra
	getQuery := global.CassSession.Session.Query(coursez.ModuleTable.Get()).BindStruct(cassandraModule)
	if err := getQuery.ExecRelease(); err != nil {
		return nil, err
	}
	if module.Name != "" {
		cassandraModule.Name = module.Name
	}
	if module.Description != "" {
		cassandraModule.Description = module.Description
	}
	cassandraModule.IsChapter = module.IsChapter
	if module.Owner != nil {
		cassandraModule.Owner = *module.Owner
	}
	if module.Duration != nil {
		cassandraModule.Duration = *module.Duration
	}
	if module.Level != nil {
		cassandraModule.Level = *module.Level
	}
	if module.Sequence != nil {
		cassandraModule.Sequence = *module.Sequence
	}
	if module.SetGlobal != nil {
		cassandraModule.SetGlobal = *module.SetGlobal
	}
	cassandraModule.UpdatedAt = time.Now().Unix()
	updated := strconv.FormatInt(cassandraModule.UpdatedAt, 10)
	created := strconv.FormatInt(cassandraModule.CreatedAt, 10)
	// set course in cassandra
	updateQuery := global.CassSession.Session.Query(coursez.ModuleTable.Update()).BindStruct(cassandraModule)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	responseModel := model.Module{
		ID:          &cassandraModule.ID,
		Name:        module.Name,
		Description: module.Description,
		CreatedAt:   &created,
		UpdatedAt:   &updated,
		IsChapter:   module.IsChapter,
		CourseID:    cassandraModule.CourseID,
		Owner:       module.Owner,
		Duration:    module.Duration,
		Level:       module.Level,
		Sequence:    module.Sequence,
		SetGlobal:   module.SetGlobal,
	}
	return &responseModel, nil
}
