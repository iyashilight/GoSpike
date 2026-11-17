import request from '@/utils/request'
import type {
  ApiResponse,
  CacheStockRequest,
  CacheStockResponse,
  FlashPurchaseRequest,
  FlashPurchaseResponse,
} from '@/types'

export function cacheStock(data: CacheStockRequest) {
  return request.post<ApiResponse<CacheStockResponse>>('/flashsale/cache-stock', data)
}

export function flashPurchase(data: FlashPurchaseRequest) {
  return request.post<ApiResponse<FlashPurchaseResponse>>('/flashsale/purchase', data)
}
