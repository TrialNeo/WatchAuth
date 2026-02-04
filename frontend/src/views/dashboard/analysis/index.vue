<template>
  <div class="flex-1 h-full flex flex-col gap-6">
    <OverviewHeader />
    <RevenueProfitAnalysis ref="revenueProfitAnalysisRef" />
    <el-row :gutter="20">
      <el-col :lg="8" :md="12" :xs="24">
        <MarketShare ref="marketShareRef" />
      </el-col>
      <el-col :lg="8" :md="12" :xs="24" class="mt-4 min-[992px]:mt-0">
        <TopCategories ref="topCategoriesRef" />
      </el-col>
      <el-col :lg="8" :md="24" :xs="24" class="mt-4 min-[1200px]:mt-0">
        <GoalsAndTodayStart />
      </el-col>
    </el-row>
    <el-row :gutter="20">
      <el-col :lg="14" :xs="24">
        <ChannelSales />
      </el-col>
      <el-col :lg="10" :xs="24" class="mt-4 min-[1200px]:mt-0">
        <OperationalEvent />
      </el-col>
    </el-row>
  </div>
</template>

<script lang="ts" setup>
import OverviewHeader from '@/views/dashboard/analysis/overviewHeader.vue'
import RevenueProfitAnalysis from '@/views/dashboard/analysis/revenueProfitAnalysis.vue'
import MarketShare from '@/views/dashboard/analysis/marketShare.vue'
import TopCategories from '@/views/dashboard/analysis/topCategories.vue'
import GoalsAndTodayStart from '@/views/dashboard/analysis/goalsAndTodayStart.vue'
import ChannelSales from '@/views/dashboard/analysis/channelSales.vue'
import OperationalEvent from '@/views/dashboard/analysis/operationalEvent.vue'

defineOptions({ name: 'AnalysisView' })

const themeStore = useThemeStore()
const marketShareRef = useTemplateRef<InstanceType<typeof MarketShare> | null>('marketShareRef')
const topCategoriesRef = useTemplateRef<InstanceType<typeof TopCategories> | null>(
  'topCategoriesRef',
)
const revenueProfitAnalysisRef = useTemplateRef<InstanceType<typeof RevenueProfitAnalysis> | null>(
  'revenueProfitAnalysisRef',
)

//  监听主题色和主题模式变化，更新图表颜色
watch(
  [() => themeStore.themeMode, () => themeStore.primaryColor],
  async () => {
    await nextTick()
    marketShareRef.value?.updateColorTrigger()
    topCategoriesRef.value?.updateColorTrigger()
    revenueProfitAnalysisRef.value?.updateColorTrigger()
  },
  { immediate: true },
)
</script>

<style></style>
