import request from '@/utils/request'
import type {
  ApiResponse,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RegisterResponse,
  GetUserResponse,
} from '@/types'

export function login(data: LoginRequest) {
  return request.post<ApiResponse<LoginResponse>>('/user/login', data)
}

export function register(data: RegisterRequest) {
  return request.post<ApiResponse<RegisterResponse>>('/user/register', data)
}

export function getUser(id: number) {
  return request.get<ApiResponse<GetUserResponse>>(`/user/${id}`)
}
