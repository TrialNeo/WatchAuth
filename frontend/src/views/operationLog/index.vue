<template>
  <div>
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>操作日志</span>
        </div>
      </template>
      <div class="mb-4 flex gap-4">
        <el-select v-model="filterAction" clearable placeholder="操作类型" style="width: 140px">
          <el-option label="全部" value="" />
          <el-option label="登录" value="login" />
          <el-option label="创建" value="create" />
          <el-option label="更新" value="update" />
          <el-option label="删除" value="delete" />
        </el-select>
        <el-input v-model="filterTarget" clearable placeholder="搜索目标" style="width: 200px" @clear="fetchData" />
        <el-button type="primary" @click="fetchData">搜索</el-button>
      </div>
      <el-table :border="TABLE_CONFIG.border" :data="list" show-overflow-tooltip>
        <el-table-column :align="TABLE_CONFIG.align" label="序号" type="index" width="55" />
        <el-table-column :align="TABLE_CONFIG.align" label="操作者" min-width="100" prop="adminName" />
        <el-table-column :align="TABLE_CONFIG.align" label="操作类型" min-width="90">
          <template #default="{ row }">
            <BaseTag :text="actionMap[row.action] ?? row.action" :type="actionType(row.action)" />
          </template>
        </el-table-column>
        <el-table-column :align="TABLE_CONFIG.align" label="目标" min-width="100" prop="target" />
        <el-table-column :align="TABLE_CONFIG.align" label="详情" min-width="240" prop="detail" show-overflow-tooltip />
        <el-table-column :align="TABLE_CONFIG.align" label="IP" min-width="130" prop="ip" />
        <el-table-column :align="TABLE_CONFIG.align" label="时间" min-width="160" prop="createdAt" />
      </el-table>
      <div class="mt-4 flex justify-center">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          background
          @current-change="fetchData"
        />
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import { getOperationLogs, type OperationLogItem } from '@/api/operationLog'
import { TABLE_CONFIG } from '@/config/elementConfig'
import { ElMessage } from 'element-plus'

defineOptions({ name: 'OperationLogView' })

const list = ref<OperationLogItem[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const filterAction = ref('')
const filterTarget = ref('')

const actionMap: Record<string, string> = {
  login: '登录',
  create: '创建',
  update: '更新',
  delete: '删除',
}

const actionType = (action: string) => {
  const map: Record<string, string> = { login: 'info', create: 'success', update: 'warning', delete: 'danger' }
  return map[action] ?? 'info'
}

const fetchData = async () => {
  const { data: res } = await getOperationLogs({
    page: page.value,
    pageSize: pageSize.value,
    action: filterAction.value || undefined,
    target: filterTarget.value || undefined,
  })
  if (res.code !== 0) {
    ElMessage.error(res.msg || '获取日志失败')
    return
  }
  list.value = res.data?.list ?? []
  total.value = res.data?.total ?? 0
}

onMounted(fetchData)
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
