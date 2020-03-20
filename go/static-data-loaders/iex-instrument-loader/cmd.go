package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type iexInstrument struct {
	Symbol    string      `json:"symbol"`
	Name      string      `json:"name"`
	Date      string      `json:"date"`
	IsEnabled bool        `json:"isEnabled"`
	Type      string      `json:"type"`
	IexID     interface{} `json:"iexId"`
}

type listing struct {
	marketId     int32
	instrumentId int
	marketSymbol string
}

func main() {

	iexUrl := "https://api.iextrading.com/1.0/ref-data/symbols"

	iexClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, iexUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := iexClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	iexInstruments := make([]iexInstrument, 2000)
	jsonErr := json.Unmarshal(body, &iexInstruments)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	db, err := sql.Open("postgres", "host=192.168.1.200 dbname=cnoms sslmode=disable user=cnomsk8s password=password")

	if err != nil {
		log.Fatal("Error: The data source arguments are not valid")
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error: Could not establish a connection with the database")
	}

	db.Exec(`set search_path="referencedata"`)

	for _, iexInst := range iexInstruments {

		sourceMap := make(map[string]iexInstrument)
		sourceMap["IEX"] = iexInst

		bytes, err := json.Marshal(sourceMap)
		if err != nil {
			log.Fatal(err)
		}

		json := string(bytes)
		json = strings.Replace(json, "'", "''", -1)

		sql := "INSERT INTO instruments (name, display_symbol, enabled, raw_sources) VALUES ($1, $2, $3, $4) RETURNING id"

		lastInsertId := 0
		err = db.QueryRow(sql, iexInst.Name, iexInst.Symbol, iexInst.IsEnabled, json).Scan(&lastInsertId)

	}

}
