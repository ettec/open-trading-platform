package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ettec/open-trading-platform/go/static-data-service/internal/model"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"strings"
)

type service struct {
	db *sql.DB
}

func NewService(driverName, dbConnString string) (*service, error) {

	s := &service{}

	db, err := sql.Open(driverName, dbConnString)
	s.db = db
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("could not establish a connection with the database: %w", err)
	}

	return s, nil
}

const listingsSelect = `SELECT listings.id, listings.market_symbol, markets.id, markets.name, markets.mic, markets.country_code,
		instruments.id, instruments.name, instruments.display_symbol, instruments.enabled FROM listings inner join instruments 
		on listings.instrument_id = instruments.id inner join markets 
		on listings.market_id = markets.id `

func (s *service) GetListingsMatching(c context.Context, m *model.MatchParameters) (*model.Listings, error) {

	lq := listingsSelect +  " where display_symbol like '" + m.SymbolMatch +"%'"

	r, err := s.db.Query(lq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch listings from database:%w", err)
	}

	result, err := hydrateListings(r, err)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) GetListing(c context.Context, id *model.ListingId) (*model.Listing, error) {
	lq := fmt.Sprintf( "%v where listings.id = %v", listingsSelect, id.ListingId)

	r, err := s.db.Query(lq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch listings from database:%w", err)
	}

	result, err := hydrateListings(r, err)
	if err != nil {
		return nil, err
	}

	if len(result.Listings) != 1 {
		return nil, fmt.Errorf("expected 1 listing, found %v", len(result.Listings))
	}

	return result.Listings[0], nil

}

func (s *service) GetListings(c context.Context, ids *model.ListingIds) (*model.Listings, error) {

	lq := listingsSelect +  " where listings.id in (" +  strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids.ListingIds)), ",") , "[]")+")"

	r, err := s.db.Query(lq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch listings from database:%w", err)
	}

	result, err := hydrateListings(r, err)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) Close() {
	if s.db != nil {
		err := s.db.Close()
		if err != nil {
			log.Printf("erron closing database connection %v", err)
		}
	}
}

func hydrateListings(r *sql.Rows, err error) (*model.Listings, error) {
	result := model.Listings{
		Listings: []*model.Listing{},
	}

	for r.Next() {
		l := &model.Listing{}

		l.Instrument = &model.Instrument{}
		l.Market = &model.Market{}

		err = r.Scan(&l.Id, &l.MarketSymbol, &l.Market.Id, &l.Market.Name, &l.Market.Mic, &l.Market.CountryCode,
			&l.Instrument.Id, &l.Instrument.Name, &l.Instrument.DisplayName, &l.Instrument.Enabled)
		result.Listings = append(result.Listings, l)
		if err != nil {
			return  nil, fmt.Errorf("failed to marshal database row into listing %w", err)
		}
	}
	return &result, nil
}



const (
	DatabaseConnectionString = "DB_CONN_STRING"
	DatabaseDriverName       = "DB_DRIVER_NAME"
)

func main() {

	dbString := getBootstrapEnvVar(DatabaseConnectionString)
	dbDriverName := getBootstrapEnvVar(DatabaseDriverName)

	port := "50551"
	fmt.Println("Starting static data service on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	service, err := NewService(dbDriverName, dbString)
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}
	defer service.Close()

	s := grpc.NewServer()
	model.RegisterStaticDataServiceServer(s, service)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}

}

// To be used to get environment variables at startup prior to using any resources as it
// exits the process if the env var is not found
func getBootstrapEnvVar(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("missing required env var %v", key)
	}

	log.Printf("%v set to %v", key, value)

	return value
}
