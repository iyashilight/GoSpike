<template>
  <div class="min-h-screen flex items-center justify-center bg-[var(--color-bg-secondary)]">
    <div class="w-full max-w-md">
      <div class="bg-[var(--color-bg-primary)] rounded-lg shadow-sm border border-[var(--color-border)] p-8">
        <div class="text-center mb-8">
          <h1 class="text-2xl font-bold text-[var(--color-primary)]">GoSpike</h1>
          <p class="text-[var(--color-text-secondary)] mt-2">创建新账号</p>
        </div>

        <a-form :model="form" @submit="handleRegister" layout="vertical">
          <a-form-item field="name" label="用户名" :rules="[{ required: true, message: '请输入用户名' }]">
            <a-input
              v-model="form.name"
              placeholder="请输入用户名"
              size="large"
            >
              <template #prefix><icon-user /></template>
            </a-input>
          </a-form-item>

          <a-form-item field="email" label="邮箱" :rules="[{ required: true, message: '请输入邮箱' }]">
            <a-input
              v-model="form.email"
              placeholder="请输入邮箱"
              size="large"
            >
              <template #prefix><icon-email /></template>
            </a-input>
          </a-form-item>

          <a-form-item field="password" label="密码" :rules="passwordRules">
            <a-input-password
              v-model="form.password"
              placeholder="请输入密码（至少8位，含大小写字母和数字）"
              size="large"
            >
              <template #prefix><icon-lock /></template>
            </a-input-password>
          </a-form-item>

          <a-form-item field="confirmPassword" label="确认密码" :rules="confirmRules">
            <a-input-password
              v-model="form.confirmPassword"
              placeholder="请再次输入密码"
              size="large"
            >
              <template #prefix><icon-lock /></template>
            </a-input-password>
          </a-form-item>

          <a-form-item>
            <a-button
              type="primary"
              html-type="submit"
              long
              size="large"
              :loading="loading"
            >
              注册
            </a-button>
          </a-form-item>
        </a-form>

        <div class="text-center">
          <span class="text-[var(--color-text-secondary)] text-sm">已有账号？</span>
          <router-link to="/login" class="text-[var(--color-primary)] text-sm hover:underline">
            立即登录
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { useAuthStore } from '@/store/auth'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const form = reactive({
  name: '',
  email: '',
  password: '',
  confirmPassword: '',
})

const passwordRules = [
  { required: true, message: '请输入密码' },
  { minLength: 8, message: '密码至少8位' },
  {
    validator: (value: string) => {
      const hasUpper = /[A-Z]/.test(value)
      const hasLower = /[a-z]/.test(value)
      const hasDigit = /[0-9]/.test(value)
      if (!hasUpper || !hasLower || !hasDigit) {
        return '密码需包含大小写字母和数字'
      }
      return true
    },
  },
]

const confirmRules = [
  { required: true, message: '请确认密码' },
  {
    validator: (value: string) => {
      if (value !== form.password) {
        return '两次密码不一致'
      }
      return true
    },
  },
]

async function handleRegister() {
  if (!form.name || !form.email || !form.password) {
    Message.warning('请填写完整信息')
    return
  }

  if (form.password !== form.confirmPassword) {
    Message.warning('两次密码不一致')
    return
  }

  loading.value = true
  try {
    await authStore.register({
      name: form.name,
      email: form.email,
      password: form.password,
    })
    Message.success('注册成功，请登录')
    router.push('/login')
  } catch {
    // Error handled by interceptor
  } finally {
    loading.value = false
  }
}
</script>
