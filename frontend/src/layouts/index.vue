<template>
  <div class="layout-container">
    <LeftMode v-if="themeStore.layout === 'leftMode'" />
    <TopMode v-else />
    <ThemeConfig />

    <!-- 移动端菜单抽屉 -->
    <el-drawer
      v-model="menuStore.isMobileMenuOpen"
      :direction="'ltr'"
      :size="220"
      :with-header="false"
      class="mobile-menu-drawer"
    >
      <MenuView />
    </el-drawer>
  </div>
</template>

<script lang="ts" setup>
import LeftMode from '@/layouts/leftMode.vue'
import TopMode from '@/layouts/topMode.vue'
import ThemeConfig from '@/components/ThemeConfig.vue'
import MenuView from '@/layouts/menu.vue'
defineOptions({ name: 'LayoutView' })

const menuStore = useMenuStore()

const themeStore = useThemeStore()
</script>

<style lang="scss" scoped>
:deep(.mobile-menu-drawer) {
  .el-drawer__body {
    padding: 0;
  }
}
.layout-container {
  width: 100%;
  height: 100vh;
}
</style>
