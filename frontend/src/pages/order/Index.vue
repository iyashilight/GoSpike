<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-semibold">订单管理</h2>
    </div>

    <div class="bg-[var(--color-bg-primary)] rounded-lg border border-[var(--color-border)]">
      <div class="p-4 border-b border-[var(--color-border)]">
        <a-space>
          <a-input
            v-model="searchId"
            placeholder="输入订单ID查询"
            style="width: 200px"
            allow-clear
          />
          <a-button type="primary" @click="handleSearch" :loading="loading">
            查询
          </a-button>
        </a-space>
      </div>

      <div v-if="order" class="p-6">
        <a-descriptions :column="2" bordered>
          <a-descriptions-item label="订单ID">{{ order.orderId }}</a-descriptions-item>
          <a-descriptions-item label="用户ID">{{ order.userId }}</a-descriptions-item>
          <a-descriptions-item label="总金额">
            <span class="text-red-500 font-semibold">¥{{ order.totalPrice }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="getStatusColor(order.status)">
              {{ getStatusLabel(order.status) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="创建时间" :span="2">{{ order.createdAt }}</a-descriptions-item>
        </a-descriptions>

        <h3 class="text-base font-semibold mt-6 mb-3">订单商品</h3>
        <a-table :data="order.items" :pagination="false" :bordered="{ cell: true }">
          <template #columns>
            <a-table-column title="商品ID" data-index="productId" />
            <a-table-column title="商品名称" data-index="productName" />
            <a-table-column title="单价">
              <template #cell="{ record }">¥{{ record.price }}</template>
            </a-table-column>
            <a-table-column title="数量" data-index="quantity" />
            <a-table-column title="小计">
              <template #cell="{ record }">
                ¥{{ (record.price * record.quantity).toFixed(2) }}
              </template>
            </a-table-column>
          </template>
        </a-table>
      </div>

      <div v-else-if="!loading" class="p-12 text-center">
        <a-empty description="请输入订单ID进行查询" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Message } from '@arco-design/web-vue'
import { useOrderStore } from '@/store/order'
import { ORDER_STATUS_MAP } from '@/constants'
import type { GetOrderResponse } from '@/types'

const orderStore = useOrderStore()

const searchId = ref('')
const loading = ref(false)
const order = ref<GetOrderResponse | null>(null)

function getStatusLabel(status: string) {
  return ORDER_STATUS_MAP[status]?.label || status
}

function getStatusColor(status: string) {
  return ORDER_STATUS_MAP[status]?.color || 'gray'
}

async function handleSearch() {
  const id = Number(searchId.value)
  if (!id || id <= 0) {
    Message.warning('请输入有效的订单ID')
    return
  }

  loading.value = true
  try {
    order.value = await orderStore.fetchOrder(id)
  } catch {
    order.value = null
  } finally {
    loading.value = false
  }
}
</script>
