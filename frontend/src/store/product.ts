import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as productApi from '@/api/product'
import type { CreateProductRequest, GetProductResponse, UpdateProductRequest } from '@/types'

export const useProductStore = defineStore('product', () => {
  const currentProduct = ref<GetProductResponse | null>(null)
  const fetchingLoading = ref(false)
  const submittingLoading = ref(false)

  async function fetchProduct(id: number) {
    fetchingLoading.value = true
    try {
      const res = await productApi.getProduct(id)
      currentProduct.value = res.data.data
      return res.data.data
    } finally {
      fetchingLoading.value = false
    }
  }

  async function createProduct(data: CreateProductRequest) {
    submittingLoading.value = true
    try {
      const res = await productApi.createProduct(data)
      return res.data.data
    } finally {
      submittingLoading.value = false
    }
  }

  async function updateProduct(id: number, data: Omit<UpdateProductRequest, 'id'>) {
    submittingLoading.value = true
    try {
      await productApi.updateProduct(id, data)
    } finally {
      submittingLoading.value = false
    }
  }

  async function deleteProduct(id: number) {
    submittingLoading.value = true
    try {
      await productApi.deleteProduct(id)
    } finally {
      submittingLoading.value = false
    }
  }

  return {
    currentProduct,
    fetchingLoading,
    submittingLoading,
    fetchProduct,
    createProduct,
    updateProduct,
    deleteProduct,
  }
})
