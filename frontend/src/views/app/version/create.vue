<template>
  <BaseDialog
    v-model="dialogState.isOpen"
    :title="`正在为《${dialogState.appInfo?.appName}》进行更新`"
    fullscreen
    @close="handleClose"
  >
    <el-scrollbar :native="true">
      <el-form ref="formRef" :model="dialogState.submitForm" class="ver-form" label-width="100px">
        <!-- 版本号（来自 appInfo，只读） -->
        <el-form-item label="版本号">
          <el-input v-model="dialogState.appInfo.version" placeholder="请输入应用版本" />
        </el-form-item>

        <!-- Sign -->
        <el-form-item
          :rules="[{ required: true, message: '请输入 Sign' }]"
          label="Sign"
          prop="sign"
        >
          <el-input v-model="dialogState.submitForm.sign" placeholder="请输入新版本 Sign" />
        </el-form-item>

        <!-- 强制更新 -->
        <el-form-item label="是否强制更新">
          <el-switch
            v-model="dialogState.submitForm.forceUpdate"
            active-text="是"
            inactive-text="否"
          />
        </el-form-item>

        <!-- 补丁地址 -->
        <el-form-item
          :rules="[{ required: true, message: '请输入补丁地址' }]"
          label="补丁地址"
          prop="patchUrl"
        >
          <el-input
            v-model="dialogState.submitForm.patchUrl"
            placeholder="https://example.com/patch.zip"
            type="url"
          />
        </el-form-item>

        <!-- 更新日志 -->
        <el-form-item
          :rules="[{ required: true, message: '请输入更新日志' }]"
          label="更新日志(Markdown)"
          prop="desc"
        >
          <el-input
            v-model="dialogState.submitForm.desc"
            :rows="6"
            placeholder="支持 Markdown 语法（提交后由后端渲染）"
            type="textarea"
          />
        </el-form-item>
      </el-form>
    </el-scrollbar>

    <!-- 自定义底部按钮 -->
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button :loading="loading" type="primary" @click="handleSubmit"> 确认更新</el-button>
      </div>
    </template>
  </BaseDialog>
</template>

<script lang="ts" setup>
import {nextTick, reactive, ref} from 'vue'
import type {FormInstance} from 'element-plus'
import {ElMessage} from 'element-plus'
import {infoApp} from '@/api/app.ts'
import type {IAppItem} from '@/types/app/app.ts'
import AppVerService from '@/api/version.ts'

defineOptions({ name: 'VerCreate' })

// 表单引用
const formRef = ref<FormInstance>()
const loading = ref(false)

// ✅ 优化后的状态结构（完整保留你需要的字段）
interface DialogState {
  isOpen: boolean
  appInfo?: IAppItem
  submitForm: {
    desc: string
    sign: string
    forceUpdate: boolean
    patchUrl: string
    version: string
  }
}

const dialogState = reactive<DialogState>({
  isOpen: false,
  appInfo: undefined,
  submitForm: {
    version: '',
    desc: '',
    sign: '',
    forceUpdate: false,
    patchUrl: '',
  },
})

// 显示对话框
const showDialog = async (appId: string) => {
  if (!appId) {
    ElMessage.error('请先选择一个应用才可进行更新')
    return
  }

  try {
    const resp = await infoApp(appId)
    if (resp.data.code !== 0) {
      ElMessage.error(resp.data.msg)
      return
    }
    dialogState.appInfo = resp.data.data.app
    dialogState.isOpen = true

    // 重置表单（可选）
    Object.assign(dialogState.submitForm, {
      desc: '',
      sign: '',
      forceUpdate: false,
      patchUrl: '',
    })

    // 重置表单验证状态
    nextTick(() => {
      formRef.value?.resetFields()
    })
  } catch (error) {
    ElMessage.error('获取应用信息失败')
  }
}

// 关闭对话框
const handleClose = () => {
  dialogState.isOpen = false
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()

    loading.value = true
    const resp = (
      await AppVerService.Update({
        appId: dialogState.appInfo?.appId,
        version: dialogState.appInfo?.version,
        ...dialogState.submitForm,
      })
    ).data
    if (resp.code !== 0) {
      ElMessage.error(resp.msg)
      return
    }
    ElMessage.success('更新成功！')
    handleClose()
  } catch (error) {
    console.warn('表单验证失败')
  } finally {
    loading.value = false
  }
}

defineExpose({
  showDialog,
})
</script>

<style scoped>
.ver-form {
  padding: 0 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 0 20px 20px;
}
</style>
