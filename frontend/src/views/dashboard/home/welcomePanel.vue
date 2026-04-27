<template>
  <BaseCard>
    <el-scrollbar :max-height="620">
      <div class="flex flex-col xl:flex-row justify-between p-6 lg:p-8">
        <div class="flex-1">
          <div class="flex flex-col lg:flex-row items-center lg:items-start xl:items-center gap-6">
            <div class="relative shrink-0">
              <el-avatar :size="110" :src="userStore.userInfo?.avatar" />
              <div
                class="absolute h-5 w-5 bottom-2 right-2 rounded-full border-3 border-(--el-bg-color) bg-(--el-color-success)"
              ></div>
            </div>

            <div class="flex flex-col gap-4 items-center lg:items-start text-center lg:text-left">
              <h2
                class="flex text-2xl md:text-3xl font-black text-(--el-text-color-primary) cursor-pointer"
              >
                <TextEllipsis
                  :clickable="false"
                  :text="`${userStore.userInfo?.name! || userStore.userInfo?.username!}，欢迎回来！`"
                  class="text-2xl md:text-3xl font-black text-(--el-text-color-primary)"
                />
                <div>👋</div>
              </h2>
              <TextEllipsis
                :text="`“ ${userStore.userInfo?.bio} ”`"
                class="text-(--el-text-color-regular) italic text-sm md:base cursor-pointer"
              />
            </div>
          </div>

          <!-- 公告栏 -->
          <div v-if="announcements.length > 0" class="mt-6 px-2 md:px-6">
            <div class="flex items-center gap-2 mb-3">
              <el-icon class="text-orange-500" size="18"><component :is="menuStore.iconComponents['HOutline:MegaphoneIcon']" /></el-icon>
              <span class="text-sm font-semibold text-(--el-text-color-primary)">系统公告</span>
            </div>
            <div class="flex flex-col gap-2">
              <div
                v-for="item in announcements.slice(0, 3)"
                :key="item.id"
                class="flex items-start gap-2 p-2 rounded-lg bg-(--el-fill-color-lighter) cursor-pointer hover:bg-(--el-fill-color-light)"
              >
                <span class="text-xs font-bold text-(--el-color-primary) shrink-0 mt-0.5">[公告]</span>
                <div class="flex-1 min-w-0">
                  <div class="text-sm font-semibold truncate">{{ item.title }}</div>
                  <div class="text-xs text-(--el-text-color-secondary) mt-0.5 truncate">{{ item.content }}</div>
                </div>
                <span class="text-xs text-(--el-text-color-secondary) shrink-0">{{ item.createdAt?.slice(0, 10) }}</span>
              </div>
            </div>
          </div>

          <!-- 今日概览 -->
          <div class="flex flex-col md:flex-row px-2 md:px-6 py-6 md:py-10 items-center gap-6 md:gap-0">
            <div class="flex w-full md:flex-1 flex-col gap-2">
              <div class="text-xs font-semibold text-(--el-text-color-secondary)">今日新增机器</div>
              <div class="flex items-center gap-2">
                <span class="text-xl font-extrabold text-(--el-color-primary)">{{ todayStats?.newMachines ?? '-' }}</span>
              </div>
            </div>
            <div class="hidden md:block mx-7"><el-divider direction="vertical" /></div>
            <div class="flex w-full md:flex-1 flex-col gap-2">
              <div class="text-xs font-semibold text-(--el-text-color-secondary)">今日新发授权</div>
              <div class="flex items-center gap-2">
                <span class="text-xl font-extrabold text-(--el-color-primary)">{{ todayStats?.newLicenses ?? '-' }}</span>
              </div>
            </div>
            <div class="hidden md:block mx-7"><el-divider direction="vertical" /></div>
            <div class="flex w-full md:flex-1 flex-col gap-2">
              <div class="text-xs font-semibold text-(--el-text-color-secondary)">当前在线</div>
              <div class="flex items-center gap-2">
                <span class="text-xl font-extrabold text-(--el-color-primary)">{{ todayStats?.onlineCount ?? '-' }}</span>
                <BaseTag text="Live" type="success" />
              </div>
            </div>
          </div>
        </div>

        <div class="hidden xl:block mx-7 my-6">
          <div class="w-px h-full border-(--el-border-color) border-l"></div>
        </div>

        <div class="flex-1 xl:flex-[0.8] grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div
            v-for="item in statCards"
            :key="item.label"
            class="flex flex-col justify-between p-4 rounded-2xl relative transition-all duration-300 cursor-pointer hover:shadow-xl hover:-translate-y-1"
          >
            <div
              :style="{ color: item.color, backgroundColor: item.color + '10' }"
              class="flex items-center justify-center w-9 h-9 rounded-[10px] p-2 mb-3"
            >
              <el-icon size="18">
                <component :is="menuStore.iconComponents[item.icon]" />
              </el-icon>
            </div>
            <div>
              <div class="text-[13px] font-semibold text-(--el-text-color-secondary) mb-1">{{ item.label }}</div>
              <div class="flex items-baseline gap-2">
                <span class="text-[20px] font-extrabold text-(--el-text-color-primary)">{{ item.value }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 最近操作 -->
      <div class="px-6 pb-6" v-if="recentLogs.length > 0">
        <div class="flex items-center gap-2 mb-3">
          <el-icon class="text-indigo-500" size="18"><component :is="menuStore.iconComponents['HOutline:ClipboardDocumentListIcon']" /></el-icon>
          <span class="text-sm font-semibold text-(--el-text-color-primary)">最近操作</span>
        </div>
        <el-table :data="recentLogs" :border="false" size="small" show-overflow-tooltip>
          <el-table-column label="操作" min-width="100" prop="action" />
          <el-table-column label="目标" min-width="100" prop="target" />
          <el-table-column label="详情" min-width="200" prop="detail" show-overflow-tooltip />
          <el-table-column label="时间" min-width="140" prop="createdAt" />
        </el-table>
      </div>
    </el-scrollbar>
  </BaseCard>
</template>

<script lang="ts" setup>
import VChart from 'vue-echarts'
import LottieAnimation from '@/components/animation/LottieAnimation.vue'
import workTimeLottie from '@/assets/lotties/welcome.json'
import { getTodayStats, getRecentLogs, type TodayStats, type OperationLogItem } from '@/api/operationLog'
import { getActiveAnnouncements, type AnnouncementItem } from '@/api/announcement'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const menuStore = useMenuStore()

const todayStats = ref<TodayStats | null>(null)
const recentLogs = ref<OperationLogItem[]>([])
const announcements = ref<AnnouncementItem[]>([])

const statCards = computed(() => [
  { label: '今日新增用户', value: todayStats.value?.newUsers ?? '-', color: '#6366f1', icon: 'HOutline:UserPlusIcon' },
  { label: '新增授权', value: todayStats.value?.newLicenses ?? '-', color: '#10b981', icon: 'HOutline:KeyIcon' },
  { label: '代理总数', value: todayStats.value?.totalAgents ?? '-', color: '#f59e0b', icon: 'HOutline:UserGroupIcon' },
  { label: '在线机器', value: todayStats.value?.onlineCount ?? '-', color: '#ef4444', icon: 'HOutline:SignalIcon' },
])

const fetchData = async () => {
  const [statsRes, logsRes, annRes] = await Promise.all([
    getTodayStats(),
    getRecentLogs(),
    getActiveAnnouncements(),
  ])

  if (statsRes.data.code === 0) todayStats.value = statsRes.data.data
  else ElMessage.error(statsRes.data.msg || '获取今日统计失败')

  if (logsRes.data.code === 0) recentLogs.value = logsRes.data.data?.list ?? []

  if (annRes.data.code === 0) announcements.value = annRes.data.data?.list ?? []
}

onMounted(fetchData)
</script>

<style lang="scss" scoped>
.el-divider--vertical {
  height: 2.5rem;
}
</style>
