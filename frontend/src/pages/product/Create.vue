<template>
  <div class="space-y-6">
    <div class="flex items-center gap-3">
      <router-link to="/product">
        <a-button type="text">
          <template #icon><icon-arrow-left /></template>
          返回
        </a-button>
      </router-link>
      <h2 class="text-xl font-semibold">创建商品</h2>
    </div>

    <div class="bg-[var(--color-bg-primary)] rounded-lg border border-[var(--color-border)] p-6 max-w-xl">
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
            <a-button type="primary" html-type="submit" :loading="loading">
              创建
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
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { useProductStore } from '@/store/product'

const router = useRouter()
const productStore = useProductStore()

const loading = ref(false)
const form = reactive({
  name: '',
  description: '',
  price: 0,
  stock: 0,
})

async function handleSubmit() {
  if (!form.name || !form.description || form.price <= 0) {
    Message.warning('请填写完整信息')
    return
  }

  loading.value = true
  try {
    await productStore.createProduct(form)
    Message.success('创建成功')
    router.push('/product')
  } catch {
    // Error handled by interceptor
  } finally {
    loading.value = false
  }
}
</script>
