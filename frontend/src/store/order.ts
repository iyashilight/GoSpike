import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as orderApi from '@/api/order'
import type { CreateOrderRequest, GetOrderResponse } from '@/types'

export const useOrderStore = defineStore('order', () => {
  const currentOrder = ref<GetOrderResponse | null>(null)
  const fetchingLoading = ref(false)
  const submittingLoading = ref(false)

  async function createOrder(data: CreateOrderRequest) {
    submittingLoading.value = true
    try {
      const res = await orderApi.createOrder(data)
      return res.data.data
    } finally {
      submittingLoading.value = false
    }
  }

  async function fetchOrder(orderId: number) {
    fetchingLoading.value = true
    try {
      const res = await orderApi.getOrder(orderId)
      currentOrder.value = res.data.data
      return res.data.data
    } finally {
      fetchingLoading.value = false
    }
  }

  return {
    currentOrder,
    fetchingLoading,
    submittingLoading,
    createOrder,
    fetchOrder,
  }
})
