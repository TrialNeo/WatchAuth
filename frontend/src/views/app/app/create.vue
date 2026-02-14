<template>
  <BaseDialog
    v-model="open"
    :title="submitForm.appID == '' ? '编辑应用' : '新增应用'"
    style="height: 60vh"
    width="650"
    @close="close"
  >
    <el-scrollbar>
      <el-form
        ref="submitFormRef"
        :model="submitForm"
        :rules="formRules"
        label-position="right"
        label-width="100px"
      >
        <el-form-item label="应用名称" prop="appName">
          <el-input
            v-model="submitForm.appName"
            :readonly="!!submitForm.appID"
            placeholder="请输入应用名称"
          />
        </el-form-item>
        <el-form-item label="应用描述" prop="description">
          <el-input
            v-model="submitForm.description"
            :rows="3"
            placeholder="请输入应用描述"
            type="textarea"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="submitForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="计费方式">
          <el-select v-model="submitForm.feeType" placeholder="请选择">
            <el-option :value="0" label="免费" />
            <el-option :value="1" label="时长" />
            <el-option :value="2" label="次数" />
          </el-select>
        </el-form-item>
        <template v-if="submitForm.feeType === 1">
          <el-form-item label="价格">
            <el-input-number v-model="submitForm.fee" :min="0" :step="0.001" style="width: 100%">
              <template #suffix>
                <span>元/小时</span>
              </template>
            </el-input-number>
          </el-form-item>
        </template>
        <template v-if="submitForm.feeType === 2">
          <el-form-item label="单价">
            <el-input-number v-model="submitForm.fee" :min="0" :step="0.001" style="width: 100%">
              <template #suffix>
                <span>元/次</span>
              </template>
            </el-input-number>
          </el-form-item>
        </template>
        <el-form-item label="加密方式" prop="encType">
          <el-select v-model="submitForm.encType" placeholder="请选择">
            <el-option :value="0" label="None" />
            <el-option :value="1" label="RSA" />
            <el-option :value="2" label="AesGcm" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-scrollbar>

    <template #footer>
      <el-button @click="close">取消</el-button>
      <el-button :loading="submitLoading" type="primary" @click="confirm">确定</el-button>
    </template>
  </BaseDialog>
</template>

<script lang="ts" setup>
import {ElMessage, type ElTree, type FormInstance, type FormRules} from 'element-plus'
import type {IMenuItem} from '@/types/system/menu'
import {createApp2, infoApp} from '@/api/app.ts'
import type {ICreateOrUpdateAppParams} from '@/types/app/app.ts'

defineOptions({ name: 'AppCreate' })

const emits = defineEmits(['refresh'])

const submitFormRef = useTemplateRef<FormInstance>('submitFormRef')
const menuTreeRef = useTemplateRef<InstanceType<typeof ElTree> | null>('menuTreeRef')

// 对话框开关
const open = ref(false)

// 提交按钮加载状态
const submitLoading = ref(false)

// 菜单列表
const menuList = ref<IMenuItem[]>([])

// 表单数据
const submitForm = ref({
  appID: undefined as string | undefined,
  appName: '',
  description: '',
  encType: undefined as number | undefined,
  feeType: undefined as number | undefined,
  fee: 0,
  status: 0 | 1,
})

// 取消
const close = () => {
  open.value = false
  menuTreeRef.value?.setCheckedKeys([])
  submitFormRef.value?.resetFields()
  menuList.value = []
}

// 确定
const confirm = async () => {
  await submitFormRef.value?.validate()
  const payload: ICreateOrUpdateAppParams = {
    appid: submitForm.value.appID!,
    appName: submitForm.value.appName,
    description: submitForm.value.description,
    encType: submitForm.value.encType!,
    feeType: submitForm.value.feeType!,
    fee: submitForm.value.fee,
    status: submitForm.value.status,
  }
  const { data: res } = submitForm.value.appID
    ? await createApp2(payload)
    : await createApp2(payload)
  if (res.code !== 0) {
    ElMessage.error(res.msg)
    return
  }
  ElMessage.success(submitForm.value.appID ? '编辑成功' : '新增成功')
  emits('refresh', submitForm.value.appID ? 'update' : 'create')
  close()
}
// 获取应用信息
const getAppInfo = async (appID: string) => {
  const { data: res } = await infoApp(appID)
  switch (res.code) {
    case 500:
      return
    case 0:
      const { appName, description, status, feeType, fee, enctype } = res.data.app
      submitForm.value = {
        appID: appID,
        appName: appName,
        description: description,
        status: status,
        feeType: feeType,
        fee: fee,
        encType: enctype,
      }
      break
    default:
      ElMessage.error(res.msg)
      return
  }
}

// 表单验证规则
const formRules: FormRules = {
  appName: [{ required: true, message: '请输入应用名称', trigger: 'blur' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }],
}

// 显示对话框
const showDialog = async (appID: string | undefined) => {
  open.value = true
  await getAppInfo(appID as string)
}

defineExpose({
  showDialog,
})
</script>

<style></style>
