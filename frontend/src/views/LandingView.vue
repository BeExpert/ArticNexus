<template>
  <div class="min-h-screen bg-white font-sans">

    <!-- ── Navigation ──────────────────────────────────────────── -->
    <nav class="bg-white/90 backdrop-blur-sm shadow-sm sticky top-0 z-50 border-b border-slate-100">
      <div class="container-responsive">
        <div class="flex justify-between items-center py-4">
          <div class="flex items-center gap-3">
            <img
              src="@/assets/images/articdev-logo.svg"
              alt="ArticDev"
              class="w-10 h-8"
              @error="e => e.target.style.display = 'none'"
            />
            <span class="text-xl font-bold text-dark-gray">ArticDev</span>
          </div>

          <div class="hidden md:flex items-center gap-8">
            <a
              v-for="item in navItems"
              :key="item.href"
              :href="item.href"
              class="text-cool-gray hover:text-dark-gray text-sm font-medium transition-colors"
            >{{ item.label }}</a>
            <RouterLink
              to="/login"
              class="btn btn-primary py-2 px-5 rounded-lg font-semibold text-sm"
            >Accede a tu empresa</RouterLink>
          </div>

          <button
            class="md:hidden p-2 text-dark-gray"
            :aria-expanded="mobileOpen"
            aria-label="Abrir menú"
            @click="mobileOpen = !mobileOpen"
          >
            <X v-if="mobileOpen" :size="24" />
            <Menu v-else :size="24" />
          </button>
        </div>

        <div v-if="mobileOpen" class="md:hidden pb-4 border-t border-slate-100 pt-3">
          <div class="flex flex-col gap-2">
            <a
              v-for="item in navItems"
              :key="item.href"
              :href="item.href"
              class="text-cool-gray hover:text-dark-gray py-2 px-2 rounded text-sm"
              @click="mobileOpen = false"
            >{{ item.label }}</a>
            <RouterLink
              to="/login"
              class="btn btn-primary text-center py-2.5 px-5 rounded-lg font-semibold text-sm mt-2"
              @click="mobileOpen = false"
            >Accede a tu empresa</RouterLink>
          </div>
        </div>
      </div>
    </nav>

    <!-- ── Hero: dual audience ────────────────────────────────── -->
    <section id="inicio" class="bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 py-28">
      <div class="container-responsive">
        <div class="text-center mb-16">
          <span class="inline-block text-nordic-cyan/80 text-xs font-bold uppercase tracking-widest mb-5 bg-nordic-cyan/10 px-4 py-1.5 rounded-full border border-nordic-cyan/20">
            ArticDev S.A.
          </span>
          <h1 class="text-5xl md:text-6xl font-bold text-white mb-6 leading-tight">
            Construimos las aplicaciones
            <span class="text-nordic-cyan block mt-1">que agilizan tu negocio</span>
          </h1>
          <p class="text-slate-400 text-xl max-w-2xl mx-auto leading-relaxed">
            Innovación, calidad, compromiso y la atención que tu empresa merece.
          </p>
        </div>

        <div class="grid md:grid-cols-2 gap-6 max-w-3xl mx-auto">
          <a
            href="#empresas"
            class="group block bg-white/5 hover:bg-white/10 border border-white/10 hover:border-nordic-cyan/40 rounded-2xl p-10 text-center transition-all duration-300 hover:-translate-y-1 cursor-pointer"
          >
            <Building2 :size="44" class="text-nordic-cyan mx-auto mb-5" />
            <h2 class="text-2xl font-bold text-white mb-3">Soy empresa o negocio</h2>
            <p class="text-slate-400 text-sm leading-relaxed mb-6">
              Quiero digitalizar mis procesos o construir un sistema a medida.
            </p>
            <span class="inline-flex items-center gap-2 text-nordic-cyan text-sm font-semibold group-hover:gap-3 transition-all">
              Ver cómo trabajamos <ChevronRight :size="16" />
            </span>
          </a>

          <RouterLink
            to="/unete"
            class="group block bg-white/5 hover:bg-white/10 border border-white/10 hover:border-nordic-cyan/40 rounded-2xl p-10 text-center transition-all duration-300 hover:-translate-y-1"
          >
            <Users :size="44" class="text-nordic-cyan mx-auto mb-5" />
            <h2 class="text-2xl font-bold text-white mb-3">Quiero colaborar</h2>
            <p class="text-slate-400 text-sm leading-relaxed mb-6">
              Tengo una idea o quiero unirme a un proyecto activo del ecosistema.
            </p>
            <span class="inline-flex items-center gap-2 text-nordic-cyan text-sm font-semibold group-hover:gap-3 transition-all">
              Explorar proyectos <ChevronRight :size="16" />
            </span>
          </RouterLink>
        </div>
      </div>
    </section>

    <!-- ── Empresas: Process + Products ──────────────────────── -->
    <section id="empresas" class="py-24 bg-white">
      <div class="container-responsive">

        <div class="text-center mb-20">
          <span class="text-nordic-cyan text-sm font-semibold uppercase tracking-wider">Para empresas</span>
          <h2 class="text-4xl font-bold text-dark-gray mt-2 mb-4">Cómo trabajamos contigo</h2>
          <p class="text-cool-gray max-w-xl mx-auto">
            De la idea a la pantalla, con un proceso claro y un equipo comprometido.
          </p>
        </div>

        <div class="grid md:grid-cols-3 gap-8 mb-24">
          <div
            v-for="(step, index) in processSteps"
            :key="step.title"
            class="relative text-center group"
          >
            <div
              v-if="index < processSteps.length - 1"
              class="hidden md:block absolute top-8 left-[58%] w-full h-px bg-slate-100 z-0"
            ></div>
            <div class="relative z-10 w-16 h-16 bg-ice-blue group-hover:bg-nordic-cyan/10 rounded-2xl flex items-center justify-center mx-auto mb-5 transition-colors border border-nordic-cyan/20">
              <component :is="step.icon" :size="28" class="text-nordic-cyan" />
            </div>
            <div class="text-xs font-bold text-nordic-cyan uppercase tracking-wider mb-2">Paso {{ index + 1 }}</div>
            <h3 class="text-lg font-bold text-dark-gray mb-2">{{ step.title }}</h3>
            <p class="text-cool-gray text-sm leading-relaxed">{{ step.description }}</p>
          </div>
        </div>

        <div class="text-center mb-12">
          <h2 class="text-3xl font-bold text-dark-gray mb-4">Nuestros productos y servicios</h2>
          <p class="text-cool-gray max-w-xl mx-auto text-sm">
            Plataformas activas y desarrollo a medida. Clic en cada tarjeta para más detalle.
          </p>
        </div>

        <!-- Flip cards -->
        <div class="grid md:grid-cols-3 gap-8 mb-12">
          <div
            v-for="product in products"
            :key="product.id"
            class="flip-card hover-lift cursor-pointer"
            @click="product.flipped = !product.flipped"
          >
            <div class="flip-card-inner" :class="{ flipped: product.flipped }">
              <div class="flip-card-front bg-white border border-slate-100 shadow-sm">
                <div class="p-8 h-full flex flex-col justify-between">
                  <div>
                    <div class="w-14 h-14 bg-ice-blue rounded-2xl flex items-center justify-center mb-5 border border-nordic-cyan/15">
                      <component :is="product.icon" :size="28" class="text-nordic-cyan" />
                    </div>
                    <h3 class="text-xl font-bold text-dark-gray mb-2">{{ product.title }}</h3>
                    <p class="text-cool-gray text-sm leading-relaxed">{{ product.shortDescription }}</p>
                  </div>
                  <p class="text-xs text-nordic-cyan mt-4 font-medium">Clic para más detalle →</p>
                </div>
              </div>
              <div class="flip-card-back bg-gradient-to-br from-slate-800 to-slate-900 text-white">
                <div class="p-8 h-full flex flex-col justify-between">
                  <div>
                    <h3 class="text-xl font-bold mb-3">{{ product.title }}</h3>
                    <p class="text-sm mb-5 leading-relaxed text-slate-300">{{ product.fullDescription }}</p>
                    <div class="flex flex-wrap gap-1.5">
                      <span
                        v-for="tech in product.technologies"
                        :key="tech"
                        class="px-2.5 py-1 bg-white/10 rounded-md text-xs text-slate-300"
                      >{{ tech }}</span>
                    </div>
                  </div>
                  <a
                    v-if="product.url && product.url !== '#contacto'"
                    :href="product.url"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="w-full bg-nordic-cyan text-slate-900 px-4 py-3 rounded-lg font-semibold text-center hover:bg-cyan-300 transition-colors block text-sm mt-5"
                    @click.stop
                  >Visitar sitio →</a>
                  <a
                    v-else
                    href="#contacto"
                    class="w-full bg-nordic-cyan text-slate-900 px-4 py-3 rounded-lg font-semibold text-center hover:bg-cyan-300 transition-colors block text-sm mt-5"
                    @click.stop
                  >Contáctanos →</a>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Service cards -->
        <div class="grid md:grid-cols-3 gap-8">
          <div
            v-for="service in services"
            :key="service.title"
            class="bg-slate-50 rounded-xl p-6 text-center border border-slate-100 hover:border-nordic-cyan/30 hover:shadow-md transition-all"
          >
            <div class="w-12 h-12 bg-white rounded-xl flex items-center justify-center mx-auto mb-4 shadow-sm border border-slate-100">
              <component :is="service.icon" :size="22" class="text-cool-gray" />
            </div>
            <h4 class="text-base font-bold text-dark-gray mb-2">{{ service.title }}</h4>
            <p class="text-cool-gray text-sm leading-relaxed">{{ service.description }}</p>
          </div>
        </div>

      </div>
    </section>

    <!-- ── About ──────────────────────────────────────────────── -->
    <section id="nosotros" class="py-20 bg-ice-blue/40">
      <div class="container-responsive">
        <div class="max-w-4xl mx-auto">
          <div class="grid md:grid-cols-2 gap-14 items-center">
            <div>
              <span class="text-nordic-cyan text-sm font-semibold uppercase tracking-wider">Nuestra empresa</span>
              <h2 class="text-4xl font-bold text-dark-gray mt-2 mb-6">¿Qué es ArticDev S.A.?</h2>
              <p class="text-cool-gray text-lg leading-relaxed mb-5">
                ArticDev S.A. es una empresa de desarrollo de software especializada en el
                sector salud y empresas en crecimiento. Construimos plataformas SaaS con el
                mismo rigor técnico y atención al detalle en cada proyecto.
              </p>
              <p class="text-cool-gray text-lg leading-relaxed">
                Nuestro equipo combina experiencia en Go, Vue.js y PostgreSQL para entregar
                sistemas rápidos, seguros y escalables que realmente resuelven los problemas
                del día a día de nuestros clientes.
              </p>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div
                v-for="stat in stats"
                :key="stat.label"
                class="bg-white rounded-2xl p-6 text-center shadow-sm border border-slate-100"
              >
                <div class="w-10 h-10 bg-ice-blue rounded-xl flex items-center justify-center mx-auto mb-3">
                  <component :is="stat.icon" :size="20" class="text-nordic-cyan" />
                </div>
                <div class="text-2xl font-bold text-dark-gray mb-1">{{ stat.value }}</div>
                <div class="text-sm text-cool-gray">{{ stat.label }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ── Business Contact Form ──────────────────────────────── -->
    <section id="contacto" class="py-24 bg-white">
      <div class="container-responsive">
        <div class="max-w-2xl mx-auto">
          <div class="text-center mb-12">
            <span class="text-nordic-cyan text-sm font-semibold uppercase tracking-wider">Contacto</span>
            <h2 class="text-4xl font-bold text-dark-gray mt-2 mb-4">Hablemos de tu proyecto</h2>
            <p class="text-cool-gray">
              Cuéntanos qué necesitas. Te respondemos en menos de 24 horas.
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
            <p class="opacity-90">Nos pondremos en contacto contigo pronto.</p>
          </div>

          <form
            v-if="!formSuccess"
            class="bg-white rounded-2xl border border-slate-100 shadow-sm p-8 space-y-5"
            @submit.prevent="submitApplication"
          >
            <!-- Honeypot -->
            <div style="position:absolute;left:-9999px;visibility:hidden">
              <input v-model="appForm.honeypot" type="text" tabindex="-1" autocomplete="off" />
            </div>

            <div class="grid md:grid-cols-2 gap-5">
              <div>
                <label class="block text-sm font-semibold text-dark-gray mb-1.5">Nombre completo *</label>
                <input
                  v-model="appForm.nombre"
                  type="text"
                  maxlength="100"
                  required
                  placeholder="Tu nombre"
                  :class="inputClass('nombre')"
                  @blur="validateAppField('nombre')"
                />
                <p v-if="appErrors.nombre" class="text-red-500 text-xs mt-1">{{ appErrors.nombre }}</p>
              </div>
              <div>
                <label class="block text-sm font-semibold text-dark-gray mb-1.5">Correo electrónico *</label>
                <input
                  v-model="appForm.email"
                  type="email"
                  maxlength="150"
                  required
                  placeholder="tu@empresa.com"
                  :class="inputClass('email')"
                  @blur="validateAppField('email')"
                />
                <p v-if="appErrors.email" class="text-red-500 text-xs mt-1">{{ appErrors.email }}</p>
              </div>
            </div>

            <div>
              <label class="block text-sm font-semibold text-dark-gray mb-1.5">¿Qué necesita tu empresa? *</label>
              <textarea
                v-model="appForm.descripcion"
                rows="5"
                maxlength="1000"
                required
                placeholder="Describe el sistema o proceso que quieres digitalizar, el tamaño de tu equipo, cualquier detalle relevante..."
                :class="inputClass('descripcion')"
                @blur="validateAppField('descripcion')"
              ></textarea>
              <p v-if="appErrors.descripcion" class="text-red-500 text-xs mt-1">{{ appErrors.descripcion }}</p>
            </div>

            <p v-if="appErrors.general" class="text-red-500 text-sm text-center">{{ appErrors.general }}</p>

            <button
              type="submit"
              :disabled="appSubmitting"
              class="w-full btn btn-primary py-4 rounded-xl text-base font-bold flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="!appSubmitting" class="flex items-center gap-2">
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

    <!-- ── Nexus Gateway ──────────────────────────────────────── -->
    <section id="nexus" class="py-20 bg-gradient-to-br from-slate-900 to-slate-800 text-center">
      <div class="container-responsive">
        <h2 class="text-3xl md:text-4xl font-bold text-white mb-4">
          ¿Ya formas parte del equipo?
        </h2>
        <p class="text-slate-400 mb-10 max-w-xl mx-auto leading-relaxed">
          Accede a la plataforma interna ArticNexus o
          <RouterLink to="/unete" class="text-nordic-cyan hover:underline">únete al ecosistema</RouterLink>.
        </p>
        <RouterLink
          to="/login"
          class="inline-flex items-center gap-2 bg-nordic-cyan hover:bg-cyan-300 text-slate-900 font-bold py-4 px-10 rounded-xl text-lg transition-all duration-300 hover:scale-105"
        >
          <LogIn :size="20" />
          Accede a tu empresa
        </RouterLink>
      </div>
    </section>

    <!-- ── Footer ─────────────────────────────────────────────── -->
    <footer class="bg-slate-900 text-white py-12 border-t border-slate-800">
      <div class="container-responsive">
        <div class="grid md:grid-cols-3 gap-8">
          <div>
            <div class="flex items-center gap-3 mb-4">
              <img
                src="@/assets/images/articdev-icon.svg"
                alt="ArticDev"
                class="w-9 h-7"
                @error="e => e.target.style.display = 'none'"
              />
              <h3 class="text-lg font-bold">ArticDev S.A.</h3>
            </div>
            <p class="text-slate-400 text-sm leading-relaxed">
              Construimos software especializado para el sector salud y empresas en crecimiento.
            </p>
          </div>

          <div>
            <h4 class="font-semibold mb-4 text-nordic-cyan text-sm">Productos</h4>
            <ul class="space-y-2 text-slate-400 text-sm">
              <li>
                <a
                  :href="vetDataUrl || '#'"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="hover:text-white transition-colors"
                >VetData — Gestión Veterinaria</a>
              </li>
              <li>
                <a
                  :href="oftaDataUrl || '#'"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="hover:text-white transition-colors"
                >OftaData — Gestión Oftalmológica</a>
              </li>
            </ul>
          </div>

          <div>
            <h4 class="font-semibold mb-4 text-nordic-cyan text-sm">Contacto</h4>
            <ul class="space-y-2 text-slate-400 text-sm">
              <li>
                <a href="mailto:articdevsa@gmail.com" class="hover:text-white transition-colors">
                  articdevsa@gmail.com
                </a>
              </li>
              <li>
                <RouterLink to="/unete" class="hover:text-white transition-colors">
                  Únete al ecosistema
                </RouterLink>
              </li>
            </ul>
          </div>
        </div>

        <div class="border-t border-slate-800 mt-8 pt-8 text-center text-slate-500 text-sm">
          <p>&copy; {{ currentYear }} ArticDev S.A. Todos los derechos reservados.</p>
        </div>
      </div>
    </footer>

  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { RouterLink } from 'vue-router'
import { te } from '@/i18n'
import {
  Menu, X,
  Building2, Users, ChevronRight,
  MessageSquare, Code2, Rocket,
  PawPrint, Eye, Monitor,
  Bot, Globe, Settings2,
  Layers, ShieldCheck,
  CheckCircle2, Send, LogIn,
} from 'lucide-vue-next'

const vetDataUrl  = import.meta.env.VITE_VETDATA_URL  || ''
const oftaDataUrl = import.meta.env.VITE_OFTADATA_URL || ''
const apiBase     = import.meta.env.VITE_API_BASE_URL || '/api/v1'
const currentYear = new Date().getFullYear()

const mobileOpen = ref(false)

const navItems = [
  { href: '#inicio',   label: 'Inicio'    },
  { href: '#empresas', label: 'Servicios' },
  { href: '#nosotros', label: 'Nosotros'  },
  { href: '#contacto', label: 'Contacto'  },
]

const processSteps = [
  {
    icon: MessageSquare,
    title: 'Reunión y diagnóstico',
    description: 'Entendemos tu negocio y tus procesos actuales. Te asesoramos sin tecnicismos.',
  },
  {
    icon: Code2,
    title: 'Diseño y desarrollo',
    description: 'Construimos la solución con tecnología robusta, iterando junto a vos en cada etapa.',
  },
  {
    icon: Rocket,
    title: 'Entrega y soporte',
    description: 'Desplegamos el sistema, capacitamos al equipo y acompañamos el crecimiento.',
  },
]

const products = reactive([
  {
    id: 1,
    icon: PawPrint,
    title: 'VetData',
    shortDescription: 'Plataforma integral para clínicas veterinarias: pacientes, historial clínico, vacunas y más.',
    fullDescription:
      'Gestión completa de mascotas, historial clínico, carnet de vacunas con código QR, portal del cliente, agenda de citas e inventario.',
    technologies: ['Vue 3', 'Go', 'PostgreSQL'],
    url: vetDataUrl,
    flipped: false,
  },
  {
    id: 2,
    icon: Eye,
    title: 'OftaData',
    shortDescription: 'Sistema de gestión integral para clínicas oftalmológicas con formularios especializados.',
    fullDescription:
      'Gestión de pacientes, formularios clínicos especializados, agenda, control de cirugías y catálogos de diagnósticos y tratamientos.',
    technologies: ['Vue 3', 'Go', 'PostgreSQL'],
    url: oftaDataUrl,
    flipped: false,
  },
  {
    id: 3,
    icon: Monitor,
    title: 'Desarrollo a medida',
    shortDescription: 'Sistemas de gestión, aplicaciones web y APIs personalizadas para tu operación.',
    fullDescription:
      'Desde sistemas internos de gestión hasta plataformas SaaS completas. Analizamos tu caso y entregamos una solución escalable.',
    technologies: ['Vue.js', 'Go', 'React', 'Node.js'],
    url: '#contacto',
    flipped: false,
  },
])

const services = [
  {
    icon: Bot,
    title: 'Bot de facturación electrónica',
    description: 'Automatización del proceso de emisión de facturas electrónicas según la normativa local.',
  },
  {
    icon: Globe,
    title: 'Páginas corporativas',
    description: 'Sitios de presentación optimizados para conversión, velocidad y posicionamiento.',
  },
  {
    icon: Settings2,
    title: 'Web Services y APIs',
    description: 'Integraciones, APIs RESTful y servicios backend robustos construidos en Go.',
  },
]

const stats = [
  { icon: Layers,      value: '3+',   label: 'Productos SaaS'     },
  { icon: ShieldCheck, value: '100%', label: 'Código Propio'       },
  { icon: Code2,       value: 'Go',   label: 'Backend Principal'   },
  { icon: Rocket,      value: '24h',  label: 'Tiempo de respuesta' },
]

// ── Contact form (tipo fijo: Solicitud de proyecto) ────────────
const formSuccess   = ref(false)
const appSubmitting = ref(false)

const appForm = reactive({
  nombre: '', email: '', descripcion: '', honeypot: '',
})

const appErrors = reactive({})

function validateAppField(field) {
  delete appErrors[field]
  const v = appForm[field]?.trim?.() ?? ''
  if (field === 'nombre') {
    if (!v) appErrors.nombre = 'El nombre es requerido'
    else if (v.length < 2) appErrors.nombre = 'Mínimo 2 caracteres'
  }
  if (field === 'email') {
    if (!v) appErrors.email = 'El correo es requerido'
    else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v)) appErrors.email = 'Correo no válido'
  }
  if (field === 'descripcion') {
    if (!v) appErrors.descripcion = 'La descripción es requerida'
    else if (v.length < 20) appErrors.descripcion = 'Mínimo 20 caracteres'
  }
}

