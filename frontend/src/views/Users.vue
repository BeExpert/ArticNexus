<template>
  <div>
    <!-- Header -->
    <div class="an-page-header">
      <div>
        <h1 class="an-page-header-title">Usuarios</h1>
        <p class="an-page-header-subtitle">Administra las cuentas de usuario del sistema.</p>
      </div>
      <div class="flex items-center gap-2">
        <input
          v-model="search"
          type="text"
          placeholder="Buscar usuario…"
          class="an-input flex-1 sm:w-56 sm:flex-none text-sm py-1.5"
        />
        <button v-if="authStore.can('usuarios.crear')" @click="openCreate" class="an-btn-primary text-sm py-2 px-4 shrink-0">Nuevo Usuario</button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="py-20 text-center text-sm text-slate-400">Cargando usuarios…</div>

    <!-- Error -->
    <div v-else-if="error" class="rounded border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">{{ error }}</div>

    <!-- Empty -->
    <div v-else-if="users.length === 0" class="an-card py-16 text-center">
      <p class="text-sm font-medium text-slate-600">Sin usuarios registrados</p>
      <p class="text-xs text-slate-400 mt-1">Crea el primer usuario para comenzar.</p>
    </div>

    <!-- Table -->
    <div v-else class="an-card overflow-hidden">
      <div class="overflow-x-auto">
      <table class="w-full text-sm border-collapse" style="min-width:600px">
        <thead>
          <tr class="bg-slate-800">
            <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-12 hidden sm:table-cell">#</th>
            <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300">Nombre completo</th>
            <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-32 hidden md:table-cell">Identificación</th>
            <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-36">Usuario</th>
            <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 hidden lg:table-cell">Email</th>
            <th class="px-4 py-2.5 text-left text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-24">Estado</th>
            <th v-if="authStore.canAny('usuarios.editar', 'usuarios.eliminar', 'usuarios.reset_contrasena')" class="px-4 py-2.5 text-right text-[11px] font-semibold uppercase tracking-wide text-slate-300 w-40">Acciones</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="u in users" :key="u.id" class="border-b border-slate-100 last:border-0 hover:bg-slate-50 transition-colors">
            <td class="px-4 py-3 text-slate-400 tabular-nums hidden sm:table-cell">{{ u.id }}</td>
            <td class="px-4 py-3 font-medium text-slate-900">{{ fullName(u.person) }}</td>
            <td class="px-4 py-3 text-slate-500 hidden md:table-cell">{{ u.person?.nationalId ?? '—' }}</td>
            <td class="px-4 py-3 font-mono text-xs text-slate-500">{{ u.username }}</td>
            <td class="px-4 py-3 text-slate-500 hidden lg:table-cell">{{ u.email }}</td>
            <td class="px-4 py-3">
              <div class="flex flex-col gap-1">
                <span :class="u.status === 'active' ? 'an-badge-active' : 'an-badge-inactive'">
                  {{ u.status === 'active' ? 'Activo' : 'Inactivo' }}
                </span>
                <span v-if="u.passwordExpiresAt" :class="[
                  'inline-flex items-center gap-1 text-[10px] font-medium px-1.5 py-0.5 rounded',
                  isExpired(u.passwordExpiresAt) ? 'bg-red-100 text-red-600' : 'bg-amber-100 text-amber-700'
                ]">
                  {{ isExpired(u.passwordExpiresAt) ? '&#9888; Demo expirado' : '&#9201; Demo ' + formatExpiry(u.passwordExpiresAt) }}
                </span>
              </div>
            </td>
            <td v-if="authStore.canAny('usuarios.editar', 'usuarios.eliminar', 'usuarios.reset_contrasena')" class="px-4 py-3">
              <div class="flex items-center justify-end gap-2">
                <button v-if="authStore.can('usuarios.reset_contrasena')" @click="openResetPwd(u)" title="Resetear contraseña" class="p-1.5 rounded text-slate-500 hover:text-slate-800 hover:bg-slate-100 transition-colors"><KeyRound :size="15" /></button>
                <button v-if="authStore.can('usuarios.editar')" @click="openEdit(u)" title="Editar" class="p-1.5 rounded text-slate-500 hover:text-slate-800 hover:bg-slate-100 transition-colors"><Pencil :size="15" /></button>
                <button v-if="authStore.can('usuarios.eliminar')" @click="confirmDelete(u)" title="Eliminar" class="p-1.5 rounded text-red-400 hover:text-red-600 hover:bg-red-50 transition-colors"><Trash2 :size="15" /></button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      </div>
    </div>

    <!-- Pagination footer -->
    <div v-if="!loading && totalItems > 0" class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between mt-4 text-sm text-slate-500">
      <span class="text-center sm:text-left">Mostrando {{ rangeStart }}–{{ rangeEnd }} de {{ totalItems }}</span>
      <div class="flex items-center justify-center gap-2">
        <button
          @click="prevPage"
          :disabled="page <= 1"
          class="px-3 py-1 rounded border border-slate-200 text-slate-600 hover:bg-slate-100 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
        >Anterior</button>
        <span class="text-xs text-slate-400">Página {{ page }} de {{ totalPages }}</span>
        <button
          @click="nextPage"
          :disabled="page >= totalPages"
          class="px-3 py-1 rounded border border-slate-200 text-slate-600 hover:bg-slate-100 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
        >Siguiente</button>
      </div>
    </div>

    <!-- ══ Create modal (2-step wizard) ═════════════════════════════════════ -->
    <Teleport to="body">
      <div v-if="createOpen" class="an-overlay" @click.self="closeCreate">
        <div class="an-modal" style="max-width:520px">
          <div class="an-modal-header">
            <h2 class="text-base font-semibold text-slate-900">Nuevo Usuario</h2>
            <button @click="closeCreate" class="an-close-btn">&times;</button>
          </div>

          <!-- Step indicator -->
          <div class="flex items-center gap-3 px-5 pt-4 pb-1">
            <div class="flex items-center gap-2 text-xs font-medium"
                 :class="step === 1 ? 'text-indigo-600' : 'text-slate-400'">
              <span class="flex items-center justify-center w-5 h-5 rounded-full text-[10px] font-bold"
                    :class="step === 1 ? 'bg-indigo-600 text-white' : 'bg-slate-200 text-slate-500'">1</span>
              Persona
            </div>
            <div class="flex-1 h-px bg-slate-200"></div>
            <div class="flex items-center gap-2 text-xs font-medium"
                 :class="step === 2 ? 'text-indigo-600' : 'text-slate-400'">
              <span class="flex items-center justify-center w-5 h-5 rounded-full text-[10px] font-bold"
                    :class="step === 2 ? 'bg-indigo-600 text-white' : 'bg-slate-200 text-slate-500'">2</span>
              Credenciales
            </div>
          </div>

          <form @submit.prevent="step === 1 ? goStep2() : submitCreate()" class="an-modal-body">
            <div v-if="formError" class="mb-4 rounded border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">{{ formError }}</div>

            <!-- Step 1: Person data -->
            <template v-if="step === 1">
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
                <div>
                  <label class="an-label">Nombre <span class="text-red-500">*</span></label>
                  <input v-model="pForm.firstName" type="text" required maxlength="100" class="an-input" placeholder="Ej. María" />
                </div>
                <div>
                  <label class="an-label">Primer Apellido <span class="text-red-500">*</span></label>
                  <input v-model="pForm.firstSurname" type="text" required maxlength="100" class="an-input" placeholder="Ej. García" />
                </div>
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
                <div>
                  <label class="an-label">Segundo Apellido</label>
                  <input v-model="pForm.secondSurname" type="text" class="an-input" placeholder="Opcional" />
                </div>
                <div>
                  <label class="an-label">Identificación</label>
                  <input v-model="pForm.nationalId" type="text" class="an-input" placeholder="Ej. 1-0234-0567" />
                </div>
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
                <div>
                  <label class="an-label">Fecha de Nacimiento</label>
                  <input v-model="pForm.birthDate" type="date" class="an-input" />
                </div>
                <div>
                  <label class="an-label">Código de Área</label>
                  <input v-model="pForm.phoneAreaCode" type="text" class="an-input" placeholder="Ej. +506" />
                </div>
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
                <div>
                  <label class="an-label">Teléfono Principal</label>
                  <input v-model="pForm.primaryPhone" type="text" class="an-input" placeholder="Ej. 8888-8888" />
                </div>
                <div>
                  <label class="an-label">Teléfono Secundario</label>
                  <input v-model="pForm.secondaryPhone" type="text" class="an-input" placeholder="Opcional" />
                </div>
              </div>
              <div class="mb-4">
                <label class="an-label">Dirección</label>
                <input v-model="pForm.address" type="text" class="an-input" placeholder="Dirección física" />
              </div>
            </template>

            <!-- Step 2: User credentials -->
            <template v-if="step === 2">
              <!-- Person summary card -->
              <div class="mb-4 rounded-lg border border-slate-100 bg-slate-50 px-4 py-3">
                <p class="text-[11px] uppercase tracking-wide font-semibold text-slate-400 mb-1">Persona</p>
                <p class="text-sm font-medium text-slate-900">{{ pForm.firstName }} {{ pForm.firstSurname }} {{ pForm.secondSurname }}</p>
                <p v-if="pForm.nationalId" class="text-xs text-slate-500">{{ pForm.nationalId }}</p>
              </div>

              <div class="mb-4">
                <label class="an-label">Nombre de usuario <span class="text-red-500">*</span></label>
                <input v-model="uForm.username" type="text" required maxlength="100" class="an-input" placeholder="Ej. mgarcia" />
              </div>
              <div class="mb-4">
                <label class="an-label">Correo electrónico <span class="text-red-500">*</span></label>
                <input v-model="uForm.email" type="email" required class="an-input" placeholder="Ej. maria@empresa.com" />
              </div>

              <!-- Send credentials checkbox -->
              <div class="mb-4 rounded-lg border border-slate-100 bg-slate-50 px-4 py-3">
                <label class="flex items-start gap-3 cursor-pointer">
                  <input
                    v-model="uForm.sendCredentials"
                    type="checkbox"
                    class="mt-0.5 h-4 w-4 rounded border-slate-300 text-indigo-600 accent-indigo-600"
                  />
                  <span class="flex flex-col">
                    <span class="text-sm font-medium text-slate-800">Enviar credenciales al usuario</span>
                    <span class="text-xs text-slate-500 mt-0.5">Se generará una contraseña aleatoria y se enviará al correo indicado junto con el nombre de usuario y el link de acceso.</span>
                  </span>
                </label>
              </div>

              <div v-if="!uForm.sendCredentials" class="mb-4">
                <label class="an-label">Contraseña <span class="text-red-500">*</span></label>
                <input v-model="uForm.password" type="password" :required="!uForm.sendCredentials" minlength="8" class="an-input" placeholder="Mínimo 8 caracteres" />
              </div>
              <div class="mb-4">
                <label class="an-label">Empresa <span class="text-slate-400 font-normal">(opcional)</span></label>
                <select v-model="uForm.companyId" :disabled="saving || loadingCompanies" class="an-input">
                  <option value="">Sin asignar</option>
                  <option v-for="c in companies" :key="c.id" :value="c.id">{{ c.name }}</option>
                </select>
                <p class="text-[11px] text-slate-400 mt-1">Asocia este usuario a una empresa desde la creación.</p>
              </div>
            </template>

            <div class="an-modal-footer">
              <button v-if="step === 2" type="button" @click="step = 1" class="an-btn-ghost">Atrás</button>
              <button v-else type="button" @click="closeCreate" class="an-btn-ghost">Cancelar</button>
              <button type="submit" :disabled="saving" class="an-btn-primary">
                {{ step === 1 ? 'Siguiente' : (saving ? 'Guardando…' : 'Crear Usuario') }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- ══ Edit modal ═══════════════════════════════════════════════════════ -->
    <Teleport to="body">
      <div v-if="editTarget" class="an-overlay" @click.self="closeEdit">
        <div class="an-modal flex flex-col" style="max-width:520px;max-height:90dvh">

          <!-- Header -->
          <div class="an-modal-header shrink-0">
            <h2 class="text-base font-semibold text-slate-900">
              {{ editShowDemo ? 'Acceso Demo' : 'Editar Usuario' }}
            </h2>
            <div class="flex items-center gap-2">
              <button
                v-if="editShowDemo"
                type="button"
                @click="editShowDemo = false"
                class="text-xs text-slate-400 hover:text-slate-600"
              >← Volver</button>
              <button
                v-if="!editShowDemo && authStore.can('usuarios.demo')"
                type="button"
                @click="editShowDemo = true"
                class="text-[11px] font-medium px-2 py-0.5 rounded border border-amber-200 bg-amber-50 text-amber-700 hover:bg-amber-100 transition-colors"
              >Demo</button>
              <button @click="closeEdit" class="an-close-btn">&times;</button>
            </div>
          </div>

          <!-- Form -->
          <form @submit.prevent="submitEdit" class="flex flex-col min-h-0 flex-1">

            <!-- Scrollable content -->
            <div class="flex-1 overflow-y-auto px-6 py-5">
              <div v-if="editError" class="mb-4 rounded border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">{{ editError }}</div>

              <!-- ── Panel principal ───────────────────────────────── -->
              <template v-if="!editShowDemo">
                <p class="text-[11px] uppercase tracking-wide font-semibold text-slate-400 mb-2">Datos de la persona</p>
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
                  <div>
                    <label class="an-label">Nombre</label>
                    <input v-model="editForm.firstName" type="text" maxlength="100" class="an-input" />
                  </div>
                  <div>
                    <label class="an-label">Primer Apellido</label>
                    <input v-model="editForm.firstSurname" type="text" maxlength="100" class="an-input" />
                  </div>
                </div>
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
                  <div>
                    <label class="an-label">Segundo Apellido</label>
                    <input v-model="editForm.secondSurname" type="text" class="an-input" />
                  </div>
                  <div>
                    <label class="an-label">Identificación</label>
                    <input v-model="editForm.nationalId" type="text" class="an-input" />
                  </div>
                </div>
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
                  <div>
                    <label class="an-label">Fecha de Nacimiento</label>
                    <input v-model="editForm.birthDate" type="date" class="an-input" />
                  </div>
                  <div>
                    <label class="an-label">Código de Área</label>
                    <input v-model="editForm.phoneAreaCode" type="text" class="an-input" />
                  </div>
                </div>
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
                  <div>
                    <label class="an-label">Teléfono Principal</label>
                    <input v-model="editForm.primaryPhone" type="text" class="an-input" />
                  </div>
                  <div>
                    <label class="an-label">Teléfono Secundario</label>
                    <input v-model="editForm.secondaryPhone" type="text" class="an-input" />
                  </div>
                </div>
                <div class="mb-5">
                  <label class="an-label">Dirección</label>
                  <input v-model="editForm.address" type="text" class="an-input" />
                </div>

                <div class="border-t border-slate-100 my-1"></div>

                <p class="text-[11px] uppercase tracking-wide font-semibold text-slate-400 mb-2 mt-4">Cuenta de usuario</p>
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
                  <div>
                    <label class="an-label">Usuario</label>
                    <input v-model="editForm.username" type="text" maxlength="100" class="an-input" />
                  </div>
                  <div>
                    <label class="an-label">Estado</label>
                    <select v-model="editForm.status" class="an-input">
                      <option value="active">Activo</option>
                      <option value="inactive">Inactivo</option>
                    </select>
                  </div>
                </div>
                <div class="mb-4">
                  <label class="an-label">Correo electrónico</label>
                  <input v-model="editForm.email" type="email" class="an-input" />
                </div>
                <div>
                  <label class="an-label">Nueva contraseña</label>
                  <input v-model="editForm.password" type="password" minlength="8" class="an-input" placeholder="Dejar vacío para no cambiar" />
                </div>
              </template>

              <!-- ── Panel demo ────────────────────────────────────── -->
              <template v-else>
                <p class="text-sm text-slate-500 mb-5">
                  Configura una fecha límite de acceso para
                  <span class="font-semibold text-slate-800">{{ editTarget?.username }}</span>.
                  Cuando la fecha pase, el usuario no podrá iniciar sesión en ningún sistema.
                </p>

                <div class="mb-4">
                  <label class="an-label">Expira el</label>
                  <input v-model="editForm.passwordExpiresAt" type="datetime-local" class="an-input" />
                  <p class="text-[11px] text-slate-400 mt-1">Deja vacío para acceso permanente.</p>
                </div>

                <button
                  v-if="editForm.passwordExpiresAt"
                  type="button"
                  @click="editForm.passwordExpiresAt = ''; editForm.clearExpiry = true"
                  class="text-xs text-red-500 hover:text-red-700"
                >Quitar vencimiento</button>

                <div v-if="editTarget?.passwordExpiresAt" class="mt-5 rounded border border-slate-100 bg-slate-50 p-3 text-xs text-slate-500">
                  Vencimiento actual:
                  <span class="font-semibold" :class="isExpired(editTarget.passwordExpiresAt) ? 'text-red-600' : 'text-slate-700'">
                    {{ formatExpiry(editTarget.passwordExpiresAt) }}
                  </span>
                  <span v-if="isExpired(editTarget.passwordExpiresAt)" class="ml-1 text-red-500 font-medium">· Ya expiró</span>
                </div>
              </template>
            </div>

            <!-- Footer fijo -->
            <div class="an-modal-footer shrink-0">
              <button type="button" @click="closeEdit" class="an-btn-ghost">Cancelar</button>
              <button type="submit" :disabled="editSaving" class="an-btn-primary">{{ editSaving ? 'Guardando…' : 'Guardar' }}</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- ══ Delete confirm ═══════════════════════════════════════════════════ -->
    <Teleport to="body">
      <div v-if="deleteTarget" class="an-overlay" @click.self="deleteTarget = null">
        <div class="an-modal an-modal--sm">
          <div class="an-modal-header">
            <h2 class="text-base font-semibold text-slate-900">Eliminar usuario</h2>
            <button @click="deleteTarget = null" class="an-close-btn">&times;</button>
          </div>
          <div class="an-modal-body">
            <p class="text-sm text-slate-600">
              ¿Confirmas la eliminación de
              <span class="font-semibold text-slate-900">{{ deleteTarget.username }}</span>
              ({{ fullName(deleteTarget.person) }})?
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

    <!-- ══ Reset Password modal ══════════════════════════════════════════════ -->
    <Teleport to="body">
      <div v-if="resetPwdTarget" class="an-overlay" @click.self="closeResetPwd">
        <div class="an-modal" style="max-width:460px">
          <div class="an-modal-header">
            <h2 class="text-base font-semibold text-slate-900">Restablecer Contraseña</h2>
            <button @click="closeResetPwd" class="an-close-btn">&times;</button>
          </div>
          <div class="an-modal-body">
            <p class="text-sm text-slate-500 mb-4">
              Usuario: <span class="font-semibold text-slate-800">{{ resetPwdTarget?.username }}</span>
              &mdash; {{ fullName(resetPwdTarget?.person) }}
            </p>

            <div v-if="resetPwdError" class="mb-4 rounded border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">{{ resetPwdError }}</div>

            <!-- Success: show generated password (shown only once) -->
            <div v-if="resetPwdGenerated" class="mb-4 rounded border border-emerald-200 bg-emerald-50 px-4 py-3">
              <p class="text-xs font-semibold text-emerald-700 mb-1">Contraseña generada — cópiela ahora, no volverá a mostrarse</p>
              <div class="flex items-center gap-2 mt-1">
                <code class="flex-1 text-sm font-mono bg-white border border-emerald-200 rounded px-3 py-1.5 text-emerald-800">{{ resetPwdGenerated }}</code>
                <button @click="copyToClipboard(resetPwdGenerated)" class="text-xs px-2 py-1.5 rounded border border-emerald-300 text-emerald-700 hover:bg-emerald-100 transition-colors">{{ copied ? '✔ Copiado' : 'Copiar' }}</button>
              </div>
            </div>

            <div v-else>
              <!-- Manual password entry -->
              <div class="mb-4">
                <label class="an-label">Nueva contraseña</label>
                <div class="relative">
                  <input
                    v-model="resetPwdForm.password"
                    :type="showResetPwd ? 'text' : 'password'"
                    minlength="8"
                    class="an-input pr-20"
                    placeholder="Mínimo 8 caracteres"
                    :disabled="resetPwdSaving"
                  />
                  <button
                    type="button"
                    @click="showResetPwd = !showResetPwd"
                    class="absolute right-2 top-1/2 -translate-y-1/2 text-xs text-slate-400 hover:text-slate-600 px-1"
                  >{{ showResetPwd ? 'Ocultar' : 'Ver' }}</button>
                </div>
              </div>

              <button
                type="button"
                @click="submitResetPwd(true)"
                :disabled="resetPwdSaving"
                class="w-full text-sm py-2 rounded border border-dashed border-slate-300 text-slate-600 hover:border-slate-400 hover:bg-slate-50 transition-colors mb-4"
              >Generar contraseña aleatoria segura</button>
            </div>

            <div class="an-modal-footer">
              <button type="button" @click="closeResetPwd" class="an-btn-ghost">{{ resetPwdGenerated ? 'Cerrar' : 'Cancelar' }}</button>
              <button
                v-if="!resetPwdGenerated"
                type="button"
                @click="submitResetPwd(false)"
                :disabled="resetPwdSaving || !resetPwdForm.password"
                class="an-btn-primary"
              >{{ resetPwdSaving ? 'Guardando…' : 'Guardar contraseña' }}</button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useAuthStore } from '@/store/auth'
