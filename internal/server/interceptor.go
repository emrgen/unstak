package server

import (
	"context"
	"errors"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/x"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
	"time"
)

var (
	// NoAuthHeaderError is returned when the authorization header is not found
	NoAuthHeaderError = errors.New("no auth header found")
)

// VerifyTokenInterceptor is a server interceptor that verifies the jwt token for each RPC call.
func VerifyTokenInterceptor(jwtSecret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		switch info.FullMethod {
		case v1.AccountService_CreateAccount_FullMethodName:
			return handler(ctx, req)
		default:
			return tokenInterceptor(ctx, jwtSecret, req, info, handler)
		}
	}
}

func tokenInterceptor(ctx context.Context, jwtSecret string, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	jwtToken, err := TokenFromHeader(ctx, "Bearer")
	if err != nil {
		logrus.Errorf("authbase: interceptor error getting token from header: %v", err)
		return nil, err
	}
	if len(jwtToken) == 0 {
		return nil, errors.New("token is empty")
	}

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid token")
	}

	ctx = x.ContextWithUserID(ctx, userID)

	return handler(ctx, req)
}

func TokenFromHeader(ctx context.Context, expectedScheme string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata not found")
	}

	val, ok := md["authorization"]
	if !ok {
		return "", NoAuthHeaderError
	}

	if len(val) == 0 {
		return "", errors.New("no token found")
	}

	scheme, token, found := strings.Cut(val[0], " ")
	if !found {
		return "", errors.New("bad authorization string")
	}

	if !strings.EqualFold(scheme, expectedScheme) {
		return "", errors.New("request unauthenticated with " + expectedScheme)
	}

	return token, nil
}

func UnaryGrpcRequestTimeInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		reqTime := time.Since(start)
		logrus.Infof("request time: %v: %v", info.FullMethod, reqTime)
		return resp, err
	}
}

func UnaryRequestTimeInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		reqTime := time.Since(start)
		logrus.Infof("request time: %v: %v", method, reqTime)
		return err
	}
}