function validateAll() {
  ;['nombre', 'email', 'descripcion'].forEach(validateAppField)
  return Object.keys(appErrors).filter(k => k !== 'general').length === 0
}

function inputClass(field) {
  const base = 'w-full px-4 py-3 border rounded-xl focus:ring-2 focus:ring-nordic-cyan focus:border-transparent transition-all text-sm'
  if (appErrors[field]) return `${base} border-red-400`
  if (appForm[field]?.trim?.()) return `${base} border-green-400`
  return `${base} border-slate-200`
}

async function submitApplication() {
  if (!validateAll()) return
  if (appForm.honeypot) return

  appSubmitting.value = true
  delete appErrors.general

  try {
    const response = await fetch(`${apiBase}/contact`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'X-Requested-With': 'XMLHttpRequest' },
      body: JSON.stringify({
        nombre:      appForm.nombre.trim(),
        email:       appForm.email.trim(),
        tipo:        'Solicitud de proyecto',
        descripcion: appForm.descripcion.trim(),
        honeypot:    appForm.honeypot,
      }),
    })

    if (!response.ok) {
      const data = await response.json().catch(() => ({}))
      appErrors.general = data.code ? te(data.code) : te('contact.send_failed')
      appSubmitting.value = false
      return
    }

    formSuccess.value = true
  } catch {
    appErrors.general = te('network.connection_error')
  } finally {
    appSubmitting.value = false
  }
}
</script>

<style scoped>
.font-sans { font-family: 'Inter', system-ui, sans-serif; }
</style>
