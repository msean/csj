import service from '@/utils/request'
// @Tags BotCmdConfig
// @Summary 创建机器人命令配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotCmdConfig true "创建机器人命令配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /botCmdConfig/createBotCmdConfig [post]
export const createBotCmdConfig = (data) => {
  return service({
    url: '/botCmdConfig/createBotCmdConfig',
    method: 'post',
    data
  })
}

// @Tags BotCmdConfig
// @Summary 删除机器人命令配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotCmdConfig true "删除机器人命令配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /botCmdConfig/deleteBotCmdConfig [delete]
export const deleteBotCmdConfig = (params) => {
  return service({
    url: '/botCmdConfig/deleteBotCmdConfig',
    method: 'delete',
    params
  })
}

// @Tags BotCmdConfig
// @Summary 批量删除机器人命令配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除机器人命令配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /botCmdConfig/deleteBotCmdConfig [delete]
export const deleteBotCmdConfigByIds = (params) => {
  return service({
    url: '/botCmdConfig/deleteBotCmdConfigByIds',
    method: 'delete',
    params
  })
}

// @Tags BotCmdConfig
// @Summary 更新机器人命令配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotCmdConfig true "更新机器人命令配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /botCmdConfig/updateBotCmdConfig [put]
export const updateBotCmdConfig = (data) => {
  return service({
    url: '/botCmdConfig/updateBotCmdConfig',
    method: 'put',
    data
  })
}

// @Tags BotCmdConfig
// @Summary 用id查询机器人命令配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.BotCmdConfig true "用id查询机器人命令配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /botCmdConfig/findBotCmdConfig [get]
export const findBotCmdConfig = (params) => {
  return service({
    url: '/botCmdConfig/findBotCmdConfig',
    method: 'get',
    params
  })
}

// @Tags BotCmdConfig
// @Summary 分页获取机器人命令配置列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取机器人命令配置列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /botCmdConfig/getBotCmdConfigList [get]
export const getBotCmdConfigList = (params) => {
  return service({
    url: '/botCmdConfig/getBotCmdConfigList',
    method: 'get',
    params
  })
}

// @Tags BotCmdConfig
// @Summary 不需要鉴权的机器人命令配置接口
// @Accept application/json
// @Produce application/json
// @Param data query botReq.BotCmdConfigSearch true "分页获取机器人命令配置列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /botCmdConfig/getBotCmdConfigPublic [get]
export const getBotCmdConfigPublic = () => {
  return service({
    url: '/botCmdConfig/getBotCmdConfigPublic',
    method: 'get',
  })
}
