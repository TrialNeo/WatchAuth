<template>
  <div class="flex-1 h-full flex flex-col gap-6">
    <OverviewHeader />
    <RevenueProfitAnalysis ref="revenueProfitAnalysisRef" />
    <el-row :gutter="20">
      <!--全球市场份额分布-->
      <el-col :lg="8" :md="12" :xs="24">
        <MarketShare ref="marketShareRef" />
      </el-col>
      <!--热销商品类目 TOP 5-->
      <el-col :lg="8" :md="12" :xs="24" class="mt-4 min-[992px]:mt-0">
        <TopCategories ref="topCategoriesRef" />
      </el-col>
      <!--营销目标和今日之星-->
      <el-col :lg="8" :md="24" :xs="24" class="mt-4 min-[1200px]:mt-0">
        <GoalsAndTodayStart />
      </el-col>
    </el-row>
    <el-row :gutter="20">
      <el-col :lg="24" :xs="24">
        <!--各渠道销售表现实时榜单-->
        <ChannelSales />
      </el-col>
      <!--近期运营大事件-->
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
