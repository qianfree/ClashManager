import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/setup',
    name: 'Setup',
    component: () => import('@/views/Setup.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    redirect: '/nodes',
    meta: { requiresAuth: true }
  },
  {
    path: '/nodes',
    name: 'Nodes',
    component: () => import('@/views/Nodes.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/rules',
    name: 'Rules',
    component: () => import('@/views/Rules.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/groups',
    name: 'Groups',
    component: () => import('@/views/Groups.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/subscription',
    name: 'Subscription',
    component: () => import('@/views/Subscription.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/subscription-logs',
    name: 'SubscriptionLogs',
    component: () => import('@/views/SubscriptionLogs.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/sources',
    name: 'Sources',
    component: () => import('@/views/Sources.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/views/Settings.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/password',
    name: 'Password',
    component: () => import('@/views/Password.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  const requiresAuth = to.meta.requiresAuth !== false

  if (requiresAuth && !userStore.isLoggedIn()) {
    // Need to be logged in
    next('/login')
  } else if ((to.path === '/login' || to.path === '/setup') && userStore.isLoggedIn()) {
    // Already logged in, don't show login/setup
    next('/')
  } else {
    next()
  }
})

export default router
