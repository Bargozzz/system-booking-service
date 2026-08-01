import { defineStore } from 'pinia'

export interface User {
  id: number
  name: string
  email: string
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(null)
  const user = ref<User | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => !!token.value)

  const restoreSession = () => {
    if (process.client) {
      const savedToken = localStorage.getItem('mbs_jwt_token')
      const savedUser = localStorage.getItem('mbs_user')
      if (savedToken && savedUser) {
        token.value = savedToken
        try {
          user.value = JSON.parse(savedUser)
        } catch {
          user.value = null
        }
      }
    }
  }

  const setAuth = (newToken: string, newUser: User) => {
    token.value = newToken
    user.value = newUser
    if (process.client) {
      localStorage.setItem('mbs_jwt_token', newToken)
      localStorage.setItem('mbs_user', JSON.stringify(newUser))
    }
  }

  const logout = () => {
    token.value = null
    user.value = null
    if (process.client) {
      localStorage.removeItem('mbs_jwt_token')
      localStorage.removeItem('mbs_user')
    }
  }

  return {
    token,
    user,
    loading,
    error,
    isAuthenticated,
    restoreSession,
    setAuth,
    logout,
  }
})
