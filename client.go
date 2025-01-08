package unpost

import (
	v1 "github.com/emrgen/unpost/apis/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
)

type Client interface {
	v1.CollectionServiceClient
	v1.CourseServiceClient
	v1.PageServiceClient
	v1.PostServiceClient
	v1.SpaceServiceClient
	v1.TagServiceClient
	v1.TierServiceClient
	io.Closer
}

type client struct {
	conn *grpc.ClientConn
	v1.CollectionServiceClient
	v1.CourseServiceClient
	v1.PageServiceClient
	v1.PostServiceClient
	v1.SpaceServiceClient
	v1.TagServiceClient
	v1.TierServiceClient
}

func NewClient(port string) (Client, error) {
	conn, err := grpc.NewClient(":8030", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &client{
		conn:                    conn,
		CollectionServiceClient: v1.NewCollectionServiceClient(conn),
		CourseServiceClient:     v1.NewCourseServiceClient(conn),
		PageServiceClient:       v1.NewPageServiceClient(conn),
		PostServiceClient:       v1.NewPostServiceClient(conn),
		SpaceServiceClient:      v1.NewSpaceServiceClient(conn),
		TagServiceClient:        v1.NewTagServiceClient(conn),
		TierServiceClient:       v1.NewTierServiceClient(conn),
	}, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}
