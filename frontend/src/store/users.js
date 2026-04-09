import { defineStore } from 'pinia'
import { ref } from 'vue'
import { userService } from '@/services/userService'

export const useUserStore = defineStore('users', () => {
  // State
  const users = ref([])
  const selectedUser = ref(null)
  const isLoading = ref(false)
  const error = ref(null)

  // Actions
  async function fetchUsers(params = {}) {
    isLoading.value = true
    error.value = null
    try {
      // TODO: Implement actual API call
      const data = await userService.getUsers(params)
      users.value = data
    } catch (err) {
      error.value = err.message
    } finally {
      isLoading.value = false
    }
  }

  async function fetchUser(id) {
    isLoading.value = true
    error.value = null
    try {
      // TODO: Implement actual API call
      const data = await userService.getUser(id)
      selectedUser.value = data
      return data
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function createUser(userData) {
    isLoading.value = true
    error.value = null
    try {
      // TODO: Implement actual API call
      const newUser = await userService.createUser(userData)
      users.value.push(newUser)
      return newUser
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function updateUser(id, userData) {
    isLoading.value = true
    error.value = null
    try {
      // TODO: Implement actual API call
      const updatedUser = await userService.updateUser(id, userData)
      const index = users.value.findIndex(u => u.id === id)
      if (index !== -1) {
        users.value[index] = updatedUser
      }
      if (selectedUser.value?.id === id) {
        selectedUser.value = updatedUser
      }
      return updatedUser
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function deleteUser(id) {
    isLoading.value = true
    error.value = null
    try {
      // TODO: Implement actual API call
      await userService.deleteUser(id)
      users.value = users.value.filter(u => u.id !== id)
      if (selectedUser.value?.id === id) {
        selectedUser.value = null
      }
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      isLoading.value = false
    }
  }

  return {
    // State
    users,
    selectedUser,
    isLoading,
    error,
    // Actions
    fetchUsers,
    fetchUser,
    createUser,
    updateUser,
    deleteUser
  }
})