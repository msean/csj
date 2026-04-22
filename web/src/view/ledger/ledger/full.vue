<template>
  <div class="ledger-full">
    <h2>入款 ({{ income.count }})</h2>
    <el-table :data="income.list" stripe style="margin-bottom: 20px">
      <el-table-column prop="time" label="时间" width="120" />
      <el-table-column prop="amount" label="金额" width="120" />
      <el-table-column prop="remark" label="备注" />
      <el-table-column prop="replyUser" label="回复人" />
      <el-table-column prop="operator" label="操作人" />
      <el-table-column prop="currentFeeRate" label="费率" />
      <el-table-column prop="afterNote" label="后备注" />
    </el-table>

    <h2>下发 ({{ payout.count }})</h2>
    <el-table :data="payout.list" stripe style="margin-bottom: 20px">
      <el-table-column prop="time" label="时间" width="120" />
      <el-table-column prop="amount" label="金额" width="120" />
      <el-table-column prop="remark" label="备注" />
      <el-table-column prop="replyUser" label="回复人" />
      <el-table-column prop="operator" label="操作人" />
      <el-table-column prop="afterNote" label="后备注" />
    </el-table>

    <h3>总计</h3>
    <p>总入款: {{ summary.totalIncome }}</p>
    <p>应下发: {{ summary.shouldPaid }}</p>
    <p>总下发: {{ summary.totalPayout }}</p>
    <p>未下发: {{ summary.unpaid }}</p>
  </div>
</template>


<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getLedgerFull } from '@/api/usage/ledger'

const income = ref({ list: [], count: 0 })
const payout = ref({ list: [], count: 0 })
const summary = ref({ totalIncome: 0, totalPayout: 0, unpaid: 0 })

const route = useRoute()  // 获取 URL query

onMounted(async () => {
  try {
    // 从 URL query 里拿参数
    const botID = Number(route.query.bot_id)
    const chatGroupID = Number(route.query.chat_group_id)
    const idMin = Number(route.query.idmin)
    const idMax = Number(route.query.idmax)

    // 调接口
    const res = await getLedgerFull({
      bot_id: botID,
      chat_group_id: chatGroupID,
      idmin: idMin,
      idmax: idMax
    })

    // 直接用接口返回的数据
    // 你的接口返回的是 { income: {...}, payout: {...}, summary: {...} }
    income.value = res.data.income || { list: [], count: 0 }
    payout.value = res.data.payout || { list: [], count: 0 }
    summary.value = res.data.summary || { totalIncome: 0, totalPayout: 0, unpaid: 0 }

  } catch (err) {
    console.error('获取数据失败', err)
  }
})
</script>
