export const useApi = () => {
  const config = useRuntimeConfig()
  const authStore = useAuthStore()

  const fetchApi = async <T = any>(endpoint: string, options: any = {}) => {
    const headers: Record<string, string> = {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      ...options.headers,
    }

    if (authStore.token) {
      headers['Authorization'] = `Bearer ${authStore.token}`
    }

    const baseUrl = config.public.apiBase || 'http://localhost:8000/api'
    const url = `${baseUrl}${endpoint.startsWith('/') ? endpoint : '/' + endpoint}`

    try {
      const data = await $fetch<T>(url, {
        ...options,
        headers,
      })
      return { data, error: null }
    } catch (err: any) {
      return {
        data: null,
        error: err.data || { message: err.message || 'Terjadi kesalahan koneksi' },
        status: err.status || err.statusCode || 500,
      }
    }
  }

  return { fetchApi }
}
