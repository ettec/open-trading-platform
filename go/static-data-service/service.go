package main

import (
	"context"
	"database/sql"
	"fmt"
	api "github.com/ettec/open-trading-platform/go/static-data-service/api/staticdataservice"
	"github.com/ettec/otp-model"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"strconv"
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
		log.Fatal("could not establish a connection with the database: ", err)
	}

	return s, nil
}

const listingsSelect = `SELECT listings.id, listings.market_symbol, markets.id, markets.name, markets.mic, markets.country_code,
		instruments.id, instruments.name, instruments.display_symbol, instruments.enabled FROM referencedata.listings inner join referencedata.instruments 
		on listings.instrument_id = instruments.id inner join referencedata.markets 
		on listings.market_id = markets.id `

func (s *service) GetListingsWithSameInstrument(c context.Context, id *api.ListingId) (*api.Listings, error) {

	listing, err := s.GetListing(c, id)
	if err != nil {
		return nil, err
	}

	lq := listingsSelect + " where instruments.id = " + strconv.Itoa(int(listing.Instrument.Id))

	r, err := s.db.Query(lq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch listings from database:%w", err)
	}

	result, err := hydrateListings(r)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (s *service) GetListingMatching(c context.Context, m *api.ExactMatchParameters) (*model.Listing, error) {

	lq := listingsSelect + " where display_symbol = '" + m.Symbol + "' and markets.mic = '" + m.Mic + "'"

	r, err := s.db.Query(lq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch listings from database:%w", err)
	}

	result, err := hydrateListings(r)
	if err != nil {
		return nil, err
	}

	if len(result.Listings) == 0 {
		return nil, status.Error(codes.NotFound, "no listing found for symbol "+m.Symbol+" and market "+m.Mic)
	}

	if len(result.Listings) > 1 {
		return nil, status.Error(codes.NotFound, "more than one listing found for symbol "+m.Symbol+" and market "+m.Mic)
	}

	return result.Listings[0], nil
}

func (s *service) GetListingsMatching(c context.Context, m *api.MatchParameters) (*api.Listings, error) {

	lq := listingsSelect + " where display_symbol like '" + m.SymbolMatch + "%'"

	r, err := s.db.Query(lq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch listings from database:%w", err)
	}

	result, err := hydrateListings(r)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) GetListing(c context.Context, id *api.ListingId) (*model.Listing, error) {
	lq := fmt.Sprintf("%v where listings.id = %v", listingsSelect, id.ListingId)

	r, err := s.db.Query(lq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch listings from database:%w", err)
	}

	result, err := hydrateListings(r)
	if err != nil {
		return nil, err
	}

	if len(result.Listings) != 1 {
		return nil, status.Error(codes.NotFound, "no listing found for listing id "+strconv.Itoa(int(id.ListingId)))
	}

	return result.Listings[0], nil

}

func (s *service) GetListings(c context.Context, ids *api.ListingIds) (*api.Listings, error) {

	lq := listingsSelect + " where listings.id in (" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids.ListingIds)), ","), "[]") + ")"

	r, err := s.db.Query(lq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch listings from database:%w", err)
	}

	result, err := hydrateListings(r)
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

func hydrateListings(r *sql.Rows) (*api.Listings, error) {
	result := api.Listings{
		Listings: []*model.Listing{},
	}

	for r.Next() {
		l := &model.Listing{}

		l.Instrument = &model.Instrument{}
		l.Market = &model.Market{}

		err := r.Scan(&l.Id, &l.MarketSymbol, &l.Market.Id, &l.Market.Name, &l.Market.Mic, &l.Market.CountryCode,
			&l.Instrument.Id, &l.Instrument.Name, &l.Instrument.DisplaySymbol, &l.Instrument.Enabled)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal database row into listing %w", err)
		}

		l.SizeIncrement = &model.Decimal64{
			Mantissa: 1,
			Exponent: 0,
		}

		te := model.TickSizeEntry{
			LowerPriceBound: &model.Decimal64{
				Mantissa: 0,
				Exponent: 0,
			},
			UpperPriceBound: &model.Decimal64{
				Mantissa: 1,
				Exponent: 10,
			},
			TickSize: &model.Decimal64{
				Mantissa: 1,
				Exponent: -2,
			},
		}

		entries := make([]*model.TickSizeEntry, 0)
		entries = append(entries, &te)

		tt := model.TickSizeTable{
			Entries: entries,
		}

		l.TickSize = &tt

		result.Listings = append(result.Listings, l)

	}
	return &result, nil
}

const (
	Port                     = "PORT"
	DatabaseConnectionString = "DB_CONN_STRING"
	DatabaseDriverName       = "DB_DRIVER_NAME"
)

func main() {

	dbString := getBootstrapEnvVar(DatabaseConnectionString)
	dbDriverName := getBootstrapEnvVar(DatabaseDriverName)
	port := getOptionalBootstrapEnvVar(Port, "50551")

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
	api.RegisterStaticDataServiceServer(s, service)

	reflection.Register(s)

	fmt.Println("Starting static data service on port:" + port)

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

func getOptionalBootstrapEnvVar(key string, def string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = def
	}

	log.Printf("%v set to %v", key, value)

	return value
}
