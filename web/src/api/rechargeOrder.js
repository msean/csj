import service from '@/utils/request'

// @Tags RechargeOrder
// @Summary 用uid查询充值订单
// @Router /rechargeOrder/findRechargeOrder [get]
export const findRechargeOrder = (params) => {
  return service({
    url: '/rechargeOrder/findRechargeOrder',
    method: 'get',
    params
  })
}

// @Tags RechargeOrder
// @Summary 分页获取充值订单列表
// @Router /rechargeOrder/getRechargeOrderList [get]
export const getRechargeOrderList = (params) => {
  return service({
    url: '/rechargeOrder/getRechargeOrderList',
    method: 'get',
    params
  })
}
