import service from '@/utils/request'
// @Tags RechargeConfig
// @Summary 创建充值配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.RechargeConfig true "创建充值配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /rechargeConfig/createRechargeConfig [post]
export const createRechargeConfig = (data) => {
  return service({
    url: '/rechargeConfig/createRechargeConfig',
    method: 'post',
    data
  })
}

// @Tags RechargeConfig
// @Summary 删除充值配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.RechargeConfig true "删除充值配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /rechargeConfig/deleteRechargeConfig [delete]
export const deleteRechargeConfig = (params) => {
  return service({
    url: '/rechargeConfig/deleteRechargeConfig',
    method: 'delete',
    params
  })
}

// @Tags RechargeConfig
// @Summary 批量删除充值配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除充值配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /rechargeConfig/deleteRechargeConfig [delete]
export const deleteRechargeConfigByIds = (params) => {
  return service({
    url: '/rechargeConfig/deleteRechargeConfigByIds',
    method: 'delete',
    params
  })
}

// @Tags RechargeConfig
// @Summary 更新充值配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.RechargeConfig true "更新充值配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /rechargeConfig/updateRechargeConfig [put]
export const updateRechargeConfig = (data) => {
  return service({
    url: '/rechargeConfig/updateRechargeConfig',
    method: 'put',
    data
  })
}

// @Tags RechargeConfig
// @Summary 用id查询充值配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.RechargeConfig true "用id查询充值配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /rechargeConfig/findRechargeConfig [get]
export const findRechargeConfig = (params) => {
  return service({
    url: '/rechargeConfig/findRechargeConfig',
    method: 'get',
    params
  })
}

// @Tags RechargeConfig
// @Summary 分页获取充值配置列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取充值配置列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /rechargeConfig/getRechargeConfigList [get]
export const getRechargeConfigList = (params) => {
  return service({
    url: '/rechargeConfig/getRechargeConfigList',
    method: 'get',
    params
  })
}

// @Tags RechargeConfig
// @Summary 不需要鉴权的充值配置接口
// @Accept application/json
// @Produce application/json
// @Param data query rechargeReq.RechargeConfigSearch true "分页获取充值配置列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /rechargeConfig/getRechargeConfigPublic [get]
export const getRechargeConfigPublic = () => {
  return service({
    url: '/rechargeConfig/getRechargeConfigPublic',
    method: 'get',
  })
}
