// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

type Chapter struct {
	ID          *string `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ModuleID    *string `json:"moduleId"`
	CourseID    *string `json:"courseId"`
	CreatedAt   *string `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
	Sequence    *int    `json:"sequence"`
}

type ChapterInput struct {
	ID          *string `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ModuleID    *string `json:"moduleId"`
	CourseID    *string `json:"courseId"`
	Sequence    *int    `json:"sequence"`
}

type Course struct {
	ID                 *string          `json:"id"`
	Name               *string          `json:"name"`
	Description        *string          `json:"description"`
	Summary            *string          `json:"summary"`
	Instructor         *string          `json:"instructor"`
	Image              *string          `json:"image"`
	PreviewVideo       *string          `json:"previewVideo"`
	TileImage          *string          `json:"tileImage"`
	Owner              *string          `json:"owner"`
	Duration           *int             `json:"duration"`
	ExpertiseLevel     *string          `json:"expertise_level"`
	Language           []*string        `json:"language"`
	Benefits           []*string        `json:"benefits"`
	Outcomes           []*string        `json:"outcomes"`
	CreatedAt          *string          `json:"created_at"`
	UpdatedAt          *string          `json:"updated_at"`
	Type               *string          `json:"type"`
	Prequisites        []*string        `json:"prequisites"`
	GoodFor            []*string        `json:"goodFor"`
	MustFor            []*string        `json:"mustFor"`
	RelatedSkills      []*string        `json:"related_skills"`
	PublishDate        *string          `json:"publish_date"`
	ExpiryDate         *string          `json:"expiry_date"`
	ExpectedCompletion *string          `json:"expected_completion"`
	QaRequired         *bool            `json:"qa_required"`
	Approvers          []*string        `json:"approvers"`
	CreatedBy          *string          `json:"created_by"`
	UpdatedBy          *string          `json:"updated_by"`
	Status             *Status          `json:"status"`
	IsDisplay          *bool            `json:"is_display"`
	Category           *string          `json:"category"`
	SubCategory        *string          `json:"sub_category"`
	SubCategories      []*SubCategories `json:"sub_categories"`
	IsActive           *bool            `json:"is_active"`
}

type CourseFile struct {
	File     *graphql.Upload `json:"file"`
	CourseID *string         `json:"courseId"`
}

type CourseInput struct {
	ID                 *string               `json:"id"`
	Name               *string               `json:"name"`
	Description        *string               `json:"description"`
	Summary            *string               `json:"summary"`
	Instructor         *string               `json:"instructor"`
	Image              *string               `json:"image"`
	PreviewVideo       *string               `json:"previewVideo"`
	TileImage          *string               `json:"tileImage"`
	Owner              *string               `json:"owner"`
	Duration           *int                  `json:"duration"`
	ExpertiseLevel     *string               `json:"expertise_level"`
	Language           []*string             `json:"language"`
	Benefits           []*string             `json:"benefits"`
	Outcomes           []*string             `json:"outcomes"`
	Type               *string               `json:"type"`
	Prequisites        []*string             `json:"prequisites"`
	GoodFor            []*string             `json:"goodFor"`
	MustFor            []*string             `json:"mustFor"`
	RelatedSkills      []*string             `json:"related_skills"`
	PublishDate        *string               `json:"publish_date"`
	ExpiryDate         *string               `json:"expiry_date"`
	ExpectedCompletion *string               `json:"expected_completion"`
	QaRequired         *bool                 `json:"qa_required"`
	Approvers          []*string             `json:"approvers"`
	CreatedBy          *string               `json:"created_by"`
	UpdatedBy          *string               `json:"updated_by"`
	Status             *Status               `json:"status"`
	IsDisplay          *bool                 `json:"is_display"`
	Category           *string               `json:"category"`
	SubCategory        *string               `json:"sub_category"`
	SubCategories      []*SubCategoriesInput `json:"sub_categories"`
	IsActive           *bool                 `json:"is_active"`
}

type Exam struct {
	ID           *string `json:"id"`
	Name         *string `json:"Name"`
	Description  *string `json:"Description"`
	Code         *string `json:"Code"`
	QpID         *string `json:"QpId"`
	CreatedAt    *string `json:"CreatedAt"`
	UpdatedAt    *string `json:"UpdatedAt"`
	CreatedBy    *string `json:"CreatedBy"`
	UpdatedBy    *string `json:"UpdatedBy"`
	IsActive     *bool   `json:"IsActive"`
	Type         *string `json:"Type"`
	ScheduleType *string `json:"ScheduleType"`
	Duration     *int    `json:"Duration"`
	Status       *string `json:"Status"`
	Category     *string `json:"Category"`
	SubCategory  *string `json:"SubCategory"`
}

type ExamCohort struct {
	ID        *string `json:"id"`
	ExamID    *string `json:"ExamId"`
	CohortID  *string `json:"CohortId"`
	CreatedAt *string `json:"CreatedAt"`
	UpdatedAt *string `json:"UpdatedAt"`
	CreatedBy *string `json:"CreatedBy"`
	UpdatedBy *string `json:"UpdatedBy"`
	IsActive  *bool   `json:"IsActive"`
}

type ExamCohortInput struct {
	ID        *string `json:"id"`
	ExamID    *string `json:"ExamId"`
	CohortID  *string `json:"CohortId"`
	CreatedAt *string `json:"CreatedAt"`
	UpdatedAt *string `json:"UpdatedAt"`
	CreatedBy *string `json:"CreatedBy"`
	UpdatedBy *string `json:"UpdatedBy"`
	IsActive  *bool   `json:"IsActive"`
}

type ExamConfiguration struct {
	ID           *string `json:"id"`
	ExamID       *string `json:"ExamId"`
	Shuffle      *bool   `json:"Shuffle"`
	DisplayHints *bool   `json:"DisplayHints"`
	ShowAnswer   *bool   `json:"ShowAnswer"`
	ShowResult   *bool   `json:"ShowResult"`
	CreatedBy    *string `json:"CreatedBy"`
	UpdatedBy    *string `json:"UpdatedBy"`
	IsActive     *bool   `json:"IsActive"`
	CreatedAt    *string `json:"CreatedAt"`
	UpdatedAt    *string `json:"UpdatedAt"`
}

type ExamConfigurationInput struct {
	ID           *string `json:"id"`
	ExamID       *string `json:"ExamId"`
	Shuffle      *bool   `json:"Shuffle"`
	DisplayHints *bool   `json:"DisplayHints"`
	ShowAnswer   *bool   `json:"ShowAnswer"`
	ShowResult   *bool   `json:"ShowResult"`
	CreatedBy    *string `json:"CreatedBy"`
	UpdatedBy    *string `json:"UpdatedBy"`
	IsActive     *bool   `json:"IsActive"`
	CreatedAt    *string `json:"CreatedAt"`
	UpdatedAt    *string `json:"UpdatedAt"`
}

type ExamInput struct {
	ID           *string `json:"id"`
	Name         *string `json:"Name"`
	Description  *string `json:"Description"`
	Code         *string `json:"Code"`
	QpID         *string `json:"QpId"`
	CreatedAt    *string `json:"CreatedAt"`
	UpdatedAt    *string `json:"UpdatedAt"`
	CreatedBy    *string `json:"CreatedBy"`
	UpdatedBy    *string `json:"UpdatedBy"`
	IsActive     *bool   `json:"IsActive"`
	Type         *string `json:"Type"`
	ScheduleType *string `json:"ScheduleType"`
	Duration     *int    `json:"Duration"`
	Status       *string `json:"Status"`
	Category     *string `json:"Category"`
	SubCategory  *string `json:"SubCategory"`
}

type ExamInstruction struct {
	ID              *string `json:"id"`
	ExamID          *string `json:"ExamId"`
	Instructions    *string `json:"Instructions"`
	PassingCriteria *string `json:"PassingCriteria"`
	NoAttempts      *int    `json:"NoAttempts"`
	AccessType      *string `json:"AccessType"`
	CreatedAt       *string `json:"CreatedAt"`
	UpdatedAt       *string `json:"UpdatedAt"`
	CreatedBy       *string `json:"CreatedBy"`
	UpdatedBy       *string `json:"UpdatedBy"`
	IsActive        *bool   `json:"IsActive"`
}

type ExamInstructionInput struct {
	ID              *string `json:"id"`
	ExamID          *string `json:"ExamId"`
	Instructions    *string `json:"Instructions"`
	PassingCriteria *string `json:"PassingCriteria"`
	NoAttempts      *int    `json:"NoAttempts"`
	AccessType      *string `json:"AccessType"`
	CreatedAt       *string `json:"CreatedAt"`
	UpdatedAt       *string `json:"UpdatedAt"`
	CreatedBy       *string `json:"CreatedBy"`
	UpdatedBy       *string `json:"UpdatedBy"`
	IsActive        *bool   `json:"IsActive"`
}

type ExamSchedule struct {
	ID         *string `json:"id"`
	ExamID     *string `json:"ExamId"`
	Start      *int    `json:"Start"`
	End        *int    `json:"End"`
	BufferTime *int    `json:"BufferTime"`
	CreatedAt  *string `json:"CreatedAt"`
	UpdatedAt  *string `json:"UpdatedAt"`
	CreatedBy  *string `json:"CreatedBy"`
	UpdatedBy  *string `json:"UpdatedBy"`
	IsActive   *bool   `json:"IsActive"`
}

type ExamScheduleInput struct {
	ID         *string `json:"id"`
	ExamID     *string `json:"ExamId"`
	Start      *int    `json:"Start"`
	End        *int    `json:"End"`
	BufferTime *int    `json:"BufferTime"`
	CreatedAt  *string `json:"CreatedAt"`
	UpdatedAt  *string `json:"UpdatedAt"`
	CreatedBy  *string `json:"CreatedBy"`
	UpdatedBy  *string `json:"UpdatedBy"`
	IsActive   *bool   `json:"IsActive"`
}

type MapSectionToBankInput struct {
	ID              *string `json:"id"`
	QbID            *string `json:"QbId"`
	SectionID       *string `json:"SectionId"`
	DifficultyLevel *string `json:"DifficultyLevel"`
	TotalQuestions  *int    `json:"TotalQuestions"`
	QuestionMarks   *string `json:"QuestionMarks"`
	QuestionType    *string `json:"QuestionType"`
	RetrieveType    *string `json:"RetrieveType"`
	CreatedAt       *string `json:"CreatedAt"`
	UpdatedAt       *string `json:"UpdatedAt"`
	CreatedBy       *string `json:"CreatedBy"`
	UpdatedBy       *string `json:"UpdatedBy"`
	IsActive        *bool   `json:"IsActive"`
}

type Module struct {
	ID          *string `json:"id"`
	Name        *string `json:"name"`
	IsChapter   *bool   `json:"isChapter"`
	Description *string `json:"description"`
	CourseID    *string `json:"courseId"`
	Owner       *string `json:"owner"`
	Duration    *int    `json:"duration"`
	CreatedAt   *string `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
	Level       *string `json:"level"`
	Sequence    *int    `json:"sequence"`
	SetGlobal   *bool   `json:"setGlobal"`
}

