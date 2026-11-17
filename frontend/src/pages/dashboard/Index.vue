<template>
  <div class="space-y-6">
    <h2 class="text-xl font-semibold">仪表盘</h2>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-[var(--color-bg-primary)] rounded-lg border border-[var(--color-border)] p-5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-[var(--color-text-secondary)] text-sm">商品总数</p>
            <p class="text-2xl font-bold mt-1">--</p>
          </div>
          <div class="w-10 h-10 rounded-lg bg-blue-50 flex items-center justify-center">
            <icon-shopping-cart class="text-blue-500" />
          </div>
        </div>
      </div>

      <div class="bg-[var(--color-bg-primary)] rounded-lg border border-[var(--color-border)] p-5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-[var(--color-text-secondary)] text-sm">订单总数</p>
            <p class="text-2xl font-bold mt-1">--</p>
          </div>
          <div class="w-10 h-10 rounded-lg bg-green-50 flex items-center justify-center">
            <icon-file class="text-green-500" />
          </div>
        </div>
      </div>

      <div class="bg-[var(--color-bg-primary)] rounded-lg border border-[var(--color-border)] p-5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-[var(--color-text-secondary)] text-sm">秒杀活动</p>
            <p class="text-2xl font-bold mt-1">--</p>
          </div>
          <div class="w-10 h-10 rounded-lg bg-orange-50 flex items-center justify-center">
            <icon-fire class="text-orange-500" />
          </div>
        </div>
      </div>

      <div class="bg-[var(--color-bg-primary)] rounded-lg border border-[var(--color-border)] p-5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-[var(--color-text-secondary)] text-sm">系统状态</p>
            <p class="text-2xl font-bold mt-1 text-green-500">正常</p>
          </div>
          <div class="w-10 h-10 rounded-lg bg-purple-50 flex items-center justify-center">
            <icon-check-circle class="text-purple-500" />
          </div>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="bg-[var(--color-bg-primary)] rounded-lg border border-[var(--color-border)] p-6">
        <h3 class="text-base font-semibold mb-4">快速操作</h3>
        <div class="space-y-3">
          <router-link to="/product/create">
            <a-button type="primary" long>
              <template #icon><icon-plus /></template>
              创建商品
            </a-button>
          </router-link>
          <router-link to="/flashsale">
            <a-button type="outline" status="warning" long>
              <template #icon><icon-fire /></template>
              秒杀管理
            </a-button>
          </router-link>
        </div>
      </div>

      <div class="bg-[var(--color-bg-primary)] rounded-lg border border-[var(--color-border)] p-6">
        <h3 class="text-base font-semibold mb-4">系统信息</h3>
        <div class="space-y-2 text-sm">
          <div class="flex justify-between">
            <span class="text-[var(--color-text-secondary)]">系统名称</span>
            <span>GoSpike 秒杀平台</span>
          </div>
          <div class="flex justify-between">
            <span class="text-[var(--color-text-secondary)]">技术栈</span>
            <span>Go + Vue 3</span>
          </div>
          <div class="flex justify-between">
            <span class="text-[var(--color-text-secondary)]">架构</span>
            <span>微服务</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
</script>
