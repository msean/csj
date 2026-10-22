import service from '@/utils/request'

// @Tags MemberConfig
// @Summary 创建会员收费配置
// @Router /memberConfig/createMemberConfig [post]
export const createMemberConfig = (data) => {
  return service({
    url: '/memberConfig/createMemberConfig',
    method: 'post',
    data
  })
}

// @Tags MemberConfig
// @Summary 删除会员收费配置
// @Router /memberConfig/deleteMemberConfig [delete]
export const deleteMemberConfig = (params) => {
  return service({
    url: '/memberConfig/deleteMemberConfig',
    method: 'delete',
    params
  })
}

// @Tags MemberConfig
// @Summary 批量删除会员收费配置
// @Router /memberConfig/deleteMemberConfigByIds [delete]
export const deleteMemberConfigByIds = (params) => {
  return service({
    url: '/memberConfig/deleteMemberConfigByIds',
    method: 'delete',
    params
  })
}

// @Tags MemberConfig
// @Summary 更新会员收费配置
// @Router /memberConfig/updateMemberConfig [put]
export const updateMemberConfig = (data) => {
  return service({
    url: '/memberConfig/updateMemberConfig',
    method: 'put',
    data
  })
}

// @Tags MemberConfig
// @Summary 用uid查询会员收费配置
// @Router /memberConfig/findMemberConfig [get]
export const findMemberConfig = (params) => {
  return service({
    url: '/memberConfig/findMemberConfig',
    method: 'get',
    params
  })
}

// @Tags MemberConfig
// @Summary 分页获取会员收费配置列表
// @Router /memberConfig/getMemberConfigList [get]
export const getMemberConfigList = (params) => {
  return service({
    url: '/memberConfig/getMemberConfigList',
    method: 'get',
    params
  })
}
