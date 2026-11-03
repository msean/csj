<template>
  <div class="gva-table-box">
    <div class="gva-btn-list">
      <el-button @click="openSupport">联系客服配置</el-button>
      <el-button type="primary" @click="open()">新增官方消息</el-button>
      <el-select v-model="status" clearable placeholder="全部状态" style="width: 160px" @change="search">
        <el-option label="草稿" :value="0" /><el-option label="已发布" :value="1" />
      </el-select>
    </div>
    <el-table :data="rows" row-key="uid">
      <el-table-column prop="title" label="标题" min-width="200" />
      <el-table-column label="状态" width="100"><template #default="{ row }">{{ row.status === 1 ? '已发布' : '草稿' }}</template></el-table-column>
      <el-table-column label="等级" width="100"><template #default="{ row }"><el-tag :type="row.level === 2 ? 'danger' : 'warning'">{{ row.level === 2 ? '重要' : '一般' }}</el-tag></template></el-table-column>
      <el-table-column prop="sort" label="排序值" width="90" />
      <el-table-column label="发布时间" width="180"><template #default="{ row }">{{ formatDate(row.publishedAt) }}</template></el-table-column>
      <el-table-column label="失效时间" width="180"><template #default="{ row }">{{ row.expiresAt ? formatDate(row.expiresAt) : '长期有效' }}</template></el-table-column>
      <el-table-column label="操作" width="160"><template #default="{ row }">
        <el-button link type="primary" @click="open(row)">编辑</el-button>
        <el-button link type="danger" @click="remove(row)">删除</el-button>
      </template></el-table-column>
    </el-table>
    <el-pagination class="gva-pagination" v-model:current-page="page" :page-size="pageSize" :total="total" layout="total, prev, pager, next" @current-change="load" />
    <el-drawer v-model="visible" title="官方消息" size="600px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="标题" prop="title"><el-input v-model="form.title" maxlength="200" show-word-limit /></el-form-item>
        <el-form-item label="内容" prop="content"><el-input v-model="form.content" type="textarea" :rows="8" /></el-form-item>
        <el-form-item label="通知等级"><el-radio-group v-model="form.level"><el-radio :value="1">一般（橙色）</el-radio><el-radio :value="2">重要（红色）</el-radio></el-radio-group></el-form-item>
        <el-form-item label="发布状态"><el-radio-group v-model="form.status"><el-radio :value="0">草稿</el-radio><el-radio :value="1">发布</el-radio></el-radio-group></el-form-item>
        <el-form-item label="发布时间" prop="publishedAt"><el-date-picker v-model="form.publishedAt" type="datetime" :clearable="false" /></el-form-item>
        <el-form-item label="失效时间（不填则长期有效）"><el-date-picker v-model="form.expiresAt" type="datetime" /></el-form-item>
        <el-form-item label="排序值（越大越靠前）"><el-input-number v-model="form.sort" :precision="0" /></el-form-item>
        <p class="mb-4">官方通知在移动端首页喇叭及独立列表显示，不进入个人消息中心。通知不区分已读和未读。</p>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </el-form>
    </el-drawer>
    <el-dialog v-model="supportVisible" title="联系客服配置" width="560px">
      <el-form label-position="top">
        <el-form-item label="客服电话"><el-input v-model="support.phone" maxlength="50" /></el-form-item>
        <el-form-item label="客服微信"><el-input v-model="support.wechat" maxlength="100" /></el-form-item>
        <el-form-item label="服务说明 / 工作时间"><el-input v-model="support.description" type="textarea" :rows="4" maxlength="2000" /></el-form-item>
      </el-form>
      <template #footer><el-button :loading="supportSaving" type="primary" @click="saveSupport">保存</el-button></template>
    </el-dialog>
  </div>
</template>
<script setup>
import { ref } from 'vue'
import service from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate } from '@/utils/format'
import { createOfficialMessage, updateOfficialMessage, deleteOfficialMessage, getOfficialMessageList } from '@/api/officialMessage'
defineOptions({ name: 'OfficialMessage' })
const rows = ref([]), total = ref(0), page = ref(1), status = ref(), pageSize = 10
const visible = ref(false), saving = ref(false), formRef = ref(), form = ref({})
const rules = {
  title: [{ required: true, whitespace: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, whitespace: true, message: '请输入内容', trigger: 'blur' }],
  publishedAt: [{ required: true, message: '请选择发布时间', trigger: 'change' }]
}
async function load() {
  const res = await getOfficialMessageList({ page: page.value, pageSize, status: status.value })
  if (res.code === 0) { rows.value = res.data.list || []; total.value = res.data.total }
}
function search() { page.value = 1; load() }
function open(row) {
  form.value = row ? { ...row, level: row.level || 1 } : { level: 1, title: '', content: '', status: 0, sort: 0, publishedAt: new Date(), expiresAt: null }
  visible.value = true
}
async function save() {
  if (saving.value || !await formRef.value.validate().catch(() => false)) return
  const data = { ...form.value, publishedAt: new Date(form.value.publishedAt).toISOString(), expiresAt: form.value.expiresAt ? new Date(form.value.expiresAt).toISOString() : null }
  if (data.expiresAt && new Date(data.expiresAt) <= new Date(data.publishedAt)) { ElMessage.error('失效时间必须晚于发布时间'); return }
  saving.value = true
  try {
    const res = await (data.uid ? updateOfficialMessage(data) : createOfficialMessage(data))
    if (res.code === 0) { visible.value = false; ElMessage.success('保存成功'); await load() }
  } finally { saving.value = false }
}
async function remove(row) {
  try { await ElMessageBox.confirm('确定删除这条官方消息？', '提示', { type: 'warning' }) } catch { return }
  const res = await deleteOfficialMessage({ uid: row.uid })
  if (res.code === 0) { if (rows.value.length === 1 && page.value > 1) page.value--; await load() }
}
const supportVisible = ref(false), supportSaving = ref(false), support = ref({ phone: '', wechat: '', description: '' })
async function openSupport() {
 const res = await service({ url: '/officialMessage/supportConfig', method: 'get' })
 if (res.code === 0) { support.value = { phone: '', wechat: '', description: '', ...res.data }; supportVisible.value = true }
}
async function saveSupport() {
 if (supportSaving.value) return
 supportSaving.value = true
 try {
  const res = await service({ url: '/officialMessage/supportConfig', method: 'put', data: support.value })
  if (res.code === 0) { ElMessage.success('保存成功'); supportVisible.value = false }
 } finally { supportSaving.value = false }
}
load()
</script>
