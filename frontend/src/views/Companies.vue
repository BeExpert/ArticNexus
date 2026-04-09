<template>
  <div class="flex h-full overflow-hidden">

    <!-- ══ LEFT PANEL: Company list ════════════════════════════════════════ -->
    <aside :class="['shrink-0 flex flex-col border-r border-slate-200 bg-white', 'w-full sm:w-72', selected ? 'hidden sm:flex' : 'flex']">

      <!-- Header — admin mode banner (visible when user can create/edit/delete companies) -->
      <div v-if="authStore.canAny('empresas.crear', 'empresas.editar', 'empresas.eliminar')" class="px-4 py-2 bg-amber-50 border-b border-amber-200 flex items-center gap-1.5 shrink-0">
        <svg class="w-3 h-3 text-amber-600 shrink-0" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 1a.75.75 0 0 1 .68.433l1.63 3.305 3.647.53a.75.75 0 0 1 .416 1.279l-2.638 2.571.623 3.63a.75.75 0 0 1-1.088.79L10 11.876l-3.27 1.72a.75.75 0 0 1-1.088-.79l.622-3.63L3.627 6.547a.75.75 0 0 1 .416-1.28l3.646-.529L9.32 1.433A.75.75 0 0 1 10 1Z" clip-rule="evenodd"/>
        </svg>
        <span class="text-[11px] font-semibold text-amber-700 uppercase tracking-wide">Modo administrador</span>
      </div>

      <!-- Header — normal -->
      <div class="px-4 py-3.5 border-b border-slate-700 bg-slate-800 flex items-center justify-between shrink-0">
        <div class="flex items-center gap-2">
          <span class="text-sm font-semibold text-white">Empresas</span>
          <span class="text-xs text-slate-400 tabular-nums">({{ companies.length }})</span>
        </div>
      </div>

      <!-- Search -->
      <div class="px-3 py-2 border-b border-slate-100 shrink-0">
        <input
          v-model="search"
          type="text"
          placeholder="Buscar empresa…"
          class="an-input text-xs py-1.5"
        />
      </div>

      <!-- List -->
      <div class="flex-1 overflow-y-auto">
        <div v-if="loadingCompanies" class="py-12 flex flex-col items-center gap-2">
          <div class="w-5 h-5 rounded-full border-2 border-slate-200 border-t-slate-500 animate-spin"></div>
          <span class="text-xs text-slate-400">Cargando…</span>
        </div>
        <div v-else-if="filteredCompanies.length === 0" class="py-12 text-center">
          <p class="text-xs text-slate-400">No hay resultados.</p>
        </div>
        <div v-for="c in filteredCompanies" :key="c.id">
          <!-- Company row -->
          <div
            @click="selectCompany(c)"
            :class="[
              'w-full text-left px-3 py-2.5 border-b border-slate-100 transition-colors flex items-center gap-3 cursor-pointer select-none',
              selected?.id === c.id ? 'bg-slate-900 border-slate-800' : 'hover:bg-slate-50'
            ]"
          >
            <!-- Avatar -->
            <div :class="[
              'shrink-0 w-8 h-8 rounded-lg flex items-center justify-center text-xs font-bold uppercase',
              selected?.id === c.id ? 'bg-white/15 text-white' : companyAvatarClass(c.name)
            ]">
              {{ c.name.charAt(0) }}
            </div>
            <!-- Info -->
            <div class="min-w-0">
              <div :class="['text-sm font-medium truncate', selected?.id === c.id ? 'text-white' : 'text-slate-900']">
                {{ c.name }}
              </div>
              <div :class="['text-[11px] mt-0.5', selected?.id === c.id ? 'text-slate-400' : 'text-slate-400']">
                {{ c.status === 'active' ? 'Activa' : 'Inactiva' }}
              </div>
            </div>
            <!-- Active dot -->
            <div class="ml-auto shrink-0">
              <span v-if="c.status === 'active'" :class="['w-1.5 h-1.5 rounded-full inline-block', selected?.id === c.id ? 'bg-emerald-400' : 'bg-emerald-500']"></span>
              <span v-else class="w-1.5 h-1.5 rounded-full inline-block bg-slate-300"></span>
            </div>
          </div>
          <!-- Inline actions — visible only for selected company when user has edit/delete perms -->
          <div v-if="selected?.id === c.id && authStore.canAny('empresas.editar', 'empresas.eliminar')" class="bg-slate-800 border-b border-slate-700 px-3 py-2 flex items-center gap-2">
            <button v-if="authStore.can('empresas.editar')" @click.stop="openEditCompany" class="flex-1 text-center text-xs font-medium text-slate-300 hover:text-white hover:bg-slate-700 rounded py-1.5 transition-colors">Editar</button>
            <div v-if="authStore.can('empresas.editar') && authStore.can('empresas.eliminar')" class="w-px h-4 bg-slate-600 shrink-0"></div>
            <button v-if="authStore.can('empresas.eliminar')" @click.stop="confirmDeleteCompany(selected)" class="flex-1 text-center text-xs font-medium text-red-400 hover:text-red-300 hover:bg-slate-700 rounded py-1.5 transition-colors">Eliminar</button>
          </div>
        </div>
      </div>

      <!-- List error -->
      <div v-if="listError" class="px-3 py-2 text-xs text-red-600 border-t border-red-100 bg-red-50 shrink-0">
        {{ listError }}
      </div>

      <!-- Footer: Nueva empresa -->
      <div v-if="authStore.can('empresas.crear')" class="p-3 border-t border-slate-200 shrink-0">
        <button @click="openCreate" class="w-full an-btn-primary text-sm py-2 flex items-center justify-center gap-1.5">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/></svg>
          Nueva Empresa
        </button>
      </div>
    </aside>

    <!-- ══ RIGHT PANEL: Company detail ═════════════════════════════════════ -->
    <main :class="['flex-1 overflow-y-auto bg-slate-50', selected ? 'block' : 'hidden sm:block']">

      <!-- Empty state -->
      <div v-if="!selected" class="flex flex-col items-center justify-center h-full gap-2 text-center px-8">
        <p class="text-sm font-medium text-slate-600">Selecciona una empresa</p>
        <p class="text-xs text-slate-400">Elige una empresa de la lista para ver y gestionar sus sucursales y roles.</p>
      </div>

      <div v-else>

        <!-- Company hero header -->
        <div class="bg-gradient-to-r from-slate-900 to-slate-800 border-b border-slate-700 px-4 sm:px-8 py-4 sm:py-5">
          <div class="flex items-center justify-between gap-4">
            <div class="flex items-center gap-3">
              <!-- Back button — mobile only -->
              <button
                @click="selected = null"
                class="sm:hidden shrink-0 p-1 -ml-1 rounded text-slate-400 hover:text-white hover:bg-slate-700 transition-colors"
                aria-label="Volver"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7"/></svg>
              </button>
              <!-- Large avatar -->
              <div :class="['w-12 h-12 rounded-xl flex items-center justify-center text-lg font-bold uppercase shrink-0', companyAvatarClass(selected.name)]">
                {{ selected.name.charAt(0) }}
              </div>
              <div>
                <div class="flex items-center gap-2.5">
                  <h1 class="text-lg font-semibold text-white">{{ selected.name }}</h1>
                  <span :class="selected.status === 'active' ? 'an-badge-active' : 'an-badge-inactive'">
                    {{ selected.status === 'active' ? 'Activa' : 'Inactiva' }}
                  </span>
                </div>
                <p class="text-xs text-slate-400 mt-0.5 tabular-nums">
                  {{ branches.length }} {{ branches.length === 1 ? 'sucursal' : 'sucursales' }}
                </p>
              </div>
            </div>
          </div>

          <!-- Tabs — pill style inside header -->
          <div class="flex gap-1 mt-4">
            <button
              v-for="tab in visibleTabs" :key="tab.key"
              @click="activeTab = tab.key"
              :class="[
                'px-3.5 py-1.5 text-xs font-medium rounded-full transition-colors',
                activeTab === tab.key
                  ? 'bg-white text-slate-900'
                  : 'text-slate-400 hover:bg-slate-700 hover:text-white'
              ]"
            >
              {{ tab.label }}
            </button>
          </div>
        </div>

        <!-- Tab content -->
        <div class="px-4 sm:px-8 py-6">

        <!-- ── Sucursales tab ─────────────────────────────────────────────── -->
        <div v-if="activeTab === 'branches'">
          <div class="flex items-center justify-between mb-4">
            <p class="text-sm text-slate-500">
              Sucursales de <strong class="text-slate-700">{{ selected.name }}</strong>
            </p>
            <button v-if="authStore.can('sucursales.crear')" @click="openCreateBranch" class="an-btn-primary text-xs py-1.5 px-3">Nueva Sucursal</button>
          </div>

          <div v-if="detailError" class="mb-4 rounded border border-red-200 bg-red-50 px-4 py-2 text-sm text-red-700">{{ detailError }}</div>

          <div class="an-card overflow-hidden">
            <div v-if="loadingBranches" class="py-16 text-center text-sm text-slate-400">
              Cargando sucursales…
            </div>
            <div v-else-if="branches.length === 0" class="py-16 text-center">
              <p class="text-sm font-medium text-slate-600">Sin sucursales registradas</p>
              <p class="text-xs text-slate-400 mt-1">Esta empresa aún no tiene sucursales.</p>
              <button v-if="authStore.can('sucursales.crear')" @click="openCreateBranch" class="mt-4 an-btn-primary text-xs py-1.5 px-4">Registrar primera sucursal</button>
            </div>
            <div class="overflow-x-auto"><table class="w-full text-sm border-collapse" style="min-width:480px">
              <thead>
                <tr class="bg-slate-800">
                  <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-28">Código</th>
                  <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300">Nombre</th>
                  <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-36">Teléfono</th>
                  <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-52">Email</th>
                  <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-24">Estado</th>
                  <th v-if="authStore.canAny('sucursales.editar', 'sucursales.eliminar')" class="px-4 py-2.5 text-right text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-40">Acciones</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="b in branches" :key="b.id"
                  class="border-b border-slate-100 last:border-0 hover:bg-slate-50 transition-colors"
                >
                  <td class="px-4 py-3 font-mono text-xs text-slate-400">{{ b.code }}</td>
                  <td class="px-4 py-3 font-medium text-slate-900">{{ b.name }}</td>
                  <td class="px-4 py-3 text-slate-500">{{ b.phoneNumber ?? '—' }}</td>
                  <td class="px-4 py-3 text-slate-500">{{ b.email ?? '—' }}</td>
                  <td class="px-4 py-3">
                    <span :class="b.status === 'active' ? 'an-badge-active' : 'an-badge-inactive'">
                      {{ b.status === 'active' ? 'Activa' : 'Inactiva' }}
                    </span>
                  </td>
                  <td v-if="authStore.canAny('sucursales.editar', 'sucursales.eliminar')" class="px-4 py-3">
                    <div class="flex items-center justify-end gap-2">
                      <button v-if="authStore.can('sucursales.editar')" @click="openEditBranch(b)" class="px-3 py-1 text-xs font-medium rounded border border-slate-200 text-slate-600 hover:bg-slate-100 hover:border-slate-300 transition-colors whitespace-nowrap">Editar</button>
                      <button v-if="authStore.can('sucursales.eliminar')" @click="confirmDeleteBranch(b)" class="px-3 py-1 text-xs font-medium rounded border border-red-100 text-red-500 hover:bg-red-50 hover:border-red-200 transition-colors whitespace-nowrap">Eliminar</button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
            </div>
          </div>
        </div>

        <!-- ── Roles tab ──────────────────────────────────────────────────── -->
        <div v-if="activeTab === 'roles'">
          <div class="flex items-center justify-between mb-4">
            <p class="text-sm text-slate-500">Roles disponibles para asignar a usuarios de esta empresa.</p>
            <button v-if="authStore.can('roles.crear')" @click="openCreateRole" class="an-btn-primary text-xs py-1.5 px-3">Nuevo Rol</button>
          </div>

          <div class="an-card overflow-hidden">
            <div v-if="loadingRoles" class="py-16 text-center text-sm text-slate-400">Cargando roles…</div>
            <div v-else-if="roles.length === 0" class="py-16 text-center">
              <p class="text-sm font-medium text-slate-600">Sin roles configurados</p>
              <p class="text-xs text-slate-400 mt-1">Aún no se han creado roles para asignar a los usuarios.</p>
            </div>
            <div class="overflow-x-auto"><table class="w-full text-sm border-collapse" style="min-width:480px">
              <thead>
                <tr class="bg-slate-800">
                  <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-16">#</th>
                  <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300">Nombre del rol</th>
                  <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-40">Aplicacion</th>
                  <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-24">Estado</th>
                  <th v-if="authStore.can('roles.asignar_modulos')" class="px-4 py-2.5 text-right text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-48">Acciones</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="r in roles" :key="r.id" class="border-b border-slate-100 last:border-0 hover:bg-slate-50 transition-colors">
                  <td class="px-4 py-3 text-slate-400 tabular-nums">{{ r.id }}</td>
                  <td class="px-4 py-3 font-medium text-slate-900">{{ r.name }}</td>
                  <td class="px-4 py-3 text-slate-500 text-xs">{{ appName(r.applicationId) }}</td>
                  <td class="px-4 py-3">
                    <span :class="r.status === 'active' ? 'an-badge-active' : 'an-badge-inactive'">
                      {{ r.status === 'active' ? 'Activo' : 'Inactivo' }}
                    </span>
                  </td>
                  <td v-if="authStore.can('roles.asignar_modulos')" class="px-4 py-3">
                    <div class="flex items-center justify-end gap-2">
                      <button @click="openRoleModules(r)" class="px-3 py-1 text-xs font-medium rounded border border-slate-200 text-slate-600 hover:bg-slate-100 hover:border-slate-300 transition-colors whitespace-nowrap">Modulos</button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
            </div>
          </div>
        </div>

        <!-- ── Personas tab ───────────────────────────────────────────────── -->
        <div v-if="activeTab === 'personas'">
          <div class="flex items-center justify-between mb-4">
            <p class="text-sm text-slate-500">
              Personas asignadas a <strong class="text-slate-700">{{ selected.name }}</strong>
            </p>
            <button v-if="authStore.can('personas.crear')" @click="openAddPerson" class="an-btn-primary text-xs py-1.5 px-3">Agregar Persona</button>
          </div>

          <div class="an-card overflow-hidden">
            <div v-if="loadingPersons" class="py-16 text-center text-sm text-slate-400">Cargando personas…</div>
            <div v-else-if="companyUsers.length === 0" class="py-16 text-center">
              <p class="text-sm font-medium text-slate-600">Sin personas asignadas</p>
              <p class="text-xs text-slate-400 mt-1">Agrega usuarios existentes a esta empresa para asignarles roles.</p>
            </div>
            <div class="overflow-x-auto"><table class="w-full text-sm border-collapse" style="min-width:480px">
              <thead>
                <tr class="bg-slate-800">
                  <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-16">#</th>
                  <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300">Nombre completo</th>
                  <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-32">Usuario</th>
                  <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300">Roles asignados</th>
                  <th v-if="authStore.canAny('personas.eliminar', 'roles.asignar_modulos')" class="px-4 py-2.5 text-right text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-48">Acciones</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="cu in companyUsers" :key="cu.id" class="border-b border-slate-100 last:border-0 hover:bg-slate-50 transition-colors">
                  <td class="px-4 py-3 text-slate-400 tabular-nums">{{ cu.id }}</td>
                  <td class="px-4 py-3 font-medium text-slate-900">{{ personFullName(cu.person) }}</td>
                  <td class="px-4 py-3 font-mono text-xs text-slate-500">{{ cu.username }}</td>
                  <td class="px-4 py-3">
                    <div v-if="cu.roles.length === 0" class="text-xs text-slate-400">Sin roles</div>
                    <div v-else class="flex flex-wrap gap-1">
                      <span
                        v-for="(ur, idx) in cu.roles" :key="idx"
                        class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-medium bg-indigo-50 text-indigo-700"
                      >
                        {{ ur.roleName || 'Rol desconocido' }}
                        <span class="text-indigo-400">· {{ ur.branchName || 'Sucursal desconocida' }}</span>
                        <button
                          v-if="authStore.can('roles.asignar_modulos')"
                          @click.stop="confirmRemoveRole(cu, ur)"
                          class="ml-0.5 text-indigo-400 hover:text-red-500 leading-none transition-colors"
                          title="Desasignar este rol"
                        >&times;</button>
                      </span>
                    </div>
                  </td>
                  <td v-if="authStore.canAny('personas.eliminar', 'roles.asignar_modulos')" class="px-4 py-3">
                    <div class="flex items-center justify-end gap-2">
                      <button v-if="authStore.can('roles.asignar_modulos')" @click="openAssignRole(cu)" class="px-3 py-1 text-xs font-medium rounded border border-slate-200 text-slate-600 hover:bg-slate-100 hover:border-slate-300 transition-colors whitespace-nowrap">Asignar Rol</button>
                      <button v-if="authStore.can('personas.eliminar')" @click="confirmRemovePerson(cu)" class="px-3 py-1 text-xs font-medium rounded border border-red-100 text-red-500 hover:bg-red-50 hover:border-red-200 transition-colors whitespace-nowrap">Quitar</button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
            </div>
          </div>
        </div>

        </div><!-- /px-8 py-6 tab content -->
      </div><!-- /v-else selected -->
    </main>

    <!-- ══ MODALS ═══════════════════════════════════════════════════════════ -->

    <!-- Create company -->
    <Teleport to="body">
      <div v-if="createOpen" class="an-overlay" @click.self="createOpen = false">
        <div class="an-modal">
          <div class="an-modal-header">
            <h2 class="text-base font-semibold text-slate-900">Nueva Empresa</h2>
            <button @click="createOpen = false" class="an-close-btn">&times;</button>
          </div>
          <form @submit.prevent="submitCreate" class="an-modal-body">
            <div v-if="createError" class="mb-4 rounded border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">{{ createError }}</div>
            <div class="mb-4">
              <label class="an-label">Nombre de la empresa <span class="text-red-500">*</span></label>
              <input v-model="createForm.name" type="text" required maxlength="255" :disabled="createSaving" class="an-input" placeholder="Ej. Acme Corporation" />
            </div>
            <div class="mb-4">
              <label class="an-label">Estado</label>
              <select v-model="createForm.status" :disabled="createSaving" class="an-input">
                <option value="active">Activa</option>
                <option value="inactive">Inactiva</option>
              </select>
            </div>

            <!-- Toggle: create admin -->
            <div class="border-t border-slate-200 pt-4 mt-4">
              <label class="flex items-center gap-2.5 cursor-pointer select-none">
                <input type="checkbox" v-model="createForm.withAdmin" :disabled="createSaving" class="rounded border-slate-300 text-indigo-600 focus:ring-indigo-500" />
                <span class="text-sm font-medium text-slate-700">Crear administrador para esta empresa</span>
              </label>
              <p class="text-[11px] text-slate-400 mt-1 ml-6">Se creará un usuario con el rol "Administrador de Empresa" que podrá gestionar su propia empresa.</p>
            </div>

            <!-- Admin fields (conditional) -->
            <div v-if="createForm.withAdmin" class="mt-4 space-y-3 rounded-lg border border-slate-200 bg-slate-50 p-4">
              <p class="text-xs font-semibold uppercase tracking-wide text-slate-500 mb-2">Datos del administrador</p>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label class="an-label">Nombre <span class="text-red-500">*</span></label>
                  <input v-model="createForm.admin.firstName" type="text" :required="createForm.withAdmin" maxlength="100" :disabled="createSaving" class="an-input" placeholder="Juan" />
                </div>
                <div>
                  <label class="an-label">Apellido <span class="text-red-500">*</span></label>
                  <input v-model="createForm.admin.firstSurname" type="text" :required="createForm.withAdmin" maxlength="100" :disabled="createSaving" class="an-input" placeholder="Pérez" />
                </div>
              </div>
              <div>
                <label class="an-label">Usuario <span class="text-red-500">*</span></label>
                <input v-model="createForm.admin.username" type="text" :required="createForm.withAdmin" maxlength="100" :disabled="createSaving" class="an-input" placeholder="jperez" />
              </div>
              <div>
                <label class="an-label">Email <span class="text-red-500">*</span></label>
                <input v-model="createForm.admin.email" type="email" :required="createForm.withAdmin" :disabled="createSaving" class="an-input" placeholder="admin@empresa.com" />
              </div>
              <div>
                <label class="an-label">Contraseña <span class="text-red-500">*</span></label>
                <input v-model="createForm.admin.password" type="password" :required="createForm.withAdmin" minlength="8" :disabled="createSaving" class="an-input" placeholder="Mínimo 8 caracteres" />
              </div>
            </div>

            <div class="an-modal-footer">
              <button type="button" @click="createOpen = false" class="an-btn-ghost">Cancelar</button>
              <button type="submit" :disabled="createSaving" class="an-btn-primary">{{ createSaving ? 'Creando…' : 'Crear Empresa' }}</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

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

    <!-- Delete company confirm -->
    <Teleport to="body">
      <div v-if="deleteCompanyTarget" class="an-overlay" @click.self="deleteCompanyTarget = null">
        <div class="an-modal an-modal--sm">
          <div class="an-modal-header">
            <h2 class="text-base font-semibold text-slate-900">Confirmar eliminación</h2>
            <button @click="deleteCompanyTarget = null" class="an-close-btn">&times;</button>
          </div>
          <div class="an-modal-body">
            <p class="text-sm text-slate-600">
              ¿Estás seguro de que deseas eliminar
              <span class="font-semibold text-slate-900">{{ deleteCompanyTarget.name }}</span>?
              Esta acción eliminará también todas sus sucursales y no se puede deshacer.
            </p>
            <div class="an-modal-footer">
              <button @click="deleteCompanyTarget = null" class="an-btn-ghost">Cancelar</button>
              <button @click="doDeleteCompany" :disabled="deletingCompany" class="an-btn-danger">{{ deletingCompany ? 'Eliminando…' : 'Eliminar' }}</button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Create / Edit branch -->
    <Teleport to="body">
      <div v-if="branchModalOpen" class="an-overlay" @click.self="closeBranchModal">
        <div class="an-modal">
          <div class="an-modal-header">
            <h2 class="text-base font-semibold text-slate-900">{{ editingBranch ? 'Editar sucursal' : 'Nueva sucursal' }}</h2>
            <button @click="closeBranchModal" class="an-close-btn">&times;</button>
          </div>
          <form @submit.prevent="submitBranch" class="an-modal-body">
            <div v-if="branchFormError" class="mb-4 rounded border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">{{ branchFormError }}</div>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
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
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
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
            <h2 class="text-base font-semibold text-slate-900">Eliminar sucursal</h2>
            <button @click="deleteBranchTarget = null" class="an-close-btn">&times;</button>
          </div>
          <div class="an-modal-body">
            <p class="text-sm text-slate-600">
              ¿Confirmas la eliminación de
              <span class="font-semibold text-slate-900">{{ deleteBranchTarget.name }}</span>?
            </p>
            <div class="an-modal-footer">
              <button @click="deleteBranchTarget = null" class="an-btn-ghost">Cancelar</button>
              <button @click="doDeleteBranch" :disabled="deletingBranch" class="an-btn-danger">{{ deletingBranch ? 'Eliminando…' : 'Eliminar' }}</button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Create role -->
    <Teleport to="body">
      <div v-if="roleModalOpen" class="an-overlay" @click.self="closeRoleModal">
        <div class="an-modal">
          <div class="an-modal-header">
            <h2 class="text-base font-semibold text-slate-900">Nuevo Rol</h2>
            <button @click="closeRoleModal" class="an-close-btn">&times;</button>
          </div>
          <form @submit.prevent="submitRole" class="an-modal-body">
            <div v-if="roleFormError" class="mb-4 rounded border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">{{ roleFormError }}</div>
            <div class="mb-4">
              <label class="an-label">Nombre del rol <span class="text-red-500">*</span></label>
              <input v-model="roleForm.name" type="text" required maxlength="100" :disabled="roleSaving" class="an-input" placeholder="Ej. Administrador, Vendedor" />
            </div>
            <div>
              <label class="an-label">Aplicación <span class="text-red-500">*</span></label>
              <select v-model="roleForm.applicationId" required :disabled="roleSaving || loadingApps" class="an-input">
                <option value="" disabled>{{ loadingApps ? 'Cargando…' : 'Seleccionar aplicación' }}</option>
                <option v-for="app in applications" :key="app.id" :value="app.id">{{ app.name }}</option>
              </select>
            </div>
            <div class="an-modal-footer">
              <button type="button" @click="closeRoleModal" class="an-btn-ghost">Cancelar</button>
              <button type="submit" :disabled="roleSaving" class="an-btn-primary">{{ roleSaving ? 'Guardando…' : 'Guardar' }}</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Add person + assign role wizard -->
    <Teleport to="body">
      <div v-if="wizardOpen" class="an-overlay" @click.self="closeWizard">
        <div class="an-modal">
          <div class="an-modal-header">
            <h2 class="text-base font-semibold text-slate-900">{{ wizardTitle }}</h2>
            <button @click="closeWizard" class="an-close-btn">&times;</button>
          </div>

          <!-- Step indicator -->
          <div class="flex items-center gap-2 px-6 pt-4">
            <div class="flex items-center gap-2">
              <span
                class="inline-flex h-6 w-6 items-center justify-center rounded-full text-xs font-semibold"
                :class="wizardStep === 1 ? 'bg-indigo-600 text-white' : 'bg-indigo-100 text-indigo-600'"
              >1</span>
              <span class="text-xs font-medium" :class="wizardStep === 1 ? 'text-slate-900' : 'text-slate-400'">Persona</span>
            </div>
            <div class="h-px w-6 bg-slate-200"></div>
            <div class="flex items-center gap-2">
              <span
                class="inline-flex h-6 w-6 items-center justify-center rounded-full text-xs font-semibold"
                :class="wizardStep === 2 ? 'bg-indigo-600 text-white' : 'bg-slate-200 text-slate-500'"
              >2</span>
              <span class="text-xs font-medium" :class="wizardStep === 2 ? 'text-slate-900' : 'text-slate-400'">Rol</span>
            </div>
          </div>

          <!-- Step 1: select user -->
          <form v-if="wizardStep === 1" @submit.prevent="wizardNext" class="an-modal-body">
            <div v-if="wizardError" class="mb-4 rounded border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">{{ wizardError }}</div>
            <div>
              <label class="an-label">Usuario <span class="text-red-500">*</span></label>
              <select v-model="wizardForm.userId" required :disabled="wizardSaving || loadingAllUsers" class="an-input">
                <option value="" disabled>{{ loadingAllUsers ? 'Cargando...' : 'Seleccionar usuario' }}</option>
                <option v-for="u in availableUsers" :key="u.id" :value="u.id">
                  {{ personFullName(u.person) }} ({{ u.username }})
                </option>
              </select>
            </div>
            <div class="an-modal-footer">
              <button type="button" @click="closeWizard" class="an-btn-ghost">Cancelar</button>
              <button type="submit" :disabled="wizardSaving" class="an-btn-primary">{{ wizardSaving ? 'Agregando...' : 'Siguiente' }}</button>
            </div>
          </form>

          <!-- Step 2: assign role -->
          <form v-else @submit.prevent="wizardFinish" class="an-modal-body">
            <div v-if="wizardError" class="mb-4 rounded border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 whitespace-pre-wrap">{{ wizardError }}</div>

            <div class="mb-4 rounded-lg border border-slate-100 bg-slate-50 px-4 py-3">
              <p class="text-[11px] uppercase tracking-wide font-semibold text-slate-400 mb-1">Persona</p>
              <p class="text-sm font-medium text-slate-900">{{ wizardPersonLabel }}</p>
            </div>

            <div class="mb-4">
              <label class="an-label">Sucursales <span class="text-red-500">*</span></label>
              <p class="text-[11px] text-slate-400 mb-2">Selecciona una o más sucursales donde aplicará este rol.</p>
              <div class="max-h-36 overflow-y-auto rounded-lg border border-slate-200 bg-white divide-y divide-slate-100">
                <label
                  v-for="b in branches" :key="b.id"
                  class="flex items-center gap-2.5 px-3 py-2 cursor-pointer hover:bg-slate-50 transition-colors"
                  :class="{ 'opacity-50 pointer-events-none': wizardSaving }"
                >
                  <input
                    type="checkbox"
                    :value="b.id"
                    v-model="wizardForm.branchIds"
                    :disabled="wizardSaving"
                    class="rounded border-slate-300 text-indigo-600 focus:ring-indigo-500"
                  />
                  <span class="text-sm text-slate-700">{{ b.name }}</span>
                  <span class="text-[11px] text-slate-400 font-mono ml-auto">{{ b.code }}</span>
                </label>
              </div>
              <p v-if="wizardBranchError" class="text-xs text-red-500 mt-1">{{ te('form.select_branch') }}</p>
            </div>

            <!-- Multi-select roles grouped by app -->
            <div>
              <label class="an-label">Roles <span class="text-red-500">*</span></label>
              <p class="text-[11px] text-slate-400 mb-2">Selecciona uno o más roles a asignar.</p>
              <div v-if="!rolesLoaded" class="py-4 text-center text-xs text-slate-400">Cargando roles…</div>
              <div v-else class="max-h-52 overflow-y-auto rounded-lg border border-slate-200 bg-white divide-y divide-slate-100 space-y-0">
                <template v-for="group in rolesGroupedByApp" :key="group.appName">
                  <!-- App group header -->
                  <div class="sticky top-0 bg-slate-50 px-3 py-1.5 border-b border-slate-200 z-10 flex items-center gap-2">
                    <span class="text-[10px] font-bold uppercase tracking-widest text-slate-500">{{ group.appName || 'Sin aplicación' }}</span>
                  </div>
                  <label
                    v-for="r in group.roles" :key="r.id"
                    class="flex items-center gap-2.5 px-3 py-2.5 cursor-pointer hover:bg-indigo-50 transition-colors"
                    :class="{ 'opacity-50 pointer-events-none': wizardSaving }"
                  >
                    <input
                      type="checkbox"
                      :value="r.id"
                      v-model="wizardForm.roleIds"
                      :disabled="wizardSaving"
                      class="rounded border-slate-300 text-indigo-600 focus:ring-indigo-500"
                    />
                    <span class="text-sm text-slate-800 font-medium">{{ r.name }}</span>
                  </label>
                </template>
              </div>
              <p v-if="wizardRoleError" class="text-xs text-red-500 mt-1">{{ te('form.select_role') }}</p>
            </div>

            <div class="an-modal-footer">
              <button v-if="wizardIsNew" type="button" @click="wizardStep = 1" class="an-btn-ghost">Atrás</button>
              <button v-else type="button" @click="closeWizard" class="an-btn-ghost">Cancelar</button>
              <button type="submit" :disabled="wizardSaving || wizardForm.branchIds.length === 0" class="an-btn-primary">{{ wizardSaving ? 'Asignando...' : 'Asignar' }}</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Remove person from company -->
    <Teleport to="body">
      <div v-if="removePersonTarget" class="an-overlay" @click.self="removePersonTarget = null">
        <div class="an-modal an-modal--sm">
          <div class="an-modal-header">
            <h2 class="text-base font-semibold text-slate-900">Quitar persona</h2>
            <button @click="removePersonTarget = null" class="an-close-btn">&times;</button>
          </div>
          <div class="an-modal-body">
            <p class="text-sm text-slate-600">
              ¿Quitar a
              <span class="font-semibold text-slate-900">{{ personFullName(removePersonTarget.person) }}</span>
              de {{ selected?.name }}? Se eliminarán también sus roles en esta empresa.
            </p>
            <div class="an-modal-footer">
              <button @click="removePersonTarget = null" class="an-btn-ghost">Cancelar</button>
              <button @click="doRemovePerson" :disabled="removingPerson" class="an-btn-danger">{{ removingPerson ? 'Quitando…' : 'Quitar' }}</button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Confirm unassign individual role -->
    <Teleport to="body">
      <div v-if="removeRoleTarget" class="an-overlay" @click.self="removeRoleTarget = null">
        <div class="an-modal an-modal--sm">
          <div class="an-modal-header">
            <h2 class="text-base font-semibold text-slate-900">Desasignar rol</h2>
            <button @click="removeRoleTarget = null" class="an-close-btn">&times;</button>
          </div>
          <div class="an-modal-body">
            <p class="text-sm text-slate-600">
              ¿Quitar el rol
              <span class="font-semibold text-slate-900">{{ removeRoleTarget?.role?.roleName }}</span>
              en <span class="font-semibold text-slate-900">{{ removeRoleTarget?.role?.branchName }}</span>
              a <span class="font-semibold text-slate-900">{{ personFullName(removeRoleTarget?.person?.person) }}</span>?
            </p>
            <div class="an-modal-footer">
              <button @click="removeRoleTarget = null" class="an-btn-ghost">Cancelar</button>
              <button @click="doRemoveUserRole" :disabled="removingRole" class="an-btn-danger">{{ removingRole ? 'Quitando…' : 'Desasignar' }}</button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Role module assignment (checkbox list) -->
    <Teleport to="body">
      <div v-if="roleModulesOpen" class="an-overlay" @click.self="closeRoleModules">
        <div class="an-modal" style="max-width:520px">
          <!-- Header with dark strip -->
          <div class="an-modal-header" style="background:linear-gradient(135deg,#0f172a,#1e293b); border-bottom:none; border-radius:0.5rem 0.5rem 0 0">
            <div>
              <h2 class="text-base font-semibold text-white">{{ roleModulesTarget?.name }}</h2>
              <p class="text-xs text-slate-400 mt-0.5">Módulos del rol · {{ roleModulesTarget?.appName || roleModulesAppName }}</p>
            </div>
            <button @click="closeRoleModules" class="text-slate-400 hover:text-white text-xl leading-none">&times;</button>
          </div>
          <div class="an-modal-body">
            <div v-if="roleModulesError" class="mb-4 rounded border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">{{ roleModulesError }}</div>

            <div v-if="loadingRoleModules" class="py-8 text-center text-sm text-slate-400">Cargando módulos...</div>

            <div v-else-if="roleModulesAppModules.length === 0" class="py-8 text-center">
              <p class="text-sm text-slate-500">No hay módulos para esta aplicación.</p>
              <p class="text-xs text-slate-400 mt-1">Crea módulos en la sección Aplicaciones primero.</p>
            </div>

            <!-- Grouped by prefix -->
            <div v-else class="max-h-[60vh] overflow-y-auto space-y-4">
              <div v-for="group in roleModulesGrouped" :key="group.prefix">
                <!-- Group header -->
                <div class="flex items-center gap-2 mb-1">
                  <span class="text-[10px] font-bold uppercase tracking-widest text-slate-400">{{ group.label }}</span>
                  <div class="flex-1 h-px bg-slate-100"></div>
                  <span class="text-[10px] text-slate-300">{{ group.modules.filter(m => roleModulesAssigned.has(m.id)).length }}/{{ group.modules.length }}</span>
                </div>
                <label
                  v-for="m in group.modules" :key="m.id"
                  class="flex items-center gap-3 px-3 py-2.5 rounded-lg hover:bg-slate-50 cursor-pointer transition-colors border border-transparent hover:border-slate-100"
                  :class="{ 'opacity-60': roleModulesSaving }"
                >
                  <input
                    type="checkbox"
                    :checked="roleModulesAssigned.has(m.id)"
                    @change="toggleRoleModule(m.id, $event.target.checked)"
                    :disabled="roleModulesSaving"
                    class="h-4 w-4 rounded border-slate-300 text-indigo-600 focus:ring-indigo-500 shrink-0"
                  />
                  <div class="flex-1 min-w-0">
                    <p class="text-sm font-medium text-slate-900 leading-snug">{{ m.displayName || m.name }}</p>
                    <p class="text-[11px] text-slate-400 font-mono mt-0.5">{{ m.name }}</p>
                  </div>
                  <span
                    v-if="roleModulesAssigned.has(m.id)"
                    class="shrink-0 w-1.5 h-1.5 rounded-full bg-indigo-500"
                  ></span>
                </label>
              </div>
            </div>

            <div class="an-modal-footer mt-4">
              <button @click="closeRoleModules" class="an-btn-primary">Cerrar</button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { companyService } from '@/services/companyService'
