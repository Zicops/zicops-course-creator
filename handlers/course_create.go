package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/scylladb/gocqlx/v2/qb"
	log "github.com/sirupsen/logrus"

	"github.com/rs/xid"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/lib/db/bucket"
	"github.com/zicops/zicops-course-creator/lib/googleprojectlib"
)

// func to get context and course input and set it in cassandra
func CourseCreator(ctx context.Context, courseInput *model.CourseInput) (*model.Course, error) {
	log.Info("CourseCreator called")
	// set course input in cassandra
	guid := xid.New()
	language := []string{}
	takeaways := []string{}
	outcomes := []string{}
	prequisites := []string{}
	goodFor := []string{}
	mustFor := []string{}
	relatedSkills := []string{}
	approvers := []string{}
	subCats := make([]coursez.SubCat, 0)
	subCatsRes := make([]*model.SubCategories, 0)

	for _, lang := range courseInput.Language {
		language = append(language, *lang)
	}
	for _, take := range courseInput.Benefits {
		takeaways = append(takeaways, *take)
	}
	for _, out := range courseInput.Outcomes {
		outcomes = append(outcomes, *out)
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
	for _, relSkill := range courseInput.RelatedSkills {
		relatedSkills = append(relatedSkills, *relSkill)
	}
	for _, approver := range courseInput.Approvers {
		approvers = append(approvers, *approver)
	}
	for _, subCat := range courseInput.SubCategories {
		var subC coursez.SubCat
		var subCR model.SubCategories
		subC.Name = *subCat.Name
		subC.Rank = *subCat.Rank
		subCR.Name = subCat.Name
		subCR.Rank = subCat.Rank
		subCats = append(subCats, subC)
		subCatsRes = append(subCatsRes, &subCR)
	}
	cassandraCourse := coursez.Course{
		ID:                 guid.String(),
		Name:               *courseInput.Name,
		Description:        *courseInput.Description,
		Image:              "https://storage.googleapis.com/zicops.com/school-board-ge1701ca8f_640.jpg",
		PreviewVideo:       "https://storage.googleapis.com/zicops.com/school-board-ge1701ca8f_640.jpg",
		TileImage:          "https://storage.googleapis.com/zicops.com/school-board-ge1701ca8f_640.jpg",
		TileImageBucket:    "",
		ImageBucket:        "",
		PreviewVideoBucket: "",
		Benefits:           takeaways,
		Outcomes:           outcomes,
		Language:           language,
		CreatedAt:          time.Now().Unix(),
		UpdatedAt:          time.Now().Unix(),
		Prequisites:        prequisites,
		GoodFor:            goodFor,
		MustFor:            mustFor,
		RelatedSkills:      relatedSkills,
		Approvers:          approvers,
		Status:             courseInput.Status.String(),
		IsActive:           *courseInput.IsActive,
		SubCategories:      subCats,
	}
	if courseInput.Summary != nil {
		cassandraCourse.Summary = *courseInput.Summary
	}
	if courseInput.Instructor != nil {
		cassandraCourse.Instructor = *courseInput.Instructor
	}
	if courseInput.Type != nil {
		cassandraCourse.Type = *courseInput.Type
	}
	if courseInput.Owner != nil {
		cassandraCourse.Owner = *courseInput.Owner
	}
	if courseInput.Duration != nil {
		cassandraCourse.Duration = *courseInput.Duration
	}
	if courseInput.PublishDate != nil {
		cassandraCourse.PublishDate = *courseInput.PublishDate
	}
	if courseInput.ExpiryDate != nil {
		cassandraCourse.ExpiryDate = *courseInput.ExpiryDate
	}
	if courseInput.IsDisplay != nil {
		cassandraCourse.IsDisplay = *courseInput.IsDisplay
	}
	if courseInput.QaRequired != nil {
		cassandraCourse.QARequired = *courseInput.QaRequired
	}
	if courseInput.CreatedBy != nil {
		cassandraCourse.CreatedBy = *courseInput.CreatedBy
	}
	if courseInput.UpdatedBy != nil {
		cassandraCourse.UpdatedBy = *courseInput.UpdatedBy
	}
	if courseInput.Category != nil {
		cassandraCourse.Category = *courseInput.Category
	}
	if courseInput.ExpectedCompletion != nil {
		cassandraCourse.ExpectedCompletion = *courseInput.ExpectedCompletion
	}
	if courseInput.SubCategory != nil {
		cassandraCourse.SubCategory = *courseInput.SubCategory
	}
	// set course in cassandra
	insertQuery := global.CassSession.Session.Query(coursez.CourseTable.Insert()).BindStruct(cassandraCourse)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}
	if courseInput.IsActive == nil {
		active := false
		courseInput.IsActive = &active
	}
	created := strconv.FormatInt(cassandraCourse.CreatedAt, 10)
	responseModel := model.Course{
		ID:                 &cassandraCourse.ID,
		Name:               courseInput.Name,
		Description:        courseInput.Description,
		Summary:            courseInput.Summary,
		Instructor:         courseInput.Instructor,
		Image:              &cassandraCourse.Image,
		PreviewVideo:       &cassandraCourse.PreviewVideo,
		Owner:              courseInput.Owner,
		Duration:           &cassandraCourse.Duration,
		ExpertiseLevel:     courseInput.ExpertiseLevel,
		Language:           courseInput.Language,
		Benefits:           courseInput.Benefits,
		CreatedAt:          &created,
		UpdatedAt:          &created,
		Type:               courseInput.Type,
		Prequisites:        courseInput.Prequisites,
		GoodFor:            courseInput.GoodFor,
		MustFor:            courseInput.MustFor,
		RelatedSkills:      courseInput.RelatedSkills,
		PublishDate:        courseInput.PublishDate,
		ExpiryDate:         courseInput.ExpiryDate,
		QaRequired:         courseInput.QaRequired,
		Approvers:          courseInput.Approvers,
		CreatedBy:          courseInput.CreatedBy,
		UpdatedBy:          courseInput.UpdatedBy,
		Status:             courseInput.Status,
		IsDisplay:          courseInput.IsDisplay,
		ExpectedCompletion: courseInput.ExpectedCompletion,
		Category:           courseInput.Category,
		SubCategory:        courseInput.SubCategory,
		SubCategories:      subCatsRes,
		Outcomes:           courseInput.Outcomes,
		IsActive:           courseInput.IsActive,
	}
	return &responseModel, nil
}

