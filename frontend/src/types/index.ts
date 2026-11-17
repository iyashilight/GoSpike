// Unified response envelope from backend
export interface ApiResponse<T = unknown> {
  code: number
  msg: string
  data: T
}

// ===== User Module =====
export interface RegisterRequest {
  name: string
  email: string
  password: string
}

export interface RegisterResponse {
  id: number
}

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  id: number
  name: string
  token: string
}

export interface GetUserRequest {
  id: number
}

export interface GetUserResponse {
  id: number
  name: string
  email: string
}

// ===== Product Module =====
export interface CreateProductRequest {
  name: string
  description: string
  price: number
  stock: number
}

export interface CreateProductResponse {
  id: number
}

export interface GetProductRequest {
  id: number
}

export interface GetProductResponse {
  id: number
  name: string
  description: string
  price: number
  stock: number
  createdAt: string
}

export interface UpdateProductRequest {
  id: number
  name: string
  description: string
  price: number
  stock: number
}

export interface DeleteProductRequest {
  id: number
}

// ===== Order Module =====
export interface OrderItemParam {
  productId: number
  productName: string
  price: number
  quantity: number
}

export interface CreateOrderRequest {
  userId: number
  items: OrderItemParam[]
}

export interface CreateOrderResponse {
  orderId: number
}

export interface GetOrderRequest {
  orderId: number
}

export interface OrderItem {
  id: number
  orderId: number
  productId: number
  productName: string
  price: number
  quantity: number
  createdAt: string
}

export interface GetOrderResponse {
  orderId: number
  userId: number
  totalPrice: number
  status: string
  items: OrderItem[]
  createdAt: string
}

// ===== Flash Sale Module =====
export interface CacheStockRequest {
  productId: number
  stock: number
}

export interface CacheStockResponse {
  success: boolean
}

export interface FlashPurchaseRequest {
  userId: number
  productId: number
  quantity: number
}

export interface FlashPurchaseResponse {
  orderId: number
  status: string
}

// ===== Error Codes =====
export const ErrorCode = {
  OK: 0,
  ErrUserNotFound: 1001,
  ErrInvalidCredential: 1002,
  ErrUserAlreadyExists: 1003,
  ErrProductNotFound: 2001,
  ErrStockNotEnough: 2002,
  ErrCreateProductFailed: 2003,
  ErrUpdateProductFailed: 2004,
  ErrDeleteProductFailed: 2005,
  ErrCacheStockFailed: 2006,
  ErrOrderNotFound: 3001,
  ErrCreateOrderFailed: 3002,
  ErrUnknown: 9999,
} as const

export type ErrorCodeType = (typeof ErrorCode)[keyof typeof ErrorCode]
