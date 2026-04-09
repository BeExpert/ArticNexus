<template>
  <div>
    <!-- Header -->
    <div class="an-page-header">
      <div>
        <h1 class="an-page-header-title">Demo Links</h1>
        <p class="an-page-header-subtitle">Genera links de acceso temporal para demostrar las apps a clientes.</p>
      </div>
      <button
        v-if="authStore.can('demo_links.crear')"
        @click="openCreate"
        class="an-btn-primary text-sm py-2 px-4 shrink-0"
      >
        Nuevo Demo Link
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="py-20 text-center text-sm text-slate-400">Cargando links…</div>

    <!-- Error -->
    <div v-else-if="error" class="rounded border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">{{ error }}</div>

    <!-- Empty -->
    <div v-else-if="links.length === 0" class="an-card py-16 text-center">
      <p class="text-sm font-medium text-slate-600">Sin demo links</p>
      <p class="text-xs text-slate-400 mt-1">Genera el primer link para que un cliente pruebe una app.</p>
    </div>

    <!-- Table -->
    <div v-else class="an-card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-sm border-collapse" style="min-width:620px">
          <thead>
            <tr class="bg-slate-800">
              <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300">App</th>
              <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 hidden sm:table-cell">Destinatario</th>
              <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300">Expira</th>
              <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300">Estado</th>
              <th class="px-4 py-2.5 text-right text-[11px] font-semibold uppercase tracking-wide text-slate-300">Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="l in links"
              :key="l.id"
              class="border-b border-slate-100 last:border-0 hover:bg-slate-50 transition-colors"
            >
              <td class="px-4 py-3 font-medium text-slate-900">{{ l.appCode }}</td>
              <td class="px-4 py-3 text-slate-500 hidden sm:table-cell">{{ l.recipientEmail ?? '—' }}</td>
              <td class="px-4 py-3 text-slate-500">{{ fmtDate(l.expiresAt) }}</td>
              <td class="px-4 py-3">
                <span
                  class="inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-semibold"
                  :class="l.isActive && !isExpired(l.expiresAt)
                    ? 'bg-emerald-50 text-emerald-700'
                    : 'bg-slate-100 text-slate-400'"
                >
                  {{ l.isActive && !isExpired(l.expiresAt) ? 'Activo' : 'Inactivo' }}
                </span>
              </td>
              <td class="px-4 py-3 text-right">
                <div class="flex items-center justify-end gap-2">
                  <span
                    v-if="copiedId === l.id"
                    class="text-[11px] text-emerald-600 flex items-center gap-1"
                  >
                    ✓ Copiado
                  </span>
                  <span
                    v-else-if="noTokenId === l.id"
                    class="text-[11px] text-amber-600"
                    title="Recarga la página para regenerar un link nuevo"
                  >
                    Recarga para copiar
                  </span>
                  <button
                    v-else-if="l.isActive && !isExpired(l.expiresAt)"
                    @click="copyLink(l)"
                    title="Copiar URL"
                    class="p-1.5 rounded text-slate-500 hover:text-slate-800 hover:bg-slate-100 transition-colors"
                  ><Copy :size="15" /></button>
                  <button
                    v-if="authStore.can('demo_links.eliminar') && l.isActive"
                    @click="confirmRevoke(l)"
                    title="Revocar"
                    class="p-1.5 rounded text-red-400 hover:text-red-600 hover:bg-red-50 transition-colors"
                  ><Ban :size="15" /></button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create Modal -->
    <div v-if="showCreate" class="an-overlay" @click.self="closeCreate">
      <div class="an-modal w-full max-w-md">
        <div class="an-modal-header">
          <h2 class="text-base font-semibold text-slate-900">Nuevo Demo Link</h2>
          <button @click="closeCreate" class="text-slate-400 hover:text-slate-600 transition-colors">✕</button>
        </div>
        <div class="an-modal-body px-6 py-4 space-y-4">

          <!-- App -->
          <div>
            <label class="block text-xs font-semibold text-slate-500 mb-1">Aplicación</label>
            <select v-model="form.appCode" class="an-input text-sm w-full">
              <option value="" disabled>Selecciona una app…</option>
              <option value="OFTADATA">OftaData</option>
              <option value="VETDATA">VetData</option>
            </select>
          </div>

          <!-- Mode toggle -->
          <div>
            <label class="block text-xs font-semibold text-slate-500 mb-2">Tipo de destinatario</label>
            <div class="grid grid-cols-2 gap-1 bg-slate-100 rounded-lg p-1">
              <button
                type="button"
                @click="demoMode = 'sistema'"
                :class="demoMode === 'sistema'
                  ? 'bg-white text-slate-900 shadow-sm'
                  : 'text-slate-500 hover:text-slate-700'"
                class="rounded-md py-1.5 text-xs font-medium transition-all"
              >
                Usuario del sistema
              </button>
              <button
                type="button"
                @click="demoMode = 'externo'"
                :class="demoMode === 'externo'
                  ? 'bg-white text-slate-900 shadow-sm'
                  : 'text-slate-500 hover:text-slate-700'"
                class="rounded-md py-1.5 text-xs font-medium transition-all"
              >
                Invitado externo
              </button>
            </div>
            <p class="text-[11px] text-slate-400 mt-1.5">
              <template v-if="demoMode === 'sistema'">
                Busca un cliente ya registrado — el demo link se enviará a su email.
              </template>
              <template v-else>
                Escribe el nombre y email del prospecto que quieres invitar.
              </template>
            </p>
          </div>

          <!-- Mode A: system user search -->
          <div v-if="demoMode === 'sistema'" class="relative">
            <label class="block text-xs font-semibold text-slate-500 mb-1">Buscar usuario</label>
            <!-- Selected -->
            <div v-if="selectedUser" class="an-input text-sm">
              <div class="flex items-start justify-between gap-2">
                <div class="min-w-0">
                  <p class="font-medium text-slate-900 truncate">{{ selectedUser.displayName }}</p>
                  <p class="text-[11px] text-slate-400">
                    {{ selectedUser.username }}
                    <span v-if="selectedUser.email"> · {{ selectedUser.email }}</span>
                  </p>
                </div>
                <button type="button" @click="clearUser" class="text-slate-400 hover:text-slate-600 shrink-0 text-base leading-none">&times;</button>
              </div>
            </div>
            <!-- Search -->
            <template v-else>
              <input
                v-model="userQuery"
                type="text"
                placeholder="Escribe nombre o usuario…"
                class="an-input text-sm w-full"
                autocomplete="off"
                @blur="hideDropdown"
              />
              <p class="text-[11px] text-slate-400 mt-0.5">Mínimo 2 caracteres para buscar.</p>
              <ul
                v-if="userDropdown.length"
                class="absolute z-20 mt-1 w-full bg-white border border-slate-200 rounded-lg shadow-lg overflow-hidden"
              >
                <li
                  v-for="u in userDropdown"
                  :key="u.id"
                  @mousedown.prevent="pickUser(u)"
                  class="px-3 py-2.5 hover:bg-slate-50 cursor-pointer border-b border-slate-100 last:border-0"
                >
                  <p class="text-sm font-medium text-slate-900">{{ userFullName(u) }}</p>
                  <p class="text-[11px] text-slate-400">
                    {{ u.username }}
                    <span v-if="u.email"> · {{ u.email }}</span>
                  </p>
                </li>
              </ul>
              <p v-if="userSearchLoading" class="text-[11px] text-slate-400 mt-1">Buscando…</p>
            </template>
          </div>

          <!-- Mode B: guest name + email -->
          <template v-else>
            <div>
              <label class="block text-xs font-semibold text-slate-500 mb-1">Nombre del invitado</label>
              <input
                v-model="form.guestName"
                type="text"
                placeholder="Ej. María López"
                class="an-input text-sm w-full"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-500 mb-1">Email del invitado</label>
              <input
                v-model="form.guestEmail"
                type="email"
                placeholder="cliente@ejemplo.com"
                class="an-input text-sm w-full"
              />
            </div>
          </template>

          <!-- Duration -->
          <div>
            <label class="block text-xs font-semibold text-slate-500 mb-1">Duración (horas)</label>
            <input
              v-model.number="form.expiresInHours"
              type="number"
              min="1"
              max="720"
              placeholder="24"
              class="an-input text-sm w-full"
            />
          </div>

          <div v-if="createError" class="rounded border border-red-200 bg-red-50 px-3 py-2 text-xs text-red-700">{{ createError }}</div>
        </div>
        <div class="an-modal-footer flex justify-end gap-2 px-6 pb-5">
          <button @click="closeCreate" class="an-btn-secondary text-sm">Cancelar</button>
          <button @click="submitCreate" :disabled="creating" class="an-btn-primary text-sm">
            {{ creating ? 'Generando…' : 'Generar link' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Token Display Modal (shown after creation) -->
    <div v-if="newToken" class="an-overlay" @click.self="newToken = null">
      <div class="an-modal w-full max-w-lg">
        <div class="an-modal-header">
          <h2 class="text-base font-semibold text-slate-900">¡Link generado!</h2>
          <button @click="newToken = null" class="text-slate-400 hover:text-slate-600">✕</button>
        </div>
        <div class="an-modal-body px-6 py-4 space-y-3">
          <p class="text-sm text-slate-600">Copia y comparte esta URL. Solo se muestra <strong>una vez</strong>.</p>
          <div class="bg-slate-50 rounded border border-slate-200 px-3 py-2 break-all text-xs font-mono text-slate-700 select-all">
            {{ buildDemoUrl(newToken.appCode, newToken.token) }}
          </div>
          <button @click="copyTokenUrl(newToken)" class="an-btn-secondary text-sm w-full">
            {{ tokenCopied ? '✓ Copiado' : 'Copiar URL' }}
          </button>
        </div>
        <div class="an-modal-footer flex justify-end px-6 pb-5">
          <button @click="newToken = null" class="an-btn-primary text-sm">Listo</button>
        </div>
      </div>
    </div>

    <!-- Revoke confirm modal -->
    <div v-if="revokeTarget" class="an-overlay" @click.self="revokeTarget = null">
      <div class="an-modal w-full max-w-sm">
        <div class="an-modal-header">
          <h2 class="text-base font-semibold text-slate-900">Revocar link</h2>
        </div>
        <div class="an-modal-body px-6 py-4">
          <p class="text-sm text-slate-600">¿Seguro que quieres revocar este link? El acceso demo quedará desactivado inmediatamente.</p>
        </div>
        <div class="an-modal-footer flex justify-end gap-2 px-6 pb-5">
          <button @click="revokeTarget = null" class="an-btn-secondary text-sm">Cancelar</button>
          <button @click="submitRevoke" :disabled="revoking" class="an-btn-danger text-sm">
            {{ revoking ? 'Revocando…' : 'Revocar' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useAuthStore } from '@/store/auth'
import api from '@/services/api'
import { te, teError } from '@/i18n'
import { Copy, Ban } from 'lucide-vue-next'

const authStore = useAuthStore()

// ─── State ────────────────────────────────────────────────────────────────────
const links   = ref([])
const loading = ref(false)
const error   = ref(null)

const showCreate  = ref(false)
const creating    = ref(false)
const createError = ref(null)
const demoMode = ref('externo') // 'sistema' | 'externo'
const form = ref({ appCode: '', expiresInHours: 24, guestName: '', guestEmail: '' })

// ─── User Search (inline combobox) ─────────────────────────────────────────────
const selectedUser      = ref(null)  // { id, displayName, username }
const userQuery         = ref('')
const userDropdown      = ref([])
const userSearchLoading = ref(false)
let   userDebounce

function userFullName(u) {
  return [u.person?.firstSurname, u.person?.secondSurname, u.person?.firstName].filter(Boolean).join(' ') || u.username
}

watch(userQuery, (val) => {
  clearTimeout(userDebounce)
  userDropdown.value = []
  if (val.trim().length < 2) { userSearchLoading.value = false; return }
  userSearchLoading.value = true
  userDebounce = setTimeout(async () => {
    try {
      const res = await api.get('/users', { params: { search: val.trim(), page_size: 8 } })
      userDropdown.value = res.data.data?.items ?? res.data.data ?? []
    } catch { userDropdown.value = [] }
    finally { userSearchLoading.value = false }
  }, 300)
})

function hideDropdown() { setTimeout(() => { userDropdown.value = [] }, 150) }
function clearUser()   { selectedUser.value = null; userQuery.value = ''; userDropdown.value = [] }
function pickUser(u)   {
  selectedUser.value = {
    id: u.id,
    displayName: userFullName(u),
    username: u.username,
    email: u.email ?? null,
  }
  userQuery.value    = ''
  userDropdown.value = []
}

const newToken    = ref(null) // { appCode, token }
const tokenCopied = ref(false)

const revokeTarget = ref(null)
const revoking     = ref(false)

const copiedId  = ref(null)
const noTokenId = ref(null) // id of a link whose token is not in session (page was reloaded)

// Session-scoped store of raw tokens: Map<linkId, rawJWT>.
// Populated on creation; cleared if the page is reloaded (by design — JWT is never persisted).
const sessionTokens = new Map()

// ─── App → frontend URL map ───────────────────────────────────────────────────
const APP_DEMO_BASE = {
  OFTADATA: import.meta.env.VITE_OFTADATA_URL ?? 'http://localhost:5174',
  VETDATA:  import.meta.env.VITE_VETDATA_URL  ?? 'http://localhost:5175',
}

function buildDemoUrl(appCode, token) {
  const base = APP_DEMO_BASE[appCode] ?? ''
  return `${base}/demo/enter?token=${token}`
}

// ─── Helpers ──────────────────────────────────────────────────────────────────
function fmtDate(dt) {
  return new Date(dt).toLocaleString('es', { dateStyle: 'short', timeStyle: 'short' })
}
function isExpired(dt) {
  return new Date(dt) < new Date()
}

// ─── Load ──────────────────────────────────────────────────────────────────────
async function load() {
  loading.value = true
  error.value   = null
  try {
    const res = await api.get('/demo-links')
    links.value = res.data.data ?? []
  } catch (e) {
    error.value = teError(e)
  } finally {
    loading.value = false
  }
}

onMounted(load)

// ─── Create ───────────────────────────────────────────────────────────────────
function openCreate() {
  form.value        = { appCode: '', expiresInHours: 24, guestName: '', guestEmail: '' }
  demoMode.value    = 'externo'
  selectedUser.value  = null
  createError.value   = null
  showCreate.value    = true
}
function closeCreate() { showCreate.value = false }

async function submitCreate() {
  createError.value = null
  if (!form.value.appCode) {
    createError.value = te('form.select_app')
    return
  }
  if (demoMode.value === 'sistema' && !selectedUser.value) {
    createError.value = te('form.select_user')
    return
  }
  if (demoMode.value === 'externo' && !form.value.guestEmail) {
    createError.value = te('form.field_required')
    return
  }
  creating.value = true
  try {
    const payload = {
      appCode: form.value.appCode,
      demoUserId: 0, // backend auto-resolves demo_[app] account
      expiresInHours: form.value.expiresInHours || 24,
    }
    if (demoMode.value === 'sistema' && selectedUser.value) {
      payload.guestName = selectedUser.value.displayName
      if (selectedUser.value.email) payload.recipientEmail = selectedUser.value.email
    } else {
      if (form.value.guestName)  payload.guestName = form.value.guestName
      payload.recipientEmail = form.value.guestEmail
    }
    const res = await api.post('/demo-links', payload)
    const created = res.data.data
    // Remember the raw token for this session so the list's "Copiar URL" button works.
    if (created.id && created.token) sessionTokens.set(created.id, created.token)
    closeCreate()
    await load()
    // Show token modal
    tokenCopied.value = false
    newToken.value = { appCode: created.appCode, token: created.token }
  } catch (e) {
    createError.value = teError(e)
  } finally {
    creating.value = false
  }
}

// ─── Copy ──────────────────────────────────────────────────────────────────────
async function copyLink(l) {
  const token = sessionTokens.get(l.id)
  if (!token) {
    // Token not available after a page reload — JWT is never persisted by design.
    noTokenId.value = l.id
    setTimeout(() => { noTokenId.value = null }, 3000)
    return
  }
  const url = buildDemoUrl(l.appCode, token)
  try {
    await navigator.clipboard.writeText(url)
    copiedId.value = l.id
    setTimeout(() => { copiedId.value = null }, 2000)
  } catch {
    // Fallback: create a temporary input and use execCommand
    const input = document.createElement('input')
    input.value = url
    document.body.appendChild(input)
    input.select()
    document.execCommand('copy')
    document.body.removeChild(input)
    copiedId.value = l.id
    setTimeout(() => { copiedId.value = null }, 2000)
  }
}

async function copyTokenUrl(item) {
  const url = buildDemoUrl(item.appCode, item.token)
  await navigator.clipboard.writeText(url)
  tokenCopied.value = true
  setTimeout(() => { tokenCopied.value = false }, 2000)
}

// ─── Revoke ───────────────────────────────────────────────────────────────────
function confirmRevoke(l) { revokeTarget.value = l }

async function submitRevoke() {
  if (!revokeTarget.value) return
  revoking.value = true
  try {
    await api.delete(`/demo-links/${revokeTarget.value.id}`)
    revokeTarget.value = null
    await load()
  } catch (e) {
    // just close and reload
    revokeTarget.value = null
    await load()
  } finally {
    revoking.value = false
  }
}
</script>
