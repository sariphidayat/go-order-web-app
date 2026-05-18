import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import CustomersView from '../views/CustomersView.vue'
import OrdersView from '../views/OrdersView.vue'
import TracksView from '../views/TracksView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', name: 'login', component: LoginView },
    { path: '/', name: 'dashboard', component: DashboardView, meta: { requiresAuth: true } },
    { path: '/customers', name: 'customers', component: CustomersView, meta: { requiresAuth: true } },
    { path: '/orders', name: 'orders', component: OrdersView, meta: { requiresAuth: true } },
    { path: '/tracks', name: 'tracks', component: TracksView, meta: { requiresAuth: true } },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return '/login'
  }
  if (to.name === 'login' && auth.isAuthenticated) {
    return '/'
  }
  return true
})

export default router
