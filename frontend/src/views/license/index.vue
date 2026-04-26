<template>
  <div>
    <el-card shadow="never">
      <el-table :border="TABLE_CONFIG.border" :data="list" show-overflow-tooltip>
        <el-table-column :align="TABLE_CONFIG.align" label="序号" type="index" width="55" />
        <el-table-column :align="TABLE_CONFIG.align" label="授权Key" min-width="220" prop="licenseKey" />
        <el-table-column :align="TABLE_CONFIG.align" label="应用" min-width="120" prop="appName" />
        <el-table-column :align="TABLE_CONFIG.align" label="用户" min-width="100" prop="username" />
        <el-table-column :align="TABLE_CONFIG.align" label="机器ID" min-width="80" prop="machineId" />
        <el-table-column :align="TABLE_CONFIG.align" label="状态" min-width="80">
          <template #default="{ row }">
            <BaseTag :text="statusMap[row.status]?.text ?? '未知'" :type="statusMap[row.status]?.type ?? 'info'" />
          </template>
        </el-table-column>
        <el-table-column :align="TABLE_CONFIG.align" label="签发时间" min-width="160" prop="issuedAt" />
        <el-table-column :align="TABLE_CONFIG.align" label="过期时间" min-width="160" prop="expireAt" />
      </el-table>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import { getLicenseList, type LicenseItem } from '@/api/license'
import { TABLE_CONFIG } from '@/config/elementConfig'
import { ElMessage } from 'element-plus'

defineOptions({ name: 'LicenseView' })

const list = ref<LicenseItem[]>([])

const statusMap: Record<number, { text: string; type: string }> = {
  0: { text: '有效', type: 'success' },
  1: { text: '已吊销', type: 'danger' },
  2: { text: '已过期', type: 'info' },
}

const fetchData = async () => {
  const { data: res } = await getLicenseList()
  if (res.code !== 0) {
    ElMessage.error(res.msg || '获取授权列表失败')
    return
  }
  list.value = res.data
}

onMounted(fetchData)
</script>
