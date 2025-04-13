package service

import (
	"context"
	"errors"
	"github.com/emrgen/authbase"
	authx "github.com/emrgen/authbase/x"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/model"
	"github.com/emrgen/unpost/internal/store"
	"github.com/google/uuid"
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
	accountID, err := authx.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
	}
	spaceID := uuid.MustParse(request.GetSpaceId())
	if spaceID == uuid.Nil {
		return nil, errors.New("space id is required")
	}
	member, err := s.store.GetSpaceMember(ctx, accountID, spaceID)
	if err != nil {
		return nil, err
	}

	if member != nil && member.Role != model.UserRoleAdmin && member.Role != model.UserRoleOwner {
		return nil, errors.New("only admin or owner can add space member")
	}

	role := model.UserRoleViewer
	switch request.GetRole() {
	case v1.UserRole_Admin:
		role = model.UserRoleAdmin
	case v1.UserRole_Editor:
		role = model.UserRoleEditor
	case v1.UserRole_Contributor:
		role = model.UserRoleContributor
	}

	newMember := &model.SpaceMember{
		UserID:  request.GetUserId(),
		SpaceID: request.GetSpaceId(),
		Role:    role,
	}

	err = s.store.UpdateSpaceMember(ctx, newMember)
	if err != nil {
		return nil, err
	}

	return &v1.AddSpaceMemberResponse{
		SpaceMember: &v1.SpaceMember{
			SpaceId: request.GetSpaceId(),
			UserId:  request.GetUserId(),
			Role:    request.GetRole(),
		},
	}, nil
}

func (s *SpaceMemberService) GetSpaceMember(ctx context.Context, request *v1.GetSpaceMemberRequest) (*v1.GetSpaceMemberResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SpaceMemberService) ListSpaceMember(ctx context.Context, request *v1.ListSpaceMemberRequest) (*v1.ListSpaceMemberResponse, error) {
	accountID, err := authx.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
	}

	spaceID := uuid.MustParse(request.GetSpaceId())
	if spaceID == uuid.Nil {
		return nil, errors.New("space id is required")
	}

	member, err := s.store.GetSpaceMember(ctx, accountID, spaceID)
	if err != nil {
		return nil, err
	}

	if member == nil {
		return nil, errors.New("user is not a member of the space")
	}

	if member.Role != model.UserRoleAdmin && member.Role != model.UserRoleOwner {
		return nil, errors.New("only admin or owner can list space members")
	}

	members, err := s.store.ListSpaceMembers(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	var responseMembers []*v1.SpaceMember
	for _, member := range members {
		responseMembers = append(responseMembers, &v1.SpaceMember{
			SpaceId: member.SpaceID,
			UserId:  member.UserID,
		})
	}

	return &v1.ListSpaceMemberResponse{
		SpaceMembers: responseMembers,
	}, nil
}

func (s *SpaceMemberService) UpdateSpaceMember(ctx context.Context, request *v1.UpdateSpaceMemberRequest) (*v1.UpdateSpaceMemberResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SpaceMemberService) DeleteSpaceMember(ctx context.Context, request *v1.DeleteSpaceMemberRequest) (*v1.DeleteSpaceMemberResponse, error) {
	//TODO implement me
	panic("implement me")
}
