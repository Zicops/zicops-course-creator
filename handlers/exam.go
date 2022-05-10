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

func ExamCreate(ctx context.Context, exam *model.ExamInput) (*model.Exam, error) {
	log.Info("ExamCreate called")
	guid := xid.New()
	cassandraQuestionBank := qbankz.Exam{
		ID:           guid.String(),
		Name:         *exam.Name,
		Category:     *exam.Category,
		SubCategory:  *exam.SubCategory,
		IsActive:     *exam.IsActive,
		CreatedBy:    *exam.CreatedBy,
		UpdatedBy:    *exam.UpdatedBy,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
		Description:  *exam.Description,
		Code:         *exam.Code,
		QPID:         *exam.QpID,
		Type:         *exam.Type,
		ScheduleType: *exam.ScheduleType,
		Duration:     *exam.Duration,
		Status:       *exam.Status,
	}

	insertQuery := global.CassSession.Session.Query(qbankz.ExamTable.Insert()).BindStruct(cassandraQuestionBank)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.Exam{
		ID:           &cassandraQuestionBank.ID,
		Name:         exam.Name,
		Category:     exam.Category,
		SubCategory:  exam.SubCategory,
		IsActive:     exam.IsActive,
		CreatedBy:    exam.CreatedBy,
		UpdatedBy:    exam.UpdatedBy,
		CreatedAt:    &created,
		UpdatedAt:    &updated,
		Description:  exam.Description,
		Code:         exam.Code,
		QpID:         exam.QpID,
		Type:         exam.Type,
		ScheduleType: exam.ScheduleType,
		Duration:     exam.Duration,
		Status:       exam.Status,
	}
	return &responseModel, nil
}

func ExamUpdate(ctx context.Context, input *model.ExamInput) (*model.Exam, error) {
	if input.ID == nil {
		return nil, fmt.Errorf("question paper id not found")
	}
	cassandraQuestionBank := qbankz.Exam{
		ID: *input.ID,
	}
	banks := []qbankz.Exam{}
	getQuery := global.CassSession.Session.Query(qbankz.ExamTable.Get()).BindMap(qb.M{"id": cassandraQuestionBank.ID})
	if err := getQuery.SelectRelease(&banks); err != nil {
		return nil, err
	}
	if len(banks) == 0 {
		return nil, fmt.Errorf("exams not found")
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
	if input.Description != nil {
		cassandraQuestionBank.Description = *input.Description
		updatedCols = append(updatedCols, "description")
	}
	if input.Code != nil {
		cassandraQuestionBank.Code = *input.Code
		updatedCols = append(updatedCols, "code")
	}
	if input.QpID != nil {
		cassandraQuestionBank.QPID = *input.QpID
		updatedCols = append(updatedCols, "qp_id")
	}
	if input.Type != nil {
		cassandraQuestionBank.Type = *input.Type
		updatedCols = append(updatedCols, "type")
	}
	if input.ScheduleType != nil {
		cassandraQuestionBank.ScheduleType = *input.ScheduleType
		updatedCols = append(updatedCols, "schedule_type")
	}
	if input.Duration != nil {
		cassandraQuestionBank.Duration = *input.Duration
		updatedCols = append(updatedCols, "duration")
	}
	if input.Status != nil {
		cassandraQuestionBank.Status = *input.Status
		updatedCols = append(updatedCols, "status")
	}
	updatedAt := time.Now().Unix()
	cassandraQuestionBank.UpdatedAt = updatedAt
	updatedCols = append(updatedCols, "updated_at")
	if len(updatedCols) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}
	upStms, uNames := qbankz.ExamTable.Update(updatedCols...)
	updateQuery := global.CassSession.Session.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.Exam{
		ID:           &cassandraQuestionBank.ID,
		Name:         input.Name,
		Category:     input.Category,
		SubCategory:  input.SubCategory,
		IsActive:     input.IsActive,
		CreatedBy:    input.CreatedBy,
		UpdatedBy:    input.UpdatedBy,
		CreatedAt:    &created,
		UpdatedAt:    &updated,
		Description:  input.Description,
		Code:         input.Code,
		QpID:         input.QpID,
		Type:         input.Type,
		ScheduleType: input.ScheduleType,
		Duration:     input.Duration,
		Status:       input.Status,
	}
	return &responseModel, nil
}
