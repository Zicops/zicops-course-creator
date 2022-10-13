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
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/helpers"
)

func ExamScheduleCreate(ctx context.Context, exam *model.ExamScheduleInput) (*model.ExamSchedule, error) {
	log.Info("ExamCreate called")
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
	guid := xid.New()
	// prase *exam.Start to int64
	startString := strconv.Itoa(*exam.Start)
	start, err := strconv.ParseInt(startString, 10, 64)
	if err != nil {
		return nil, err
	}
	// prase *exam.End to int64
	endString := strconv.Itoa(*exam.End)
	end, err := strconv.ParseInt(endString, 10, 64)
	if err != nil {
		return nil, err
	}
	cassandraQuestionBank := qbankz.ExamSchedule{
		ID:         guid.String(),
		LspId:      lspID,
		ExamID:     *exam.ExamID,
		IsActive:   *exam.IsActive,
		CreatedBy:  email_creator,
		UpdatedBy:  email_creator,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
		Start:      start,
		End:        end,
		BufferTime: *exam.BufferTime,
	}

	insertQuery := CassSession.Query(qbankz.ExamScheduleTable.Insert()).BindStruct(cassandraQuestionBank)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.ExamSchedule{
		ID:         &cassandraQuestionBank.ID,
		ExamID:     exam.ExamID,
		IsActive:   exam.IsActive,
		CreatedBy:  exam.CreatedBy,
		UpdatedBy:  exam.UpdatedBy,
		CreatedAt:  &created,
		UpdatedAt:  &updated,
		Start:      exam.Start,
		End:        exam.End,
		BufferTime: exam.BufferTime,
	}
	return &responseModel, nil
}

func ExamScheduleUpdate(ctx context.Context, input *model.ExamScheduleInput) (*model.ExamSchedule, error) {
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
	cassandraQuestionBank := qbankz.ExamSchedule{
		ID: *input.ID,
	}
	banks := []qbankz.ExamSchedule{}
	getQuery := CassSession.Query(qbankz.ExamScheduleTable.Get()).BindMap(qb.M{"id": cassandraQuestionBank.ID, "lsp_id": lspID, "is_active": true})
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
	if email_creator != "" {
		cassandraQuestionBank.UpdatedBy = email_creator
		updatedCols = append(updatedCols, "updated_by")
	}
	if input.Start != nil {
		startString := strconv.Itoa(*input.Start)
		start, err := strconv.ParseInt(startString, 10, 64)
		if err != nil {
			return nil, err
		}
		cassandraQuestionBank.Start = start
		updatedCols = append(updatedCols, "start")
	}
	if input.End != nil {
		endString := strconv.Itoa(*input.End)
		end, err := strconv.ParseInt(endString, 10, 64)
		if err != nil {
			return nil, err
		}
		cassandraQuestionBank.End = end
		updatedCols = append(updatedCols, "end")
	}
	if input.BufferTime != nil {
		cassandraQuestionBank.BufferTime = *input.BufferTime
		updatedCols = append(updatedCols, "buffer_time")
	}
	updatedAt := time.Now().Unix()
	cassandraQuestionBank.UpdatedAt = updatedAt
	updatedCols = append(updatedCols, "updated_at")
	if len(updatedCols) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}
	upStms, uNames := qbankz.ExamScheduleTable.Update(updatedCols...)
	updateQuery := CassSession.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraQuestionBank.CreatedAt, 10)
	updated := strconv.FormatInt(cassandraQuestionBank.UpdatedAt, 10)
	responseModel := model.ExamSchedule{
		ID:         &cassandraQuestionBank.ID,
		ExamID:     input.ExamID,
		IsActive:   input.IsActive,
		CreatedBy:  input.CreatedBy,
		UpdatedBy:  input.UpdatedBy,
		CreatedAt:  &created,
		UpdatedAt:  &updated,
		Start:      input.Start,
		End:        input.End,
		BufferTime: input.BufferTime,
	}
	return &responseModel, nil
}
