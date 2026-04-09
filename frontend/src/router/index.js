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
      component: () => import('@/views/Login.vue')
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
      component: () => import('@/views/ForgotPassword.vue')
    },
    {
      path: '/reset-password',
      name: 'reset-password',
      component: () => import('@/views/ResetPassword.vue')
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

const publicRoutes = ['landing', 'login', 'unete', 'forgot-password', 'reset-password']

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const isAuthenticated = !!token

  if (!publicRoutes.includes(to.name) && !isAuthenticated) {
    next({ name: 'login' })
  } else if (to.name === 'landing' && isAuthenticated) {
    next({ name: 'dashboard' })
  } else if (to.name === 'login' && isAuthenticated) {
    next({ name: 'dashboard' })
  } else {
    next()
  }
})

export default router