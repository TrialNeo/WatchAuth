<!-- 登录日志 -->
<template>
  <BaseCard>
    <el-empty v-if="!userStore.userInfo?.loginLogs?.length" description="暂无登录日志" />
    <div v-else>
      <div class="flex justify-end mb-4">
        <el-button
          :disabled="!userStore.userInfo?.loginLogs.length"
          type="primary"
          @click="exportLoginLogsExcel"
          >导出日志</el-button
        >
      </div>
      <el-table
        :border="TABLE_CONFIG.border"
        :data="userStore.userInfo?.loginLogs"
        class="custom-modern-table"
        show-overflow-tooltip
      >
        <el-table-column label="设备型号" min-width="150" prop="device" />
        <el-table-column label="浏览器/版本" min-width="200" prop="browser" />
        <el-table-column label="IP 地址" min-width="150" prop="ip" />
        <el-table-column label="地理位置" min-width="180" prop="location" />
        <el-table-column label="登录时间" min-width="170" prop="time" />
        <el-table-column label="结果" width="100">
          <template #default="{ row }">
            <BaseTag :text="row.status === 'success' ? '成功' : '失败'" :type="row.status" />
          </template>
        </el-table-column>
      </el-table>
    </div>
  </BaseCard>
</template>

<script lang="ts" setup>
import { TABLE_CONFIG } from '@/config/elementConfig'
import { exportToExcel } from '@/utils/exportExcel'

const userStore = useUserStore()

// 导出登录日志为Excel
const exportLoginLogsExcel = async () => {
  exportToExcel({
    fileName: '登录日志.xlsx',
    sheetName: '登录日志',
    data: userStore.userInfo?.loginLogs || [],
    columns: {
      device: '设备型号',
      browser: '浏览器/版本',
      ip: 'IP 地址',
      location: '地理位置',
      time: '登录时间',
      status: '结果',
    },
  })
}
</script>

<style lang="scss" scoped></style>
