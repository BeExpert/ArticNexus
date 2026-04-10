<template>
  <div class="max-w-6xl mx-auto px-6 py-8">

    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-xl font-bold text-slate-900">Licencias</h1>
        <p class="text-sm text-slate-500 mt-1">Gestiona qué aplicaciones puede usar cada empresa.</p>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="py-24 text-center text-sm text-slate-400">Cargando datos…</div>

    <!-- Error -->
    <div v-else-if="error" class="mb-4 rounded border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">{{ error }}</div>

    <!-- Matrix table -->
    <div v-else class="an-card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-sm border-collapse" style="min-width: 600px">
          <thead>
            <tr class="bg-slate-800">
              <th class="px-4 py-3 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 sticky left-0 bg-slate-800 z-10">Empresa</th>
              <th class="px-4 py-3 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 sticky left-0 bg-slate-800 z-10 w-24">Estado</th>
              <th
                v-for="app in apps" :key="app.id"
                class="px-4 py-3 text-center text-[11px] font-semibold uppercase tracking-wide text-slate-300 min-w-[120px]"
              >
                <div>{{ app.name }}</div>
                <div class="text-[10px] text-slate-500 font-mono font-normal mt-0.5">{{ app.code }}</div>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="company in companies" :key="company.id"
              class="border-b border-slate-100 last:border-0 hover:bg-slate-50 transition-colors"
            >
              <td class="px-4 py-3 font-medium text-slate-900 sticky left-0 bg-white z-10">
                <div class="flex items-center gap-2.5">
                  <div class="shrink-0 w-7 h-7 rounded-lg flex items-center justify-center text-[11px] font-bold uppercase bg-indigo-100 text-indigo-700">
                    {{ company.name.charAt(0) }}
                  </div>
                  {{ company.name }}
                </div>
              </td>
              <td class="px-4 py-3 sticky left-0 bg-white z-10">
                <span :class="company.status === 'active' ? 'an-badge-active' : 'an-badge-inactive'">
                  {{ company.status === 'active' ? 'Activa' : 'Inactiva' }}
                </span>
              </td>
              <td
                v-for="app in apps" :key="app.id"
                class="px-4 py-3 text-center"
              >
                <label class="inline-flex items-center justify-center cursor-pointer">
                  <input
                    type="checkbox"
                    :checked="isLicensed(company.id, app.id)"
                    @change="toggleLicense(company, app, $event)"
                    :disabled="toggling[`${company.id}-${app.id}`]"
                    class="rounded border-slate-300 text-indigo-600 focus:ring-indigo-500 w-4 h-4"
                  />
                </label>
                <div v-if="isLicensed(company.id, app.id)" class="text-[10px] mt-0.5" :class="getLicenseStatus(company.id, app.id) === 'active' ? 'text-emerald-500' : 'text-amber-500'">
                  {{ getLicenseStatus(company.id, app.id) === 'active' ? 'Activa' : 'Inactiva' }}
                </div>
              </td>
            </tr>
            <tr v-if="companies.length === 0">
              <td :colspan="apps.length + 2" class="py-16 text-center text-sm text-slate-400">No hay empresas registradas.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Legend -->
    <div v-if="!loading && !error" class="mt-4 flex items-center gap-4 text-[11px] text-slate-400">
      <span class="flex items-center gap-1.5">
        <input type="checkbox" checked disabled class="rounded border-slate-300 text-indigo-600 w-3.5 h-3.5" /> Licenciada
      </span>
      <span class="flex items-center gap-1.5">
        <input type="checkbox" disabled class="rounded border-slate-300 w-3.5 h-3.5" /> Sin licencia
      </span>
      <span class="text-slate-300">|</span>
      <span>Los cambios se guardan automáticamente al marcar/desmarcar.</span>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/services/api'

const loading    = ref(true)
const error      = ref('')
const companies  = ref([])
const apps       = ref([])
const licenseMap = ref({})   // { "companyId-appId": { status: "active" } }
const toggling   = ref({})   // { "companyId-appId": true } while request in flight

function licenseKey(companyId, appId) {
  return `${companyId}-${appId}`
}

function isLicensed(companyId, appId) {
  return !!licenseMap.value[licenseKey(companyId, appId)]
}

function getLicenseStatus(companyId, appId) {
  return licenseMap.value[licenseKey(companyId, appId)]?.status ?? ''
}

async function loadData() {
  loading.value = true
  error.value   = ''
  try {
    const [companiesRes, appsRes] = await Promise.all([
      api.get('/companies', { params: { pageSize: 500 } }),
      api.get('/applications', { params: { pageSize: 100 } }),
    ])
    companies.value = (companiesRes.data?.data ?? companiesRes.data ?? [])
    apps.value = (appsRes.data?.data ?? appsRes.data ?? []).filter(a => a.code !== 'ARTICNEXUS')

    // Load licenses for all companies in parallel
    const allLicenses = await Promise.all(
      companies.value.map(c =>
        api.get(`/companies/${c.id}/applications`)
          .then(r => ({ companyId: c.id, apps: r.data?.data ?? r.data ?? [] }))
          .catch(() => ({ companyId: c.id, apps: [] }))
      )
    )
    const map = {}
    for (const { companyId, apps: capps } of allLicenses) {
      for (const a of capps) {
        map[licenseKey(companyId, a.appId)] = { status: a.status }
      }
    }
    licenseMap.value = map
  } catch (e) {
    if (e.response?.status === 401) {
      // Session expired — user will be redirected to login
      return
    }
    error.value = e.response?.data?.message || e.message || 'Error al cargar datos'
  } finally {
    loading.value = false
  }
}

async function toggleLicense(company, app, event) {
  const key     = licenseKey(company.id, app.id)
  const checked = event.target.checked
  toggling.value[key] = true

  try {
    if (checked) {
      await api.post(`/companies/${company.id}/applications`, { appId: app.id })
      licenseMap.value[key] = { status: 'active' }
    } else {
      await api.delete(`/companies/${company.id}/applications/${app.id}`)
      delete licenseMap.value[key]
    }
  } catch (e) {
    // Revert checkbox
    event.target.checked = !checked
    error.value = e.response?.data?.message || 'Error al actualizar licencia'
    setTimeout(() => { error.value = '' }, 4000)
  } finally {
    delete toggling.value[key]
  }
}

onMounted(loadData)
</script>
