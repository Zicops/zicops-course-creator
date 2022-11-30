package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/scylladb/gocqlx/v2"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/qbankz"
	"github.com/zicops/zicops-cass-pool/cassandra"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/helpers"
)

func QuestionSectionMap(ctx context.Context, input *model.MapSectionToBankInput) (*model.SectionQBMapping, error) {
	log.Info("QuestionSectionMap called")

	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	email_creator := claims["email"].(string)
	lspID := claims["lsp_id"].(string)
	cassandraQuestionBank := qbankz.SectionQBMapping{
		ID:              uuid.New().String(),
		QBId:            *input.QbID,
		SectionID:       *input.SectionID,
		IsActive:        *input.IsActive,
		CreatedBy:       email_creator,
		UpdatedBy:       email_creator,
		CreatedAt:       time.Now().Unix(),
		UpdatedAt:       time.Now().Unix(),
		DifficultyLevel: *input.DifficultyLevel,
		TotalQuestions:  *input.TotalQuestions,
		QuestionType:    *input.QuestionType,
		QuestionMarks:   *input.QuestionMarks,
		RetrievalType:   *input.RetrieveType,
		LspId:           lspID,
	}

	insertQuery := CassSession.Query(qbankz.SectionQBMappingTable.Insert()).BindStruct(cassandraQuestionBank)
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
		QbID:            input.QbID,
		SectionID:       input.SectionID,
	}
	return &responseModel, nil
}

func QuestionSectionMapUpdate(ctx context.Context, input *model.MapSectionToBankInput) (*model.SectionQBMapping, error) {
	log.Info("QuestionPaperUpdate called")
	if input.ID == nil {
		return nil, fmt.Errorf("section id not found")
	}
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	email_creator := claims["email"].(string)
	lspID := claims["lsp_id"].(string)
	cassandraQuestionBank := *GetSecQBMap(ctx, *input.ID, lspID, CassSession)
	updatedCols := []string{}
	if email_creator != "" && cassandraQuestionBank.CreatedBy != email_creator {
		cassandraQuestionBank.UpdatedBy = email_creator
		updatedCols = append(updatedCols, "updated_by")
	}
	if input.DifficultyLevel != nil && cassandraQuestionBank.DifficultyLevel != *input.DifficultyLevel {
		cassandraQuestionBank.DifficultyLevel = *input.DifficultyLevel
		updatedCols = append(updatedCols, "difficulty_level")
	}
	if input.TotalQuestions != nil && cassandraQuestionBank.TotalQuestions != *input.TotalQuestions {
		cassandraQuestionBank.TotalQuestions = *input.TotalQuestions
		updatedCols = append(updatedCols, "total_questions")
	}
	if input.QuestionType != nil && cassandraQuestionBank.QuestionType != *input.QuestionType {
		cassandraQuestionBank.QuestionType = *input.QuestionType
		updatedCols = append(updatedCols, "question_type")
	}
	if input.QuestionMarks != nil && cassandraQuestionBank.QuestionMarks != *input.QuestionMarks {
		cassandraQuestionBank.QuestionMarks = *input.QuestionMarks
		updatedCols = append(updatedCols, "question_marks")
	}
	if input.RetrieveType != nil && cassandraQuestionBank.RetrievalType != *input.RetrieveType {
		cassandraQuestionBank.RetrievalType = *input.RetrieveType
		updatedCols = append(updatedCols, "retrieval_type")
	}
	if input.QbID != nil && cassandraQuestionBank.QBId != *input.QbID {
		cassandraQuestionBank.QBId = *input.QbID
		updatedCols = append(updatedCols, "qb_id")
	}
	if input.SectionID != nil && cassandraQuestionBank.SectionID != *input.SectionID {
		cassandraQuestionBank.SectionID = *input.SectionID
		updatedCols = append(updatedCols, "section_id")
	}
	updatedAt := time.Now().Unix()
	if len(updatedCols) > 0 {
		cassandraQuestionBank.UpdatedAt = updatedAt
		updatedCols = append(updatedCols, "updated_at")
		upStms, uNames := qbankz.SectionQBMappingTable.Update(updatedCols...)
		updateQuery := CassSession.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
		if err := updateQuery.ExecRelease(); err != nil {
			return nil, err
		}
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
		QbID:            input.QbID,
		SectionID:       input.SectionID,
	}
	return &responseModel, nil
}

