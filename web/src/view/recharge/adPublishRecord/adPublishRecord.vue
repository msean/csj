
<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
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
      
       <el-form-item label="用户ID" prop="userID">
        <el-input
          v-model.number="searchInfo.userID"
          clearable
          placeholder="请输入用户ID"
          style="width: 180px"
        />
      </el-form-item>

        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
        >
        <el-table-column type="selection" width="55" />
        
        <el-table-column sortable align="left" label="日期" prop="createdAt" width="180">
            <template #default="scope">{{ formatDate(scope.row.createdAt) }}</template>
        </el-table-column>
            <el-table-column align="left" label="机器人" prop="botName" width="120" />
            <el-table-column align="left" label="发布次数" prop="publishTimes" width="120" />
            <el-table-column align="left" label="发布用户ID" prop="userID" width="120" />
            <el-table-column align="left" label="发布价格" prop="price" width="120" />
            <el-table-column align="left" label="用户名" prop="userName" width="120" />
            <el-table-column align="left" label="频道" prop="channelName" width="120" />
            <el-table-column label="发布内容" prop="content" width="200">
                <template #default="scope">
                  <el-button type="primary" size="mini" @click="showContent(scope.row.content)">
                    查看
                  </el-button>
                </template>
            </el-table-column>
        </el-table>
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
    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
       <template #header>
              <div class="flex justify-between items-center">
                <span class="text-lg">{{type==='create'?'新增':'编辑'}}</span>
                <div>
                  <el-button :loading="btnLoading" type="primary" @click="enterDialog">确 定</el-button>
                  <el-button @click="closeDialog">取 消</el-button>
                </div>
              </div>
            </template>

          <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
            <el-form-item label="机器人ID:" prop="botID">
    <el-input v-model.number="formData.botID" :clearable="true" placeholder="请输入机器人ID" />
</el-form-item>
            <el-form-item label="发布次数:" prop="publishTimes">
    <el-input v-model.number="formData.publishTimes" :clearable="true" placeholder="请输入发布次数" />
</el-form-item>
            <el-form-item label="发布用户ID:" prop="userID">
    <el-input v-model.number="formData.userID" :clearable="true" placeholder="请输入发布用户ID" />
</el-form-item>
            <el-form-item label="发布价格:" prop="price">
    <el-input-number v-model="formData.price" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
            <el-form-item label="发布内容:" prop="content">
    <RichEdit v-model="formData.content"/>
</el-form-item>
          </el-form>
    </el-drawer>

    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="detailShow" :show-close="true" :before-close="closeDetailShow" title="查看">
            <el-descriptions :column="1" border>
                    <el-descriptions-item label="机器人ID">
                        {{ detailForm.botID }}
                    </el-descriptions-item>
                    <el-descriptions-item label="发布次数">
                      {{ detailForm.publishTimes }}
                  </el-descriptions-item>
                    <el-descriptions-item label="发布用户ID">
    {{ detailForm.userID }}
</el-descriptions-item>
                    <el-descriptions-item label="发布价格">
    {{ detailForm.price }}
</el-descriptions-item>
                    <el-descriptions-item label="发布内容">
    <RichView v-model="detailForm.content" />
</el-descriptions-item>
            </el-descriptions>
        </el-drawer>

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
import {
  createAdPublishRecord,
  deleteAdPublishRecord,
  deleteAdPublishRecordByIds,
  updateAdPublishRecord,
  findAdPublishRecord,
  getAdPublishRecordList
} from '@/api/recharge/adPublishRecord'
// 富文本组件
import RichEdit from '@/components/richtext/rich-edit.vue'
import RichView from '@/components/richtext/rich-view.vue'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict ,filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import { useAppStore } from "@/pinia"




defineOptions({
    name: 'AdPublishRecord'
})

const dialogVisible = ref(false)
const dialogContent = ref([])

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


// 提交按钮loading
const btnLoading = ref(false)
const appStore = useAppStore()

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
            botID: undefined,
            publishTimes: undefined,
            userID: undefined,
            price: 0,
            content: '',
        })



// 验证规则
const rule = reactive({
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    page.value = 1
    getTableData()
  })
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async() => {
  const table = await getAdPublishRecordList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () =>{
}

// 获取需要的字典 可能为空 按需保留
setOptions()


// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
    multipleSelection.value = val
}

// 删除行
const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
            deleteAdPublishRecordFunc(row)
        })
    }

// 多选删除
const onDelete = async() => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
      const IDs = []
      if (multipleSelection.value.length === 0) {
        ElMessage({
          type: 'warning',
          message: '请选择要删除的数据'
        })
        return
      }
      multipleSelection.value &&
        multipleSelection.value.map(item => {
          IDs.push(item.ID)
        })
      const res = await deleteAdPublishRecordByIds({ IDs })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        if (tableData.value.length === IDs.length && page.value > 1) {
          page.value--
        }
        getTableData()
      }
      })
    }

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updateAdPublishRecordFunc = async(row) => {
    const res = await findAdPublishRecord({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteAdPublishRecordFunc = async (row) => {
    const res = await deleteAdPublishRecord({ ID: row.ID })
    if (res.code === 0) {
        ElMessage({
                type: 'success',
                message: '删除成功'
            })
            if (tableData.value.length === 1 && page.value > 1) {
            page.value--
        }
        getTableData()
    }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)

// 打开弹窗
const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        botID: undefined,
        publishTimes: undefined,
        userID: undefined,
        price: 0,
        content: '',
        }
}
// 弹窗确定
const enterDialog = async () => {
     btnLoading.value = true
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return btnLoading.value = false
              let res
              switch (type.value) {
                case 'create':
                  res = await createAdPublishRecord(formData.value)
                  break
                case 'update':
                  res = await updateAdPublishRecord(formData.value)
                  break
                default:
                  res = await createAdPublishRecord(formData.value)
                  break
              }
              btnLoading.value = false
              if (res.code === 0) {
                ElMessage({
                  type: 'success',
                  message: '创建/更改成功'
                })
                closeDialog()
                getTableData()
              }
      })
}

const detailForm = ref({})

// 查看详情控制标记
const detailShow = ref(false)


// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}


// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findAdPublishRecord({ ID: row.ID })
  if (res.code === 0) {
    detailForm.value = res.data
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  detailForm.value = {}
}


</script>

<style>

</style>