type ModuleInput struct {
	ID          *string `json:"id"`
	Name        *string `json:"name"`
	IsChapter   *bool   `json:"isChapter"`
	Description *string `json:"description"`
	CourseID    *string `json:"courseId"`
	Owner       *string `json:"owner"`
	Duration    *int    `json:"duration"`
	Level       *string `json:"level"`
	Sequence    *int    `json:"sequence"`
	SetGlobal   *bool   `json:"setGlobal"`
}

type QuestionBank struct {
	ID          *string `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	SubCategory *string `json:"sub_category"`
	CreatedAt   *string `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
	CreatedBy   *string `json:"created_by"`
	UpdatedBy   *string `json:"updated_by"`
	IsActive    *bool   `json:"is_active"`
	IsDefault   *bool   `json:"is_default"`
	Owner       *string `json:"owner"`
}

type QuestionBankInput struct {
	ID          *string `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	SubCategory *string `json:"sub_category"`
	CreatedAt   *string `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
	CreatedBy   *string `json:"created_by"`
	UpdatedBy   *string `json:"updated_by"`
	IsActive    *bool   `json:"is_active"`
	IsDefault   *bool   `json:"is_default"`
	Owner       *string `json:"owner"`
}

type QuestionBankQuestion struct {
	ID             *string `json:"id"`
	Name           *string `json:"Name"`
	Description    *string `json:"Description"`
	Type           *string `json:"Type"`
	Difficulty     *int    `json:"Difficulty"`
	Attachment     *string `json:"Attachment"`
	AttachmentType *string `json:"AttachmentType"`
	Hint           *string `json:"Hint"`
	QbmID          *string `json:"QbmId"`
	Status         *string `json:"Status"`
	CreatedAt      *string `json:"CreatedAt"`
	UpdatedAt      *string `json:"UpdatedAt"`
	CreatedBy      *string `json:"CreatedBy"`
	UpdatedBy      *string `json:"UpdatedBy"`
}

