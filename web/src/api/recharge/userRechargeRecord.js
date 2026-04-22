import service from '@/utils/request'
// @Tags UserRechargeRecord
// @Summary 创建用户充值记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.UserRechargeRecord true "创建用户充值记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /userRechargeRecord/createUserRechargeRecord [post]
export const createUserRechargeRecord = (data) => {
  return service({
    url: '/userRechargeRecord/createUserRechargeRecord',
    method: 'post',
    data
  })
}

// @Tags UserRechargeRecord
// @Summary 删除用户充值记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.UserRechargeRecord true "删除用户充值记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /userRechargeRecord/deleteUserRechargeRecord [delete]
export const deleteUserRechargeRecord = (params) => {
  return service({
    url: '/userRechargeRecord/deleteUserRechargeRecord',
    method: 'delete',
    params
  })
}

// @Tags UserRechargeRecord
// @Summary 批量删除用户充值记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除用户充值记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /userRechargeRecord/deleteUserRechargeRecord [delete]
export const deleteUserRechargeRecordByIds = (params) => {
  return service({
    url: '/userRechargeRecord/deleteUserRechargeRecordByIds',
    method: 'delete',
    params
  })
}

// @Tags UserRechargeRecord
// @Summary 更新用户充值记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.UserRechargeRecord true "更新用户充值记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /userRechargeRecord/updateUserRechargeRecord [put]
export const updateUserRechargeRecord = (data) => {
  return service({
    url: '/userRechargeRecord/updateUserRechargeRecord',
    method: 'put',
    data
  })
}

// @Tags UserRechargeRecord
// @Summary 用id查询用户充值记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.UserRechargeRecord true "用id查询用户充值记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /userRechargeRecord/findUserRechargeRecord [get]
export const findUserRechargeRecord = (params) => {
  return service({
    url: '/userRechargeRecord/findUserRechargeRecord',
    method: 'get',
    params
  })
}

// @Tags UserRechargeRecord
// @Summary 分页获取用户充值记录列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取用户充值记录列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /userRechargeRecord/getUserRechargeRecordList [get]
export const getUserRechargeRecordList = (params) => {
  return service({
    url: '/userRechargeRecord/getUserRechargeRecordList',
    method: 'get',
    params
  })
}

// @Tags UserRechargeRecord
// @Summary 不需要鉴权的用户充值记录接口
// @Accept application/json
// @Produce application/json
// @Param data query rechargeReq.UserRechargeRecordSearch true "分页获取用户充值记录列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /userRechargeRecord/getUserRechargeRecordPublic [get]
export const getUserRechargeRecordPublic = () => {
  return service({
    url: '/userRechargeRecord/getUserRechargeRecordPublic',
    method: 'get',
  })
}
