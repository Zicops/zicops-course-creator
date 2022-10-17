package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/xid"
	"github.com/scylladb/gocqlx/v2"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/qbankz"
	"github.com/zicops/zicops-cass-pool/cassandra"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/helpers"
)

func ExamInstructionsCreate(ctx context.Context, exam *model.ExamInstructionInput) (*model.ExamInstruction, error) {
	log.Info("ExamInstructionsCreate called")
	guid := xid.New()
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspID := claims["lsp_id"].(string)
	email_creator := claims["email"].(string)
	cassandraQuestionBank := qbankz.ExamInstructions{
		ID:              guid.String(),
		Instructions:    *exam.Instructions,
		ExamID:          *exam.ExamID,
		PassingCriteria: *exam.PassingCriteria,
		NoAttempts:      *exam.NoAttempts,
		AccessType:      *exam.AccessType,
		CreatedBy:       email_creator,
		UpdatedBy:       email_creator,
		CreatedAt:       time.Now().Unix(),
		UpdatedAt:       time.Now().Unix(),
		IsActive:        *exam.IsActive,
		LspId:           lspID,
	}
	insertQuery := CassSession.Query(qbankz.ExamInstructionsTable.Insert()).BindStruct(cassandraQuestionBank)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.ExamInstruction{
		ID:              &cassandraQuestionBank.ID,
		ExamID:          exam.ExamID,
		Instructions:    exam.Instructions,
		PassingCriteria: exam.PassingCriteria,
		NoAttempts:      exam.NoAttempts,
		AccessType:      exam.AccessType,
		UpdatedBy:       exam.UpdatedBy,
		CreatedBy:       exam.CreatedBy,
		CreatedAt:       &created,
		UpdatedAt:       &updated,
		IsActive:        exam.IsActive,
	}
	return &responseModel, nil
}

func ExamInstructionsUpdate(ctx context.Context, input *model.ExamInstructionInput) (*model.ExamInstruction, error) {
	if input.ID == nil {
		return nil, fmt.Errorf("exam schedule id not found")
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
	cassandraQuestionBank := *GetExamInstructions(ctx, *input.ID, lspID, CassSession)
	updatedCols := []string{}
	if input.Instructions != nil && cassandraQuestionBank.Instructions != *input.Instructions {
		cassandraQuestionBank.Instructions = *input.Instructions
		updatedCols = append(updatedCols, "instructions")
	}
	if input.ExamID != nil && cassandraQuestionBank.ExamID != *input.ExamID {
		cassandraQuestionBank.ExamID = *input.ExamID
		updatedCols = append(updatedCols, "exam_id")
	}

	if email_creator != "" && cassandraQuestionBank.UpdatedBy != email_creator {
		cassandraQuestionBank.UpdatedBy = email_creator
		updatedCols = append(updatedCols, "updated_by")
	}
	if input.PassingCriteria != nil && cassandraQuestionBank.PassingCriteria != *input.PassingCriteria {
		cassandraQuestionBank.PassingCriteria = *input.PassingCriteria
		updatedCols = append(updatedCols, "passing_criteria")
	}
	if input.NoAttempts != nil && cassandraQuestionBank.NoAttempts != *input.NoAttempts {
		cassandraQuestionBank.NoAttempts = *input.NoAttempts
		updatedCols = append(updatedCols, "no_attempts")
	}
	if input.AccessType != nil && cassandraQuestionBank.AccessType != *input.AccessType {
		cassandraQuestionBank.AccessType = *input.AccessType
		updatedCols = append(updatedCols, "access_type")
	}
	updatedAt := time.Now().Unix()
	cassandraQuestionBank.UpdatedAt = updatedAt
	updatedCols = append(updatedCols, "updated_at")
	if len(updatedCols) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}
	upStms, uNames := qbankz.ExamInstructionsTable.Update(updatedCols...)
	updateQuery := CassSession.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.ExamInstruction{
		ID:              &cassandraQuestionBank.ID,
		Instructions:    &cassandraQuestionBank.Instructions,
		ExamID:          input.ExamID,
		PassingCriteria: input.PassingCriteria,
		NoAttempts:      input.NoAttempts,
		AccessType:      input.AccessType,
		UpdatedBy:       input.UpdatedBy,
		CreatedAt:       &created,
		UpdatedAt:       &updated,
		IsActive:        input.IsActive,
		CreatedBy:       input.CreatedBy,
	}
	return &responseModel, nil
}

func GetExamInstructions(ctx context.Context, id string, lspID string, session *gocqlx.Session) *qbankz.ExamInstructions {
	chapters := []qbankz.ExamInstructions{}
	getQueryStr := fmt.Sprintf("SELECT * FROM qbankz.exam_instructions WHERE id='%s' and lsp_id='%s' and is_active=true", id, lspID)
	getQuery := session.Query(getQueryStr, nil)
	if err := getQuery.SelectRelease(&chapters); err != nil {
		return nil
	}
	return &chapters[0]
}
