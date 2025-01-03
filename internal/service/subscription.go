package service

import (
	"context"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/store"
	"github.com/google/uuid"
)

func NewSubscriptionService() v1.SubscriptionServiceServer {
	return &SubscriptionService{}
}

var _ v1.SubscriptionServiceServer = (*SubscriptionService)(nil)

type SubscriptionService struct {
	store store.UnPostStore
	v1.UnimplementedSubscriptionServiceServer
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, request *v1.CreateSubscriptionRequest) (*v1.CreateSubscriptionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SubscriptionService) GetSubscription(ctx context.Context, request *v1.GetSubscriptionRequest) (*v1.GetSubscriptionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SubscriptionService) ListSubscriptions(ctx context.Context, request *v1.ListSubscriptionsRequest) (*v1.ListSubscriptionsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SubscriptionService) UpdateSubscription(ctx context.Context, request *v1.UpdateSubscriptionRequest) (*v1.UpdateSubscriptionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, request *v1.DeleteSubscriptionRequest) (*v1.DeleteSubscriptionResponse, error) {
	subID := uuid.MustParse(request.GetId())

	if err := s.store.DeleteSubscription(ctx, subID); err != nil {
		return nil, err
	}

	return &v1.DeleteSubscriptionResponse{}, nil
}