type QuestionBankQuestionInput struct {
	ID             *string         `json:"id"`
	Name           *string         `json:"Name"`
	Description    *string         `json:"Description"`
	Type           *string         `json:"Type"`
	Difficulty     *int            `json:"Difficulty"`
	File           *graphql.Upload `json:"File"`
	AttachmentType *string         `json:"AttachmentType"`
	Hint           *string         `json:"Hint"`
	QbmID          *string         `json:"QbmId"`
	Status         *string         `json:"Status"`
	CreatedAt      *string         `json:"CreatedAt"`
	UpdatedAt      *string         `json:"UpdatedAt"`
	CreatedBy      *string         `json:"CreatedBy"`
	UpdatedBy      *string         `json:"UpdatedBy"`
}

type QuestionOption struct {
	ID             *string `json:"id"`
	QmID           *string `json:"QmId"`
	Description    *string `json:"Description"`
	IsCorrect      *bool   `json:"IsCorrect"`
	CreatedAt      *string `json:"CreatedAt"`
	UpdatedAt      *string `json:"UpdatedAt"`
	CreatedBy      *string `json:"CreatedBy"`
	UpdatedBy      *string `json:"UpdatedBy"`
	AttachmentType *string `json:"AttachmentType"`
	Attachment     *string `json:"Attachment"`
	IsActive       *bool   `json:"IsActive"`
}

