package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/xid"
	"github.com/scylladb/gocqlx/v2/qb"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/qbankz"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/graph/model"
)

func QuestionPaperCreate(ctx context.Context, input *model.QuestionPaperInput) (*model.QuestionPaper, error) {
	log.Info("QuestionPaperCreate called")
	guid := xid.New()
	cassandraQuestionBank := qbankz.QuestionPaperMain{
		ID:                guid.String(),
		Name:              *input.Name,
		Category:          *input.Category,
		SubCategory:       *input.SubCategory,
		IsActive:          *input.IsActive,
		CreatedBy:         *input.CreatedBy,
		UpdatedBy:         *input.UpdatedBy,
		CreatedAt:         time.Now().Unix(),
		UpdatedAt:         time.Now().Unix(),
		Description:       *input.Description,
		DifficultyLevel:   *input.DifficultyLevel,
		SectionWise:       *input.SectionWise,
		SuggestedDuration: *input.SuggestedDuration,
	}

	insertQuery := global.CassSessioQBank.Session.Query(qbankz.QuestionPaperMainTable.Insert()).BindStruct(cassandraQuestionBank)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.QuestionPaper{
		ID:                &cassandraQuestionBank.ID,
		Name:              input.Name,
		Category:          input.Category,
		SubCategory:       input.SubCategory,
		IsActive:          input.IsActive,
		CreatedBy:         input.CreatedBy,
		UpdatedBy:         input.UpdatedBy,
		CreatedAt:         &created,
		UpdatedAt:         &updated,
		Description:       input.Description,
		DifficultyLevel:   input.DifficultyLevel,
		SectionWise:       input.SectionWise,
		SuggestedDuration: input.SuggestedDuration,
	}
	return &responseModel, nil
}

func QuestionPaperUpdate(ctx context.Context, input *model.QuestionPaperInput) (*model.QuestionPaper, error) {
	log.Info("QuestionPaperUpdate called")
	if input.ID == nil {
		return nil, fmt.Errorf("question paper id not found")
	}
	cassandraQuestionBank := qbankz.QuestionPaperMain{
		ID: *input.ID,
	}
	banks := []qbankz.QuestionPaperMain{}
	getQuery := global.CassSessioQBank.Session.Query(qbankz.QuestionPaperMainTable.Get()).BindMap(qb.M{"id": cassandraQuestionBank.ID})
	if err := getQuery.SelectRelease(&banks); err != nil {
		return nil, err
	}
	if len(banks) == 0 {
		return nil, fmt.Errorf("question bank not found")
	}
	cassandraQuestionBank = banks[0]
	updatedCols := []string{}
	if input.Name != nil {
		cassandraQuestionBank.Name = *input.Name
		updatedCols = append(updatedCols, "name")
	}
	if input.Category != nil {
		cassandraQuestionBank.Category = *input.Category
		updatedCols = append(updatedCols, "category")
	}
	if input.SubCategory != nil {
		cassandraQuestionBank.SubCategory = *input.SubCategory
		updatedCols = append(updatedCols, "sub_category")
	}
	if input.IsActive != nil {
		cassandraQuestionBank.IsActive = *input.IsActive
		updatedCols = append(updatedCols, "is_active")
	}
	if input.UpdatedBy != nil {
		cassandraQuestionBank.UpdatedBy = *input.UpdatedBy
		updatedCols = append(updatedCols, "updated_by")
	}
	if input.CreatedBy != nil {
		cassandraQuestionBank.CreatedBy = *input.CreatedBy
		updatedCols = append(updatedCols, "created_by")
	}
	if input.Description != nil {
		cassandraQuestionBank.Description = *input.Description
		updatedCols = append(updatedCols, "description")
	}
	if input.DifficultyLevel != nil {
		cassandraQuestionBank.DifficultyLevel = *input.DifficultyLevel
		updatedCols = append(updatedCols, "difficulty_level")
	}
	if input.SectionWise != nil {
		cassandraQuestionBank.SectionWise = *input.SectionWise
		updatedCols = append(updatedCols, "section_wise")
	}
	if input.SuggestedDuration != nil {
		cassandraQuestionBank.SuggestedDuration = *input.SuggestedDuration
		updatedCols = append(updatedCols, "suggested_duration")
	}
	updatedAt := time.Now().Unix()
	cassandraQuestionBank.UpdatedAt = updatedAt
	updatedCols = append(updatedCols, "updated_at")
	if len(updatedCols) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}
	upStms, uNames := qbankz.QuestionPaperMainTable.Update(updatedCols...)
	updateQuery := global.CassSessioQBank.Session.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.QuestionPaper{
		ID:                &cassandraQuestionBank.ID,
		Name:              input.Name,
		Category:          input.Category,
		SubCategory:       input.SubCategory,
		IsActive:          input.IsActive,
		CreatedBy:         input.CreatedBy,
		UpdatedBy:         input.UpdatedBy,
		CreatedAt:         &created,
		UpdatedAt:         &updated,
		Description:       input.Description,
		DifficultyLevel:   input.DifficultyLevel,
		SectionWise:       input.SectionWise,
		SuggestedDuration: input.SuggestedDuration,
	}
	return &responseModel, nil
}