func UploadCourseImage(ctx context.Context, file model.CourseFile) (*model.UploadResult, error) {
	log.Info("UploadCourseImage called")
	isSuccess := model.UploadResult{}
	if *file.CourseID == "" {
		return &isSuccess, fmt.Errorf("course id is required")
	}
	storageC := bucket.NewStorageHandler()
	gproject := googleprojectlib.GetGoogleProjectID()
	err := storageC.InitializeStorageClient(ctx, gproject)
	if err != nil {
		log.Errorf("Failed to upload image to course: %v", err.Error())
		return &isSuccess, nil
	}
	bucketPath := *file.CourseID + "/" + file.File.Filename
	writer, err := storageC.UploadToGCS(ctx, bucketPath)
	if err != nil {
		log.Errorf("Failed to upload image to course: %v", err.Error())
		return &isSuccess, nil
	}
	defer writer.Close()
	fileBuffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(fileBuffer, file.File.File); err != nil {
		return &isSuccess, nil
	}
	currentBytes := fileBuffer.Bytes()
	_, err = io.Copy(writer, bytes.NewReader(currentBytes))
	if err != nil {
		return &isSuccess, err
	}
	getUrl := storageC.GetSignedURLForObject(bucketPath)
	// update course image in cassandra
	where := qb.Eq("id")
	updateQB := qb.Update("coursez.course").Set("imagebucket").Set("image").Where(where)
	updateQuery := updateQB.Query(*global.CassSession.Session).BindMap(qb.M{"id": file.CourseID, "imagebucket": bucketPath, "image": getUrl})
	if err := updateQuery.ExecRelease(); err != nil {
		return &isSuccess, err
	}
	isSuccessRes := true
	isSuccess.Success = &isSuccessRes
	isSuccess.URL = &getUrl
	return &isSuccess, nil
}

