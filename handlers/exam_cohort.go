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

func AddExamCohort(ctx context.Context, input *model.ExamCohortInput) (*model.ExamCohort, error) {
	log.Info("ExamInstructionsCreate called")
	guid := xid.New()
	cassandraQuestionBank := qbankz.ExamCohort{
		ID:        guid.String(),
		ExamID:    *input.ExamID,
		IsActive:  *input.IsActive,
		CreatedBy: *input.CreatedBy,
		UpdatedBy: *input.UpdatedBy,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		CohortID:  *input.CohortID,
	}
	insertQuery := global.CassSessioQBank.Session.Query(qbankz.ExamCohortTable.Insert()).BindStruct(cassandraQuestionBank)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.ExamCohort{
		ID:        &cassandraQuestionBank.ID,
		ExamID:    input.ExamID,
		IsActive:  input.IsActive,
		CreatedBy: input.CreatedBy,
		UpdatedBy: input.UpdatedBy,
		CreatedAt: &created,
		UpdatedAt: &updated,
		CohortID:  input.CohortID,
	}
	return &responseModel, nil
}

func UpdateExamCohort(ctx context.Context, input *model.ExamCohortInput) (*model.ExamCohort, error) {
	if input.ID == nil {
		return nil, fmt.Errorf("exam schedule id not found")
	}
	cassandraQuestionBank := qbankz.ExamCohort{
		ID: *input.ID,
	}
	banks := []qbankz.ExamCohort{}
	getQuery := global.CassSessioQBank.Session.Query(qbankz.ExamCohortTable.Get()).BindMap(qb.M{"id": cassandraQuestionBank.ID})
	if err := getQuery.SelectRelease(&banks); err != nil {
		return nil, err
	}
	if len(banks) == 0 {
		return nil, fmt.Errorf("exams not found")
	}
	cassandraQuestionBank = banks[0]
	updatedCols := []string{}
	if input.ExamID != nil {
		cassandraQuestionBank.ExamID = *input.ExamID
		updatedCols = append(updatedCols, "exam_id")
	}
	if input.CohortID != nil {
		cassandraQuestionBank.CohortID = *input.CohortID
		updatedCols = append(updatedCols, "cohort_id")
	}
	if input.IsActive != nil {
		cassandraQuestionBank.IsActive = *input.IsActive
		updatedCols = append(updatedCols, "is_active")
	}
	if input.CreatedBy != nil {
		cassandraQuestionBank.CreatedBy = *input.CreatedBy
		updatedCols = append(updatedCols, "created_by")
	}
	if input.UpdatedBy != nil {
		cassandraQuestionBank.UpdatedBy = *input.UpdatedBy
		updatedCols = append(updatedCols, "updated_by")
	}
	updatedAt := time.Now().Unix()
	cassandraQuestionBank.UpdatedAt = updatedAt
	updatedCols = append(updatedCols, "updated_at")
	if len(updatedCols) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}
	upStms, uNames := qbankz.ExamCohortTable.Update(updatedCols...)
	updateQuery := global.CassSessioQBank.Session.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.ExamCohort{
		ID:        &cassandraQuestionBank.ID,
		ExamID:    input.ExamID,
		IsActive:  input.IsActive,
		CreatedBy: input.CreatedBy,
		UpdatedBy: input.UpdatedBy,
		CreatedAt: &created,
		UpdatedAt: &updated,
		CohortID:  input.CohortID,
	}
	return &responseModel, nil
}