import api from '@/services/api'
import { te, teError } from '@/i18n'
import { Pencil, Trash2, KeyRound } from 'lucide-vue-next'

const authStore = useAuthStore()

// ── Helpers ───────────────────────────────────────────────────────────────────
function fullName(p) {
  if (!p) return '—'
  return [p.firstName, p.firstSurname, p.secondSurname].filter(Boolean).join(' ')
}

// ── List + search + pagination ────────────────────────────────────────────────
const users      = ref([])
const loading    = ref(true)
const error      = ref('')
const search     = ref('')
const page       = ref(1)
const pageSize   = 20
const totalItems = ref(0)
const totalPages = ref(0)

const rangeStart = computed(() => totalItems.value === 0 ? 0 : (page.value - 1) * pageSize + 1)
const rangeEnd   = computed(() => Math.min(page.value * pageSize, totalItems.value))

async function loadUsers() {
  loading.value = true
  error.value   = ''
  try {
    const params = { page: page.value, pageSize }
    if (search.value.trim()) params.search = search.value.trim()
    const res = await api.get('/users', { params })
    users.value      = res.data?.data ?? []
    totalItems.value = res.data?.pagination?.totalItems ?? users.value.length
    totalPages.value = res.data?.pagination?.totalPages ?? 1
  } catch (e) {
    error.value = teError(e)
  } finally {
    loading.value = false
  }
}

