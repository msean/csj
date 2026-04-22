package trongrid

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	limiterOnce sync.Once
	trxLimiter  *rate.Limiter
	clientOnce  sync.Once
	httpClient  *http.Client
)

func getTronLimiter() *rate.Limiter {
	limiterOnce.Do(func() {
		trxLimiter = rate.NewLimiter(rate.Limit(1), 1)
	})
	return trxLimiter
}

func getHTTPClient() *http.Client {
	clientOnce.Do(func() {
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	})
	return httpClient
}

// addIsUnInvalid 是
func doRequest(url string, result interface{}) (addIsUnInvalid bool, err error) {

	limiter := getTronLimiter()
	client := getHTTPClient()

	var lastErr error

	for i := 0; i < 3; i++ {

		if err := limiter.Wait(context.Background()); err != nil {
			return false, err
		}

		resp, err := client.Get(url)
		if err != nil {
			lastErr = err
			time.Sleep(time.Millisecond * 1100)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		// ===== 1. 先解析错误结构 =====
		var errResp struct {
			Success    bool   `json:"success"`
			Error      string `json:"error"`
			StatusCode int    `json:"statusCode"`
		}

		_ = json.Unmarshal(body, &errResp)

		// 🚨 地址无效（你要的逻辑）
		if errResp.Error == "A valid account address is required." {
			return true, nil
		}

		// ===== 2. HTTP错误 =====
		if resp.StatusCode != 200 {
			lastErr = fmt.Errorf("http error: %d", resp.StatusCode)
			time.Sleep(time.Millisecond * 1100)
			continue
		}

		// ===== 3. 正常解析 =====
		if err := json.Unmarshal(body, result); err != nil {
			lastErr = err
			continue
		}

		return false, nil
	}

	return false, lastErr
}
