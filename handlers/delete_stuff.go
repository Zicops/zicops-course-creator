package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/contracts/qbankz"
	"github.com/zicops/zicops-cass-pool/cassandra"
	"github.com/zicops/zicops-course-creator/helpers"
)

// in memory go cache
var (
	CourseCache = cache.New(5*time.Minute, 10*time.Minute)
)

func GetCourse(session *gocqlx.Session, id string, lspId string) (*coursez.Course, error) {
	// look up the course in memory cache
	course, found := CourseCache.Get(id)
	if found {
		return course.(*coursez.Course), nil
	}
	courses := []coursez.Course{}
	getQuery := session.Query(coursez.CourseTable.Get()).BindMap(qb.M{"id": id, "lsp_id": lspId, "is_active": true})
	if err := getQuery.SelectRelease(&courses); err != nil {
		return nil, err
	}
	if len(courses) < 1 {
		return nil, fmt.Errorf("course not found")
	}
	return &courses[0], nil
}

func GetModule(session *gocqlx.Session, id string, lspId string) (*coursez.Module, error) {
	// look up the course in memory cache
	course, found := CourseCache.Get(id)
	if found {
		return course.(*coursez.Module), nil
	}
	courses := []coursez.Module{}
	getQuery := session.Query(coursez.ModuleTable.Get()).BindMap(qb.M{"id": id, "lsp_id": lspId, "is_active": true})
	if err := getQuery.SelectRelease(&courses); err != nil {
		return nil, err
	}
	if len(courses) < 1 {
		return nil, fmt.Errorf("module not found")
	}
	return &courses[0], nil
}

func GetChapter(session *gocqlx.Session, id string, lspId string) (*coursez.Chapter, error) {
	// look up the course in memory cache
	course, found := CourseCache.Get(id)
	if found {
		return course.(*coursez.Chapter), nil
	}
	courses := []coursez.Chapter{}
	getQuery := session.Query(coursez.ChapterTable.Get()).BindMap(qb.M{"id": id, "lsp_id": lspId, "is_active": true})
	if err := getQuery.SelectRelease(&courses); err != nil {
		return nil, err
	}
	if len(courses) < 1 {
		return nil, fmt.Errorf("chapter not found")
	}
	return &courses[0], nil
}

func GetTopic(session *gocqlx.Session, id string, lspId string) (*coursez.Topic, error) {
	// look up the course in memory cache
	course, found := CourseCache.Get(id)
	if found {
		return course.(*coursez.Topic), nil
	}
	courses := []coursez.Topic{}
	getQuery := session.Query(coursez.TopicTable.Get()).BindMap(qb.M{"id": id, "lsp_id": lspId, "is_active": true})
	if err := getQuery.SelectRelease(&courses); err != nil {
		return nil, err
	}
	if len(courses) < 1 {
		return nil, fmt.Errorf("topic not found")
	}
	return &courses[0], nil
}

func GetTopicContent(session *gocqlx.Session, id string, lspId string) (*coursez.TopicContent, error) {
	// look up the course in memory cache
	course, found := CourseCache.Get(id)
	if found {
		return course.(*coursez.TopicContent), nil
	}
	courses := []coursez.TopicContent{}
	getQuery := session.Query(coursez.TopicContentTable.Get()).BindMap(qb.M{"id": id, "lsp_id": lspId, "is_active": true})
	if err := getQuery.SelectRelease(&courses); err != nil {
		return nil, err
	}
	if len(courses) < 1 {
		return nil, fmt.Errorf("topic content not found")
	}
	return &courses[0], nil
}