function prevPage() { if (page.value > 1) { page.value--; loadUsers() } }
function nextPage() { if (page.value < totalPages.value) { page.value++; loadUsers() } }

let searchTimer = null
watch(search, () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    page.value = 1
    loadUsers()
  }, 300)
})

onMounted(loadUsers)

// ── Create (2-step) ──────────────────────────────────────────────────────────
const createOpen = ref(false)
const step       = ref(1)
const saving     = ref(false)
const formError  = ref('')

const pForm = reactive({
  firstName: '', firstSurname: '', secondSurname: '',
  nationalId: '', birthDate: '', phoneAreaCode: '',
  primaryPhone: '', secondaryPhone: '', address: '',
})

const uForm = reactive({ username: '', email: '', password: '', companyId: '', sendCredentials: false })

// ── Companies (for create wizard) ─────────────────────────────────────────────
const companies       = ref([])
const loadingCompanies = ref(false)

async function loadCompanies() {
  loadingCompanies.value = true
  try {
    const res = await api.get('/companies', { params: { pageSize: 200 } })
    companies.value = res.data?.data ?? []
  } catch { /* silently ignore — optional field */ }
  finally { loadingCompanies.value = false }
}

onMounted(loadCompanies)

function openCreate() {
  Object.assign(pForm, {
    firstName: '', firstSurname: '', secondSurname: '',
    nationalId: '', birthDate: '', phoneAreaCode: '',
    primaryPhone: '', secondaryPhone: '', address: '',
  })
  Object.assign(uForm, { username: '', email: '', password: '', companyId: '', sendCredentials: false })
  step.value      = 1
  formError.value = ''
  createOpen.value = true
}

