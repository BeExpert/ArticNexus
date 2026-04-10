import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'landing',
      component: () => import('@/views/LandingView.vue'),
      meta: { public: true }
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('@/views/Home.vue')
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/Login.vue'),
      meta: { public: true }
    },
    {
      path: '/unete',
      name: 'unete',
      component: () => import('@/views/JoinView.vue'),
      meta: { public: true }
    },
    {
      path: '/forgot-password',
      name: 'forgot-password',
      component: () => import('@/views/ForgotPassword.vue'),
      meta: { public: true }
    },
    {
      path: '/reset-password',
      name: 'reset-password',
      component: () => import('@/views/ResetPassword.vue'),
      meta: { public: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('@/views/Profile.vue')
    },
    {
      path: '/users',
      name: 'users',
      component: () => import('@/views/Users.vue')
    },
    {
      path: '/companies',
      name: 'companies',
      component: () => import('@/views/Companies.vue'),
      meta: { fullWidth: true }
    },
    {
      path: '/companies/:id',
      redirect: { name: 'companies' }
    },
    {
      path: '/applications',
      name: 'applications',
      component: () => import('@/views/Applications.vue'),
      meta: { fullWidth: true }
    },
    {
      path: '/licenses',
      name: 'licenses',
      component: () => import('@/views/Licenses.vue')
    },
    {
      path: '/demo-links',
      name: 'demo-links',
      component: () => import('@/views/DemoLinks.vue')
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/views/NotFound.vue')
    }
  ]
})

router.beforeEach(async (to, from, next) => {
  const token = localStorage.getItem('token')
  const isAuthenticated = !!token

  // Lazy-import the store to avoid circular dependency at module load time.
  const { useAuthStore } = await import('@/store/auth')
  const authStore = useAuthStore()

  // A user in the middle of company selection must stay on /login.
  if (authStore.pendingCompanySelection) {
    if (to.name !== 'login') return next({ name: 'login' })
    return next()
  }

  if (!to.meta.public && !isAuthenticated) {
    return next({ name: 'login' })
  }
  if ((to.name === 'landing' || to.name === 'login') && isAuthenticated) {
    return next({ name: 'dashboard' })
  }
  next()
})

export default router