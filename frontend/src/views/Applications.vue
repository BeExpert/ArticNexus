<template>
  <div class="flex h-full overflow-hidden">

    <!-- ══ LEFT PANEL: Application list ════════════════════════════════════ -->
    <!-- Hidden on mobile when an app is selected; visible on sm+ always -->
    <aside :class="[
      'shrink-0 flex flex-col border-r border-slate-200 bg-white',
      'w-full sm:w-72',
      selected ? 'hidden sm:flex' : 'flex'
    ]">

      <!-- Header -->
      <div class="px-4 py-3.5 border-b border-slate-700 bg-slate-800 flex items-center justify-between shrink-0">
        <div class="flex items-center gap-2">
          <span class="text-sm font-semibold text-white">Aplicaciones</span>
          <span class="text-xs text-slate-400 tabular-nums">({{ apps.length }})</span>
        </div>
      </div>

      <!-- Search -->
      <div class="px-3 py-2 border-b border-slate-100 shrink-0">
        <input
          v-model="search"
          type="text"
          placeholder="Buscar aplicación…"
          class="an-input text-xs py-1.5"
        />
      </div>

      <!-- List -->
      <div class="flex-1 overflow-y-auto">
        <div v-if="loading" class="py-12 flex flex-col items-center gap-2">
          <div class="w-5 h-5 rounded-full border-2 border-slate-200 border-t-slate-500 animate-spin"></div>
          <span class="text-xs text-slate-400">Cargando…</span>
        </div>
        <div v-else-if="filteredApps.length === 0" class="py-12 text-center">
          <p class="text-xs text-slate-400">No hay resultados.</p>
        </div>
        <div v-for="app in filteredApps" :key="app.id">
          <!-- App row -->
          <div
            @click="selectApp(app)"
            :class="[
              'w-full text-left px-3 py-2.5 border-b border-slate-100 transition-colors flex items-center gap-3 cursor-pointer select-none',
              selected?.id === app.id ? 'bg-slate-900 border-slate-800' : 'hover:bg-slate-50'
            ]"
          >
            <!-- Avatar -->
            <div :class="[
              'shrink-0 w-8 h-8 rounded-lg flex items-center justify-center text-xs font-bold uppercase',
              selected?.id === app.id ? 'bg-white/15 text-white' : appAvatarClass(app.name)
            ]">
              {{ app.name.charAt(0) }}
            </div>
            <!-- Info -->
            <div class="min-w-0">
              <div :class="['text-sm font-medium truncate', selected?.id === app.id ? 'text-white' : 'text-slate-900']">
                {{ app.name }}
              </div>
              <div :class="['text-[11px] mt-0.5 font-mono', selected?.id === app.id ? 'text-slate-400' : 'text-slate-400']">
                {{ app.code }}
              </div>
            </div>
            <!-- Active dot -->
            <div class="ml-auto shrink-0">
              <span v-if="app.status === 'active'" :class="['w-1.5 h-1.5 rounded-full inline-block', selected?.id === app.id ? 'bg-emerald-400' : 'bg-emerald-500']"></span>
              <span v-else class="w-1.5 h-1.5 rounded-full inline-block bg-slate-300"></span>
            </div>
          </div>
          <!-- Inline actions for selected app -->
          <div v-if="selected?.id === app.id && authStore.canAny('aplicaciones.editar', 'aplicaciones.eliminar')" class="bg-slate-800 border-b border-slate-700 px-3 py-2 flex items-center gap-2">
            <button v-if="authStore.can('aplicaciones.editar')" @click.stop="openEdit(selected)" title="Editar" class="p-1.5 rounded text-slate-300 hover:text-white hover:bg-slate-700 transition-colors"><Pencil :size="15" /></button>
            <div v-if="authStore.can('aplicaciones.editar') && authStore.can('aplicaciones.eliminar')" class="w-px h-4 bg-slate-600 shrink-0"></div>
            <button v-if="authStore.can('aplicaciones.eliminar')" @click.stop="confirmDelete(selected)" title="Eliminar" class="p-1.5 rounded text-red-400 hover:text-red-300 hover:bg-slate-700 transition-colors"><Trash2 :size="15" /></button>
          </div>
        </div>
      </div>

      <!-- List error -->
      <div v-if="error" class="px-3 py-2 text-xs text-red-600 border-t border-red-100 bg-red-50 shrink-0">
        {{ error }}
      </div>

      <!-- Footer: Nueva aplicación -->
      <div v-if="authStore.can('aplicaciones.crear')" class="p-3 border-t border-slate-200 shrink-0">
        <button @click="openCreate" class="w-full an-btn-primary text-sm py-2 flex items-center justify-center gap-1.5">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/></svg>
          Nueva Aplicación
        </button>
      </div>
    </aside>

    <!-- ══ RIGHT PANEL: Application detail ═════════════════════════════════ -->
    <main class="flex-1 overflow-y-auto bg-slate-50">

      <!-- Empty state -->
      <div v-if="!selected" class="flex flex-col items-center justify-center h-full gap-2 text-center px-8">
        <p class="text-sm font-medium text-slate-600">Selecciona una aplicación</p>
        <p class="text-xs text-slate-400">Elige una aplicación de la lista para ver y gestionar sus módulos.</p>
      </div>

      <div v-else>

        <!-- App hero header -->
        <div class="bg-gradient-to-r from-slate-900 to-slate-800 border-b border-slate-700 px-4 sm:px-8 py-4 sm:py-5">
          <div class="flex items-center gap-4">
            <!-- Back button — mobile only -->
            <button
              @click="selected = null"
              class="sm:hidden shrink-0 p-1 -ml-1 rounded text-slate-400 hover:text-white hover:bg-slate-700 transition-colors"
              aria-label="Volver"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7"/></svg>
            </button>
            <!-- Large avatar -->
            <div :class="['w-12 h-12 rounded-xl flex items-center justify-center text-lg font-bold uppercase shrink-0', appAvatarClass(selected.name)]">
              {{ selected.name.charAt(0) }}
            </div>
            <div>
              <div class="flex items-center gap-2.5">
                <h1 class="text-lg font-semibold text-white">{{ selected.name }}</h1>
                <span :class="selected.status === 'active' ? 'an-badge-active' : 'an-badge-inactive'">
                  {{ selected.status === 'active' ? 'Activa' : 'Inactiva' }}
                </span>
              </div>
              <p class="text-xs text-slate-400 mt-0.5 font-mono">{{ selected.code }}</p>
            </div>
          </div>
        </div>

        <!-- Modules section -->
        <div class="px-4 sm:px-8 py-6">

          <div class="flex items-center justify-between mb-4">
            <p class="text-sm text-slate-500">
              Módulos de <strong class="text-slate-700">{{ selected.name }}</strong>
            </p>
            <button v-if="authStore.can('modulos.crear')" @click="openCreateModule" class="an-btn-primary text-xs py-1.5 px-3">Nuevo Módulo</button>
          </div>

          <div v-if="detailError" class="mb-4 rounded border border-red-200 bg-red-50 px-4 py-2 text-sm text-red-700">{{ detailError }}</div>

          <div class="an-card overflow-hidden">
            <div v-if="loadingModules" class="py-16 text-center text-sm text-slate-400">
              Cargando módulos…
            </div>
            <div v-else-if="modules.length === 0" class="py-16 text-center">
              <p class="text-sm font-medium text-slate-600">Sin módulos configurados</p>
              <p class="text-xs text-slate-400 mt-1">Aún no se han creado módulos para esta aplicación.</p>
              <button v-if="authStore.can('modulos.crear')" @click="openCreateModule" class="mt-4 an-btn-primary text-xs py-1.5 px-4">Crear primer módulo</button>
            </div>
            <div v-else class="overflow-x-auto">
            <table class="w-full text-sm border-collapse" style="min-width:540px">
              <thead>
                <tr class="bg-slate-800">
                  <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-12 hidden sm:table-cell">#</th>
                  <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300">Nombre</th>
                  <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-32 hidden md:table-cell">Opción de menú</th>
                  <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-32 hidden md:table-cell">Sub-función</th>
                  <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-20">Estado</th>
                  <th v-if="authStore.canAny('modulos.editar', 'modulos.eliminar')" class="px-4 py-2.5 text-right text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-32">Acciones</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="m in modules" :key="m.id" class="border-b border-slate-100 last:border-0 hover:bg-slate-50 transition-colors">
                  <td class="px-4 py-3 text-slate-400 tabular-nums hidden sm:table-cell">{{ m.id }}</td>
                  <td class="px-4 py-3 font-medium text-slate-900">
                    {{ m.displayName || m.name }}
                    <p class="text-[11px] text-slate-400 mt-0.5 font-mono">{{ m.name }}</p>
                    <p v-if="m.description" class="text-[11px] text-slate-500 mt-0.5">{{ m.description }}</p>
                  </td>
                  <td class="px-4 py-3 text-slate-500 text-xs hidden md:table-cell">{{ m.menuOption || '—' }}</td>
                  <td class="px-4 py-3 text-slate-500 text-xs hidden md:table-cell">{{ m.subFunction || '—' }}</td>
                  <td class="px-4 py-3">
                    <span :class="m.status === 'active' ? 'an-badge-active' : 'an-badge-inactive'">
                      {{ m.status === 'active' ? 'Activo' : 'Inactivo' }}
                    </span>
                  </td>
                  <td v-if="authStore.canAny('modulos.editar', 'modulos.eliminar')" class="px-4 py-3">
                    <div class="flex items-center justify-end gap-2">
                      <button v-if="authStore.can('modulos.editar')" @click="openEditModule(m)" title="Editar" class="p-1.5 rounded text-slate-500 hover:text-slate-800 hover:bg-slate-100 transition-colors"><Pencil :size="15" /></button>
                      <button v-if="authStore.can('modulos.eliminar')" @click="confirmDeleteModule(m)" title="Eliminar" class="p-1.5 rounded text-red-400 hover:text-red-600 hover:bg-red-50 transition-colors"><Trash2 :size="15" /></button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- ══ MODALS ═══════════════════════════════════════════════════════════ -->

    <!-- Create / Edit application -->
    <Teleport to="body">
      <div v-if="modalOpen" class="an-overlay" @click.self="closeModal">
        <div class="an-modal">
          <div class="an-modal-header">
            <h2 class="text-base font-semibold text-slate-900">{{ editing ? 'Editar Aplicación' : 'Nueva Aplicación' }}</h2>
            <button @click="closeModal" class="an-close-btn">&times;</button>
          </div>
          <form @submit.prevent="submitForm" class="an-modal-body">
            <div v-if="formError" class="mb-4 rounded border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">{{ formError }}</div>

            <div class="mb-4">
              <label class="an-label">Código <span class="text-red-500">*</span></label>
              <input
                v-model="form.code"
                type="text"
                required
                maxlength="50"
                :disabled="saving || !!editing"
                class="an-input"
                :class="{ 'opacity-60 cursor-not-allowed': !!editing }"
                placeholder="Ej. ERP, CRM, POS"
              />
              <p v-if="editing" class="text-[11px] text-slate-400 mt-1">El código no puede cambiarse después de la creación.</p>
            </div>

            <div class="mb-4">
              <label class="an-label">Nombre <span class="text-red-500">*</span></label>
              <input v-model="form.name" type="text" required maxlength="255" :disabled="saving" class="an-input" placeholder="Ej. Sistema de Planificación" />
            </div>

            <div v-if="editing">
              <label class="an-label">Estado</label>
              <select v-model="form.status" :disabled="saving" class="an-input">
                <option value="active">Activa</option>
                <option value="inactive">Inactiva</option>
              </select>
            </div>

            <div class="an-modal-footer">
              <button type="button" @click="closeModal" class="an-btn-ghost">Cancelar</button>
              <button type="submit" :disabled="saving" class="an-btn-primary">{{ saving ? 'Guardando…' : 'Guardar' }}</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Delete application confirm -->
    <Teleport to="body">
      <div v-if="deleteTarget" class="an-overlay" @click.self="deleteTarget = null">
        <div class="an-modal an-modal--sm">
          <div class="an-modal-header">
            <h2 class="text-base font-semibold text-slate-900">Eliminar aplicación</h2>
            <button @click="deleteTarget = null" class="an-close-btn">&times;</button>
          </div>
          <div class="an-modal-body">
            <p class="text-sm text-slate-600">
              ¿Confirmas la eliminación de
              <span class="font-semibold text-slate-900">{{ deleteTarget.name }}</span>?
              Esta acción no se puede deshacer.
            </p>
            <div class="an-modal-footer">
              <button @click="deleteTarget = null" class="an-btn-ghost">Cancelar</button>
              <button @click="doDelete" :disabled="deleting" class="an-btn-danger">{{ deleting ? 'Eliminando…' : 'Eliminar' }}</button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Create / Edit module -->
    <Teleport to="body">
      <div v-if="moduleModalOpen" class="an-overlay" @click.self="closeModuleModal">
        <div class="an-modal">
          <div class="an-modal-header">
            <h2 class="text-base font-semibold text-slate-900">{{ editingModule ? 'Editar Módulo' : 'Nuevo Módulo' }}</h2>
            <button @click="closeModuleModal" class="an-close-btn">&times;</button>
          </div>
          <form @submit.prevent="submitModule" class="an-modal-body">
            <div v-if="moduleFormError" class="mb-4 rounded border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">{{ moduleFormError }}</div>
            <div class="mb-4">
              <label class="an-label">Clave del módulo <span class="text-red-500">*</span></label>
              <input v-model="moduleForm.name" required class="an-input font-mono" placeholder="Ej. usuarios.crear" :disabled="moduleSaving" />
              <p class="text-[11px] text-slate-400 mt-1">Identificador interno usado en el código para verificar permisos.</p>
            </div>
            <div class="mb-4">
              <label class="an-label">Nombre para mostrar</label>
              <input v-model="moduleForm.displayName" class="an-input" placeholder="Ej. Crear usuarios" :disabled="moduleSaving" />
              <p class="text-[11px] text-slate-400 mt-1">Nombre legible que se muestra en la interfaz.</p>
            </div>
            <div class="mb-4">
              <label class="an-label">Opción de menú</label>
              <input v-model="moduleForm.menuOption" class="an-input" placeholder="Ej. Usuarios" :disabled="moduleSaving" />
            </div>
            <div class="mb-4">
              <label class="an-label">Sub-función</label>
              <input v-model="moduleForm.subFunction" class="an-input" placeholder="Ej. Crear usuario" :disabled="moduleSaving" />
            </div>
            <div class="mb-4">
              <label class="an-label">Descripción</label>
              <textarea v-model="moduleForm.description" class="an-input" rows="2" placeholder="Descripción opcional" :disabled="moduleSaving"></textarea>
            </div>
            <div v-if="editingModule" class="mb-4">
              <label class="an-label">Estado</label>
              <select v-model="moduleForm.status" class="an-input" :disabled="moduleSaving">
                <option value="active">Activo</option>
                <option value="inactive">Inactivo</option>
              </select>
            </div>
            <div class="an-modal-footer">
              <button type="button" @click="closeModuleModal" class="an-btn-ghost">Cancelar</button>
              <button type="submit" :disabled="moduleSaving" class="an-btn-primary">{{ moduleSaving ? 'Guardando…' : 'Guardar' }}</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Delete module confirm -->
    <Teleport to="body">
      <div v-if="deleteModuleTarget" class="an-overlay" @click.self="deleteModuleTarget = null">
        <div class="an-modal an-modal--sm">
          <div class="an-modal-header">
            <h2 class="text-base font-semibold text-slate-900">Eliminar módulo</h2>
            <button @click="deleteModuleTarget = null" class="an-close-btn">&times;</button>
          </div>
          <div class="an-modal-body">
            <p class="text-sm text-slate-600">
              ¿Eliminar el módulo
              <span class="font-semibold text-slate-900">{{ deleteModuleTarget.name }}</span>?
              Se eliminarán también sus asignaciones a roles.
            </p>
            <div class="an-modal-footer">
              <button @click="deleteModuleTarget = null" class="an-btn-ghost">Cancelar</button>
              <button @click="doDeleteModule" :disabled="deletingModule" class="an-btn-danger">{{ deletingModule ? 'Eliminando…' : 'Eliminar' }}</button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useAuthStore } from '@/store/auth'
