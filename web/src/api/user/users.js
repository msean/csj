import service from '@/utils/request'

// @Tags Users
// @Summary 创建users表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Users true "创建users表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /users/createUsers [post]
export const createUsers = (data) => {
  return service({
    url: '/users/createUsers',
    method: 'post',
    data
  })
}

// @Tags Users
// @Summary 删除users表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Users true "删除users表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /users/deleteUsers [delete]
export const deleteUsers = (params) => {
  return service({
    url: '/users/deleteUsers',
    method: 'delete',
    params
  })
}

// @Tags Users
// @Summary 批量删除users表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除users表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /users/deleteUsers [delete]
export const deleteUsersByIds = (params) => {
  return service({
    url: '/users/deleteUsersByIds',
    method: 'delete',
    params
  })
}

// @Tags Users
// @Summary 更新users表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Users true "更新users表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /users/updateUsers [put]
export const updateUsers = (data) => {
  return service({
    url: '/users/updateUsers',
    method: 'put',
    data
  })
}

// @Tags Users
// @Summary 用id查询users表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Users true "用id查询users表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /users/findUsers [get]
export const findUsers = (params) => {
  return service({
    url: '/users/findUsers',
    method: 'get',
    params
  })
}

// @Tags Users
// @Summary 分页获取users表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取users表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /users/getUsersList [get]
export const getUsersList = (params) => {
  return service({
    url: '/users/getUsersList',
    method: 'get',
    params
  })
}
