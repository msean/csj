import service from '@/utils/request'
// @Tags Ledger
// @Summary 创建帐薄
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Ledger true "创建帐薄"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /ledger/createLedger [post]
export const getListenChoice = (data) => {
  return service({
    url: '/listen/choice',
    method: 'get',
    data
  })
}


export const getListenList = (params) => {
  return service({
    url: '/listen/query',
    method: 'get',
    params
  })
}

export const exportListen = (data) => {
  return service({
    url: '/listen/export',
    method: 'post',
    data
  })
}

