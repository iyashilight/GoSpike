export const TOKEN_KEY = 'gospike_token'
export const USER_ID_KEY = 'gospike_user_id'
export const USER_NAME_KEY = 'gospike_user_name'

export const ORDER_STATUS = {
  PENDING: 'pending',
  PAID: 'paid',
  FAILED: 'failed',
} as const

export const ORDER_STATUS_MAP: Record<string, { label: string; color: string }> = {
  pending: { label: '待处理', color: 'orange' },
  paid: { label: '已支付', color: 'green' },
  failed: { label: '失败', color: 'red' },
}
