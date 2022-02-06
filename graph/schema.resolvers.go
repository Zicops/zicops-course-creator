package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/zicops/zicops-course-creator/graph/generated"
	"github.com/zicops/zicops-course-creator/graph/model"
	"github.com/zicops/zicops-course-creator/handlers"
)

func (r *mutationResolver) AddCatergories(ctx context.Context, category []*string) (*bool, error) {
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

func (r *mutationResolver) UploadCourseImage(ctx context.Context, file model.CourseFile) (*bool, error) {
	resp, err := handlers.UploadCourseImage(ctx, file)
	if err != nil {
		log.Errorf("error upload course image: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UploadCoursePreviewVideo(ctx context.Context, file model.CourseFile) (*bool, error) {
	resp, err := handlers.UploadCoursePreviewVideo(ctx, file)
	if err != nil {
		log.Errorf("error upload course video: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UploadCourseTileImage(ctx context.Context, file model.CourseFile) (*bool, error) {
	resp, err := handlers.UploadCourseTileImage(ctx, file)
	if err != nil {
		log.Errorf("error upload course tile image: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) AddCourseModule(ctx context.Context, courseID string, module *model.ModuleInput) (*model.Module, error) {
	resp, err := handlers.ModuleCreate(ctx, courseID, module)
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

func (r *mutationResolver) AddCourseChapter(ctx context.Context, courseID string, chapter *model.ChapterInput) (*model.Chapter, error) {
	resp, err := handlers.ChapterCreate(ctx, courseID, chapter)
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

func (r *mutationResolver) AddCourseTopic(ctx context.Context, courseID string, topic *model.TopicInput) (*model.Topic, error) {
	resp, err := handlers.TopicCreate(ctx, courseID, topic)
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

func (r *mutationResolver) AddTopicContent(ctx context.Context, topicID string, topicContent *model.TopicContentInput) (*model.TopicContent, error) {
	resp, err := handlers.TopicContentCreate(ctx, topicID, topicContent)
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

func (r *mutationResolver) UploadTopicContentVideo(ctx context.Context, file model.TopicVideo) (*bool, error) {
	resp, err := handlers.UploadTopicVideo(ctx, file)
	if err != nil {
		log.Errorf("error uploading topic video: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UploadTopicContentSubtitle(ctx context.Context, file model.TopicSubtitle) (*bool, error) {
	resp, err := handlers.UploadTopicSubtitle(ctx, file)
	if err != nil {
		log.Errorf("error uploading topic subtitle: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UploadTopicStaticContent(ctx context.Context, file model.StaticContent) (*bool, error) {
	return nil, nil
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

func (r *mutationResolver) UploadQuizFile(ctx context.Context, courseID string, file model.QuizFile) (*bool, error) {
	resp, err := handlers.UploadQuizFile(ctx, courseID, file)
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

func (r *mutationResolver) UploadTopicResource(ctx context.Context, courseID string, resource *model.TopicResourceInput) (*bool, error) {
	resp, err := handlers.AddTopicResources(ctx, courseID, resource)
	if err != nil {
		log.Errorf("error adding topic resource: %v", err)
		return nil, err
	}
	return resp, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
