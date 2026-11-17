import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as flashsaleApi from '@/api/flashsale'
import type { FlashPurchaseRequest } from '@/types'

export const useFlashsaleStore = defineStore('flashsale', () => {
  const cacheLoading = ref(false)
  const purchaseLoading = ref(false)
  const lastPurchaseResult = ref<{ orderId: number; status: string } | null>(null)

  async function cacheStock(productId: number, stock: number) {
    cacheLoading.value = true
    try {
      await flashsaleApi.cacheStock({ productId, stock })
    } finally {
      cacheLoading.value = false
    }
  }

  async function purchase(data: FlashPurchaseRequest) {
    purchaseLoading.value = true
    try {
      const res = await flashsaleApi.flashPurchase(data)
      lastPurchaseResult.value = res.data.data
      return res.data.data
    } finally {
      purchaseLoading.value = false
    }
  }

  return {
    cacheLoading,
    purchaseLoading,
    lastPurchaseResult,
    cacheStock,
    purchase,
  }
})
