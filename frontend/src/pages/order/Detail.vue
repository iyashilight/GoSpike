<template>
  <div class="space-y-6">
    <div class="flex items-center gap-3">
      <router-link to="/order">
        <a-button type="text">
          <template #icon><icon-arrow-left /></template>
          返回
        </a-button>
      </router-link>
      <h2 class="text-xl font-semibold">订单详情</h2>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <a-spin :size="36" />
    </div>

    <div v-else-if="order" class="space-y-6">
      <div class="bg-[var(--color-bg-primary)] rounded-lg border border-[var(--color-border)] p-6">
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
      </div>

      <div class="bg-[var(--color-bg-primary)] rounded-lg border border-[var(--color-border)] p-6">
        <h3 class="text-base font-semibold mb-4">订单商品</h3>
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
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { useOrderStore } from '@/store/order'
import { ORDER_STATUS_MAP } from '@/constants'
import type { GetOrderResponse } from '@/types'

const route = useRoute()
const orderStore = useOrderStore()

const loading = ref(true)
const order = ref<GetOrderResponse | null>(null)

function getStatusLabel(status: string) {
  return ORDER_STATUS_MAP[status]?.label || status
}

function getStatusColor(status: string) {
  return ORDER_STATUS_MAP[status]?.color || 'gray'
}

onMounted(async () => {
  const orderId = Number(route.params.orderId)
  if (!orderId) {
    Message.error('无效的订单ID')
    return
  }

  try {
    order.value = await orderStore.fetchOrder(orderId)
  } catch {
    Message.error('获取订单信息失败')
  } finally {
    loading.value = false
  }
})
</script>
