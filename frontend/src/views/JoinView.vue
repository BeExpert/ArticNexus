<template>
  <div class="min-h-screen bg-white font-sans">

    <!-- ── Navigation ──────────────────────────────────────────── -->
    <nav class="bg-white/90 backdrop-blur-sm shadow-sm sticky top-0 z-50 border-b border-slate-100">
      <div class="container-responsive">
        <div class="flex justify-between items-center py-4">
          <RouterLink to="/" class="flex items-center gap-3 group">
            <img
              src="@/assets/images/articdev-logo.svg"
              alt="ArticDev"
              class="w-10 h-8"
              @error="e => e.target.style.display = 'none'"
            />
            <span class="text-xl font-bold text-dark-gray group-hover:text-nordic-cyan transition-colors">ArticDev</span>
          </RouterLink>

          <div class="flex items-center gap-4">
            <RouterLink
              to="/"
              class="hidden sm:inline-flex items-center gap-1.5 text-cool-gray hover:text-dark-gray text-sm font-medium transition-colors"
            >
              <ArrowLeft :size="16" />
              Volver al inicio
            </RouterLink>
            <RouterLink
              to="/login"
              class="btn btn-primary py-2 px-5 rounded-lg font-semibold text-sm"
            >Accede a tu empresa</RouterLink>
          </div>
        </div>
      </div>
    </nav>

    <!-- ── Hero ──────────────────────────────────────────────────── -->
    <section class="bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 py-24">
      <div class="container-responsive text-center">
        <span class="inline-block text-nordic-cyan/80 text-xs font-bold uppercase tracking-widest mb-5 bg-nordic-cyan/10 px-4 py-1.5 rounded-full border border-nordic-cyan/20">
          Ecosistema ArticDev
        </span>
        <h1 class="text-5xl md:text-6xl font-bold text-white mb-6 leading-tight">
          Construyamos
          <span class="text-nordic-cyan block mt-1">juntos</span>
        </h1>
        <p class="text-slate-400 text-xl max-w-2xl mx-auto leading-relaxed">
          Somos un equipo técnico en crecimiento. Si tenés una idea o querés
          contribuir a uno de nuestros proyectos activos, este es el lugar.
        </p>
      </div>
    </section>

    <!-- ── Options ───────────────────────────────────────────────── -->
    <section class="py-20 bg-white">
      <div class="container-responsive">
        <div class="text-center mb-14">
          <h2 class="text-3xl font-bold text-dark-gray mb-3">¿Cómo querés participar?</h2>
          <p class="text-cool-gray max-w-lg mx-auto">
            Hay dos formas de ser parte del ecosistema ArticDev.
          </p>
        </div>

        <div class="grid md:grid-cols-2 gap-8 max-w-3xl mx-auto">
          <div class="bg-slate-50 border border-slate-100 rounded-2xl p-8 hover:border-nordic-cyan/30 hover:shadow-md transition-all">
            <div class="w-14 h-14 bg-white rounded-2xl flex items-center justify-center mb-5 shadow-sm border border-slate-100">
              <Lightbulb :size="26" class="text-nordic-cyan" />
            </div>
            <h3 class="text-xl font-bold text-dark-gray mb-3">Proponer un proyecto</h3>
            <p class="text-cool-gray text-sm leading-relaxed mb-5">
              Tiene una idea de producto o sistema y busca un equipo técnico
              para desarrollarla. Nos comparte su visión, evaluamos el alcance
              y avanzamos juntos si hay buena compatibilidad.
            </p>
            <ul class="space-y-2 text-sm text-cool-gray">
              <li class="flex items-center gap-2">
                <CheckCircle2 :size="16" class="text-nordic-cyan flex-shrink-0" />
                Evaluamos la idea sin costo
              </li>
              <li class="flex items-center gap-2">
                <CheckCircle2 :size="16" class="text-nordic-cyan flex-shrink-0" />
                Definimos el alcance inicial juntos
              </li>
              <li class="flex items-center gap-2">
                <CheckCircle2 :size="16" class="text-nordic-cyan flex-shrink-0" />
                Participación en el producto final
              </li>
            </ul>
          </div>

          <div class="bg-slate-50 border border-slate-100 rounded-2xl p-8 hover:border-nordic-cyan/30 hover:shadow-md transition-all">
            <div class="w-14 h-14 bg-white rounded-2xl flex items-center justify-center mb-5 shadow-sm border border-slate-100">
              <UserPlus :size="26" class="text-nordic-cyan" />
            </div>
            <h3 class="text-xl font-bold text-dark-gray mb-3">Unirse a un proyecto</h3>
            <p class="text-cool-gray text-sm leading-relaxed mb-5">
              Querés contribuir a VetData, OftaData u otro proyecto activo del
              ecosistema. Podés sumar como desarrollador, diseñador o en otra
              área que agregue valor al equipo.
            </p>
            <ul class="space-y-2 text-sm text-cool-gray">
              <li class="flex items-center gap-2">
                <CheckCircle2 :size="16" class="text-nordic-cyan flex-shrink-0" />
                Proyectos reales en producción
              </li>
              <li class="flex items-center gap-2">
                <CheckCircle2 :size="16" class="text-nordic-cyan flex-shrink-0" />
                Stack moderno: Vue 3, Go, PostgreSQL
              </li>
              <li class="flex items-center gap-2">
                <CheckCircle2 :size="16" class="text-nordic-cyan flex-shrink-0" />
                Equipo pequeño con autonomía real
              </li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <!-- ── Active Projects ───────────────────────────────────────── -->
    <section class="py-20 bg-ice-blue/40">
      <div class="container-responsive">
        <div class="text-center mb-14">
          <span class="text-nordic-cyan text-sm font-semibold uppercase tracking-wider">Proyectos activos</span>
          <h2 class="text-3xl font-bold text-dark-gray mt-2 mb-3">En desarrollo activo</h2>
          <p class="text-cool-gray max-w-lg mx-auto text-sm">
            Estos son los proyectos donde podés contribuir hoy.
          </p>
        </div>

        <div class="grid md:grid-cols-2 gap-8 max-w-3xl mx-auto">
          <div
            v-for="project in activeProjects"
            :key="project.name"
            class="bg-white rounded-2xl border border-slate-100 shadow-sm p-8 hover:border-nordic-cyan/30 transition-all"
          >
            <div class="flex items-start gap-4 mb-5">
              <div class="w-12 h-12 bg-ice-blue rounded-xl flex items-center justify-center flex-shrink-0 border border-nordic-cyan/15">
                <component :is="project.icon" :size="24" class="text-nordic-cyan" />
              </div>
              <div>
                <h3 class="text-xl font-bold text-dark-gray">{{ project.name }}</h3>
                <span class="inline-flex items-center gap-1.5 bg-green-50 text-green-700 text-xs font-semibold px-2.5 py-1 rounded-full mt-1">
                  <span class="w-1.5 h-1.5 bg-green-500 rounded-full"></span>
                  {{ project.status }}
                </span>
              </div>
            </div>
            <p class="text-cool-gray text-sm leading-relaxed mb-5">{{ project.description }}</p>
            <div class="flex flex-wrap gap-1.5 mb-5">
              <span
                v-for="tech in project.stack"
                :key="tech"
                class="inline-flex items-center gap-1 bg-slate-50 text-slate-600 text-xs font-medium px-2.5 py-1 rounded-md border border-slate-100"
              >
                <Code2 :size="11" />
                {{ tech }}
              </span>
            </div>
            <div class="space-y-1.5">
              <div
                v-for="area in project.openAreas"
                :key="area"
                class="flex items-center gap-2 text-sm text-cool-gray"
              >
                <ChevronRight :size="14" class="text-nordic-cyan flex-shrink-0" />
                {{ area }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ── Application Form ──────────────────────────────────────── -->
    <section class="py-24 bg-white">
      <div class="container-responsive">
        <div class="max-w-2xl mx-auto">
          <div class="text-center mb-12">
            <span class="text-nordic-cyan text-sm font-semibold uppercase tracking-wider">Solicitud</span>
            <h2 class="text-4xl font-bold text-dark-gray mt-2 mb-4">Cuéntanos tu propuesta</h2>
            <p class="text-cool-gray">
              Completá el formulario y te respondemos en menos de 48 horas.
            </p>
          </div>

          <div
            v-if="formSuccess"
            class="bg-gradient-to-r from-nordic-cyan to-blue-500 text-white p-10 rounded-2xl text-center"
          >
            <div class="w-14 h-14 bg-white/20 rounded-full flex items-center justify-center mx-auto mb-4">
              <CheckCircle2 :size="28" class="text-white" />
            </div>
            <h3 class="text-xl font-bold mb-2">¡Solicitud enviada!</h3>
            <p class="opacity-90">Revisamos tu mensaje y nos ponemos en contacto pronto.</p>
          </div>

          <form
            v-if="!formSuccess"
            class="bg-white rounded-2xl border border-slate-100 shadow-sm p-8 space-y-5"
            @submit.prevent="submitForm"
          >
            <!-- Honeypot -->
            <div style="position:absolute;left:-9999px;visibility:hidden">
              <input v-model="form.honeypot" type="text" tabindex="-1" autocomplete="off" />
            </div>

            <div class="grid md:grid-cols-2 gap-5">
              <div>
                <label class="block text-sm font-semibold text-dark-gray mb-1.5">Nombre completo *</label>
                <input
                  v-model="form.nombre"
                  type="text"
                  maxlength="100"
                  required
                  placeholder="Tu nombre"
                  :class="inputClass('nombre')"
                  @blur="validateField('nombre')"
                />
                <p v-if="errors.nombre" class="text-red-500 text-xs mt-1">{{ errors.nombre }}</p>
              </div>
              <div>
                <label class="block text-sm font-semibold text-dark-gray mb-1.5">Correo electrónico *</label>
                <input
                  v-model="form.email"
                  type="email"
                  maxlength="150"
                  required
                  placeholder="tu@email.com"
                  :class="inputClass('email')"
                  @blur="validateField('email')"
                />
                <p v-if="errors.email" class="text-red-500 text-xs mt-1">{{ errors.email }}</p>
              </div>
            </div>

            <div>
              <label class="block text-sm font-semibold text-dark-gray mb-1.5">Tipo de solicitud *</label>
              <select
                v-model="form.tipo"
                required
                class="w-full px-4 py-3 border border-slate-200 rounded-xl focus:ring-2 focus:ring-nordic-cyan focus:border-transparent transition-all text-sm bg-white"
              >
                <option value="" disabled>Seleccioná una opción</option>
                <option value="Proponer un nuevo proyecto">Proponer un nuevo proyecto</option>
                <option value="Unirme a un proyecto existente">Unirme a un proyecto existente</option>
              </select>
              <p v-if="errors.tipo" class="text-red-500 text-xs mt-1">{{ errors.tipo }}</p>
            </div>

            <div>
              <label class="block text-sm font-semibold text-dark-gray mb-1.5">
                {{ form.tipo === 'Proponer un nuevo proyecto' ? 'Describe tu idea *' : 'Cuéntanos cómo querés contribuir *' }}
              </label>
              <textarea
                v-model="form.descripcion"
                rows="5"
                maxlength="1000"
                required
                :placeholder="form.tipo === 'Unirme a un proyecto existente'
                  ? 'Describe tu experiencia, las tecnologías que manejás y en qué proyecto te interesa participar...'
                  : 'Describe la idea, el problema que resuelve, el público objetivo y cualquier detalle relevante...'"
                :class="inputClass('descripcion')"
                @blur="validateField('descripcion')"
              ></textarea>
              <p v-if="errors.descripcion" class="text-red-500 text-xs mt-1">{{ errors.descripcion }}</p>
            </div>

            <p v-if="errors.general" class="text-red-500 text-sm text-center">{{ errors.general }}</p>

            <button
              type="submit"
              :disabled="submitting"
              class="w-full btn btn-primary py-4 rounded-xl text-base font-bold flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="!submitting" class="flex items-center gap-2">
                <Send :size="18" /> Enviar solicitud
              </span>
              <span v-else class="flex items-center gap-2">
                <svg class="animate-spin w-5 h-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
                </svg>
                Enviando...
              </span>
            </button>
          </form>
        </div>
      </div>
    </section>

    <!-- ── Footer minimal ────────────────────────────────────────── -->
    <footer class="bg-slate-900 text-slate-400 py-8 border-t border-slate-800">
      <div class="container-responsive flex flex-col sm:flex-row justify-between items-center gap-3 text-sm">
        <div class="flex items-center gap-2">
          <img
            src="@/assets/images/articdev-icon.svg"
            alt="ArticDev"
            class="w-6 h-5 opacity-70"
            @error="e => e.target.style.display = 'none'"
          />
          <span class="font-semibold text-white">ArticDev S.A.</span>
        </div>
        <p>&copy; {{ currentYear }} ArticDev S.A. Todos los derechos reservados.</p>
        <RouterLink to="/" class="hover:text-white transition-colors">Volver al inicio</RouterLink>
      </div>
    </footer>

  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { RouterLink } from 'vue-router'
import { te } from '@/i18n'
import {
  ArrowLeft,
  Lightbulb, UserPlus, CheckCircle2,
  PawPrint, Eye, Code2, ChevronRight,
  Send,
} from 'lucide-vue-next'

const apiBase     = import.meta.env.VITE_API_BASE_URL || '/api/v1'
const currentYear = new Date().getFullYear()

const activeProjects = [
  {
    icon: PawPrint,
    name: 'VetData',
    status: 'En desarrollo activo',
    description:
      'Plataforma SaaS para clínicas veterinarias. Gestión de pacientes, historial clínico, carnet de vacunas con QR, portal del cliente, agenda e inventario.',
    stack: ['Vue 3', 'Go', 'PostgreSQL', 'Tailwind'],
    openAreas: [
      'Desarrollo frontend (Vue 3)',
      'Desarrollo backend (Go)',
      'UI/UX y diseño de interfaces',
    ],
  },
  {
    icon: Eye,
    name: 'OftaData',
    status: 'En desarrollo activo',
    description:
      'Sistema de gestión para clínicas oftalmológicas. Formularios clínicos especializados, agenda, cirugías, catálogos y notificaciones.',
    stack: ['Vue 3', 'Go', 'PostgreSQL', 'Tailwind'],
    openAreas: [
      'Desarrollo frontend (Vue 3)',
      'Desarrollo backend (Go)',
      'Modelado de datos clínicos',
    ],
  },
]

// ── Form ──────────────────────────────────────────────────────
const formSuccess = ref(false)
const submitting  = ref(false)

const form = reactive({
  nombre: '', email: '', tipo: '', descripcion: '', honeypot: '',
})

const errors = reactive({})

function validateField(field) {
  delete errors[field]
  const v = form[field]?.trim?.() ?? ''
  if (field === 'nombre') {
    if (!v) errors.nombre = te('form.name_required')
    else if (v.length < 2) errors.nombre = te('form.min_2_chars')
  }
  if (field === 'email') {
    if (!v) errors.email = te('form.field_required')
    else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v)) errors.email = te('form.invalid_email')
  }
  if (field === 'tipo') {
    if (!v) errors.tipo = te('form.select_option')
  }
  if (field === 'descripcion') {
    if (!v) errors.descripcion = te('form.field_required')
    else if (v.length < 20) errors.descripcion = te('form.min_20_chars')
  }
}