type QuestionOptionInput struct {
	ID             *string         `json:"id"`
	QmID           *string         `json:"QmId"`
	Description    *string         `json:"Description"`
	IsCorrect      *bool           `json:"IsCorrect"`
	CreatedAt      *string         `json:"CreatedAt"`
	UpdatedAt      *string         `json:"UpdatedAt"`
	CreatedBy      *string         `json:"CreatedBy"`
	UpdatedBy      *string         `json:"UpdatedBy"`
	AttachmentType *string         `json:"AttachmentType"`
	File           *graphql.Upload `json:"File"`
	IsActive       *bool           `json:"IsActive"`
}

type QuestionPaper struct {
	ID                *string `json:"id"`
	Name              *string `json:"name"`
	Category          *string `json:"Category"`
	SubCategory       *string `json:"SubCategory"`
	CreatedAt         *string `json:"CreatedAt"`
	UpdatedAt         *string `json:"UpdatedAt"`
	CreatedBy         *string `json:"CreatedBy"`
	UpdatedBy         *string `json:"UpdatedBy"`
	IsActive          *bool   `json:"IsActive"`
	DifficultyLevel   *string `json:"DifficultyLevel"`
	SectionWise       *bool   `json:"SectionWise"`
	Description       *string `json:"Description"`
	SuggestedDuration *string `json:"SuggestedDuration"`
	Status            *string `json:"Status"`
}

type QuestionPaperInput struct {
	ID                *string `json:"id"`
	Name              *string `json:"name"`
	Category          *string `json:"Category"`
	SubCategory       *string `json:"SubCategory"`
	CreatedAt         *string `json:"CreatedAt"`
	UpdatedAt         *string `json:"UpdatedAt"`
	CreatedBy         *string `json:"CreatedBy"`
	UpdatedBy         *string `json:"UpdatedBy"`
	IsActive          *bool   `json:"IsActive"`
	DifficultyLevel   *string `json:"DifficultyLevel"`
	SectionWise       *bool   `json:"SectionWise"`
	Description       *string `json:"Description"`
	SuggestedDuration *string `json:"SuggestedDuration"`
	Status            *string `json:"Status"`
}

type QuestionPaperSection struct {
	ID              *string `json:"id"`
	QpID            *string `json:"QpId"`
	Name            *string `json:"Name"`
	Description     *string `json:"Description"`
	CreatedAt       *string `json:"CreatedAt"`
	UpdatedAt       *string `json:"UpdatedAt"`
	CreatedBy       *string `json:"CreatedBy"`
	UpdatedBy       *string `json:"UpdatedBy"`
	IsActive        *bool   `json:"IsActive"`
	Type            *string `json:"Type"`
	DifficultyLevel *string `json:"DifficultyLevel"`
	TotalQuestions  *int    `json:"TotalQuestions"`
}

