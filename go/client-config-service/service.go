package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ettec/otp-common/bootstrap"
	"github.com/ettec/otp-model"
	api "github.com/ettech/open-trading-platform/go/client-config-service/api/clientconfigservice"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

type service struct {
	db *sql.DB
}

func (s *service) StoreClientConfig(_ context.Context, params *api.StoreConfigParams) (*model.Empty, error) {

	lq := fmt.Sprintf("select count(*) from clientconfig.reactclientconfig where userid = '%v'", params.UserId)

	r, err := s.db.Query(lq)

	if err != nil {
		return nil, fmt.Errorf("failed to check for existing config in database:%w", err)
	}

	var count int
	r.Next()
	err = r.Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing config in database:%w", err)
	}

	var q string
	if count == 0 {
		q = fmt.Sprintf("insert into clientconfig.reactclientconfig (userid, config) values ('%v', '%v')", params.UserId, params.Config)
	} else {
		q = fmt.Sprintf("update  clientconfig.reactclientconfig set config = '%v' where  userid = '%v'", params.Config, params.UserId)
	}

	_, err = s.db.Exec(q)

	if err != nil {
		return nil, fmt.Errorf("failed to update config in database:%w", err)
	}

	return &model.Empty{}, nil

}

const configSelect = "select config from clientconfig.reactclientconfig where userid = '%v'"

func (s *service) GetClientConfig(_ context.Context, parameters *api.GetConfigParameters) (*api.Config, error) {

	lq := fmt.Sprintf(configSelect, parameters.UserId)

	r, err := s.db.Query(lq)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch config from database:%w", err)
	}

	if r.Next() {
		var config string
		err = r.Scan(&config)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch config from database:%w", err)
		}


		return &api.Config{
			Config: config,
		}, nil
	} else {
		return nil, fmt.Errorf("failed to find configuration for user id:%v", parameters.UserId)
	}

}

func newService(driverName, dbConnString string) (*service, error) {

	s := &service{}

	db, err := sql.Open(driverName, dbConnString)
	s.db = db
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("could not establish a connection with the database: ", err)
	}

	return s, nil
}

func (s *service) Close() {
	if s.db != nil {
		err := s.db.Close()
		if err != nil {
			log.Printf("erron closing database connection %v", err)
		}
	}
}


func main() {

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime|log.Lshortfile)

	dbString := bootstrap.GetEnvVar("DB_CONN_STRING")
	dbDriverName := bootstrap.GetEnvVar("DB_DRIVER_NAME")
	port := bootstrap.GetOptionalEnvVar("PORT", "50551")

	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Panicf("Error while listening : %v", err)
	}

	service, err := newService(dbDriverName, dbString)
	if err != nil {
		log.Panicf("failed to create service: %v", err)
	}
	defer service.Close()

	s := grpc.NewServer()
	api.RegisterClientConfigServiceServer(s, service)

	reflection.Register(s)
	fmt.Println("Started client config service on port:" + port)
	if err := s.Serve(lis); err != nil {
		log.Panicf("Error while serving : %v", err)
	}

}
