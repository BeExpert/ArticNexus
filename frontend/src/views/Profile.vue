<template>
  <div class="max-w-4xl mx-auto">
    <h1 class="text-xl font-semibold text-slate-900 mb-6">Mi Perfil</h1>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Left: User card -->
      <div class="overflow-hidden rounded-xl border border-slate-200 shadow-sm lg:col-span-1">
        <!-- Dark avatar header -->
        <div class="relative px-6 py-8 text-center" style="background: linear-gradient(135deg, #0f172a, #1e293b)">
          <!-- Subtle dot pattern -->
          <div class="absolute inset-0 opacity-[0.15]" style="background-image: radial-gradient(circle, #94a3b8 1px, transparent 1px); background-size: 20px 20px;"></div>
          <div class="relative z-10">
            <div class="w-20 h-20 rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-xl text-3xl font-extrabold text-white" style="background: linear-gradient(135deg, #6366f1, #7c3aed)">
              {{ initial }}
            </div>
            <p class="text-lg font-bold text-white">{{ displayName }}</p>
            <p class="text-sm text-slate-400 mt-0.5 font-mono">{{ authStore.user?.username }}</p>
          </div>
        </div>
        <!-- Meta info -->
        <div class="px-6 py-4 bg-white space-y-2">
          <div class="flex items-center gap-2">
            <svg class="w-4 h-4 text-slate-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21.75 6.75v10.5a2.25 2.25 0 0 1-2.25 2.25h-15a2.25 2.25 0 0 1-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0 0 19.5 4.5h-15a2.25 2.25 0 0 0-2.25 2.25m19.5 0v.243a2.25 2.25 0 0 1-1.07 1.916l-7.5 4.615a2.25 2.25 0 0 1-2.36 0L3.32 8.91a2.25 2.25 0 0 1-1.07-1.916V6.75" />
            </svg>
            <span class="text-sm text-slate-600 truncate">{{ authStore.user?.email ?? '—' }}</span>
          </div>
          <div class="flex items-center gap-2">
            <svg class="w-4 h-4 text-slate-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 0 1 2.25-2.25h13.5A2.25 2.25 0 0 1 21 7.5v11.25m-18 0A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75m-18 0v-7.5A2.25 2.25 0 0 1 5.25 9h13.5A2.25 2.25 0 0 1 21 11.25v7.5" />
            </svg>
            <span class="text-sm text-slate-500">Miembro desde {{ memberSince }}</span>
          </div>
        </div>
      </div>

      <!-- Right: Editable form -->
      <div class="an-card p-6 lg:col-span-2">
        <h2 class="text-sm font-semibold text-slate-900 mb-4">Información personal</h2>

        <div v-if="successMsg" class="mb-4 rounded border border-green-200 bg-green-50 px-3 py-2 text-sm text-green-700">{{ successMsg }}</div>
        <div v-if="errorMsg" class="mb-4 rounded border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">{{ errorMsg }}</div>

        <form @submit.prevent="saveProfile">
          <div class="grid grid-cols-2 gap-4 mb-4">
            <div>
              <label class="an-label">Nombre</label>
              <input v-model="form.firstName" type="text" class="an-input" />
            </div>
            <div>
              <label class="an-label">Primer Apellido</label>
              <input v-model="form.firstSurname" type="text" class="an-input" />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-4 mb-4">
            <div>
              <label class="an-label">Segundo Apellido</label>
              <input v-model="form.secondSurname" type="text" class="an-input" />
            </div>
            <div>
              <label class="an-label">Identificación</label>
              <input v-model="form.nationalId" type="text" class="an-input" />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-4 mb-4">
            <div>
              <label class="an-label">Fecha de Nacimiento</label>
              <input v-model="form.birthDate" type="date" class="an-input" />
            </div>
            <div>
              <label class="an-label">Código de Área</label>
              <input v-model="form.phoneAreaCode" type="text" class="an-input" />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-4 mb-4">
            <div>
              <label class="an-label">Teléfono Principal</label>
              <input v-model="form.primaryPhone" type="text" class="an-input" />
            </div>
            <div>
              <label class="an-label">Teléfono Secundario</label>
              <input v-model="form.secondaryPhone" type="text" class="an-input" />
            </div>
          </div>
          <div class="mb-4">
            <label class="an-label">Dirección</label>
            <input v-model="form.address" type="text" class="an-input" />
          </div>

          <div class="border-t border-slate-100 my-4"></div>

          <div class="mb-4">
            <label class="an-label">Correo electrónico</label>
            <input v-model="form.email" type="email" class="an-input" />
          </div>

          <div class="flex items-center justify-between">
            <button type="submit" :disabled="saving" class="an-btn-primary">
              {{ saving ? 'Guardando…' : 'Guardar cambios' }}
            </button>
          </div>
        </form>

        <p class="text-xs text-slate-400 mt-6">
          Para cambiar tu contraseña, usa la opción "¿Olvidaste tu contraseña?" en la pantalla de inicio de sesión.
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/store/auth'
import api from '@/services/api'
import { teError } from '@/i18n'

const authStore = useAuthStore()

const initial = computed(() => (authStore.user?.username ?? '?')[0])
const displayName = computed(() => {
  const p = authStore.user?.person
  if (!p) return authStore.user?.username ?? ''
  return [p.firstName, p.firstSurname, p.secondSurname].filter(Boolean).join(' ')
})
const memberSince = computed(() => {
  const d = authStore.user?.createdAt
  if (!d) return '—'
  return new Date(d).toLocaleDateString('es-CR', { year: 'numeric', month: 'long' })
})

const form = reactive({
  email: '', firstName: '', firstSurname: '', secondSurname: '',
  nationalId: '', birthDate: '', phoneAreaCode: '',
  primaryPhone: '', secondaryPhone: '', address: '',
})

const saving = ref(false)
const successMsg = ref('')
const errorMsg = ref('')

function populateForm() {
  const u = authStore.user
  const p = u?.person ?? {}
  form.email          = u?.email ?? ''
  form.firstName      = p.firstName ?? ''
  form.firstSurname   = p.firstSurname ?? ''
  form.secondSurname  = p.secondSurname ?? ''
  form.nationalId     = p.nationalId ?? ''
  form.birthDate      = p.birthDate ? p.birthDate.substring(0, 10) : ''
  form.phoneAreaCode  = p.phoneAreaCode ?? ''
  form.primaryPhone   = p.primaryPhone ?? ''
  form.secondaryPhone = p.secondaryPhone ?? ''
  form.address        = p.address ?? ''
}

onMounted(populateForm)

async function saveProfile() {
  saving.value = true
  successMsg.value = ''
  errorMsg.value = ''
  try {
    await api.put('/auth/me', form)
    await authStore.getCurrentUser()
    populateForm()
    successMsg.value = 'Perfil actualizado correctamente.'
  } catch (e) {
    errorMsg.value = teError(e)
  } finally {
    saving.value = false
  }
}
</script>
