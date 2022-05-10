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

func QuestionSectionMap(ctx context.Context, input *model.MapSectionToBankInput) (*model.SectionQBMapping, error) {
	log.Info("QuestionSectionMap called")
	guid := xid.New()
	cassandraQuestionBank := qbankz.SectionQBMapping{
		ID:              guid.String(),
		QBId:            *input.QbID,
		SectionID:       *input.SectionID,
		IsActive:        *input.IsActive,
		CreatedBy:       *input.CreatedBy,
		UpdatedBy:       *input.UpdatedBy,
		CreatedAt:       time.Now().Unix(),
		UpdatedAt:       time.Now().Unix(),
		DifficultyLevel: *input.DifficultyLevel,
		TotalQuestions:  *input.TotalQuestions,
		QuestionType:    *input.QuestionType,
		QuestionMarks:   *input.QuestionMarks,
		RetrievalType:   *input.RetrieveType,
	}

	insertQuery := global.CassSession.Session.Query(qbankz.SectionQBMappingTable.Insert()).BindStruct(cassandraQuestionBank)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.SectionQBMapping{
		ID:              &cassandraQuestionBank.ID,
		IsActive:        input.IsActive,
		CreatedBy:       input.CreatedBy,
		UpdatedBy:       input.UpdatedBy,
		CreatedAt:       &created,
		UpdatedAt:       &updated,
		DifficultyLevel: input.DifficultyLevel,
		TotalQuestions:  input.TotalQuestions,
		QuestionType:    input.QuestionType,
		QuestionMarks:   input.QuestionMarks,
		RetrieveType:    input.RetrieveType,
	}
	return &responseModel, nil
}

func QuestionSectionMapUpdate(ctx context.Context, input *model.MapSectionToBankInput) (*model.SectionQBMapping, error) {
	log.Info("QuestionPaperUpdate called")
	if input.ID == nil {
		return nil, fmt.Errorf("section id not found")
	}
	cassandraQuestionBank := qbankz.SectionQBMapping{
		ID: *input.ID,
	}
	banks := []qbankz.SectionQBMapping{}
	getQuery := global.CassSession.Session.Query(qbankz.SectionQBMappingTable.Get()).BindMap(qb.M{"id": cassandraQuestionBank.ID})
	if err := getQuery.SelectRelease(&banks); err != nil {
		return nil, err
	}
	if len(banks) == 0 {
		return nil, fmt.Errorf("question bank not found")
	}
	cassandraQuestionBank = banks[0]
	updatedCols := []string{}
	if input.IsActive != nil {
		cassandraQuestionBank.IsActive = *input.IsActive
		updatedCols = append(updatedCols, "is_active")
	}
	if input.UpdatedBy != nil {
		cassandraQuestionBank.UpdatedBy = *input.UpdatedBy
		updatedCols = append(updatedCols, "updated_by")
	}
	if input.DifficultyLevel != nil {
		cassandraQuestionBank.DifficultyLevel = *input.DifficultyLevel
		updatedCols = append(updatedCols, "difficulty_level")
	}
	if input.TotalQuestions != nil {
		cassandraQuestionBank.TotalQuestions = *input.TotalQuestions
		updatedCols = append(updatedCols, "total_questions")
	}
	if input.QuestionType != nil {
		cassandraQuestionBank.QuestionType = *input.QuestionType
		updatedCols = append(updatedCols, "question_type")
	}
	if input.QuestionMarks != nil {
		cassandraQuestionBank.QuestionMarks = *input.QuestionMarks
		updatedCols = append(updatedCols, "question_marks")
	}
	if input.RetrieveType != nil {
		cassandraQuestionBank.RetrievalType = *input.RetrieveType
		updatedCols = append(updatedCols, "retrieval_type")
	}
	
	updatedAt := time.Now().Unix()
	cassandraQuestionBank.UpdatedAt = updatedAt
	updatedCols = append(updatedCols, "updated_at")
	if len(updatedCols) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}
	upStms, uNames := qbankz.SectionQBMappingTable.Update(updatedCols...)
	updateQuery := global.CassSession.Session.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.SectionQBMapping{
		ID:              &cassandraQuestionBank.ID,
		IsActive:        input.IsActive,
		CreatedBy:       input.CreatedBy,
		UpdatedBy:       input.UpdatedBy,
		CreatedAt:       &created,
		UpdatedAt:       &updated,
		DifficultyLevel: input.DifficultyLevel,
		TotalQuestions:  input.TotalQuestions,
		QuestionType:    input.QuestionType,
		QuestionMarks:   input.QuestionMarks,
		RetrieveType:    input.RetrieveType,
	}
	return &responseModel, nil
}
