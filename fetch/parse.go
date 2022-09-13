package fetch

import "encoding/json"

type Fetched struct {
	Records Records `json:"records"`
}

type Records struct {
	Data            []OptionData `json:"data"`
	ExpiryDates     []string     `json:"expiryDates"`
	StrikePrices    []int        `json:"strikePrices"`
	Timestamp       string       `json:"timestamp"`
	UnderlyingValue float64      `json:"underlyingValue"`
}

type OptionData struct {
	StrikePrice int     `json:"strikePrice"`
	ExpiryDate  string  `json:"expiryDate"`
	CE          *Option `json:"CE,omitempty"`
	PE          *Option `json:"PE,omitempty"`
}

type Option struct {
	AskPrice                        float64 `json:"askPrice"`              // best open sell order price
	AskQuantity                     int     `json:"askQty"`                // open sell order quantity
	BidQuantity                     int     `json:"bidQty"`                // open buy order quantity
	BidPrice                        float64 `json:"bidPrice"`              // best open buy order price
	Change                          float64 `json:"change"`                // change in LTP since the previous closing
	ChangeInOpenInterest            float64 `json:"changeInOpenInterest"`  // change in unique contracts since the previous closing
	ExpiryDate                      string  `json:"expiryDate"`            //
	Identifier                      string  `json:"identifier"`            //
	ImpliedVolatility               float64 `json:"impliedVolatility"`     // volatility of the underlying
	LastTradedPrice                 float64 `json:"lastPrice"`             // current price
	OpenInterest                    float64 `json:"openInterest"`          // open unique contracts
	PercentangeChangeInPrice        float64 `json:"pChange"`               //
	PercentangeChangeInOpenInterest float64 `json:"pchangeInOpenInterest"` //
	StrikePrice                     int     `json:"strikePrice"`           //
	TotalBuyQuantity                int     `json:"totalBuyQuantity"`      //
	TotalSellQuantity               int     `json:"totalSellQuantity"`     //
	TotalTradedVolume               int     `json:"totalTradedVolume"`     // traded (not necessarily unique) contracts
	Underlying                      string  `json:"underlying"`            //
	UnderlyingValue                 float64 `json:"underlyingValue"`       //
}

func ParseJson(body []byte) (Fetched, error) {
	data := Fetched{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}
