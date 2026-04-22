import service from '@/utils/request'
// @Tags BotChannel
// @Summary 创建机器人渠道
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotChannel true "创建机器人渠道"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /botChannel/createBotChannel [post]
export const createBotChannel = (data) => {
  return service({
    url: '/botChannel/createBotChannel',
    method: 'post',
    data
  })
}

// @Tags BotChannel
// @Summary 删除机器人渠道
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotChannel true "删除机器人渠道"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /botChannel/deleteBotChannel [delete]
export const deleteBotChannel = (params) => {
  return service({
    url: '/botChannel/deleteBotChannel',
    method: 'delete',
    params
  })
}

// @Tags BotChannel
// @Summary 批量删除机器人渠道
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除机器人渠道"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /botChannel/deleteBotChannel [delete]
export const deleteBotChannelByIds = (params) => {
  return service({
    url: '/botChannel/deleteBotChannelByIds',
    method: 'delete',
    params
  })
}

// @Tags BotChannel
// @Summary 更新机器人渠道
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotChannel true "更新机器人渠道"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /botChannel/updateBotChannel [put]
export const updateBotChannel = (data) => {
  return service({
    url: '/botChannel/updateBotChannel',
    method: 'put',
    data
  })
}

// @Tags BotChannel
// @Summary 用id查询机器人渠道
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.BotChannel true "用id查询机器人渠道"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /botChannel/findBotChannel [get]
export const findBotChannel = (params) => {
  return service({
    url: '/botChannel/findBotChannel',
    method: 'get',
    params
  })
}

// @Tags BotChannel
// @Summary 分页获取机器人渠道列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取机器人渠道列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /botChannel/getBotChannelList [get]
export const getBotChannelList = (params) => {
  return service({
    url: '/botChannel/getBotChannelList',
    method: 'get',
    params
  })
}

// @Tags BotChannel
// @Summary 不需要鉴权的机器人渠道接口
// @Accept application/json
// @Produce application/json
// @Param data query botReq.BotChannelSearch true "分页获取机器人渠道列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /botChannel/getBotChannelPublic [get]
export const getBotChannelPublic = () => {
  return service({
    url: '/botChannel/getBotChannelPublic',
    method: 'get',
  })
}