import { branchService } from '@/services/branchService'
import { useAuthStore } from '@/store/auth'
import api from '@/services/api'
import { te, teError } from '@/i18n'

const authStore = useAuthStore()

// Returns a deterministic colour class for a company avatar based on its name initial.
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
function companyAvatarClass(name = '') {
  const idx = (name.charCodeAt(0) || 0) % AVATAR_COLOURS.length
  return AVATAR_COLOURS[idx]
}

// ── Tabs ──────────────────────────────────────────────────────────────────────
const allTabs = [
  { key: 'branches', label: 'Sucursales',  prefix: 'sucursales' },
  { key: 'personas', label: 'Personas',    prefix: 'personas' },
  { key: 'roles',    label: 'Roles',       prefix: 'roles' },
]
const visibleTabs = computed(() =>
  allTabs.filter(t => authStore.canPrefix(t.prefix))
)
const activeTab = ref('branches')

// ── Companies list ────────────────────────────────────────────────────────────
const companies        = ref([])
const loadingCompanies = ref(true)
const listError        = ref('')
const search           = ref('')
const selected         = ref(null)

const filteredCompanies = computed(() =>
  companies.value.filter(c =>
    c.name.toLowerCase().includes(search.value.toLowerCase())
  )
)

async function loadCompanies() {
  loadingCompanies.value = true
  listError.value = ''
  try {
    const res = await companyService.getMyCompanies()
    companies.value = res.data ?? []
  } catch (e) {
    listError.value = teError(e)
  } finally {
    loadingCompanies.value = false
  }
}

