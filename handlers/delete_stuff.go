package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-cass-pool/cassandra"
	"github.com/zicops/zicops-course-creator/helpers"
)

func DeleteCatMain(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	_, err = helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	deleteSrt := fmt.Sprintf("DELETE FROM coursez.cat_main WHERE id='%s' AND is_active=true", *id)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil

}

func DeleteSubCatMain(ctx context.Context, id *string, parent_id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	_, err = helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	deleteSrt := fmt.Sprintf("DELETE FROM coursez.sub_cat_main WHERE id='%s' AND is_active=true AND parent_id='%s' AND is_active=true", *id, *parent_id)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteCourse(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("coursez")
	if err != nil || id == nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	cassandraCourse := GetCourse(ctx, *id, lspId, CassSession)
	deleteSrt := fmt.Sprintf("DELETE FROM coursez.course WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	if strings.ToLower(cassandraCourse.Status) == "published" {
		cassandraCourse.IsActive = false
		insertQuery := CassSession.Query(coursez.CourseTable.Insert()).BindStruct(&cassandraCourse)
		if err := insertQuery.ExecRelease(); err != nil {
			return &isSuccess, err
		}
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteCourseModule(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	deleteSrt := fmt.Sprintf("DELETE FROM coursez.module WHERE id='%s' AND lsp_id='%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteCourseChapter(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	deleteSrt := fmt.Sprintf("DELETE FROM coursez.chapter WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteCourseTopic(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	deleteSrt := fmt.Sprintf("DELETE FROM coursez.topic WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteTopicContent(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	deleteSrt := fmt.Sprintf("DELETE FROM coursez.topic_content WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteTopicExam(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	deleteSrt := fmt.Sprintf("DELETE FROM coursez.topic_exam WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteQuiz(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	deleteSrt := fmt.Sprintf("DELETE FROM coursez.quiz WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteTopicResource(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	deleteSrt := fmt.Sprintf("DELETE FROM coursez.resource WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteQuestionBank(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.question_bank_main WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteQuestionBankQuestion(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.question_main WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteQuestionOptions(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.options_main WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteQuestionPaper(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.question_paper_main WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteQuestionPaperSection(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.section_main WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteSectionToBank(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.section_qb_mapping WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteSectionFixedQuestions(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.section_fixed_questions WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteExam(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.exam WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteExamSchedule(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.exam_schedule WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteExamInstruction(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.exam_instructions WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteExamCohort(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.exam_cohort WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteExamConfiguration(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("qbankz")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.exam_config WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteCourseCohort(ctx context.Context, id *string) (*bool, error) {
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspId := claims["lsp_id"].(string)
	deleteSrt := fmt.Sprintf("DELETE FROM coursez.course_cohort_mapping WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}
