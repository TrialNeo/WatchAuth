<template>
  <el-card class="card-clear-mb" shadow="never">
    <el-form
      ref="formRef"
      :model="formData"
      label-width="auto"
      @keyup.enter="handleSearch"
    >
      <el-row :gutter="10">
        <el-col :lg="6" :md="12" :sm="12" :xl="6" :xs="24">
                  <el-form-item label="设备ID" prop="deviceId">
                    <el-input v-model="formData.deviceId" placeholder="请输入设备ID" />
                  </el-form-item>
                </el-col>
                <el-col :lg="6" :md="12" :sm="12" :xl="6" :xs="24">
                  <el-form-item label="所属用户ID" prop="belong">
                    <el-input v-model="formData.belong" placeholder="请输入所属用户ID" type="number" />
                  </el-form-item>
                </el-col>
        <el-col :lg="6" :md="12" :sm="12" :xl="6" :xs="24">
          <el-form-item>
            <el-button
              type="primary"
              @click="handleSearch"
              >搜索
            </el-button>
            <el-button @click="handleReset">重置</el-button>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </el-card>
</template>

<script lang="ts" setup>
import {defineEmits, defineProps, ref} from 'vue'
import type {FormInstance} from 'element-plus'

const props = defineProps<{
  modelValue: {
    deviceId: string
    belong: string
  }
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: { deviceId: string; belong: string }): void
  (e: 'search'): void
  (e: 'reset'): void
}>()

const formRef = ref<FormInstance>()
const formData = ref({ ...props.modelValue })

const handleSearch = () => {
  emit('search')
}

const handleReset = () => {
  formData.value = {
    deviceId: '',
    belong: '',
  }
  emit('update:modelValue', formData.value)
  emit('reset')
}
</script>

<style lang="scss" scoped>
</style>