function selectCompany(c) {
  if (selected.value?.id === c.id) return
  selected.value = c
  activeTab.value = 'branches'
  detailError.value = ''
  loadBranches(c.id)
}

onMounted(loadCompanies)

// ── Create company ────────────────────────────────────────────────────────────
const createOpen   = ref(false)
const createSaving = ref(false)
const createError  = ref('')
const createForm   = reactive({
  name: '', status: 'active', withAdmin: false,
  admin: { firstName: '', firstSurname: '', username: '', email: '', password: '' },
})

function openCreate() {
  createForm.name   = ''
  createForm.status = 'active'
  createForm.withAdmin = false
  Object.assign(createForm.admin, { firstName: '', firstSurname: '', username: '', email: '', password: '' })
  createError.value = ''
  createOpen.value  = true
}

async function submitCreate() {
  createError.value  = ''
  createSaving.value = true
  try {
    const payload = { name: createForm.name.trim(), status: createForm.status }
    if (createForm.withAdmin) {
      payload.admin = {
        firstName:    createForm.admin.firstName.trim(),
        firstSurname: createForm.admin.firstSurname.trim(),
        username:     createForm.admin.username.trim(),
        email:        createForm.admin.email.trim(),
        password:     createForm.admin.password,
      }
    }
    const res = await companyService.create(payload)
    const newCompany = res.data
    companies.value.unshift(newCompany)
    createOpen.value = false
    selectCompany(newCompany)
  } catch (e) {
    createError.value = teError(e)
  } finally {
    createSaving.value = false
  }
}

