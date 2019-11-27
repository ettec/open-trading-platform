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

type canonicalInstrument struct {
	Date      string `json:"date"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	IsEnabled bool   `json:"isEnabled"`
}

type otpInstrument struct {
	Raw     map[string]interface{} `json:"raw"`
	Canon   canonicalInstrument    `json:"canon"`
	Symbols map[string]string      `json:"symbols"`
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
	if (err != nil) {
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

	normalisedInstruments := make([]string, 0,2000)
	for _, iex := range iexInstruments {

		inst := otpInstrument{Symbols: map[string]string{"IEX": iex.Symbol}, Raw: map[string]interface{}{"IEX": iex},
			Canon: canonicalInstrument{Name: iex.Name, Date: iex.Date, Type: iex.Type, IsEnabled: iex.IsEnabled}}

		jsonBytes, err := json.Marshal(inst)
		if err != nil {
			log.Fatal(err);
		}

		jsonString := string(jsonBytes);
		normalisedInstruments = append(normalisedInstruments, jsonString)
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

	for _, inst := range normalisedInstruments {
		inst = strings.Replace(inst, "'", "''", -1);

		sql := "INSERT INTO instruments (json) VALUES ('" + inst + "'::jsonb)"

		_, err := db.Exec(sql)
		if err != nil  {
			log.Printf("Error: Failed to insert row error:%v  row sql:%v", err, sql)
		}


	}



	//db.Exec()

}
