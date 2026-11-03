<template>
  <div class="gva-table-box" v-loading="loading">
    <h2 class="mb-4">联系客服配置</h2>
    <el-alert class="mb-4" title="保存后，移动端“我的 → 联系客服”将展示这些信息；留空保存可清除对应内容。" type="info" :closable="false" />
    <el-form :model="form" label-position="top" style="max-width: 640px">
      <el-form-item label="客服电话"><el-input v-model="form.phone" maxlength="50" show-word-limit /></el-form-item>
      <el-form-item label="客服微信"><el-input v-model="form.wechat" maxlength="100" show-word-limit /></el-form-item>
      <el-form-item label="服务说明 / 工作时间"><el-input v-model="form.description" type="textarea" :rows="5" maxlength="2000" show-word-limit /></el-form-item>
      <el-button type="primary" :loading="saving" :disabled="loading || !loaded" @click="save">保存</el-button>
      <el-button :disabled="saving || loading" @click="load">重新加载</el-button>
    </el-form>
  </div>
</template>
<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import service from '@/utils/request'
defineOptions({ name: 'SupportConfig' })
const form = ref({ phone: '', wechat: '', description: '' })
const loading = ref(false), saving = ref(false), loaded = ref(false)
async function load() {
  if (loading.value || saving.value) return
  loading.value = true
  loaded.value = false
  try {
    const res = await service({ url: '/officialMessage/supportConfig', method: 'get' })
    if (res.code === 0) { form.value = { phone: '', wechat: '', description: '', ...res.data }; loaded.value = true }
  } finally { loading.value = false }
}
async function save() {
  if (!loaded.value || saving.value || loading.value) return
  saving.value = true
  try {
    const res = await service({ url: '/officialMessage/supportConfig', method: 'put', data: form.value })
    if (res.code === 0) ElMessage.success('保存成功')
  } finally { saving.value = false }
}
load()
</script>
