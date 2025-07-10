package server

import (
	"context"
	"errors"
	"fmt"
	gatewayfile "github.com/black-06/grpc-gateway-file"
	authx "github.com/emrgen/authbase/x"
	v1 "github.com/emrgen/unpost/apis/v1"
	"github.com/emrgen/unpost/internal/config"
	"github.com/emrgen/unpost/internal/model"
	"github.com/emrgen/unpost/internal/service"
	"github.com/emrgen/unpost/internal/store"
	"github.com/gobuffalo/packr"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcvalidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/supabase-community/auth-go"
	"github.com/supabase-community/auth-go/types"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

// Start starts the grpc and http servers
func Start(grpcPort, httpPort string) error {
	var err error

	grpcPort = ":" + grpcPort
	httpPort = ":" + httpPort

	cfg := config.LoadConfig()
	rdb := config.GetDb(cfg)

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

	// authClient provides the auth service
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcvalidator.UnaryServerInterceptor(),
			VerifyTokenInterceptor(cfg.SupabaseConfig.JwtSecret),
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

	authClient := auth.New(cfg.SupabaseConfig.ProjectRef, cfg.SupabaseConfig.ApiKey)

	// update the admin role on app start up if not set
	client := authClient.WithToken(cfg.SupabaseConfig.ApiKey)
	user, err := client.AdminGetUser(types.AdminGetUserRequest{UserID: cfg.AdminUserID})
	if err != nil {
		logrus.Error(err)
	}
	if _, ok := user.AppMetadata["role"]; !ok {
		_, err := client.AdminUpdateUser(types.AdminUpdateUserRequest{
			UserID: cfg.AdminUserID,
			AppMetadata: map[string]interface{}{
				"role": model.UserRoleOwner,
			},
		})
		if err != nil {
			return err
		}
	}

	users, err := client.AdminListUsers(types.AdminListUsersRequest{})
	if err != nil {
		return err
	}

	logrus.Infof("Admin users: %v", users)

	unpostStore := store.NewGormStore(rdb)
	err = unpostStore.Migrate()
	if err != nil {
		return err
	}

	// Register the grpc server
	v1.RegisterAccountServiceServer(grpcServer, service.NewAccountService(unpostStore, authClient))
	v1.RegisterPostServiceServer(grpcServer, service.NewPostService(authConfig, unpostStore))
	//v1.RegisterTagServiceServer(grpcServer, service.NewTagService(unpostStore))
	//v1.RegisterCourseServiceServer(grpcServer, service.NewCourseService(authConfig, unpostStore))
	//v1.RegisterPageServiceServer(grpcServer, service.NewPageService(authConfig, unpostStore))

	// Register the rest gateway
	if err = v1.RegisterAccountServiceHandlerFromEndpoint(context.TODO(), mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterPostServiceHandlerFromEndpoint(context.TODO(), mux, endpoint, opts); err != nil {
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

	apiMux := http.NewServeMux()
	openapiDocs := packr.NewBox("../../docs/v1")
	docsPath := "/v1/docs/"
	apiMux.Handle(docsPath, http.StripPrefix(docsPath, http.FileServer(openapiDocs)))
	apiMux.Handle("/", mux)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // All origins are allowed
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
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