// ── Edit company ──────────────────────────────────────────────────────────────
const editCompanyOpen   = ref(false)
const editCompanySaving = ref(false)
const editCompanyError  = ref('')
const editCompanyForm   = reactive({ name: '', status: 'active' })

function openEditCompany() {
  editCompanyForm.name   = selected.value.name
  editCompanyForm.status = selected.value.status
  editCompanyError.value = ''
  editCompanyOpen.value  = true
}

async function submitEditCompany() {
  editCompanyError.value  = ''
  editCompanySaving.value = true
  try {
    const res = await companyService.update(selected.value.id, {
      name:   editCompanyForm.name.trim(),
      status: editCompanyForm.status,
    })
    const updated = res.data
    const idx = companies.value.findIndex(c => c.id === updated.id)
    if (idx !== -1) companies.value[idx] = updated
    selected.value = updated
    editCompanyOpen.value = false
  } catch (e) {
    editCompanyError.value = teError(e)
  } finally {
    editCompanySaving.value = false
  }
}

// ── Delete company ────────────────────────────────────────────────────────────
const deleteCompanyTarget = ref(null)
const deletingCompany     = ref(false)

function confirmDeleteCompany(c) { deleteCompanyTarget.value = c }

async function doDeleteCompany() {
  deletingCompany.value = true
  try {
    await companyService.remove(deleteCompanyTarget.value.id)
    companies.value = companies.value.filter(c => c.id !== deleteCompanyTarget.value.id)
    if (selected.value?.id === deleteCompanyTarget.value.id) selected.value = null
    deleteCompanyTarget.value = null
  } catch (e) {
    listError.value = teError(e)
  } finally {
    deletingCompany.value = false
  }
}

