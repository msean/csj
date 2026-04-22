package trongrid

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const USDTContract = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"

type (
	TronResponse struct {
		Data []TronResponseData `json:"data"`
		Meta struct {
			Fingerprint string `json:"fingerprint"`
		} `json:"meta"`
		Success bool `json:"success"`
	}
	TronResponseData struct {
		TransactionID  string `json:"transaction_id"`
		BlockTimestamp int64  `json:"block_timestamp"`
		From           string `json:"from"`
		To             string `json:"to"`
		Value          string `json:"value"`
		Type           string `json:"type"`
		TokenInfo      struct {
			Decimals int    `json:"decimals"`
			Symbol   string `json:"symbol"`
			Address  string `json:"address"`
			Name     string `json:"name"`
		} `json:"token_info"`
	}
)

type DayStat struct {
	In  float64
	Out float64
}

func FetchTransactions(account string, limit int, contract string) (*TronResponse, error) {

	var url string
	if contract != "" {
		url = fmt.Sprintf(
			"https://api.trongrid.io/v1/accounts/%s/transactions/trc20?limit=%d&contract_address=%s",
			account, limit, contract,
		)
	} else {
		url = fmt.Sprintf(
			"https://api.trongrid.io/v1/accounts/%s/transactions/trc20?limit=%d",
			account, limit,
		)
	}

	limiter := getTronLimiter()
	client := getHTTPClient()

	var lastErr error

	for attempt := 1; attempt <= 3; attempt++ {

		limiter.Wait(context.Background())

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("http error: %v (attempt %d/3)", err, attempt)
			time.Sleep(1000 * time.Millisecond)
			continue
		}

		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("status %d from API (attempt %d/3)", resp.StatusCode, attempt)
			time.Sleep(1000 * time.Millisecond)
			continue
		}

		var result TronResponse
		if err := json.Unmarshal(body, &result); err != nil {
			lastErr = fmt.Errorf("json parse error: %v", err)
			continue
		}

		if !result.Success {
			lastErr = fmt.Errorf("api returned success=false (attempt %d/3)", attempt)
			time.Sleep(1000 * time.Millisecond)
			continue
		}

		return &result, nil
	}

	return nil, lastErr
}
