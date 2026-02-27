<template>
  <div>
    <el-card class="card-clear-mb" shadow="never">
      <el-form
        ref="queryFormRef"
        :model="queryForm"
        label-width="auto"
        @keyup.enter="getUserList"
      >
        <el-row :gutter="10">
          <el-col :lg="6" :md="12" :sm="12" :xl="6" :xs="24">
            <el-form-item label="用户名" prop="userName">
              <el-input v-model="queryForm.userName" placeholder="请输入用户名" />
            </el-form-item>
          </el-col>
          <el-col :lg="6" :md="12" :sm="12" :xl="6" :xs="24">
            <el-form-item label="账号" prop="account">
              <el-input v-model="queryForm.account" placeholder="请输入账号" />
            </el-form-item>
          </el-col>
          <el-col :lg="6" :md="12" :sm="12" :xl="6" :xs="24">
            <el-form-item>
              <el-button
                type="primary"
                @click="getUserList"
                >搜索
              </el-button>
              <el-button @click="resetForm">重置</el-button>
              <el-button type="success" @click="handleAddUser">添加用户</el-button>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>
    <el-card class="card-mt-16" shadow="never">
      <el-table :border="TABLE_CONFIG.border" :data="userList" show-overflow-tooltip>
        <el-table-column :align="TABLE_CONFIG.align" type="selection" width="55" />
        <el-table-column :align="TABLE_CONFIG.align" fixed label="序号" type="index" width="55" />
        <el-table-column
          :align="TABLE_CONFIG.align"
          fixed
          label="用户ID"
          min-width="100"
          prop="userID"
        />
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="用户名"
          min-width="100"
          prop="userName"
        />
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="邮箱"
          min-width="150"
          prop="email"
        />
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="权限"
          min-width="80"
        >
          <template #default="{ row }">
            <BaseTag :text="row.privilege === 0 ? '普通用户' : '管理员'" :type="row.privilege === 0 ? 'info' : 'success'" />
          </template>
        </el-table-column>
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="账号"
          min-width="100"
          prop="account"
        />
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="余额"
          min-width="100"
          prop="balance"
        />
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="消费"
          min-width="100"
          prop="spend"
        />
        <el-table-column
          :align="TABLE_CONFIG.align"
          label="创建时间"
          min-width="180"
          prop="createTime"
        />
        <el-table-column :align="TABLE_CONFIG.align" label="操作" min-width="150">
          <template #default="{ row }">
            <el-button size="small" @click="handleEditUser(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleBanUser(row)">封禁</el-button>
            <el-button size="small" type="primary" @click="handleTopUp(row)">充值</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, ref} from 'vue'
import {ElMessage, type FormInstance} from 'element-plus'
import {TABLE_CONFIG} from '@/config/elementConfig.ts'
import {userManager} from '@/api/userManager.ts'
import type {IUserItem} from '@/types/userManager/user.ts'
import BaseTag from '@/components/tag/BaseTag.vue'

defineOptions({ name: 'UserManager' })

const queryFormRef = ref<FormInstance>()

// 查询表单
const queryForm = ref({
  userName: '',
  account: '',
})

const userList = ref<IUserItem[]>([])

const getUserList = async () => {
  try {
    const resp = await userManager.userList()
    if (resp.data.code !== 0) {
      ElMessage.error(resp.data.msg)
      return
    }
    userList.value = resp.data.data?.list || []
  } catch (error) {
    ElMessage.error('获取用户列表失败')
    console.error('获取用户列表失败:', error)
  }
}

const resetForm = () => {
  queryForm.value = {
    userName: '',
    account: '',
  }
  getUserList()
}

const handleAddUser = () => {
  // 实现添加用户逻辑
  ElMessage.info('添加用户功能待实现')
}

const handleEditUser = (row: IUserItem) => {
  // 实现编辑用户逻辑
  ElMessage.info(`编辑用户: ${row.userName}`)
}

const handleBanUser = async (row: IUserItem) => {
  try {
    const resp = await userManager.ban(row.userID)
    if (resp.data.code !== 0) {
      ElMessage.error(resp.data.msg)
      return
    }
    ElMessage.success('用户封禁成功')
    getUserList()
  } catch (error) {
    ElMessage.error('封禁用户失败')
    console.error('封禁用户失败:', error)
  }
}

const handleTopUp = (row: IUserItem) => {
  // 实现充值逻辑
  ElMessage.info(`为用户 ${row.userName} 充值`)
}

onMounted(() => {
  getUserList()
})
</script>

<style lang="scss" scoped></style>
