package server

import (
	"context"
	"errors"
	"fmt"
	gatewayfile "github.com/black-06/grpc-gateway-file"
	"github.com/emrgen/authbase"
	authv1 "github.com/emrgen/authbase/apis/v1"
	authx "github.com/emrgen/authbase/x"
	docv1 "github.com/emrgen/document/apis/v1"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/config"
	"github.com/emrgen/unpost/internal/model"
	"github.com/emrgen/unpost/internal/service"
	"github.com/emrgen/unpost/internal/store"
	"github.com/gobuffalo/packr"
	"github.com/google/uuid"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcvalidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

// Start starts the grpc and http servers
func Start(grpcPort, httpPort string) error {
	var err error

	grpcPort = ":" + grpcPort
	httpPort = ":" + httpPort

	cnf := config.LoadConfig()
	rdb := config.GetDb(cnf)

	gl, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return err
	}

	rl, err := net.Listen("tcp", httpPort)
	if err != nil {
		return err
	}

	authConfig, err := authx.ConfigFromEnv()
	if err != nil {
		return err
	}

	//tinyConn, err := grpc.NewClient(":4010", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//defer tinyConn.Close()
	//// tinyClient provides the membership service
	//tinyClient := tinysv1.NewMembershipServiceClient(tinyConn)

	// tinyClient provides the membership service
	authClient, err := authbase.NewClient("4000")

	docConn, err := grpc.NewClient(":4020", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer docConn.Close()
	// tinyClient provides the membership service
	docClient := docv1.NewDocumentServiceClient(docConn)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcvalidator.UnaryServerInterceptor(),
			// verify the token
			authx.VerifyTokenInterceptor(authx.NewUnverifiedKeyProvider(), authClient),
			UnaryGrpcRequestTimeInterceptor(),
		)),
	)

	// connect the rest gateway to the grpc server
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
		gatewayfile.WithHTTPBodyMarshaler(),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(UnaryRequestTimeInterceptor()),
	}
	endpoint := "localhost" + grpcPort

	unpostStore := store.NewGormStore(rdb)
	err = unpostStore.Migrate()
	if err != nil {
		return err
	}

	res, err := authClient.GetTokenFromAccessKey(context.TODO(), &authv1.GetTokenFromAccessKeyRequest{
		AccessKey: authConfig.AccessKey,
	})
	if err != nil {
		return err
	}
	claims, err := authx.GetTokenClaims(res.AccessToken)
	if err != nil {
		return err
	}

	// create master space and the owner user
	// create owner user
	spaceID := uuid.New().String()
	err = unpostStore.CreateUser(context.TODO(), &model.User{
		ID:      claims.AccountID,
		SpaceID: spaceID,
	})
	if err != nil {
		return err
	}

	space := &model.Space{
		ID:                spaceID,
		OwnerID:           claims.AccountID,
		AuthbaseProjectID: claims.ProjectID,
		Name:              "unpost",
		Master:            true,
	}

	// create master space
	err = unpostStore.CreateSpace(context.TODO(), space)
	if err != nil && !strings.Contains(err.Error(), "UNIQUE constraint failed: spaces.name") {
		return err
	}

	if err == nil {
		// create master space member
		err = unpostStore.AddSpaceMember(context.TODO(), &model.SpaceMember{
			SpaceID: space.ID,
			UserID:  claims.AccountID,
			Role:    model.SpaceRoleOwner,
		})
		if err != nil && !strings.Contains(err.Error(), "UNIQUE constraint failed: space_members.space_id") {
			return err
		}
	}

	// Register the grpc server
	v1.RegisterTagServiceServer(grpcServer, service.NewTagService(unpostStore))
	v1.RegisterPostServiceServer(grpcServer, service.NewPostService(authConfig, unpostStore, docClient, authClient))
	v1.RegisterCollectionServiceServer(grpcServer, service.NewCollectionService(unpostStore))
	v1.RegisterCourseServiceServer(grpcServer, service.NewCourseService(authConfig, unpostStore, docClient))
	v1.RegisterPageServiceServer(grpcServer, service.NewPageService(authConfig, unpostStore, docClient))
	v1.RegisterSpaceServiceServer(grpcServer, service.NewSpaceService(authConfig, unpostStore, authClient))
	v1.RegisterSpaceMemberServiceServer(grpcServer, service.NewSpaceMemberService(unpostStore, authClient))

	// Register the rest gateway
	if err = v1.RegisterPostServiceHandlerFromEndpoint(context.TODO(), mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterCollectionServiceHandlerFromEndpoint(context.TODO(), mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterTagServiceHandlerFromEndpoint(context.TODO(), mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterCourseServiceHandlerFromEndpoint(context.TODO(), mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterPageServiceHandlerFromEndpoint(context.TODO(), mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterSpaceServiceHandlerFromEndpoint(context.TODO(), mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterSpaceMemberServiceHandlerFromEndpoint(context.TODO(), mux, endpoint, opts); err != nil {
		return err
	}

	apiMux := http.NewServeMux()
	openapiDocs := packr.NewBox("../../docs/v1")
	docsPath := "/v1/docs/"
	apiMux.Handle(docsPath, http.StripPrefix(docsPath, http.FileServer(openapiDocs)))
	apiMux.Handle("/", mux)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // All origins are allowed
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
	})

	restServer := &http.Server{
		Addr:    httpPort,
		Handler: c.Handler(apiMux),
	}

	// make sure to wait for the servers to stop before exiting
	var wg sync.WaitGroup

	wg.Add(1)
	// Start the grpc server
	go func() {
		defer wg.Done()
		logrus.Info("starting rest gateway on: ", httpPort)
		logrus.Info("click on the following link to view the API documentation: http://localhost", httpPort, "/v1/docs/")
		if err := restServer.Serve(rl); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logrus.Errorf("error starting rest gateway: %v", err)
			}
		}
		logrus.Infof("rest gateway stopped")
	}()

	// Start the rest gateway
	wg.Add(1)
	go func() {
		defer wg.Done()
		logrus.Info("starting grpc server on: ", grpcPort)
		if err := grpcServer.Serve(gl); err != nil {
			logrus.Infof("grpc failed to start: %v", err)
		}
		logrus.Infof("grpc server stopped")
	}()

	time.Sleep(1 * time.Second)
	logrus.Infof("Press Ctrl+C to stop the server")

	// listen for interrupt signal to gracefully shut down the server
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, unix.SIGTERM, unix.SIGINT, unix.SIGTSTP)
	<-sigs
	// clean Ctrl+C output
	fmt.Println()

	grpcServer.Stop()
	err = restServer.Shutdown(context.Background())
	if err != nil {
		logrus.Errorf("error stopping rest gateway: %v", err)
	}

	wg.Wait()

	return nil
}

func createMasterSpace() {}
