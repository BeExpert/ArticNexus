import axios from 'axios'
import router from '@/router'

// Create axios instance with base configuration
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  }
})

// Request interceptor — attach JWT token when present
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor — handle 401 globally (skip login endpoint so it can show proper error)
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const isLoginEndpoint = error.config?.url?.includes('/auth/login')
    if (error.response?.status === 401 && !isLoginEndpoint) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      // Lazy import to avoid circular dependency (api → auth → userService → api)
      const { useAuthStore } = await import('@/store/auth')
      const authStore = useAuthStore()
      authStore.user = null
      authStore.token = null
      authStore.isAuthenticated = false
      router.push({ name: 'login' })
    }
    return Promise.reject(error)
  }
)

export default api