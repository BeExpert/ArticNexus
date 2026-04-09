<template>
  <div id="app" class="min-h-screen bg-gray-50">

    <!-- ── Top bar (authenticated only) ─────────────────────────── -->
    <header v-if="authStore.isAuthenticated" class="fixed top-0 inset-x-0 z-30 h-14 bg-slate-900 flex items-center px-5" style="box-shadow: 0 1px 0 rgba(255,255,255,0.06), 0 4px 16px rgba(0,0,0,0.3)">
      <!-- Brand -->
      <RouterLink to="/dashboard" class="flex items-center gap-2.5 mr-8 shrink-0 group">
        <div class="w-7 h-7 rounded-lg flex items-center justify-center text-white text-xs font-black shrink-0 transition-transform group-hover:scale-105" style="background: linear-gradient(135deg, #6366f1, #4338ca)">
          N
        </div>
        <span class="text-white font-semibold text-[15px] tracking-tight">ArticNexus</span>
      </RouterLink>

      <!-- Nav links inline -->
      <nav class="hidden sm:flex items-center">
        <RouterLink
          v-for="item in navItems" :key="item.to"
          :to="item.to"
          class="relative px-3 h-14 flex items-center text-sm font-medium text-slate-400 hover:text-white transition-colors"
          active-class="!text-white nav-active"
        >
          {{ item.label }}
        </RouterLink>
      </nav>

      <!-- Spacer -->
      <div class="flex-1"></div>

      <!-- User badge + logout -->
      <div class="flex items-center gap-3">
        <RouterLink to="/profile" class="flex items-center gap-2 group">
          <span class="text-slate-400 text-sm hidden sm:inline group-hover:text-white transition-colors">{{ authStore.user?.username }}</span>
          <div class="w-8 h-8 rounded-full flex items-center justify-center text-white text-sm font-bold uppercase transition-opacity group-hover:opacity-80" style="background: linear-gradient(135deg, #6366f1, #7c3aed)">
            {{ (authStore.user?.username ?? '?')[0] }}
          </div>
        </RouterLink>
        <div class="w-px h-4 bg-slate-700 hidden sm:block"></div>
        <button @click="handleLogout" class="text-slate-400 hover:text-white transition-colors text-sm font-medium hidden sm:inline">
          Salir
        </button>
      </div>

      <!-- Mobile hamburger -->
      <button
        @click="mobileMenuOpen = !mobileMenuOpen"
        class="sm:hidden text-white p-1 rounded hover:bg-white/10 ml-2"
        aria-label="Menú"
      >
        <svg v-if="mobileMenuOpen" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
        <svg v-else class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </button>
    </header>

    <!-- Mobile dropdown menu -->
    <transition name="fade">
      <div v-if="authStore.isAuthenticated && mobileMenuOpen" class="fixed inset-0 z-40 bg-black/40 sm:hidden" @click="mobileMenuOpen = false" />
    </transition>
    <transition name="slide-down">
      <div v-if="authStore.isAuthenticated && mobileMenuOpen" class="fixed top-14 left-0 right-0 z-50 bg-white shadow-lg border-b border-slate-200 sm:hidden">
        <nav class="py-2 px-3 space-y-1">
          <RouterLink
            v-for="item in navItems" :key="item.to"
            :to="item.to"
            @click="mobileMenuOpen = false"
            class="block px-3 py-2 rounded text-sm font-medium text-slate-600 hover:bg-slate-100 hover:text-slate-900 transition-colors"
            active-class="bg-slate-100 text-slate-900 font-semibold"
          >
            {{ item.label }}
          </RouterLink>
        </nav>
        <div class="border-t px-4 py-3 flex items-center justify-between">
          <span class="text-sm text-slate-600">{{ authStore.user?.username }}</span>
          <button @click="handleLogout" class="text-sm text-red-600 hover:text-red-700 font-medium">Cerrar sesión</button>
        </div>
      </div>
    </transition>

    <!-- ── Page content ──────────────────────────────────────────── -->
    <main :class="authStore.isAuthenticated ? 'pt-14' : ''">
      <div :class="authStore.isAuthenticated ? (route.meta.fullWidth ? 'h-[calc(100vh-3.5rem)] overflow-hidden' : 'max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8') : ''">  
        <RouterView />
      </div>
    </main>

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/store/auth'

const router    = useRouter()
const route     = useRoute()
const authStore = useAuthStore()
const mobileMenuOpen = ref(false)

const allNavItems = [
  { to: '/dashboard',     label: 'Inicio',        prefix: null             },
  { to: '/companies',     label: 'Empresas',      prefix: 'empresas'       },
  { to: '/users',         label: 'Usuarios',      prefix: 'usuarios'       },
  { to: '/applications',  label: 'Aplicaciones',  prefix: 'aplicaciones'   },
  { to: '/demo-links',    label: 'Demo Links',     prefix: 'demo_links'     },
]

const navItems = computed(() =>
  allNavItems.filter(item => !item.prefix || authStore.canPrefix(item.prefix))
)

onMounted(() => authStore.initialize())

const handleLogout = async () => {
  mobileMenuOpen.value = false
  await authStore.logout()
  router.push({ name: 'login' })
}
</script>

<style scoped>
/* Active nav link: bottom indigo underline */
.nav-active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0.5rem;
  right: 0.5rem;
  height: 2px;
  background: #6366f1;
  border-radius: 1px;
}

/* Mobile dropdown slide */
.slide-down-enter-active,
.slide-down-leave-active {
  transition: transform 0.2s ease, opacity 0.2s ease;
}
.slide-down-enter-from,
.slide-down-leave-to {
  transform: translateY(-0.5rem);
  opacity: 0;
}

/* Overlay fade */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>