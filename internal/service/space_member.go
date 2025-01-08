package service

import (
	"context"
	"github.com/emrgen/authbase"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/store"
)

func NewSpaceMemberService(store store.UnPostStore, authClient authbase.Client) *SpaceMemberService {
	return &SpaceMemberService{store: store, authClient: authClient}
}

var _ v1.SpaceMemberServiceServer = (*SpaceMemberService)(nil)

type SpaceMemberService struct {
	store      store.UnPostStore
	authClient authbase.Client
	v1.UnimplementedSpaceMemberServiceServer
}

func (s *SpaceMemberService) AddSpaceMember(ctx context.Context, request *v1.AddSpaceMemberRequest) (*v1.AddSpaceMemberResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SpaceMemberService) GetSpaceMember(ctx context.Context, request *v1.GetSpaceMemberRequest) (*v1.GetSpaceMemberResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SpaceMemberService) ListSpaceMember(ctx context.Context, request *v1.ListSpaceMemberRequest) (*v1.ListSpaceMemberResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SpaceMemberService) UpdateSpaceMember(ctx context.Context, request *v1.UpdateSpaceMemberRequest) (*v1.UpdateSpaceMemberResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SpaceMemberService) DeleteSpaceMember(ctx context.Context, request *v1.DeleteSpaceMemberRequest) (*v1.DeleteSpaceMemberResponse, error) {
	//TODO implement me
	panic("implement me")
}
