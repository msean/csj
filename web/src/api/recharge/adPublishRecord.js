import service from '@/utils/request'
// @Tags AdPublishRecord
// @Summary 创建广告发布记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AdPublishRecord true "创建广告发布记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /adPublishRecord/createAdPublishRecord [post]
export const createAdPublishRecord = (data) => {
  return service({
    url: '/adPublishRecord/createAdPublishRecord',
    method: 'post',
    data
  })
}

// @Tags AdPublishRecord
// @Summary 删除广告发布记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AdPublishRecord true "删除广告发布记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /adPublishRecord/deleteAdPublishRecord [delete]
export const deleteAdPublishRecord = (params) => {
  return service({
    url: '/adPublishRecord/deleteAdPublishRecord',
    method: 'delete',
    params
  })
}

// @Tags AdPublishRecord
// @Summary 批量删除广告发布记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除广告发布记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /adPublishRecord/deleteAdPublishRecord [delete]
export const deleteAdPublishRecordByIds = (params) => {
  return service({
    url: '/adPublishRecord/deleteAdPublishRecordByIds',
    method: 'delete',
    params
  })
}

// @Tags AdPublishRecord
// @Summary 更新广告发布记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AdPublishRecord true "更新广告发布记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /adPublishRecord/updateAdPublishRecord [put]
export const updateAdPublishRecord = (data) => {
  return service({
    url: '/adPublishRecord/updateAdPublishRecord',
    method: 'put',
    data
  })
}

// @Tags AdPublishRecord
// @Summary 用id查询广告发布记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.AdPublishRecord true "用id查询广告发布记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /adPublishRecord/findAdPublishRecord [get]
export const findAdPublishRecord = (params) => {
  return service({
    url: '/adPublishRecord/findAdPublishRecord',
    method: 'get',
    params
  })
}

// @Tags AdPublishRecord
// @Summary 分页获取广告发布记录列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取广告发布记录列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /adPublishRecord/getAdPublishRecordList [get]
export const getAdPublishRecordList = (params) => {
  return service({
    url: '/adPublishRecord/getAdPublishRecordList',
    method: 'get',
    params
  })
}

// @Tags AdPublishRecord
// @Summary 不需要鉴权的广告发布记录接口
// @Accept application/json
// @Produce application/json
// @Param data query rechargeReq.AdPublishRecordSearch true "分页获取广告发布记录列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /adPublishRecord/getAdPublishRecordPublic [get]
export const getAdPublishRecordPublic = () => {
  return service({
    url: '/adPublishRecord/getAdPublishRecordPublic',
    method: 'get',
  })
}
