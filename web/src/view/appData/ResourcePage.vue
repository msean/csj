<template>
  <div>
    <div class="gva-search-box">
      <h2 class="text-lg mb-4">{{ config.title }}</h2>
      <el-form :inline="true" @submit.prevent="search">
        <el-form-item label="所属用户"><OwnerPicker v-model="filters.ownerUser" /></el-form-item>
        <el-form-item label="关键词"><el-input v-model="filters.keyword" clearable placeholder="名称、号码或 UID" @keyup.enter="search" /></el-form-item>
        <el-form-item v-if="config.customer" label="客户 UID"><el-input v-model="filters.customerUUID" clearable /></el-form-item>
        <el-form-item v-if="config.batch" label="批次 UID"><el-input v-model="filters.batchUUID" clearable /></el-form-item>
        <el-form-item v-if="config.states" label="状态"><el-select v-model="filters.status" clearable placeholder="全部" style="width: 140px"><el-option v-for="state in config.states" :key="state.value" :value="state.value" :label="state.label" /></el-select></el-form-item>
        <el-form-item label="创建时间"><el-date-picker v-model="dates" type="datetimerange" start-placeholder="开始时间" end-placeholder="结束时间" /></el-form-item>
        <el-form-item><el-button type="primary" native-type="submit" :loading="loading">查询</el-button><el-button @click="reset">重置</el-button></el-form-item>
      </el-form>
      <el-alert v-if="config.sharded && !filters.ownerUser" type="info" :closable="false" title="请先选择所属用户，再查询订单或还款记录。" />
    </div>
    <div class="gva-table-box">
      <el-table :data="rows" v-loading="loading" row-key="uid" empty-text="暂无记录">
        <el-table-column prop="uid" label="UID" min-width="220" show-overflow-tooltip />
        <el-table-column v-for="column in config.columns" :key="column.key" :label="column.label" min-width="150" show-overflow-tooltip>
          <template #default="{ row }">{{ display(row[column.key], column.kind) }}</template>
        </el-table-column>
        <el-table-column prop="ownerUser" label="所属用户 UID" min-width="220" show-overflow-tooltip />
        <el-table-column label="创建时间" min-width="180"><template #default="{ row }">{{ formatDate(row.createdAt) }}</template></el-table-column>
        <el-table-column label="操作" fixed="right" width="90"><template #default="{ row }"><el-button type="primary" link @click="openDetail(row)">详情</el-button></template></el-table-column>
      </el-table>
      <div class="gva-pagination"><el-pagination v-model:current-page="page" v-model:page-size="pageSize" :total="total" :page-sizes="[10, 30, 50, 100]" layout="total, sizes, prev, pager, next, jumper" @current-change="load" @size-change="search" /></div>
    </div>
    <el-drawer v-model="detailVisible" :title="`${config.title}详情`" size="min(900px, 95vw)" destroy-on-close>
      <div v-loading="detailLoading">
        <RecordDetail v-if="detail.record" :record="detail.record" :enums="detailEnums" />
        <div v-for="(value, key) in detail.related || {}" :key="key" class="mt-5">
          <h3 class="mb-3">{{ relatedLabels[key] || key }}</h3>
          <template v-if="Array.isArray(value)"><RecordDetail v-for="(record, index) in value" :key="record.uid || index" :record="record" class="mb-3" /><el-empty v-if="!value.length" description="暂无记录" /></template>
          <RecordDetail v-else :record="value" />
        </div>
        <div class="mt-5" v-if="detail.record">
          <el-button v-if="resource === 'customers'" @click="relatedRoute('orders', { customerUUID: detail.record.uid })">查看客户订单</el-button>
          <el-button v-if="resource === 'customers'" @click="relatedRoute('payments', { customerUUID: detail.record.uid })">查看客户还款</el-button>
          <el-button v-if="resource === 'batches'" @click="relatedRoute('orders', { batchUUID: detail.record.uid })">查看批次订单</el-button>
        </div>
      </div>
    </el-drawer>
  </div>
