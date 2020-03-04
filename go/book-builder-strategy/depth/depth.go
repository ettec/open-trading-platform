package depth


type Depths []Depth

type Depth struct {
	Symbol        string  `json:"symbol"`
	MarketPercent float64 `json:"marketPercent"`
	Volume        int     `json:"volume"`
	LastSalePrice float64 `json:"lastSalePrice"`
	LastSaleSize  int     `json:"lastSaleSize"`
	LastSaleTime  int64   `json:"lastSaleTime"`
	LastUpdated   int64   `json:"lastUpdated"`
	Bids          []struct {
		Price     float64 `json:"price"`
		Size      int     `json:"size"`
		Timestamp int64   `json:"timestamp"`
	} `json:"bids"`
	Asks []struct {
		Price     float64 `json:"price"`
		Size      int     `json:"size"`
		Timestamp int64   `json:"timestamp"`
	} `json:"asks"`
	SystemEvent struct {
		SystemEvent string `json:"systemEvent"`
		Timestamp   int64  `json:"timestamp"`
	} `json:"systemEvent"`
	TradingStatus struct {
		Status    string `json:"status"`
		Reason    string `json:"reason"`
		Timestamp int64  `json:"timestamp"`
	} `json:"tradingStatus"`
	OpHaltStatus struct {
		IsHalted  bool  `json:"isHalted"`
		Timestamp int64 `json:"timestamp"`
	} `json:"opHaltStatus"`
	SsrStatus struct {
		IsSSR     bool   `json:"isSSR"`
		Detail    string `json:"detail"`
		Timestamp int64  `json:"timestamp"`
	} `json:"ssrStatus"`
	SecurityEvent struct {
		SecurityEvent string `json:"securityEvent"`
		Timestamp     int64  `json:"timestamp"`
	} `json:"securityEvent"`
	Trades []struct {
		Price                 float64 `json:"price"`
		Size                  int     `json:"size"`
		TradeID               int     `json:"tradeId"`
		IsISO                 bool    `json:"isISO"`
		IsOddLot              bool    `json:"isOddLot"`
		IsOutsideRegularHours bool    `json:"isOutsideRegularHours"`
		IsSinglePriceCross    bool    `json:"isSinglePriceCross"`
		IsTradeThroughExempt  bool    `json:"isTradeThroughExempt"`
		Timestamp             int64   `json:"timestamp"`
	} `json:"trades"`
	TradeBreaks []interface{} `json:"tradeBreaks"`
}

