package fetch

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseUrl = "https://www.nseindia.com/api/option-chain-indices?symbol="

type fetcher struct {
	symbol          string
	url             string
	client          http.Client
	request         *http.Request
	nRetries        int
	refetchInterval time.Duration
}

func newFetcher(symbol string, nRetries int, refetchInterval time.Duration) (fetcher, error) {
	url := baseUrl + symbol
	client := http.Client{}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fetcher{}, err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")

	fetcher := fetcher{
		symbol:          symbol,
		url:             url,
		client:          client,
		request:         request,
		nRetries:        nRetries,
		refetchInterval: refetchInterval,
	}

	return fetcher, nil
}

func (f *fetcher) Symbol() string {
	return f.symbol
}

func (fetcher *fetcher) RawFetch() (Fetched, error) {
	resp, err := fetcher.client.Do(fetcher.request)
	if err != nil {
		return Fetched{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Fetched{}, err
	}

	fetched, err := ParseJson(body)
	if err != nil {
		return Fetched{}, err
	}

	return fetched, nil
}

func (f *fetcher) Fetch() (Fetched, error) {
	refetchTimer := time.NewTimer(f.refetchInterval)
	for r := 0; r <= f.nRetries; r++ {
		fetched, err := f.RawFetch()
		if err == nil {
			return fetched, nil
		}

		refetchTimer.Reset(f.refetchInterval)
		<-refetchTimer.C
	}

	return Fetched{}, fmt.Errorf("Unable to fetch data for %v after %v retries.", f.symbol, f.nRetries)
}
