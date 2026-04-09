<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50">
    <div class="w-full max-w-md bg-white rounded-lg shadow-md p-8">
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-gray-900">ArticNexus</h1>
        <p class="text-gray-500 mt-2">Recuperar contraseña</p>
      </div>

      <!-- Success state -->
      <div v-if="sent" class="text-center">
        <div class="rounded-md bg-green-50 border border-green-200 px-4 py-4 mb-6">
          <p class="text-sm text-green-700">
            Si el usuario existe y tiene un correo asociado, recibirás un enlace para restablecer tu contraseña.
            Revisa tu bandeja de entrada y carpeta de spam.
          </p>
        </div>
        <RouterLink to="/login" class="text-sm text-indigo-600 hover:text-indigo-500 transition-colors">
          Volver al inicio de sesión
        </RouterLink>
      </div>

      <!-- Form state -->
      <form v-else @submit.prevent="handleSubmit" class="space-y-5">
        <div v-if="errorMsg" class="rounded-md bg-red-50 border border-red-200 px-4 py-3 text-sm text-red-700">
          {{ errorMsg }}
        </div>

        <p class="text-sm text-slate-600">
          Ingresa tu nombre de usuario y te enviaremos un enlace al correo asociado para restablecer tu contraseña.
        </p>

        <div>
          <label for="username" class="block text-sm font-medium text-gray-700 mb-1">Nombre de usuario</label>
          <input
            id="username"
            v-model="username"
            type="text"
            autocomplete="username"
            required
            :disabled="isLoading"
            class="block w-full rounded-md border border-gray-300 px-3 py-2 text-gray-900 placeholder-gray-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500 disabled:opacity-50"
            placeholder="admin"
          />
        </div>

        <button
          type="submit"
          :disabled="isLoading"
          class="w-full flex justify-center rounded-md bg-indigo-600 px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          <span v-if="isLoading">Enviando…</span>
          <span v-else>Enviar enlace</span>
        </button>

        <div class="text-center">
          <RouterLink to="/login" class="text-sm text-indigo-600 hover:text-indigo-500 transition-colors">
            Volver al inicio de sesión
          </RouterLink>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import api from '@/services/api'
import { teError } from '@/i18n'

const username = ref('')
const isLoading = ref(false)
const errorMsg = ref('')
const sent = ref(false)

async function handleSubmit() {
  isLoading.value = true
  errorMsg.value = ''
  try {
    await api.post('/auth/forgot-password', { username: username.value })
    sent.value = true
  } catch (err) {
    errorMsg.value = teError(err)
  } finally {
    isLoading.value = false
  }
}
</script>
