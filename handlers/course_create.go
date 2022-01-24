package handlers

import (
	"context"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/rs/xid"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/graph/model"
)

// func to get context and course input and set it in cassandra
func CourseCreator(ctx context.Context, courseInput *model.CourseInput) (*model.Course, error) {
	log.Info("CourseCreator called")
	// set course input in cassandra
	guid := xid.New()
	language := []string{}
	takeaways := []string{}
	prequisites := []string{}
	goodFor := []string{}
	mustFor := []string{}
	for _, lang := range courseInput.Language {
		language = append(language, *lang)
	}
	for _, take := range courseInput.Takeaways {
		takeaways = append(takeaways, *take)
	}
	for _, preq := range courseInput.Prequisites {
		prequisites = append(prequisites, *preq)
	}
	for _, good := range courseInput.GoodFor {
		goodFor = append(goodFor, *good)
	}
	for _, must := range courseInput.MustFor {
		mustFor = append(mustFor, *must)
	}

	cassandraCourse := coursez.Course{
		ID:           guid.String(),
		Name:         courseInput.Name,
		Description:  courseInput.Description,
		Instructor:   "",
		Image:        "https://storage.googleapis.com/zicops.com/school-board-ge1701ca8f_640.jpg",
		PreviewVideo: "https://storage.googleapis.com/zicops.com/school-board-ge1701ca8f_640.jpg",
		Duration:     0,
		Language:     language,
		TakeAways:    takeaways,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
		Prequisites:  prequisites,
		GoodFor:      goodFor,
		MustFor:      mustFor,
		IsDeleted:    false,
	}
	if courseInput.Type != nil {
		cassandraCourse.Type = *courseInput.Type
	}
	if courseInput.Owner != nil {
		cassandraCourse.Owner = *courseInput.Owner
	}
	if courseInput.Level != nil {
		cassandraCourse.Level = *courseInput.Level
	}
	if courseInput.CreatedBy != nil {
		cassandraCourse.CreatedBy = *courseInput.CreatedBy
	}
	if courseInput.UpdatedBy != nil {
		cassandraCourse.UpdatedBy = *courseInput.UpdatedBy
	}
	if courseInput.Status != nil {
		cassandraCourse.Status = *courseInput.Status
	}

	// set course in cassandra
	insertQuery := global.CassSession.Session.Query(coursez.CourseTable.Insert()).BindStruct(cassandraCourse)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	created := strconv.FormatInt(cassandraCourse.CreatedAt, 10)
	responseModel := model.Course{
		ID:           &cassandraCourse.ID,
		Name:         courseInput.Name,
		Description:  courseInput.Description,
		Instructor:   courseInput.Instructor,
		Image:        &cassandraCourse.Image,
		PreviewVideo: &cassandraCourse.PreviewVideo,
		Owner:        courseInput.Owner,
		Duration:     &cassandraCourse.Duration,
		Level:        courseInput.Level,
		Language:     courseInput.Language,
		Takeaways:    courseInput.Takeaways,
		CreatedAt:    &created,
		UpdatedAt:    &created,
		Type:         courseInput.Type,
		Prequisites:  courseInput.Prequisites,
		GoodFor:      courseInput.GoodFor,
		MustFor:      courseInput.MustFor,
		CreatedBy:    courseInput.CreatedBy,
		UpdatedBy:    courseInput.UpdatedBy,
		Status:       courseInput.Status,
	}
	return &responseModel, nil
}
