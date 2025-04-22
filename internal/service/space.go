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
	"github.com/sirupsen/logrus"
)

// NewSpaceService creates a new space service
func NewSpaceService(cfg *authx.AuthbaseConfig, store store.UnstakStore, authClient authbase.Client) *SpaceService {
	return &SpaceService{cfg: cfg, store: store, authClient: authClient}
}

var _ v1.SpaceServiceServer = (*SpaceService)(nil)

type SpaceService struct {
	cfg        *authx.AuthbaseConfig
	store      store.UnstakStore
	authClient authbase.Client
	v1.UnimplementedSpaceServiceServer
}

func (s *SpaceService) CreateSpace(ctx context.Context, request *v1.CreateSpaceRequest) (*v1.CreateSpaceResponse, error) {
	accountID, err := authx.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
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
		OwnerID: accountID.String(),
		Private: request.GetPrivate(),
	}

	if poolID != uuid.Nil {
		space.PoolID = poolID.String()
	}

	member := &model.SpaceMember{
		SpaceID: space.ID,
		UserID:  accountID.String(),
		Role:    model.UserRoleOwner,
	}

	err = s.store.Transaction(ctx, func(ctx context.Context, tx store.UnstakStore) error {
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
	accountID, err := authx.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
	}

	masterSpace, err := s.store.GetMasterSpace(ctx)
	if err != nil {
		return nil, err
	}

	logrus.Infof("accountID: %s", accountID.String())
	spaces, err := s.store.ListSpaces(ctx, accountID)
	if err != nil {
		return nil, err
	}

	var res []*v1.Space
	res = append(res, &v1.Space{
		Id:   masterSpace.ID,
		Name: masterSpace.Name,
	})
	for _, space := range spaces {
		res = append(res, &v1.Space{
			Id:   space.ID,
			Name: space.Name,
		})
	}

	return &v1.ListSpaceResponse{
		Spaces: res,
	}, nil
}

func (s *SpaceService) UpdateSpace(ctx context.Context, request *v1.UpdateSpaceRequest) (*v1.UpdateSpaceResponse, error) {
	accountID, err := authx.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
	}

	spaceID, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}
	space, err := s.store.GetSpace(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	// check if user is in a space member with admin role
	member, err := s.store.GetSpaceMember(ctx, spaceID, accountID)
	if err != nil {
		return nil, err
	}

	if member.Role != model.UserRoleAdmin {
		return nil, errors.New("only admin can update space")
	}

	// update space
	if request.Name != nil {
		space.Name = request.GetName()
	}
	err = s.store.UpdateSpace(ctx, space)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateSpaceResponse{
		Space: &v1.Space{
			Id:          space.ID,
			Name:        space.Name,
			Description: "",
			Thumbnail:   "",
			Private:     false,
		},
	}, nil
}

func (s *SpaceService) DeleteSpace(ctx context.Context, request *v1.DeleteSpaceRequest) (*v1.DeleteSpaceResponse, error) {
	//TODO implement me
	panic("implement me")
}
