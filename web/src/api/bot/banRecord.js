import service from '@/utils/request'
// @Tags BanRecord
// @Summary 创建封禁记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BanRecord true "创建封禁记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /banRecord/createBanRecord [post]
export const createBanRecord = (data) => {
  return service({
    url: '/banRecord/createBanRecord',
    method: 'post',
    data
  })
}

// @Tags BanRecord
// @Summary 删除封禁记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BanRecord true "删除封禁记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /banRecord/deleteBanRecord [delete]
export const deleteBanRecord = (params) => {
  return service({
    url: '/banRecord/deleteBanRecord',
    method: 'delete',
    params
  })
}

// @Tags BanRecord
// @Summary 批量删除封禁记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除封禁记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /banRecord/deleteBanRecord [delete]
export const deleteBanRecordByIds = (params) => {
  return service({
    url: '/banRecord/deleteBanRecordByIds',
    method: 'delete',
    params
  })
}

// @Tags BanRecord
// @Summary 更新封禁记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BanRecord true "更新封禁记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /banRecord/updateBanRecord [put]
export const updateBanRecord = (data) => {
  return service({
    url: '/banRecord/updateBanRecord',
    method: 'put',
    data
  })
}

// @Tags BanRecord
// @Summary 用id查询封禁记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.BanRecord true "用id查询封禁记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /banRecord/findBanRecord [get]
export const findBanRecord = (params) => {
  return service({
    url: '/banRecord/findBanRecord',
    method: 'get',
    params
  })
}

// @Tags BanRecord
// @Summary 分页获取封禁记录列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取封禁记录列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /banRecord/getBanRecordList [get]
export const getBanRecordList = (params) => {
  return service({
    url: '/banRecord/getBanRecordList',
    method: 'get',
    params
  })
}

// @Tags BanRecord
// @Summary 不需要鉴权的封禁记录接口
// @Accept application/json
// @Produce application/json
// @Param data query botReq.BanRecordSearch true "分页获取封禁记录列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /banRecord/getBanRecordPublic [get]
export const getBanRecordPublic = () => {
  return service({
    url: '/banRecord/getBanRecordPublic',
    method: 'get',
  })
}
