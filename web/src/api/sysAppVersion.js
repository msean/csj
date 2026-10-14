import service from '@/utils/request'

// @Tags SysAppVersion
// @Summary 创建App版本
// @Router /sysAppVersion/createSysAppVersion [post]
export const createSysAppVersion = (data) => {
  return service({
    url: '/sysAppVersion/createSysAppVersion',
    method: 'post',
    data
  })
}

// @Tags SysAppVersion
// @Summary 删除App版本
// @Router /sysAppVersion/deleteSysAppVersion [delete]
export const deleteSysAppVersion = (params) => {
  return service({
    url: '/sysAppVersion/deleteSysAppVersion',
    method: 'delete',
    params
  })
}

// @Tags SysAppVersion
// @Summary 批量删除App版本
// @Router /sysAppVersion/deleteSysAppVersionByIds [delete]
export const deleteSysAppVersionByIds = (params) => {
  return service({
    url: '/sysAppVersion/deleteSysAppVersionByIds',
    method: 'delete',
    params
  })
}

// @Tags SysAppVersion
// @Summary 更新App版本
// @Router /sysAppVersion/updateSysAppVersion [put]
export const updateSysAppVersion = (data) => {
  return service({
    url: '/sysAppVersion/updateSysAppVersion',
    method: 'put',
    data
  })
}

// @Tags SysAppVersion
// @Summary 用id查询App版本
// @Router /sysAppVersion/findSysAppVersion [get]
export const findSysAppVersion = (params) => {
  return service({
    url: '/sysAppVersion/findSysAppVersion',
    method: 'get',
    params
  })
}

// @Tags SysAppVersion
// @Summary 分页获取App版本列表
// @Router /sysAppVersion/getSysAppVersionList [get]
export const getSysAppVersionList = (params) => {
  return service({
    url: '/sysAppVersion/getSysAppVersionList',
    method: 'get',
    params
  })
}
