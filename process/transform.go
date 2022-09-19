package process

import (
	"log"
	"math"
	"time"

	"github.com/ArchitBhonsle/go-fpw/fetch"
)

type Data struct {
	Underlying      string
	UnderlyingValue float64
	ExpiryDate      time.Time
	Timestamp       time.Time
	Records         []Record
}

type Record struct {
	StrikePrice float64
	PE          *Option
	CE          *Option
}

type Option struct {
	AskPrice                        float64
	AskQuantity                     int
	BidQuantity                     int
	BidPrice                        float64
	Change                          float64
	ChangeInOpenInterest            float64
	ImpliedVolatility               float64
	LastTradedPrice                 float64
	OpenInterest                    float64
	PercentangeChangeInPrice        float64
	PercentangeChangeInOpenInterest float64
	TotalBuyQuantity                int
	TotalSellQuantity               int
	TotalTradedVolume               int
}

func Transform(fetched fetch.Fetched) (Data, error) {
	log.Println(fetched.Records.Data[0].PE.Underlying, "processing")

	res := Data{}

	timestamp, err := time.Parse("02-Jan-2006 15:04:05", fetched.Records.Timestamp)
	if err != nil {
		return res, err
	}
	res.Timestamp = timestamp
	res.UnderlyingValue = fetched.Records.UnderlyingValue

	expiryDates := make([]time.Time, 0, len(fetched.Records.ExpiryDates))
	for _, ed := range fetched.Records.ExpiryDates {
		edParsed, err := time.Parse("02-Jan-2006", ed)
		if err != nil {
			return res, err
		}
		expiryDates = append(expiryDates, edParsed)
	}

	var expiryDate time.Time // for now it's the first expiry date after the current timestamp
	for _, ed := range expiryDates {
		if ed.After(timestamp) {
			expiryDate = ed
			break
		}
	}
	res.ExpiryDate = expiryDate
	targetExpiryDate := expiryDate.Format("02-Jan-2006")

	strikePrices := make(map[float64]bool)
	baseIndex, difference := -1, math.MaxFloat64
	for i, sp := range fetched.Records.StrikePrices {
		diff := math.Abs(float64(sp) - res.UnderlyingValue)
		if diff < difference {
			baseIndex = i
			difference = diff
		}
	}
	for i, sp := range fetched.Records.StrikePrices {
		diff := i - baseIndex
		if -5 < diff && diff <= 5 {
			strikePrices[float64(sp)] = true
		}
	}

	records := make([]Record, 0)
	for _, optionData := range fetched.Records.Data {
		if optionData.ExpiryDate != targetExpiryDate ||
			!strikePrices[float64(optionData.StrikePrice)] {
			continue
		}

		if res.Underlying == "" {
			if optionData.CE != nil {
				res.Underlying = optionData.CE.Underlying
			}
			if optionData.PE != nil {
				res.Underlying = optionData.PE.Underlying
			}
		}

		record := Record{
			StrikePrice: float64(optionData.StrikePrice),
			CE:          transformOption(optionData.CE),
			PE:          transformOption(optionData.PE),
		}

		records = append(records, record)
	}
	res.Records = records

	return res, nil
}

func transformOption(o *fetch.Option) *Option {
	if o == nil {
		return nil
	}

	return &Option{
		AskPrice:                        o.AskPrice,
		AskQuantity:                     o.AskQuantity,
		BidQuantity:                     o.BidQuantity,
		BidPrice:                        o.BidPrice,
		Change:                          o.Change,
		ChangeInOpenInterest:            o.ChangeInOpenInterest,
		ImpliedVolatility:               o.ImpliedVolatility,
		LastTradedPrice:                 o.LastTradedPrice,
		OpenInterest:                    o.OpenInterest,
		PercentangeChangeInPrice:        o.PercentangeChangeInPrice,
		PercentangeChangeInOpenInterest: o.PercentangeChangeInOpenInterest,
		TotalBuyQuantity:                o.TotalBuyQuantity,
		TotalSellQuantity:               o.TotalSellQuantity,
		TotalTradedVolume:               o.TotalTradedVolume,
	}
}