function validateAll() {
  ;['nombre', 'email', 'tipo', 'descripcion'].forEach(validateField)
  return Object.keys(errors).filter(k => k !== 'general').length === 0
}

function inputClass(field) {
  const base = 'w-full px-4 py-3 border rounded-xl focus:ring-2 focus:ring-nordic-cyan focus:border-transparent transition-all text-sm'
  if (errors[field]) return `${base} border-red-400`
  if (form[field]?.trim?.()) return `${base} border-green-400`
  return `${base} border-slate-200`
}

async function submitForm() {
  if (!validateAll()) return
  if (form.honeypot) return

  submitting.value = true
  delete errors.general

  try {
    const response = await fetch(`${apiBase}/contact`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'X-Requested-With': 'XMLHttpRequest' },
      body: JSON.stringify({
        nombre:      form.nombre.trim(),
        email:       form.email.trim(),
        tipo:        form.tipo,
        descripcion: form.descripcion.trim(),
        honeypot:    form.honeypot,
      }),
    })

    if (!response.ok) {
      const data = await response.json().catch(() => ({}))
      errors.general = data.code ? te(data.code) : te('contact.send_failed')
      submitting.value = false
      return
    }

    formSuccess.value = true
  } catch (err) {
    errors.general = te('network.connection_error')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.font-sans { font-family: 'Inter', system-ui, sans-serif; }
</style>
