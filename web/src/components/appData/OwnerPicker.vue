<template>
  <el-select :model-value="modelValue" @update:model-value="$emit('update:modelValue', $event || '')"
    filterable remote clearable allow-create default-first-option :remote-method="search" :loading="loading"
    placeholder="搜索用户姓名、手机，或输入用户 UID" style="width: 330px" @visible-change="onVisible">
    <el-option v-for="user in users" :key="user.uid" :value="user.uid" :label="`${user.name || '未命名'} / ${user.phone} / ${user.uid}`" />
  </el-select>
</template>
<script setup>
import { ref } from 'vue'
import { getAppUserList } from '@/api/appUser'
defineProps({ modelValue: { type: String, default: '' } })
defineEmits(['update:modelValue'])
const loading = ref(false), users = ref([])
let sequence = 0
async function search(term = '') {
  const current = ++sequence
  loading.value = true
  try {
    const filter = /^\d+$/.test(term) ? { phone: term } : { name: term }
    const res = await getAppUserList({ page: 1, pageSize: 30, ...filter })
    if (current === sequence && res.code === 0) users.value = res.data.list || []
  } finally { if (current === sequence) loading.value = false }
}
function onVisible(open) { if (open && !users.value.length) search() }
</script>
