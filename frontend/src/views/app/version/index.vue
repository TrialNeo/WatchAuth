<template>
  <div>
    <el-card class="card-clear-mb" shadow="never">
      <el-form
        ref="queryFormRef"
        :model="queryForm"
        label-width="auto"
        @keyup.enter="getAppNameList"
      >
        <el-row :gutter="10">
          <el-col :lg="6" :md="12" :sm="12" :xl="6" :xs="24">
            <el-form-item label="应用名称" prop="status">
              <el-select v-model="queryForm.appId" filterable placeholder="请选择">
                <el-option v-for="item in appNames" :key="item.appId" :label="item.appName" :value="item.appId" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :lg="6" :md="12" :sm="12" :xl="6" :xs="24">
            <el-form-item>
              <el-button
                :icon="menuStore.iconComponents.Search"
                type="primary"
                @click="getVerInfoList"
                >搜索
              </el-button>
              <el-button :icon="menuStore.iconComponents.Refresh">重置</el-button>
              <el-button type="success" @click="verUpdateRef?.showDialog(queryForm.appId)"
                >应用更新
              </el-button>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>
    <el-card class="card-mt-16" shadow="never">
      <el-table :border="TABLE_CONFIG.border" :data="verList" show-overflow-tooltip>
        <el-table-column :align="TABLE_CONFIG.align" type="selection" width="55" />
        <el-table-column :align="TABLE_CONFIG.align" fixed label="序号" type="index" width="55" />
        <el-table-column
          :align="TABLE_CONFIG.align"
          fixed
          label="应用名称"
          min-width="100"
          prop="appName"
        />
        <el-table-column
          :align="TABLE_CONFIG.align"
          fixed
          label="版本"
          min-width="75"
          prop="version"
        />
        <el-table-column :align="TABLE_CONFIG.align" label="版本描述" min-width="100" prop="desc" />
        <el-table-column :align="TABLE_CONFIG.align" label="sign" min-width="200" prop="sign" />
        <el-table-column :align="TABLE_CONFIG.align" label="强制更新" min-width="75">
          <template #default="{ row }">
            <BaseTag :text="row.status ? '是' : '否'" :type="row.status ? 'success' : 'danger'" />
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
          label="更新时间"
          min-width="180"
          prop="updatedTime"
        />
        <el-table-column :align="TABLE_CONFIG.align" label="补丁地址" min-width="180">
          <template #default="{ row }">
            <a :href="row.patch_url" :text="row.patch_url" />
          </template>
        </el-table-column>
        <el-table-column :align="TABLE_CONFIG.align" fixed="right" label="操作" min-width="130">
          <template #default="{ row }: { row: IAppItem }">
            <el-button
              v-permission="['role:edit']"
              :icon="menuStore.iconComponents.Edit"
              link
              type="primary"
            >
              编辑
            </el-button>
            <el-popconfirm
              :placement="POPCONFIRM_CONFIG.placement"
              :width="POPCONFIRM_CONFIG.width"
              title="确定要删除选中的应用吗？"
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
    </el-card>
    <VerCreate ref="verUpdateRef" />
  </div>
</template>

<script lang="ts" setup>
import {ElMessage, type FormInstance} from 'element-plus'
import {POPCONFIRM_CONFIG, TABLE_CONFIG} from '@/config/elementConfig.ts'
import type {IAppItem} from '@/types/app/app.ts'
import AppVerService from '@/api/version.ts'
import type {IKAppIDVAppName, IVerInfoItem} from '@/types/app/version.ts'
import VerCreate from '@/views/app/version/create.vue'

defineOptions({ name: 'VerView' })

const menuStore = useMenuStore()
const queryFormRef = useTemplateRef<FormInstance>('queryFormRef')
const verUpdateRef = ref<InstanceType<typeof VerCreate> | null>(null)

// 查询表单
const queryForm = ref({
  appId: '',
  appName: '',
})

const appNames = ref<IKAppIDVAppName[]>([])

const getAppNameList = async () => {
  const { data: resp } = await AppVerService.GetAppNames()
  if (resp.code !== 0) {
    ElMessage.error(resp.msg)
    return
  }
  appNames.value = resp.data?.appNames
}
// --------------------------------------------------------
const verList = ref<IVerInfoItem[]>([])

const getVerInfoList = async () => {
  const resp = await AppVerService.GetVerInfoList()
  //验证一下数据
  if (resp.data.code !== 0) {
    ElMessage.error(resp.data.msg)
    return
  }
  verList.value = resp.data.data?.infos || []
}

onMounted(() => {
  getAppNameList()
  getVerInfoList()
})
</script>

<style lang="scss" scoped></style>
