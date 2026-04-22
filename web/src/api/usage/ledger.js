import service from '@/utils/request'
// @Tags Ledger
// @Summary 创建帐薄
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Ledger true "创建帐薄"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /ledger/createLedger [post]
export const createLedger = (data) => {
  return service({
    url: '/ledger/createLedger',
    method: 'post',
    data
  })
}

// @Tags Ledger
// @Summary 删除帐薄
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Ledger true "删除帐薄"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /ledger/deleteLedger [delete]
export const deleteLedger = (params) => {
  return service({
    url: '/ledger/deleteLedger',
    method: 'delete',
    params
  })
}

// @Tags Ledger
// @Summary 批量删除帐薄
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除帐薄"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /ledger/deleteLedger [delete]
export const deleteLedgerByIds = (params) => {
  return service({
    url: '/ledger/deleteLedgerByIds',
    method: 'delete',
    params
  })
}

// @Tags Ledger
// @Summary 更新帐薄
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Ledger true "更新帐薄"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /ledger/updateLedger [put]
export const updateLedger = (data) => {
  return service({
    url: '/ledger/updateLedger',
    method: 'put',
    data
  })
}

// @Tags Ledger
// @Summary 用id查询帐薄
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.Ledger true "用id查询帐薄"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /ledger/findLedger [get]
export const findLedger = (params) => {
  return service({
    url: '/ledger/findLedger',
    method: 'get',
    params
  })
}

// @Tags Ledger
// @Summary 分页获取帐薄列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取帐薄列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ledger/getLedgerList [get]
export const getLedgerList = (params) => {
  return service({
    url: '/ledger/getLedgerList',
    method: 'get',
    params
  })
}

// @Tags Ledger
// @Summary 不需要鉴权的帐薄接口
// @Accept application/json
// @Produce application/json
// @Param data query usageReq.LedgerSearch true "分页获取帐薄列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /ledger/getLedgerPublic [get]
export const getLedgerPublic = () => {
  return service({
    url: '/ledger/getLedgerPublic',
    method: 'get',
  })
}


export const getLedgerFull = (params) => {
  return service({
    url: '/ledger/full',
    method: 'get',
    params
  })
}
