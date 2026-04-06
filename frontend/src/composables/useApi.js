import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

// Criar instância do axios
const api = axios.create({
  baseURL: API_URL,
  withCredentials: true,
})

// Request interceptor - adicionar token automaticamente
api.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.accessToken) {
      config.headers.Authorization = `Bearer ${authStore.accessToken}`
    }

    // Definir Content-Type como application/json se não for FormData
    if (!(config.data instanceof FormData)) {
      config.headers['Content-Type'] = 'application/json'
    }

    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor - renovar token automaticamente em caso de 401
let isRefreshing = false
let failedQueue = []

const processQueue = (error, token = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error)
    } else {
      prom.resolve(token)
    }
  })

  failedQueue = []
}

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config

    // Se o erro não for 401 ou já tentamos renovar, rejeitar
    if (error.response?.status !== 401 || originalRequest._retry) {
      return Promise.reject(error)
    }

    if (isRefreshing) {
      // Se já estiver renovando, adicionar à fila
      return new Promise((resolve, reject) => {
        failedQueue.push({ resolve, reject })
      })
        .then((token) => {
          originalRequest.headers.Authorization = `Bearer ${token}`
          return api(originalRequest)
        })
        .catch((err) => {
          return Promise.reject(err)
        })
    }

    originalRequest._retry = true
    isRefreshing = true

    const authStore = useAuthStore()

    try {
      // Tentar renovar o token
      const success = await authStore.refreshToken()

      if (success) {
        processQueue(null, authStore.accessToken)
        originalRequest.headers.Authorization = `Bearer ${authStore.accessToken}`
        return api(originalRequest)
      } else {
        processQueue(error, null)
        router.push({ name: 'login' })
        return Promise.reject(error)
      }
    } catch (refreshError) {
      processQueue(refreshError, null)
      authStore.logout()
      router.push({ name: 'login' })
      return Promise.reject(refreshError)
    } finally {
      isRefreshing = false
    }
  }
)

export function useApi() {
  return api
}
