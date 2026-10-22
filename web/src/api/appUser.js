import service from '@/utils/request'

// @Tags AppUser
// @Summary 用uid查询App注册用户
// @Router /appUser/findAppUser [get]
export const findAppUser = (params) => {
  return service({
    url: '/appUser/findAppUser',
    method: 'get',
    params
  })
}

// @Tags AppUser
// @Summary 分页获取App注册用户列表
// @Router /appUser/getAppUserList [get]
export const getAppUserList = (params) => {
  return service({
    url: '/appUser/getAppUserList',
    method: 'get',
    params
  })
}