import api from '@/services/api'
import { teError } from '@/i18n'
import { Pencil, Trash2 } from 'lucide-vue-next'

const authStore = useAuthStore()

// Avatar colour palette (deterministic by first letter).
const AVATAR_COLOURS = [
  'bg-violet-100 text-violet-700',
  'bg-sky-100 text-sky-700',
  'bg-emerald-100 text-emerald-700',
  'bg-rose-100 text-rose-700',
  'bg-amber-100 text-amber-700',
  'bg-teal-100 text-teal-700',
  'bg-indigo-100 text-indigo-700',
  'bg-fuchsia-100 text-fuchsia-700',
]
function appAvatarClass(name = '') {
  const idx = (name.charCodeAt(0) || 0) % AVATAR_COLOURS.length
  return AVATAR_COLOURS[idx]
}

// ── Applications list ─────────────────────────────────────────────────────────
const apps     = ref([])
const loading  = ref(true)
const error    = ref('')
const search   = ref('')
const selected = ref(null)

const filteredApps = computed(() =>
  apps.value.filter(a =>
    a.name.toLowerCase().includes(search.value.toLowerCase()) ||
    a.code.toLowerCase().includes(search.value.toLowerCase())
  )
)

async function loadApps() {
  loading.value = true
  error.value   = ''
  try {
    const res = await api.get('/applications', { params: { pageSize: 200 } })
    apps.value = res.data?.data ?? []
  } catch (e) {
    error.value = teError(e)
  } finally {
    loading.value = false
  }
}

