<!-- 我的消息 -->
<template>
  <div>
    <BaseCard>
      <div class="flex items-center gap-4">
        <el-avatar :size="32" :src="userStore.userInfo?.avatar" />
        <span class="text-sm font-medium text-(--el-text-color-secondary)"
          >Hi, {{ userStore.userInfo?.name }}，可以在这里发送新通知哦。</span
        >
      </div>
      <el-input
        v-model="postContent"
        :rows="3"
        class="mt-4"
        placeholder="输入消息内容..."
        type="textarea"
      />
      <div class="flex items-center justify-between mt-4">
        <div class="text-xs text-(--el-text-color-secondary)">将推送给所有相关人员</div>
        <el-button :disabled="!postContent.trim()" type="primary" @click="sendMessage"
          >发布消息</el-button
        >
      </div>
    </BaseCard>

    <BaseCard class="mt-4">
      <template #header>
        <div class="flex items-center justify-between">
          <div class="flex-1">
            <BadgeTabsMenu v-model="activeName" :tabs-item-height="30" :tabs-menu-data="tabsMenu" />
          </div>
          <div class="flex items-center">
            <IconButton
              v-if="menuStore.isMobile"
              :disabled="!userStore.unreadCount"
              icon="Element:Check"
              iconSize="1rem"
              size="1.5rem"
              tooltip="一键已读"
              type="primary"
              @click="userStore.markAllAsRead()"
            />
            <el-button
              v-else
              :disabled="!userStore.unreadCount"
              link
              type="primary"
              @click="userStore.markAllAsRead()"
            >
              一键已读
            </el-button>
            <el-divider direction="vertical" />
            <IconButton
              v-if="menuStore.isMobile"
              :disabled="!userStore.userMessages.length"
              icon="Element:Delete"
              iconSize="1rem"
              size="1.5rem"
              tooltip="清空全部"
              type="danger"
              @click="clearAllMessages"
            />
            <el-button
              v-else
              :disabled="!userStore.userMessages.length"
              link
              type="danger"
              @click="clearAllMessages"
            >
              清空全部
            </el-button>
          </div>
        </div>
      </template>
      <div>
        <Transition mode="out-in" name="zoom">
          <el-empty
            v-if="messageList.length === 0"
            :description="activeName === 'unread' ? '暂无未读消息' : '暂无消息'"
          />
          <TransitionGroup v-else name="group-slide-right" tag="div">
            <div v-for="message in messageList" :key="message.appID">
              <HoverAnimateWrapper class="w-full" intensity="light" name="lift">
                <div
                  class="group p-4 mb-3 flex items-center gap-4 border border-(--el-border-color-light) rounded-xl cursor-pointer hover:border-(--el-border-color) hover:bg-(--el-bg-color-page)"
                >
                  <div class="relative">
                    <el-avatar :size="48" :src="message.avatar" />
                    <span
                      v-if="!message.read"
                      class="absolute h-3 w-3 bottom-1.5 right-0.5 rounded-full border-3 border-(--el-bg-color) bg-(--el-color-danger)"
                    ></span>
                  </div>

                  <div class="flex-1">
                    <div class="flex justify-between">
                      <TextEllipsis :clickable="false" :text="message.title" tooltipType="none" />
                      <div
                        class="flex items-center opacity-100 lg:opacity-0 group-hover:opacity-100"
                      >
                        <IconButton
                          v-if="!message.read"
                          icon="Element:Check"
                          iconSize="1rem"
                          size="1.5rem"
                          tooltip="设为已读"
                          type="primary"
                          @click="userStore.markAsRead(message.appID)"
                        />
                        <el-divider v-if="!message.read" direction="vertical" />
                        <el-popconfirm
                          title="确定删除这条消息吗？"
                          @confirm="
                            (userStore.deleteMessage(message.appID), ElMessage.success('删除成功'))
                          "
                        >
                          <template #reference>
                            <div>
                              <IconButton
                                icon="Element:Delete"
                                iconSize="1rem"
                                size="1.5rem"
                                tooltip="删除"
                                type="danger"
                              />
                            </div>
                          </template>
                        </el-popconfirm>
                      </div>
                    </div>
                    <div
                      class="mt-2 text-sm text-(--el-text-color-regular) leading-relaxed wrap-break-word"
                    >
                      {{ message.content }}
                    </div>
                    <div class="text-xs text-(--el-text-color-secondary) mt-2">
                      {{ message.time }}
                    </div>
                  </div>
                </div>
              </HoverAnimateWrapper>
            </div>
          </TransitionGroup>
        </Transition>
      </div>
    </BaseCard>
  </div>
</template>

<script lang="ts" setup>
import {Dialog} from '@/utils/dialog'
import {delay} from '@/utils/utils'
import BadgeTabsMenu from '@/components/tabs/BadgeTabsMenu.vue'
import {ElMessage} from 'element-plus'

const userStore = useUserStore()
const menuStore = useMenuStore()

// 消息内容
const postContent = ref('')
// 当前菜单
const activeName = ref<'all' | 'unread'>('all')

// 菜单
const tabsMenu = computed(() => [
  { key: 'all', label: '全部消息', badge: 0 },
  { key: 'unread', label: '未读消息', badge: userStore.unreadCount },
])

// 消息列表
const messageList = computed(() => {
  if (activeName.value === 'unread') {
    return userStore.userMessages.filter((item) => !item.read)
  }
  return userStore.userMessages
})

// 发送消息
const sendMessage = () => {
  userStore.sendMessage(postContent.value)
  ElMessage.success('发送成功')
  postContent.value = ''
}

// 清空全部消息
const clearAllMessages = () => {
  Dialog.confirm({
    title: '确认清空？',
    content: '这一操作会删除所有消息，手滑之后可就找不回来了哦～',
    onConfirm: async () => {
      await delay(1000)
      userStore.deleteAllMessages()
      ElMessage.success('消息清空完成')
    },
  })
}
</script>

<style lang="scss" scoped></style>
