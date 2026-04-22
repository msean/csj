import service from '@/utils/request'
// @Tags BotTask
// @Summary 创建任务列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotTask true "创建任务列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /task/createBotTask [post]
export const createBotTask = (data) => {
  return service({
    url: '/task/createBotTask',
    method: 'post',
    data
  })
}

// @Tags BotTask
// @Summary 删除任务列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotTask true "删除任务列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /task/deleteBotTask [delete]
export const deleteBotTask = (params) => {
  return service({
    url: '/task/deleteBotTask',
    method: 'delete',
    params
  })
}

// @Tags BotTask
// @Summary 批量删除任务列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除任务列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /task/deleteBotTask [delete]
export const deleteBotTaskByIds = (params) => {
  return service({
    url: '/task/deleteBotTaskByIds',
    method: 'delete',
    params
  })
}

// @Tags BotTask
// @Summary 更新任务列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotTask true "更新任务列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /task/updateBotTask [put]
export const updateBotTask = (data) => {
  return service({
    url: '/task/updateBotTask',
    method: 'put',
    data
  })
}

// @Tags BotTask
// @Summary 用id查询任务列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.BotTask true "用id查询任务列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /task/findBotTask [get]
export const findBotTask = (params) => {
  return service({
    url: '/task/findBotTask',
    method: 'get',
    params
  })
}

// @Tags BotTask
// @Summary 分页获取任务列表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取任务列表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/getBotTaskList [get]
export const getBotTaskList = (params) => {
  return service({
    url: '/task/getBotTaskList',
    method: 'get',
    params
  })
}

// @Tags BotTask
// @Summary 不需要鉴权的任务列表接口
// @Accept application/json
// @Produce application/json
// @Param data query botReq.BotTaskSearch true "分页获取任务列表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /task/getBotTaskPublic [get]
export const getBotTaskPublic = () => {
  return service({
    url: '/task/getBotTaskPublic',
    method: 'get',
  })
}
