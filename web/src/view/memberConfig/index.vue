<template>
  <div>
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        class="demo-form-inline"
        @keyup.enter="onSubmit"
      >
        <el-form-item label="类型" prop="type">
          <el-select v-model="searchInfo.type" placeholder="全部" clearable style="width: 160px">
            <el-option
              v-for="item in memberTypeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
        <el-button
          icon="delete"
          style="margin-left: 10px"
          :disabled="!multipleSelection.length"
          @click="onDelete"
          >删除</el-button
        >
      </div>
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="uid"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="类型" prop="type" width="160">
          <template #default="scope">
            <el-tag>{{ memberTypeLabel(scope.row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="收费金额" prop="price" width="180" />
        <el-table-column align="left" label="创建时间" prop="createdAt" width="180">
          <template #default="scope">{{ formatDate(scope.row.createdAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="200">
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="edit"
              class="table-button"
              @click="updateRow(scope.row)"
              >变更</el-button
            >
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
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

    <el-drawer
      destroy-on-close
      size="600"
      v-model="dialogFormVisible"
      :show-close="false"
      :before-close="closeDialog"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '添加会员收费配置' : '修改会员收费配置' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>

      <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="类型:" prop="type">
          <el-select v-model="formData.type" placeholder="请选择类型" style="width: 100%">
            <el-option
              v-for="item in memberTypeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="收费金额:" prop="price">
          <el-input-number v-model="formData.price" :min="0" :precision="2" :step="1" style="width: 100%" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    createMemberConfig,
    deleteMemberConfig,
    deleteMemberConfigByIds,
    updateMemberConfig,
    findMemberConfig,
    getMemberConfigList
  } from '@/api/memberConfig'
  import { formatDate } from '@/utils/format'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { ref, reactive } from 'vue'

  defineOptions({ name: 'MemberConfig' })

  const memberTypeOptions = [
    { label: '月', value: 1 },
    { label: '季度', value: 2 },
    { label: '年', value: 3 }
  ]

  const memberTypeLabel = (type) => {
    const item = memberTypeOptions.find((opt) => opt.value === type)
    return item ? item.label : type
  }

  const formData = ref({
    type: '',
    price: 0
  })

  const rule = reactive({
    type: [{ required: true, message: '请选择类型', trigger: ['change', 'blur'] }],
    price: [{ required: true, message: '请输入收费金额', trigger: ['input', 'blur'] }]
  })

  const elFormRef = ref()
  const elSearchFormRef = ref()

  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const tableData = ref([])
  const searchInfo = ref({})

  const onReset = () => {
    searchInfo.value = {}
    getTableData()
  }

  const onSubmit = () => {
    page.value = 1
    getTableData()
  }

  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }

  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  const getTableData = async () => {
    const table = await getMemberConfigList({
      page: page.value,
      pageSize: pageSize.value,
      ...searchInfo.value
    })
    if (table.code === 0) {
      tableData.value = table.data.list
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
    }
  }

  getTableData()

  const multipleSelection = ref([])
  const handleSelectionChange = (val) => {
    multipleSelection.value = val
  }

  const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      deleteMemberConfigFunc(row)
    })
  }

  const onDelete = async () => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      const IDs = []
      if (multipleSelection.value.length === 0) {
        ElMessage({ type: 'warning', message: '请选择要删除的数据' })
        return
      }
      multipleSelection.value.map((item) => IDs.push(item.uid))
      const res = await deleteMemberConfigByIds({ uids: IDs })
      if (res.code === 0) {
        ElMessage({ type: 'success', message: '删除成功' })
        if (tableData.value.length === IDs.length && page.value > 1) page.value--
        getTableData()
      }
    })
  }

  const type = ref('')

  const updateRow = async (row) => {
    const res = await findMemberConfig({ uid: row.uid })
    type.value = 'update'
    if (res.code === 0) {
      formData.value = res.data
      dialogFormVisible.value = true
    }
  }

  const deleteMemberConfigFunc = async (row) => {
    const res = await deleteMemberConfig({ uid: row.uid })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '删除成功' })
      if (tableData.value.length === 1 && page.value > 1) page.value--
      getTableData()
    }
  }

  const dialogFormVisible = ref(false)

  const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
  }

  const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = { type: '', price: 0 }
  }

  const enterDialog = async () => {
    elFormRef.value?.validate(async (valid) => {
      if (!valid) return
      let res
      if (type.value === 'update') {
        res = await updateMemberConfig(formData.value)
      } else {
        res = await createMemberConfig(formData.value)
      }
      if (res.code === 0) {
        ElMessage({ type: 'success', message: '操作成功' })
        closeDialog()
        getTableData()
      }
    })
  }
</script>

<style></style>
