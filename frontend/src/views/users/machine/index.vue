<template>
  <div>
    <SearchForm
      v-model="queryForm"
      @reset="resetForm"
      @search="getMachineList"
    />

    <el-card class="card-mt-16" shadow="never">
      <el-table :border="TABLE_CONFIG.border" :data="machineList" show-overflow-tooltip>
        <el-table-column :align="TABLE_CONFIG.align" type="selection" width="55" />
        <el-table-column :align="TABLE_CONFIG.align" fixed label="序号" type="index" width="55" />
        <el-table-column
          :align="TABLE_CONFIG.align"
          fixed
          label="机器ID"
          min-width="100"
          prop="machineId"
        />
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="所属用户ID"
          min-width="100"
          prop="belong"
        />
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="机器名称"
          min-width="150"
        >
          <template #default="{ row }">
            {{ row.machine.machineName }}
          </template>
        </el-table-column>
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="平台"
          min-width="120"
        >
          <template #default="{ row }">
            {{ row.machine.platform }}
          </template>
        </el-table-column>
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="设备ID"
          min-width="180"
        >
          <template #default="{ row }">
            {{ row.machine.deviceId }}
          </template>
        </el-table-column>
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="CPU"
          min-width="150"
        >
          <template #default="{ row }">
            {{ row.machine.cpu }}
          </template>
        </el-table-column>
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="GPU"
          min-width="150"
        >
          <template #default="{ row }">
            {{ row.machine.gpu }}
          </template>
        </el-table-column>
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="内存"
          min-width="80"
        >
          <template #default="{ row }">
            {{ row.machine.ram }}
          </template>
        </el-table-column>
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="使用应用"
          min-width="200"
        >
          <template #default="{ row }">
            <el-tag v-for="app in row.usedApps" :key="app.appId" size="small" style="margin-right: 4px;">
              App {{ app.appId }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="在线状态"
          min-width="100"
        >
          <template #default="{ row }">
            <BaseTag
              :text="row.usedApps.some(app => app.online) ? '在线' : '离线'"
              :type="row.usedApps.some(app => app.online) ? 'success' : 'danger'"
            />
          </template>
        </el-table-column>
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="封禁状态"
          min-width="100"
        >
          <template #default="{ row }">
            <BaseTag
              :text="row.isBan ? '已封禁' : '正常'"
              :type="row.isBan ? 'danger' : 'success'"
            />
          </template>
        </el-table-column>
        <el-table-column :align="TABLE_CONFIG.align" fixed="right" label="操作" min-width="220">
          <template #default="{ row }">
            <el-button
              :type="row.isBan ? 'success' : 'danger'"
              size="small"
              @click="handleBanMachine(row)"
            >
              {{ row.isBan ? '解封' : '封禁' }}
            </el-button>
            <el-button size="small" type="warning" @click="handleOfflineMachine(row)">下线</el-button>
            <el-button size="small" type="primary" @click="handleReadLog(row)">读取日志</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 日志显示对话框 -->
    <LogDialog
      v-if="logDialogVisible"
      v-model:dialogVisible="logDialogVisible"
      :appLogs="appLogs"
      :machineId="currentMachineId"
      :networkLogs="networkLogs"
      :systemLogs="systemLogs"
    />

  </div>
</template>

<script lang="ts" setup>
import {onMounted, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {TABLE_CONFIG} from '@/config/elementConfig.ts'
import {MachineManager} from '@/api/userManager.ts'
import type {IMachineItem} from '@/types/userManager/user.ts'
import BaseTag from '@/components/tag/BaseTag.vue'
import LogDialog from '@/components/machine/LogDialog.vue'
import SearchForm from '@/components/machine/SearchForm.vue'

defineOptions({ name: 'MachineManager' })

// 查询表单
const queryForm = ref({
  deviceId: '',
  belong: '',
})

const machineList = ref<IMachineItem[]>([])

// 日志对话框相关
const logDialogVisible = ref(false)
const currentMachineId = ref(0)
const activeLogTab = ref('system')
const systemLogs = ref('')
const appLogs = ref('')
const networkLogs = ref('')

const getMachineList = async () => {
  try {
    const resp = await MachineManager.list(queryForm.value)
    if (resp.data.code !== 0) {
      ElMessage.error(resp.data.msg)
      return
    }
    machineList.value = resp.data.data || []
  } catch (error) {
    ElMessage.error('获取机器列表失败')
    console.error('获取机器列表失败:', error)
  }
}

const resetForm = () => {
  queryForm.value = {
    deviceId: '',
    belong: '',
  }
  getMachineList()
}

const handleBanMachine = async (row: IMachineItem) => {
  try {
    const resp = await MachineManager.ban(row.machineId)
    if (resp.data.code !== 0) {
      ElMessage.error(resp.data.msg)
      return
    }
    ElMessage.success('机器封禁成功')
    getMachineList()
  } catch (error) {
    ElMessage.error('封禁机器失败')
    console.error('封禁机器失败:', error)
  }
}

const handleOfflineMachine = async (row: IMachineItem) => {
  try {
    const resp = await MachineManager.offline(row.machineId)
    if (resp.data.code !== 0) {
      ElMessage.error(resp.data.msg)
      return
    }
    ElMessage.success('机器下线成功')
    getMachineList()
  } catch (error) {
    ElMessage.error('机器下线失败')
    console.error('机器下线失败:', error)
  }
}

const handleReadLog = async (row: IMachineItem) => {
  try {
    const resp = await MachineManager.readLog(row.machineId)
    if (resp.data.code !== 0) {
      ElMessage.error(resp.data.msg)
      return
    }
    ElMessage.success('读取日志成功')

    // 设置当前机器ID
    currentMachineId.value = row.machineId

    // 处理日志数据，显示在对话框中
    console.log('日志数据:', resp.data.data)

    // 模拟日志数据，实际项目中应该根据 API 返回的数据来设置
    systemLogs.value = `[2026-02-27 10:00:00] 系统启动\n[2026-02-27 10:01:00] 加载配置文件\n[2026-02-27 10:02:00] 初始化服务\n[2026-02-27 10:03:00] 系统正常运行`
    appLogs.value = `[2026-02-27 10:05:00] 应用启动\n[2026-02-27 10:06:00] 加载插件\n[2026-02-27 10:07:00] 应用正常运行\n[2026-02-27 10:08:00] 处理用户请求`
    networkLogs.value = `[2026-02-27 10:00:00] 网络连接建立\n[2026-02-27 10:01:00] 连接服务器成功\n[2026-02-27 10:02:00] 心跳包发送成功\n[2026-02-27 10:03:00] 数据同步完成`

    // 显示日志对话框
    logDialogVisible.value = true
  } catch (error) {
    ElMessage.error('读取日志失败')
    console.error('读取日志失败:', error)
  }
}



onMounted(() => {
  getMachineList()
})
</script>

<style lang="scss" scoped>
</style>