function closeCreate() {
  if (!saving.value) createOpen.value = false
}

function goStep2() {
  formError.value = ''
  if (!pForm.firstName.trim() || !pForm.firstSurname.trim()) {
    formError.value = te('form.first_name_required')
    return
  }
  step.value = 2
}

async function submitCreate() {
  formError.value = ''
  saving.value    = true
  try {
    const payload = {
      firstName:    pForm.firstName.trim(),
      firstSurname: pForm.firstSurname.trim(),
      username:     uForm.username.trim(),
      email:        uForm.email.trim(),
      sendCredentials: uForm.sendCredentials,
    }
    if (!uForm.sendCredentials) {
      payload.password = uForm.password
    }
    // Optional person fields — only send non-empty values
    if (pForm.secondSurname.trim())  payload.secondSurname  = pForm.secondSurname.trim()
    if (pForm.nationalId.trim())     payload.nationalId     = pForm.nationalId.trim()
    if (pForm.birthDate)             payload.birthDate      = pForm.birthDate
    if (pForm.phoneAreaCode.trim())  payload.phoneAreaCode  = pForm.phoneAreaCode.trim()
    if (pForm.primaryPhone.trim())   payload.primaryPhone   = pForm.primaryPhone.trim()
    if (pForm.secondaryPhone.trim()) payload.secondaryPhone = pForm.secondaryPhone.trim()
    if (pForm.address.trim())        payload.address        = pForm.address.trim()
    if (uForm.companyId)             payload.companyId       = Number(uForm.companyId)

    await api.post('/users', payload)
    createOpen.value = false
    page.value = 1
    await loadUsers()
  } catch (e) {
    formError.value = teError(e)
  } finally {
    saving.value = false
  }
}

