import service from '@/utils/request'
// @Tags LedgerAccountGroup
// @Summary 创建记账账号组
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.LedgerAccountGroup true "创建记账账号组"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /ledgerAccountGroup/createLedgerAccountGroup [post]
export const createLedgerAccountGroup = (data) => {
  return service({
    url: '/ledgerAccountGroup/createLedgerAccountGroup',
    method: 'post',
    data
  })
}

// @Tags LedgerAccountGroup
// @Summary 删除记账账号组
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.LedgerAccountGroup true "删除记账账号组"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /ledgerAccountGroup/deleteLedgerAccountGroup [delete]
export const deleteLedgerAccountGroup = (params) => {
  return service({
    url: '/ledgerAccountGroup/deleteLedgerAccountGroup',
    method: 'delete',
    params
  })
}

// @Tags LedgerAccountGroup
// @Summary 批量删除记账账号组
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除记账账号组"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /ledgerAccountGroup/deleteLedgerAccountGroup [delete]
export const deleteLedgerAccountGroupByIds = (params) => {
  return service({
    url: '/ledgerAccountGroup/deleteLedgerAccountGroupByIds',
    method: 'delete',
    params
  })
}

// @Tags LedgerAccountGroup
// @Summary 更新记账账号组
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.LedgerAccountGroup true "更新记账账号组"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /ledgerAccountGroup/updateLedgerAccountGroup [put]
export const updateLedgerAccountGroup = (data) => {
  return service({
    url: '/ledgerAccountGroup/updateLedgerAccountGroup',
    method: 'put',
    data
  })
}

// @Tags LedgerAccountGroup
// @Summary 用id查询记账账号组
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.LedgerAccountGroup true "用id查询记账账号组"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /ledgerAccountGroup/findLedgerAccountGroup [get]
export const findLedgerAccountGroup = (params) => {
  return service({
    url: '/ledgerAccountGroup/findLedgerAccountGroup',
    method: 'get',
    params
  })
}

// @Tags LedgerAccountGroup
// @Summary 分页获取记账账号组列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取记账账号组列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ledgerAccountGroup/getLedgerAccountGroupList [get]
export const getLedgerAccountGroupList = (params) => {
  return service({
    url: '/ledgerAccountGroup/getLedgerAccountGroupList',
    method: 'get',
    params
  })
}

// @Tags LedgerAccountGroup
// @Summary 不需要鉴权的记账账号组接口
// @Accept application/json
// @Produce application/json
// @Param data query ledgerReq.LedgerAccountGroupSearch true "分页获取记账账号组列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /ledgerAccountGroup/getLedgerAccountGroupPublic [get]
export const getLedgerAccountGroupPublic = () => {
  return service({
    url: '/ledgerAccountGroup/getLedgerAccountGroupPublic',
    method: 'get',
  })
}
