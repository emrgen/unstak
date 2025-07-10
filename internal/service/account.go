package service

import (
	"context"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/store"
	"github.com/supabase-community/auth-go"
	"github.com/supabase-community/auth-go/types"
)

// NewAccountService creates a new auth service
func NewAccountService(store store.UnstakStore, authClient auth.Client) *AccountService {
	return &AccountService{
		store:      store,
		authClient: authClient,
	}
}

var (
	_ v1.AccountServiceServer = new(AccountService)
)

type AccountService struct {
	store      store.UnstakStore
	authClient auth.Client
	v1.UnimplementedAccountServiceServer
}

func (a *AccountService) CreateAccount(ctx context.Context, req *v1.CreateAccountRequest) (*v1.CreateAccountResponse, error) {
	signup, err := a.authClient.Signup(types.SignupRequest{
		Email:         req.Email,
		Password:      req.Password,
		SecurityEmbed: types.SecurityEmbed{},
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateAccountResponse{
		Account: &v1.Account{
			Id: signup.ID.String(),
		},
	}, nil
}

func (a *AccountService) LoginUsingPassword(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	token, err := a.authClient.Token(types.TokenRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	user := token.User

	return &v1.LoginResponse{
		Token: &v1.AuthToken{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			TokenType:    token.TokenType,
			ExpiresIn:    int32(token.ExpiresIn),
			ExpiresAt:    int32(token.ExpiresAt),
		},
		Account: &v1.Account{
			Id:    user.ID.String(),
			Email: user.Email,
		},
	}, nil
}
