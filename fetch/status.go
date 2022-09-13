package fetch

import (
	"encoding/json"
	"io"
	"net/http"
)

type marketState struct {
	State []marketStateItem `json:"marketState"`
}

type marketStateItem struct {
	Market       string `json:"market"`
	MarketStatus string `json:"marketStatus"`
}

func parseMarketStatus(body []byte) (marketState, error) {
	data := marketState{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

const url string = "https://www.nseindia.com/api/marketStatus"

type StatusChecker struct {
	client  http.Client
	request *http.Request
}

func NewStatusChecker() (*StatusChecker, error) {
	client := http.Client{}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")

	statusChecker := &StatusChecker{
		client:  client,
		request: request,
	}

	return statusChecker, nil
}

func (sc *StatusChecker) Check() bool {
	resp, err := sc.client.Do(sc.request)
	if err != nil {
		return false
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	marketStatus, err := parseMarketStatus(body)
	if err != nil {
		return false
	}

	for _, ms := range marketStatus.State {
		if ms.Market != "Capital Market" {
			continue
		}

		if ms.MarketStatus == "Open" {
			return true
		}
	}

	return false
}
