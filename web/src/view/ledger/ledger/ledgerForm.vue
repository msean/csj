
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="操作用户ID:" prop="oprUserID">
    <el-input v-model.number="formData.oprUserID" :clearable="true" placeholder="请输入操作用户ID" />
</el-form-item>
        <el-form-item label="操作人的用户名称:" prop="oprUsername">
    <el-input v-model="formData.oprUsername" :clearable="true" placeholder="请输入操作人的用户名称" />
</el-form-item>
        <el-form-item label="操作人昵称:" prop="oprUserNickname">
    <el-input v-model="formData.oprUserNickname" :clearable="true" placeholder="请输入操作人昵称" />
</el-form-item>
        <el-form-item label="操作类型:" prop="actionType">
    <el-input v-model.number="formData.actionType" :clearable="true" placeholder="请输入操作类型" />
</el-form-item>
        <el-form-item label="操作金额:" prop="amount">
    <el-input-number v-model="formData.amount" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="所在群组:" prop="chatGroupID">
    <el-input v-model.number="formData.chatGroupID" :clearable="true" placeholder="请输入所在群组" />
</el-form-item>
        <el-form-item label="消息ID:" prop="messageID">
    <el-input v-model.number="formData.messageID" :clearable="true" placeholder="请输入消息ID" />
</el-form-item>
        <el-form-item label="原始输入:" prop="rawInput">
    <el-input v-model="formData.rawInput" :clearable="true" placeholder="请输入原始输入" />
</el-form-item>
        <el-form-item>
          <el-button :loading="btnLoading" type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {
  createLedger,
  updateLedger,
  findLedger
} from '@/api/usage/ledger'

defineOptions({
    name: 'LedgerForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'


const route = useRoute()
const router = useRouter()

// 提交按钮loading
const btnLoading = ref(false)

const type = ref('')
const formData = ref({
            oprUserID: undefined,
            oprUsername: '',
            oprUserNickname: '',
            actionType: undefined,
            amount: 0,
            chatGroupID: undefined,
            messageID: undefined,
            rawInput: '',
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findLedger({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
}

init()
// 保存按钮
const save = async() => {
      btnLoading.value = true
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return btnLoading.value = false
            let res
           switch (type.value) {
             case 'create':
               res = await createLedger(formData.value)
               break
             case 'update':
               res = await updateLedger(formData.value)
               break
             default:
               res = await createLedger(formData.value)
               break
           }
           btnLoading.value = false
           if (res.code === 0) {
             ElMessage({
               type: 'success',
               message: '创建/更改成功'
             })
           }
       })
}

// 返回按钮
const back = () => {
    router.go(-1)
}

</script>

<style>
</style>
