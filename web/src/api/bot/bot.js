import service from '@/utils/request'

export const createBot = (data) => {
  return service({
    url: '/bot_mgr/create',
    method: 'post',
    data
  })
}


export const deleteBot = (params) => {
  return service({
    url: '/bot_mgr/delete',
    method: 'delete',
    params
  })
}

export const deleteBotByIds = (params) => {
  return service({
    url: '/bot_mgr/delete_by_ids',
    method: 'delete',
    params
  })
}


export const updateBot = (data) => {
  return service({
    url: '/bot_mgr/update',
    method: 'put',
    data
  })
}


export const findBot = (params) => {
  return service({
    url: '/bot_mgr/get',
    method: 'get',
    params
  })
}


export const getBotList = (params) => {
  return service({
    url: '/bot_mgr/list',
    method: 'get',
    params
  })
}


export const getBotPublic = () => {
  return service({
    url: '/bot_mgr/public_get',
    method: 'get',
  })
}

export const getBotChoice = (params) => {
  return service({
    url: '/bot_mgr/choice',
    method: 'get',
    params
  })
}

export const getBotChoiceWithChatGroup = (params) => {
  return service({
    url: '/bot_mgr/choice_with_chat_group',
    method: 'get',
    params
  })
}

export const unBanUser = (params) => {
  return service({
    url: '/bot_mgr/unban_user',
    method: 'post',
    data: params,
  })
}

export const getUserChoice = (params) => {
  return service({
    url: '/bot_mgr/choice',
    method: 'get',
    params
  })
}