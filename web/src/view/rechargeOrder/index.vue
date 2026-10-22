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
        <el-form-item label="用户UUID" prop="userUUID">
          <el-input v-model="searchInfo.userUUID" placeholder="用户UUID" clearable />
        </el-form-item>
        <el-form-item label="订单号" prop="orderNo">
          <el-input v-model="searchInfo.orderNo" placeholder="订单号" clearable />
        </el-form-item>
        <el-form-item label="支付状态" prop="payStatus">
          <el-select v-model="searchInfo.payStatus" placeholder="全部" clearable style="width: 140px">
            <el-option
              v-for="item in payStatusOptions"
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
      <el-table
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="uid"
      >
        <el-table-column align="left" label="订单号" prop="orderNo" width="220" show-overflow-tooltip />
        <el-table-column align="left" label="用户UUID" prop="userUUID" width="200" show-overflow-tooltip />
        <el-table-column align="left" label="会员类型" prop="memberType" width="120">
          <template #default="scope">
            <el-tag>{{ memberTypeLabel(scope.row.memberType) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="充值金额" prop="amount" width="120" />
        <el-table-column align="left" label="支付状态" prop="payStatus" width="110">
          <template #default="scope">
            <el-tag :type="payStatusTagType(scope.row.payStatus)">
              {{ payStatusLabel(scope.row.payStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="微信交易号" prop="transactionID" width="220" show-overflow-tooltip />
        <el-table-column align="left" label="支付时间" prop="payTime" width="180">
          <template #default="scope">{{ scope.row.payTime ? formatDate(scope.row.payTime) : '-' }}</template>
        </el-table-column>
        <el-table-column align="left" label="创建时间" prop="createdAt" width="180">
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
  import { findRechargeOrder, getRechargeOrderList } from '@/api/rechargeOrder'
  import { formatDate } from '@/utils/format'
  import { ref } from 'vue'

  defineOptions({ name: 'RechargeOrder' })

  const memberTypeOptions = [
    { label: '月', value: 1 },
    { label: '季度', value: 2 },
    { label: '年', value: 3 }
  ]

  const memberTypeLabel = (type) => {
    const item = memberTypeOptions.find((opt) => opt.value === type)
    return item ? item.label : type
  }

  const payStatusOptions = [
    { label: '待支付', value: 0 },
    { label: '已支付', value: 1 },
    { label: '已取消', value: 2 },
    { label: '已退款', value: 3 }
  ]

  const payStatusLabel = (status) => {
    const item = payStatusOptions.find((opt) => opt.value === status)
    return item ? item.label : status
  }

  const payStatusTagType = (status) => {
    switch (status) {
      case 1:
        return 'success'
      case 2:
        return 'info'
      case 3:
        return 'warning'
      default:
        return 'danger'
    }
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
    const table = await getRechargeOrderList({
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
