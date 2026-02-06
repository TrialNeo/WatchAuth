<template>
  <div>
    <el-card class="card-clear-mb" shadow="never">
      <el-form ref="queryFormRef" :model="queryForm" label-width="auto" @keyup.enter="getRoleList">
        <el-row :gutter="10">
          <el-col :lg="6" :md="12" :sm="12" :xl="6" :xs="24">
            <el-form-item label="角色名称" prop="name">
              <el-input v-model="queryForm.name" placeholder="请输入" />
            </el-form-item>
          </el-col>
          <el-col :lg="6" :md="12" :sm="12" :xl="6" :xs="24">
            <el-form-item label="角色编码" prop="code">
              <el-input v-model="queryForm.code" placeholder="请输入" />
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
              <el-button :icon="menuStore.iconComponents.Search" type="primary" @click="getRoleList"
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
          v-permission="['role:add']"
          :icon="menuStore.iconComponents.Plus"
          type="primary"
          @click="roleCreateRef?.showDialog(undefined)"
          >新增角色</el-button
        >
        <el-popconfirm
          :placement="POPCONFIRM_CONFIG.placement"
          :width="POPCONFIRM_CONFIG.width"
          title="确定要删除选中的角色吗？"
          @confirm="deleteRoleHandle(deleteRoleIds)"
        >
          <template #reference>
            <el-button
              :disabled="
                !useButtonPermission(['role:delete'], [() => !!deleteRoleIds.length]).value
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
        :data="roleList"
        show-overflow-tooltip
        @selection-change="tableSelectionChange"
        @sort-change="tableSortChange"
      >
        <el-table-column :align="TABLE_CONFIG.align" type="selection" width="55" />
        <el-table-column :align="TABLE_CONFIG.align" fixed label="序号" type="index" width="55" />
        <el-table-column
          :align="TABLE_CONFIG.align"
          fixed
          label="角色名称"
          min-width="160"
          prop="name"
        />
        <el-table-column :align="TABLE_CONFIG.align" label="角色编码" min-width="160" prop="code" />
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="角色描述"
          min-width="200"
          prop="description"
        />
        <el-table-column :align="TABLE_CONFIG.align" label="类型" prop="isBuiltIn">
          <template #default="{ row }">
            <BaseTag
              :text="row.isBuiltIn ? '内置' : '自定义'"
              :type="row.isBuiltIn ? 'warning' : 'success'"
            />
          </template>
        </el-table-column>
        <el-table-column :align="TABLE_CONFIG.align" label="状态" prop="status">
          <template #default="{ row }">
            <BaseTag
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
          <template #default="{ row }: { row: IRoleItem }">
            <el-button
              v-permission="['role:edit']"
              :icon="menuStore.iconComponents.Edit"
              link
              type="primary"
              @click="roleCreateRef?.showDialog(row.id)"
            >
              编辑
            </el-button>
            <el-popconfirm
              :placement="POPCONFIRM_CONFIG.placement"
              :width="POPCONFIRM_CONFIG.width"
              title="确定要删除选中的角色吗？"
              @confirm="deleteRoleHandle([row.id])"
            >
              <template #reference>
                <el-button
                  v-permission="['role:delete']"
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
          @change="getRoleList"
        />
      </div>
    </el-card>

    <RoleCreate ref="roleCreateRef" @refresh="refresh" />
  </div>
</template>

<script lang="ts" setup>
import {deleteRole, rolePage} from '@/api/role'
import {useButtonPermission} from '@/composables/useButtonPermission'
import {PAGINATION_CONFIG, POPCONFIRM_CONFIG, TABLE_CONFIG} from '@/config/elementConfig'
import RoleCreate from '@/views/system/role/create.vue'
import type {IRoleItem} from '@/types/system/role'
import type {FormInstance} from 'element-plus'

defineOptions({ name: 'RoleView' })

const menuStore = useMenuStore()
const queryFormRef = useTemplateRef<FormInstance>('queryFormRef')
const roleCreateRef = useTemplateRef<InstanceType<typeof RoleCreate> | null>('roleCreateRef')

// 删除角色的ids
const deleteRoleIds = ref<string[]>([])

// 查询表单
const queryForm = ref({
  name: '',
  code: '',
  status: undefined,
  sortOrder: 'desc' as 'asc' | 'desc',
})

// 角色列表
const roleList = ref<IRoleItem[]>([])

// 分页
const pagination = ref({
  page: 1,
  pageSize: 10,
  total: 0,
})

// 重置查询表单
const reset = () => {
  queryFormRef.value?.resetFields()
  getRoleList()
}

// 获取角色列表
const getRoleList = async () => {
  const params = {
    ...queryForm.value,
    page: pagination.value.page,
    pageSize: pagination.value.pageSize,
  }
  const { data: res } = await rolePage(params)
  if (res.code !== 0) return
  roleList.value = res.data?.list || []
  pagination.value.total = res.data?.total || 0
}

// 表格选择变化
const tableSelectionChange = (selection: IRoleItem[]) => {
  deleteRoleIds.value = selection.map((item) => item.id)
}

// 表格排序变化
const tableSortChange = ({ order }: { order: 'ascending' | 'descending' | null }) => {
  queryForm.value.sortOrder = order === 'ascending' ? 'asc' : 'desc'
  getRoleList()
}

// 删除角色
const deleteRoleHandle = async (ids: string[]) => {
  const { data: res } = await deleteRole(ids)
  if (res.code !== 0) return
  ElMessage.success('删除成功')
  getRoleList()
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
  getRoleList()
}

onMounted(() => {
  getRoleList()
})
</script>

<style lang="scss" scoped></style>
