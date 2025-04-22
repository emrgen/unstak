package server

import (
	"context"
	"errors"
	authv1 "github.com/emrgen/authbase/apis/v1"
	authx "github.com/emrgen/authbase/x"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

func VerifyTokenInterceptor(keyProvider authx.VerifierProvider, authClient authv1.AccessKeyServiceClient) grpc.UnaryServerInterceptor {
	interceptor := authx.VerifyTokenInterceptor(authx.NewUnverifiedKeyProvider(), authClient)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		switch info.FullMethod {
		case v1.AccountService_LoginUsingPassword_FullMethodName,
			v1.AccountService_CreateAccount_FullMethodName:
			return handler(ctx, req)
		default:
			return interceptor(ctx, req, info, handler)
		}
	}
}

func CheckPermissionInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		//permission, ok := ctx.Value("project_permission").(tinysv1.MemberPermission)
		//if !ok {
		//	return nil, errors.New("missing project permission, check if the user is a member of the project")
		//}
		//
		//switch info.FullMethod {
		//	if permission >= tinysv1.MemberPermission_MEMBER_READ {
		//		return handler(ctx, req)
		//	}
		//
		//default:
		//	return nil, errors.New("unknown method")
		//}

		return nil, errors.New("permission denied")
	}
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
