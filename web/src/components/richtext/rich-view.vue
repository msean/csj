<template>
  <div class="rich-view border border-gray-200 p-2 rounded" v-html="sanitizedContent"></div>
</template>

<script setup>
import { computed } from 'vue'
import DOMPurify from 'dompurify'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

// 使用 DOMPurify 对 HTML 做安全处理，防止 XSS
const sanitizedContent = computed(() => {
  return DOMPurify.sanitize(props.modelValue || '')
})
</script>

<style scoped lang="scss">
.rich-view {
  min-height: 200px;
  overflow-y: auto;
  word-break: break-word;
}

.rich-view img {
  max-width: 100%;
  display: block;
  margin-bottom: 5px;
}

.rich-view video {
  max-width: 100%;
  display: block;
  margin-bottom: 5px;
}
</style>
