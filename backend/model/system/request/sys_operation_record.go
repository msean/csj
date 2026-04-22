package request

import (
	"github.com/msean/csj/backend/model/common/request"
	"github.com/msean/csj/backend/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