// ── Branches ──────────────────────────────────────────────────────────────────
const branches        = ref([])
const loadingBranches = ref(false)
const detailError     = ref('')

async function loadBranches(companyId) {
  loadingBranches.value = true
  detailError.value = ''
  try {
    const res = await branchService.getAll(companyId)
    branches.value = res.data ?? []
  } catch (e) {
    detailError.value = teError(e)
  } finally {
    loadingBranches.value = false
  }
}

// Branch form
const branchModalOpen    = ref(false)
const branchSaving       = ref(false)
const branchFormError    = ref('')
const editingBranch      = ref(null)
const deleteBranchTarget = ref(null)
const deletingBranch     = ref(false)

const branchForm = reactive({ code: '', name: '', address: '', phoneNumber: '', email: '', status: 'active' })

function openCreateBranch() {
  editingBranch.value   = null
  branchFormError.value = ''
  Object.assign(branchForm, { code: '', name: '', address: '', phoneNumber: '', email: '', status: 'active' })
  branchModalOpen.value = true
}

function openEditBranch(b) {
  editingBranch.value   = b
  branchFormError.value = ''
  Object.assign(branchForm, {
    code: b.code, name: b.name,
    address: b.address ?? '', phoneNumber: b.phoneNumber ?? '',
    email: b.email ?? '', status: b.status,
  })
  branchModalOpen.value = true
}

