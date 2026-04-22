import service from '@/utils/request'

export const createBotBanContent = (data) => {
  return service({
    url: '/bot_ban_content/create',
    method: 'post',
    data
  })
}


export const deleteBotBanContent = (params) => {
  return service({
    url: '/bot_ban_content/delete',
    method: 'delete',
    params
  })
}


export const deleteBotBanContentByIds = (params) => {
  return service({
    url: '/bot_ban_content/delete_by_ids',
    method: 'delete',
    params
  })
}


export const updateBotBanContent = (data) => {
  return service({
    url: '/bot_ban_content/update',
    method: 'put',
    data
  })
}


export const findBotBanContent = (params) => {
  return service({
    url: '/bot_ban_content/get',
    method: 'get',
    params
  })
}


export const getBotBanContentList = (params) => {
  return service({
    url: '/bot_ban_content/list',
    method: 'get',
    params
  })
}


export const getBotBanContentPublic = () => {
  return service({
    url: '/bot_ban_content/public_get',
    method: 'get',
  })
}
