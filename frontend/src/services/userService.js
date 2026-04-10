import api from './api'

// User service for API communication
export const userService = {
  // TODO: Implement user CRUD operations
  async getUsers(params = {}) {
    try {
      const response = await api.get('/users', { params })
      return response.data
    } catch (error) {
      throw new Error('Error fetching users')
    }
  },

  async getUser(id) {
    try {
      const response = await api.get(`/users/${id}`)
      return response.data
    } catch (error) {
      throw new Error('Error fetching user')
    }
  },

  async createUser(userData) {
    try {
      const response = await api.post('/users', userData)
      return response.data
    } catch (error) {
      throw new Error('Error creating user')
    }
  },

  async updateUser(id, userData) {
    try {
      const response = await api.put(`/users/${id}`, userData)
      return response.data
    } catch (error) {
      throw new Error('Error updating user')
    }
  },

  async deleteUser(id) {
    try {
      await api.delete(`/users/${id}`)
      return true
    } catch (error) {
      throw new Error('Error deleting user')
    }
  }
}

// Authentication service
export const authService = {
  async login(credentials) {
    const response = await api.post('/auth/login', credentials)
    return response.data
  },

  async logout() {
    try {
      await api.post('/auth/logout')
    } catch {
      // Ignore logout endpoint errors — we clear the session client-side anyway
    }
    return true
  },

  async getCurrentUser() {
    const response = await api.get('/auth/me')
    return response.data
  },

  async getMyCompanies() {
    const response = await api.get('/auth/me/companies')
    return response.data
  },

  async selectCompany(companyId) {
    const response = await api.post('/auth/select-company', { companyId })
    return response.data
  }
}