func UploadCoursePreviewVideo(ctx context.Context, file model.CourseFile) (*model.UploadResult, error) {
	log.Info("UploadCoursePreviewVideo called")
	isSuccess := model.UploadResult{}
	if *file.CourseID == "" {
		return &isSuccess, fmt.Errorf("course id is required")
	}
	storageC := bucket.NewStorageHandler()
	gproject := googleprojectlib.GetGoogleProjectID()
	err := storageC.InitializeStorageClient(ctx, gproject)
	if err != nil {
		log.Errorf("Failed to upload image to course: %v", err.Error())
		return &isSuccess, nil
	}
	bucketPath := *file.CourseID + "/" + file.File.Filename
	writer, err := storageC.UploadToGCS(ctx, bucketPath)
	if err != nil {
		log.Errorf("Failed to upload image to course: %v", err.Error())
		return &isSuccess, nil
	}
	defer writer.Close()
	fileBuffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(fileBuffer, file.File.File); err != nil {
		return &isSuccess, nil
	}
	currentBytes := fileBuffer.Bytes()
	_, err = io.Copy(writer, bytes.NewReader(currentBytes))
	if err != nil {
		return &isSuccess, err
	}
	getUrl := storageC.GetSignedURLForObject(bucketPath)
	// update course image in cassandra
	where := qb.Eq("id")
	updateQB := qb.Update("coursez.course").Set("previewvideobucket").Set("previewvideo").Where(where)
	updateQuery := updateQB.Query(*global.CassSession.Session).BindMap(qb.M{"id": file.CourseID, "previewvideobucket": bucketPath, "previewvideo": getUrl})
	if err := updateQuery.ExecRelease(); err != nil {
		return &isSuccess, err
	}
	isSuccessRes := true
	isSuccess.Success = &isSuccessRes
	isSuccess.URL = &getUrl
	return &isSuccess, nil
}

func UploadCourseTileImage(ctx context.Context, file model.CourseFile) (*model.UploadResult, error) {
	log.Info("UploadCourseTileImage called")
	isSuccess := model.UploadResult{}
	if *file.CourseID == "" {
		return &isSuccess, fmt.Errorf("course id is required")
	}
	storageC := bucket.NewStorageHandler()
	gproject := googleprojectlib.GetGoogleProjectID()
	err := storageC.InitializeStorageClient(ctx, gproject)
	if err != nil {
		log.Errorf("Failed to upload image to course: %v", err.Error())
		return &isSuccess, nil
	}
	bucketPath := *file.CourseID + "/" + file.File.Filename
	writer, err := storageC.UploadToGCS(ctx, bucketPath)
	if err != nil {
		log.Errorf("Failed to upload image to course: %v", err.Error())
		return &isSuccess, nil
	}
	defer writer.Close()
	fileBuffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(fileBuffer, file.File.File); err != nil {
		return &isSuccess, nil
	}
	currentBytes := fileBuffer.Bytes()
	_, err = io.Copy(writer, bytes.NewReader(currentBytes))
	if err != nil {
		return &isSuccess, err
	}
	getUrl := storageC.GetSignedURLForObject(bucketPath)
	where := qb.Eq("id")
	updateQB := qb.Update("coursez.course").Set("tileimagebucket").Set("tileimage").Where(where)
	updateQuery := updateQB.Query(*global.CassSession.Session).BindMap(qb.M{"id": file.CourseID, "tileimagebucket": bucketPath, "tileimage": getUrl})
	if err := updateQuery.ExecRelease(); err != nil {
		return &isSuccess, err
	}
	isSuccessRes := true
	isSuccess.Success = &isSuccessRes
	isSuccess.URL = &getUrl
	return &isSuccess, nil
}

