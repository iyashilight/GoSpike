<template>
  <div class="space-y-6">
    <div class="flex items-center gap-3">
      <router-link to="/product">
        <a-button type="text">
          <template #icon><icon-arrow-left /></template>
          返回
        </a-button>
      </router-link>
      <h2 class="text-xl font-semibold">编辑商品</h2>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <a-spin :size="36" />
    </div>

    <div v-else-if="product" class="bg-[var(--color-bg-primary)] rounded-lg border border-[var(--color-border)] p-6 max-w-xl">
      <a-form :model="form" @submit="handleSubmit" layout="vertical">
        <a-form-item field="name" label="商品名称" :rules="[{ required: true, message: '请输入商品名称' }]">
          <a-input v-model="form.name" placeholder="请输入商品名称" />
        </a-form-item>

        <a-form-item field="description" label="商品描述" :rules="[{ required: true, message: '请输入商品描述' }]">
          <a-textarea v-model="form.description" placeholder="请输入商品描述" :max-length="255" show-word-limit />
        </a-form-item>

        <a-form-item field="price" label="价格" :rules="[{ required: true, message: '请输入价格' }]">
          <a-input-number
            v-model="form.price"
            placeholder="请输入价格"
            :min="0.01"
            :precision="2"
            :step="0.01"
          >
            <template #prefix>¥</template>
          </a-input-number>
        </a-form-item>

        <a-form-item field="stock" label="库存" :rules="[{ required: true, message: '请输入库存' }]">
          <a-input-number
            v-model="form.stock"
            placeholder="请输入库存"
            :min="0"
          />
        </a-form-item>

        <a-form-item>
          <a-space>
            <a-button type="primary" html-type="submit" :loading="submitting">
              保存
            </a-button>
            <router-link to="/product">
              <a-button>取消</a-button>
            </router-link>
          </a-space>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { useProductStore } from '@/store/product'
import type { GetProductResponse } from '@/types'

const route = useRoute()
const router = useRouter()
const productStore = useProductStore()

const loading = ref(true)
const submitting = ref(false)
const product = ref<GetProductResponse | null>(null)
const form = reactive({
  name: '',
  description: '',
  price: 0,
  stock: 0,
})

onMounted(async () => {
  const id = Number(route.params.id)
  if (!id) {
    Message.error('无效的商品ID')
    router.push('/product')
    return
  }

  try {
    const data = await productStore.fetchProduct(id)
    product.value = data
    form.name = data.name
    form.description = data.description
    form.price = data.price
    form.stock = data.stock
  } catch {
    Message.error('获取商品信息失败')
    router.push('/product')
  } finally {
    loading.value = false
  }
})

async function handleSubmit() {
  if (!product.value) return

  if (!form.name || !form.description || form.price <= 0) {
    Message.warning('请填写完整信息')
    return
  }

  submitting.value = true
  try {
    await productStore.updateProduct(product.value.id, form)
    Message.success('更新成功')
    router.push('/product')
  } catch {
    // Error handled by interceptor
  } finally {
    submitting.value = false
  }
}
</script>