type QuestionPaperSectionInput struct {
	ID              *string `json:"id"`
	QpID            *string `json:"QpId"`
	Name            *string `json:"Name"`
	Description     *string `json:"Description"`
	CreatedAt       *string `json:"CreatedAt"`
	UpdatedAt       *string `json:"UpdatedAt"`
	CreatedBy       *string `json:"CreatedBy"`
	UpdatedBy       *string `json:"UpdatedBy"`
	IsActive        *bool   `json:"IsActive"`
	Type            *string `json:"Type"`
	DifficultyLevel *string `json:"DifficultyLevel"`
	TotalQuestions  *int    `json:"TotalQuestions"`
}

type Quiz struct {
	ID          *string `json:"id"`
	Name        *string `json:"name"`
	Category    *string `json:"category"`
	Type        *string `json:"type"`
	IsMandatory *bool   `json:"isMandatory"`
	CreatedAt   *string `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
	TopicID     *string `json:"topicId"`
	Sequence    *int    `json:"sequence"`
	StartTime   *int    `json:"startTime"`
}

type QuizDescriptive struct {
	QuizID        *string `json:"quizId"`
	Question      *string `json:"question"`
	CorrectAnswer *string `json:"correctAnswer"`
	Explanation   *string `json:"explanation"`
}

type QuizFile struct {
	QuizID *string         `json:"quizId"`
	Type   *string         `json:"type"`
	Name   *string         `json:"name"`
	File   *graphql.Upload `json:"file"`
}

type QuizInput struct {
	ID          *string `json:"id"`
	Name        *string `json:"name"`
	Category    *string `json:"category"`
	Type        *string `json:"type"`
	IsMandatory *bool   `json:"isMandatory"`
	CreatedAt   *string `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
	TopicID     *string `json:"topicId"`
	Sequence    *int    `json:"sequence"`
	StartTime   *int    `json:"startTime"`
}

type QuizMcq struct {
	QuizID        *string   `json:"quizId"`
	Question      *string   `json:"question"`
	Options       []*string `json:"options"`
	CorrectOption *string   `json:"correctOption"`
	Explanation   *string   `json:"explanation"`
}

type SectionFixedQuestions struct {
	ID         *string `json:"id"`
	SqbID      *string `json:"SqbId"`
	QuestionID *string `json:"QuestionId"`
	CreatedAt  *string `json:"CreatedAt"`
	UpdatedAt  *string `json:"UpdatedAt"`
	CreatedBy  *string `json:"CreatedBy"`
	UpdatedBy  *string `json:"UpdatedBy"`
	IsActive   *bool   `json:"IsActive"`
}

type SectionFixedQuestionsInput struct {
	ID         *string `json:"id"`
	SqbID      *string `json:"SqbId"`
	QuestionID *string `json:"QuestionId"`
	CreatedAt  *string `json:"CreatedAt"`
	UpdatedAt  *string `json:"UpdatedAt"`
	CreatedBy  *string `json:"CreatedBy"`
	UpdatedBy  *string `json:"UpdatedBy"`
	IsActive   *bool   `json:"IsActive"`
}

type SectionQBMapping struct {
	ID              *string `json:"id"`
	QbID            *string `json:"QbId"`
	SectionID       *string `json:"SectionId"`
	DifficultyLevel *string `json:"DifficultyLevel"`
	TotalQuestions  *int    `json:"TotalQuestions"`
	QuestionMarks   *string `json:"QuestionMarks"`
	QuestionType    *string `json:"QuestionType"`
	RetrieveType    *string `json:"RetrieveType"`
	CreatedAt       *string `json:"CreatedAt"`
	UpdatedAt       *string `json:"UpdatedAt"`
	CreatedBy       *string `json:"CreatedBy"`
	UpdatedBy       *string `json:"UpdatedBy"`
	IsActive        *bool   `json:"IsActive"`
}

type StaticContent struct {
	Type      *Type           `json:"type"`
	File      *graphql.Upload `json:"file"`
	CourseID  *string         `json:"courseId"`
	ContentID *string         `json:"contentId"`
}

type Topic struct {
	ID          *string `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Type        *string `json:"type"`
	ModuleID    *string `json:"moduleId"`
	ChapterID   *string `json:"chapterId"`
	CourseID    *string `json:"courseId"`
	CreatedAt   *string `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
	Sequence    *int    `json:"sequence"`
	CreatedBy   *string `json:"created_by"`
	UpdatedBy   *string `json:"updated_by"`
}

