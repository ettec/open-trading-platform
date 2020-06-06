package main

import (
	"context"
	"errors"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	envoytype "github.com/envoyproxy/go-control-plane/envoy/type"
	"github.com/ettech/open-trading-platform/go/authorization-service/api/loginservice"
	rpc "github.com/gogo/googleapis/google/rpc"
	status "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"log"
	"net"
)

type LoginService struct{}

func (l LoginService) Login(ctx context.Context, params *loginservice.LoginParams) (*loginservice.Token, error) {
	if params.User == "bert" && params.Password == "poop" {
		return &loginservice.Token{
			Token: "bob",
			Desk:  "Delta1",
		}, nil
	}

	return nil, errors.New("invalid username/password")
}

// empty struct because this isn't a fancy example
type AuthorizationServer struct{}

// inject a header that can be used for future rate limiting
func (a *AuthorizationServer) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {

	path, ok := req.Attributes.Request.Http.Headers[":path"]
	if ok && path == "/clientconfigservice.ClientConfigService/GetClientConfig" {
		log.Print("allowing it:", req)
		return &auth.CheckResponse{
			Status: &status.Status{
				Code: int32(rpc.OK),
			},
			HttpResponse: &auth.CheckResponse_OkResponse{
				OkResponse: &auth.OkHttpResponse{
					Headers: []*envoy_api_v2_core.HeaderValueOption{
						{
							Header: &envoy_api_v2_core.HeaderValue{
								// here is where set trading perms
								Key:   "somestuff",
								Value: "ratatouey",
							},
						},
					},
				},
			},
		}, nil
	}

	log.Print("checking it:", req)

	authHeader, ok := req.Attributes.Request.Http.Headers["auth-token"]

	if ok && authHeader == "bob" {

		// valid tokens have exactly 3 characters. #secure.
		// Normally this is where you'd go check with the system that knows if it's a valid token.
		return &auth.CheckResponse{
			Status: &status.Status{
				Code: int32(rpc.OK),
			},
			HttpResponse: &auth.CheckResponse_OkResponse{
				OkResponse: &auth.OkHttpResponse{
					Headers: []*envoy_api_v2_core.HeaderValueOption{
						{
							Header: &envoy_api_v2_core.HeaderValue{
								// here is where set trading perms
								Key:   "x-ext-auth-ratelimit",
								Value: "ratatouey",
							},
						},
					},
				},
			},
		}, nil

	}

	return &auth.CheckResponse{
		Status: &status.Status{
			Code: int32(rpc.UNAUTHENTICATED),
		},
		HttpResponse: &auth.CheckResponse_DeniedResponse{
			DeniedResponse: &auth.DeniedHttpResponse{
				Status: &envoytype.HttpStatus{
					Code: envoytype.StatusCode_Unauthorized,
				},
				Body: "No authorisation token found",
			},
		},
	}, nil
}

func main() {
	// create a TCP listener on port 4000
	authPort := "4000"
	lis, err := net.Listen("tcp", ":"+authPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listening on %s", lis.Addr())

	grpcServer := grpc.NewServer()
	authServer := &AuthorizationServer{}
	auth.RegisterAuthorizationServer(grpcServer, authServer)

	log.Print("starting authorization server of port:", authPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
