<template>
  <div>
    <el-card class="card-clear-mb" shadow="never">
      <el-form ref="queryFormRef" :model="queryForm" label-width="auto" @keyup.enter="getRoleList">
        <el-row :gutter="10">
          <el-col :lg="6" :md="12" :sm="12" :xl="6" :xs="24">
            <el-form-item label="应用名称" prop="name">
              <el-input v-model="queryForm.name" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :lg="6" :md="12" :sm="12" :xl="6" :xs="24">
            <el-form-item label="应用ID" prop="code">
              <el-input v-model="queryForm.code" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :lg="6" :md="12" :sm="12" :xl="6" :xs="24">
            <el-form-item label="状态" prop="status">
              <el-select v-model="queryForm.status" placeholder="请选择">
                <el-option :value="1" label="启用"/>
                <el-option :value="0" label="禁用"/>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :lg="6" :md="12" :sm="12" :xl="6" :xs="24">
            <el-form-item>
              <el-button :icon="menuStore.iconComponents.Search" type="primary" @click="getRoleList"
              >搜索
              </el-button
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
          @click="appCreateRef?.showDialog(undefined)"
        >新增应用
        </el-button>
        <el-popconfirm
          :placement="POPCONFIRM_CONFIG.placement"
          :width="POPCONFIRM_CONFIG.width"
          title="确定要删除选中的应用吗？"
          @confirm="deleteAppHandle(deleteAppIDs)"
        >
          <template #reference>
            <el-button
              :disabled="
                !useButtonPermission(['role:delete'], [() => !!deleteAppIDs.length]).value
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
        :data="appList"
        show-overflow-tooltip
        @selection-change="tableSelectionChange"
        @sort-change="tableSortChange"
      >
        <el-table-column :align="TABLE_CONFIG.align" type="selection" width="55"/>
        <el-table-column :align="TABLE_CONFIG.align" fixed label="序号" type="index" width="55"/>
        <el-table-column
          :align="TABLE_CONFIG.align"
          fixed
          label="appid"
          min-width="100"
          prop="appId"
        />
        <el-table-column
          :align="TABLE_CONFIG.align"
          fixed
          label="应用名称"
          min-width="130"
          prop="appName"
        />
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="应用描述"
          min-width="200"
          prop="description"
        />
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="版本"
          min-width="75"
          prop="version"
        />
        <el-table-column :align="TABLE_CONFIG.align" label="计费类型" min-width="90"
                         prop="payType">
          <template #default="{ row }">
            <BaseTag
              :text="feeTypeMap[row.feeType as keyof typeof feeTypeMap]?.text ?? '未知'"
              :type="feeTypeMap[row.feeType as keyof typeof feeTypeMap]?.type ?? 'danger'"
            />
          </template>
        </el-table-column>
        <el-table-column :align="TABLE_CONFIG.align" label="状态" prop="status">
          <template #default="{ row }">
            <BaseTag
              :text="row.status ? '启用' : '禁用'"
              :type="row.status ? 'success' : 'danger'"
            />
          </template>
        </el-table-column>
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="加密方式"
          min-width="180"
          prop="enctype">
          <template #default="{ row }">
            <BaseTag
              :text="encTypeMap[row.feeType as keyof typeof feeTypeMap]?.text ?? '未知'"
              :type="encTypeMap[row.feeType as keyof typeof feeTypeMap]?.type ?? 'danger'"
            />
          </template>
        </el-table-column>
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="SecretKeys"
          min-width="180"
          prop="secretKeys"
        />
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="创建时间"
          min-width="180"
          prop="createdAt"
          sortable="custom"
        />
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="更新时间"
          min-width="180"
          prop="updatedAt"
        />

        <el-table-column :align="TABLE_CONFIG.align" fixed="right" label="操作" min-width="130">
          <template #default="{ row }: { row: IAppItem }">
            <el-button
              v-permission="['role:edit']"
              :icon="menuStore.iconComponents.Edit"
              link
              type="primary"
              @click="appCreateRef?.showDialog(row.appId)"
            >
              编辑
            </el-button>
            <el-popconfirm
              :placement="POPCONFIRM_CONFIG.placement"
              :width="POPCONFIRM_CONFIG.width"
              title="确定要删除选中的应用吗？"
              @confirm="deleteAppHandle([row.appId])"
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

    <AppCreate ref="appCreateRef" @refresh="refresh"/>
  </div>
</template>

<script lang="ts" setup>
import {appPage, deleteApp} from '@/api/app.ts'
import {useButtonPermission} from '@/composables/useButtonPermission'
import {PAGINATION_CONFIG, POPCONFIRM_CONFIG, TABLE_CONFIG} from '@/config/elementConfig'
import AppCreate from '@/views/app/app/create.vue'
import type {IAppItem} from '@/types/app/app.ts'
import type {FormInstance} from 'element-plus'

defineOptions({name: 'AppView'})

const menuStore = useMenuStore()
const queryFormRef = useTemplateRef<FormInstance>('queryFormRef')
const appCreateRef = ref<InstanceType<typeof AppCreate> | null>(null)

// 计费方式map
const feeTypeMap = {
  0: {text: '免费', type: 'success' as const},
  1: {text: '时间', type: 'warning' as const},
  2: {text: '次数', type: 'info' as const}
}

const encTypeMap = {
  0: {text: 'None', type: '' as const},
  1: {text: 'RSA', type: '' as const},
  2: {text: 'AesGcm', type: '' as const}
}

// 删除角色的ids
const deleteAppIDs   = ref<string[]>([])

// 查询表单
const queryForm = ref({
  name: '',
  code: '',
  status: undefined,
  sortOrder: 'desc' as 'asc' | 'desc',
})

// 角色列表
const appList = ref<IAppItem[]>([])

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
  const {data: res} = await appPage(params)
  if (res.code !== 0) return
  appList.value = res.data?.apps || []
  pagination.value.total = res.data?.total || 0
}

// 表格选择变化
const tableSelectionChange = (selection: IAppItem[]) => {
  deleteAppIDs.value = selection.map((item) => item.appId)
}

// 表格排序变化
const tableSortChange = ({order}: { order: 'ascending' | 'descending' | null }) => {
  queryForm.value.sortOrder = order === 'ascending' ? 'asc' : 'desc'
  getRoleList()
}

// 删除角色
const deleteAppHandle = async (names: string[]) => {
  const {data: res} = await deleteApp(names)
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
