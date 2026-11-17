<template>
  <div class="space-y-6">
    <h2 class="text-xl font-semibold">秒杀活动</h2>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="bg-[var(--color-bg-primary)] rounded-lg border border-[var(--color-border)] p-6">
        <h3 class="text-base font-semibold mb-4">缓存库存</h3>
        <a-form :model="cacheForm" @submit="handleCacheStock" layout="vertical">
          <a-form-item field="productId" label="商品ID" :rules="[{ required: true, message: '请输入商品ID' }]">
            <a-input-number v-model="cacheForm.productId" placeholder="请输入商品ID" :min="1" />
          </a-form-item>

          <a-form-item field="stock" label="库存数量" :rules="[{ required: true, message: '请输入库存数量' }]">
            <a-input-number v-model="cacheForm.stock" placeholder="请输入库存数量" :min="1" />
          </a-form-item>

          <a-form-item>
            <a-button type="primary" html-type="submit" :loading="flashsaleStore.cacheLoading">
              缓存库存
            </a-button>
          </a-form-item>
        </a-form>
      </div>

      <div class="bg-[var(--color-bg-primary)] rounded-lg border border-[var(--color-border)] p-6">
        <h3 class="text-base font-semibold mb-4">秒杀抢购</h3>
        <a-form :model="purchaseForm" @submit="handlePurchase" layout="vertical">
          <a-form-item field="productId" label="商品ID" :rules="[{ required: true, message: '请输入商品ID' }]">
            <a-input-number v-model="purchaseForm.productId" placeholder="请输入商品ID" :min="1" />
          </a-form-item>

          <a-form-item field="quantity" label="购买数量" :rules="[{ required: true, message: '请输入购买数量' }]">
            <a-input-number v-model="purchaseForm.quantity" placeholder="请输入购买数量" :min="1" />
          </a-form-item>

          <a-form-item>
            <a-button type="primary" status="warning" html-type="submit" :loading="flashsaleStore.purchaseLoading">
              <template #icon><icon-fire /></template>
              立即抢购
            </a-button>
          </a-form-item>
        </a-form>

        <div v-if="purchaseResult" class="mt-4 p-4 rounded-lg" :class="purchaseResult.status === 'pending' ? 'bg-green-50' : 'bg-red-50'">
          <div class="flex items-center gap-2">
            <icon-check-circle v-if="purchaseResult.status === 'pending'" class="text-green-500" />
            <icon-close-circle v-else class="text-red-500" />
            <span class="font-medium">
              {{ purchaseResult.status === 'pending' ? '抢购成功，订单处理中' : '抢购失败' }}
            </span>
          </div>
          <p class="text-sm text-[var(--color-text-secondary)] mt-1">
            订单ID: {{ purchaseResult.orderId }} | 状态: {{ purchaseResult.status }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { Message } from '@arco-design/web-vue'
import { useFlashsaleStore } from '@/store/flashsale'
import { useAuthStore } from '@/store/auth'

const flashsaleStore = useFlashsaleStore()
const authStore = useAuthStore()

const purchaseResult = ref<{ orderId: number; status: string } | null>(null)

const cacheForm = reactive({
  productId: 0,
  stock: 0,
})

const purchaseForm = reactive({
  productId: 0,
  quantity: 1,
})

async function handleCacheStock() {
  if (!cacheForm.productId || cacheForm.stock <= 0) {
    Message.warning('请填写完整信息')
    return
  }

  try {
    await flashsaleStore.cacheStock(cacheForm.productId, cacheForm.stock)
    Message.success('缓存库存成功')
  } catch {
    // Error handled by interceptor
  }
}

async function handlePurchase() {
  if (!purchaseForm.productId || purchaseForm.quantity <= 0) {
    Message.warning('请填写完整信息')
    return
  }

  if (!authStore.isLoggedIn) {
    Message.error('请先登录')
    return
  }

  purchaseResult.value = null
  try {
    const result = await flashsaleStore.purchase({
      userId: authStore.userId,
      productId: purchaseForm.productId,
      quantity: purchaseForm.quantity,
    })
    purchaseResult.value = result
    if (result.status === 'pending') {
      Message.success('抢购成功！订单正在处理中')
    } else {
      Message.error('抢购失败：' + result.status)
    }
  } catch {
    // Error handled by interceptor
  }
}
</script>
