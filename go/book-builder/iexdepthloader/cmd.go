package main

import (
	"encoding/json"
	"github.com/ettec/open-trading-platform/go/book-builder/depth"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
)

type tops []top

type top struct {
	Symbol        string  `json:"symbol"`
	Sector        string  `json:"sector"`
	SecurityType  string  `json:"securityType"`
	BidPrice      float64 `json:"bidPrice"`
	BidSize       int     `json:"bidSize"`
	AskPrice      float64 `json:"askPrice"`
	AskSize       int     `json:"askSize"`
	LastUpdated   int64   `json:"lastUpdated"`
	LastSalePrice float64 `json:"lastSalePrice"`
	LastSaleSize  int     `json:"lastSaleSize"`
	LastSaleTime  int64   `json:"lastSaleTime"`
	Volume        int     `json:"volume"`
	MarketPercent float64 `json:"marketPercent"`
}

func (s tops) Len() int {
	return len(s)
}
func (s tops) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s tops) Less(i, j int) bool {
	return s[i].Volume < s[j].Volume
}

func main() {

	body := readIexJson("https://api.iextrading.com/1.0/tops")

	var t tops = make([]top, 2000)
	jsonErr := json.Unmarshal(body, &t)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	sort.Sort(t)

	var depths depth.Depths = make([]depth.Depth, 0)

	for i := 0; i < 100; i++ {
		idx := len(t) - i - 1
		body := readIexJson("https://api.iextrading.com/1.0/deep?symbols=" + t[idx].Symbol)
		depth := depth.Depth{}
		jsonErr := json.Unmarshal(body, &depth)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		depth.Trades = []struct {
			Price                 float64 `json:"price"`
			Size                  int     `json:"size"`
			TradeID               int     `json:"tradeId"`
			IsISO                 bool    `json:"isISO"`
			IsOddLot              bool    `json:"isOddLot"`
			IsOutsideRegularHours bool    `json:"isOutsideRegularHours"`
			IsSinglePriceCross    bool    `json:"isSinglePriceCross"`
			IsTradeThroughExempt  bool    `json:"isTradeThroughExempt"`
			Timestamp             int64   `json:"timestamp"`
		}{}

		depths = append(depths, depth)
	}

	var symbols []string
	for _, d := range depths {
		symbols = append(symbols, d.Symbol)
	}

	body, jsonErr = json.Marshal(depths)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	ioutil.WriteFile("./resources/depth.json", body, os.ModePerm)

	log.Println("loaded depth for symbols:", symbols)

}

func readIexJson(iexUrl string) []byte {
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
	return body
}
