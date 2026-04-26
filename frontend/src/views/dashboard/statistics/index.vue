<template>
  <div class="flex flex-col gap-4">
    <!-- 概览卡片 -->
    <el-row :gutter="20">
      <el-col v-for="item in overviewCards" :key="item.label" :lg="8" :xs="24">
        <BaseCard :class="item.type" style="height: 100%">
          <div class="h-full flex flex-col gap-2 justify-center p-5">
            <div class="w-12 h-12 bg-white/20 rounded-xl flex items-center justify-center">
              <el-icon size="24">
                <component :is="menuStore.iconComponents[item.icon]" class="text-white" />
              </el-icon>
            </div>
            <div class="text-white break-all">
              <span class="text-2xl font-bold">{{ item.value }}</span>
            </div>
            <div class="text-sm opacity-80 text-white">{{ item.label }}</div>
          </div>
        </BaseCard>
      </el-col>
    </el-row>

    <!-- 机器统计 + 许可证统计 -->
    <el-row :gutter="20">
      <el-col :lg="12" :xs="24">
        <BaseCard title="机器概览">
          <div class="grid grid-cols-2 gap-4 p-2">
            <div v-for="s in machineStatItems" :key="s.label" class="stat-item">
              <div class="stat-value">{{ s.value }}</div>
              <div class="stat-label">{{ s.label }}</div>
            </div>
          </div>
        </BaseCard>
      </el-col>
      <el-col :lg="12" :xs="24" class="mt-4 min-[1600px]:mt-0">
        <BaseCard title="许可证概览">
          <div class="grid grid-cols-2 gap-4 p-2">
            <div v-for="s in licenseStatItems" :key="s.label" class="stat-item">
              <div class="stat-value">{{ s.value }}</div>
              <div class="stat-label">{{ s.label }}</div>
            </div>
          </div>
        </BaseCard>
      </el-col>
    </el-row>

    <!-- 分布图表 -->
    <el-row :gutter="20">
      <el-col :lg="8" :md="12" :xs="24">
        <BaseCard title="平台分布">
          <div class="h-60 w-full">
            <VChart :option="platformOption" autoresize />
          </div>
        </BaseCard>
      </el-col>
      <el-col :lg="8" :md="12" :xs="24" class="mt-4 min-[992px]:mt-0">
        <BaseCard title="架构分布">
          <div class="h-60 w-full">
            <VChart :option="archOption" autoresize />
          </div>
        </BaseCard>
      </el-col>
      <el-col :lg="8" :xs="24" class="mt-4 min-[1200px]:mt-0">
        <BaseCard title="App 排行（按机器数）">
          <div class="h-60 w-full">
            <VChart :option="rankingOption" autoresize />
          </div>
        </BaseCard>
      </el-col>
    </el-row>

    <!-- 近30天趋势 -->
    <el-row :gutter="20">
      <el-col :lg="12" :xs="24">
        <BaseCard title="每日机器注册趋势（近30天）">
          <div class="h-64 w-full">
            <VChart :option="machineTrendOption" autoresize />
          </div>
        </BaseCard>
      </el-col>
      <el-col :lg="12" :xs="24" class="mt-4 min-[1600px]:mt-0">
        <BaseCard title="每日许可证签发趋势（近30天）">
          <div class="h-64 w-full">
            <VChart :option="licenseTrendOption" autoresize />
          </div>
        </BaseCard>
      </el-col>
    </el-row>

    <!-- App 排行表格 -->
    <el-card shadow="never">
      <template #header>
        <span>App 详细排行</span>
      </template>
      <el-table :border="TABLE_CONFIG.border" :data="data?.appRankings ?? []" show-overflow-tooltip>
        <el-table-column :align="TABLE_CONFIG.align" label="序号" type="index" width="55" />
        <el-table-column :align="TABLE_CONFIG.align" label="应用名称" min-width="150" prop="appName" />
        <el-table-column :align="TABLE_CONFIG.align" label="AppID" min-width="200" prop="appId" />
        <el-table-column :align="TABLE_CONFIG.align" label="绑定机器" min-width="100" prop="machines" />
        <el-table-column :align="TABLE_CONFIG.align" label="在线" min-width="80" prop="online" />
        <el-table-column :align="TABLE_CONFIG.align" label="许可证" min-width="100" prop="licenses" />
      </el-table>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import VChart from 'vue-echarts'
import { getStatistics, type StatisticsData } from '@/api/statistics'
import { TABLE_CONFIG } from '@/config/elementConfig'
import { ElMessage } from 'element-plus'

defineOptions({ name: 'StatisticsView' })

const menuStore = useMenuStore()

const data = ref<StatisticsData | null>(null)

const overviewCards = computed(() => [
  { label: '总机器数', value: data.value?.overview.totalMachines ?? '-', icon: 'HSolid:ComputerDesktopIcon', type: 'blue' },
  { label: '总应用数', value: data.value?.overview.totalApps ?? '-', icon: 'HSolid:Square3Stack3DIcon', type: 'green' },
  { label: '在线机器', value: data.value?.overview.onlineMachines ?? '-', icon: 'HSolid:SignalIcon', type: 'purple' },
])

