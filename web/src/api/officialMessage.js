import service from '@/utils/request'

// @Tags OfficialMessage
// @Summary 创建官方消息
// @Router /officialMessage/createOfficialMessage [post]
export const createOfficialMessage = (data) => {
  return service({
    url: '/officialMessage/createOfficialMessage',
    method: 'post',
    data
  })
}

// @Tags OfficialMessage
// @Summary 删除官方消息
// @Router /officialMessage/deleteOfficialMessage [delete]
export const deleteOfficialMessage = (params) => {
  return service({
    url: '/officialMessage/deleteOfficialMessage',
    method: 'delete',
    params
  })
}

// @Tags OfficialMessage
// @Summary 批量删除官方消息
// @Router /officialMessage/deleteOfficialMessageByIds [delete]
export const deleteOfficialMessageByIds = (params) => {
  return service({
    url: '/officialMessage/deleteOfficialMessageByIds',
    method: 'delete',
    params
  })
}

// @Tags OfficialMessage
// @Summary 更新官方消息
// @Router /officialMessage/updateOfficialMessage [put]
export const updateOfficialMessage = (data) => {
  return service({
    url: '/officialMessage/updateOfficialMessage',
    method: 'put',
    data
  })
}

// @Tags OfficialMessage
// @Summary 用uid查询官方消息
// @Router /officialMessage/findOfficialMessage [get]
export const findOfficialMessage = (params) => {
  return service({
    url: '/officialMessage/findOfficialMessage',
    method: 'get',
    params
  })
}

// @Tags OfficialMessage
// @Summary 分页获取官方消息列表
// @Router /officialMessage/getOfficialMessageList [get]
export const getOfficialMessageList = (params) => {
  return service({
    url: '/officialMessage/getOfficialMessageList',
    method: 'get',
    params
  })
}
