<!-- 我的权限 -->
<template>
  <BaseCard>
    <template #header>
      <div class="flex flex-col md:flex-row md:items-end justify-between">
        <div class="space-y-2">
          <div class="flex items-center gap-3">
            <h1 class="text-2xl font-semibold">我的权限</h1>
            <BaseTag :text="userStore.userRoleName" />
          </div>
          <p class="text-sm text-(--el-text-color-secondary)">
            查看您在系统中获准访问的菜单项与操作功能。如有权限变动，请联系系统管理员。
          </p>
        </div>

        <div class="flex items-center justify-center gap-10 mt-6 md:mt-0 pr-4">
          <el-statistic
            :value="currentRolePermission?.length"
            class="text-center"
            title="已开启权限"
          />

          <el-divider direction="vertical" />

          <el-statistic :value="authorizedCount" class="text-center" title="权限总数" />
        </div>
      </div>
    </template>

    <el-table
      v-loading="loading"
      :border="TABLE_CONFIG.border"
      :data="menuList"
      :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
      default-expand-all
      row-key="id"
      show-overflow-tooltip
    >
      <el-table-column
        :align="TABLE_CONFIG.align"
        label="菜单/功能名称"
        min-width="200"
        prop="title"
      />
      <el-table-column :align="TABLE_CONFIG.align" label="类型" min-width="100" prop="type">
        <template #default="{ row }">
          <BaseTag v-if="row.type === 'directory'" text="目录" type="info" />
          <BaseTag v-else-if="row.type === 'menu'" text="菜单" type="primary" />
          <BaseTag v-else-if="row.type === 'button'" text="按钮" type="warning" />
        </template>
      </el-table-column>
      <el-table-column :align="TABLE_CONFIG.align" label="菜单路径" min-width="250" prop="path" />
      <el-table-column :align="TABLE_CONFIG.align" label="图标" min-width="100" prop="icon">
        <template #default="{ row }">
          <el-icon v-if="row.icon">
            <component :is="menuStore.iconComponents[row.icon]" />
          </el-icon>
        </template>
      </el-table-column>

      <el-table-column :align="TABLE_CONFIG.align" label="权限状态" width="150">
        <template #default="{ row }: { row: IMenuItem }">
          <BaseTag :text="getPermissionTag(row.id).text" :type="getPermissionTag(row.id).type" />
        </template>
      </el-table-column>
    </el-table>
  </BaseCard>
</template>

<script lang="ts" setup>
import {menuPage} from '@/api/menu'
import {TABLE_CONFIG} from '@/config/elementConfig'
import type {IMenuItem} from '@/types/system/menu'

const menuStore = useMenuStore()
const userStore = useUserStore()

const loading = ref(false)

// 菜单列表
const menuList = ref<IMenuItem[]>([])

// 总共的权限数量
const authorizedCount = computed(() => {
  const countAuth = (list: IMenuItem[]): number => {
    let count = 0
    list.forEach((item) => {
      if (item.id) count++
      if (item.children) count += countAuth(item.children)
    })
    return count
  }
  return countAuth(menuList.value)
})

// 当前角色权限id列表
const currentRolePermission = computed(() => {
  return userStore.roleList.filter((role) => role.id === userStore.userInfo?.roleId)[0]?.menuIds
})

// 获取当前权限标签
const getPermissionTag = (menuId: string): { type: 'success' | 'danger'; text: string } => {
  const enabled = currentRolePermission.value?.some((item) => item === menuId)
  return {
    type: enabled ? 'success' : 'danger',
    text: enabled ? '已启用' : '未授权',
  }
}

// 获取菜单列表
const getMenuList = async () => {
  loading.value = true
  try {
    const { data: res } = await menuPage()
    if (res.code !== 0) return
    menuList.value = res.data || []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  getMenuList()
})
</script>

<style lang="scss" scoped></style>
