import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { TOKEN_KEY, USER_ID_KEY, USER_NAME_KEY } from '@/constants'
import * as userApi from '@/api/user'
import type { LoginRequest, RegisterRequest } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem(TOKEN_KEY) || '')
  const userId = ref(Number(localStorage.getItem(USER_ID_KEY)) || 0)
  const userName = ref(localStorage.getItem(USER_NAME_KEY) || '')

  const isLoggedIn = computed(() => !!token.value)

  async function login(data: LoginRequest) {
    const res = await userApi.login(data)
    const { id, name, token: newToken } = res.data.data
    token.value = newToken
    userId.value = id
    userName.value = name
    localStorage.setItem(TOKEN_KEY, newToken)
    localStorage.setItem(USER_ID_KEY, String(id))
    localStorage.setItem(USER_NAME_KEY, name)
    return res.data.data
  }

  async function register(data: RegisterRequest) {
    const res = await userApi.register(data)
    return res.data.data
  }

  function logout() {
    token.value = ''
    userId.value = 0
    userName.value = ''
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_ID_KEY)
    localStorage.removeItem(USER_NAME_KEY)
  }

  return {
    token,
    userId,
    userName,
    isLoggedIn,
    login,
    register,
    logout,
  }
})
