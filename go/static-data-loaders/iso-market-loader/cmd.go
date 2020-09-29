package main

import (
	"database/sql"
	"encoding/csv"
	_ "github.com/lib/pq"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

type market struct {
	countryCode string
	name        string
	mic         string
}

func main() {

	markets := make([]market, 0)

	data, err := ioutil.ReadFile("./resources/IISO10383_MIC.csv")
	if err != nil {
		log.Fatalf("failed to read markets file:%v", err)
	}

	csvString := string(data)
	r := csv.NewReader(strings.NewReader(csvString))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else {
			if err != nil {
				log.Panicf("error whilst reading file:%v", err)
			}
		}


		market := market{
			countryCode: record[1],
			name:        record[5],
			mic:         record[2],
		}

		market.name = strings.ReplaceAll(market.name, "'", "")

		markets = append(markets, market)

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

	for _, market := range markets {

		sqlstmt := "INSERT INTO markets (mic, name, country_code ) VALUES ('" + market.mic + "','" + market.name + "','" + market.countryCode + "')"

		_, err := db.Exec(sqlstmt)
		if err != nil {
			log.Printf("Error: Failed to insert row error:%v  row sql:%v", err, sqlstmt)
		}

	}

}
