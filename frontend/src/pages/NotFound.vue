<template>
  <div class="min-h-screen flex items-center justify-center bg-[var(--color-bg-secondary)]">
    <div class="text-center">
      <p class="text-6xl font-bold text-[var(--color-text-secondary)]">404</p>
      <p class="text-[var(--color-text-secondary)] mt-4">页面不存在</p>
      <router-link to="/dashboard">
        <a-button type="primary" class="mt-6">返回首页</a-button>
      </router-link>
    </div>
  </div>
</template>