func CourseUpdate(ctx context.Context, courseInput *model.CourseInput) (*model.Course, error) {
	log.Info("CourseUpdater called")
	if courseInput.ID == nil {
		return nil, fmt.Errorf("course id is required")
	}
	// set course input in cassandra
	courseID := *courseInput.ID
	// get course from cassandra
	cassandraCourse := coursez.Course{
		ID: courseID,
	}
	courses := []coursez.Course{}
	getQuery := global.CassSession.Session.Query(coursez.CourseTable.Get()).BindMap(qb.M{"id": courseID})
	if err := getQuery.SelectRelease(&courses); err != nil {
		return nil, err
	}
	if len(courses) < 1 {
		return nil, fmt.Errorf("course not found")
	}
	updateCols := make([]string, 0)
	cassandraCourse = courses[0]
	language := []string{}
	takeaways := []string{}
	outcomes := []string{}
	prequisites := []string{}
	goodFor := []string{}
	mustFor := []string{}
	relatedSkills := []string{}
	approvers := []string{}
	subCats := make([]coursez.SubCat, 0)
	subCatsRes := make([]*model.SubCategories, 0)

	for _, lang := range courseInput.Language {
		language = append(language, *lang)
	}
	for _, take := range courseInput.Benefits {
		takeaways = append(takeaways, *take)
	}
	for _, out := range courseInput.Outcomes {
		outcomes = append(outcomes, *out)
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
	for _, relSkill := range courseInput.RelatedSkills {
		relatedSkills = append(relatedSkills, *relSkill)
	}
	for _, approver := range courseInput.Approvers {
		approvers = append(approvers, *approver)
	}
	for _, subCat := range courseInput.SubCategories {
		var subC coursez.SubCat
		var subCR model.SubCategories
		subC.Name = *subCat.Name
		subC.Rank = *subCat.Rank
		subCR.Name = subCat.Name
		subCR.Rank = subCat.Rank
		subCats = append(subCats, subC)
		subCatsRes = append(subCatsRes, &subCR)
	}
	// update cassandraCourse with input
	if courseInput.Name != nil {
		updateCols = append(updateCols, "name")
		cassandraCourse.Name = *courseInput.Name
	}
	if courseInput.Description != nil {
		updateCols = append(updateCols, "description")
		cassandraCourse.Description = *courseInput.Description
	}
	if courseInput.Summary != nil {
		updateCols = append(updateCols, "summary")
		cassandraCourse.Summary = *courseInput.Summary
	}
	if courseInput.Instructor != nil {
		updateCols = append(updateCols, "instructor")
		cassandraCourse.Instructor = *courseInput.Instructor
	}
	if courseInput.Status != nil {
		updateCols = append(updateCols, "status")
		cassandraCourse.Status = (*courseInput.Status).String()
	}
	if courseInput.Language != nil {
		updateCols = append(updateCols, "language")
		cassandraCourse.Language = language
	}
	if courseInput.Benefits != nil {
		updateCols = append(updateCols, "benefits")
		cassandraCourse.Benefits = takeaways
	}
	if courseInput.Outcomes != nil {
		updateCols = append(updateCols, "outcomes")
		cassandraCourse.Outcomes = outcomes
	}
	if courseInput.Prequisites != nil {
		updateCols = append(updateCols, "prequisites")
		cassandraCourse.Prequisites = prequisites
	}
	if courseInput.GoodFor != nil {
		updateCols = append(updateCols, "goodfor")
		cassandraCourse.GoodFor = goodFor
	}
	if courseInput.MustFor != nil {
		updateCols = append(updateCols, "mustfor")
		cassandraCourse.MustFor = mustFor
	}
	if courseInput.RelatedSkills != nil {
		updateCols = append(updateCols, "related_skills")
		cassandraCourse.RelatedSkills = relatedSkills
	}
	if courseInput.Approvers != nil {
		updateCols = append(updateCols, "approvers")
		cassandraCourse.Approvers = approvers
	}
	if courseInput.Category != nil {
		updateCols = append(updateCols, "category")
		cassandraCourse.Category = *courseInput.Category
	}
	if courseInput.SubCategory != nil {
		updateCols = append(updateCols, "sub_category")
		cassandraCourse.SubCategory = *courseInput.SubCategory
	}
	if courseInput.SubCategories != nil {
		updateCols = append(updateCols, "sub_categories")
		cassandraCourse.SubCategories = subCats
	}
	if courseInput.Owner != nil {
		updateCols = append(updateCols, "owner")
		cassandraCourse.Owner = *courseInput.Owner
	}
	if courseInput.ExpertiseLevel != nil {
		updateCols = append(updateCols, "expertise_level")
		cassandraCourse.ExpertiseLevel = *courseInput.ExpertiseLevel
	}
	cassandraCourse.UpdatedAt = time.Now().Unix()
	if courseInput.CreatedBy != nil {
		updateCols = append(updateCols, "created_by")
		cassandraCourse.CreatedBy = *courseInput.CreatedBy
	}
	if courseInput.UpdatedBy != nil {
		updateCols = append(updateCols, "updated_by")
		cassandraCourse.UpdatedBy = *courseInput.UpdatedBy
	}
	if courseInput.Type != nil {
		updateCols = append(updateCols, "type")
		cassandraCourse.Type = *courseInput.Type
	}
	if courseInput.Duration != nil {
		updateCols = append(updateCols, "duration")
		cassandraCourse.Duration = *courseInput.Duration
	}
	if courseInput.IsDisplay != nil {
		updateCols = append(updateCols, "is_display")
		cassandraCourse.IsDisplay = *courseInput.IsDisplay
	}
	if courseInput.ExpectedCompletion != nil {
		updateCols = append(updateCols, "expected_completion_time")
		cassandraCourse.ExpectedCompletion = *courseInput.ExpectedCompletion
	}
	if courseInput.QaRequired != nil {
		updateCols = append(updateCols, "quality_control_check_reqd")
		cassandraCourse.QARequired = *courseInput.QaRequired
	}
	if courseInput.IsActive != nil {
		updateCols = append(updateCols, "is_active")
		cassandraCourse.IsActive = *courseInput.IsActive
	}
	updateCols = append(updateCols, "updated_at")
	// set course in cassandra
	upStms, uNames := coursez.CourseTable.Update(updateCols...)
	updateQuery := global.CassSession.Session.Query(upStms, uNames).BindStruct(&cassandraCourse)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	updated := strconv.FormatInt(cassandraCourse.UpdatedAt, 10)
	created := strconv.FormatInt(cassandraCourse.CreatedAt, 10)
	responseModel := model.Course{
		ID:                 &cassandraCourse.ID,
		Name:               courseInput.Name,
		Description:        courseInput.Description,
		Summary:            courseInput.Summary,
		Instructor:         courseInput.Instructor,
		Image:              &cassandraCourse.Image,
		PreviewVideo:       &cassandraCourse.PreviewVideo,
		TileImage:          &cassandraCourse.TileImage,
		Owner:              courseInput.Owner,
		Duration:           &cassandraCourse.Duration,
		ExpertiseLevel:     courseInput.ExpertiseLevel,
		Language:           courseInput.Language,
		Benefits:           courseInput.Benefits,
		CreatedAt:          &created,
		UpdatedAt:          &updated,
		Type:               courseInput.Type,
		Prequisites:        courseInput.Prequisites,
		GoodFor:            courseInput.GoodFor,
		MustFor:            courseInput.MustFor,
		RelatedSkills:      courseInput.RelatedSkills,
		PublishDate:        courseInput.PublishDate,
		ExpiryDate:         courseInput.ExpiryDate,
		QaRequired:         courseInput.QaRequired,
		Approvers:          courseInput.Approvers,
		CreatedBy:          courseInput.CreatedBy,
		UpdatedBy:          courseInput.UpdatedBy,
		Status:             courseInput.Status,
		IsDisplay:          courseInput.IsDisplay,
		ExpectedCompletion: courseInput.ExpectedCompletion,
		Category:           courseInput.Category,
		SubCategory:        courseInput.SubCategory,
		SubCategories:      subCatsRes,
		Outcomes:           courseInput.Outcomes,
		IsActive:           courseInput.IsActive,
	}
	return &responseModel, nil
}
