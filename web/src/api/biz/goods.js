import service from '@/utils/request'

export const getGoodsList = (params) => {
  return service({
    url: '/goods/list_goods',
    method: 'get',
    params
  })
}
