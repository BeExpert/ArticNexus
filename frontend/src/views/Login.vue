<template>
  <div class="min-h-screen relative flex items-center justify-center overflow-hidden bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">

    <!-- Decorative circles -->
    <div class="absolute top-[-80px] right-[-80px] w-96 h-96 bg-nordic-cyan/10 rounded-full pointer-events-none"></div>
    <div class="absolute bottom-[-120px] left-[-60px] w-[500px] h-[500px] bg-blue-500/5 rounded-full pointer-events-none"></div>
    <div class="absolute top-1/2 left-[-40px] w-48 h-48 bg-nordic-cyan/5 rounded-full pointer-events-none"></div>

    <!-- Card -->
    <div class="relative z-10 w-full max-w-md mx-4">
      <div class="bg-white/95 backdrop-blur-sm rounded-2xl shadow-2xl p-8 border border-white/20">

        <!-- Header -->
        <div class="text-center mb-8">
          <div class="flex justify-center mb-4">
            <img
              src="@/assets/images/articdev-logo.svg"
              alt="ArticDev"
              class="h-10"
              @error="e => e.target.style.display = 'none'"
            />
          </div>
          <span class="inline-block bg-slate-100 text-slate-500 text-xs font-semibold px-3 py-1 rounded-full mb-4 uppercase tracking-wider">
            Acceso restringido
          </span>
          <h1 class="text-2xl font-bold text-slate-900">Accede a tu empresa</h1>
          <p class="text-slate-500 text-sm mt-1">Solo miembros del equipo ArticDev</p>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-5">
          <!-- Error banner -->
          <div v-if="errorMsg" class="rounded-xl bg-red-50 border border-red-200 px-4 py-3 text-sm text-red-700">
            {{ errorMsg }}
          </div>

          <div>
            <label for="username" class="block text-sm font-semibold text-slate-700 mb-1.5">Usuario</label>
            <input
              id="username"
              v-model="form.username"
              type="text"
              autocomplete="username"
              :disabled="isLoading"
              :class="[
                'block w-full rounded-xl border px-4 py-3 text-slate-900 placeholder-slate-400 text-sm focus:outline-none focus:ring-2 disabled:opacity-50 transition-all',
                fieldErr.username
                  ? 'border-red-400 focus:border-red-400 focus:ring-red-300/30'
                  : 'border-slate-200 focus:border-nordic-cyan focus:ring-nordic-cyan/30'
              ]"
              placeholder="usuario"
              @input="fieldErr.username = ''"
            />
            <p v-if="fieldErr.username" class="mt-1.5 text-xs text-red-600">{{ fieldErr.username }}</p>
          </div>

          <div>
            <label for="password" class="block text-sm font-semibold text-slate-700 mb-1.5">Contraseña</label>
            <input
              id="password"
              v-model="form.password"
              type="password"
              autocomplete="current-password"
              :disabled="isLoading"
              :class="[
                'block w-full rounded-xl border px-4 py-3 text-slate-900 placeholder-slate-400 text-sm focus:outline-none focus:ring-2 disabled:opacity-50 transition-all',
                fieldErr.password
                  ? 'border-red-400 focus:border-red-400 focus:ring-red-300/30'
                  : 'border-slate-200 focus:border-nordic-cyan focus:ring-nordic-cyan/30'
              ]"
              placeholder="••••••••"
              @input="fieldErr.password = ''"
            />
            <p v-if="fieldErr.password" class="mt-1.5 text-xs text-red-600">{{ fieldErr.password }}</p>
          </div>

          <button
            type="submit"
            :disabled="isLoading"
            class="w-full flex justify-center items-center gap-2 rounded-xl bg-gradient-to-r from-nordic-cyan to-blue-500 hover:from-blue-500 hover:to-nordic-cyan px-4 py-3 text-sm font-bold text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-nordic-cyan focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-300"
          >
            <span v-if="isLoading">Ingresando…</span>
            <span v-else>Iniciar sesión</span>
          </button>

          <div class="text-center">
            <RouterLink to="/forgot-password" class="text-sm text-slate-500 hover:text-nordic-cyan transition-colors">
              ¿Olvidaste tu contraseña?
            </RouterLink>
          </div>
        </form>

        <!-- Footer link -->
        <div class="mt-8 pt-6 border-t border-slate-100 text-center">
          <p class="text-xs text-slate-400">
            ¿No eres parte del equipo?
            <RouterLink to="/unete" class="text-nordic-cyan hover:underline font-medium">
              Únete al ecosistema
            </RouterLink>
          </p>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { te, teError } from '@/i18n'

const router    = useRouter()
const authStore = useAuthStore()

const form      = reactive({ username: '', password: '' })
const fieldErr  = reactive({ username: '', password: '' })
const errorMsg  = ref('')
const isLoading = ref(false)

// ── Validación de campos ───────────────────────────────────────────────────
function validate() {
  fieldErr.username = form.username.trim() ? '' : te('form.username_required')
  fieldErr.password = form.password        ? '' : te('form.password_required')
  return !fieldErr.username && !fieldErr.password
}

// ── Submit ─────────────────────────────────────────────────────────────────
const handleLogin = async () => {
  errorMsg.value = ''
  if (!validate()) return

  isLoading.value = true
  try {
    await authStore.login({ username: form.username.trim(), password: form.password })
    router.push({ name: 'dashboard' })
  } catch (err) {
    errorMsg.value = teError(err)
  } finally {
    isLoading.value = false
  }
}
</script>