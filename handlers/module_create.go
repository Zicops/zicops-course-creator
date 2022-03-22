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
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/graph/model"
)

func ModuleCreate(ctx context.Context, courseID string, module *model.ModuleInput) (*model.Module, error) {
	log.Info("ModuleCreate called")
	guid := xid.New()
	cassandraModule := coursez.Module{
		ID:          guid.String(),
		Name:        *module.Name,
		Description: *module.Description,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		IsChapter:   *module.IsChapter,
		CourseID:    courseID,
		IsActive:    false,
	}
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
		CourseID:    &cassandraModule.CourseID,
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
	modules := []coursez.Module{}
	getQuery := global.CassSession.Session.Query(coursez.ModuleTable.Get()).BindMap(qb.M{"id": cassandraModule.ID})
	if err := getQuery.SelectRelease(&modules); err != nil {
		return nil, err
	}
	if len(modules) < 1 {
		return nil, fmt.Errorf("module not found")
	}
	cassandraModule = modules[0]
	updateCols := []string{}
	if *module.Name != "" {
		updateCols = append(updateCols, "name")
		cassandraModule.Name = *module.Name
	}
	if *module.Description != "" {
		updateCols = append(updateCols, "description")
		cassandraModule.Description = *module.Description
	}
	cassandraModule.IsChapter = *module.IsChapter
	if module.Owner != nil {
		updateCols = append(updateCols, "owner")
		cassandraModule.Owner = *module.Owner
	}
	if module.Duration != nil {
		updateCols = append(updateCols, "duration")
		cassandraModule.Duration = *module.Duration
	}
	if module.Level != nil {
		updateCols = append(updateCols, "level")
		cassandraModule.Level = *module.Level
	}
	if module.Sequence != nil {
		updateCols = append(updateCols, "sequence")
		cassandraModule.Sequence = *module.Sequence
	}
	if module.SetGlobal != nil {
		updateCols = append(updateCols, "setglobal")
		cassandraModule.SetGlobal = *module.SetGlobal
	}
	cassandraModule.UpdatedAt = time.Now().Unix()
	updateCols = append(updateCols, "updated_at")
	updated := strconv.FormatInt(cassandraModule.UpdatedAt, 10)
	created := strconv.FormatInt(cassandraModule.CreatedAt, 10)
	upStms, uNames := coursez.ModuleTable.Update(updateCols...)
	updateQuery := global.CassSession.Session.Query(upStms, uNames).BindStruct(&cassandraModule)
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
		CourseID:    &cassandraModule.CourseID,
		Owner:       module.Owner,
		Duration:    module.Duration,
		Level:       module.Level,
		Sequence:    module.Sequence,
		SetGlobal:   module.SetGlobal,
	}
	return &responseModel, nil
}
