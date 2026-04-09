<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50">
    <div class="w-full max-w-md bg-white rounded-lg shadow-md p-8">
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-gray-900">ArticNexus</h1>
        <p class="text-gray-500 mt-2">Restablecer contraseña</p>
      </div>

      <!-- Success state -->
      <div v-if="success" class="text-center">
        <div class="rounded-md bg-green-50 border border-green-200 px-4 py-4 mb-6">
          <p class="text-sm text-green-700">
            Tu contraseña ha sido restablecida correctamente.
          </p>
        </div>
        <RouterLink to="/login" class="text-sm text-indigo-600 hover:text-indigo-500 transition-colors">
          Ir al inicio de sesión
        </RouterLink>
      </div>

      <!-- No token -->
      <div v-else-if="!token" class="text-center">
        <div class="rounded-md bg-red-50 border border-red-200 px-4 py-4 mb-6">
          <p class="text-sm text-red-700">
            {{ te('auth.reset_link_invalid') }}
          </p>
        </div>
        <RouterLink to="/forgot-password" class="text-sm text-indigo-600 hover:text-indigo-500 transition-colors">
          Solicitar nuevo enlace
        </RouterLink>
      </div>

      <!-- Form -->
      <form v-else @submit.prevent="handleSubmit" class="space-y-5">
        <div v-if="errorMsg" class="rounded-md bg-red-50 border border-red-200 px-4 py-3 text-sm text-red-700">
          {{ errorMsg }}
        </div>

        <div>
          <label for="newPassword" class="block text-sm font-medium text-gray-700 mb-1">Nueva contraseña</label>
          <input
            id="newPassword"
            v-model="newPassword"
            type="password"
            required
            minlength="8"
            :disabled="isLoading"
            class="block w-full rounded-md border border-gray-300 px-3 py-2 text-gray-900 placeholder-gray-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500 disabled:opacity-50"
            placeholder="Mínimo 8 caracteres"
          />
        </div>

        <div>
          <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-1">Confirmar contraseña</label>
          <input
            id="confirmPassword"
            v-model="confirmPassword"
            type="password"
            required
            minlength="8"
            :disabled="isLoading"
            class="block w-full rounded-md border border-gray-300 px-3 py-2 text-gray-900 placeholder-gray-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500 disabled:opacity-50"
            placeholder="Repite la contraseña"
          />
        </div>

        <button
          type="submit"
          :disabled="isLoading"
          class="w-full flex justify-center rounded-md bg-indigo-600 px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          <span v-if="isLoading">Restableciendo…</span>
          <span v-else>Restablecer contraseña</span>
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/services/api'
import { te, teError } from '@/i18n'

const route = useRoute()
const token = route.query.token || ''

const newPassword = ref('')
const confirmPassword = ref('')
const isLoading = ref(false)
const errorMsg = ref('')
const success = ref(false)

async function handleSubmit() {
  errorMsg.value = ''

  if (newPassword.value !== confirmPassword.value) {
    errorMsg.value = te('form.passwords_mismatch')
    return
  }

  isLoading.value = true
  try {
    await api.post('/auth/reset-password', {
      token,
      newPassword: newPassword.value,
    })
    success.value = true
  } catch (e) {
    errorMsg.value = teError(e)
  } finally {
    isLoading.value = false
  }
}
</script>
