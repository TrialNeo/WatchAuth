<template>
  <div>
    <el-card shadow="never">
      <template #header>
        <span>系统配置</span>
      </template>
      <el-form :model="configMap" label-width="160px">
        <el-form-item label="站点名称">
          <el-input v-model="configMap.site_name" placeholder="WatchAuth" />
        </el-form-item>
        <el-form-item label="站点 Logo">
          <el-input v-model="configMap.site_logo" placeholder="Logo URL" />
        </el-form-item>
        <el-form-item label="开放注册">
          <el-switch v-model="allowRegister" active-text="开启" inactive-text="关闭" />
        </el-form-item>
        <el-form-item label="License 默认有效期(天)">
          <el-input-number v-model="licenseDays" :min="1" :max="36500" />
        </el-form-item>
        <el-form-item label="SDK 心跳间隔(秒)">
          <el-input-number v-model="heartbeatInterval" :min="30" :max="3600" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saving" @click="handleSave">保存配置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import { getSystemConfigs, updateSystemConfig } from '@/api/systemConfig'
import { ElMessage } from 'element-plus'

defineOptions({ name: 'SystemConfigView' })

const configMap = reactive<Record<string, string>>({})
const saving = ref(false)

const allowRegister = computed({
  get: () => configMap.allow_register === 'true',
  set: (v: boolean) => { configMap.allow_register = v ? 'true' : 'false' },
})

const licenseDays = computed({
  get: () => Number(configMap.default_license_days ?? '365'),
  set: (v: number) => { configMap.default_license_days = String(v) },
})

const heartbeatInterval = computed({
  get: () => Number(configMap.license_check_interval ?? '300'),
  set: (v: number) => { configMap.license_check_interval = String(v) },
})

const fetchData = async () => {
  const { data: res } = await getSystemConfigs()
  if (res.code !== 0) {
    ElMessage.error(res.msg || '获取配置失败')
    return
  }
  for (const item of (res.data?.configs ?? [])) {
    configMap[item.key] = item.value
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    for (const [key, value] of Object.entries(configMap)) {
      const { data: res } = await updateSystemConfig(key, value)
      if (res.code !== 0) {
        ElMessage.error(`${key}: ${res.msg || '保存失败'}`)
        return
      }
    }
    ElMessage.success('配置已保存')
  } finally {
    saving.value = false
  }
}

onMounted(fetchData)
</script>