function selectApp(app) {
  if (selected.value?.id === app.id) return
  selected.value = app
  detailError.value = ''
  loadModules(app.id)
}

onMounted(loadApps)

// ── Create / Edit application ─────────────────────────────────────────────────
const modalOpen = ref(false)
const saving    = ref(false)
const formError = ref('')
const editing   = ref(null)
const form      = reactive({ code: '', name: '', status: 'active' })

function openCreate() {
  editing.value   = null
  form.code       = ''
  form.name       = ''
  form.status     = 'active'
  formError.value = ''
  modalOpen.value = true
}

function openEdit(app) {
  editing.value   = app
  form.code       = app.code
  form.name       = app.name
  form.status     = app.status
  formError.value = ''
  modalOpen.value = true
}

function closeModal() {
  if (!saving.value) modalOpen.value = false
}

async function submitForm() {
  formError.value = ''
  saving.value    = true
  try {
    if (editing.value) {
      const res = await api.put(`/applications/${editing.value.id}`, {
        name:   form.name.trim(),
        status: form.status,
      })
      const updated = res.data?.data ?? res.data
      const idx = apps.value.findIndex(a => a.id === updated.id)
      if (idx !== -1) apps.value[idx] = updated
      if (selected.value?.id === updated.id) selected.value = updated
    } else {
      const res = await api.post('/applications', {
        code: form.code.trim(),
        name: form.name.trim(),
      })
      const newApp = res.data?.data ?? res.data
      apps.value.unshift(newApp)
      selectApp(newApp)
    }
    modalOpen.value = false
  } catch (e) {
    formError.value = teError(e)
  } finally {
    saving.value = false
  }
}

