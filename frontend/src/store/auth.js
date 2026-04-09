import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authService } from '@/services/userService'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || null)
  const isLoading = ref(false)
  const isAuthenticated = ref(!!token.value)

  // Permissions helpers
  const permissions = computed(() => user.value?.permissions ?? [])
  function hasPermission(name) {
    if (user.value?.isSuperAdmin) return true
    return permissions.value.includes(name)
  }
  function can(mod) {
    if (user.value?.isSuperAdmin) return true
    return permissions.value.includes(mod)
  }
  function canAny(...mods) {
    if (user.value?.isSuperAdmin) return true
    return mods.some(m => permissions.value.includes(m))
  }
  function canPrefix(prefix) {
    if (user.value?.isSuperAdmin) return true
    return permissions.value.some(p => p.startsWith(prefix + '.'))
  }

  // Actions
  async function login(credentials) {
    isLoading.value = true
    try {
      const response = await authService.login(credentials)
      // Backend wraps the payload: { success, data: { token, user } }
      const { token: jwt, user: userData } = response.data ?? response

      token.value = jwt
      user.value = userData
      isAuthenticated.value = true

      localStorage.setItem('token', jwt)
      if (userData) localStorage.setItem('user', JSON.stringify(userData))
    } catch (error) {
      // Re-throw so callers (Login.vue) can show a proper error message
      throw error
    } finally {
      isLoading.value = false
    }
  }

  async function logout() {
    try {
      // TODO: Call logout endpoint
      await authService.logout()
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      // Clear state regardless of API call result
      token.value = null
      user.value = null
      isAuthenticated.value = false
      localStorage.removeItem('token')
    }
  }

  async function getCurrentUser() {
    if (!token.value) return

    try {
      const response = await authService.getCurrentUser()
      // Backend wraps payload: { success, data: { user } }
      user.value = response.data ?? response
    } catch (error) {
      console.error('Error fetching current user:', error)
      logout()
    }
  }

  // Initialize store — restore session from localStorage
  function initialize() {
    const savedUser = localStorage.getItem('user')
    if (savedUser) {
      try { user.value = JSON.parse(savedUser) } catch { /* ignore */ }
    }
    if (token.value) {
      getCurrentUser()
    }
  }

  return {
    // State
    user,
    token,
    isLoading,
    isAuthenticated,
    permissions,
    // Actions
    login,
    logout,
    getCurrentUser,
    initialize,
    hasPermission,
    can,
    canAny,
    canPrefix
  }
})