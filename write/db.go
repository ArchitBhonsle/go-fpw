package write

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const createRecordsSQL = `
CREATE TABLE Records (
	id                                INTEGER PRIMARY KEY AUTOINCREMENT,
	underlying                        TEXT NOT NULL,
	underlyingValue                   REAL NOT NULL,
	expiryDate                        DATE NOT NULL,
	timestamp                         DATE NOT NULL,
	strikePrice                       REAL NOT NULL,
	peAskPrice                        REAL NOT NULL,
	peAskQuantity                     INT  NOT NULL,
	peBidQuantity                     INT  NOT NULL,
	peBidPrice                        REAL NOT NULL,
	peChange                          REAL NOT NULL,
	peChangeInOpenInterest            REAL NOT NULL,
	peImpliedVolatility               REAL NOT NULL,
	peLastTradedPrice                 REAL NOT NULL,
	peOpenInterest                    REAL NOT NULL,
	pePercentangeChangeInPrice        REAL NOT NULL,
	pePercentangeChangeInOpenInterest REAL NOT NULL,
	peTotalBuyQuantity                INT  NOT NULL,
	peTotalSellQuantity               INT  NOT NULL,
	peTotalTradedVolume               INT  NOT NULL,
	ceAskPrice                        REAL NOT NULL,
	ceAskQuantity                     INT  NOT NULL,
	ceBidQuantity                     INT  NOT NULL,
	ceBidPrice                        REAL NOT NULL,
	ceChange                          REAL NOT NULL,
	ceChangeInOpenInterest            REAL NOT NULL,
	ceImpliedVolatility               REAL NOT NULL,
	ceLastTradedPrice                 REAL NOT NULL,
	ceOpenInterest                    REAL NOT NULL,
	cePercentangeChangeInPrice        REAL NOT NULL,
	cePercentangeChangeInOpenInterest REAL NOT NULL,
	ceTotalBuyQuantity                INT  NOT NULL,
	ceTotalSellQuantity               INT  NOT NULL,
	ceTotalTradedVolume               INT  NOT NULL
);`

const insertRecordSQL = `
INSERT INTO Records (
	underlying,
	underlyingValue,
	expiryDate,
	timestamp,
	strikePrice,
	peAskPrice,
	peAskQuantity,
	peBidQuantity,
	peBidPrice,
	peChange,
	peChangeInOpenInterest,
	peImpliedVolatility,
	peLastTradedPrice,
	peOpenInterest,
	pePercentangeChangeInPrice,
	pePercentangeChangeInOpenInterest,
	peTotalBuyQuantity,
	peTotalSellQuantity,
	peTotalTradedVolume,
	ceAskPrice,
	ceAskQuantity,
	ceBidQuantity,
	ceBidPrice,
	ceChange,
	ceChangeInOpenInterest,
	ceImpliedVolatility,
	ceLastTradedPrice,
	ceOpenInterest,
	cePercentangeChangeInPrice,
	cePercentangeChangeInOpenInterest,
	ceTotalBuyQuantity,
	ceTotalSellQuantity,
	ceTotalTradedVolume
) VALUES (
  ?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?
);`

func newDB(dbPath string) *sql.DB {
	createTables := false
	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		createTables = true
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	if createTables {
		_, err = db.Exec(createRecordsSQL)
		if err != nil {
			panic(err)
		}
	}

	return db
}
