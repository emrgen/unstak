package cmd

import (
	v1 "github.com/emrgen/unpost/apis/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getConnection() (*grpc.ClientConn, func() error) {
	conn, err := grpc.NewClient(":8030", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
		return nil, nil
	}

	return conn, conn.Close
}

func outletClient() (v1.OutletServiceClient, func() error) {
	conn, closer := getConnection()
	client := v1.NewOutletServiceClient(conn)

	return client, closer
}

func postClient() (v1.PostServiceClient, func() error) {
	conn, closer := getConnection()
	client := v1.NewPostServiceClient(conn)

	return client, closer
}

func tagClient() (v1.TagServiceClient, func() error) {
	conn, closer := getConnection()

	client := v1.NewTagServiceClient(conn)
	return client, closer
}
