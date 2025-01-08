package server

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

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
