package service

import (
	"context"
	v1 "github.com/emrgen/unpost/apis/v1"
)

// NewPageService creates a new book service
func NewPageService() *PageService {
	return &PageService{}
}

var _ v1.PageServiceServer = new(PageService)

type PageService struct {
	v1.UnimplementedPageServiceServer
}

func (p *PageService) CreatePage(ctx context.Context, request *v1.CreatePageRequest) (*v1.CreatePageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PageService) GetPage(ctx context.Context, request *v1.GetPageRequest) (*v1.GetPageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PageService) ListPage(ctx context.Context, request *v1.ListPageRequest) (*v1.ListPageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PageService) UpdatePage(ctx context.Context, request *v1.UpdatePageRequest) (*v1.UpdatePageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PageService) DeletePage(ctx context.Context, request *v1.DeletePageRequest) (*v1.DeletePageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PageService) AddPageTag(ctx context.Context, request *v1.AddPageTagRequest) (*v1.AddPageTagResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PageService) RemovePageTag(ctx context.Context, request *v1.RemovePageTagRequest) (*v1.RemovePageTagResponse, error) {
	//TODO implement me
	panic("implement me")
}
