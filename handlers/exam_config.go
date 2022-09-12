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
	"github.com/zicops/zicops-cass-pool/cassandra"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/graph/model"
)

func AddExamConfiguration(ctx context.Context, input *model.ExamConfigurationInput) (*model.ExamConfiguration, error) {
	log.Info("ExamInstructionsCreate called")
	guid := xid.New()
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	global.CassSessioQBank = session
	cassandraQuestionBank := qbankz.ExamConfig{
		ID:           guid.String(),
		ExamID:       *input.ExamID,
		IsActive:     *input.IsActive,
		CreatedBy:    *input.CreatedBy,
		UpdatedBy:    *input.UpdatedBy,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
		Shuffle:      *input.Shuffle,
		DisplayHints: *input.DisplayHints,
		ShowAnswer:   *input.ShowAnswer,
		ShowResult:   *input.ShowResult,
	}
	insertQuery := global.CassSessioQBank.Query(qbankz.ExamConfigTable.Insert()).BindStruct(cassandraQuestionBank)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.ExamConfiguration{
		ID:           &cassandraQuestionBank.ID,
		ExamID:       input.ExamID,
		IsActive:     input.IsActive,
		CreatedBy:    input.CreatedBy,
		UpdatedBy:    input.UpdatedBy,
		CreatedAt:    &created,
		UpdatedAt:    &updated,
		Shuffle:      input.Shuffle,
		DisplayHints: input.DisplayHints,
		ShowAnswer:   input.ShowAnswer,
		ShowResult:   input.ShowResult,
	}
	return &responseModel, nil
}

func UpdateExamConfiguration(ctx context.Context, input *model.ExamConfigurationInput) (*model.ExamConfiguration, error) {
	if input.ID == nil {
		return nil, fmt.Errorf("exam schedule id not found")
	}
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	global.CassSessioQBank = session
	cassandraQuestionBank := qbankz.ExamConfig{
		ID: *input.ID,
	}
	banks := []qbankz.ExamConfig{}
	getQuery := global.CassSessioQBank.Query(qbankz.ExamConfigTable.Get()).BindMap(qb.M{"id": cassandraQuestionBank.ID})
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
	if input.IsActive != nil {
		cassandraQuestionBank.IsActive = *input.IsActive
		updatedCols = append(updatedCols, "is_active")
	}
	if input.Shuffle != nil {
		cassandraQuestionBank.Shuffle = *input.Shuffle
		updatedCols = append(updatedCols, "shuffle_questions")
	}
	if input.DisplayHints != nil {
		cassandraQuestionBank.DisplayHints = *input.DisplayHints
		updatedCols = append(updatedCols, "display_hints")
	}
	if input.ShowAnswer != nil {
		cassandraQuestionBank.ShowAnswer = *input.ShowAnswer
		updatedCols = append(updatedCols, "show_answer")
	}
	if input.ShowResult != nil {
		cassandraQuestionBank.ShowResult = *input.ShowResult
		updatedCols = append(updatedCols, "show_result")
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
	upStms, uNames := qbankz.ExamConfigTable.Update(updatedCols...)
	updateQuery := global.CassSessioQBank.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.ExamConfiguration{
		ID:           &cassandraQuestionBank.ID,
		ExamID:       input.ExamID,
		IsActive:     input.IsActive,
		CreatedBy:    input.CreatedBy,
		UpdatedBy:    input.UpdatedBy,
		CreatedAt:    &created,
		UpdatedAt:    &updated,
		Shuffle:      input.Shuffle,
		DisplayHints: input.DisplayHints,
		ShowAnswer:   input.ShowAnswer,
		ShowResult:   input.ShowResult,
	}
	return &responseModel, nil
}