// ── Edit ──────────────────────────────────────────────────────────────────────
const editTarget   = ref(null)
const editSaving   = ref(false)
const editError    = ref('')
const editShowDemo = ref(false)
const editForm   = reactive({
  // Person
  firstName: '', firstSurname: '', secondSurname: '',
  nationalId: '', birthDate: '', phoneAreaCode: '',
  primaryPhone: '', secondaryPhone: '', address: '',
  // User
  username: '', email: '', password: '', status: 'active',
  // Demo
  passwordExpiresAt: '', clearExpiry: false,
})

function openEdit(u) {
  editTarget.value   = u
  editShowDemo.value = false
  const p = u.person || {}
  editForm.firstName      = p.firstName ?? ''
  editForm.firstSurname   = p.firstSurname ?? ''
  editForm.secondSurname  = p.secondSurname ?? ''
  editForm.nationalId     = p.nationalId ?? ''
  editForm.birthDate      = p.birthDate ? p.birthDate.substring(0, 10) : ''
  editForm.phoneAreaCode  = p.phoneAreaCode ?? ''
  editForm.primaryPhone   = p.primaryPhone ?? ''
  editForm.secondaryPhone = p.secondaryPhone ?? ''
  editForm.address        = p.address ?? ''
  editForm.username          = u.username
  editForm.email             = u.email
  editForm.password          = ''
  editForm.status            = u.status
  editForm.passwordExpiresAt = u.passwordExpiresAt
    ? new Date(u.passwordExpiresAt).toISOString().substring(0, 16) : ''
  editForm.clearExpiry       = false
  editError.value            = ''
}

