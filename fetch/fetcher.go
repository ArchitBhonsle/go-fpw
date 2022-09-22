package fetch

import (
	"log"
	"time"
)

func NewFetcher(
	symbol string,
	nRetries int,
	refetchInterval time.Duration,
) func(time.Time) (Fetched, error) {
	fetcher, err := newFetcher(symbol, nRetries, refetchInterval)
	if err != nil {
		panic(err)
	}

	return func(t time.Time) (Fetched, error) {
		log.Println(symbol, "fetching")

		return fetcher.Fetch()
	}
}

func NewFetcherAny(
	symbol string,
	nRetries int,
	refetchInterval time.Duration,
) func(any) (any, error) {
	fetcher, err := newFetcher(symbol, nRetries, refetchInterval)
	if err != nil {
		panic(err)
	}

	return func(tAny any) (any, error) {
		log.Println(symbol, "fetching")

		return fetcher.Fetch()
	}
}
