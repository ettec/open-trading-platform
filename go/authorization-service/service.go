package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	envoytype "github.com/envoyproxy/go-control-plane/envoy/type"
	"github.com/ettec/otp-common/bootstrap"
	"github.com/ettech/open-trading-platform/go/authorization-service/api/loginservice"
	"github.com/gogo/googleapis/google/rpc"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type user struct {
	id              string
	desk            string
	permissionFlags string
	token           string
}

type authService struct {
	users map[string]user
}

func (a *authService) Login(_ context.Context, params *loginservice.LoginParams) (*loginservice.Token, error) {

	log.Printf("logging in")

	if user, ok := a.users[params.User]; ok {
		return &loginservice.Token{
			Token: user.token,
			Desk:  user.desk,
		}, nil
	}

	return nil, errors.New("user not found")
}

func (a *authService) Check(_ context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {

	path, ok := req.Attributes.Request.Http.Headers[":path"]

	if ok && strings.HasPrefix(path, "/loginservice.LoginService") {
		slog.Info("permitted login for path", "path", path)
		return newOkResponse(), nil
	}

	authHeader, ok := req.Attributes.Request.Http.Headers["auth-token"]
	if !ok {
		return newPermissionDeniedResponse("auth-token header is required"), nil
	}

	username, ok := req.Attributes.Request.Http.Headers["user-name"]
	if !ok {
		return newUnauthenticatedResponse("No user-name found on request"), nil
	}

	user, ok := a.users[username]
	if !ok {
		return newUnauthenticatedResponse("user not found"), nil
	}

	if user.token != authHeader {
		return newUnauthenticatedResponse("invalid token"), nil
	}

	// Authorisation
	if ok && strings.HasPrefix(path, "/executionvenue.ExecutionVenue") {
		if strings.Contains(user.permissionFlags, "T") {
			return newOkResponse(), nil
		} else {
			return newPermissionDeniedResponse("trading permissions required"), nil
		}
	}

	return newOkResponse(), nil
}

func newOkResponse() *auth.CheckResponse {
	return &auth.CheckResponse{
		Status: &status.Status{
			Code: int32(rpc.OK),
		},
		HttpResponse: &auth.CheckResponse_OkResponse{
			OkResponse: &auth.OkHttpResponse{
				Headers: []*envoy_api_v2_core.HeaderValueOption{
					{
						Header: &envoy_api_v2_core.HeaderValue{
							Key:   "authorised",
							Value: "true",
						},
					},
				},
			},
		},
	}
}

func newPermissionDeniedResponse(message string) *auth.CheckResponse {
	return &auth.CheckResponse{
		Status: &status.Status{
			Code: int32(rpc.PERMISSION_DENIED),
		},
		HttpResponse: &auth.CheckResponse_DeniedResponse{
			DeniedResponse: &auth.DeniedHttpResponse{
				Status: &envoytype.HttpStatus{
					Code: envoytype.StatusCode_Unauthorized,
				},
				Body: message,
			},
		},
	}
}

func newUnauthenticatedResponse(message string) *auth.CheckResponse {
	return &auth.CheckResponse{
		Status: &status.Status{
			Code: int32(rpc.UNAUTHENTICATED),
		},
		HttpResponse: &auth.CheckResponse_DeniedResponse{
			DeniedResponse: &auth.DeniedHttpResponse{
				Status: &envoytype.HttpStatus{
					Code: envoytype.StatusCode_Unauthorized,
				},
				Body: message,
			},
		},
	}
}

const (
	DatabaseConnectionString = "DB_CONN_STRING"
	DatabaseDriverName       = "DB_DRIVER_NAME"
)

func main() {

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))

	dbString := bootstrap.GetEnvVar(DatabaseConnectionString)
	dbDriverName := bootstrap.GetEnvVar(DatabaseDriverName)

	db, err := sql.Open(dbDriverName, dbString)
	if err != nil {
		log.Panicf("failed to open database connection: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("error when closing database connection", "error", err)
		}
	}()

	err = db.Ping()
	if err != nil {
		log.Panicf("could not establish a connection with the database: %v", err)
	}

	r, err := db.Query("SELECT id, desk, permissionflags FROM users.users")
	if err != nil {
		log.Panicf("failed to get users from database")
	}

	users := map[string]user{}
	for r.Next() {
		u := user{}
		err := r.Scan(&u.id, &u.desk, &u.permissionFlags)
		if err != nil {
			log.Panicf("failed to scan user row: %v", err)
		}
		u.token = uuid.New().String()
		users[u.id] = u
	}

	slog.Info("loaded users", "userCount", len(users))

	authServer := &authService{users: users}

	go func() {

		loginPort := "50551"
		lis, err := net.Listen("tcp", ":"+loginPort)
		if err != nil {
			log.Panicf("failed to listen: %v", err)
		}
		slog.Info("authentication server listening", "listenAddress", lis.Addr())

		authenticationGrpcServer := grpc.NewServer()
		loginservice.RegisterLoginServiceServer(authenticationGrpcServer, authServer)

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh,
			syscall.SIGKILL,
			syscall.SIGTERM,
			syscall.SIGQUIT)
		go func() {
			<-sigCh
			authenticationGrpcServer.GracefulStop()
		}()

		slog.Info("starting authentication server", "port", loginPort)
		if err := authenticationGrpcServer.Serve(lis); err != nil {
			log.Panicf("Failed to start authentication server: %v", err)
		}
	}()

	authPort := "4000"
	lis, err := net.Listen("tcp", ":"+authPort)
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}
	slog.Info("authorisation server listening", "listenAddress", lis.Addr())

	grpcServer := grpc.NewServer()

	auth.RegisterAuthorizationServer(grpcServer, authServer)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigCh
		grpcServer.GracefulStop()
	}()

	slog.Info("starting authorization server", "port", authPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Panicf("Failed to start authorization server: %v", err)
	}

}
