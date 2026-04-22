import service from '@/utils/request'
// @Tags BotBanGroupMem
// @Summary 创建封禁成员设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotBanGroupMem true "创建封禁成员设置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /botBanGroupMem/createBotBanGroupMem [post]
export const createBotBanGroupMem = (data) => {
  return service({
    url: '/botBanGroupMem/createBotBanGroupMem',
    method: 'post',
    data
  })
}

// @Tags BotBanGroupMem
// @Summary 删除封禁成员设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotBanGroupMem true "删除封禁成员设置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /botBanGroupMem/deleteBotBanGroupMem [delete]
export const deleteBotBanGroupMem = (params) => {
  return service({
    url: '/botBanGroupMem/deleteBotBanGroupMem',
    method: 'delete',
    params
  })
}

// @Tags BotBanGroupMem
// @Summary 批量删除封禁成员设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除封禁成员设置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /botBanGroupMem/deleteBotBanGroupMem [delete]
export const deleteBotBanGroupMemByIds = (params) => {
  return service({
    url: '/botBanGroupMem/deleteBotBanGroupMemByIds',
    method: 'delete',
    params
  })
}

// @Tags BotBanGroupMem
// @Summary 更新封禁成员设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotBanGroupMem true "更新封禁成员设置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /botBanGroupMem/updateBotBanGroupMem [put]
export const updateBotBanGroupMem = (data) => {
  return service({
    url: '/botBanGroupMem/updateBotBanGroupMem',
    method: 'put',
    data
  })
}

// @Tags BotBanGroupMem
// @Summary 用id查询封禁成员设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.BotBanGroupMem true "用id查询封禁成员设置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /botBanGroupMem/findBotBanGroupMem [get]
export const findBotBanGroupMem = (params) => {
  return service({
    url: '/botBanGroupMem/findBotBanGroupMem',
    method: 'get',
    params
  })
}

// @Tags BotBanGroupMem
// @Summary 分页获取封禁成员设置列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取封禁成员设置列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /botBanGroupMem/getBotBanGroupMemList [get]
export const getBotBanGroupMemList = (params) => {
  return service({
    url: '/botBanGroupMem/getBotBanGroupMemList',
    method: 'get',
    params
  })
}

// @Tags BotBanGroupMem
// @Summary 不需要鉴权的封禁成员设置接口
// @Accept application/json
// @Produce application/json
// @Param data query botReq.BotBanGroupMemSearch true "分页获取封禁成员设置列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /botBanGroupMem/getBotBanGroupMemPublic [get]
export const getBotBanGroupMemPublic = () => {
  return service({
    url: '/botBanGroupMem/getBotBanGroupMemPublic',
    method: 'get',
  })
}
