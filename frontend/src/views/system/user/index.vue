<template>
  <div>
    <el-card class="card-clear-mb" shadow="never">
      <el-form ref="queryFormRef" :model="queryForm" label-width="auto" @keyup.enter="getUserList">
        <el-row :gutter="10">
          <el-col :lg="6" :md="12" :sm="12" :xl="6" :xs="24">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="queryForm.username" placeholder="请输入" />
            </el-form-item>
          </el-col>
          <el-col :lg="6" :md="12" :sm="12" :xl="6" :xs="24">
            <el-form-item label="姓名" prop="name">
              <el-input v-model="queryForm.name" placeholder="请输入" />
            </el-form-item>
          </el-col>
          <el-col :lg="6" :md="12" :sm="12" :xl="6" :xs="24">
            <el-form-item label="状态" prop="status">
              <el-select v-model="queryForm.status" placeholder="请选择">
                <el-option label="启用" value="active" />
                <el-option label="禁用" value="inactive" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :lg="6" :md="12" :sm="12" :xl="6" :xs="24">
            <el-form-item>
              <el-button :icon="menuStore.iconComponents.Search" type="primary" @click="getUserList"
                >搜索</el-button
              >
              <el-button :icon="menuStore.iconComponents.Refresh" @click="reset">重置</el-button>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>
    <el-card class="card-mt-16" shadow="never">
      <div class="operation-container">
        <el-button
          v-permission="['user:add']"
          :icon="menuStore.iconComponents.Plus"
          type="primary"
          @click="userCreateRef?.showDialog(undefined)"
          >新增用户</el-button
        >
        <el-popconfirm
          :placement="POPCONFIRM_CONFIG.placement"
          :width="POPCONFIRM_CONFIG.width"
          title="确定要删除选中的用户吗？"
          @confirm="deleteUserHandle(deleteUserIds)"
        >
          <template #reference>
            <el-button
              :disabled="
                !useButtonPermission(['user:delete'], [() => !!deleteUserIds.length]).value
              "
              :icon="menuStore.iconComponents.Delete"
              type="danger"
            >
              批量删除
            </el-button>
          </template>
        </el-popconfirm>
      </div>
      <el-table
        :border="TABLE_CONFIG.border"
        :data="userList"
        show-overflow-tooltip
        @selection-change="tableSelectionChange"
        @sort-change="tableSortChange"
      >
        <el-table-column :align="TABLE_CONFIG.align" type="selection" width="55" />
        <el-table-column :align="TABLE_CONFIG.align" fixed label="序号" type="index" width="55" />
        <el-table-column
          :align="TABLE_CONFIG.align"
          fixed
          label="用户名"
          min-width="160"
          prop="username"
        />
        <el-table-column :align="TABLE_CONFIG.align" label="姓名" min-width="120" prop="name" />
        <el-table-column :align="TABLE_CONFIG.align" label="手机号" min-width="120" prop="phone" />
        <el-table-column :align="TABLE_CONFIG.align" label="邮箱" min-width="180" prop="email" />
        <el-table-column :align="TABLE_CONFIG.align" label="角色" min-width="150" prop="roleId">
          <template #default="{ row }">
            <BaseTag v-if="row.roleId" :text="getRoleName(row.roleId)" />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :align="TABLE_CONFIG.align" label="类型" prop="isBuiltIn">
          <template #default="{ row }">
            <BaseTag v-if="row.isBuiltIn" text="内置" type="warning" />
            <BaseTag v-else text="自定义" type="success" />
          </template>
        </el-table-column>
        <el-table-column :align="TABLE_CONFIG.align" label="状态" prop="status">
          <template #default="{ row }">
            <BaseTag
              v-if="row.status === 'active'"
              :text="row.status === 'active' ? '启用' : '禁用'"
              :type="row.status === 'active' ? 'success' : 'danger'"
            />
          </template>
        </el-table-column>
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="创建时间"
          min-width="180"
          prop="createTime"
          sortable="custom"
        />
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="更新时间"
          min-width="180"
          prop="updateTime"
        />
        <el-table-column :align="TABLE_CONFIG.align" fixed="right" label="操作" width="150">
          <template #default="{ row }: { row: IUserItem }">
            <el-button
              v-permission="['user:edit']"
              :icon="menuStore.iconComponents.Edit"
              link
              type="primary"
              @click="userCreateRef?.showDialog(row.id)"
            >
              编辑
            </el-button>
            <el-popconfirm
              :placement="POPCONFIRM_CONFIG.placement"
              :width="POPCONFIRM_CONFIG.width"
              title="确定要删除选中的用户吗？"
              @confirm="deleteUserHandle([row.id])"
            >
              <template #reference>
                <el-button
                  v-permission="['user:delete']"
                  :icon="menuStore.iconComponents.Delete"
                  link
                  type="danger"
                >
                  删除
                </el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :layout="
            menuStore.isMobile ? PAGINATION_CONFIG.mobileLayout : PAGINATION_CONFIG.desktopLayout
          "
          :page-sizes="PAGINATION_CONFIG.pageSizes"
          :pager-count="
            menuStore.isMobile
              ? PAGINATION_CONFIG.mobilePagerCount
              : PAGINATION_CONFIG.desktopPagerCount
          "
          :total="pagination.total"
          @change="getUserList"
        />
      </div>
    </el-card>

    <UserCreate ref="userCreateRef" @refresh="refresh" />
  </div>
