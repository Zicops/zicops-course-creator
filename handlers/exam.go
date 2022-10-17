package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/scylladb/gocqlx/v2"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/qbankz"
	"github.com/zicops/zicops-cass-pool/cassandra"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/helpers"
)

func ExamCreate(ctx context.Context, exam *model.ExamInput) (*model.Exam, error) {
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
	qpId := *exam.QpID
	questionsIDs, err := GetQuestionIDsFromPaperId(CassSession, ctx, lspID, qpId)
	if err != nil {
		return nil, err
	}
	words := []string{}
	if exam.Name != nil {
		name := strings.ToLower(*exam.Name)
		wordsLocal := strings.Split(name, " ")
		words = append(words, wordsLocal...)
	}
	cassandraQuestionBank := qbankz.Exam{
		ID:           guid.String(),
		Name:         *exam.Name,
		Words:        words,
		Category:     *exam.Category,
		SubCategory:  *exam.SubCategory,
		IsActive:     *exam.IsActive,
		CreatedBy:    email_creator,
		UpdatedBy:    email_creator,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
		Description:  *exam.Description,
		Code:         *exam.Code,
		QPID:         *exam.QpID,
		Type:         *exam.Type,
		ScheduleType: *exam.ScheduleType,
		Duration:     *exam.Duration,
		Status:       *exam.Status,
		LSPID:        lspID,
		QuestionIDs:  questionsIDs,
	}

	insertQuery := CassSession.Query(qbankz.ExamTable.Insert()).BindStruct(cassandraQuestionBank)
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
	cassandraQuestionBank := *GetExam(ctx, *input.ID, lspID, CassSession)
	updatedCols := []string{}
	if input.Name != nil && *input.Name != cassandraQuestionBank.Name {
		cassandraQuestionBank.Name = *input.Name
		words := []string{}
		if input.Name != nil {
			name := strings.ToLower(*input.Name)
			wordsLocal := strings.Split(name, " ")
			words = append(words, wordsLocal...)
		}
		cassandraQuestionBank.Words = words
		updatedCols = append(updatedCols, "words")
		updatedCols = append(updatedCols, "name")
	}
	if input.Category != nil && cassandraQuestionBank.Category != *input.Category {
		cassandraQuestionBank.Category = *input.Category
		updatedCols = append(updatedCols, "category")
	}
	if input.SubCategory != nil && cassandraQuestionBank.SubCategory != *input.SubCategory {
		cassandraQuestionBank.SubCategory = *input.SubCategory
		updatedCols = append(updatedCols, "sub_category")
	}
	if email_creator != "" && cassandraQuestionBank.UpdatedBy != email_creator {
		cassandraQuestionBank.UpdatedBy = email_creator
		updatedCols = append(updatedCols, "updated_by")
	}
	if input.Description != nil && cassandraQuestionBank.Description != *input.Description {
		cassandraQuestionBank.Description = *input.Description
		updatedCols = append(updatedCols, "description")
	}
	if input.Code != nil && cassandraQuestionBank.Code != *input.Code {
		cassandraQuestionBank.Code = *input.Code
		updatedCols = append(updatedCols, "code")
	}
	if input.QpID != nil && cassandraQuestionBank.QPID != *input.QpID {
		cassandraQuestionBank.QPID = *input.QpID
		updatedCols = append(updatedCols, "qp_id")
	}
	if input.Type != nil && cassandraQuestionBank.Type != *input.Type {
		cassandraQuestionBank.Type = *input.Type
		updatedCols = append(updatedCols, "type")
	}
	if input.ScheduleType != nil && cassandraQuestionBank.ScheduleType != *input.ScheduleType {
		cassandraQuestionBank.ScheduleType = *input.ScheduleType
		updatedCols = append(updatedCols, "schedule_type")
	}
	if input.Duration != nil && cassandraQuestionBank.Duration != *input.Duration {
		cassandraQuestionBank.Duration = *input.Duration
		updatedCols = append(updatedCols, "duration")
	}
	if input.Status != nil && cassandraQuestionBank.Status != *input.Status {
		cassandraQuestionBank.Status = *input.Status
		updatedCols = append(updatedCols, "status")
	}
	updatedAt := time.Now().Unix()
	if input.QpID != nil && *input.QpID != "" && *input.QpID != cassandraQuestionBank.QPID {
		questionsIDs, err := GetQuestionIDsFromPaperId(CassSession, ctx, lspID, *input.QpID)
		if err != nil {
			return nil, err
		}
		cassandraQuestionBank.QuestionIDs = questionsIDs
		updatedCols = append(updatedCols, "question_ids")
	}
	if len(updatedCols) > 0 {
		cassandraQuestionBank.UpdatedAt = updatedAt
		updatedCols = append(updatedCols, "updated_at")
		upStms, uNames := qbankz.ExamTable.Update(updatedCols...)
		updateQuery := CassSession.Query(upStms, uNames).BindStruct(&cassandraQuestionBank)
		if err := updateQuery.ExecRelease(); err != nil {
			return nil, err
		}
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

func GetQuestionIDsFromPaperId(session *gocqlx.Session, ctx context.Context, lspID string, qpID string) ([]string, error) {
	qryStr := fmt.Sprintf(`SELECT * from qbankz.section_main where lsp_id='%s' AND is_active=true  AND qp_id='%s' ALLOW FILTERING`, lspID, qpID)
	getSectionsMap := func() (banks []qbankz.SectionMain, err error) {
		q := session.Query(qryStr, nil)
		defer q.Release()
		iter := q.Iter()
		return banks, iter.Select(&banks)
	}
	sectionsMap, err := getSectionsMap()
	if err != nil {
		return nil, err
	}
	questionsIDs := []string{}
	for _, section := range sectionsMap {
		currentSectionId := section.ID
		sectionQbQuery := fmt.Sprintf(`SELECT * from qbankz.section_qb_mapping where lsp_id='%s' AND is_active=true  AND section_id='%s' ALLOW FILTERING`, lspID, currentSectionId)
		getSections := func() (banks []qbankz.SectionQBMapping, err error) {
			q := session.Query(sectionQbQuery, nil)
			defer q.Release()
			iter := q.Iter()
			return banks, iter.Select(&banks)
		}
		sections, err := getSections()
		if err != nil {
			return nil, err
		}
		for _, section := range sections {
			sectionID := section.ID
			qryStr := fmt.Sprintf(`SELECT * from qbankz.section_fixed_questions where sqb_id='%s' AND lsp_id='%s' AND is_active=true ALLOW FILTERING`, sectionID, lspID)
			getQuestions := func() (banks []qbankz.SectionFixedQuestions, err error) {
				q := session.Query(qryStr, nil)
				defer q.Release()
				iter := q.Iter()
				return banks, iter.Select(&banks)
			}
			questions, err := getQuestions()
			if err != nil {
				return nil, err
			}
			for _, question := range questions {
				questionsIDs = append(questionsIDs, question.QuestionID)
			}
		}
	}
	return questionsIDs, nil
}

func GetExam(ctx context.Context, courseID string, lspID string, session *gocqlx.Session) *qbankz.Exam {
	chapters := []qbankz.Exam{}
	getQueryStr := fmt.Sprintf("SELECT * FROM qbankz.exam WHERE id='%s' and lsp_id='%s' and is_active=true", courseID, lspID)
	getQuery := session.Query(getQueryStr, nil)
	if err := getQuery.SelectRelease(&chapters); err != nil {
		return nil
	}
	return &chapters[0]
}
