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
        <el-form-item label="手机号码" prop="phone">
          <el-input v-model="searchInfo.phone" placeholder="手机号码" clearable />
        </el-form-item>
        <el-form-item label="用户名" prop="name">
          <el-input v-model="searchInfo.name" placeholder="用户名" clearable />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="uid">
        <el-table-column align="left" label="UID" prop="uid" width="220" show-overflow-tooltip />
        <el-table-column align="left" label="用户名" prop="name" width="140" show-overflow-tooltip />
        <el-table-column align="left" label="手机号码" prop="phone" width="140" />
        <el-table-column align="left" label="体验截止时间" prop="experienceDeadline" width="180">
          <template #default="scope">{{ formatDate(scope.row.experienceDeadline) }}</template>
        </el-table-column>
        <el-table-column align="left" label="会员截止时间" prop="vipDeadline" width="180">
          <template #default="scope">
            <el-tag :type="isVipActive(scope.row.vipDeadline) ? 'success' : 'info'">
              {{ formatDate(scope.row.vipDeadline) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="注册时间" prop="createdAt" width="180">
          <template #default="scope">{{ formatDate(scope.row.createdAt) }}</template>
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
  </div>
</template>

<script setup>
  import { getAppUserList } from '@/api/appUser'
  import { formatDate } from '@/utils/format'
  import { ref } from 'vue'

  defineOptions({ name: 'AppUser' })

  const isVipActive = (vipDeadline) => {
    if (!vipDeadline) return false
    return new Date(vipDeadline) > new Date()
  }

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
    const table = await getAppUserList({
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
</script>

<style></style>