</template>
<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getAppDataList, getAppDataDetail } from '@/api/appData'
import { formatDate } from '@/utils/format'
import OwnerPicker from '@/components/appData/OwnerPicker.vue'
import RecordDetail from '@/components/appData/RecordDetail.vue'
import { resources } from './config'
const props = defineProps({ resource: { type: String, required: true } })
const route = useRoute(), router = useRouter(), config = computed(() => resources[props.resource])
const filters = ref({}), dates = ref([]), rows = ref([]), page = ref(1), pageSize = ref(10), total = ref(0), loading = ref(false)
const detailVisible = ref(false), detailLoading = ref(false), detail = ref({})
const detailEnums = computed(() => {
  const result = {}
  if (config.value.states) result[props.resource === 'payments' ? 'isRevoked' : props.resource === 'messages' ? 'isRead' : 'status'] = Object.fromEntries(config.value.states.map((s) => [s.value, s.label]))
  if (props.resource === 'goods') result.type = { 1: '定装', 2: '散装' }
  if (props.resource === 'messages') result.type = { 1: '还款', 2: '撤销还款' }
  return result
})
const relatedLabels = { goodsList: '商品明细', history: '操作历史', customer: '客户信息' }
let sequence = 0, detailSequence = 0
async function load() {
  const current = ++sequence
  if (config.value.sharded && !filters.value.ownerUser) { rows.value = []; total.value = 0; loading.value = false; return }
  loading.value = true; rows.value = []; total.value = 0
  try {
    const params = { ...filters.value, page: page.value, pageSize: pageSize.value }
    if (dates.value?.length === 2) { params.from = new Date(dates.value[0]).toISOString(); params.to = new Date(dates.value[1]).toISOString() }
    const res = await getAppDataList(props.resource, params)
    if (current === sequence && res.code === 0) { rows.value = res.data.list || []; total.value = res.data.total }
  } finally { if (current === sequence) loading.value = false }
}
function search() {
  if (config.value.sharded && !filters.value.ownerUser) { ElMessage.info('请先选择所属用户'); return }
  page.value = 1; load()
}
function reset() { filters.value = {}; dates.value = []; page.value = 1; load() }
async function openDetail(row) {
  const current = ++detailSequence
  detail.value = {}; detailVisible.value = true; detailLoading.value = true
  try { const res = await getAppDataDetail(props.resource, { uid: row.uid, ownerUser: row.ownerUser }); if (current === detailSequence && res.code === 0) detail.value = res.data }
  finally { if (current === detailSequence) detailLoading.value = false }
}
function relatedRoute(resource, query) { detailVisible.value = false; router.push({ name: resources[resource].name, query: { ownerUser: detail.value.record.ownerUser, ...query } }) }
function display(value, kind) {
  if (value === null || value === undefined || value === '') return '—'
  if (kind === 'money') return Number(value).toFixed(2)
  if (kind === 'unix') return Number(value) > 0 ? formatDate(new Date(Number(value) * 1000)) : '—'
  const enums = { goodsType: { 1: '定装', 2: '散装' }, payType: { 1: '微信', 2: '支付宝', 3: '银行', 4: '现金' }, messageType: { 1: '还款', 2: '撤销还款' } }
  if (kind === 'status') return config.value.states?.find((s) => s.value === value)?.label || value
  return enums[kind]?.[value] || value
}
watch(() => [props.resource, route.fullPath], () => {
  sequence++; detailSequence++; detailVisible.value = false; detailLoading.value = false
  filters.value = { ownerUser: String(route.query.ownerUser || ''), customerUUID: String(route.query.customerUUID || ''), batchUUID: String(route.query.batchUUID || '') }
  dates.value = []; page.value = 1; load()
}, { immediate: true })
</script>
