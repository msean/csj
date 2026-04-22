<template>
  <div>
    <!-- 搜索区域 -->
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        class="demo-form-inline"
        @keyup.enter="onSubmit"
      >
        <el-form-item label="创建日期" prop="createdAtRange">
          <template #label>
            <span>
              创建日期
              <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
                <el-icon><QuestionFilled /></el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-date-picker
            v-model="searchInfo.createdAtRange"
            class="!w-380px"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
          />
        </el-form-item>

        <el-form-item label="机器人" prop="botID">
          <el-select
            v-model="searchInfo.botID"
            filterable
            clearable
            placeholder="请选择机器人"
            style="width: 220px"
          >
            <el-option
              v-for="item in botList"
              :key="item.value"
              :label="item.name"
              :value="item.botID"
            />
          </el-select>
        </el-form-item>

         <el-form-item label="用户ID" prop="userID">
        <el-input
          v-model.number="searchInfo.userID"
          clearable
          placeholder="请输入用户ID"
          style="width: 180px"
        />
      </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-select
            v-model="searchInfo.status"
            clearable
            placeholder="请选择状态"
            style="width: 160px"
          >
            <el-option label="创建" :value="1" />
            <el-option label="完成" :value="2" />
            <el-option label="超时" :value="3" />
            <el-option label="取消" :value="5" />
          </el-select>
        </el-form-item>

        <template v-if="showAllQuery">
          <!-- 可控显示的更多查询条件 -->
        </template>
        <el-form-item>
          <el-button type="primary" icon="search" size="medium" @click="onSubmit">
            查询
          </el-button>
          <el-button icon="refresh" size="medium" @click="onReset">
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 表格操作区 -->
    <div class="gva-table-box">
      <!-- 数据表格 -->
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="机器人" prop="botName" width="180" />
        <el-table-column align="left" label="用户名称" prop="userName" width="180" />
        <el-table-column align="left" label="用户ID" prop="userID" width="120" />
        <el-table-column align="left" label="价格" prop="price" width="120" />
        <el-table-column
          align="left"
          label="状态"
          prop="status"
          width="120"
          :formatter="statusFormatter"
        />
        <el-table-column align="left" label="收款地址" prop="paymentAddr" width="400" />
        <el-table-column align="left" label="交易ID" prop="txID" width="400" />
        <!-- <el-table-column label="发布内容" width="120">
          <template #default="scope">
            <el-button type="primary" size="mini" @click="showContent(scope.row.publishContent)">
              查看
            </el-button>
          </template>
        </el-table-column> -->
        <el-table-column sortable align="left" label="日期" prop="createdAt" width="180">
          <template #default="scope">{{ formatDate(scope.row.createdAt) }}</template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <!-- 查看弹窗 -->
    <el-dialog v-model="dialogVisible" title="发布内容" width="50%">
      <div v-for="(item, index) in dialogContent" :key="index" style="margin-bottom: 10px;">
        <div v-if="item.type === 'text'" style="white-space: pre-wrap;">{{ item.text }}</div>
        <div v-else-if="item.type === 'photo'">
          <img :src="item.file_id" alt="" style="max-width: 100%;" />
        </div>
        <div v-else-if="item.type === 'video'">
          <video controls :src="item.file_id" style="max-width: 100%;"></video>
        </div>
      </div>
      <template #footer>
        <el-button @click="dialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { QuestionFilled } from '@element-plus/icons-vue'
import { getUserRechargeRecordList, deleteUserRechargeRecordByIds } from '@/api/recharge/userRechargeRecord'
import { formatDate as formatDateUtil } from '@/utils/format'
import { getBotChoice } from '@/api/bot/bot'


const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const searchInfo = ref({})
const showAllQuery = ref(false)
const multipleSelection = ref([])

const dialogVisible = ref(false)
const dialogContent = ref([])

const statusMap = { 1: '创建', 2: '完成', 3: '超时', 5: '取消' }

const elSearchFormRef = ref()

function formatDate(dateStr) {
  return formatDateUtil(dateStr)
}

function statusFormatter(row, column, cellValue) {
  return statusMap[cellValue] || '未知'
}

function showContent(content) {
  if (!content) {
    ElMessage.warning('没有发布内容')
    return
  }

  let data = []
  try {
    data = Array.isArray(content) ? content : JSON.parse(content)
  } catch (e) {
    console.error('解析失败', e)
    data = [{ type: 'text', text: String(content) }]
  }
  dialogContent.value = data
  dialogVisible.value = true
}

async function getTableData() {
  const res = await getUserRechargeRecordList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  })
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
    page.value = res.data.page
    pageSize.value = res.data.pageSize
  }
}

const botList = ref([])

// ======= 获取机器人列表 =======
const getBotList = async () => {
  const res = await getBotChoice()
  if (res.code === 0) {
    botList.value = res.data || []
  }
}

function onSubmit() {
  elSearchFormRef.value?.validate(valid => {
    if (!valid) return
    page.value = 1
    getTableData()
  })
}

function onReset() {
  searchInfo.value = {}
  getTableData()
}

function handleSizeChange(val) {
  pageSize.value = val
  getTableData()
}

function handleCurrentChange(val) {
  page.value = val
  getTableData()
}

function handleSelectionChange(val) {
  multipleSelection.value = val
}

async function onDelete() {
  if (!multipleSelection.value.length) return
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const IDs = multipleSelection.value.map(i => i.ID)
    const res = await deleteUserRechargeRecordByIds({ IDs })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '删除成功' })
      if (tableData.value.length === IDs.length && page.value > 1) page.value--
      getTableData()
    }
  })
}

onMounted(() => {
  getBotList()
})


// 初始化数据
getTableData()
</script>

<style>
.gva-search-box {
  margin-bottom: 20px;
}
.gva-pagination {
  margin-top: 10px;
  text-align: right;
}
</style>
