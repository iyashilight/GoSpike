import request from '@/utils/request'
import type {
  ApiResponse,
  CreateProductRequest,
  CreateProductResponse,
  GetProductResponse,
  UpdateProductRequest,
  DeleteProductRequest,
} from '@/types'

export function createProduct(data: CreateProductRequest) {
  return request.post<ApiResponse<CreateProductResponse>>('/product/create', data)
}

export function getProduct(id: number) {
  return request.get<ApiResponse<GetProductResponse>>(`/product/${id}`)
}

export function updateProduct(id: number, data: Omit<UpdateProductRequest, 'id'>) {
  return request.put<ApiResponse<void>>(`/product/${id}`, data)
}

export function deleteProduct(id: number) {
  return request.delete<ApiResponse<void>>(`/product/${id}`)
}
