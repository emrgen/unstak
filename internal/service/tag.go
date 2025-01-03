package service

import (
	"context"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/model"
	"github.com/emrgen/unpost/internal/store"
	"github.com/google/uuid"
)

// NewTagService creates a new tag service
func NewTagService(store store.TinyPostStore) *TagService {
	return &TagService{
		store: store,
	}
}

var _ v1.TagServiceServer = new(TagService)

// TagService is the service that provides tag operations
type TagService struct {
	store store.TinyPostStore
	v1.UnimplementedTagServiceServer
}

func (t *TagService) CreateTag(ctx context.Context, request *v1.CreateTagRequest) (*v1.CreateTagResponse, error) {
	tag := &model.Tag{
		ID:   uuid.New().String(),
		Name: request.GetName(),
	}

	err := t.store.CreateTag(ctx, tag)
	if err != nil {
		return nil, err
	}

	return &v1.CreateTagResponse{
		Tag: &v1.Tag{
			Id: tag.ID,
		},
	}, nil
}

func (t *TagService) GetTag(ctx context.Context, request *v1.GetTagRequest) (*v1.GetTagResponse, error) {
	tag, err := t.store.GetTag(ctx, uuid.MustParse(request.GetId()))
	if err != nil {
		return nil, err
	}

	return &v1.GetTagResponse{
		Tag: &v1.Tag{
			Id: tag.ID,
		},
	}, nil
}

func (t *TagService) ListTag(ctx context.Context, request *v1.ListTagRequest) (*v1.ListTagResponse, error) {

	tags, err := t.store.ListTags(ctx, 0, 100)
	if err != nil {
		return nil, err
	}

	var tagResponses []*v1.Tag
	for _, tag := range tags {
		tagResponses = append(tagResponses, &v1.Tag{
			Id:   tag.ID,
			Name: tag.Name,
		})
	}

	return &v1.ListTagResponse{
		Tags: tagResponses,
	}, nil
}

func (t *TagService) UpdateTag(ctx context.Context, request *v1.UpdateTagRequest) (*v1.UpdateTagResponse, error) {
	tag := &model.Tag{
		ID:   request.GetId(),
		Name: request.GetName(),
	}

	err := t.store.UpdateTag(ctx, tag)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateTagResponse{
		Tag: &v1.Tag{
			Id: tag.ID,
		},
	}, nil
}

func (t *TagService) DeleteTag(ctx context.Context, request *v1.DeleteTagRequest) (*v1.DeleteTagResponse, error) {
	err := t.store.DeleteTag(ctx, uuid.MustParse(request.GetId()))
	if err != nil {
		return nil, err
	}

	return &v1.DeleteTagResponse{}, nil
}
