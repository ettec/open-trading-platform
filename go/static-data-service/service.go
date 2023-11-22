package main

import (
	"context"
	"database/sql"
	"fmt"
	api "github.com/ettec/open-trading-platform/go/static-data-service/api/staticdataservice"
	"github.com/ettec/otp-common/bootstrap"
	"github.com/ettec/otp-common/model"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
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

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("could not establish a connection with the database: %w", err)
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
		return nil, fmt.Errorf("failed to fetch listing:%w", err)
	}

	lq := listingsSelect + " where instruments.id = " + strconv.Itoa(int(listing.Instrument.Id))

	r, err := s.db.Query(lq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch listings:%w", err)
	}

	result, err := hydrateListings(r)
	if err != nil {
		return nil, fmt.Errorf("failed to hydrate listings:%w", err)
	}

	return result, nil

}

func (s *service) GetListingMatching(_ context.Context, m *api.ExactMatchParameters) (*model.Listing, error) {

	lq := listingsSelect + " where display_symbol = '" + m.Symbol + "' and markets.mic = '" + m.Mic + "'"

	r, err := s.db.Query(lq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch listings:%w", err)
	}

	result, err := hydrateListings(r)
	if err != nil {
		return nil, fmt.Errorf("failed to hydrate listings:%w", err)
	}

	if len(result.Listings) == 0 {
		return nil, status.Error(codes.NotFound, "no listing found for symbol "+m.Symbol+" and market "+m.Mic)
	}

	if len(result.Listings) > 1 {
		return nil, status.Error(codes.NotFound, "more than one listing found for symbol "+m.Symbol+" and market "+m.Mic)
	}

	return result.Listings[0], nil
}

func (s *service) GetListingsMatching(_ context.Context, m *api.MatchParameters) (*api.Listings, error) {

	lq := listingsSelect + " where display_symbol like '" + m.SymbolMatch + "%'"

	r, err := s.db.Query(lq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch listings from database:%w", err)
	}

	result, err := hydrateListings(r)
	if err != nil {
		return nil, fmt.Errorf("failed to hydrate listings:%w", err)
	}

	return result, nil
}

func (s *service) GetListing(_ context.Context, id *api.ListingId) (*model.Listing, error) {
	lq := fmt.Sprintf("%v where listings.id = %v", listingsSelect, id.ListingId)

	r, err := s.db.Query(lq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch listings from database:%w", err)
	}

	result, err := hydrateListings(r)
	if err != nil {
		return nil, fmt.Errorf("failed to hydrate listings:%w", err)
	}

	if len(result.Listings) != 1 {
		return nil, status.Error(codes.NotFound, "no listing found for listing id "+strconv.Itoa(int(id.ListingId)))
	}

	return result.Listings[0], nil

}

func (s *service) GetListings(_ context.Context, ids *api.ListingIds) (*api.Listings, error) {

	lq := listingsSelect + " where listings.id in (" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids.ListingIds)), ","), "[]") + ")"

	r, err := s.db.Query(lq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch listings from database:%w", err)
	}

	result, err := hydrateListings(r)
	if err != nil {
		return nil, fmt.Errorf("failed to hydrate listings:%w", err)
	}

	return result, nil
}

func (s *service) Close() {
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			slog.Error("error when closing database connection", "error", err)
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
			return nil, fmt.Errorf("failed to scan database row into listing %w", err)
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

func main() {

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))

	dbString := bootstrap.GetEnvVar("DB_CONN_STRING")
	dbDriverName := bootstrap.GetEnvVar("DB_DRIVER_NAME")
	port := bootstrap.GetOptionalEnvVar("PORT", "50551")

	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Panicf("Error while listening : %v", err)
	}

	staticDataService, err := NewService(dbDriverName, dbString)
	if err != nil {
		log.Panicf("failed to create service: %v", err)
	}
	defer staticDataService.Close()

	s := grpc.NewServer()
	api.RegisterStaticDataServiceServer(s, staticDataService)

	reflection.Register(s)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigCh
		s.GracefulStop()
	}()

	slog.Info("Starting static data service", "port", port)

	if err := s.Serve(lis); err != nil {
		log.Panicf("Error while serving : %v", err)
	}

}
