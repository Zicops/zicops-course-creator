package handlers

import (
	"context"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/graph/model"
)

func TopicContentCreate(ctx context.Context, topicID string, topicConent *model.TopicContentInput) (*model.TopicContent, error) {
	log.Info("TopicContentCreate called")
	cassandraTopicContent := coursez.TopicContent{
		TopicId:            topicID,
		Language:           topicConent.Language,
		CreatedAt:          time.Now().Unix(),
		UpdatedAt:          time.Now().Unix(),
		StartTime:          *topicConent.StartTime,
		Duration:           *topicConent.Duration,
		SkipIntro:          *topicConent.SkipIntro,
		NextShowtime:       *topicConent.NextShowTime,
		FromEndTime:        *topicConent.FromEndTime,
		TopicContentBucket: "",
		Url:                "",
		IsDeleted:          false,
	}
	// set course in cassandra
	insertQuery := global.CassSession.Session.Query(coursez.TopicContentTable.Insert()).BindStruct(cassandraTopicContent)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraTopicContent.CreatedAt, 10)
	responseModel := model.TopicContent{
		Language:     topicConent.Language,
		StartTime:    topicConent.StartTime,
		CreatedAt:    &created,
		UpdatedAt:    &created,
		Duration:     topicConent.Duration,
		SkipIntro:    topicConent.SkipIntro,
		NextShowTime: topicConent.NextShowTime,
		FromEndTime:  topicConent.FromEndTime,
		TopicID:      topicID,
	}
	return &responseModel, nil
}