// ── Delete application ────────────────────────────────────────────────────────
const deleteTarget = ref(null)
const deleting     = ref(false)

function confirmDelete(app) {
  deleteTarget.value = app
}

async function doDelete() {
  deleting.value = true
  try {
    await api.delete(`/applications/${deleteTarget.value.id}`)
    apps.value = apps.value.filter(a => a.id !== deleteTarget.value.id)
    if (selected.value?.id === deleteTarget.value.id) selected.value = null
    deleteTarget.value = null
  } catch (e) {
    error.value = teError(e)
    deleteTarget.value = null
  } finally {
    deleting.value = false
  }
}

// ── Modules ───────────────────────────────────────────────────────────────────
const modules        = ref([])
const loadingModules = ref(false)
const detailError    = ref('')

async function loadModules(appId) {
  loadingModules.value = true
  detailError.value = ''
  try {
    const res = await api.get(`/applications/${appId}/modules`, { params: { pageSize: 500 } })
    modules.value = res.data?.data ?? []
  } catch (e) {
    detailError.value = teError(e)
  } finally {
    loadingModules.value = false
  }
}

// ── Module CRUD ───────────────────────────────────────────────────────────────
const moduleModalOpen  = ref(false)
const moduleSaving     = ref(false)
const moduleFormError  = ref('')
const editingModule    = ref(null)
const moduleForm       = reactive({ name: '', displayName: '', menuOption: '', subFunction: '', description: '', status: 'active' })

