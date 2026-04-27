<template>
  <div>
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>公告管理</span>
          <el-button type="primary" size="default" @click="openCreateDialog()">发布公告</el-button>
        </div>
      </template>
      <el-table :border="TABLE_CONFIG.border" :data="list" show-overflow-tooltip>
        <el-table-column :align="TABLE_CONFIG.align" label="序号" type="index" width="55" />
        <el-table-column :align="TABLE_CONFIG.align" label="标题" min-width="200" prop="title" />
        <el-table-column :align="TABLE_CONFIG.align" label="内容" min-width="300" prop="content" show-overflow-tooltip />
        <el-table-column :align="TABLE_CONFIG.align" label="状态" min-width="80">
          <template #default="{ row }">
            <BaseTag :text="row.status === 0 ? '已发布' : '草稿'" :type="row.status === 0 ? 'success' : 'info'" />
          </template>
        </el-table-column>
        <el-table-column :align="TABLE_CONFIG.align" label="发布时间" min-width="160" prop="createdAt" />
        <el-table-column :align="TABLE_CONFIG.align" fixed="right" label="操作" min-width="140">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="openEditDialog(row)">编辑</el-button>
            <el-popconfirm title="确定删除此公告？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button link type="danger" size="small">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑公告' : '发布公告'" width="600px">
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="80px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="formData.title" placeholder="请输入公告标题" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="4" placeholder="请输入公告内容" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="formData.status">
            <el-radio :value="0">立即发布</el-radio>
            <el-radio :value="1">存为草稿</el-radio>
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
import { getAnnouncementList, createAnnouncement, updateAnnouncement, deleteAnnouncement, type AnnouncementItem } from '@/api/announcement'
import { TABLE_CONFIG } from '@/config/elementConfig'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'

defineOptions({ name: 'AnnouncementView' })

const list = ref<AnnouncementItem[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()

const formData = reactive({
  id: 0,
  title: '',
  content: '',
  status: 0,
})

const formRules: FormRules = {
  title: [{ required: true, message: '请输入公告标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入公告内容', trigger: 'blur' }],
}

const fetchData = async () => {
  const { data: res } = await getAnnouncementList()
  if (res.code !== 0) {
    ElMessage.error(res.msg || '获取公告列表失败')
    return
  }
  list.value = res.data?.list ?? []
}

const openCreateDialog = () => {
  isEdit.value = false
  formData.id = 0
  formData.title = ''
  formData.content = ''
  formData.status = 0
  dialogVisible.value = true
}

const openEditDialog = (row: AnnouncementItem) => {
  isEdit.value = true
  formData.id = row.id
  formData.title = row.title
  formData.content = row.content
  formData.status = row.status
  dialogVisible.value = true
}

const handleSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    if (isEdit.value) {
      const { data: res } = await updateAnnouncement({
        id: formData.id,
        title: formData.title,
        content: formData.content,
        status: formData.status,
      })
      if (res.code !== 0) {
        ElMessage.error(res.msg || '更新失败')
        return
      }
      ElMessage.success('更新成功')
    } else {
      const { data: res } = await createAnnouncement({
        title: formData.title,
        content: formData.content,
        status: formData.status,
      })
      if (res.code !== 0) {
        ElMessage.error(res.msg || '发布失败')
        return
      }
      ElMessage.success('发布成功')
    }
    dialogVisible.value = false
    await fetchData()
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (id: number) => {
  const { data: res } = await deleteAnnouncement(id)
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
