<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-semibold">商品管理</h2>
      <router-link to="/product/create">
        <a-button type="primary">
          <template #icon><icon-plus /></template>
          创建商品
        </a-button>
      </router-link>
    </div>

    <div class="bg-[var(--color-bg-primary)] rounded-lg border border-[var(--color-border)]">
      <div class="p-4 border-b border-[var(--color-border)]">
        <a-space>
          <a-input
            v-model="searchId"
            placeholder="输入商品ID查询"
            style="width: 200px"
            allow-clear
          />
          <a-button type="primary" @click="handleSearch" :loading="loading">
            查询
          </a-button>
        </a-space>
      </div>

      <div v-if="product" class="p-6">
        <a-descriptions :column="2" bordered>
          <a-descriptions-item label="ID">{{ product.id }}</a-descriptions-item>
          <a-descriptions-item label="名称">{{ product.name }}</a-descriptions-item>
          <a-descriptions-item label="描述" :span="2">{{ product.description }}</a-descriptions-item>
          <a-descriptions-item label="价格">
            <span class="text-red-500 font-semibold">¥{{ product.price }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="库存">
            <a-tag :color="product.stock > 0 ? 'green' : 'red'">
              {{ product.stock > 0 ? `剩余 ${product.stock}` : '已售罄' }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="创建时间" :span="2">{{ product.createdAt }}</a-descriptions-item>
        </a-descriptions>

        <div class="mt-4 flex gap-2">
          <router-link :to="`/product/${product.id}/edit`">
            <a-button type="primary" status="warning">
              <template #icon><icon-edit /></template>
              编辑
            </a-button>
          </router-link>
          <a-popconfirm content="确定要删除此商品吗？" @ok="handleDelete">
            <a-button type="primary" status="danger">
              <template #icon><icon-delete /></template>
              删除
            </a-button>
          </a-popconfirm>
        </div>
      </div>

      <div v-else-if="!loading" class="p-12 text-center">
        <a-empty description="请输入商品ID进行查询" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Message } from '@arco-design/web-vue'
import { useProductStore } from '@/store/product'
import type { GetProductResponse } from '@/types'

const productStore = useProductStore()

const searchId = ref('')
const loading = ref(false)
const product = ref<GetProductResponse | null>(null)

async function handleSearch() {
  const id = Number(searchId.value)
  if (!id || id <= 0) {
    Message.warning('请输入有效的商品ID')
    return
  }

  loading.value = true
  try {
    product.value = await productStore.fetchProduct(id)
  } catch {
    product.value = null
  } finally {
    loading.value = false
  }
}

async function handleDelete() {
  if (!product.value) return

  try {
    await productStore.deleteProduct(product.value.id)
    Message.success('删除成功')
    product.value = null
    searchId.value = ''
  } catch {
    // Error handled by interceptor
  }
}
</script>