type TopicContent struct {
	ID                *string `json:"id"`
	Language          *string `json:"language"`
	TopicID           *string `json:"topicId"`
	StartTime         *int    `json:"startTime"`
	Duration          *int    `json:"duration"`
	SkipIntroDuration *int    `json:"skipIntroDuration"`
	NextShowTime      *int    `json:"nextShowTime"`
	FromEndTime       *int    `json:"fromEndTime"`
	CreatedAt         *string `json:"created_at"`
	UpdatedAt         *string `json:"updated_at"`
	Type              *string `json:"type"`
	IsDefault         *bool   `json:"is_default"`
}

type TopicContentInput struct {
	Language          *string `json:"language"`
	ContentID         *string `json:"contentId"`
	StartTime         *int    `json:"startTime"`
	Duration          *int    `json:"duration"`
	SkipIntroDuration *int    `json:"skipIntroDuration"`
	NextShowTime      *int    `json:"nextShowTime"`
	FromEndTime       *int    `json:"fromEndTime"`
	Type              *string `json:"type"`
	IsDefault         *bool   `json:"is_default"`
}

type TopicExam struct {
	ID        *string `json:"id"`
	TopicID   *string `json:"topicId"`
	ExamID    *string `json:"examId"`
	CourseID  *string `json:"courseId"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
	Language  *string `json:"language"`
}

type TopicExamInput struct {
	ID        *string `json:"id"`
	TopicID   *string `json:"topicId"`
	ExamID    *string `json:"examId"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
	Language  *string `json:"language"`
}

type TopicInput struct {
	ID          *string `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Type        *string `json:"type"`
	ModuleID    *string `json:"moduleId"`
	ChapterID   *string `json:"chapterId"`
	CourseID    *string `json:"courseId"`
	Sequence    *int    `json:"sequence"`
	CreatedBy   *string `json:"created_by"`
	UpdatedBy   *string `json:"updated_by"`
}

type TopicResourceInput struct {
	Name      *string         `json:"name"`
	Type      *string         `json:"type"`
	TopicID   *string         `json:"topicId"`
	CreatedAt *string         `json:"created_at"`
	UpdatedAt *string         `json:"updated_at"`
	CreatedBy *string         `json:"created_by"`
	UpdatedBy *string         `json:"updated_by"`
	URL       *string         `json:"url"`
	File      *graphql.Upload `json:"file"`
}

type TopicSubtitle struct {
	File     *graphql.Upload `json:"file"`
	CourseID *string         `json:"courseId"`
	TopicID  *string         `json:"topicId"`
	Language *string         `json:"language"`
}

type TopicVideo struct {
	File      *graphql.Upload `json:"file"`
	CourseID  *string         `json:"courseId"`
	ContentID *string         `json:"contentId"`
}

type UploadResult struct {
	Success *bool   `json:"success"`
	URL     *string `json:"url"`
}

type UploadResultSubtitles struct {
	Success  *bool   `json:"success"`
	URL      *string `json:"url"`
	Language *string `json:"language"`
}

type SubCategories struct {
	Name *string `json:"name"`
	Rank *int    `json:"rank"`
}

type SubCategoriesInput struct {
	Name *string `json:"name"`
	Rank *int    `json:"rank"`
}

type Status string

const (
	StatusSaved           Status = "SAVED"
	StatusApprovalPending Status = "APPROVAL_PENDING"
	StatusOnHold          Status = "ON_HOLD"
	StatusApproved        Status = "APPROVED"
	StatusPublished       Status = "PUBLISHED"
	StatusRejected        Status = "REJECTED"
)

var AllStatus = []Status{
	StatusSaved,
	StatusApprovalPending,
	StatusOnHold,
	StatusApproved,
	StatusPublished,
	StatusRejected,
}

func (e Status) IsValid() bool {
	switch e {
	case StatusSaved, StatusApprovalPending, StatusOnHold, StatusApproved, StatusPublished, StatusRejected:
		return true
	}
	return false
}

func (e Status) String() string {
	return string(e)
}

func (e *Status) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Status(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Status", str)
	}
	return nil
}

func (e Status) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Type string

const (
	TypeScorm  Type = "SCORM"
	TypeTincan Type = "TINCAN"
	TypeCmi5   Type = "CMI5"
	TypeHTML5  Type = "HTML5"
)

var AllType = []Type{
	TypeScorm,
	TypeTincan,
	TypeCmi5,
	TypeHTML5,
}

func (e Type) IsValid() bool {
	switch e {
	case TypeScorm, TypeTincan, TypeCmi5, TypeHTML5:
		return true
	}
	return false
}

func (e Type) String() string {
	return string(e)
}

func (e *Type) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Type(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Type", str)
	}
	return nil
}

func (e Type) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
