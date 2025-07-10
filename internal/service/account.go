package service

import (
	"context"
	"errors"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/store"
	"github.com/emrgen/unpost/internal/x"
	"github.com/sirupsen/logrus"
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

func (a *AccountService) Logout(ctx context.Context, request *v1.LogoutRequest) (*v1.LogoutResponse, error) {
	token, _ := x.TokenFromContext(ctx)
	logrus.Printf("token: %s", token)
	err := a.authClient.WithToken(token).Logout()
	if err != nil {
		return nil, err
	}

	return &v1.LogoutResponse{}, nil
}

func (a *AccountService) LoginUsingPassword(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	token, err := a.authClient.Token(types.TokenRequest{
		GrantType: "password",
		Email:     email,
		Password:  password,
	})
	if err != nil {
		return nil, err
	}

	user := token.User
	var userRole string
	if role, ok := token.User.AppMetadata["role"]; ok {
		userRole = role.(string)
	} else {
		return nil, errors.New("role not found in token")
	}

	var role *v1.UserRole
	switch userRole {
	case "viewer":
		role = v1.UserRole_Viewer.Enum()
	case "author":
		role = v1.UserRole_Author.Enum()
	case "admin":
		role = v1.UserRole_Admin.Enum()
	case "owner":
		role = v1.UserRole_Owner.Enum()
	}

	if role == nil {
		return nil, errors.New("role not found in token")
	}

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
			Role:  *role,
		},
	}, nil
}
