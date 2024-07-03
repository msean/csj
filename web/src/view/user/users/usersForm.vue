<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="uid字段:" prop="uid">
          <el-input v-model="formData.uid" :clearable="true"  placeholder="请输入uid字段" />
       </el-form-item>
        <el-form-item label="用户名:" prop="name">
          <el-input v-model="formData.name" :clearable="true"  placeholder="请输入用户名" />
       </el-form-item>
        <el-form-item label="手机号码:" prop="phone">
          <el-input v-model="formData.phone" :clearable="true"  placeholder="请输入手机号码" />
       </el-form-item>
        <el-form-item label="体验截止时间:" prop="vipExpireTime">
          <el-date-picker v-model="formData.vipExpireTime" type="date" placeholder="选择日期" :clearable="true"></el-date-picker>
       </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {
  createUsers,
  updateUsers,
  findUsers
} from '@/api/user/users'

defineOptions({
    name: 'UsersForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'

const route = useRoute()
const router = useRouter()

const type = ref('')
const formData = ref({
            uid: '',
            name: '',
            phone: '',
            vipExpireTime: new Date(),
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findUsers({ ID: route.query.id })
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
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return
            let res
           switch (type.value) {
             case 'create':
               res = await createUsers(formData.value)
               break
             case 'update':
               res = await updateUsers(formData.value)
               break
             default:
               res = await createUsers(formData.value)
               break
           }
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
