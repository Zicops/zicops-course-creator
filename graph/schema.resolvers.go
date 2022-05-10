package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/zicops/zicops-course-creator/graph/generated"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/handlers"
)

func (r *mutationResolver) AddCategories(ctx context.Context, category []*string) (*bool, error) {
	resp, err := handlers.AddCategory(ctx, category)
	if err != nil {
		log.Errorf("error adding categotries: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) AddSubCategories(ctx context.Context, subCategory []*string) (*bool, error) {
	resp, err := handlers.AddSubCategory(ctx, subCategory)
	if err != nil {
		log.Errorf("error adding sub categotries: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) AddCourse(ctx context.Context, course *model.CourseInput) (*model.Course, error) {
	resp, err := handlers.CourseCreator(ctx, course)
	if err != nil {
		log.Errorf("error creating course: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UpdateCourse(ctx context.Context, course *model.CourseInput) (*model.Course, error) {
	resp, err := handlers.CourseUpdate(ctx, course)
	if err != nil {
		log.Errorf("error creating course: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UploadCourseImage(ctx context.Context, file *model.CourseFile) (*model.UploadResult, error) {
	resp, err := handlers.UploadCourseImage(ctx, *file)
	if err != nil {
		log.Errorf("error upload course image: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UploadCoursePreviewVideo(ctx context.Context, file *model.CourseFile) (*model.UploadResult, error) {
	resp, err := handlers.UploadCoursePreviewVideo(ctx, *file)
	if err != nil {
		log.Errorf("error upload course video: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UploadCourseTileImage(ctx context.Context, file *model.CourseFile) (*model.UploadResult, error) {
	resp, err := handlers.UploadCourseTileImage(ctx, *file)
	if err != nil {
		log.Errorf("error upload course tile image: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) AddCourseModule(ctx context.Context, courseID *string, module *model.ModuleInput) (*model.Module, error) {
	resp, err := handlers.ModuleCreate(ctx, *courseID, module)
	if err != nil {
		log.Errorf("error creating module: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UpdateCourseModule(ctx context.Context, module *model.ModuleInput) (*model.Module, error) {
	resp, err := handlers.UpdateModule(ctx, module)
	if err != nil {
		log.Errorf("error updating module: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) AddCourseChapter(ctx context.Context, courseID *string, chapter *model.ChapterInput) (*model.Chapter, error) {
	resp, err := handlers.ChapterCreate(ctx, *courseID, chapter)
	if err != nil {
		log.Errorf("error creating chapter: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UpdateCourseChapter(ctx context.Context, chapter *model.ChapterInput) (*model.Chapter, error) {
	resp, err := handlers.UpdateChapter(ctx, chapter)
	if err != nil {
		log.Errorf("error updating chapter: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) AddCourseTopic(ctx context.Context, courseID *string, topic *model.TopicInput) (*model.Topic, error) {
	resp, err := handlers.TopicCreate(ctx, *courseID, topic)
	if err != nil {
		log.Errorf("error creating topic: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UpdateCourseTopic(ctx context.Context, topic *model.TopicInput) (*model.Topic, error) {
	resp, err := handlers.TopicUpdate(ctx, topic)
	if err != nil {
		log.Errorf("error updating topic: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) AddTopicContent(ctx context.Context, topicID *string, courseID *string, topicContent *model.TopicContentInput) (*model.TopicContent, error) {
	resp, err := handlers.TopicContentCreate(ctx, *topicID, *courseID, topicContent)
	if err != nil {
		log.Errorf("error creating topic content: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UpdateTopicContent(ctx context.Context, topicContent *model.TopicContentInput) (*model.TopicContent, error) {
	resp, err := handlers.UpdateTopicContent(ctx, topicContent)
	if err != nil {
		log.Errorf("error updating topic content: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UploadTopicContentVideo(ctx context.Context, file *model.TopicVideo) (*model.UploadResult, error) {
	resp, err := handlers.UploadTopicVideo(ctx, *file)
	if err != nil {
		log.Errorf("error uploading topic video: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UploadTopicContentSubtitle(ctx context.Context, file []*model.TopicSubtitle) ([]*model.UploadResultSubtitles, error) {
	resp, err := handlers.UploadTopicSubtitle(ctx, file)
	if err != nil {
		log.Errorf("error uploading topic subtitle: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UploadTopicStaticContent(ctx context.Context, file *model.StaticContent) (*model.UploadResult, error) {
	resp, err := handlers.UploadTopicStaticContent(ctx, file)
	if err != nil {
		log.Errorf("error uploading topic subtitle: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) AddQuiz(ctx context.Context, quiz *model.QuizInput) (*model.Quiz, error) {
	resp, err := handlers.CreateTopicQuiz(ctx, quiz)
	if err != nil {
		log.Errorf("error creating quiz: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UpdateQuiz(ctx context.Context, quiz *model.QuizInput) (*model.Quiz, error) {
	resp, err := handlers.UpdateQuiz(ctx, quiz)
	if err != nil {
		log.Errorf("error updating quiz: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UploadQuizFile(ctx context.Context, courseID *string, file *model.QuizFile) (*model.UploadResult, error) {
	resp, err := handlers.UploadQuizFile(ctx, *courseID, *file)
	if err != nil {
		log.Errorf("error uploading quiz file: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) AddQuizMcq(ctx context.Context, quiz *model.QuizMcq) (*bool, error) {
	resp, err := handlers.AddMCQQuiz(ctx, quiz)
	if err != nil {
		log.Errorf("error adding mcq quiz: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) AddQuizDescriptive(ctx context.Context, quiz *model.QuizDescriptive) (*bool, error) {
	resp, err := handlers.AddQuizDescriptive(ctx, quiz)
	if err != nil {
		log.Errorf("error adding descriptive quiz: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UploadTopicResource(ctx context.Context, courseID *string, resource *model.TopicResourceInput) (*model.UploadResult, error) {
	resp, err := handlers.AddTopicResources(ctx, *courseID, resource)
	if err != nil {
		log.Errorf("error adding topic resource: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) CreateQuestionBank(ctx context.Context, input *model.QuestionBankInput) (*model.QuestionBank, error) {
	resp, err := handlers.QuestionBankCreate(ctx, input)
	if err != nil {
		log.Errorf("error adding question bank: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UpdateQuestionBank(ctx context.Context, input *model.QuestionBankInput) (*model.QuestionBank, error) {
	resp, err := handlers.QuestionBankUpdate(ctx, input)
	if err != nil {
		log.Errorf("error updating question bank: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) AddQuestionBankQuestion(ctx context.Context, input *model.QuestionBankQuestionInput) (*model.QuestionBankQuestion, error) {
	resp, err := handlers.AddQuestionBankQuestion(ctx, input)
	if err != nil {
		log.Errorf("error adding question: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UpdateQuestionBankQuestion(ctx context.Context, input *model.QuestionBankQuestionInput) (*model.QuestionBankQuestion, error) {
	resp, err := handlers.UpdateQuestionBankQuestion(ctx, input)
	if err != nil {
		log.Errorf("error updating question: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) AddQuestionOptions(ctx context.Context, input *model.QuestionOptionInput) (*model.QuestionOption, error) {
	resp, err := handlers.AddQuestionOptions(ctx, input)
	if err != nil {
		log.Errorf("error adding question option: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UpdateQuestionOptions(ctx context.Context, input *model.QuestionOptionInput) (*model.QuestionOption, error) {
	resp, err := handlers.UpdateQuestionOptions(ctx, input)
	if err != nil {
		log.Errorf("error updating question option: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) AddQuestionPaper(ctx context.Context, input *model.QuestionPaperInput) (*model.QuestionPaper, error) {
	resp, err := handlers.QuestionPaperCreate(ctx, input)
	if err != nil {
		log.Errorf("error updating question option: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UpdateQuestionPaper(ctx context.Context, input *model.QuestionPaperInput) (*model.QuestionPaper, error) {
	resp, err := handlers.QuestionPaperUpdate(ctx, input)
	if err != nil {
		log.Errorf("error updating question option: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) AddQuestionPaperSection(ctx context.Context, input *model.QuestionPaperSectionInput) (*model.QuestionPaperSection, error) {
	resp, err := handlers.QuestionPaperSectionCreate(ctx, input)
	if err != nil {
		log.Errorf("error creating question section: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UpdateQuestionPaperSection(ctx context.Context, input *model.QuestionPaperSectionInput) (*model.QuestionPaperSection, error) {
	resp, err := handlers.QuestionPaperSectionCreate(ctx, input)
	if err != nil {
		log.Errorf("error updating question paper section: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) MapSectionToBank(ctx context.Context, input *model.MapSectionToBankInput) (*model.SectionQBMapping, error) {
	resp, err := handlers.QuestionSectionMap(ctx, input)
	if err != nil {
		log.Errorf("error creating mapping: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UpdateSectionToBank(ctx context.Context, input *model.MapSectionToBankInput) (*model.SectionQBMapping, error) {
	resp, err := handlers.QuestionSectionMapUpdate(ctx, input)
	if err != nil {
		log.Errorf("error updating mapping: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) AddSectionFixedQuestions(ctx context.Context, input *model.SectionFixedQuestionsInput) (*model.SectionFixedQuestions, error) {
	resp, err := handlers.QuestionFixed(ctx, input)
	if err != nil {
		log.Errorf("error adding fixed : %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UpdateSectionFixedQuestions(ctx context.Context, input *model.SectionFixedQuestionsInput) (*model.SectionFixedQuestions, error) {
	resp, err := handlers.QuestionFixedUpdate(ctx, input)
	if err != nil {
		log.Errorf("error updating fixed : %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) AddExam(ctx context.Context, input *model.ExamInput) (*model.Exam, error) {
	resp, err := handlers.ExamCreate(ctx, input)
	if err != nil {
		log.Errorf("error adding exam : %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UpdateExam(ctx context.Context, input *model.ExamInput) (*model.Exam, error) {
	resp, err := handlers.ExamUpdate(ctx, input)
	if err != nil {
		log.Errorf("error updating exam : %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) AddExamSchedule(ctx context.Context, input *model.ExamScheduleInput) (*model.ExamSchedule, error) {
	resp, err := handlers.ExamScheduleCreate(ctx, input)
	if err != nil {
		log.Errorf("error adding exam schedule : %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UpdateExamSchedule(ctx context.Context, input *model.ExamScheduleInput) (*model.ExamSchedule, error) {
	resp, err := handlers.ExamScheduleUpdate(ctx, input)
	if err != nil {
		log.Errorf("error updating exam schedule: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) AddExamInstruction(ctx context.Context, input *model.ExamInstructionInput) (*model.ExamInstruction, error) {
	resp, err := handlers.ExamInstructionsCreate(ctx, input)
	if err != nil {
		log.Errorf("error adding exam ins : %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UpdateExamInstruction(ctx context.Context, input *model.ExamInstructionInput) (*model.ExamInstruction, error) {
	resp, err := handlers.ExamInstructionsUpdate(ctx, input)
	if err != nil {
		log.Errorf("error updating exam ins : %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) AddExamCohort(ctx context.Context, input *model.ExamCohortInput) (*model.ExamCohort, error) {
	resp, err := handlers.AddExamCohort(ctx, input)
	if err != nil {
		log.Errorf("error adding exam cohort : %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UpdateExamCohort(ctx context.Context, input *model.ExamCohortInput) (*model.ExamCohort, error) {
	resp, err := handlers.UpdateExamCohort(ctx, input)
	if err != nil {
		log.Errorf("error updating exam cohort : %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) AddExamConfiguration(ctx context.Context, input *model.ExamConfigurationInput) (*model.ExamConfiguration, error) {
	resp, err := handlers.AddExamConfiguration(ctx, input)
	if err != nil {
		log.Errorf("error adding exam config : %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UpdateExamConfiguration(ctx context.Context, input *model.ExamConfigurationInput) (*model.ExamConfiguration, error) {
	resp, err := handlers.UpdateExamConfiguration(ctx, input)
	if err != nil {
		log.Errorf("error updating exam config : %v", err)
		return nil, err
	}
	return resp, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
