import service from '@/utils/request'

// @Tags Customers
// @Summary 创建customers表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Customers true "创建customers表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /customers/createCustomers [post]
export const createCustomers = (data) => {
  return service({
    url: '/customers/createCustomers',
    method: 'post',
    data
  })
}

// @Tags Customers
// @Summary 删除customers表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Customers true "删除customers表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /customers/deleteCustomers [delete]
export const deleteCustomers = (params) => {
  return service({
    url: '/customers/deleteCustomers',
    method: 'delete',
    params
  })
}

// @Tags Customers
// @Summary 批量删除customers表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除customers表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /customers/deleteCustomers [delete]
export const deleteCustomersByIds = (params) => {
  return service({
    url: '/customers/deleteCustomersByIds',
    method: 'delete',
    params
  })
}

// @Tags Customers
// @Summary 更新customers表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Customers true "更新customers表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /customers/updateCustomers [put]
export const updateCustomers = (data) => {
  return service({
    url: '/customers/updateCustomers',
    method: 'put',
    data
  })
}

// @Tags Customers
// @Summary 用id查询customers表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Customers true "用id查询customers表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /customers/findCustomers [get]
export const findCustomers = (params) => {
  return service({
    url: '/customers/findCustomers',
    method: 'get',
    params
  })
}

// @Tags Customers
// @Summary 分页获取customers表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取customers表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /customers/getCustomersList [get]
export const getCustomersList = (params) => {
  return service({
    url: '/customers/getCustomersList',
    method: 'get',
    params
  })
}
