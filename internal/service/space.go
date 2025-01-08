package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/emrgen/authbase"
	authv1 "github.com/emrgen/authbase/apis/v1"
	authx "github.com/emrgen/authbase/x"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/model"
	"github.com/emrgen/unpost/internal/store"
	"github.com/emrgen/unpost/internal/x"
	"github.com/google/uuid"
)

// NewSpaceService creates a new space service
func NewSpaceService(cfg *authx.AuthbaseConfig, store store.UnPostStore, authClient authbase.Client) *SpaceService {
	return &SpaceService{cfg: cfg, store: store, authClient: authClient}
}

var _ v1.SpaceServiceServer = (*SpaceService)(nil)

type SpaceService struct {
	cfg        *authx.AuthbaseConfig
	store      store.UnPostStore
	authClient authbase.Client
	v1.UnimplementedSpaceServiceServer
}

func (s *SpaceService) CreateSpace(ctx context.Context, request *v1.CreateSpaceRequest) (*v1.CreateSpaceResponse, error) {
	userID, err := authx.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := s.store.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !user.Space.Master {
		return nil, errors.New("only user of master space can create a new space")
	}

	var poolID = uuid.Nil
	// create a new pool in authbase project
	if request.GetPoolName() != "" {
		projectID, err := authx.GetAuthbaseProjectID(ctx)
		if err != nil {
			return nil, err
		}

		res, err := s.authClient.CreatePool(s.cfg.IntoContext(), &authv1.CreatePoolRequest{
			ProjectId: projectID.String(),
			Name:      fmt.Sprintf("%s-%s", request.GetPoolName(), x.RandomString(6)),
		})
		if err != nil {
			return nil, err
		}
		poolID = uuid.MustParse(res.Pool.Id)
	}

	space := &model.Space{
		ID:      uuid.New().String(),
		Name:    request.GetName(),
		OwnerID: userID.String(),
		Private: request.GetPrivate(),
	}

	if poolID != uuid.Nil {
		space.PoolID = poolID.String()
	}

	member := &model.SpaceMember{
		SpaceID: space.ID,
		UserID:  userID.String(),
		Role:    model.UserRoleOwner,
	}

	err = s.store.Transaction(ctx, func(ctx context.Context, tx store.UnPostStore) error {
		if err := tx.CreateSpace(ctx, space); err != nil {
			return err
		}

		if err := tx.AddSpaceMember(ctx, member); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateSpaceResponse{
		Space: &v1.Space{
			Id:   space.ID,
			Name: space.Name,
		},
	}, nil
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
