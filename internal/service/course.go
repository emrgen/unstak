package service

import (
	"context"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/store"
)

func NewCourseService(store store.UnPostStore) *CourseService {
	return &CourseService{
		store: store,
	}

}

var _ v1.CourseServiceServer = new(CourseService)

type CourseService struct {
	store store.UnPostStore
	v1.UnimplementedCourseServiceServer
}

func (c *CourseService) CreateCourse(ctx context.Context, request *v1.CreateCourseRequest) (*v1.CreateCourseResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CourseService) GetCourse(ctx context.Context, request *v1.GetCourseRequest) (*v1.GetCourseResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CourseService) ListCourse(ctx context.Context, request *v1.ListCourseRequest) (*v1.ListCourseResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CourseService) UpdateCourse(ctx context.Context, request *v1.UpdateCourseRequest) (*v1.UpdateCourseResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CourseService) DeleteCourse(ctx context.Context, request *v1.DeleteCourseRequest) (*v1.DeleteCourseResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CourseService) AddCourseTag(ctx context.Context, request *v1.AddCourseTagRequest) (*v1.AddCourseTagResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CourseService) RemoveCourseTag(ctx context.Context, request *v1.RemoveCourseTagRequest) (*v1.RemoveCourseTagResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CourseService) mustEmbedUnimplementedCourseServiceServer() {
	//TODO implement me
	panic("implement me")
}