function closeBranchModal() {
  if (!branchSaving.value) branchModalOpen.value = false
}

async function submitBranch() {
  branchFormError.value = ''
  branchSaving.value    = true
  const payload = {
    name:        branchForm.name.trim(),
    address:     branchForm.address.trim()      || undefined,
    phoneNumber: branchForm.phoneNumber.trim()  || undefined,
    email:       branchForm.email.trim()        || undefined,
  }
  try {
    if (editingBranch.value) {
      const res = await branchService.update(selected.value.id, editingBranch.value.id, { ...payload, status: branchForm.status })
      const idx = branches.value.findIndex(b => b.id === editingBranch.value.id)
      if (idx !== -1) branches.value[idx] = res.data
    } else {
      const res = await branchService.create(selected.value.id, { ...payload, code: branchForm.code.trim() })
      branches.value.unshift(res.data)
    }
    branchModalOpen.value = false
  } catch (e) {
    branchFormError.value = teError(e)
  } finally {
    branchSaving.value = false
  }
}

function confirmDeleteBranch(b) { deleteBranchTarget.value = b }

async function doDeleteBranch() {
  deletingBranch.value = true
  try {
    await branchService.remove(selected.value.id, deleteBranchTarget.value.id)
    branches.value = branches.value.filter(b => b.id !== deleteBranchTarget.value.id)
    deleteBranchTarget.value = null
  } catch (e) {
    detailError.value = teError(e)
    deleteBranchTarget.value = null
  } finally {
    deletingBranch.value = false
  }
}

