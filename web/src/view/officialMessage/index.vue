<template>
  <div class="gva-table-box">
    <div class="gva-btn-list">
      <el-button type="primary" @click="open()">新增官方消息</el-button>
      <el-select v-model="status" clearable placeholder="全部状态" style="width: 160px" @change="search">
        <el-option label="草稿" :value="0" /><el-option label="已发布" :value="1" />
      </el-select>
    </div>
    <el-table :data="rows" row-key="uid">
      <el-table-column prop="title" label="标题" min-width="200" />
      <el-table-column label="状态" width="100"><template #default="{ row }">{{ row.status === 1 ? '已发布' : '草稿' }}</template></el-table-column>
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
        <el-form-item label="发布状态"><el-radio-group v-model="form.status"><el-radio :value="0">草稿</el-radio><el-radio :value="1">发布</el-radio></el-radio-group></el-form-item>
        <el-form-item label="发布时间" prop="publishedAt"><el-date-picker v-model="form.publishedAt" type="datetime" :clearable="false" /></el-form-item>
        <el-form-item label="失效时间（不填则长期有效）"><el-date-picker v-model="form.expiresAt" type="datetime" /></el-form-item>
        <el-form-item label="排序值（越大越靠前）"><el-input-number v-model="form.sort" :precision="0" /></el-form-item>
        <p class="mb-4">官方消息置顶显示。修改内容不会重置用户已读状态，需要重新提醒时请新增消息。</p>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </el-form>
    </el-drawer>
  </div>
</template>
<script setup>
import { ref } from 'vue'
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
  form.value = row ? { ...row } : { title: '', content: '', status: 0, sort: 0, publishedAt: new Date(), expiresAt: null }
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
load()
</script>
