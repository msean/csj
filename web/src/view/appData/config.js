const field = (key, label, kind = 'text') => ({ key, label, kind })
const money = (key, label) => field(key, label, 'money')
const states = (values) => Object.entries(values).map(([value, label]) => ({ value: Number(value), label }))
export const resources = {
  customers: { name: 'AppCustomers', title: '客户管理', columns: [field('name', '客户名称'), field('phone', '手机'), field('carNo', '车牌号'), money('credit', '欠款余额'), money('debt', '原始欠款'), money('totalAmount', '累计订单金额'), field('orderCount', '订单数')] },
  'goods-categories': { name: 'AppGoodsCategories', title: '商品分类', columns: [field('name', '分类名称')] },
  goods: { name: 'AppGoods', title: '商品管理', columns: [field('name', '商品名称'), field('categoryID', '分类 UID'), field('type', '商品类型', 'goodsType'), money('price', '单价'), field('weight', '重量'), field('status', '状态')] },
  batches: { name: 'AppBatches', title: '批次管理', states: states({ 1: '在售', 2: '售罄', 3: '已结算' }), columns: [field('serialNo', '批次序号'), field('serialID', '批次编号'), field('storageTime', '入库时间', 'unix'), field('status', '状态', 'status'), money('statSalesAmount', '销售金额'), money('statSellProfit', '利润')] },
  orders: { name: 'AppOrders', title: '订单查询', sharded: true, customer: true, batch: true, states: states({ 1: '记账单', 2: '历史赊欠单', 3: '已结算', 100: '作废', 101: '退款', 102: '退货' }), columns: [field('customerUUID', '客户 UID'), field('batchUUID', '批次 UID'), field('status', '状态', 'status'), money('totalAmount', '订单金额'), money('creditAmount', '赊欠金额'), money('discount', '优惠'), money('tricycleFee', '三轮车费'), field('orderTime', '下单时间', 'unix')] },
  payments: { name: 'AppPayments', title: '还款查询', sharded: true, customer: true, states: states({ 0: '正常', 1: '已撤销' }), columns: [field('customerUUID', '客户 UID'), field('batchOrderUUID', '订单 UID'), money('amount', '还款金额'), money('discount', '优惠'), field('payType', '付款方式', 'payType'), field('isRevoked', '撤销状态', 'status'), money('beforeCredit', '还款前欠款'), money('afterCredit', '还款后欠款')] },
  messages: { name: 'AppMessages', title: '个人消息', customer: true, states: states({ 0: '未读', 1: '已读' }), columns: [field('event', '事件'), field('content', '内容'), field('type', '消息类型', 'messageType'), field('isRead', '已读状态', 'status'), field('customerUUID', '客户 UID')] },
  shares: { name: 'AppShares', title: '分享记录', customer: true, columns: [field('customerUUID', '客户 UID'), field('shareTime', '最近分享时间', 'unix')] }
}
