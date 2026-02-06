<template>
  <el-dropdown
    ref="notificationDropdownRef"
    :show-arrow="false"
    placement="bottom-end"
    popper-class="notification-dropdown-popper"
    trigger="click"
  >
    <el-badge
      :hidden="userStore.unreadCount === 0"
      :max="99"
      :offset="[-5, 5]"
      :value="userStore.unreadCount"
    >
      <HoverAnimateWrapper :duration="600" intensity="light" name="bell">
        <IconButton icon="HOutline:BellAlertIcon" />
      </HoverAnimateWrapper>
    </el-badge>

    <template #dropdown>
      <div class="notification-dropdown">
        <div class="notification-header">
          <span class="title">消息通知</span>
          <el-button
            v-if="userStore.unreadCount > 0"
            link
            size="small"
            type="primary"
            @click.stop="userStore.markAllAsRead()"
          >
            <el-icon class="button-icon"
              ><component :is="menuStore.iconComponents['Check']"
            /></el-icon>
            全部已读
          </el-button>
        </div>
        <div class="notification-list">
          <el-scrollbar max-height="400px">
            <Transition :style="{ '--animation-duration': '0.5s' }" mode="out-in" name="zoom">
              <div v-if="unreadMessageList.length === 0" class="empty-message">
                <el-empty :image-size="80" description="暂无消息" />
              </div>
              <TransitionGroup v-else name="group-slide-right" tag="div">
                <div
                  v-for="message in unreadMessageList"
                  :key="message.appID"
                  class="notification-item"
                  @click="userStore.markAsRead(message.appID)"
                >
                  <el-avatar :size="32" :src="message.avatar" />
                  <div class="message-content">
                    <div class="message-title">{{ message.title }}</div>
                    <div class="message-text">{{ message.content }}</div>
                    <div class="message-time">{{ message.time }}</div>
                  </div>
                </div>
              </TransitionGroup>
            </Transition>
          </el-scrollbar>
        </div>
        <div class="notification-footer">
          <el-button link type="primary" @click="goToProfile">查看全部消息</el-button>
        </div>
      </div>
    </template>
  </el-dropdown>
</template>

<script lang="ts" setup>
import type {ElDropdown} from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const menuStore = useMenuStore()
const notificationDropdownRef = ref<InstanceType<typeof ElDropdown>>()
// 未读消息列表
const unreadMessageList = computed(() => {
  return userStore.userMessages.filter((msg) => !msg.read)
})

// 跳转到个人中心
const goToProfile = () => {
  router.push('/profile')
  userStore.currentTab = 'messages'
  notificationDropdownRef.value?.handleClose()
}
</script>

<style lang="scss" scoped>
.notification-dropdown {
  width: 22rem;
  background: var(--el-bg-color);
  .notification-header {
    padding: 1rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-bottom: 1px solid var(--el-border-color-lighter);
    .title {
      font-size: 1rem;
      font-weight: 600;
      color: var(--el-text-color-primary);
    }
    .button-icon {
      margin-right: 0.25rem;
    }
  }
  .notification-list {
    max-height: 25rem;
    .notification-item {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 1rem;
      border-bottom: 1px solid var(--el-border-color-lighter);
      cursor: pointer;

      &:hover {
        background: var(--el-fill-color-light);
      }

      .message-content {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 6px;
        min-width: 0;

        .message-title {
          font-size: 14px;
          font-weight: 600;
          color: var(--el-text-color-primary);
          line-height: 1.4;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .message-text {
          font-size: 13px;
          color: var(--el-text-color-regular);
          line-height: 1.5;
          overflow: hidden;
          text-overflow: ellipsis;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          line-clamp: 2;
          -webkit-box-orient: vertical;
        }

        .message-time {
          font-size: 12px;
          color: var(--el-text-color-secondary);
          line-height: 1.4;
        }
      }
    }
  }
  .notification-footer {
    padding: 1rem;
    border-top: 1px solid var(--el-border-color-lighter);
    text-align: center;
  }
}
</style>

<style lang="scss">
.notification-dropdown-popper {
  border-radius: 8px !important;
  overflow: hidden;
  .el-dropdown-menu {
    border-radius: 8px !important;
    overflow: hidden;
  }
}
</style>
