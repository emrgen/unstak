package service

import (
	"context"
	authx "github.com/emrgen/authbase/x"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/model"
	"github.com/emrgen/unpost/internal/store"
	"github.com/google/uuid"
)

func NewTierService(store store.UnPostStore) v1.TierServiceServer {
	return &TierService{
		store: store,
	}
}

var _ v1.TierServiceServer = (*TierService)(nil)

type TierService struct {
	store store.UnPostStore
	v1.UnimplementedTierServiceServer
}

func (s *TierService) CreateTier(ctx context.Context, request *v1.CreateTierRequest) (*v1.CreateTierResponse, error) {
	var err error
	userID, err := authx.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
	}

	tier := &model.Tier{
		Name:        request.GetName(),
		CreatedByID: userID.String(),
	}

	err = s.store.CreateTier(ctx, tier)
	if err != nil {
		return nil, err
	}

	return &v1.CreateTierResponse{Tier: &v1.Tier{
		Id: tier.ID,
	}}, nil
}

func (s *TierService) GetTier(ctx context.Context, request *v1.GetTierRequest) (*v1.GetTierResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *TierService) ListTiers(ctx context.Context, request *v1.ListTiersRequest) (*v1.ListTiersResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *TierService) UpdateTier(ctx context.Context, request *v1.UpdateTierRequest) (*v1.UpdateTierResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *TierService) DeleteTier(ctx context.Context, request *v1.DeleteTierRequest) (*v1.DeleteTierResponse, error) {
	subID := uuid.MustParse(request.GetId())

	if err := s.store.DeleteTier(ctx, subID); err != nil {
		return nil, err
	}

	return &v1.DeleteTierResponse{}, nil
}
