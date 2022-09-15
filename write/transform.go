package write

import (
	"time"

	"github.com/ArchitBhonsle/go-fpw/process"
)

type Record struct {
	Underlying                        string
	UnderlyingValue                   float64
	ExpiryDate                        time.Time
	Timestamp                         time.Time
	StrikePrice                       float64
	PEAskPrice                        float64
	PEAskQuantity                     int
	PEBidQuantity                     int
	PEBidPrice                        float64
	PEChange                          float64
	PEChangeInOpenInterest            float64
	PEImpliedVolatility               float64
	PELastTradedPrice                 float64
	PEOpenInterest                    float64
	PEPercentangeChangeInPrice        float64
	PEPercentangeChangeInOpenInterest float64
	PETotalBuyQuantity                int
	PETotalSellQuantity               int
	PETotalTradedVolume               int
	CEAskPrice                        float64
	CEAskQuantity                     int
	CEBidQuantity                     int
	CEBidPrice                        float64
	CEChange                          float64
	CEChangeInOpenInterest            float64
	CEImpliedVolatility               float64
	CELastTradedPrice                 float64
	CEOpenInterest                    float64
	CEPercentangeChangeInPrice        float64
	CEPercentangeChangeInOpenInterest float64
	CETotalBuyQuantity                int
	CETotalSellQuantity               int
	CETotalTradedVolume               int
}

func transform(data process.Data) []Record {
	records := make([]Record, 0, len(data.Records))

	for _, dataRecord := range data.Records {
		r := Record{
			Underlying:                        data.Underlying,
			UnderlyingValue:                   data.UnderlyingValue,
			ExpiryDate:                        data.ExpiryDate,
			Timestamp:                         data.Timestamp,
			StrikePrice:                       dataRecord.StrikePrice,
			PEAskPrice:                        dataRecord.PE.AskPrice,
			PEAskQuantity:                     dataRecord.PE.AskQuantity,
			PEBidQuantity:                     dataRecord.PE.BidQuantity,
			PEBidPrice:                        dataRecord.PE.BidPrice,
			PEChange:                          dataRecord.PE.Change,
			PEChangeInOpenInterest:            dataRecord.PE.ChangeInOpenInterest,
			PEImpliedVolatility:               dataRecord.PE.ImpliedVolatility,
			PELastTradedPrice:                 dataRecord.PE.LastTradedPrice,
			PEOpenInterest:                    dataRecord.PE.OpenInterest,
			PEPercentangeChangeInPrice:        dataRecord.PE.PercentangeChangeInPrice,
			PEPercentangeChangeInOpenInterest: dataRecord.PE.PercentangeChangeInOpenInterest,
			PETotalBuyQuantity:                dataRecord.PE.TotalBuyQuantity,
			PETotalSellQuantity:               dataRecord.PE.TotalSellQuantity,
			PETotalTradedVolume:               dataRecord.PE.TotalTradedVolume,
			CEAskPrice:                        dataRecord.CE.AskPrice,
			CEAskQuantity:                     dataRecord.CE.AskQuantity,
			CEBidQuantity:                     dataRecord.CE.BidQuantity,
			CEBidPrice:                        dataRecord.CE.BidPrice,
			CEChange:                          dataRecord.CE.Change,
			CEChangeInOpenInterest:            dataRecord.CE.ChangeInOpenInterest,
			CEImpliedVolatility:               dataRecord.CE.ImpliedVolatility,
			CELastTradedPrice:                 dataRecord.CE.LastTradedPrice,
			CEOpenInterest:                    dataRecord.CE.OpenInterest,
			CEPercentangeChangeInPrice:        dataRecord.CE.PercentangeChangeInPrice,
			CEPercentangeChangeInOpenInterest: dataRecord.CE.PercentangeChangeInOpenInterest,
			CETotalBuyQuantity:                dataRecord.CE.TotalBuyQuantity,
			CETotalSellQuantity:               dataRecord.CE.TotalSellQuantity,
			CETotalTradedVolume:               dataRecord.CE.TotalTradedVolume,
		}
		records = append(records, r)
	}

	return records
}
