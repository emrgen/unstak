package service

import (
	"context"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/model"
	"github.com/emrgen/unpost/internal/store"
	"github.com/google/uuid"
)

func NewSubscriptionMemberService() v1.SubscriptionMemberServiceServer {
	return &SubscriptionMemberService{}

}

var _ v1.SubscriptionMemberServiceServer = (*SubscriptionMemberService)(nil)

type SubscriptionMemberService struct {
	store store.UnPostStore
	v1.UnimplementedSubscriptionMemberServiceServer
}

func (s *SubscriptionMemberService) CreateSubscriptionMember(ctx context.Context, request *v1.CreateSubscriptionMemberRequest) (*v1.CreateSubscriptionMemberResponse, error) {
	userID := uuid.MustParse(request.GetUserId())
	subscriptionID := uuid.MustParse(request.GetSubscriptionId())

	err := s.store.AddSubscriptionMember(ctx, &model.SubscriptionMember{
		UserID:         userID.String(),
		SubscriptionID: subscriptionID.String(),
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateSubscriptionMemberResponse{}, nil
}

func (s *SubscriptionMemberService) GetSubscriptionMember(ctx context.Context, request *v1.GetSubscriptionMemberRequest) (*v1.GetSubscriptionMemberResponse, error) {
	subMemberID := uuid.MustParse(request.GetId())

	subscriptionMember, err := s.store.GetSubscriptionMember(ctx, subMemberID)
	if err != nil {
		return nil, err
	}

	return &v1.GetSubscriptionMemberResponse{
		Member: &v1.SubscriptionMember{
			Id:             subscriptionMember.ID,
			UserId:         subscriptionMember.UserID,
			SubscriptionId: subscriptionMember.SubscriptionID,
			Subscription: &v1.Subscription{
				Id:   subscriptionMember.Subscription.ID,
				Name: subscriptionMember.Subscription.Name,
			},
		},
	}, nil
}

func (s *SubscriptionMemberService) ListSubscriptionMember(ctx context.Context, request *v1.ListSubscriptionMemberRequest) (*v1.ListSubscriptionMemberResponse, error) {
	subID := uuid.MustParse(request.GetSubscriptionId())

	subscriptionMembers, err := s.store.ListSubscriptionMembers(ctx, subID)
	if err != nil {
		return nil, err
	}

	members := make([]*v1.SubscriptionMember, 0, len(subscriptionMembers))
	for _, member := range subscriptionMembers {

		members = append(members, &v1.SubscriptionMember{
			Id:             member.ID,
			UserId:         member.UserID,
			SubscriptionId: member.SubscriptionID,
			Subscription: &v1.Subscription{
				Id:   member.Subscription.ID,
				Name: member.Subscription.Name,
			},
		})
	}

	return &v1.ListSubscriptionMemberResponse{Members: members}, nil
}

func (s *SubscriptionMemberService) UpdateSubscriptionMember(ctx context.Context, request *v1.UpdateSubscriptionMemberRequest) (*v1.UpdateSubscriptionMemberResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SubscriptionMemberService) DeleteSubscriptionMember(ctx context.Context, request *v1.DeleteSubscriptionMemberRequest) (*v1.DeleteSubscriptionMemberResponse, error) {
	subID := uuid.MustParse(request.GetId())

	if err := s.store.RemoveSubscriptionMember(ctx, subID); err != nil {
		return nil, err
	}

	return &v1.DeleteSubscriptionMemberResponse{}, nil
}
