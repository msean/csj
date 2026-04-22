<template>
  <div class="button-group-editor">

    <div class="group-row" v-for="(row, rowIndex) in groups" :key="rowIndex">
      <!-- 每行按钮展示 -->
      <el-button
        v-for="(btn, btnIndex) in row"
        :key="btnIndex"
        @click="editButton(rowIndex, btnIndex)"
        size="small"
      >
        {{ btn.name }}
      </el-button>

      <!-- 删除行 -->
      <el-button
        type="danger"
        size="small"
        @click="removeRow(rowIndex)"
      >删除该行</el-button>

      <!-- 在本行添加按钮 -->
      <el-button
        type="primary"
        size="small"
        @click="addButton(rowIndex)"
      >新增按钮</el-button>
    </div>

    <!-- 新增新的行 -->
    <el-button
      type="success"
      size="small"
      @click="addRow"
      style="margin-top: 10px;"
    >
      新增一行
    </el-button>

    <!-- 编辑弹窗 -->
    <el-dialog v-model="dialogVisible" title="编辑按钮" width="400px">
      <el-form :model="editForm" label-width="90px">
        <el-form-item label="按钮名称">
          <el-input v-model="editForm.name" />
        </el-form-item>
        <el-form-item label="跳转链接">
          <el-input v-model="editForm.url" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveEdit">保存</el-button>
      </template>
    </el-dialog>

  </div>
</template>

<script setup>
import { ref, watch } from "vue";

// 接收父组件的 v-model
const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  }
});
const emit = defineEmits(["update:modelValue"]);

// 转成二维数组的结构
const groups = ref(props.modelValue.length ? props.modelValue : [[]]);

watch(
  () => props.modelValue,
  (val) => {
    if (val) groups.value = val;
  }
);

// 弹窗相关
const dialogVisible = ref(false);
const editForm = ref({ name: "", url: "" });
let currentRow = null;
let currentIndex = null;

// 新增一行
function addRow() {
  groups.value.push([]);
  emit("update:modelValue", groups.value);
}

// 删除一行
function removeRow(rowIndex) {
  groups.value.splice(rowIndex, 1);
  emit("update:modelValue", groups.value);
}

// 在本行新增按钮
function addButton(rowIndex) {
  currentRow = rowIndex;
  currentIndex = null;
  editForm.value = { name: "", url: "" };
  dialogVisible.value = true;
}

// 编辑按钮
function editButton(rowIndex, btnIndex) {
  currentRow = rowIndex;
  currentIndex = btnIndex;
  editForm.value = { ...groups.value[rowIndex][btnIndex] };
  dialogVisible.value = true;
}

// 保存按钮
function saveEdit() {
  if (!editForm.value.name) return;

  if (currentIndex === null) {
    // 新增
    groups.value[currentRow].push({ ...editForm.value });
  } else {
    // 修改
    groups.value[currentRow][currentIndex] = { ...editForm.value };
  }

  emit("update:modelValue", groups.value);
  dialogVisible.value = false;
}
</script>

<style scoped>
.group-row {
  margin-bottom: 10px;
}
</style>
