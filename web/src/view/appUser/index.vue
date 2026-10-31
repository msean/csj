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
        <el-table-column label="市场名称" prop="marketName" min-width="140" show-overflow-tooltip />
 <el-table-column label="档口名称" prop="stallName" min-width="140" show-overflow-tooltip />
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
        <el-table-column label="操作" fixed="right" width="90"><template #default="{ row }"><el-button type="primary" link @click="openDetail(row)">详情</el-button></template></el-table-column>
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
    <el-drawer v-model="detailVisible" title="记录详情" size="min(850px, 95vw)" destroy-on-close>
      <div v-loading="detailLoading"><RecordDetail :record="detail" /><div class="mt-5"><el-button v-for="(config, resource) in resources" :key="resource" @click="router.push({ name: config.name, query: { ownerUser: detail.uid } }); detailVisible = false">{{ config.title }}</el-button></div></div></el-drawer>
  </div>
</template>

<script setup>
  import RecordDetail from '@/components/appData/RecordDetail.vue'

  import { findAppUser, getAppUserList } from '@/api/appUser'
  import { formatDate } from '@/utils/format'
  import { ref } from 'vue'
 import { useRouter } from 'vue-router'
 import { resources } from '@/view/appData/config'
 const router = useRouter()

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
 page.value = 1
    getTableData()
  }

  const onSubmit = () => {
    page.value = 1
    getTableData()
  }

  const handleSizeChange = (val) => {
    pageSize.value = val
 page.value = 1
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


  const detailVisible = ref(false), detailLoading = ref(false), detail = ref({})
  let detailSequence = 0
  const openDetail = async (row) => {
    const current = ++detailSequence
    detail.value = {}; detailVisible.value = true; detailLoading.value = true
    try { const res = await findAppUser({ uid: row.uid }); if (current === detailSequence && res.code === 0) detail.value = res.data }
    finally { if (current === detailSequence) detailLoading.value = false }
  }

  getTableData()
</script>

<style></style>