const machineStatItems = computed(() => [
  { label: '今日新增', value: data.value?.machineStats.todayNew ?? '-' },
  { label: '本月新增', value: data.value?.machineStats.monthNew ?? '-' },
  { label: '7天活跃', value: data.value?.machineStats.active7d ?? '-' },
  { label: '30天活跃', value: data.value?.machineStats.active30d ?? '-' },
])

const licenseStatItems = computed(() => [
  { label: '有效', value: data.value?.licenseStats.valid ?? '-' },
  { label: '已吊销', value: data.value?.licenseStats.revoked ?? '-' },
  { label: '已过期', value: data.value?.licenseStats.expired ?? '-' },
  { label: '今日签发', value: data.value?.licenseStats.todayNew ?? '-' },
  { label: '本月签发', value: data.value?.licenseStats.monthNew ?? '-' },
])

const style = () => getComputedStyle(document.documentElement)

const platformOption = computed(() => ({
  tooltip: { trigger: 'item' as const },
  series: [{
    type: 'pie',
    radius: ['40%', '70%'],
    center: ['50%', '50%'],
    itemStyle: { borderRadius: 6, borderColor: style().getPropertyValue('--el-bg-color-overlay'), borderWidth: 3 },
    label: { show: true, formatter: '{b}: {c}' },
    data: (data.value?.distributions.platform ?? []).map((item, i) => ({
      value: item.count,
      name: item.name || '未知',
      itemStyle: { color: [style().getPropertyValue('--el-color-primary'), style().getPropertyValue('--el-color-success'), style().getPropertyValue('--el-color-warning'), style().getPropertyValue('--el-color-danger'), style().getPropertyValue('--el-color-info')][i % 5] },
    })),
  }],
}))

const archOption = computed(() => ({
  tooltip: { trigger: 'item' as const },
  series: [{
    type: 'pie',
    radius: ['40%', '70%'],
    center: ['50%', '50%'],
    itemStyle: { borderRadius: 6, borderColor: style().getPropertyValue('--el-bg-color-overlay'), borderWidth: 3 },
    label: { show: true, formatter: '{b}: {c}' },
    data: (data.value?.distributions.arch ?? []).map((item, i) => ({
      value: item.count,
      name: item.name || '未知',
      itemStyle: { color: [style().getPropertyValue('--el-color-primary'), style().getPropertyValue('--el-color-warning'), style().getPropertyValue('--el-color-success'), style().getPropertyValue('--el-color-info')][i % 4] },
    })),
  }],
}))

const rankingOption = computed(() => ({
  tooltip: { trigger: 'axis' as const, axisPointer: { type: 'shadow' as const } },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { type: 'value' as const },
  yAxis: { type: 'category' as const, data: (data.value?.appRankings ?? []).map(r => r.appName).reverse() },
  series: [{
    type: 'bar',
    data: (data.value?.appRankings ?? []).map(r => r.machines).reverse(),
    itemStyle: { color: style().getPropertyValue('--el-color-primary'), borderRadius: [0, 6, 6, 0] },
  }],
}))

const machineTrendOption = computed(() => ({
  tooltip: { trigger: 'axis' as const },
  grid: { left: '3%', right: '4%', bottom: '10%', containLabel: true },
  xAxis: { type: 'category' as const, data: (data.value?.machineTrend ?? []).map(r => r.date.slice(5)), axisLabel: { rotate: 45 } },
  yAxis: { type: 'value' as const },
  series: [{
    type: 'line',
    smooth: true,
    data: (data.value?.machineTrend ?? []).map(r => r.count),
    itemStyle: { color: style().getPropertyValue('--el-color-primary') },
    areaStyle: { color: style().getPropertyValue('--el-color-primary-light-7') },
  }],
}))

const licenseTrendOption = computed(() => ({
  tooltip: { trigger: 'axis' as const },
  grid: { left: '3%', right: '4%', bottom: '10%', containLabel: true },
  xAxis: { type: 'category' as const, data: (data.value?.licenseTrend ?? []).map(r => r.date.slice(5)), axisLabel: { rotate: 45 } },
  yAxis: { type: 'value' as const },
  series: [{
    type: 'line',
    smooth: true,
    data: (data.value?.licenseTrend ?? []).map(r => r.count),
    itemStyle: { color: style().getPropertyValue('--el-color-success') },
    areaStyle: { color: style().getPropertyValue('--el-color-success-light-7') },
  }],
}))

const fetchData = async () => {
  const { data: res } = await getStatistics()
  if (res.code !== 0) {
    ElMessage.error(res.msg || '获取统计数据失败')
    return
  }
  data.value = res.data
}

onMounted(fetchData)
</script>

<style scoped>
.stat-item {
  background: var(--el-fill-color-lighter);
  border-radius: 12px;
  padding: 16px;
  text-align: center;
}
.stat-value {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--el-color-primary);
}
.stat-label {
  font-size: 0.85rem;
  color: var(--el-text-color-secondary);
  margin-top: 4px;
}
.blue {
  background: linear-gradient(135deg, #5bbff9 0%, #2563eb 100%);
}
.green {
  background: linear-gradient(135deg, #86efac 0%, #16a34a 100%);
}
.purple {
  background: linear-gradient(135deg, #c4b5fd 0%, #7e22ce 100%);
}
</style>