// ── Roles (lazy, once per session) ───────────────────────────────────────────
const roles        = ref([])
const loadingRoles = ref(false)
const rolesLoaded  = ref(false)

watch(activeTab, (tab) => {
  if (tab === 'roles' && !rolesLoaded.value) loadRoles()
  if (tab === 'roles') loadApplications()
})

async function loadRoles() {
  loadingRoles.value = true
  try {
    const res = await api.get('/roles', { params: { pageSize: 200 } })
    roles.value = res.data?.data ?? []
    rolesLoaded.value = true
  } catch {
    roles.value = []
  } finally {
    loadingRoles.value = false
  }
}

// ── Create role ─────────────────────────────────────────────────────────────

// Helper: get app name for display
function appName(appId) {
  const a = applications.value.find(x => x.id === appId)
  return a ? a.name : `App #${appId}`
}
const roleModalOpen  = ref(false)
const roleSaving     = ref(false)
const roleFormError  = ref('')
const roleForm       = reactive({ name: '', applicationId: '' })

// Applications for the dropdown
const applications = ref([])
const loadingApps  = ref(false)
const appsLoaded   = ref(false)

async function loadApplications() {
  if (appsLoaded.value) return
  loadingApps.value = true
  try {
    const res = await api.get('/applications', { params: { pageSize: 200 } })
    applications.value = res.data?.data ?? []
    appsLoaded.value = true
  } catch {
    applications.value = []
  } finally {
    loadingApps.value = false
  }
}

function openCreateRole() {
  roleForm.name          = ''
  roleForm.applicationId = ''
  roleFormError.value    = ''
  roleModalOpen.value    = true
  loadApplications()
}

function closeRoleModal() {
  if (!roleSaving.value) roleModalOpen.value = false
}

async function submitRole() {
  roleFormError.value = ''
  roleSaving.value    = true
  try {
    const res = await api.post('/roles', {
      name:          roleForm.name.trim(),
      applicationId: Number(roleForm.applicationId),
    })
    roles.value.unshift(res.data?.data ?? res.data)
    roleModalOpen.value = false
  } catch (e) {
    roleFormError.value = teError(e)
  } finally {
    roleSaving.value = false
  }
}

// ── Personas (company users) ────────────────────────────────────────────────
const companyUsers     = ref([])
const loadingPersons   = ref(false)
const allUsers         = ref([])
const loadingAllUsers  = ref(false)

// Computed list: only users not already in the company
const availableUsers = computed(() => {
  const existingIds = new Set(companyUsers.value.map(cu => cu.id))
  return allUsers.value.filter(u => !existingIds.has(u.id))
})

// Load when tab activates or when company changes while tab is active
watch(activeTab, (tab) => {
  if (tab === 'personas' && selected.value) loadCompanyUsers()
})

function personFullName(p) {
  if (!p) return '—'
  return [p.firstName, p.firstSurname, p.secondSurname].filter(Boolean).join(' ') || '—'
}

async function loadCompanyUsers() {
  loadingPersons.value = true
  try {
    const res = await api.get(`/companies/${selected.value.id}/users`)
    companyUsers.value = res.data?.data ?? res.data ?? []
  } catch {
    companyUsers.value = []
  } finally {
    loadingPersons.value = false
  }
}

async function loadAllUsers() {
  if (allUsers.value.length > 0) return
  loadingAllUsers.value = true
  try {
    const res = await api.get('/users', { params: { pageSize: 500 } })
    allUsers.value = res.data?.data ?? []
  } catch {
    allUsers.value = []
  } finally {
    loadingAllUsers.value = false
  }
}

// ── Add person + assign role wizard ──────────────────────────────────────────
const wizardOpen    = ref(false)
const wizardStep    = ref(1)
const wizardSaving  = ref(false)
const wizardError   = ref('')
const wizardIsNew   = ref(true)        // true = adding new person, false = assigning role to existing
const wizardTargetUser = ref(null)     // set when opening from table "Asignar Rol"
const wizardBranchError = ref(false)
const wizardRoleError   = ref(false)
const wizardForm    = reactive({ userId: '', branchIds: [], roleIds: [] })

const wizardTitle = computed(() => {
  if (wizardIsNew.value) return `Agregar Persona a ${selected.value?.name ?? ''}`
  return 'Asignar Rol'
})

const wizardPersonLabel = computed(() => {
  if (wizardTargetUser.value) return personFullName(wizardTargetUser.value.person)
  const u = allUsers.value.find(u => String(u.id) === String(wizardForm.userId))
  return u ? `${personFullName(u.person)} (${u.username})` : ''
})

// Called from "Agregar Persona" button — starts at step 1
function openAddPerson() {
  wizardIsNew.value = true
  wizardTargetUser.value = null
  wizardStep.value  = 1
  wizardError.value = ''
  Object.assign(wizardForm, { userId: '', branchIds: [], roleIds: [] })
  wizardBranchError.value = false
  wizardRoleError.value   = false
  wizardOpen.value  = true
  loadAllUsers()
}

// Called from table "Asignar Rol" button — opens directly at step 2
function openAssignRole(cu) {
  wizardIsNew.value = false
  wizardTargetUser.value = cu
  wizardStep.value  = 2
  wizardError.value = ''
  Object.assign(wizardForm, { userId: String(cu.id), branchIds: [], roleIds: [] })
  wizardBranchError.value = false
  wizardRoleError.value   = false
  wizardOpen.value  = true
  if (!rolesLoaded.value) loadRoles()
  if (!branches.value.length && selected.value) loadBranches(selected.value.id)
}

