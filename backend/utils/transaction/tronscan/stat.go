package tronscan

import (
	"strconv"
	"time"
)

type DayStat struct {
	In  float64
	Out float64
}

func CalcTRXStat(address string, start, end int64) (*DayStat, error) {

	list, err := GetTRXTransfers(address, start, end)
	if err != nil {
		return nil, err
	}

	stat := &DayStat{}

	for _, tx := range list {
		amount := float64(tx.Amount) / 1e6

		if tx.To == address {
			stat.In += amount
		}
		if tx.From == address {
			stat.Out += amount
		}
	}

	return stat, nil
}

func CalcZeroBalance(current float64, todayIn, todayOut float64) float64 {
	return current - (todayIn - todayOut)
}

func GetDayRange(t time.Time) (int64, int64) {
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	end := start.Add(24*time.Hour - time.Millisecond)
	return start.UnixMilli(), end.UnixMilli()
}

func CalcUSDTZeroBalance(current float64, todayIn, todayOut float64) float64 {
	return current - (todayIn - todayOut)
}

func CalcUSDTStat(address string, start, end int64) (*DayStat, error) {

	list, err := GetUSDTTransfers(address, start, end)
	if err != nil {
		return nil, err
	}

	stat := &DayStat{}

	for _, tx := range list {

		val, _ := strconv.ParseFloat(tx.Amount, 64)
		amount := val / 1e6

		if tx.To == address {
			stat.In += amount
		}

		if tx.From == address {
			stat.Out += amount
		}
	}

	return stat, nil
}
