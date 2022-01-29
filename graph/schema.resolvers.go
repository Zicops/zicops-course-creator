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
	panic(fmt.Errorf("not implemented"))
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
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddTopicContent(ctx context.Context, topicID string, topicConent *model.TopicContentInput) (*model.TopicContent, error) {
	resp, err := handlers.TopicContentCreate(ctx, topicID, topicConent)
	if err != nil {
		log.Errorf("error creating topic content: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UpdateTopicContent(ctx context.Context, topicContent *model.TopicContentInput) (*model.TopicContent, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UploadTopicContentVideo(ctx context.Context, file model.TopicVideo) (*bool, error) {
	resp, err := handlers.UploadTopicVideo(ctx, file)
	if err != nil {
		log.Errorf("error creating topic content: %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UploadTopicStaticContent(ctx context.Context, file model.StaticContent) (*bool, error) {
	return nil, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