function closeWizard() {
  if (!wizardSaving.value) wizardOpen.value = false
}

// Step 1 → Step 2: add user to company, then advance
async function wizardNext() {
  wizardError.value  = ''
  wizardSaving.value = true
  try {
    await api.post(`/companies/${selected.value.id}/users`, {
      userId: Number(wizardForm.userId),
    })
    // Keep the selected user reference for step 2 label
    const u = allUsers.value.find(u => String(u.id) === String(wizardForm.userId))
    if (u) wizardTargetUser.value = u
    wizardStep.value = 2
    // Pre-load roles & branches for step 2
    if (!rolesLoaded.value) loadRoles()
    if (!branches.value.length && selected.value) loadBranches(selected.value.id)
  } catch (e) {
    wizardError.value = teError(e)
  } finally {
    wizardSaving.value = false
  }
}

// Step 2: assign roles (multiple) to all selected branches and close
async function wizardFinish() {
  wizardBranchError.value = wizardForm.branchIds.length === 0
  wizardRoleError.value   = wizardForm.roleIds.length === 0
  if (wizardBranchError.value || wizardRoleError.value) return

  wizardError.value  = ''
  wizardSaving.value = true
  const userId = wizardTargetUser.value?.id ?? Number(wizardForm.userId)
  const errors = []

  // Double loop: every branchId × every roleId
  for (const bid of wizardForm.branchIds) {
    for (const roleId of wizardForm.roleIds) {
      try {
        await api.post(
          `/companies/${selected.value.id}/users/${userId}/roles`,
          { branchId: Number(bid), roleId: Number(roleId) },
        )
      } catch (e) {
        const msg = e?.response?.data?.message ?? 'Error desconocido'
        const branchName = branches.value.find(b => b.id === bid)?.name ?? `#${bid}`
        const roleName   = roles.value.find(r => r.id === roleId)?.name ?? `#${roleId}`
        errors.push(`${branchName} / ${roleName}: ${msg}`)
      }
    }
  }

  const totalOps = wizardForm.branchIds.length * wizardForm.roleIds.length
  if (errors.length === totalOps) {
    wizardError.value = errors.join('\n')
    wizardSaving.value = false
    return
  }
  if (errors.length > 0) {
    wizardError.value = `Algunas asignaciones fallaron:\n${errors.join('\n')}`
  }
  wizardOpen.value = false
  wizardSaving.value = false
  await loadCompanyUsers()
}

// ── Remove person ─────────────────────────────────────────────────────────────
const removePersonTarget = ref(null)
const removingPerson     = ref(false)

function confirmRemovePerson(cu) {
  removePersonTarget.value = cu
}

// ── Remove individual role from a person ─────────────────────────────────────
const removeRoleTarget = ref(null) // { person: CompanyUser, role: CompanyUserRoleResponse }
const removingRole     = ref(false)

function confirmRemoveRole(cu, role) {
  removeRoleTarget.value = { person: cu, role }
}

async function doRemoveUserRole() {
  removingRole.value = true
  const { person, role } = removeRoleTarget.value
  try {
    await api.delete(
      `/companies/${selected.value.id}/users/${person.id}/roles`,
      { data: { branchId: role.branchId, roleId: role.roleId } },
    )
    removeRoleTarget.value = null
    await loadCompanyUsers()
  } catch {
    removeRoleTarget.value = null
  } finally {
    removingRole.value = false
  }
}

async function doRemovePerson() {
  removingPerson.value = true
  try {
    await api.delete(`/companies/${selected.value.id}/users/${removePersonTarget.value.id}`)
    removePersonTarget.value = null
    await loadCompanyUsers()
  } catch {
    removePersonTarget.value = null
  } finally {
    removingPerson.value = false
  }
}

// ── Role module assignment (checkbox modal) ─────────────────────────────────
const roleModulesOpen       = ref(false)
const roleModulesTarget     = ref(null)
const roleModulesAssigned   = ref(new Set())
const roleModulesAppModules = ref([])
const loadingRoleModules    = ref(false)
const roleModulesSaving     = ref(false)
const roleModulesError      = ref('')

async function openRoleModules(role) {
  roleModulesTarget.value = role
  roleModulesError.value  = ''
  roleModulesOpen.value   = true
  loadingRoleModules.value = true

  try {
    // Load all modules for the role's application
    const [modsRes, assignedRes] = await Promise.all([
      api.get(`/applications/${role.applicationId}/modules`, { params: { pageSize: 500 } }),
      api.get(`/roles/${role.id}/modules`),
    ])
    roleModulesAppModules.value = modsRes.data?.data ?? []
    const assignedList = assignedRes.data?.data ?? assignedRes.data ?? []
    roleModulesAssigned.value = new Set(assignedList.map(m => m.id))
  } catch {
    roleModulesAppModules.value = []
    roleModulesAssigned.value = new Set()
  } finally {
    loadingRoleModules.value = false
  }
}

function closeRoleModules() {
  roleModulesOpen.value = false
}

async function toggleRoleModule(moduleId, checked) {
  roleModulesSaving.value = true
  roleModulesError.value  = ''
  try {
    if (checked) {
      await api.post(`/roles/${roleModulesTarget.value.id}/modules`, { moduleIds: [moduleId] })
      roleModulesAssigned.value.add(moduleId)
    } else {
      await api.delete(`/roles/${roleModulesTarget.value.id}/modules`, { data: { moduleIds: [moduleId] } })
      roleModulesAssigned.value.delete(moduleId)
    }
    // Force reactivity
    roleModulesAssigned.value = new Set(roleModulesAssigned.value)
  } catch (e) {
    roleModulesError.value = teError(e)
    // Revert checkbox state
    if (checked) roleModulesAssigned.value.delete(moduleId)
    else roleModulesAssigned.value.add(moduleId)
    roleModulesAssigned.value = new Set(roleModulesAssigned.value)
  } finally {
    roleModulesSaving.value = false
  }
}

// Computed: group roles by appName for the wizard multi-select
const rolesGroupedByApp = computed(() => {
  const map = new Map()
  for (const r of roles.value) {
    const key = r.appName || 'Sin aplicación'
    if (!map.has(key)) map.set(key, [])
    map.get(key).push(r)
  }
  return [...map.entries()].map(([appName, appRoles]) => ({ appName, roles: appRoles }))
})

// Computed: group roleModulesAppModules by the prefix before the first dot
// e.g. "empresas.crear" → group "empresas" with label "Empresas"
const roleModulesAppName = computed(() => roleModulesTarget.value?.appName || '')

const roleModulesGrouped = computed(() => {
  const map = new Map()
  for (const m of roleModulesAppModules.value) {
    const prefix = m.name.includes('.') ? m.name.split('.')[0] : 'general'
    if (!map.has(prefix)) map.set(prefix, [])
    map.get(prefix).push(m)
  }
  return [...map.entries()].map(([prefix, modules]) => ({
    prefix,
    label: prefix.charAt(0).toUpperCase() + prefix.slice(1).replace(/_/g, ' '),
    modules,
  }))
})
</script>
