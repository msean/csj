<template>
  <div class="border border-solid border-gray-100 h-full z-10">
    <Toolbar
      :editor="editorRef"
      :default-config="toolbarConfig"
      mode="default"
    />
    <Editor
      v-model="valueHtml"
      class="overflow-y-hidden mt-0.5"
      style="height: 18rem"
      :default-config="editorConfig"
      mode="default"
      @onCreated="handleCreated"
      @onChange="handleEditorChange"
    />
  </div>
</template>

<script setup>
import '@wangeditor/editor/dist/css/style.css'
import { ref, shallowRef, onBeforeUnmount, watch } from 'vue'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/pinia/modules/user'

const props = defineProps({
  modelValue: { type: String, default: '' },
})
const emits = defineEmits(['update:modelValue'])

const base = import.meta.env.VITE_BASE_API
const userStore = useUserStore()

const editorRef = shallowRef()
const valueHtml = ref(props.modelValue)

const toolbarConfig = {
  // 可自定义工具栏，避免报错 translate-page 等
  excludeKeys: ['save-page', 'translate-page']
}

const editorConfig = {
  placeholder: '请输入内容...',
  MENU_CONF: {}
}

editorConfig.MENU_CONF['uploadImage'] = {
  fieldName: 'file',
  server: base + '/public/uploadMedia',
  headers: { 'x-token': userStore.token },
  maxFileSize: 20 * 1024 * 1024,
  customInsert(res, insertFn) {
    if (!res.url) {
      ElMessage.error(res.msg || '上传失败')
      return
    }
    insertFn(res.url)
    emits('update:modelValue', editorRef.value.getHtml())
  }
}

editorConfig.MENU_CONF['uploadVideo'] = {
  fieldName: 'file',
  server: base + '/public/uploadMedia',
  headers: { 'x-token': userStore.token },
  maxFileSize: 100 * 1024 * 1024,
  customInsert(res, insertFn) {
    if (!res.url) {
      ElMessage.error(res.msg || '上传失败')
      return
    }
    insertFn(res.url)
    emits('update:modelValue', editorRef.value.getHtml())
  }
}

const handleCreated = (editor) => {
  editorRef.value = editor
  editorRef.value?.txt.html(props.modelValue || '')
}

const handleEditorChange = (editor) => {
  valueHtml.value = editor.getHtml()
  emits('update:modelValue', valueHtml.value)
}

onBeforeUnmount(() => {
  editorRef.value?.destroy()
})

watch(() => props.modelValue, (val) => {
  if (val !== valueHtml.value) {
    valueHtml.value = val
  }
})
</script>

<style scoped lang="scss">
/* 可根据项目自定义样式 */
</style>
