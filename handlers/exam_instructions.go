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

func ExamInstructionsCreate(ctx context.Context, exam *model.ExamInstructionInput) (*model.ExamInstruction, error) {
	log.Info("ExamInstructionsCreate called")
	guid := xid.New()
	cassandraQuestionBank := qbankz.ExamInstructions{
		ID:              guid.String(),
		Instructions:    *exam.Instructions,
		ExamID:          *exam.ExamID,
		PassingCriteria: *exam.PassingCriteria,
		NoAttempts:      *exam.NoAttempts,
		AccessType:      *exam.AccessType,
		CreatedBy:       *exam.CreatedBy,
		UpdatedBy:       *exam.UpdatedBy,
		CreatedAt:       time.Now().Unix(),
		UpdatedAt:       time.Now().Unix(),
		IsActive:        *exam.IsActive,
	}
	insertQuery := global.CassSessioQBank.Session.Query(qbankz.ExamInstructionsTable.Insert()).BindStruct(cassandraQuestionBank)
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
	cassandraQuestionBank := qbankz.ExamInstructions{
		ID: *input.ID,
	}
	banks := []qbankz.ExamInstructions{}
	getQuery := global.CassSessioQBank.Session.Query(qbankz.ExamInstructionsTable.Get()).BindMap(qb.M{"id": cassandraQuestionBank.ID})
	if err := getQuery.SelectRelease(&banks); err != nil {
		return nil, err
	}
	if len(banks) == 0 {
		return nil, fmt.Errorf("exams not found")
	}
	cassandraQuestionBank = banks[0]
	updatedCols := []string{}
	if input.Instructions != nil {
		cassandraQuestionBank.Instructions = *input.Instructions
		updatedCols = append(updatedCols, "instructions")
	}
	if input.ExamID != nil {
		cassandraQuestionBank.ExamID = *input.ExamID
		updatedCols = append(updatedCols, "exam_id")
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
	if input.PassingCriteria != nil {
		cassandraQuestionBank.PassingCriteria = *input.PassingCriteria
		updatedCols = append(updatedCols, "passing_criteria")
	}
	if input.NoAttempts != nil {
		cassandraQuestionBank.NoAttempts = *input.NoAttempts
		updatedCols = append(updatedCols, "no_attempts")
	}
	if input.AccessType != nil {
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
	updateQuery := global.CassSessioQBank.Session.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.ExamInstruction{
		ID:              &cassandraQuestionBank.ID,
		Instructions:   &cassandraQuestionBank.Instructions,
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
