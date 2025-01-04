package service

import (
	"context"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/model"
	"github.com/emrgen/unpost/internal/store"
	"github.com/google/uuid"
)

func NewTierMemberService() v1.TierMemberServiceServer {
	return &TierMemberService{}

}

var _ v1.TierMemberServiceServer = (*TierMemberService)(nil)

type TierMemberService struct {
	store store.UnPostStore
	v1.UnimplementedTierMemberServiceServer
}

func (s *TierMemberService) CreateTierMember(ctx context.Context, request *v1.CreateTierMemberRequest) (*v1.CreateTierMemberResponse, error) {
	userID := uuid.MustParse(request.GetUserId())
	tierID := uuid.MustParse(request.GetTierId())

	err := s.store.AddTierMember(ctx, &model.TierMember{
		UserID: userID.String(),
		TierID: tierID.String(),
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateTierMemberResponse{}, nil
}

func (s *TierMemberService) GetTierMember(ctx context.Context, request *v1.GetTierMemberRequest) (*v1.GetTierMemberResponse, error) {
	subMemberID := uuid.MustParse(request.GetId())

	tierMember, err := s.store.GetTierMember(ctx, subMemberID)
	if err != nil {
		return nil, err
	}

	return &v1.GetTierMemberResponse{
		Member: &v1.TierMember{
			Id:     tierMember.ID,
			UserId: tierMember.UserID,
			TierId: tierMember.TierID,
			Tier: &v1.Tier{
				Id:   tierMember.Tier.ID,
				Name: tierMember.Tier.Name,
			},
		},
	}, nil
}

func (s *TierMemberService) ListTierMember(ctx context.Context, request *v1.ListTierMemberRequest) (*v1.ListTierMemberResponse, error) {
	subID := uuid.MustParse(request.GetTierId())

	tierMembers, err := s.store.ListTierMembers(ctx, subID)
	if err != nil {
		return nil, err
	}

	members := make([]*v1.TierMember, 0, len(tierMembers))
	for _, member := range tierMembers {

		members = append(members, &v1.TierMember{
			Id:     member.ID,
			UserId: member.UserID,
			TierId: member.TierID,
			Tier: &v1.Tier{
				Id:   member.Tier.ID,
				Name: member.Tier.Name,
			},
		})
	}

	return &v1.ListTierMemberResponse{Members: members}, nil
}

func (s *TierMemberService) UpdateTierMember(ctx context.Context, request *v1.UpdateTierMemberRequest) (*v1.UpdateTierMemberResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *TierMemberService) DeleteTierMember(ctx context.Context, request *v1.DeleteTierMemberRequest) (*v1.DeleteTierMemberResponse, error) {
	subID := uuid.MustParse(request.GetId())

	if err := s.store.RemoveTierMember(ctx, subID); err != nil {
		return nil, err
	}

	return &v1.DeleteTierMemberResponse{}, nil
}