func GetTopicExam(session *gocqlx.Session, id string, lspId string) (*coursez.TopicExam, error) {
	// look up the course in memory cache
	course, found := CourseCache.Get(id)
	if found {
		return course.(*coursez.TopicExam), nil
	}
	courses := []coursez.TopicExam{}
	getQuery := session.Query(coursez.TopicExamTable.Get()).BindMap(qb.M{"id": id, "lsp_id": lspId, "is_active": true})
	if err := getQuery.SelectRelease(&courses); err != nil {
		return nil, err
	}
	if len(courses) < 1 {
		return nil, fmt.Errorf("topic exam not found")
	}
	return &courses[0], nil
}

func GetQuiz(session *gocqlx.Session, id string, lspId string) (*coursez.Quiz, error) {
	// look up the course in memory cache
	course, found := CourseCache.Get(id)
	if found {
		return course.(*coursez.Quiz), nil
	}
	courses := []coursez.Quiz{}
	getQuery := session.Query(coursez.QuizTable.Get()).BindMap(qb.M{"id": id, "lsp_id": lspId, "is_active": true})
	if err := getQuery.SelectRelease(&courses); err != nil {
		return nil, err
	}
	if len(courses) < 1 {
		return nil, fmt.Errorf("quiz not found")
	}
	return &courses[0], nil
}

func GetResource(session *gocqlx.Session, id string, lspId string) (*coursez.Resource, error) {
	// look up the course in memory cache
	course, found := CourseCache.Get(id)
	if found {
		return course.(*coursez.Resource), nil
	}
	courses := []coursez.Resource{}
	getQuery := session.Query(coursez.ResourceTable.Get()).BindMap(qb.M{"id": id, "lsp_id": lspId, "is_active": true})
	if err := getQuery.SelectRelease(&courses); err != nil {
		return nil, err
	}
	if len(courses) < 1 {
		return nil, fmt.Errorf("resource not found")
	}
	return &courses[0], nil
}

func GetQuestionBank(session *gocqlx.Session, id string, lspId string) (*qbankz.QuestionBankMain, error) {
	// look up the course in memory cache
	course, found := CourseCache.Get(id)
	if found {
		return course.(*qbankz.QuestionBankMain), nil
	}
	courses := []qbankz.QuestionBankMain{}
	getQuery := session.Query(qbankz.QuestionBankMainTable.Get()).BindMap(qb.M{"id": id, "lsp_id": lspId, "is_active": true})
	if err := getQuery.SelectRelease(&courses); err != nil {
		return nil, err
	}
	if len(courses) < 1 {
		return nil, fmt.Errorf("question bank not found")
	}
	return &courses[0], nil
}