function openCreateModule() {
  editingModule.value   = null
  moduleFormError.value = ''
  Object.assign(moduleForm, { name: '', displayName: '', menuOption: '', subFunction: '', description: '', status: 'active' })
  moduleModalOpen.value = true
}

function openEditModule(m) {
  editingModule.value   = m
  moduleFormError.value = ''
  Object.assign(moduleForm, {
    name:        m.name,
    displayName: m.displayName ?? '',
    menuOption:  m.menuOption ?? '',
    subFunction: m.subFunction ?? '',
    description: m.description ?? '',
    status:      m.status,
  })
  moduleModalOpen.value = true
}

function closeModuleModal() {
  if (!moduleSaving.value) moduleModalOpen.value = false
}

async function submitModule() {
  moduleFormError.value = ''
  moduleSaving.value    = true
  try {
    if (editingModule.value) {
      await api.put(`/applications/${selected.value.id}/modules/${editingModule.value.id}`, {
        name:        moduleForm.name.trim(),
        displayName: moduleForm.displayName.trim() || null,
        menuOption:  moduleForm.menuOption.trim() || null,
        subFunction: moduleForm.subFunction.trim() || null,
        description: moduleForm.description.trim() || null,
        status:      moduleForm.status,
      })
    } else {
      await api.post(`/applications/${selected.value.id}/modules`, {
        name:        moduleForm.name.trim(),
        displayName: moduleForm.displayName.trim() || null,
        menuOption:  moduleForm.menuOption.trim() || null,
        subFunction: moduleForm.subFunction.trim() || null,
        description: moduleForm.description.trim() || null,
      })
    }
    moduleModalOpen.value = false
    await loadModules(selected.value.id)
  } catch (e) {
    moduleFormError.value = teError(e)
  } finally {
    moduleSaving.value = false
  }
}

// ── Delete module ─────────────────────────────────────────────────────────────
const deleteModuleTarget = ref(null)
const deletingModule     = ref(false)

function confirmDeleteModule(m) { deleteModuleTarget.value = m }

async function doDeleteModule() {
  deletingModule.value = true
  try {
    await api.delete(`/applications/${selected.value.id}/modules/${deleteModuleTarget.value.id}`)
    deleteModuleTarget.value = null
    await loadModules(selected.value.id)
  } catch (e) {
    detailError.value = teError(e)
    deleteModuleTarget.value = null
  } finally {
    deletingModule.value = false
  }
}
</script>
