package write

import (
	"database/sql"
	"log"

	"github.com/ArchitBhonsle/go-fpw/process"
)

func NewWriter(dbPath string) func(data process.Data) (struct{}, error) {
	db := newDB(dbPath)

	return func(data process.Data) (struct{}, error) {
		log.Println(data.Underlying, "writing")

		records := transform(data)
		for _, record := range records {
			err := writeRecord(record, db)
			if err != nil {
				return struct{}{}, err
			}
		}
		return struct{}{}, nil
	}
}

func NewWriterAny(dbPath string) func(dataAny any) (any, error) {
	db := newDB(dbPath)

	return func(dataAny any) (any, error) {
		data := dataAny.(process.Data)

		log.Println(data.Underlying, "writing")

		records := transform(data)
		for _, record := range records {
			err := writeRecord(record, db)
			if err != nil {
				return struct{}{}, err
			}
		}
		return struct{}{}, nil
	}
}

func writeRecord(record Record, db *sql.DB) error {
	_, err := db.Exec(
		insertRecordSQL,
		record.Underlying,
		record.UnderlyingValue,
		record.ExpiryDate,
		record.Timestamp,
		record.StrikePrice,
		record.PEAskPrice,
		record.PEAskQuantity,
		record.PEBidQuantity,
		record.PEBidPrice,
		record.PEChange,
		record.PEChangeInOpenInterest,
		record.PEImpliedVolatility,
		record.PELastTradedPrice,
		record.PEOpenInterest,
		record.PEPercentangeChangeInPrice,
		record.PEPercentangeChangeInOpenInterest,
		record.PETotalBuyQuantity,
		record.PETotalSellQuantity,
		record.PETotalTradedVolume,
		record.CEAskPrice,
		record.CEAskQuantity,
		record.CEBidQuantity,
		record.CEBidPrice,
		record.CEChange,
		record.CEChangeInOpenInterest,
		record.CEImpliedVolatility,
		record.CELastTradedPrice,
		record.CEOpenInterest,
		record.CEPercentangeChangeInPrice,
		record.CEPercentangeChangeInOpenInterest,
		record.CETotalBuyQuantity,
		record.CETotalSellQuantity,
		record.CETotalTradedVolume,
	)

	return err
}
