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
		ID:         uuid.New().String(),
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
	cassandraQuestionBank := *GetExamSchedule(ctx, *input.ID, lspID, CassSession)
	updatedCols := []string{}
	if input.ExamID != nil && cassandraQuestionBank.ExamID != *input.ExamID {
		cassandraQuestionBank.ExamID = *input.ExamID
		updatedCols = append(updatedCols, "exam_id")
	}
	if email_creator != "" && cassandraQuestionBank.UpdatedBy != email_creator {
		cassandraQuestionBank.UpdatedBy = email_creator
		updatedCols = append(updatedCols, "updated_by")
	}
	if input.Start != nil && int(cassandraQuestionBank.Start) != *input.Start {
		startString := strconv.Itoa(*input.Start)
		start, err := strconv.ParseInt(startString, 10, 64)
		if err != nil {
			return nil, err
		}
		cassandraQuestionBank.Start = start
		updatedCols = append(updatedCols, "start")
	}
	if input.End != nil && int(cassandraQuestionBank.End) != *input.End {
		endString := strconv.Itoa(*input.End)
		end, err := strconv.ParseInt(endString, 10, 64)
		if err != nil {
			return nil, err
		}
		cassandraQuestionBank.End = end
		updatedCols = append(updatedCols, "end")
	}
	if input.BufferTime != nil && cassandraQuestionBank.BufferTime != *input.BufferTime {
		cassandraQuestionBank.BufferTime = *input.BufferTime
		updatedCols = append(updatedCols, "buffer_time")
	}
	if len(updatedCols) > 0 {
		updatedAt := time.Now().Unix()
		cassandraQuestionBank.UpdatedAt = updatedAt
		updatedCols = append(updatedCols, "updated_at")
		upStms, uNames := qbankz.ExamScheduleTable.Update(updatedCols...)
		updateQuery := CassSession.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
		if err := updateQuery.ExecRelease(); err != nil {
			return nil, err
		}
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

func GetExamSchedule(ctx context.Context, courseID string, lspID string, session *gocqlx.Session) *qbankz.ExamSchedule {
	chapters := []qbankz.ExamSchedule{}
	getQueryStr := fmt.Sprintf("SELECT * FROM qbankz.exam_schedule WHERE id='%s' and lsp_id='%s' and is_active=true", courseID, lspID)
	getQuery := session.Query(getQueryStr, nil)
	if err := getQuery.SelectRelease(&chapters); err != nil {
		return nil
	}
	return &chapters[0]
}
