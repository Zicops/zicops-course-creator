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

func ModuleCreate(ctx context.Context, courseID string, module *model.ModuleInput) (*model.Module, error) {
	log.Info("ModuleCreate called")
	guid := xid.New()
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspID := claims["lsp_id"].(string)
	cassandraModule := coursez.Module{
		ID:          guid.String(),
		Name:        *module.Name,
		Description: *module.Description,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		IsChapter:   *module.IsChapter,
		CourseID:    courseID,
		IsActive:    true,
		LspId:       lspID,
	}
	if module.Owner != nil {
		cassandraModule.Owner = *module.Owner
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
	insertQuery := CassSession.Query(coursez.ModuleTable.Insert()).BindStruct(cassandraModule)
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
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	cassandraModule := *GetModule(ctx, *module.ID, lspId, CassSession)
	updateCols := []string{}
	if module.Name != nil && cassandraModule.Name != *module.Name {
		updateCols = append(updateCols, "name")
		cassandraModule.Name = *module.Name
	}
	if module.Description != nil && cassandraModule.Description != *module.Description {
		updateCols = append(updateCols, "description")
		cassandraModule.Description = *module.Description
	}
	if module.IsChapter != nil && cassandraModule.IsChapter != *module.IsChapter {
		updateCols = append(updateCols, "ischapter")
		cassandraModule.IsChapter = *module.IsChapter
	}
	if module.Owner != nil && cassandraModule.Owner != *module.Owner {
		updateCols = append(updateCols, "owner")
		cassandraModule.Owner = *module.Owner
	}
	if module.Level != nil && cassandraModule.Level != *module.Level {
		updateCols = append(updateCols, "level")
		cassandraModule.Level = *module.Level
	}
	if module.Sequence != nil && cassandraModule.Sequence != *module.Sequence {
		updateCols = append(updateCols, "sequence")
		cassandraModule.Sequence = *module.Sequence
	}
	if module.SetGlobal != nil && cassandraModule.SetGlobal != *module.SetGlobal {
		updateCols = append(updateCols, "setglobal")
		cassandraModule.SetGlobal = *module.SetGlobal
	}
	if len(updateCols) > 0 {

		cassandraModule.UpdatedAt = time.Now().Unix()
		updateCols = append(updateCols, "updated_at")
		upStms, uNames := coursez.ModuleTable.Update(updateCols...)
		updateQuery := CassSession.Query(upStms, uNames).BindStruct(&cassandraModule)
		if err := updateQuery.ExecRelease(); err != nil {
			return nil, err
		}
	}
	updated := strconv.FormatInt(cassandraModule.UpdatedAt, 10)
	created := strconv.FormatInt(cassandraModule.CreatedAt, 10)
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

func GetModule(ctx context.Context, courseID string, lspID string, session *gocqlx.Session) *coursez.Module {
	chapters := []coursez.Module{}
	getQueryStr := fmt.Sprintf("SELECT * FROM coursez.module WHERE id='%s' and lsp_id='%s' and is_active=true", courseID, lspID)
	getQuery := session.Query(getQueryStr, nil)
	if err := getQuery.SelectRelease(&chapters); err != nil {
		return nil
	}
	return &chapters[0]
}
