import service from '@/utils/request'

// @Tags AppVersion
// @Summary 创建App版本
// @Router /appVersion/createAppVersion [post]
export const createAppVersion = (data) => {
  return service({
    url: '/appVersion/createAppVersion',
    method: 'post',
    data
  })
}

// @Tags AppVersion
// @Summary 删除App版本
// @Router /appVersion/deleteAppVersion [delete]
export const deleteAppVersion = (params) => {
  return service({
    url: '/appVersion/deleteAppVersion',
    method: 'delete',
    params
  })
}

// @Tags AppVersion
// @Summary 批量删除App版本
// @Router /appVersion/deleteAppVersionByIds [delete]
export const deleteAppVersionByIds = (params) => {
  return service({
    url: '/appVersion/deleteAppVersionByIds',
    method: 'delete',
    params
  })
}

// @Tags AppVersion
// @Summary 更新App版本
// @Router /appVersion/updateAppVersion [put]
export const updateAppVersion = (data) => {
  return service({
    url: '/appVersion/updateAppVersion',
    method: 'put',
    data
  })
}

// @Tags AppVersion
// @Summary 用id查询App版本
// @Router /appVersion/findAppVersion [get]
export const findAppVersion = (params) => {
  return service({
    url: '/appVersion/findAppVersion',
    method: 'get',
    params
  })
}

// @Tags AppVersion
// @Summary 分页获取App版本列表
// @Router /appVersion/getAppVersionList [get]
export const getAppVersionList = (params) => {
  return service({
    url: '/appVersion/getAppVersionList',
    method: 'get',
    params
  })
}
