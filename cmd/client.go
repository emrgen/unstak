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

func tierClient() (v1.TierServiceClient, func() error) {
	conn, closer := getConnection()
	client := v1.NewTierServiceClient(conn)

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

func collectionClient() (v1.CollectionServiceClient, func() error) {
	conn, closer := getConnection()

	client := v1.NewCollectionServiceClient(conn)
	return client, closer
}

func spaceClient() (v1.SpaceServiceClient, func() error) {
	conn, closer := getConnection()

	client := v1.NewSpaceServiceClient(conn)
	return client, closer
}