func QuestionPaperSectionCreate(ctx context.Context, input *model.QuestionPaperSectionInput) (*model.QuestionPaperSection, error) {
	log.Info("QuestionPaperSectionCreate called")
	guid := xid.New()
	cassandraQuestionBank := qbankz.SectionMain{
		ID:              guid.String(),
		Name:            *input.Name,
		IsActive:        *input.IsActive,
		CreatedBy:       *input.CreatedBy,
		UpdatedBy:       *input.UpdatedBy,
		CreatedAt:       time.Now().Unix(),
		UpdatedAt:       time.Now().Unix(),
		Description:     *input.Description,
		DifficultyLevel: *input.DifficultyLevel,
		QPID:            *input.QpID,
		Type:            "",
		TotalQuestions:  0,
	}

	insertQuery := global.CassSessioQBank.Session.Query(qbankz.SectionMainTable.Insert()).BindStruct(cassandraQuestionBank)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.QuestionPaperSection{
		ID:              &cassandraQuestionBank.ID,
		Name:            input.Name,
		IsActive:        input.IsActive,
		CreatedBy:       input.CreatedBy,
		UpdatedBy:       input.UpdatedBy,
		CreatedAt:       &created,
		UpdatedAt:       &updated,
		Description:     input.Description,
		DifficultyLevel: input.DifficultyLevel,
		QpID:            input.QpID,
		Type:            input.Type,
		TotalQuestions:  input.TotalQuestions,
	}
	return &responseModel, nil
}

func QuestionPaperSectionUpdate(ctx context.Context, input *model.QuestionPaperSectionInput) (*model.QuestionPaperSection, error) {
	log.Info("QuestionPaperUpdate called")
	if input.ID == nil {
		return nil, fmt.Errorf("section id not found")
	}
	cassandraQuestionBank := qbankz.SectionMain{
		ID: *input.ID,
	}
	banks := []qbankz.SectionMain{}
	getQuery := global.CassSessioQBank.Session.Query(qbankz.SectionMainTable.Get()).BindMap(qb.M{"id": cassandraQuestionBank.ID})
	if err := getQuery.SelectRelease(&banks); err != nil {
		return nil, err
	}
	if len(banks) == 0 {
		return nil, fmt.Errorf("question bank not found")
	}
	cassandraQuestionBank = banks[0]
	updatedCols := []string{}
	if input.Name != nil {
		cassandraQuestionBank.Name = *input.Name
		updatedCols = append(updatedCols, "name")
	}

	if input.IsActive != nil {
		cassandraQuestionBank.IsActive = *input.IsActive
		updatedCols = append(updatedCols, "is_active")
	}
	if input.UpdatedBy != nil {
		cassandraQuestionBank.UpdatedBy = *input.UpdatedBy
		updatedCols = append(updatedCols, "updated_by")
	}
	if input.CreatedBy != nil {
		cassandraQuestionBank.CreatedBy = *input.CreatedBy
		updatedCols = append(updatedCols, "created_by")
	}
	if input.Description != nil {
		cassandraQuestionBank.Description = *input.Description
		updatedCols = append(updatedCols, "description")
	}
	if input.DifficultyLevel != nil {
		cassandraQuestionBank.DifficultyLevel = *input.DifficultyLevel
		updatedCols = append(updatedCols, "difficulty_level")
	}
	if input.Type != nil {
		cassandraQuestionBank.Type = *input.Type
		updatedCols = append(updatedCols, "type")
	}
	if input.TotalQuestions != nil {
		cassandraQuestionBank.TotalQuestions = *input.TotalQuestions
		updatedCols = append(updatedCols, "total_questions")
	}
	if input.QpID != nil {
		cassandraQuestionBank.QPID = *input.QpID
		updatedCols = append(updatedCols, "qp_id")
	}
	updatedAt := time.Now().Unix()
	cassandraQuestionBank.UpdatedAt = updatedAt
	updatedCols = append(updatedCols, "updated_at")
	if len(updatedCols) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}
	upStms, uNames := qbankz.SectionMainTable.Update(updatedCols...)
	updateQuery := global.CassSessioQBank.Session.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.QuestionPaperSection{
		ID:              &cassandraQuestionBank.ID,
		Name:            input.Name,
		IsActive:        input.IsActive,
		CreatedBy:       input.CreatedBy,
		UpdatedBy:       input.UpdatedBy,
		CreatedAt:       &created,
		UpdatedAt:       &updated,
		Description:     input.Description,
		DifficultyLevel: input.DifficultyLevel,
		QpID:            input.QpID,
		Type:            input.Type,
		TotalQuestions:  input.TotalQuestions,
	}
	return &responseModel, nil
}
