<template>
  <!-- ═══════════════════════════════════════════════════════════════
       MODO A — usuario con permisos de administración
  ═══════════════════════════════════════════════════════════════ -->
  <div v-if="hasManageAccess">
    <!-- Header strip -->
    <div class="mb-8 flex items-center gap-4">
      <div class="w-12 h-12 rounded-xl bg-indigo-600 flex items-center justify-center text-white text-xl font-bold shrink-0 shadow-lg shadow-indigo-200">
        {{ userInitial }}
      </div>
      <div>
        <h1 class="text-2xl font-bold text-slate-900 leading-tight">
          {{ greeting }}, {{ firstName || authStore.user?.username }}
        </h1>
        <p class="text-sm text-slate-500 mt-0.5">Resumen general de ArticNexus</p>
      </div>
    </div>

    <!-- KPI Cards -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <div
        v-for="card in kpiCards" :key="card.label"
        class="relative overflow-hidden rounded-xl border border-slate-200 bg-white p-5 shadow-sm hover:shadow-md transition-shadow"
      >
        <div class="absolute inset-0 opacity-[0.03]" :style="{ background: card.bg }"></div>
        <p class="text-[11px] uppercase tracking-widest font-semibold text-slate-400 mb-2">{{ card.label }}</p>
        <p class="text-3xl font-extrabold tabular-nums" :class="card.color">
          <span v-if="statsLoading" class="inline-block w-10 h-7 bg-slate-100 rounded animate-pulse"></span>
          <span v-else>{{ card.value }}</span>
        </p>
        <div class="absolute bottom-3 right-4 text-4xl font-black opacity-[0.06]" :class="card.color">{{ card.icon }}</div>
      </div>
    </div>

    <!-- Quick Access — permission filtered (security: only show what user can access) -->
    <h2 class="text-xs font-bold text-slate-400 uppercase tracking-widest mb-3">Acceso rápido</h2>
    <div v-if="quickLinks.length" class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <RouterLink
        v-for="link in quickLinks" :key="link.to"
        :to="link.to"
        class="group flex items-center gap-4 rounded-xl border border-slate-200 bg-white p-4 shadow-sm hover:shadow-md hover:border-indigo-200 hover:-translate-y-0.5 transition-all"
      >
        <div class="w-10 h-10 rounded-lg flex items-center justify-center text-lg shrink-0 transition-colors" :class="link.iconBg">
          {{ link.icon }}
        </div>
        <div>
          <p class="text-sm font-semibold text-slate-900 group-hover:text-indigo-600 transition-colors">{{ link.label }}</p>
          <p class="text-xs text-slate-400 mt-0.5 leading-snug">{{ link.desc }}</p>
        </div>
      </RouterLink>
    </div>
    <p v-else class="text-sm text-slate-400 italic">Sin accesos rápidos disponibles.</p>
  </div>

  <!-- ═══════════════════════════════════════════════════════════════
       MODO B — usuario sin permisos de administración
  ═══════════════════════════════════════════════════════════════ -->
  <div
    v-else
    class="relative overflow-hidden rounded-2xl min-h-[calc(100vh-8rem)]"
    style="background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #0f172a 100%)"
  >
    <!-- Background: subtle dot grid -->
    <div class="absolute inset-0 opacity-[0.15]" style="background-image: radial-gradient(circle, #94a3b8 1px, transparent 1px); background-size: 28px 28px;"></div>

    <!-- Background: time-of-day accent glow -->
    <div
      class="absolute inset-0 opacity-20 pointer-events-none"
      :style="{ background: `radial-gradient(ellipse 60% 40% at 50% 0%, ${timeAccentColor} 0%, transparent 70%)` }"
    ></div>

    <div class="relative z-10 px-6 py-10 sm:px-10">
      <!-- ── Greeting banner ── -->
      <div class="flex items-start gap-5 mb-12">
        <!-- Avatar ring -->
        <div class="relative shrink-0">
          <div class="w-16 h-16 rounded-2xl flex items-center justify-center text-2xl font-extrabold text-white shadow-xl" :style="{ background: `linear-gradient(135deg, ${timeAccentColor}, hsl(${timeHue + 60}, 60%, 20%))` }">
            {{ userInitial }}
          </div>
          <!-- Pulsing ring -->
          <span class="absolute -inset-1 rounded-2xl opacity-30 animate-ping" :style="{ border: `2px solid ${timeAccentColor}` }"></span>
        </div>

        <div>
          <div class="inline-flex items-center gap-2 mb-2">
            <span class="w-2 h-2 rounded-full animate-pulse" :style="{ background: timeAccentColor }"></span>
            <span class="text-xs font-semibold uppercase tracking-widest" :style="{ color: timeAccentColor }">Sesión activa</span>
          </div>
          <h1 class="text-3xl sm:text-4xl font-extrabold text-white leading-tight">
            {{ greeting }},<br>
            <span class="text-slate-300 font-semibold">{{ firstName || authStore.user?.username }}</span>
          </h1>
          <p class="text-slate-400 text-sm mt-2">
            {{ companiesLoading ? 'Cargando membresías…' : myCompanies.length === 0 ? 'Aún no tienes empresas asignadas.' : `Eres miembro de ${myCompanies.length} empresa${myCompanies.length !== 1 ? 's' : ''}.` }}
          </p>
        </div>
      </div>

      <!-- ── Companies section ── -->
      <template v-if="!companiesLoading">
        <!-- Skeleton while loading technically already gated above, but keep clean -->
        <h2 class="text-xs font-bold uppercase tracking-widest text-slate-500 mb-4">Tus empresas</h2>

        <!-- Company cards grid -->
        <div v-if="myCompanies.length" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
          <div
            v-for="company in myCompanies"
            :key="company.id ?? company.com_id"
            class="relative overflow-hidden rounded-2xl p-6 cursor-default select-none group transition-all duration-300 hover:-translate-y-1.5 hover:shadow-2xl"
            :style="companyCardStyle(company.name ?? company.com_name)"
          >
            <!-- Watermark initial -->
            <span
              class="absolute -right-4 -bottom-6 text-[8rem] leading-none font-black pointer-events-none select-none opacity-[0.08] text-white"
            >{{ (company.name ?? company.com_name ?? '?')[0].toUpperCase() }}</span>

            <!-- Glow on hover -->
            <div
              class="absolute inset-0 opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none"
              :style="{ background: `radial-gradient(ellipse 80% 50% at 50% 0%, ${companyAccentColor(company.name ?? company.com_name)}, transparent)` }"
            ></div>

            <!-- Card content -->
            <div class="relative z-10">
              <!-- Company identity header -->
              <div class="flex items-center gap-3 mb-5">
                <div
                  class="w-11 h-11 rounded-xl flex items-center justify-center text-lg font-extrabold text-white shrink-0 shadow-lg"
                  :style="{ background: companyAccentColor(company.name ?? company.com_name) }"
                >
                  {{ (company.name ?? company.com_name ?? '?')[0].toUpperCase() }}
                </div>
                <div class="min-w-0">
                  <p class="font-bold text-white text-sm leading-tight truncate">{{ company.name ?? company.com_name }}</p>
                  <p class="text-xs text-slate-400 mt-0.5">{{ company.code ?? company.com_code ?? '—' }}</p>
                </div>
              </div>

              <!-- Divider -->
              <div class="h-px mb-4 opacity-20" :style="{ background: companyAccentColor(company.name ?? company.com_name) }"></div>

              <!-- Meta row -->
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-1.5">
                  <span class="w-2 h-2 rounded-full" :style="{ background: companyAccentColor(company.name ?? company.com_name) }"></span>
                  <span class="text-xs text-slate-400 font-medium">Miembro activo</span>
                </div>
                <span v-if="company.branches || company.branch_count" class="text-xs text-slate-500">
                  {{ company.branches ?? company.branch_count }} sucursal{{ (company.branches ?? company.branch_count) !== 1 ? 'es' : '' }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Empty state -->
        <div v-else class="flex flex-col items-center justify-center py-20 text-center">
          <div class="w-20 h-20 rounded-full bg-slate-800 flex items-center justify-center mb-5 border border-slate-700">
            <span class="text-3xl">🏢</span>
          </div>
          <p class="text-slate-300 font-semibold text-lg mb-1">Sin empresas asignadas</p>
          <p class="text-slate-500 text-sm max-w-xs">
            Contacta al administrador para que te asigne a una empresa y puedas acceder a los módulos disponibles.
          </p>
        </div>
      </template>

      <!-- Loading skeleton -->
      <template v-else>
        <h2 class="text-xs font-bold uppercase tracking-widest text-slate-500 mb-4">Tus empresas</h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
          <div v-for="i in 3" :key="i" class="h-40 rounded-2xl bg-slate-800/60 animate-pulse"></div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/store/auth'
import { companyService } from '@/services/companyService'
import api from '@/services/api'

const authStore = useAuthStore()

// ── User info ──────────────────────────────────────────────────
const firstName  = computed(() => authStore.user?.person?.firstName ?? authStore.user?.name ?? '')
const userInitial = computed(() => {
  const name = firstName.value || authStore.user?.username || '?'
  return name[0].toUpperCase()
})

// ── Permission gate: does the user have any management permission? ──
const hasManageAccess = computed(() =>
  authStore.canPrefix('empresas') ||
  authStore.canPrefix('usuarios') ||
  authStore.canPrefix('aplicaciones') ||
  authStore.canPrefix('roles') ||
  authStore.canPrefix('demo') ||
  authStore.user?.isSuperAdmin === true
)

// ── Time-based greeting ────────────────────────────────────────
const hour = new Date().getHours()
const greeting = hour < 12 ? 'Buenos días' : hour < 19 ? 'Buenas tardes' : 'Buenas noches'

// Accent color per time of day: morning=amber, afternoon=indigo, evening=violet
const timeHue   = hour < 12 ? 38 : hour < 19 ? 234 : 263
const timeAccentColor = computed(() =>
  hour < 12 ? '#f59e0b' : hour < 19 ? '#6366f1' : '#8b5cf6'
)

// ── MODO A: admin stats + permission-filtered quick links ──────
const stats        = ref({ users: 0, companies: 0, applications: 0, roles: 0 })
const statsLoading = ref(true)

const kpiCards = computed(() => [
  { label: 'Usuarios',     value: stats.value.users,        color: 'text-indigo-600',  bg: 'linear-gradient(135deg,#6366f1,#4338ca)', icon: '👤' },
  { label: 'Empresas',     value: stats.value.companies,    color: 'text-violet-600',  bg: 'linear-gradient(135deg,#7c3aed,#6d28d9)', icon: '🏢' },
  { label: 'Aplicaciones', value: stats.value.applications, color: 'text-sky-600',     bg: 'linear-gradient(135deg,#0ea5e9,#0284c7)', icon: '⚡' },
  { label: 'Roles',        value: stats.value.roles,        color: 'text-emerald-600', bg: 'linear-gradient(135deg,#10b981,#059669)', icon: '🔑' },
])

const allQuickLinks = [
  { to: '/companies',    label: 'Empresas',     desc: 'Gestionar empresas y sucursales', prefix: 'empresas',    icon: '🏢', iconBg: 'bg-violet-100 text-violet-600 group-hover:bg-violet-200' },
  { to: '/users',        label: 'Usuarios',     desc: 'Administrar cuentas de usuario',  prefix: 'usuarios',    icon: '👤', iconBg: 'bg-indigo-100 text-indigo-600 group-hover:bg-indigo-200' },
  { to: '/applications', label: 'Aplicaciones', desc: 'Configurar apps y módulos',       prefix: 'aplicaciones',icon: '⚡', iconBg: 'bg-sky-100 text-sky-600 group-hover:bg-sky-200'          },
  { to: '/roles',        label: 'Roles',        desc: 'Gestionar roles y permisos',      prefix: 'roles',       icon: '🔑', iconBg: 'bg-emerald-100 text-emerald-600 group-hover:bg-emerald-200' },
]

// Security fix: only expose quick links for sections the user can actually access
const quickLinks = computed(() =>
  allQuickLinks.filter(link => authStore.canPrefix(link.prefix) || authStore.user?.isSuperAdmin)
)

// ── MODO B: company membership cards ──────────────────────────
const myCompanies    = ref([])
const companiesLoading = ref(false)

// Deterministic hue from company name
function companyHue(name) {
  let h = 0
  for (const c of (name || '')) h = (h * 31 + c.charCodeAt(0)) & 0xFFFF
  return (h * 137) % 360
}

function companyAccentColor(name) {
  return `hsl(${companyHue(name)}, 60%, 55%)`
}

function companyCardStyle(name) {
  const hue = companyHue(name)
  return {
    background: `linear-gradient(140deg, hsl(${hue}, 35%, 13%) 0%, hsl(${hue + 25}, 25%, 8%) 100%)`,
    border: `1px solid hsl(${hue}, 35%, 22%)`,
    boxShadow: `0 4px 20px -4px hsl(${hue}, 50%, 10%)`,
  }
}

// ── Lifecycle ──────────────────────────────────────────────────
onMounted(async () => {
  if (hasManageAccess.value) {
    try {
      const res = await api.get('/stats')
      stats.value = res.data?.data ?? res.data
    } catch { /* keep defaults */ } finally {
      statsLoading.value = false
    }
  } else {
    companiesLoading.value = true
    try {
      const res = await companyService.getMyCompanies()
      myCompanies.value = res?.data ?? res ?? []
    } catch { myCompanies.value = [] } finally {
      companiesLoading.value = false
    }
  }
})
</script>