function closeEdit() {
  if (!editSaving.value) {
    editTarget.value   = null
    editShowDemo.value = false
  }
}

async function submitEdit() {
  editError.value  = ''
  editSaving.value = true
  try {
    const payload = {}
    const orig = editTarget.value
    const op   = orig.person || {}

    // User fields — only changed
    if (editForm.username.trim() !== orig.username) payload.username = editForm.username.trim()
    if (editForm.email.trim() !== orig.email)       payload.email    = editForm.email.trim()
    if (editForm.password)                          payload.password = editForm.password
    if (editForm.status !== orig.status)            payload.status   = editForm.status

    // Demo expiry
    const origExpiry = orig.passwordExpiresAt
      ? new Date(orig.passwordExpiresAt).toISOString().substring(0, 16) : ''
    if (editForm.clearExpiry) {
      payload.clearExpiry = true
    } else if (editForm.passwordExpiresAt !== origExpiry) {
      payload.passwordExpiresAt = editForm.passwordExpiresAt
        ? new Date(editForm.passwordExpiresAt).toISOString() : null
    }

    // Person fields — only changed
    if (editForm.firstName.trim() !== (op.firstName ?? ''))           payload.firstName      = editForm.firstName.trim()
    if (editForm.firstSurname.trim() !== (op.firstSurname ?? ''))     payload.firstSurname   = editForm.firstSurname.trim()
    if (editForm.secondSurname.trim() !== (op.secondSurname ?? ''))   payload.secondSurname  = editForm.secondSurname.trim()
    if (editForm.nationalId.trim() !== (op.nationalId ?? ''))         payload.nationalId     = editForm.nationalId.trim()
    if (editForm.phoneAreaCode.trim() !== (op.phoneAreaCode ?? ''))   payload.phoneAreaCode  = editForm.phoneAreaCode.trim()
    if (editForm.primaryPhone.trim() !== (op.primaryPhone ?? ''))     payload.primaryPhone   = editForm.primaryPhone.trim()
    if (editForm.secondaryPhone.trim() !== (op.secondaryPhone ?? '')) payload.secondaryPhone = editForm.secondaryPhone.trim()
    if (editForm.address.trim() !== (op.address ?? ''))               payload.address        = editForm.address.trim()
    const origBd = op.birthDate ? op.birthDate.substring(0, 10) : ''
    if (editForm.birthDate !== origBd) payload.birthDate = editForm.birthDate || null

    await api.put(`/users/${orig.id}`, payload)
    editTarget.value = null
    await loadUsers()
  } catch (e) {
    editError.value = teError(e)
  } finally {
    editSaving.value = false
  }
}

