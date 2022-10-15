package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/scylladb/gocqlx/v2/qb"
	log "github.com/sirupsen/logrus"

	"github.com/rs/xid"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/zicops-cass-pool/cassandra"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/helpers"
	"github.com/zicops/zicops-course-creator/lib/db/bucket"
	"github.com/zicops/zicops-course-creator/lib/googleprojectlib"
)

// func to get context and course input and set it in cassandra
func CourseCreator(ctx context.Context, courseInput *model.CourseInput) (*model.Course, error) {
	log.Info("CourseCreator called")
	// set course input in cassandra
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}

	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	email_creator := claims["email"].(string)
	lspId := claims["lsp_id"].(string)
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
	active := true
	if courseInput.IsActive != nil {
		active = *courseInput.IsActive
	}
	words := []string{}
	if courseInput.Name != nil {
		name := strings.ToLower(*courseInput.Name)
		wordsLocal := strings.Split(name, " ")
		words = append(words, wordsLocal...)
	}

	cassandraCourse := coursez.Course{
		ID:                 guid.String(),
		Name:               *courseInput.Name,
		Words:              words,
		LspId:              lspId,
		Publisher:          *courseInput.Publisher,
		Description:        *courseInput.Description,
		Image:              "https://storage.googleapis.com/zicops.com/default-image.jpg",
		PreviewVideo:       "",
		TileImage:          "https://storage.googleapis.com/zicops.com/default-image.jpg",
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
		IsActive:           active,
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
	cassandraCourse.CreatedBy = email_creator
	cassandraCourse.UpdatedBy = email_creator
	if courseInput.Category != nil {
		cassandraCourse.Category = *courseInput.Category
	}
	if courseInput.ExpectedCompletion != nil {
		cassandraCourse.ExpectedCompletion = *courseInput.ExpectedCompletion
	}
	if courseInput.SubCategory != nil {
		cassandraCourse.SubCategory = *courseInput.SubCategory
	}
	courseInput.CreatedBy = &email_creator
	courseInput.UpdatedBy = &email_creator
	// set course in cassandra
	insertQuery := CassSession.Query(coursez.CourseTable.Insert()).BindStruct(cassandraCourse)
	if err := insertQuery.ExecRelease(); err != nil {
		return nil, err
	}

	created := strconv.FormatInt(cassandraCourse.CreatedAt, 10)
	responseModel := model.Course{
		ID:                 &cassandraCourse.ID,
		LspID:              &lspId,
		Publisher:          &cassandraCourse.Publisher,
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
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspID := claims["lsp_id"].(string)
	if *file.CourseID == "" {
		return &isSuccess, fmt.Errorf("course id is required")
	}
	storageC := bucket.NewStorageHandler()
	gproject := googleprojectlib.GetGoogleProjectID()
	err = storageC.InitializeStorageClient(ctx, gproject, lspID)
	if err != nil {
		log.Errorf("Failed to upload image to course: %v", err.Error())
		return &isSuccess, nil
	}
	bucketPath := *file.CourseID + "/" + file.File.Filename
	writer, err := storageC.UploadToGCS(ctx, bucketPath, map[string]string{})
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
	updateQuery := fmt.Sprintf("UPDATE coursez.course SET imagebucket='%s', image='%s' WHERE id='%s' AND lsp_id='%s' AND is_active=true", bucketPath, getUrl, *file.CourseID, lspID)
	updateQ := CassSession.Query(updateQuery, nil)
	if err := updateQ.ExecRelease(); err != nil {
		return nil, err
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
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspID := claims["lsp_id"].(string)
	if lspID == "" {
		return &isSuccess, fmt.Errorf("lsp id is required")
	}
	if *file.CourseID == "" {
		return &isSuccess, fmt.Errorf("course id is required")
	}
	storageC := bucket.NewStorageHandler()
	gproject := googleprojectlib.GetGoogleProjectID()
	err = storageC.InitializeStorageClient(ctx, gproject, lspID)
	if err != nil {
		log.Errorf("Failed to upload image to course: %v", err.Error())
		return &isSuccess, nil
	}
	bucketPath := *file.CourseID + "/" + file.File.Filename
	writer, err := storageC.UploadToGCS(ctx, bucketPath, map[string]string{})
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
	updateQuery := fmt.Sprintf("UPDATE coursez.course SET previewvideobucket='%s', previewvideo='%s' WHERE id='%s' AND lsp_id='%s' AND is_active=true", bucketPath, getUrl, *file.CourseID, lspID)
	updateQ := CassSession.Query(updateQuery, nil)
	if err := updateQ.ExecRelease(); err != nil {
		return nil, err
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
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lspID := claims["lsp_id"].(string)
	if lspID == "" {
		return &isSuccess, fmt.Errorf("lsp id is required")
	}
	storageC := bucket.NewStorageHandler()
	gproject := googleprojectlib.GetGoogleProjectID()
	err = storageC.InitializeStorageClient(ctx, gproject, lspID)
	if err != nil {
		log.Errorf("Failed to upload image to course: %v", err.Error())
		return &isSuccess, nil
	}
	bucketPath := *file.CourseID + "/" + file.File.Filename
	writer, err := storageC.UploadToGCS(ctx, bucketPath, map[string]string{})
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
	updateQuery := fmt.Sprintf("UPDATE coursez.course SET tileimagebucket='%s', tileimage='%s' WHERE id='%s' AND lsp_id='%s' AND is_active=true", bucketPath, getUrl, *file.CourseID, lspID)
	updateQ := CassSession.Query(updateQuery, nil)
	if err := updateQ.ExecRelease(); err != nil {
		return nil, err
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
	session, err := cassandra.GetCassSession("coursez")
	if err != nil {
		return nil, err
	}
	CassSession := session
	claims, err := helpers.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	email_creator := claims["email"].(string)
	lspId := claims["lsp_id"].(string)
	// set course input in cassandra
	courseID := *courseInput.ID
	// get course from cassandra
	cassandraCourse := coursez.Course{
		ID: courseID,
	}
	courses := []coursez.Course{}
	getQuery := CassSession.Query(coursez.CourseTable.Get()).BindMap(qb.M{"id": courseID, "lsp_id": lspId, "is_active": true})
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
	if courseInput.Name != nil && *courseInput.Name != cassandraCourse.Name {
		name := *courseInput.Name
		wordsLocal := strings.Split(name, " ")
		words := make([]string, 0)
		words = append(words, wordsLocal...)
		updateCols = append(updateCols, "name")
		updateCols = append(updateCols, "words")
		cassandraCourse.Words = words
		cassandraCourse.Name = *courseInput.Name
	}
	if courseInput.Description != nil && *courseInput.Description != cassandraCourse.Description {
		updateCols = append(updateCols, "description")
		cassandraCourse.Description = *courseInput.Description
	}
	if courseInput.Summary != nil && *courseInput.Summary != cassandraCourse.Summary {
		updateCols = append(updateCols, "summary")
		cassandraCourse.Summary = *courseInput.Summary
	}
	if courseInput.Instructor != nil && *courseInput.Instructor != cassandraCourse.Instructor {
		updateCols = append(updateCols, "instructor")
		cassandraCourse.Instructor = *courseInput.Instructor
	}
	if courseInput.Status != nil && string(*courseInput.Status) != cassandraCourse.Status {
		updateCols = append(updateCols, "status")
		cassandraCourse.Status = (*courseInput.Status).String()
	}
	if courseInput.Language != nil && !reflect.DeepEqual(language, cassandraCourse.Language) {
		updateCols = append(updateCols, "language")
		cassandraCourse.Language = language
	}
	if courseInput.Benefits != nil && !reflect.DeepEqual(takeaways, cassandraCourse.Benefits) {
		updateCols = append(updateCols, "benefits")
		cassandraCourse.Benefits = takeaways
	}
	if courseInput.Outcomes != nil && !reflect.DeepEqual(outcomes, cassandraCourse.Outcomes) {
		updateCols = append(updateCols, "outcomes")
		cassandraCourse.Outcomes = outcomes
	}
	if courseInput.Prequisites != nil && !reflect.DeepEqual(prequisites, cassandraCourse.Prequisites) {
		updateCols = append(updateCols, "prequisites")
		cassandraCourse.Prequisites = prequisites
	}
	if courseInput.GoodFor != nil && !reflect.DeepEqual(goodFor, cassandraCourse.GoodFor) {
		updateCols = append(updateCols, "goodfor")
		cassandraCourse.GoodFor = goodFor
	}
	if courseInput.MustFor != nil && !reflect.DeepEqual(mustFor, cassandraCourse.MustFor) {
		updateCols = append(updateCols, "mustfor")
		cassandraCourse.MustFor = mustFor
	}
	if courseInput.RelatedSkills != nil && !reflect.DeepEqual(relatedSkills, cassandraCourse.RelatedSkills) {
		updateCols = append(updateCols, "related_skills")
		cassandraCourse.RelatedSkills = relatedSkills
	}
	if courseInput.Approvers != nil && !reflect.DeepEqual(approvers, cassandraCourse.Approvers) {
		updateCols = append(updateCols, "approvers")
		cassandraCourse.Approvers = approvers
	}
	if courseInput.Category != nil && *courseInput.Category != cassandraCourse.Category {
		updateCols = append(updateCols, "category")
		cassandraCourse.Category = *courseInput.Category
	}
	if courseInput.SubCategory != nil && *courseInput.SubCategory != cassandraCourse.SubCategory {
		updateCols = append(updateCols, "sub_category")
		cassandraCourse.SubCategory = *courseInput.SubCategory
	}
	if courseInput.SubCategories != nil && !reflect.DeepEqual(subCats, cassandraCourse.SubCategories) {
		updateCols = append(updateCols, "sub_categories")
		cassandraCourse.SubCategories = subCats
	}
	if courseInput.Owner != nil && *courseInput.Owner != cassandraCourse.Owner {
		updateCols = append(updateCols, "owner")
		cassandraCourse.Owner = *courseInput.Owner
	}
	if courseInput.ExpertiseLevel != nil && *courseInput.ExpertiseLevel != cassandraCourse.ExpertiseLevel {
		updateCols = append(updateCols, "expertise_level")
		cassandraCourse.ExpertiseLevel = *courseInput.ExpertiseLevel
	}
	cassandraCourse.UpdatedAt = time.Now().Unix()
	if email_creator != "" && email_creator != cassandraCourse.UpdatedBy {
		updateCols = append(updateCols, "updated_by")
		cassandraCourse.UpdatedBy = email_creator
	}
	if courseInput.Type != nil && *courseInput.Type != cassandraCourse.Type {
		updateCols = append(updateCols, "type")
		cassandraCourse.Type = *courseInput.Type
	}
	if courseInput.IsDisplay != nil && *courseInput.IsDisplay != cassandraCourse.IsDisplay {
		updateCols = append(updateCols, "is_display")
		cassandraCourse.IsDisplay = *courseInput.IsDisplay
	}
	if courseInput.ExpectedCompletion != nil && *courseInput.ExpectedCompletion != cassandraCourse.ExpectedCompletion {
		updateCols = append(updateCols, "expected_completion_time")
		cassandraCourse.ExpectedCompletion = *courseInput.ExpectedCompletion
	}
	if courseInput.QaRequired != nil && *courseInput.QaRequired != cassandraCourse.QARequired {
		updateCols = append(updateCols, "quality_control_check_reqd")
		cassandraCourse.QARequired = *courseInput.QaRequired
	}
	if courseInput.Publisher != nil && *courseInput.Publisher != cassandraCourse.Publisher {
		updateCols = append(updateCols, "publisher")
		cassandraCourse.Publisher = *courseInput.Publisher
	}
	if len(updateCols) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}
	updateCols = append(updateCols, "updated_at")
	// set course in cassandra
	upStms, uNames := coursez.CourseTable.Update(updateCols...)
	updateQuery := CassSession.Query(upStms, uNames).BindStruct(cassandraCourse)
	if err := updateQuery.ExecRelease(); err != nil {
		return nil, err
	}
	updated := strconv.FormatInt(cassandraCourse.UpdatedAt, 10)
	created := strconv.FormatInt(cassandraCourse.CreatedAt, 10)
	responseModel := model.Course{
		ID:                 &cassandraCourse.ID,
		Name:               courseInput.Name,
		LspID:              &cassandraCourse.LspId,
		Publisher:          &cassandraCourse.Publisher,
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
