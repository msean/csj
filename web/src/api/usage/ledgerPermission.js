import service from '@/utils/request'
// @Tags LedgerPermission
// @Summary 创建帐薄权限管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.LedgerPermission true "创建帐薄权限管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /ledgerPermission/createLedgerPermission [post]
export const createLedgerPermission = (data) => {
  return service({
    url: '/ledgerPermission/createLedgerPermission',
    method: 'post',
    data
  })
}

// @Tags LedgerPermission
// @Summary 删除帐薄权限管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.LedgerPermission true "删除帐薄权限管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /ledgerPermission/deleteLedgerPermission [delete]
export const deleteLedgerPermission = (params) => {
  return service({
    url: '/ledgerPermission/deleteLedgerPermission',
    method: 'delete',
    params
  })
}

// @Tags LedgerPermission
// @Summary 批量删除帐薄权限管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除帐薄权限管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /ledgerPermission/deleteLedgerPermission [delete]
export const deleteLedgerPermissionByIds = (params) => {
  return service({
    url: '/ledgerPermission/deleteLedgerPermissionByIds',
    method: 'delete',
    params
  })
}

// @Tags LedgerPermission
// @Summary 更新帐薄权限管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.LedgerPermission true "更新帐薄权限管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /ledgerPermission/updateLedgerPermission [put]
export const updateLedgerPermission = (data) => {
  return service({
    url: '/ledgerPermission/updateLedgerPermission',
    method: 'put',
    data
  })
}

// @Tags LedgerPermission
// @Summary 用id查询帐薄权限管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.LedgerPermission true "用id查询帐薄权限管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /ledgerPermission/findLedgerPermission [get]
export const findLedgerPermission = (params) => {
  return service({
    url: '/ledgerPermission/findLedgerPermission',
    method: 'get',
    params
  })
}

// @Tags LedgerPermission
// @Summary 分页获取帐薄权限管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取帐薄权限管理列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ledgerPermission/getLedgerPermissionList [get]
export const getLedgerPermissionList = (params) => {
  return service({
    url: '/ledgerPermission/getLedgerPermissionList',
    method: 'get',
    params
  })
}

// @Tags LedgerPermission
// @Summary 不需要鉴权的帐薄权限管理接口
// @Accept application/json
// @Produce application/json
// @Param data query usageReq.LedgerPermissionSearch true "分页获取帐薄权限管理列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /ledgerPermission/getLedgerPermissionPublic [get]
export const getLedgerPermissionPublic = () => {
  return service({
    url: '/ledgerPermission/getLedgerPermissionPublic',
    method: 'get',
  })
}
