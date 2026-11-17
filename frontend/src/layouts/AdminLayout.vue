<template>
  <div class="flex h-screen overflow-hidden bg-[var(--color-bg-secondary)]">
    <aside
      class="w-60 flex-shrink-0 bg-[var(--color-bg-primary)] border-r border-[var(--color-border)] flex flex-col"
    >
      <div class="h-14 flex items-center justify-center border-b border-[var(--color-border)]">
        <span class="text-lg font-bold text-[var(--color-primary)]">GoSpike</span>
      </div>
      <nav class="flex-1 py-2 overflow-y-auto">
        <router-link
          v-for="item in menuItems"
          :key="item.path"
          :to="item.path"
          class="flex items-center gap-3 px-5 py-2.5 text-sm transition-colors hover:bg-[var(--color-bg-secondary)]"
          :class="{
            'text-[var(--color-primary)] bg-[var(--color-bg-secondary)] font-medium': isActive(item.path),
            'text-[var(--color-text-secondary)]': !isActive(item.path),
          }"
        >
          <component :is="item.icon" class="w-4 h-4" />
          <span>{{ item.title }}</span>
        </router-link>
      </nav>
      <div class="border-t border-[var(--color-border)] p-4">
        <div class="flex items-center gap-3">
          <a-avatar :size="32" class="bg-[var(--color-primary)]">
            {{ authStore.userName?.charAt(0) || 'U' }}
          </a-avatar>
          <div class="flex-1 min-w-0">
            <div class="text-sm font-medium truncate">{{ authStore.userName }}</div>
          </div>
          <a-button type="text" size="small" @click="handleLogout">
            <template #icon><icon-logout /></template>
          </a-button>
        </div>
      </div>
    </aside>

    <main class="flex-1 flex flex-col overflow-hidden">
      <header class="h-14 flex items-center justify-between px-6 bg-[var(--color-bg-primary)] border-b border-[var(--color-border)]">
        <a-breadcrumb>
          <a-breadcrumb-item v-for="item in breadcrumbs" :key="item">{{ item }}</a-breadcrumb-item>
        </a-breadcrumb>
        <div class="flex items-center gap-3">
          <a-button shape="circle" size="small" @click="toggleDark">
            <template #icon>
              <icon-moon v-if="isDark" />
              <icon-sun v-else />
            </template>
          </a-button>
        </div>
      </header>

      <div class="flex-1 p-6 overflow-y-auto">
        <router-view />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const isDark = ref(localStorage.getItem('gospike_dark') === 'true' || document.documentElement.classList.contains('dark'))

const menuItems = [
  { path: '/dashboard', title: '仪表盘', icon: 'icon-dashboard' },
  { path: '/product', title: '商品管理', icon: 'icon-shopping-cart' },
  { path: '/order', title: '订单管理', icon: 'icon-file' },
  { path: '/flashsale', title: '秒杀活动', icon: 'icon-fire' },
  { path: '/user/profile', title: '个人中心', icon: 'icon-user' },
]

function isActive(path: string) {
  return route.path === path
}

const breadcrumbs = computed(() => {
  const matched = route.matched.filter((item) => item.meta?.title)
  return matched.map((item) => item.meta.title as string)
})

function toggleDark() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark')
  localStorage.setItem('gospike_dark', String(isDark.value))
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}

watch(isDark, (val) => {
  document.documentElement.classList.toggle('dark', val)
})
</script>
