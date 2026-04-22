package trongrid

import (
	"fmt"
	"strconv"
	"time"

	"github.com/msean/csj/backend/global"
	"go.uber.org/zap"
)

func GetDayStart(t time.Time) int64 {
	y, m, d := t.Date()
	loc := t.Location()
	return time.Date(y, m, d, 0, 0, 0, 0, loc).UnixMilli()
}

func pow10(n int) int64 {
	res := int64(1)
	for i := 0; i < n; i++ {
		res *= 10
	}
	return res
}

func ParseAmount(value string, decimals int) float64 {
	var v float64
	fmt.Sscanf(value, "%f", &v)
	for i := 0; i < decimals; i++ {
		v /= 10
	}
	return v
}

func CalcTodayYesterday(address string) (*DayStat, *DayStat, error) {

	now := time.Now()
	todayStart := GetDayStart(now)
	yesterdayStart := GetDayStart(now.AddDate(0, 0, -1))

	todayStat := &DayStat{}
	yesterdayStat := &DayStat{}

	fingerprint := ""
	limit := 200

	for {

		url := fmt.Sprintf(
			"https://api.trongrid.io/v1/accounts/%s/transactions/trc20?limit=%d&contract_address=%s",
			address, limit, USDTContract,
		)

		if fingerprint != "" {
			url += "&fingerprint=" + fingerprint
		}

		global.GVA_LOG.Debug("CalcTodayYesterday", zap.Any("url", url))

		var resp TronResponse
		addIsUnvalid, err := doRequest(url, &resp)
		if err != nil {
			global.GVA_LOG.Error("CalcTodayYesterday", zap.Error(err), zap.Any("url", url))
			return nil, nil, err
		}
		if addIsUnvalid {
			return nil, nil, nil
		}

		if len(resp.Data) == 0 {
			break
		}

		for _, tx := range resp.Data {

			ts := tx.BlockTimestamp

			// 🚀 提前终止（性能关键）
			if ts < yesterdayStart {
				return todayStat, yesterdayStat, nil
			}

			// 金额解析（自动 decimals）
			val, _ := strconv.ParseFloat(tx.Value, 64)
			decimals := tx.TokenInfo.Decimals
			amount := val / float64(pow10(decimals))

			from := tx.From
			to := tx.To

			// ===== 今日 =====
			if ts >= todayStart {

				if to == address {
					todayStat.In += amount
				}
				if from == address {
					todayStat.Out += amount
				}
				continue
			}

			// ===== 昨日 =====
			if ts >= yesterdayStart {

				if to == address {
					yesterdayStat.In += amount
				}
				if from == address {
					yesterdayStat.Out += amount
				}
			}
		}

		// 分页结束
		if resp.Meta.Fingerprint == "" {
			break
		}

		fingerprint = resp.Meta.Fingerprint
	}

	return todayStat, yesterdayStat, nil
}