// ── Delete ────────────────────────────────────────────────────────────────────
const deleteTarget = ref(null)
const deleting     = ref(false)

function confirmDelete(u) {
  deleteTarget.value = u
}

async function doDelete() {
  deleting.value = true
  try {
    await api.delete(`/users/${deleteTarget.value.id}`)
    deleteTarget.value = null
    await loadUsers()
  } catch (e) {
    error.value = teError(e)
    deleteTarget.value = null
  } finally {
    deleting.value = false
  }
}

// ── Reset Password ────────────────────────────────────────────────────────────
const resetPwdTarget  = ref(null)
const resetPwdSaving  = ref(false)
const resetPwdError   = ref('')
const resetPwdGenerated = ref('')
const showResetPwd    = ref(false)
const copied          = ref(false)
const resetPwdForm    = reactive({ password: '' })

function openResetPwd(u) {
  resetPwdTarget.value  = u
  resetPwdGenerated.value = ''
  resetPwdError.value   = ''
  showResetPwd.value    = false
  copied.value          = false
  resetPwdForm.password = ''
}

function closeResetPwd() {
  if (!resetPwdSaving.value) resetPwdTarget.value = null
}

async function submitResetPwd(generateRandom) {
  if (!generateRandom && !resetPwdForm.password) return
  resetPwdSaving.value = true
  resetPwdError.value  = ''
  try {
    const body = generateRandom
      ? { generateRandom: true }
      : { password: resetPwdForm.password, generateRandom: false }
    const res = await api.post(`/users/${resetPwdTarget.value.id}/reset-password`, body)
    if (generateRandom) {
      resetPwdGenerated.value = res.data?.data?.generatedPassword ?? ''
    } else {
      resetPwdTarget.value = null
    }
  } catch (e) {
    resetPwdError.value = teError(e)
  } finally {
    resetPwdSaving.value = false
  }
}

async function copyToClipboard(text) {
  try {
    await navigator.clipboard.writeText(text)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch { /* silent */ }
}

// ── Demo helpers ──────────────────────────────────────────────────────────────
function formatExpiry(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString('es', { day: '2-digit', month: 'short', year: 'numeric' })
}

function isExpired(iso) {
  if (!iso) return false
  return Date.now() > new Date(iso).getTime()
}
</script>