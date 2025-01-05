package request

type (
	CreateBatchOrderPayParam struct {
		BatchOrderUUID string  `json:"batchOrderUUID"`
		CustomerUUID   string  `json:"customerUUID"`
		Amount         float64 `json:"amount"`
		PayType        int32   `json:"payType"`
	}
	UpdateBatchOrderPayParam struct {
		BatchOrderPayUUID string  `json:"uuid"`
		Amount            float64 `json:"amount"`
		PayType           int32   `json:"payType"`
	}
	GetBatchOrderPayParam struct {
		BatchOrderPayUUID string `json:"uuid"`
	}
)
