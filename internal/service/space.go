package service

import (
	"context"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/store"
)

// NewSpaceService creates a new space service
func NewSpaceService(store store.UnPostStore) *SpaceService {
	return &SpaceService{store: store}
}

var _ v1.SpaceServiceServer = (*SpaceService)(nil)

type SpaceService struct {
	store store.UnPostStore
	v1.UnimplementedSpaceServiceServer
}

func (s *SpaceService) CreateSpace(ctx context.Context, request *v1.CreateSpaceRequest) (*v1.CreateSpaceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SpaceService) GetSpace(ctx context.Context, request *v1.GetSpaceRequest) (*v1.GetSpaceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SpaceService) ListSpace(ctx context.Context, request *v1.ListSpaceRequest) (*v1.ListSpaceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SpaceService) UpdateSpace(ctx context.Context, request *v1.UpdateSpaceRequest) (*v1.UpdateSpaceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SpaceService) DeleteSpace(ctx context.Context, request *v1.DeleteSpaceRequest) (*v1.DeleteSpaceResponse, error) {
	//TODO implement me
	panic("implement me")
}