</template>

<script lang="ts" setup>
import { userPage, deleteUser } from '@/api/user'
import { rolePage } from '@/api/role'
import { useButtonPermission } from '@/composables/useButtonPermission'
import { PAGINATION_CONFIG, POPCONFIRM_CONFIG, TABLE_CONFIG } from '@/config/elementConfig'
import UserCreate from '@/views/system/user/create.vue'
import type { IUserItem } from '@/types/system/user'
import type { IRoleItem } from '@/types/system/role'
import type { FormInstance } from 'element-plus'

defineOptions({ name: 'UserView' })

const menuStore = useMenuStore()
const queryFormRef = useTemplateRef<FormInstance>('queryFormRef')
const userCreateRef = useTemplateRef<InstanceType<typeof UserCreate> | null>('userCreateRef')

// 删除用户的ids
const deleteUserIds = ref<string[]>([])

// 角色列表（用于显示角色名称）
const roleList = ref<IRoleItem[]>([])

// 查询表单
const queryForm = ref({
  username: '',
  name: '',
  status: undefined,
  sortOrder: 'desc' as 'asc' | 'desc',
})

// 用户列表
const userList = ref<IUserItem[]>([])

// 分页
const pagination = ref({
  page: 1,
  pageSize: 10,
  total: 0,
})

// 获取角色名称
const getRoleName = (roleId: string): string => {
  const role = roleList.value.find((r) => r.id === roleId)
  return role?.name || roleId
}

// 获取角色列表
const getRoleList = async () => {
  const { data: res } = await rolePage({
    page: 1,
    pageSize: 1000, // 获取所有角色
    name: '',
    code: '',
    sortOrder: 'asc',
  })
  if (res.code !== 200) return
  roleList.value = res.data?.list || []
}

// 重置查询表单
const reset = () => {
  queryFormRef.value?.resetFields()
  getUserList()
}

// 获取用户列表
const getUserList = async () => {
  const params = {
    ...queryForm.value,
    page: pagination.value.page,
    pageSize: pagination.value.pageSize,
  }
  const { data: res } = await userPage(params)
  if (res.code !== 200) return
  userList.value = res.data?.list || []
  pagination.value.total = res.data?.total || 0
}

// 表格选择变化
const tableSelectionChange = (selection: IUserItem[]) => {
  deleteUserIds.value = selection.map((item) => item.id)
}

// 表格排序变化
const tableSortChange = ({ order }: { order: 'ascending' | 'descending' | null }) => {
  queryForm.value.sortOrder = order === 'ascending' ? 'asc' : 'desc'
  getUserList()
}

// 删除用户
const deleteUserHandle = async (ids: string[]) => {
  const { data: res } = await deleteUser(ids)
  if (res.code !== 200) return
  ElMessage.success('删除成功')
  getUserList()
}

// 刷新
const refresh = (type: 'create' | 'update') => {
  pagination.value.page = type === 'create' ? 1 : pagination.value.page
  // 如果排序为升序，则计算最后一页
  if (queryForm.value.sortOrder === 'asc' && type === 'create') {
    pagination.value.page = PAGINATION_CONFIG.calculateLastPage(
      pagination.value.total + 1,
      pagination.value.pageSize,
    )
  }
  getUserList()
}

onMounted(() => {
  getRoleList()
  getUserList()
})
</script>

<style lang="scss" scoped></style>
