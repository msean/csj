package tronscan

import (
	"fmt"

	"github.com/msean/csj/backend/global"
	"go.uber.org/zap"
)

type (
	trxTransferResp struct {
		Data []TrxTransferData `json:"data"`
	}
	TrxTransferData struct {
		From      string `json:"transferFromAddress"`
		To        string `json:"transferToAddress"`
		Amount    int64  `json:"amount"`
		Timestamp int64  `json:"timestamp"`
	}
	Trc20Resp struct {
		TokenTransfers []Trc20TransferData `json:"token_transfers"`
	}
	Trc20TransferData struct {
		From      string `json:"from_address"`
		To        string `json:"to_address"`
		Amount    string `json:"quant"`
		Timestamp int64  `json:"block_ts"`
	}
)

func GetTRXTransfers(address string, start, end int64) ([]TrxTransferData, error) {

	url := fmt.Sprintf(
		"https://apilist.tronscan.org/api/transfer/trx?address=%s&start_timestamp=%d&end_timestamp=%d&limit=200",
		address, start, end,
	)

	var resp trxTransferResp
	err := doRequest(url, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func GetUSDTTransfers(address string, start, end int64) ([]Trc20TransferData, error) {

	url := fmt.Sprintf(
		"https://apilist.tronscan.org/api/token_trc20/transfers?relatedAddress=%s&start_timestamp=%d&end_timestamp=%d&limit=200&contract_address=TXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj",
		address, start, end,
	)

	var resp Trc20Resp
	err := doRequest(url, &resp)
	global.GVA_LOG.Debug("GetUSDTTransfers", zap.Any("body", resp))
	if err != nil {
		global.GVA_LOG.Error("GetUSDTTransfers", zap.Error(err), zap.Any("url", url))
		return nil, err
	}

	return resp.TokenTransfers, nil
}
