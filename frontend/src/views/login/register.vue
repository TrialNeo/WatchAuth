<template>
  <div class="form-content-inner">
    <h2 class="title">创建账号</h2>
    <p class="subtitle">加入 {{ APP_CONFIG.name }}，开始您的管理之旅</p>

    <el-form ref="formRef" :model="form" :rules="rules" class="register-form" label-position="top">
      <el-form-item prop="username">
        <el-input v-model="form.username" placeholder="设置用户名" />
      </el-form-item>
      <el-form-item prop="email">
        <el-input v-model="form.email" placeholder="输入电子邮箱" />
      </el-form-item>
      <el-form-item prop="password">
        <el-input
          v-model="form.password"
          placeholder="设置登录密码"
          show-password
          type="password"
        />
      </el-form-item>
      <el-form-item prop="confirmPassword">
        <el-input
          v-model="form.confirmPassword"
          placeholder="确认您的密码"
          show-password
          type="password"
        />
      </el-form-item>
      <el-button :loading="loading" class="submit-btn" type="primary" @click="handleRegister">
        立即注册
      </el-button>
      <div class="back-link">
        <span class="have-account">已有账号？</span>
        <el-link :underline="false" @click="emits('goToMode', 'login')">
          <el-icon><component :is="menuStore.iconComponents['Element:ArrowLeft']" /></el-icon>
          返回登录
        </el-link>
      </div>
    </el-form>
  </div>
</template>

<script lang="ts" setup>
import { APP_CONFIG } from '@/config/app.config'
import { userRegister } from '@/api/login'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import type { IEmits } from '@/types/login'

defineOptions({ name: 'RegisterComponent' })

const emits = defineEmits<IEmits>()
const menuStore = useMenuStore()
const formRef = useTemplateRef<FormInstance>('formRef')
const loading = ref(false)

const form = ref({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
})

const rules = reactive<FormRules>({
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '邮箱格式不正确', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请设置密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (_rule: unknown, value: string, callback: (e?: Error) => void) => {
        if (value !== form.value.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
})

const handleRegister = async () => {
  await formRef.value?.validate()
  loading.value = true
  try {
    const { data: res } = await userRegister({
      username: form.value.username,
      password: form.value.password,
      email: form.value.email,
    })
    if (res.code !== 0) {
      ElMessage.error(res.msg || '注册失败')
      return
    }
    ElMessage.success('注册成功，请登录')
    emits('goToMode', 'login')
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
.form-content-inner {
  .title {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--el-text-color-primary);
    margin-bottom: 0.5rem;
  }

  .subtitle {
    font-size: 0.95rem;
    color: var(--el-text-color-secondary);
    margin-bottom: 2rem;
  }

  .register-form {
    :deep(.el-input__wrapper),
    :deep(.el-select__wrapper) {
      padding: 0.5rem 1rem;
      border-radius: 0.5rem;
      box-shadow: 0 0 0 1px var(--el-border-color) inset;
      min-height: 2.75rem;

      &.is-focus {
        box-shadow: 0 0 0 1px var(--el-color-primary) inset;
      }
    }

    .submit-btn {
      width: 100%;
      height: 2.75rem;
      border-radius: 0.75rem;
      font-size: 1rem;
      font-weight: 600;
      margin-top: 0.9rem;
      margin-bottom: 1.5rem;
    }

    .back-link {
      display: flex;
      justify-content: center;
      align-items: center;
      gap: 0.5rem;

      .have-account {
        font-size: 0.875rem;
        color: var(--el-text-color-secondary);
      }

      .el-link {
        font-size: 0.9rem;
        font-weight: 600;
        transition: all 0.3s;
        color: var(--el-text-color-secondary);

        &:hover {
          color: var(--el-color-primary);
          transform: translateX(-4px);
        }
      }
    }
  }
}
</style>
