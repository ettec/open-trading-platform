package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type instrument struct {
	instId int32
	symbol string
}

func main() {

	db, err := sql.Open("postgres", "host=192.168.1.200 dbname=cnoms sslmode=disable user=cnomsk8s password=password")

	if err != nil {
		log.Panic("Error: The data source arguments are not valid")
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Panic("Error: Could not establish a connection with the database")
	}

	db.Exec(`set search_path="referencedata"`)

	marketIds := []int32{}

	sql := "select id from markets where mic = 'IEXG' or mic = 'XNAS' or mic = 'XOSR'"
	rows, err := db.Query(sql)
	if err != nil {
		log.Panic("failed to query markets:", err)
	}
	for rows.Next() {
		var marketId int32
		err = rows.Scan(&marketId)
		if err != nil {
			log.Panic("failed to read market record:", err)
		}

		marketIds = append(marketIds, marketId)
	}

	instSql := "select id, display_symbol from instruments"
	rows, err = db.Query(instSql)
	if err != nil {
		log.Panic("failed to query instruments:", err)
	}

	instruments := []instrument{}
	for rows.Next() {
		var instId int32
		var symbol string
		err = rows.Scan(&instId, &symbol)
		if err != nil {
			log.Panicf("failed to load instruments:%v", err)
		}

		instruments = append(instruments, instrument{
			instId: instId,
			symbol: symbol,
		})
	}

	log.Printf("creating listings for market ids: %v", marketIds)

	sql = "INSERT INTO listings (instrument_id, market_id, market_symbol) VALUES ($1, $2, $3) RETURNING id"
	for _, marketId := range marketIds {
		for _, instrument := range instruments {

			_, err = db.Exec(sql, instrument.instId, marketId, instrument.symbol)

			if err != nil {
				log.Printf("Error: Failed to insert row error instid:%v marketid:%v  err: %v row sql:%v", instrument.instId, marketId, err, sql)
			}

		}
	}

}
