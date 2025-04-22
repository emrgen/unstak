package service

import (
	"context"
	"github.com/emrgen/authbase"
	authv1 "github.com/emrgen/authbase/apis/v1"
	authx "github.com/emrgen/authbase/x"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/store"
)

// NewAccountService creates a new auth service
func NewAccountService(cfg *authx.AuthbaseConfig, store store.UnstakStore, authClient authbase.Client) *AccountService {
	return &AccountService{
		cfg:        cfg,
		store:      store,
		authClient: authClient,
	}
}

var (
	_ v1.AccountServiceServer = new(AccountService)
)

type AccountService struct {
	cfg        *authx.AuthbaseConfig
	store      store.UnstakStore
	authClient authbase.Client
	v1.UnimplementedAccountServiceServer
}

func (a *AccountService) CreateAccount(ctx context.Context, request *v1.CreateAccountRequest) (*v1.CreateAccountResponse, error) {
	poolID := request.GetPoolId()
	// if poolID is empty, get the pool id from the config client
	if poolID == "" {
		res, err := a.authClient.GetClient(a.cfg.IntoContext(), &authv1.GetClientRequest{
			ClientId: a.cfg.ClientID,
		})
		if err != nil {
			return nil, err
		}

		poolID = res.GetClient().GetPoolId()
	}
	res, err := a.authClient.CreateAccount(a.cfg.IntoContext(), &authv1.CreateAccountRequest{
		PoolId:   poolID,
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateAccountResponse{
		Account: &v1.Account{
			Id: res.Account.Id,
		},
	}, nil
}

func (a *AccountService) LoginUsingPassword(ctx context.Context, request *v1.LoginRequest) (*v1.LoginResponse, error) {
	email := request.GetEmail()
	password := request.GetPassword()

	clientID := request.GetClientId()
	if clientID == "" {
		clientID = a.cfg.ClientID
	}

	res, err := a.authClient.LoginUsingPassword(ctx, &authv1.LoginUsingPasswordRequest{
		Email:    email,
		Password: password,
		ClientId: clientID,
	})
	if err != nil {
		return nil, err
	}

	token := res.GetToken()
	account := res.GetAccount()

	return &v1.LoginResponse{
		Token: &v1.AuthToken{
			AccessToken:      token.AccessToken,
			RefreshToken:     token.RefreshToken,
			ExpiresAt:        token.ExpiresAt,
			IssuedAt:         token.IssuedAt,
			RefreshExpiresAt: token.RefreshExpiresAt,
		},
		Account: &v1.Account{
			Id:    account.Id,
			Email: account.Email,
			//FirstName: account.FirstName,
			//LastName:  account.LastName,
			//Username: account.Username,
		},
	}, nil
}

func (a *AccountService) CreateOwner(ctx context.Context, request *v1.CreateOwnerRequest) (*v1.CreateOwnerResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AccountService) CheckOwnerSetup(ctx context.Context, request *v1.CheckOwnerSetupRequest) (*v1.CheckOwnerSetupResponse, error) {
	//TODO implement me
	panic("implement me")
}
