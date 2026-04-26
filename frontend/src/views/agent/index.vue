<template>
  <div>
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>代理列表</span>
          <el-button type="primary" size="default" @click="openCreateDialog()">新增代理</el-button>
        </div>
      </template>
      <el-table
        :border="TABLE_CONFIG.border"
        :data="treeData"
        row-key="id"
        default-expand-all
        show-overflow-tooltip
      >
        <el-table-column :align="TABLE_CONFIG.align" label="代理名称" min-width="160" prop="name" />
        <el-table-column :align="TABLE_CONFIG.align" label="联系方式" min-width="140" prop="contact" />
        <el-table-column :align="TABLE_CONFIG.align" label="级别" min-width="80">
          <template #default="{ row }">
            {{ levelMap[row.level] ?? '未知' }}
          </template>
        </el-table-column>
        <el-table-column :align="TABLE_CONFIG.align" label="折扣(%)" min-width="90" prop="discount" />
        <el-table-column :align="TABLE_CONFIG.align" label="余额" min-width="120" prop="balance" />
        <el-table-column :align="TABLE_CONFIG.align" label="状态" min-width="80">
          <template #default="{ row }">
            <BaseTag :text="statusMap[row.status]?.text ?? '未知'" :type="statusMap[row.status]?.type ?? 'info'" />
          </template>
        </el-table-column>
        <el-table-column :align="TABLE_CONFIG.align" label="创建时间" min-width="160" prop="createdAt" />
        <el-table-column :align="TABLE_CONFIG.align" fixed="right" label="操作" min-width="180">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="openCreateDialog(row)">新增下级</el-button>
            <el-button link type="primary" size="small" @click="openEditDialog(row)">编辑</el-button>
            <el-popconfirm title="确定删除此代理？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button link type="danger" size="small">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑代理' : '新增代理'" width="500px">
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <el-form-item label="代理名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入代理名称" />
        </el-form-item>
        <el-form-item label="上级代理" v-if="formData.parentId > 0">
          <el-input :model-value="parentName" disabled />
        </el-form-item>
        <el-form-item label="联系方式" prop="contact">
          <el-input v-model="formData.contact" placeholder="请输入联系方式" />
        </el-form-item>
        <el-form-item label="折扣" prop="discount">
          <el-input-number v-model="formData.discount" :min="0" :max="100" :precision="2" /> %
        </el-form-item>
        <el-form-item label="状态" prop="status" v-if="isEdit">
          <el-radio-group v-model="formData.status">
            <el-radio :value="0">正常</el-radio>
            <el-radio :value="1">冻结</el-radio>
            <el-radio :value="2">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { getAgentList, createAgent, updateAgent, deleteAgent, type AgentItem } from '@/api/agent'
import { TABLE_CONFIG } from '@/config/elementConfig'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'

defineOptions({ name: 'AgentView' })

const treeData = ref<AgentItem[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const parentName = ref('')

const levelMap: Record<number, string> = {
  0: '总代',
  1: '一级',
  2: '二级',
  3: '三级',
}

const statusMap: Record<number, { text: string; type: string }> = {
  0: { text: '正常', type: 'success' },
  1: { text: '冻结', type: 'warning' },
  2: { text: '禁用', type: 'danger' },
}

const formData = reactive({
  id: 0,
  name: '',
  contact: '',
  parentId: 0,
  discount: 100,
  status: 0,
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入代理名称', trigger: 'blur' }],
}

const fetchData = async () => {
  const { data: res } = await getAgentList()
  if (res.code !== 0) {
    ElMessage.error(res.msg || '获取代理列表失败')
    return
  }
  treeData.value = res.data.agents || []
}

const openCreateDialog = (parent?: AgentItem) => {
  isEdit.value = false
  formData.id = 0
  formData.name = ''
  formData.contact = ''
  formData.parentId = parent?.id ?? 0
  formData.discount = parent?.discount ?? 100
  formData.status = 0
  parentName.value = parent?.name ?? ''
  dialogVisible.value = true
}

const openEditDialog = (row: AgentItem) => {
  isEdit.value = true
  formData.id = row.id
  formData.name = row.name
  formData.contact = row.contact
  formData.parentId = row.parentId
  formData.discount = row.discount
  formData.status = row.status
  dialogVisible.value = true
}

const handleSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    if (isEdit.value) {
      const { data: res } = await updateAgent({
        id: formData.id,
        name: formData.name,
        contact: formData.contact,
        discount: formData.discount,
        status: formData.status,
      })
      if (res.code !== 0) {
        ElMessage.error(res.msg || '更新失败')
        return
      }
      ElMessage.success('更新成功')
    } else {
      const { data: res } = await createAgent({
        name: formData.name,
        contact: formData.contact,
        parentId: formData.parentId,
        discount: formData.discount,
      })
      if (res.code !== 0) {
        ElMessage.error(res.msg || '创建失败')
        return
      }
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    await fetchData()
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (id: number) => {
  const { data: res } = await deleteAgent(id)
  if (res.code !== 0) {
    ElMessage.error(res.msg || '删除失败')
    return
  }
  ElMessage.success('删除成功')
  await fetchData()
}

onMounted(fetchData)
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
