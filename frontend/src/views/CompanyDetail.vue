<template>
  <div>

    <!-- Breadcrumb -->
    <nav class="mb-6 flex items-center gap-2 text-sm text-slate-500">
      <RouterLink to="/companies" class="hover:text-slate-900 transition-colors">Empresas</RouterLink>
      <span class="text-slate-300">/</span>
      <span class="text-slate-900 font-medium">{{ company?.name ?? '…' }}</span>
    </nav>

    <!-- Company header card -->
    <div class="an-card px-6 py-5 mb-6" v-if="company">
      <div class="flex items-start justify-between gap-4">
        <div>
          <h1 class="text-xl font-semibold text-slate-900">{{ company.name }}</h1>
          <span :class="company.status === 'active' ? 'an-badge-active' : 'an-badge-inactive'" class="mt-2 inline-block">
            {{ company.status === 'active' ? 'Activa' : 'Inactiva' }}
          </span>
        </div>
        <button @click="openEditCompany" class="an-btn-ghost shrink-0">Editar empresa</button>
      </div>
    </div>

    <!-- Error banner -->
    <div v-if="error" class="mb-5 rounded border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">{{ error }}</div>

    <!-- Tabs -->
    <div class="border-b border-slate-200 mb-6">
      <nav class="flex gap-6">
        <button
          v-for="tab in tabs" :key="tab.key"
          @click="activeTab = tab.key"
          :class="[
            'pb-3 text-sm font-medium border-b-2 transition-colors',
            activeTab === tab.key
              ? 'border-slate-900 text-slate-900'
              : 'border-transparent text-slate-500 hover:text-slate-700'
          ]"
        >
          {{ tab.label }}
        </button>
      </nav>
    </div>

    <!-- ── Tab: Sucursales ───────────────────────────────────────────────── -->
    <div v-if="activeTab === 'branches'">
      <div class="flex items-center justify-between mb-5">
        <p class="text-sm text-slate-500">Sucursales registradas en esta empresa.</p>
        <button @click="openCreateBranch" class="an-btn-primary">Nueva Sucursal</button>
      </div>

      <div class="an-card">
        <div v-if="loadingBranches" class="py-16 text-center text-sm text-slate-400">Cargando…</div>
        <div v-else-if="branches.length === 0" class="py-16 text-center">
          <p class="text-sm text-slate-500">No hay sucursales registradas.</p>
          <button @click="openCreateBranch" class="mt-4 an-btn-primary">Registrar primera sucursal</button>
        </div>
        <table v-else class="an-table">
          <thead>
            <tr>
              <th>Código</th>
              <th>Nombre</th>
              <th>Teléfono</th>
              <th>Email</th>
              <th class="w-24">Estado</th>
              <th class="w-32 text-right">Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="b in branches" :key="b.id">
              <td class="font-mono text-xs text-slate-500">{{ b.code }}</td>
              <td class="font-medium text-slate-900">{{ b.name }}</td>
              <td class="text-slate-500">{{ b.phoneNumber ?? '—' }}</td>
              <td class="text-slate-500">{{ b.email ?? '—' }}</td>
              <td>
                <span :class="b.status === 'active' ? 'an-badge-active' : 'an-badge-inactive'">
                  {{ b.status === 'active' ? 'Activa' : 'Inactiva' }}
                </span>
              </td>
              <td class="text-right space-x-3">
                <button @click="openEditBranch(b)" class="text-sm text-slate-500 hover:text-slate-800 transition-colors">Editar</button>
                <button @click="confirmDeleteBranch(b)" class="text-sm text-red-500 hover:text-red-700 transition-colors">Eliminar</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- ── Tab: Roles ────────────────────────────────────────────────────── -->
    <div v-if="activeTab === 'roles'">
      <div class="flex items-center justify-between mb-5">
        <p class="text-sm text-slate-500">Roles disponibles en el sistema para asignar a usuarios de esta empresa.</p>
        <RouterLink to="/roles" class="an-btn-ghost">Gestionar Roles</RouterLink>
      </div>

      <div class="an-card">
        <div v-if="loadingRoles" class="py-16 text-center text-sm text-slate-400">Cargando…</div>
        <div v-else-if="roles.length === 0" class="py-16 text-center text-sm text-slate-500">
          No hay roles creados. <RouterLink to="/roles" class="text-slate-900 underline underline-offset-2">Crear rol</RouterLink>
        </div>
        <table v-else class="an-table">
          <thead>
            <tr>
              <th class="w-16">#</th>
              <th>Nombre del rol</th>
              <th class="w-28">Estado</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="r in roles" :key="r.id">
              <td class="text-slate-400 tabular-nums">{{ r.id }}</td>
              <td class="font-medium text-slate-900">{{ r.name }}</td>
              <td>
                <span :class="r.status === 'active' ? 'an-badge-active' : 'an-badge-inactive'">
                  {{ r.status === 'active' ? 'Activo' : 'Inactivo' }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- ══ Modals ══════════════════════════════════════════════════════════ -->

    <!-- Edit company -->
    <Teleport to="body">
      <div v-if="editCompanyOpen" class="an-overlay" @click.self="editCompanyOpen = false">
        <div class="an-modal">
          <div class="an-modal-header">
            <h2 class="text-base font-semibold text-slate-900">Editar empresa</h2>
            <button @click="editCompanyOpen = false" class="an-close-btn">&times;</button>
          </div>
          <form @submit.prevent="submitEditCompany" class="an-modal-body">
            <div v-if="editCompanyError" class="mb-4 rounded border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">{{ editCompanyError }}</div>
            <div class="mb-4">
              <label class="an-label">Nombre</label>
              <input v-model="editCompanyForm.name" type="text" required maxlength="255" :disabled="editCompanySaving" class="an-input" />
            </div>
            <div>
              <label class="an-label">Estado</label>
              <select v-model="editCompanyForm.status" :disabled="editCompanySaving" class="an-input">
                <option value="active">Activa</option>
                <option value="inactive">Inactiva</option>
              </select>
            </div>
            <div class="an-modal-footer">
              <button type="button" @click="editCompanyOpen = false" class="an-btn-ghost">Cancelar</button>
              <button type="submit" :disabled="editCompanySaving" class="an-btn-primary">{{ editCompanySaving ? 'Guardando…' : 'Guardar' }}</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Create / Edit branch modal -->
    <Teleport to="body">
      <div v-if="branchModalOpen" class="an-overlay" @click.self="closeBranchModal">
        <div class="an-modal">
          <div class="an-modal-header">
            <h2 class="text-base font-semibold text-slate-900">{{ editingBranch ? 'Editar sucursal' : 'Nueva sucursal' }}</h2>
            <button @click="closeBranchModal" class="an-close-btn">&times;</button>
          </div>
          <form @submit.prevent="submitBranch" class="an-modal-body">
            <div v-if="branchFormError" class="mb-4 rounded border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">{{ branchFormError }}</div>

            <div class="grid grid-cols-2 gap-4 mb-4">
              <div>
                <label class="an-label">Código <span class="text-red-500">*</span></label>
                <input v-model="branchForm.code" type="text" required maxlength="50" :disabled="branchSaving || !!editingBranch" class="an-input" placeholder="HQ-001" />
              </div>
              <div>
                <label class="an-label">Nombre <span class="text-red-500">*</span></label>
                <input v-model="branchForm.name" type="text" required maxlength="150" :disabled="branchSaving" class="an-input" placeholder="Casa Matriz" />
              </div>
            </div>

            <div class="mb-4">
              <label class="an-label">Dirección</label>
              <input v-model="branchForm.address" type="text" :disabled="branchSaving" class="an-input" placeholder="Calle 123, Ciudad" />
            </div>

            <div class="grid grid-cols-2 gap-4 mb-4">
              <div>
                <label class="an-label">Teléfono</label>
                <input v-model="branchForm.phoneNumber" type="text" :disabled="branchSaving" class="an-input" placeholder="+1 555 0000" />
              </div>
              <div>
                <label class="an-label">Email</label>
                <input v-model="branchForm.email" type="email" :disabled="branchSaving" class="an-input" placeholder="sucursal@empresa.com" />
              </div>
            </div>

            <div v-if="editingBranch">
              <label class="an-label">Estado</label>
              <select v-model="branchForm.status" :disabled="branchSaving" class="an-input">
                <option value="active">Activa</option>
                <option value="inactive">Inactiva</option>
              </select>
            </div>

            <div class="an-modal-footer">
              <button type="button" @click="closeBranchModal" class="an-btn-ghost">Cancelar</button>
              <button type="submit" :disabled="branchSaving" class="an-btn-primary">{{ branchSaving ? 'Guardando…' : 'Guardar' }}</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Delete branch confirm -->
    <Teleport to="body">
      <div v-if="deleteBranchTarget" class="an-overlay" @click.self="deleteBranchTarget = null">
        <div class="an-modal an-modal--sm">
          <div class="an-modal-header">
            <h2 class="text-base font-semibold text-slate-900">Confirmar eliminación</h2>
            <button @click="deleteBranchTarget = null" class="an-close-btn">&times;</button>
          </div>
          <div class="an-modal-body">
            <p class="text-sm text-slate-600">
              ¿Estás seguro de que deseas eliminar la sucursal
              <span class="font-semibold text-slate-900">{{ deleteBranchTarget.name }}</span>?
              Esta acción no se puede deshacer.
            </p>
            <div class="an-modal-footer">
              <button @click="deleteBranchTarget = null" class="an-btn-ghost">Cancelar</button>
              <button @click="doDeleteBranch" :disabled="deletingBranch" class="an-btn-danger">
                {{ deletingBranch ? 'Eliminando…' : 'Eliminar' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { companyService } from '@/services/companyService'
import { branchService } from '@/services/branchService'
import api from '@/services/api'
import { teError } from '@/i18n'

const route = useRoute()
const companyId = Number(route.params.id)

// ── State ─────────────────────────────────────────────────────────────────────
const company       = ref(null)
const branches      = ref([])
const roles         = ref([])
const error         = ref('')
const loadingBranches = ref(true)
const loadingRoles    = ref(false)
const activeTab     = ref('branches')

const tabs = [
  { key: 'branches', label: 'Sucursales' },
  { key: 'roles',    label: 'Roles' },
]

// ── Load ──────────────────────────────────────────────────────────────────────
async function loadCompany() {
  try {
    const res = await companyService.getById(companyId)
    company.value = res.data
  } catch (e) {
    error.value = teError(e)
  }
}

async function loadBranches() {
  loadingBranches.value = true
  try {
    const res = await branchService.getAll(companyId)
    branches.value = res.data ?? []
  } catch (e) {
    error.value = teError(e)
  } finally {
    loadingBranches.value = false
  }
}

async function loadRoles() {
  loadingRoles.value = true
  try {
    const res = await api.get('/roles', { params: { pageSize: 200 } })
    roles.value = res.data?.data ?? []
  } catch {
    roles.value = []
  } finally {
    loadingRoles.value = false
  }
}

// Watch tab change to lazy-load roles
watch(activeTab, (tab) => {
  if (tab === 'roles' && roles.value.length === 0) loadRoles()
})

onMounted(() => {
  loadCompany()
  loadBranches()
})

// ── Edit Company ──────────────────────────────────────────────────────────────
const editCompanyOpen   = ref(false)
const editCompanySaving = ref(false)
const editCompanyError  = ref('')
const editCompanyForm   = reactive({ name: '', status: 'active' })

function openEditCompany() {
  editCompanyForm.name   = company.value.name
  editCompanyForm.status = company.value.status
  editCompanyError.value = ''
  editCompanyOpen.value  = true
}

async function submitEditCompany() {
  editCompanyError.value = ''
  editCompanySaving.value = true
  try {
    const res = await companyService.update(companyId, {
      name:   editCompanyForm.name.trim(),
      status: editCompanyForm.status,
    })
    company.value = res.data
    editCompanyOpen.value = false
  } catch (e) {
    editCompanyError.value = teError(e)
  } finally {
    editCompanySaving.value = false
  }
}

// ── Branch CRUD ───────────────────────────────────────────────────────────────
const branchModalOpen   = ref(false)
const branchSaving      = ref(false)
const branchFormError   = ref('')
const editingBranch     = ref(null)
const deleteBranchTarget = ref(null)
const deletingBranch    = ref(false)

const branchForm = reactive({
  code: '', name: '', address: '', phoneNumber: '', email: '', status: 'active',
})

function openCreateBranch() {
  editingBranch.value  = null
  branchFormError.value = ''
  Object.assign(branchForm, { code: '', name: '', address: '', phoneNumber: '', email: '', status: 'active' })
  branchModalOpen.value = true
}

function openEditBranch(b) {
  editingBranch.value  = b
  branchFormError.value = ''
  Object.assign(branchForm, {
    code:        b.code,
    name:        b.name,
    address:     b.address ?? '',
    phoneNumber: b.phoneNumber ?? '',
    email:       b.email ?? '',
    status:      b.status,
  })
  branchModalOpen.value = true
}

function closeBranchModal() {
  if (!branchSaving.value) branchModalOpen.value = false
}

async function submitBranch() {
  branchFormError.value = ''
  branchSaving.value = true

  const payload = {
    name:        branchForm.name.trim(),
    address:     branchForm.address.trim() || undefined,
    phoneNumber: branchForm.phoneNumber.trim() || undefined,
    email:       branchForm.email.trim() || undefined,
  }

  try {
    if (editingBranch.value) {
      const res = await branchService.update(companyId, editingBranch.value.id, {
        ...payload, status: branchForm.status,
      })
      const idx = branches.value.findIndex(b => b.id === editingBranch.value.id)
      if (idx !== -1) branches.value[idx] = res.data
    } else {
      const res = await branchService.create(companyId, { ...payload, code: branchForm.code.trim() })
      branches.value.unshift(res.data)
    }
    branchModalOpen.value = false
  } catch (e) {
    branchFormError.value = teError(e)
  } finally {
    branchSaving.value = false
  }
}

function confirmDeleteBranch(b) {
  deleteBranchTarget.value = b
}

async function doDeleteBranch() {
  if (!deleteBranchTarget.value) return
  deletingBranch.value = true
  try {
    await branchService.remove(companyId, deleteBranchTarget.value.id)
    branches.value = branches.value.filter(b => b.id !== deleteBranchTarget.value.id)
    deleteBranchTarget.value = null
  } catch (e) {
    error.value = teError(e)
  } finally {
    deletingBranch.value = false
  }
}
</script>
