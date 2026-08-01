import { reactive } from 'vue'

export const authState = reactive({
  username: localStorage.getItem('username') || null,
  userId: localStorage.getItem('user_id') ? Number(localStorage.getItem('user_id')) : null,
  isAuthenticated: !!localStorage.getItem('access_token'),
})

export function setSession(username, userId, accessToken, refreshToken) {
  localStorage.setItem('username', username)
  localStorage.setItem('user_id', String(userId))
  localStorage.setItem('access_token', accessToken)
  localStorage.setItem('refresh_token', refreshToken)
  authState.username = username
  authState.userId = userId
  authState.isAuthenticated = true
}

export function clearSession() {
  localStorage.removeItem('username')
  localStorage.removeItem('user_id')
  localStorage.removeItem('access_token')
  localStorage.removeItem('refresh_token')
  authState.username = null
  authState.userId = null
  authState.isAuthenticated = false
}
