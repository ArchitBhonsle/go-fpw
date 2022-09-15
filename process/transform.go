package process

import (
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

func transform(fetched fetch.Fetched) (Data, error) {
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

	const multipleOf = 50
	const delta = 4
	smallerMultiple := int(res.UnderlyingValue / multipleOf * multipleOf)
	largerMultiple := int(smallerMultiple + multipleOf)
	lowerBound := smallerMultiple - (delta * multipleOf)
	upperBound := largerMultiple + (delta * multipleOf)
	checkStrikePrice := func(sp int) bool {
		if lowerBound <= sp && upperBound >= sp {
			return true
		} else {
			return false
		}
	}

	records := make([]Record, 0)
	for _, optionData := range fetched.Records.Data {
		if optionData.ExpiryDate != targetExpiryDate ||
			!checkStrikePrice(optionData.StrikePrice) {
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
