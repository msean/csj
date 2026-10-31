<template>
  <el-descriptions :column="1" border>
    <el-descriptions-item v-for="(value, key) in record || {}" :key="key" :label="labels[key] || key">
      <pre v-if="value !== null && typeof value === 'object'" class="record-json">{{ JSON.stringify(value, null, 2) }}</pre>
      <span v-else>{{ display(value, key) }}</span>
    </el-descriptions-item>
  </el-descriptions>
</template>
<script setup>
import { formatDate } from '@/utils/format'
const props = defineProps({ record: { type: Object, default: () => ({}) }, labels: { type: Object, default: () => ({}) }, enums: { type: Object, default: () => ({}) } })
const defaultLabels = { uid: 'UID', createdAt: '创建时间', updatedAt: '更新时间', ownerUser: '所属用户', customerUUID: '客户 UID', userUUID: '用户 UID', batchUUID: '批次 UID', goodsUUID: '商品 UID', batchOrderUUID: '订单 UID', name: '名称', phone: '手机', contactPhones: '分享联系方式', province: '省份', marketName: '市场名称', marketAddress: '市场地址', stallName: '档口名称', experienceDeadline: '体验截止时间', vipDeadline: '会员截止时间', remark: '备注', addr: '地址', carNo: '车牌号', debt: '原始欠款', paidDebt: '已还原始欠款', credit: '欠款余额', orderCount: '有效订单数', totalAmount: '总金额', categoryID: '分类 UID', price: '单价', weight: '重量', mount: '件数', status: '状态', serialNo: '批次序号', serialID: '批次编号', storageTime: '入库时间', creditAmount: '赊欠金额', discount: '优惠金额', tricycleFee: '三轮车费', orderTime: '下单时间', sharedTime: '分享时间', shared: '分享状态', amount: '金额', payType: '付款方式', isRevoked: '是否撤销', revokedAt: '撤销时间', revokedReason: '撤销原因', beforeCredit: '还款前欠款', afterCredit: '还款后欠款', payDetails: '还款明细', event: '事件', content: '内容', isRead: '已读状态', relatedUUID: '关联 UID', shareTime: '最近分享时间', orderNo: '充值订单号', memberType: '会员类型', payStatus: '支付状态', payTime: '支付时间', transactionID: '微信交易号', title: '标题', type: '类型', sort: '排序值', publishedAt: '发布时间', expiresAt: '失效时间', version: '版本号', downloadAddr: '下载地址', isForce: '强制更新', description: '更新说明' }
// Explicit labels can override the defaults used by shared entity fields.
const labels = { ...defaultLabels, total: '小计', goods_type: '商品类型', goodsWeight: '定装重量', sortNo: '批次内序号', statMount: '入库总件数', statWeight: '入库总重量', statInventoryMount: '库存件数', statInventoryWeight: '库存重量', statSalesAmount: '销售金额', statSellProfit: '销售利润', statSellWeight: '售出重量', statSellMount: '售出件数', surplus: '剩余库存', sellAmount: '销售数量', sellAvgPrice: '销售均价', sellTotalAmount: '销售金额', sellProfit: '销售利润', steps: '操作记录', goodsList: '商品明细', ...props.labels }
function display(value, key) {
  if (value === null || value === undefined || value === '') return '—'
  const enums = { payType: { 1: '微信', 2: '支付宝', 3: '银行', 4: '现金' }, isRead: { 0: '未读', 1: '已读' }, isRevoked: { 0: '正常', 1: '已撤销' }, shared: { 1: '未分享', 2: '已分享' }, memberType: { 1: '月', 2: '季度', 3: '年' }, payStatus: { 0: '待支付', 1: '已支付', 2: '已取消', 3: '已退款' }, goods_type: { 1: '定装', 2: '散装' }, ...props.enums }
  if (enums[key]?.[value] !== undefined) return enums[key][value]
  if (['createdAt', 'updatedAt', 'experienceDeadline', 'vipDeadline', 'revokedAt', 'payTime', 'publishedAt', 'expiresAt'].includes(key)) return formatDate(value)
  if (['storageTime', 'orderTime', 'sharedTime', 'shareTime'].includes(key) && Number(value) > 0) return formatDate(new Date(Number(value) * 1000))
  if (typeof value === 'boolean') return value ? '是' : '否'
  return value
}
</script>
<style scoped>
.record-json { white-space: pre-wrap; overflow-wrap: anywhere; margin: 0; max-height: 360px; overflow: auto; }
</style>
