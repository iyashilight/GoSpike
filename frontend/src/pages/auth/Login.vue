<template>
  <div class="min-h-screen flex items-center justify-center bg-[var(--color-bg-secondary)]">
    <div class="w-full max-w-md">
      <div class="bg-[var(--color-bg-primary)] rounded-lg shadow-sm border border-[var(--color-border)] p-8">
        <div class="text-center mb-8">
          <h1 class="text-2xl font-bold text-[var(--color-primary)]">GoSpike</h1>
          <p class="text-[var(--color-text-secondary)] mt-2">秒杀电商平台</p>
        </div>

        <a-form :model="form" @submit="handleLogin" layout="vertical">
          <a-form-item field="email" label="邮箱" :rules="[{ required: true, message: '请输入邮箱' }]">
            <a-input
              v-model="form.email"
              placeholder="请输入邮箱"
              size="large"
            >
              <template #prefix><icon-email /></template>
            </a-input>
          </a-form-item>

          <a-form-item field="password" label="密码" :rules="[{ required: true, message: '请输入密码' }]">
            <a-input-password
              v-model="form.password"
              placeholder="请输入密码"
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
              登录
            </a-button>
          </a-form-item>
        </a-form>

        <div class="text-center">
          <span class="text-[var(--color-text-secondary)] text-sm">还没有账号？</span>
          <router-link to="/register" class="text-[var(--color-primary)] text-sm hover:underline">
            立即注册
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { useAuthStore } from '@/store/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const loading = ref(false)
const form = reactive({
  email: '',
  password: '',
})

async function handleLogin() {
  if (!form.email || !form.password) {
    Message.warning('请填写完整信息')
    return
  }

  loading.value = true
  try {
    await authStore.login(form)
    Message.success('登录成功')
    const redirect = (route.query.redirect as string) || '/dashboard'
    router.push(redirect.startsWith('/') && !redirect.startsWith('//') ? redirect : '/dashboard')
  } catch {
    // Error handled by interceptor
  } finally {
    loading.value = false
  }
}
</script>
