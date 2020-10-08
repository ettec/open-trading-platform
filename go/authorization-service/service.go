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
	rpc "github.com/gogo/googleapis/google/rpc"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	status "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strings"
)

func (a AuthService) Login(ctx context.Context, params *loginservice.LoginParams) (*loginservice.Token, error) {

	log.Printf("logging in")

	if user, ok := a.users[params.User]; ok {
		return &loginservice.Token{
			Token: user.token,
			Desk:  user.desk,
		}, nil
	}

	return nil, errors.New("user not found")
}

type AuthService struct {
	users map[string]user
}

// inject a header that can be used for future rate limiting
func (a *AuthService) Check(_ context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {

	path, ok := req.Attributes.Request.Http.Headers[":path"]

	if ok && strings.HasPrefix(path, "/loginservice.LoginService") {
		log.Printf("permitted login for path:%v", path)
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

type user struct {
	id              string
	desk            string
	permissionFlags string
	token           string
}

func main() {

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime|log.Lshortfile)

	dbString := bootstrap.GetEnvVar(DatabaseConnectionString)
	dbDriverName := bootstrap.GetEnvVar(DatabaseDriverName)

	db, err := sql.Open(dbDriverName, dbString)
	if err != nil {
		log.Panicf("failed to open database connection: %v", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Printf("error when closing database connection: %v", err)
		}
	} ()


	err = db.Ping()
	if err != nil {
		log.Panic("could not establish a connection with the database: ", err)
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

	log.Printf("loaded %v users", len(users))

	authServer := &AuthService{users: users}

	go func() {
		loginPort := "50551"
		lis, err := net.Listen("tcp", ":"+loginPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Printf("listening on %s", lis.Addr())

		grpcServer := grpc.NewServer()
		loginservice.RegisterLoginServiceServer(grpcServer, authServer)

		log.Print("starting login server of port:", loginPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	authPort := "4000"
	lis, err := net.Listen("tcp", ":"+authPort)
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}
	log.Printf("listening on %s", lis.Addr())

	grpcServer := grpc.NewServer()

	auth.RegisterAuthorizationServer(grpcServer, authServer)

	log.Print("starting authorization server of port:", authPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Panicf("Failed to start server: %v", err)
	}

}
