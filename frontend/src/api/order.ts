import request from '@/utils/request'
import type {
  ApiResponse,
  CreateOrderRequest,
  CreateOrderResponse,
  GetOrderResponse,
} from '@/types'

export function createOrder(data: CreateOrderRequest) {
  return request.post<ApiResponse<CreateOrderResponse>>('/order/create', data)
}

export function getOrder(orderId: number) {
  return request.get<ApiResponse<GetOrderResponse>>(`/order/${orderId}`)
}