func DeleteCourse(ctx context.Context, id *string) (*bool, error) {
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
	cassandraCourse, err := GetCourse(CassSession, *id, lspId)
	if err != nil {
		return &isSuccess, err
	}
	deleteSrt := fmt.Sprintf("DELETE FROM coursez.course WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
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
	cassandraModule, err := GetModule(CassSession, *id, lspId)
	if err != nil {
		return &isSuccess, err
	}
	cassandraCourse, err := GetCourse(CassSession, cassandraModule.CourseID, lspId)
	if err != nil {
		return &isSuccess, err
	}
	deleteSrt := fmt.Sprintf("DELETE FROM coursez.module WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	if strings.ToLower(cassandraCourse.Status) == "published" {
		cassandraModule.IsActive = false
		insertQuery := CassSession.Query(coursez.ModuleTable.Insert()).BindStruct(&cassandraModule)
		if err := insertQuery.ExecRelease(); err != nil {
			return &isSuccess, err
		}
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
	deleteSrt := fmt.Sprintf("DELETE FROM coursez.chapter WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	chapter, err := GetChapter(CassSession, *id, lspId)
	if err != nil {
		return &isSuccess, err
	}
	course, err := GetCourse(CassSession, chapter.CourseID, lspId)
	if err != nil {
		return &isSuccess, err
	}
	if strings.ToLower(course.Status) == "published" {
		chapter.IsActive = false
		insertQuery := CassSession.Query(coursez.ChapterTable.Insert()).BindStruct(&chapter)
		if err := insertQuery.ExecRelease(); err != nil {
			return &isSuccess, err
		}
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
	topic, err := GetTopic(CassSession, *id, lspId)
	if err != nil {
		return &isSuccess, err
	}
	course, err := GetCourse(CassSession, topic.CourseID, lspId)
	if err != nil {
		return &isSuccess, err
	}
	deleteSrt := fmt.Sprintf("DELETE FROM coursez.topic WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	if strings.ToLower(course.Status) == "published" {
		topic.IsActive = false
		insertQuery := CassSession.Query(coursez.TopicTable.Insert()).BindStruct(&topic)
		if err := insertQuery.ExecRelease(); err != nil {
			return &isSuccess, err
		}
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
	topicContent, err := GetTopicContent(CassSession, *id, lspId)
	if err != nil {
		return &isSuccess, err
	}
	course, err := GetCourse(CassSession, topicContent.CourseId, lspId)
	if err != nil {
		return &isSuccess, err
	}
	deleteSrt := fmt.Sprintf("DELETE FROM coursez.topic_content WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	if strings.ToLower(course.Status) == "published" {
		topicContent.IsActive = false
		insertQuery := CassSession.Query(coursez.TopicContentTable.Insert()).BindStruct(&topicContent)
		if err := insertQuery.ExecRelease(); err != nil {
			return &isSuccess, err
		}
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
	topicExam, err := GetTopicExam(CassSession, *id, lspId)
	if err != nil {
		return &isSuccess, err
	}
	course, err := GetCourse(CassSession, topicExam.CourseId, lspId)
	if err != nil {
		return &isSuccess, err
	}
	deleteSrt := fmt.Sprintf("DELETE FROM coursez.topic_exam WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	if strings.ToLower(course.Status) == "published" {
		topicExam.IsActive = false
		insertQuery := CassSession.Query(coursez.TopicExamTable.Insert()).BindStruct(&topicExam)
		if err := insertQuery.ExecRelease(); err != nil {
			return &isSuccess, err
		}
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
	quiz, err := GetQuiz(CassSession, *id, lspId)
	if err != nil {
		return &isSuccess, err
	}
	course, err := GetCourse(CassSession, quiz.CourseID, lspId)
	if err != nil {
		return &isSuccess, err
	}
	deleteSrt := fmt.Sprintf("DELETE FROM coursez.quiz WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	if strings.ToLower(course.Status) == "published" {
		quiz.IsActive = false
		insertQuery := CassSession.Query(coursez.QuizTable.Insert()).BindStruct(&quiz)
		if err := insertQuery.ExecRelease(); err != nil {
			return &isSuccess, err
		}
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
	topicResource, err := GetResource(CassSession, *id, lspId)
	if err != nil {
		return &isSuccess, err
	}
	course, err := GetCourse(CassSession, topicResource.CourseId, lspId)
	if err != nil {
		return &isSuccess, err
	}
	deleteSrt := fmt.Sprintf("DELETE FROM coursez.resource WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	if strings.ToLower(course.Status) == "published" {
		topicResource.IsActive = false
		insertQuery := CassSession.Query(coursez.ResourceTable.Insert()).BindStruct(&topicResource)
		if err := insertQuery.ExecRelease(); err != nil {
			return &isSuccess, err
		}
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
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.question_bank_main WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
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
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.question_main WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
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
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.options_main WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
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
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.question_paper_main WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
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
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.section_main WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
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
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.section_qb_mapping WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
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
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.section_fixed_questions WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
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
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.exam WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
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
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.exam_schedule WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
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
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.exam_instructions WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
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
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.exam_cohort WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
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
	deleteSrt := fmt.Sprintf("DELETE FROM qbankz.exam_config WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
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
	deleteSrt := fmt.Sprintf("DELETE FROM coursez.course_cohort_mapping WHERE id = %s AND lsp_id = '%s' AND is_active=true", *id, lspId)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}
