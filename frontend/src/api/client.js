import axios from 'axios'

const BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

const client = axios.create({ baseURL: BASE_URL })

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// If an access token expires mid-session, transparently refresh it once
// and retry the original request before giving up.
client.interceptors.response.use(
  (response) => response,
  async (error) => {
    const original = error.config
    if (error.response?.status === 401 && !original._retried) {
      original._retried = true
      const refreshToken = localStorage.getItem('refresh_token')
      if (refreshToken) {
        try {
          const { data } = await axios.post(`${BASE_URL}/api/refresh`, { refresh_token: refreshToken })
          localStorage.setItem('access_token', data.access_token)
          original.headers.Authorization = `Bearer ${data.access_token}`
          return client(original)
        } catch (refreshErr) {
          localStorage.removeItem('access_token')
          localStorage.removeItem('refresh_token')
        }
      }
    }
    return Promise.reject(error)
  }
)

export const api = {
  login(username, password) {
    return client.post('/api/login', { username, password })
  },
  searchSchedules(params) {
    return client.get('/api/schedules', { params })
  },
  viewSeats(scheduleId) {
    return client.get(`/api/schedules/${scheduleId}/seats`)
  },
  lockSeat(seatId) {
    return client.post(`/api/seats/${seatId}/lock`)
  },
  unlockSeat(seatId) {
    return client.post(`/api/seats/${seatId}/unlock`)
  },
  confirmBooking(seatId) {
    return client.post(`/api/seats/${seatId}/confirm`)
  },
}

export default client
