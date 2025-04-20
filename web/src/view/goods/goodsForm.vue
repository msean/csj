<template>
    <div>
      <div class="gva-form-box">
        <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
          <el-form-item label="uid" prop="uid">
            <el-input v-model="formData.uid" :clearable="true"  placeholder="请输入uid字段" />
         </el-form-item>
          <el-form-item label="所属用户" prop="ownerUser">
            <el-input v-model="formData.ownerUser" :clearable="true"  placeholder="请输入所属用户" />
         </el-form-item>
          <el-form-item label="菜品名称" prop="name">
            <el-input v-model="formData.name" :clearable="true"  placeholder="请输入客户名字" />
         </el-form-item>
          <el-form-item label="所属分类" prop="category">
            <el-input v-model="formData.category" :clearable="true"  placeholder="请输入手机号" />
         </el-form-item>
          <el-form-item label="'类别" prop="type">
            <el-input v-model="formData.type" :clearable="true"  placeholder="请输入类别" />
         </el-form-item>
          <el-form-item label="价格" prop="price">
            <el-input-number v-model="formData.price" :precision="2" :clearable="true"></el-input-number>
         </el-form-item>
          <el-form-item label="重量" prop="weight">
            <el-input v-model="formData.weight" :clearable="true"  placeholder="请输入重量'" />
         </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-input v-model="formData.status" :clearable="true"  placeholder="请输入状态" />
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
    createCustomers,
    updateCustomers,
    findCustomers
  } from '@/api/csj_customers/customers'
  
  defineOptions({
      name: 'CustomersForm'
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
              ownerUser: '',
              name: '',
              phone: '',
              remark: '',
              debt: 0,
              addr: '',
              carNo: '',
          })
  // 验证规则
  const rule = reactive({
  })
  
  const elFormRef = ref()
  
  // 初始化方法
  const init = async () => {
   // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
      if (route.query.id) {
        const res = await findCustomers({ ID: route.query.id })
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
                 res = await createCustomers(formData.value)
                 break
               case 'update':
                 res = await updateCustomers(formData.value)
                 break
               default:
                 res = await createCustomers(formData.value)
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
  