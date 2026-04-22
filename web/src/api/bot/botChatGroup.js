import service from '@/utils/request'
// @Tags BotChatGroup
// @Summary 创建机器人群组列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotChatGroup true "创建机器人群组列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /botChatGroup/createBotChatGroup [post]
export const createBotChatGroup = (data) => {
  return service({
    url: '/botChatGroup/createBotChatGroup',
    method: 'post',
    data
  })
}

// @Tags BotChatGroup
// @Summary 删除机器人群组列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotChatGroup true "删除机器人群组列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /botChatGroup/deleteBotChatGroup [delete]
export const deleteBotChatGroup = (params) => {
  return service({
    url: '/botChatGroup/deleteBotChatGroup',
    method: 'delete',
    params
  })
}

// @Tags BotChatGroup
// @Summary 批量删除机器人群组列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除机器人群组列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /botChatGroup/deleteBotChatGroup [delete]
export const deleteBotChatGroupByIds = (params) => {
  return service({
    url: '/botChatGroup/deleteBotChatGroupByIds',
    method: 'delete',
    params
  })
}

// @Tags BotChatGroup
// @Summary 更新机器人群组列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotChatGroup true "更新机器人群组列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /botChatGroup/updateBotChatGroup [put]
export const updateBotChatGroup = (data) => {
  return service({
    url: '/botChatGroup/updateBotChatGroup',
    method: 'put',
    data
  })
}

// @Tags BotChatGroup
// @Summary 用id查询机器人群组列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.BotChatGroup true "用id查询机器人群组列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /botChatGroup/findBotChatGroup [get]
export const findBotChatGroup = (params) => {
  return service({
    url: '/botChatGroup/findBotChatGroup',
    method: 'get',
    params
  })
}

// @Tags BotChatGroup
// @Summary 分页获取机器人群组列表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取机器人群组列表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /botChatGroup/getBotChatGroupList [get]
export const getBotChatGroupList = (params) => {
  return service({
    url: '/botChatGroup/getBotChatGroupList',
    method: 'get',
    params
  })
}

// @Tags BotChatGroup
// @Summary 不需要鉴权的机器人群组列表接口
// @Accept application/json
// @Produce application/json
// @Param data query botReq.BotChatGroupSearch true "分页获取机器人群组列表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /botChatGroup/getBotChatGroupPublic [get]
export const getBotChatGroupPublic = () => {
  return service({
    url: '/botChatGroup/getBotChatGroupPublic',
    method: 'get',
  })
}

// getChatMessageList 获取聊天消息
export const getChatMessageList = (params) => {
  return service({
    url: '/botChatGroup/chatHistory',
    method: 'get',
    params
  })
}


export const getBotChatGroupClassifyList = (params) => {
  return service({
    url: '/botChatGroup/getBotChatGroupClassifyList',
    method: 'get',
    params
  })
}

export const saveBotChatGroupClassify = (data) => {
  return service({
    url: '/botChatGroup/saveBotChatGroupClassify',
    method: 'post',
    data
  })
}

export const deleteBotChatGroupClassify = (data) => {
  return service({
    url: '/botChatGroup/deleteBotChatGroupClassify',
    method: 'delete',
    data
  })
}

export const chooseChatGroupClassify = (params) => {
  return service({
    url: '/botChatGroup/chooseChatGroupClassify',
    method: 'get',
    params
  })
}
