import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { apiFetch } from '../api/client'

export interface User {
  id: string
  tenant_id: string
  email: string
  full_name: string
  role: 'admin' | 'hr_admin' | 'payroll_manager' | 'compliance_officer' | 'auditor'
  permissions: string[]
}

interface LoginResponse {
  data: {
    access_token: string
    expires_in: number
    user: User
  }
  meta: Record<string, unknown>
  success: boolean
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('complylite_token'))
  const user = ref<User | null>(null)
  const isAuthenticated = computed(() => Boolean(token.value && user.value))

  async function login(email: string, password: string) {
    const result = await apiFetch<LoginResponse>('/auth/login', {
      method: 'POST', body: JSON.stringify({ email, password }),
    })
    token.value = result.data.access_token
    user.value = result.data.user
    localStorage.setItem('complylite_token', result.data.access_token)
  }

  async function loadUser() {
    if (!token.value) return
    const result = await apiFetch<{ data: User; meta: Record<string, unknown>; success: boolean }>('/auth/me')
    user.value = result.data
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('complylite_token')
  }

  window.addEventListener('auth:expired', logout)
  return { token, user, isAuthenticated, login, loadUser, logout }
})
