import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

// Criar instância axios separada para auth (sem interceptor para evitar loop)
const authAxios = axios.create({
  baseURL: API_URL,
  withCredentials: true,
})

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const accessToken = ref(localStorage.getItem('accessToken') || null)

  const isAuthenticated = computed(() => !!accessToken.value && !!user.value)

  async function login(idToken) {
    try {
      const response = await authAxios.post('/api/auth/google', {
        idToken
      })

      accessToken.value = response.data.data.accessToken
      user.value = response.data.data.user

      localStorage.setItem('accessToken', accessToken.value)

      return { success: true }
    } catch (error) {
      console.error('Login error:', error)
      return {
        success: false,
        error: error.response?.data?.error || 'Erro ao fazer login'
      }
    }
  }

  async function refreshToken() {
    try {
      const response = await authAxios.post('/api/auth/refresh', {})

      accessToken.value = response.data.data.accessToken
      user.value = response.data.data.user
      localStorage.setItem('accessToken', accessToken.value)

      return true
    } catch (error) {
      console.error('Refresh token error:', error)
      logout()
      return false
    }
  }

  async function logout() {
    const tokenToUse = accessToken.value

    // Limpar estado primeiro
    user.value = null
    accessToken.value = null
    localStorage.removeItem('accessToken')

    // Tentar fazer logout no servidor (não crítico)
    try {
      await authAxios.post('/api/auth/logout', {}, {
        headers: {
          Authorization: `Bearer ${tokenToUse}`
        }
      })
    } catch (error) {
      console.error('Logout error:', error)
    }
  }

  async function fetchUser() {
    if (!accessToken.value) {
      return false
    }

    try {
      const response = await authAxios.get('/api/user/me', {
        headers: {
          Authorization: `Bearer ${accessToken.value}`
        }
      })

      user.value = response.data.data
      return true
    } catch (error) {
      console.error('Fetch user error:', error)
      // Não fazer logout em caso de erro - deixar o interceptor da api lidar com isso
      accessToken.value = null
      localStorage.removeItem('accessToken')
      return false
    }
  }

  return {
    user,
    accessToken,
    isAuthenticated,
    login,
    logout,
    refreshToken,
    fetchUser
  }
})
