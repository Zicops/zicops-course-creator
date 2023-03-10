package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/helpers"
	"github.com/zicops/zicops-course-creator/lib/db/bucket"
	"github.com/zicops/zicops-course-creator/lib/googleprojectlib"
)

func DeleteCatMain(ctx context.Context, id *string) (*bool, error) {
	session, err := global.CassPool.GetSession(ctx, "coursez")
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
	session, err := global.CassPool.GetSession(ctx, "coursez")
	if err != nil {
		return nil, err
	}
	isSuccess := false
	CassSession := session
	_, err = helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	deleteSrt := fmt.Sprintf("DELETE FROM coursez.sub_cat_main WHERE id='%s'AND parent_id='%s' AND is_active=true", *id, *parent_id)
	if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
		return &isSuccess, err
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteCourse(ctx context.Context, id *string) (*bool, error) {
	session, err := global.CassPool.GetSession(ctx, "coursez")
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
	session, err := global.CassPool.GetSession(ctx, "coursez")
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
	session, err := global.CassPool.GetSession(ctx, "coursez")
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
	session, err := global.CassPool.GetSession(ctx, "coursez")
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
	session, err := global.CassPool.GetSession(ctx, "coursez")
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
	cassandraTopicContent := GetTopicContent(ctx, *id, lspId, CassSession)
	if cassandraTopicContent != nil {
		deleteSrt := fmt.Sprintf("DELETE FROM coursez.topic_content WHERE id = '%s' AND lsp_id = '%s' AND is_active=true", *id, lspId)
		if err := CassSession.Query(deleteSrt, nil).Exec(); err != nil {
			return &isSuccess, err
		}
		mod := GetModule(ctx, cassandraTopicContent.ModuleId, cassandraTopicContent.LspId, CassSession)
		newDuration := mod.Duration - cassandraTopicContent.Duration
		queryStr := fmt.Sprintf("UPDATE coursez.module SET duration=%d WHERE id='%s'and lsp_id='%s' and is_active=true and created_at=%d", newDuration, cassandraTopicContent.ModuleId, cassandraTopicContent.LspId, mod.CreatedAt)
		updateQ := CassSession.Query(queryStr, nil)
		if err := updateQ.ExecRelease(); err != nil {
			return nil, err
		}

		course := GetCourse(ctx, cassandraTopicContent.CourseId, cassandraTopicContent.LspId, CassSession)
		newDuration = course.Duration - cassandraTopicContent.Duration
		queryStr = fmt.Sprintf("UPDATE coursez.course SET duration=%d WHERE id='%s' and lsp_id='%s' and is_active=true and created_at=%d", newDuration, cassandraTopicContent.CourseId, cassandraTopicContent.LspId, course.CreatedAt)
		updateQ = CassSession.Query(queryStr, nil)
		if err := updateQ.ExecRelease(); err != nil {
			return nil, err
		}
	}
	isSuccess = true
	return &isSuccess, nil
}

func DeleteTopicExam(ctx context.Context, id *string) (*bool, error) {
	session, err := global.CassPool.GetSession(ctx, "coursez")
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
	session, err := global.CassPool.GetSession(ctx, "coursez")
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
	session, err := global.CassPool.GetSession(ctx, "coursez")
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
	session, err := global.CassPool.GetSession(ctx, "qbankz")
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
	session, err := global.CassPool.GetSession(ctx, "qbankz")
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
	session, err := global.CassPool.GetSession(ctx, "qbankz")
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
	session, err := global.CassPool.GetSession(ctx, "qbankz")
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
	session, err := global.CassPool.GetSession(ctx, "qbankz")
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
	session, err := global.CassPool.GetSession(ctx, "qbankz")
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
	session, err := global.CassPool.GetSession(ctx, "qbankz")
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
	session, err := global.CassPool.GetSession(ctx, "qbankz")
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
	session, err := global.CassPool.GetSession(ctx, "qbankz")
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
	session, err := global.CassPool.GetSession(ctx, "qbankz")
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
	session, err := global.CassPool.GetSession(ctx, "qbankz")
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
	session, err := global.CassPool.GetSession(ctx, "qbankz")
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
	session, err := global.CassPool.GetSession(ctx, "coursez")
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

func DeleteTopicSubtitle(ctx context.Context, courseID string, topicID string, fileName string, lang *string) (*bool, error) {
	log.Info("DeleteTopicSubtitle called")
	resp := false
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return &resp, err
	}
	session, err := global.CassPool.GetSession(ctx, "coursez")
	if err != nil {
		log.Println("Got error while getting session: ", err)
		return nil, err
	}
	CassSession := session
	lspId := claims["lsp_id"].(string)
	if lspId == "" {
		return &resp, fmt.Errorf("lsp_id is empty")
	}
	// convet any url encoded characters to normal characters
	fileName, err = url.QueryUnescape(fileName)
	if err != nil {
		log.Errorf("Failed to delete subtitle to course topic: %v", err.Error())
		return &resp, err
	}
	storageC := bucket.NewStorageHandler()
	gproject := googleprojectlib.GetGoogleProjectID()
	err = storageC.InitializeStorageClient(ctx, gproject, lspId)
	if err != nil {
		log.Errorf("Failed to delete subtitle to course topic: %v", err.Error())
		return &resp, err
	}
	mainBucket := courseID + "/" + topicID + "/subtitles/"
	bucketPath := mainBucket + fileName
	res := storageC.DeleteObjectsFromBucket(ctx, bucketPath)
	anyitemsLeft := storageC.CheckIfBucketHasItems(ctx, mainBucket)
	r := true
	if res == "1" {
		return &r, nil
	}
	r = false

	if !anyitemsLeft {
		qryStr := fmt.Sprintf(`SELECT * FROM coursez.topic_content WHERE courseid = '%s' AND topicid='%s' ALLOW FILTERING`, courseID, topicID)
		getTopic := func() (topics []coursez.TopicContent, err error) {
			q := CassSession.Query(qryStr, nil)
			defer q.Release()
			iter := q.Iter()
			return topics, iter.Select(&topics)
		}

		topics, err := getTopic()
		if err != nil {
			return nil, err
		}
		if len(topics) == 0 {
			return nil, nil
		}
		topic := topics[0]

		//remove the subtitle bucket
		if topic.SubtitleFileBucket != "" {
			topic.SubtitleFileBucket = ""
		}
		//update the table
		upStms, uNames := coursez.TopicContentTable.Update("subtitlefilebucket")
		updateQuery := CassSession.Query(upStms, uNames).BindStruct(&topic)
		if err := updateQuery.ExecRelease(); err != nil {
			return nil, err
		}
	}

	return &r, errors.New(res)
}

func DeleteCourseMedia(ctx context.Context, courseID string, fileName string) (*bool, error) {
	resp := false
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return &resp, err
	}
	lspId := claims["lsp_id"].(string)
	if lspId == "" {
		return &resp, fmt.Errorf("lsp_id is empty")
	}
	storageC := bucket.NewStorageHandler()
	gproject := googleprojectlib.GetGoogleProjectID()
	err = storageC.InitializeStorageClient(ctx, gproject, lspId)
	if err != nil {
		log.Errorf("Failed to delete subtitle to course topic: %v", err.Error())
		return &resp, err
	}
	fileName, err = url.QueryUnescape(fileName)
	if err != nil {
		log.Errorf("Failed to delete subtitle to course topic: %v", err.Error())
		return &resp, err
	}

	bucketPath := courseID + "/" + fileName
	res := storageC.DeleteObjectsFromBucket(ctx, bucketPath)

	r := true
	if res == "1" {
		return &r, nil
	}
	r = false
	return &r, errors.New(res)

}
