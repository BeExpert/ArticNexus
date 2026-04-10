import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authService } from '@/services/userService'

// Decode the payload of a JWT without any signature verification.
// Safe for reading claims client-side (the server already validated the sig).
function parseJwtPayload(jwt) {
  try {
    const base64 = jwt.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
    return JSON.parse(atob(base64))
  } catch {
    return {}
  }
}

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || null)
  const isLoading = ref(false)
  const isAuthenticated = ref(!!token.value)
  // When the JWT has com_id = 0, the user must pick a company before entering.
  const pendingCompanySelection = ref(false)
  const availableCompanies = ref([])

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
      const { token: jwt, user: userData } = response.data ?? response

      token.value = jwt
      user.value = userData
      isAuthenticated.value = true

      localStorage.setItem('token', jwt)
      if (userData) localStorage.setItem('user', JSON.stringify(userData))

      // Detect multi-company users: com_id == 0 means no company selected yet.
      const payload = parseJwtPayload(jwt)
      if (!userData?.isSuperAdmin && (payload.com_id === 0 || payload.com_id == null)) {
        // Fetch the list of companies so the selector can render them.
        const companiesResp = await authService.getMyCompanies()
        const companies = companiesResp.data ?? companiesResp
        if (Array.isArray(companies) && companies.length > 1) {
          availableCompanies.value = companies
          pendingCompanySelection.value = true
          // Return early without redirecting; the caller detects this flag.
          return
        }
        // Edge case: 0 or 1 company — no selection needed.
      }
      pendingCompanySelection.value = false
    } catch (error) {
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // selectCompany exchanges the current JWT for one scoped to the chosen company.
  async function selectCompany(companyId) {
    isLoading.value = true
    try {
      const response = await authService.selectCompany(companyId)
      const { token: jwt, user: userData } = response.data ?? response

      token.value = jwt
      user.value = userData
      localStorage.setItem('token', jwt)
      if (userData) localStorage.setItem('user', JSON.stringify(userData))

      pendingCompanySelection.value = false
      availableCompanies.value = []
    } finally {
      isLoading.value = false
    }
  }

  async function logout() {
    try {
      await authService.logout()
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      token.value = null
      user.value = null
      isAuthenticated.value = false
      pendingCompanySelection.value = false
      availableCompanies.value = []
      localStorage.removeItem('token')
    }
  }

  async function getCurrentUser() {
    if (!token.value) return

    try {
      const response = await authService.getCurrentUser()
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
    pendingCompanySelection,
    availableCompanies,
    permissions,
    // Actions
    login,
    logout,
    selectCompany,
    getCurrentUser,
    initialize,
    hasPermission,
    can,
    canAny,
    canPrefix
  }
})