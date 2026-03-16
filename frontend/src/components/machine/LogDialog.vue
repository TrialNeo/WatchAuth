<template>
  <el-dialog
    v-model="localDialogVisible"
    :close-on-click-modal="false"
    :title="`机器 ${machineName} 日志`"
    width="80%"
  >
    <div class="log-container">
      <div class="log-header">
        <div class="log-summary">
          共 {{ logs.length }} 条日志
        </div>
      </div>
      <el-scrollbar height="400px">
        <div v-if="logs.length === 0" class="empty-log">
          暂无日志数据
        </div>
        <div v-else class="log-list">
          <div
            v-for="(log, index) in logs"
            :key="index"
            :class="`log-type-${log.type}`"
            class="log-item"
          >
            <div class="log-time">{{ formatTime(log.time) }}</div>
            <div class="log-content">
              <div class="log-info">
                <span :class="`badge-${log.type}`" class="log-type-badge">
                  {{ log.type.toUpperCase() }}
                </span>
                <span class="log-app">App: {{ log.appId }}</span>
                <span class="log-version">{{ log.appVersionId }}</span>
              </div>
              <div class="log-details">
                <span class="log-module">{{ log.module }}.{{ log.funcName }}</span>
                <span class="log-train">Train: {{ log.trainId }} Span: {{ log.spanId }}</span>
              </div>
              <div class="log-message">
                {{ log.msg }}
              </div>
            </div>
          </div>
        </div>
      </el-scrollbar>
    </div>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">关闭</el-button>
        <el-button type="primary" @click="downloadLogs">下载日志</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import {ref, watch} from 'vue'
import {ElMessage} from 'element-plus'
import type {ILogItem} from '@/types/userManager/user.ts'

const props = defineProps<{
  dialogVisible: boolean
  machineId: number
  machineName: string
  logs: ILogItem[]
}>()

const emit = defineEmits<{
  (e: 'update:dialogVisible', value: boolean): void
}>()

const localDialogVisible = ref(props.dialogVisible)
const machineId = ref(props.machineId)
const machineName = ref(props.machineName)
const logs = ref<ILogItem[]>(props.logs)

// 监听 dialogVisible 变化
watch(() => props.dialogVisible, (newValue) => {
  localDialogVisible.value = newValue
})

// 监听本地 dialogVisible 变化，通知父组件
watch(localDialogVisible, (newValue) => {
  emit('update:dialogVisible', newValue)
})

// 监听 machineId 变化
watch(() => props.machineId, (newValue) => {
  machineId.value = newValue
})

// 监听 machineName 变化
watch(() => props.machineName, (newValue) => {
  machineName.value = newValue
})

// 监听日志内容变化
watch(() => props.logs, (newValue) => {
  logs.value = newValue
}, { deep: true })

const handleClose = () => {
  localDialogVisible.value = false
}

const formatTime = (timeStr: string) => {
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

const downloadLogs = () => {
  try {
    // 生成日志内容
    let logContent = `机器 ${machineId.value} 日志\n\n`

    logs.value.forEach((log, index) => {
      logContent += `[${formatTime(log.time)}] [${log.type.toUpperCase()}] App: ${log.appId} (${log.appVersionId})\n`
      logContent += `  Module: ${log.module}.${log.funcName}\n`
      logContent += `  Message: ${log.msg}\n\n`
    })

    // 创建下载链接
    const blob = new Blob([logContent], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `machine_${machineId.value}_logs_${new Date().toISOString().split('T')[0]}.txt`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)

    ElMessage.success('日志下载成功')
  } catch (error) {
    console.error('下载日志失败:', error)
  }
}
</script>

<style lang="scss" scoped>
.log-container {
  width: 100%;
}

.log-header {
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e4e7ed;
  .log-summary {
    font-size: 14px;
    color: #606266;
  }
}

.empty-log {
  text-align: center;
  padding: 40px 0;
  color: #999;
}

.log-list {
  padding: 10px 0;
}

.log-item {
  margin-bottom: 12px;
  padding: 15px;
  border-radius: 6px;
  background-color: #f5f5f5;
  border-left: 4px solid #ccc;
  transition: all 0.3s ease;

  &:hover {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  &.log-type-info {
    border-left-color: #409eff;
    background-color: #ecf5ff;
  }

  &.log-type-debug {
    border-left-color: #67c23a;
    background-color: #f0f9eb;
  }

  &.log-type-warning {
    border-left-color: #e6a23c;
    background-color: #fdf6ec;
  }

  &.log-type-error {
    border-left-color: #f56c6c;
    background-color: #fef0f0;
  }
}

.log-time {
  font-size: 12px;
  color: #999;
  margin-bottom: 8px;
}

.log-content {
  font-size: 14px;
  line-height: 1.5;

  .log-info {
    margin-bottom: 6px;
    display: flex;
    align-items: center;
    flex-wrap: wrap;
  }

  .log-type-badge {
    padding: 2px 8px;
    border-radius: 10px;
    font-size: 12px;
    font-weight: bold;
    margin-right: 10px;

    &.badge-info {
      background-color: #ecf5ff;
      color: #409eff;
      border: 1px solid #d9ecff;
    }

    &.badge-debug {
      background-color: #f0f9eb;
      color: #67c23a;
      border: 1px solid #e1f3d8;
    }

    &.badge-warning {
      background-color: #fdf6ec;
      color: #e6a23c;
      border: 1px solid #faecd8;
    }

    &.badge-error {
      background-color: #fef0f0;
      color: #f56c6c;
      border: 1px solid #fbc4c4;
    }
  }

  .log-app {
    margin-right: 10px;
    color: #666;
  }

  .log-version {
    margin-right: 15px;
    font-size: 12px;
    color: #999;
  }

  .log-details {
    margin-bottom: 6px;
    display: flex;
    align-items: center;
    flex-wrap: wrap;
  }

  .log-module {
    margin-right: 15px;
    font-family: 'Courier New', Courier, monospace;
    color: #333;
    font-weight: 500;
  }

  .log-train {
    font-size: 12px;
    color: #999;
  }

  .log-message {
    color: #333;
    margin-top: 4px;
    padding-top: 4px;
    border-top: 1px solid rgba(0, 0, 0, 0.05);
  }
}

.dialog-footer {
  width: 100%;
  display: flex;
  justify-content: flex-end;
}
</style>