func QuestionFixed(ctx context.Context, input *model.SectionFixedQuestionsInput) (*model.SectionFixedQuestions, error) {
	log.Info("QuestionFixed called")

	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	email_creator := claims["email"].(string)
	lspID := claims["lsp_id"].(string)
	cassandraQuestionBank := qbankz.SectionFixedQuestions{
		ID:         uuid.New().String(),
		SQBId:      *input.SqbID,
		QuestionID: *input.QuestionID,
		IsActive:   *input.IsActive,
		CreatedBy:  email_creator,
		UpdatedBy:  email_creator,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
		LspId:      lspID,
	}

	insertQuery := CassSession.Query(qbankz.SectionFixedQuestionsTable.Insert()).BindStruct(cassandraQuestionBank)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.SectionFixedQuestions{
		ID:         &cassandraQuestionBank.ID,
		IsActive:   input.IsActive,
		CreatedBy:  input.CreatedBy,
		UpdatedBy:  input.UpdatedBy,
		CreatedAt:  &created,
		UpdatedAt:  &updated,
		SqbID:      input.SqbID,
		QuestionID: input.QuestionID,
	}
	return &responseModel, nil
}

func QuestionFixedUpdate(ctx context.Context, input *model.SectionFixedQuestionsInput) (*model.SectionFixedQuestions, error) {
	log.Info("QuestionFixedUpdate called")
	if input.ID == nil {
		return nil, fmt.Errorf("id not found")
	}
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	email_creator := claims["email"].(string)
	lspID := claims["lsp_id"].(string)
	cassandraQuestionBank := *GetSecFixedQs(ctx, *input.ID, lspID, CassSession)
	updatedCols := []string{}
	if email_creator != "" && cassandraQuestionBank.UpdatedBy != email_creator {
		cassandraQuestionBank.UpdatedBy = email_creator
		updatedCols = append(updatedCols, "updated_by")
	}
	if input.SqbID != nil && cassandraQuestionBank.SQBId != *input.SqbID {
		cassandraQuestionBank.SQBId = *input.SqbID
		updatedCols = append(updatedCols, "sqb_id")
	}
	if input.QuestionID != nil && cassandraQuestionBank.QuestionID != *input.QuestionID {
		cassandraQuestionBank.QuestionID = *input.QuestionID
		updatedCols = append(updatedCols, "question_id")
	}
	if len(updatedCols) > 0 {
		updatedAt := time.Now().Unix()
		cassandraQuestionBank.UpdatedAt = updatedAt
		updatedCols = append(updatedCols, "updated_at")
		upStms, uNames := qbankz.SectionFixedQuestionsTable.Update(updatedCols...)
		updateQuery := CassSession.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
		if err := updateQuery.ExecRelease(); err != nil {
			return nil, err
		}
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.SectionFixedQuestions{
		ID:         &cassandraQuestionBank.ID,
		IsActive:   input.IsActive,
		CreatedBy:  input.CreatedBy,
		UpdatedBy:  input.UpdatedBy,
		CreatedAt:  &created,
		UpdatedAt:  &updated,
		SqbID:      input.SqbID,
		QuestionID: input.QuestionID,
	}
	return &responseModel, nil
}

func GetSecQBMap(ctx context.Context, courseID string, lspID string, session *gocqlx.Session) *qbankz.SectionQBMapping {
	chapters := []qbankz.SectionQBMapping{}
	getQueryStr := fmt.Sprintf("SELECT * FROM qbankz.section_qb_mapping WHERE id='%s' and lsp_id='%s' and is_active=true", courseID, lspID)
	getQuery := session.Query(getQueryStr, nil)
	if err := getQuery.SelectRelease(&chapters); err != nil {
		return nil
	}
	return &chapters[0]
}

func GetSecFixedQs(ctx context.Context, courseID string, lspID string, session *gocqlx.Session) *qbankz.SectionFixedQuestions {
	chapters := []qbankz.SectionFixedQuestions{}
	getQueryStr := fmt.Sprintf("SELECT * FROM qbankz.section_fixed_questions WHERE id='%s' and lsp_id='%s' and is_active=true", courseID, lspID)
	getQuery := session.Query(getQueryStr, nil)
	if err := getQuery.SelectRelease(&chapters); err != nil {
		return nil
	}
	return &chapters[0]
}
