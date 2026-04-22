import service from '@/utils/request'
// @Tags BotMsgMass
// @Summary 创建机器人群发
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotMsgMass true "创建机器人群发"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /botMsgMass/createBotMsgMass [post]
export const createBotMsgMass = (data) => {
  return service({
    url: '/botMsgMass/createBotMsgMass',
    method: 'post',
    data
  })
}

// @Tags BotMsgMass
// @Summary 删除机器人群发
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotMsgMass true "删除机器人群发"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /botMsgMass/deleteBotMsgMass [delete]
export const deleteBotMsgMass = (params) => {
  return service({
    url: '/botMsgMass/deleteBotMsgMass',
    method: 'delete',
    params
  })
}

// @Tags BotMsgMass
// @Summary 批量删除机器人群发
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除机器人群发"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /botMsgMass/deleteBotMsgMass [delete]
export const deleteBotMsgMassByIds = (params) => {
  return service({
    url: '/botMsgMass/deleteBotMsgMassByIds',
    method: 'delete',
    params
  })
}

// @Tags BotMsgMass
// @Summary 更新机器人群发
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotMsgMass true "更新机器人群发"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /botMsgMass/updateBotMsgMass [put]
export const updateBotMsgMass = (data) => {
  return service({
    url: '/botMsgMass/updateBotMsgMass',
    method: 'put',
    data
  })
}

// @Tags BotMsgMass
// @Summary 用id查询机器人群发
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.BotMsgMass true "用id查询机器人群发"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /botMsgMass/findBotMsgMass [get]
export const findBotMsgMass = (params) => {
  return service({
    url: '/botMsgMass/findBotMsgMass',
    method: 'get',
    params
  })
}

// @Tags BotMsgMass
// @Summary 分页获取机器人群发列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取机器人群发列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /botMsgMass/getBotMsgMassList [get]
export const getBotMsgMassList = (params) => {
  return service({
    url: '/botMsgMass/getBotMsgMassList',
    method: 'get',
    params
  })
}

// @Tags BotMsgMass
// @Summary 不需要鉴权的机器人群发接口
// @Accept application/json
// @Produce application/json
// @Param data query botReq.BotMsgMassSearch true "分页获取机器人群发列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /botMsgMass/getBotMsgMassPublic [get]
export const getBotMsgMassPublic = () => {
  return service({
    url: '/botMsgMass/getBotMsgMassPublic',
    method: 'get',
  })
}


export const sendBotMsgMass = (data) => {
  return service({
    url: '/botMsgMass/sendBotMsgMass',
    method: 'post',
    data
  })
}


export const massMsgHistory = (params) => {
  return service({
    url: '/botMsgMass/history',
    method: 'get',
    params
  })
}
