<template>
  <div class="space-y-6">
    <h2 class="text-xl font-semibold">个人中心</h2>

    <div class="bg-[var(--color-bg-primary)] rounded-lg border border-[var(--color-border)] p-6 max-w-xl">
      <div v-if="loading" class="flex justify-center py-12">
        <a-spin :size="36" />
      </div>

      <div v-else-if="userInfo">
        <a-descriptions :column="1" bordered>
          <a-descriptions-item label="用户ID">{{ userInfo.id }}</a-descriptions-item>
          <a-descriptions-item label="用户名">{{ userInfo.name }}</a-descriptions-item>
          <a-descriptions-item label="邮箱">{{ userInfo.email }}</a-descriptions-item>
        </a-descriptions>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { useAuthStore } from '@/store/auth'
import { getUser } from '@/api/user'
import type { GetUserResponse } from '@/types'

const authStore = useAuthStore()

const loading = ref(true)
const userInfo = ref<GetUserResponse | null>(null)

onMounted(async () => {
  if (!authStore.userId) {
    Message.error('未登录')
    return
  }

  try {
    const res = await getUser(authStore.userId)
    userInfo.value = res.data.data
  } catch {
    Message.error('获取用户信息失败')
  } finally {
    loading.value = false
  }
})
</script>
