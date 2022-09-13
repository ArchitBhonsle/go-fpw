package fetch

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseUrl = "https://www.nseindia.com/api/option-chain-indices?symbol="

type Fetcher struct {
	symbol          string
	url             string
	client          http.Client
	request         *http.Request
	nRetries        int
	refetchInterval time.Duration
}

func NewFetcher(symbol string, nRetries int, refetchInterval time.Duration) (Fetcher, error) {
	url := baseUrl + symbol
	client := http.Client{}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Fetcher{}, err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")

	fetcher := Fetcher{
		symbol:          symbol,
		url:             url,
		client:          client,
		request:         request,
		nRetries:        nRetries,
		refetchInterval: refetchInterval,
	}

	return fetcher, nil
}

func (f *Fetcher) Symbol() string {
	return f.symbol
}

func (fetcher *Fetcher) RawFetch() (Fetched, error) {
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

func (f *Fetcher) Fetch() (Fetched, error) {
	fmt.Println(time.Now().Format("15:04:05.000"), f.symbol, "fetch")

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
