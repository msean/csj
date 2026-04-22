<template>
  <div>
    <!-- 搜索 -->
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <!-- 群 / 频道 -->
        <el-form-item label="群 / 频道">
          <el-select
            v-model="searchInfo.groupId"
            placeholder="请选择群 / 频道"
            filterable
            clearable
            style="width:260px"
          >
            <el-option
              v-for="item in chatOptions"
              :key="item.groupId"
              :label="item.groupName"
              :value="item.groupId"
            />
          </el-select>
        </el-form-item>

        <!-- 关键词 -->
        <el-form-item label="关键词">
          <el-input
            v-model="searchInfo.keyword"
            clearable
            placeholder="关键词 / 文本"
            style="width:200px"
          />
        </el-form-item>

        <!-- 时间 -->
        <el-form-item label="时间">
          <el-date-picker
            v-model="searchInfo.timeRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            class="!w-380px"
          />
        </el-form-item>

        <!-- 操作 -->
        <el-form-item>
          <el-button type="primary" @click="onSubmit">查询</el-button>
          <el-button @click="onReset">重置</el-button>
          <el-button type="success" :loading="exportLoading" @click="onExport">
            导出
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 表格 -->
    <div class="gva-table-box">
      <el-table :data="tableData" style="width:100%">
        <el-table-column label="用户ID" width="160">
          <template #default="scope">
            {{ scope.row.user_id }}
          </template>
        </el-table-column>
        <el-table-column label="用户" width="160">
          <template #default="scope">
            {{ scope.row.username }}
          </template>
        </el-table-column>
         <el-table-column label="昵称" width="160">
          <template #default="scope">
            {{ scope.row.nick_name}}
          </template>
        </el-table-column>
        <el-table-column label="消息内容" prop="text" min-width="300" />
        <el-table-column label="时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.timestamp) }}
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { formatDate } from '@/utils/format'
import {
  getListenChoice,
  getListenList,
  exportListen
} from '@/api/listen/listen'

defineOptions({
  name: 'TelegramListen'
})

/* ================= 搜索条件 ================= */
const searchInfo = ref({
  groupId: null,
  keyword: '',
  timeRange: []
})

/* ================= 群 / 频道 ================= */
const chatOptions = ref([])

const loadChatOptions = async () => {
  const res = await getListenChoice()
  if (res.code === 0) {
    chatOptions.value = res.data.map(item => ({
      groupId: item.group_id,
      groupName: item.group_name
    }))
  }
}

/* ================= 表格 ================= */
const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const buildParams = () => {
  const params = {
    page: page.value,
    pageSize: pageSize.value,
    groupId: searchInfo.value.groupId,
    keyword: searchInfo.value.keyword
  }

  if (searchInfo.value.timeRange?.length === 2) {
    params.startTime = searchInfo.value.timeRange[0]
    params.endTime = searchInfo.value.timeRange[1]
  }

  return params
}

const getTableData = async () => {
  if (!searchInfo.value.groupId) return

  const res = await getListenList(buildParams())
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
  }
}

/* ================= 事件 ================= */
const onSubmit = () => {
  if (!searchInfo.value.groupId) {
    ElMessage.warning('请选择群 / 频道')
    return
  }
  page.value = 1
  getTableData()
}

const onReset = () => {
  searchInfo.value = {
    groupId: null,
    keyword: '',
    timeRange: []
  }
  tableData.value = []
  total.value = 0
}

const handleCurrentChange = val => {
  page.value = val
  getTableData()
}

const handleSizeChange = val => {
  pageSize.value = val
  page.value = 1
  getTableData()
}

/* ================= 导出 ================= */
const exportLoading = ref(false)

const onExport = async () => {
  if (!searchInfo.value.groupId) {
    ElMessage.warning('请选择群 / 频道')
    return
  }

  exportLoading.value = true
  try {
    const params = buildParams()
    delete params.page
    delete params.pageSize

    const res = await exportListen(params)

    if (res.code === 0 && res.data?.file) {
      const downloadUrl = `/api/public/download?file=${res.data.file}`

      const a = document.createElement('a')
      a.href = downloadUrl
      a.download = res.data.file
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)

      ElMessage.success('开始下载')
    }
  } finally {
    exportLoading.value = false
  }
}

/* ================= 初始化 ================= */
loadChatOptions()
</script>
