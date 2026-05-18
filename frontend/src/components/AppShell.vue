<template>
  <div class="page shell">
    <aside class="sidebar">
      <div class="brand">
        <div class="brand-mark">CO</div>
        <div>
          <strong>OrderOps</strong>
          <span>Customer orders</span>
        </div>
      </div>

      <nav class="nav">
        <RouterLink v-for="item in navItems" :key="item.path" :to="item.path" class="nav-item">
          <n-icon :component="item.icon" />
          <span>{{ item.label }}</span>
        </RouterLink>
      </nav>
    </aside>

    <main class="main">
      <header class="topbar">
        <div>
          <h1>{{ title }}</h1>
          <p>{{ subtitle }}</p>
        </div>
        <div class="profile">
          <span>{{ auth.user?.name || 'User' }}</span>
          <n-button secondary @click="logout">
            <template #icon><n-icon :component="LogOutOutline" /></template>
            Logout
          </n-button>
        </div>
      </header>

      <section class="content-wrap main-content">
        <slot />
      </section>
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NButton, NIcon } from 'naive-ui'
import {
  AnalyticsOutline,
  BagHandleOutline,
  LogOutOutline,
  PeopleOutline,
  TrailSignOutline,
} from '@vicons/ionicons5'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const navItems = [
  { path: '/', label: 'Dashboard', icon: AnalyticsOutline },
  { path: '/customers', label: 'Customers', icon: PeopleOutline },
  { path: '/orders', label: 'Orders', icon: BagHandleOutline },
  { path: '/tracks', label: 'Tracking', icon: TrailSignOutline },
]

const pageCopy = {
  dashboard: ['Dashboard', 'Ringkasan customer, order, dan tracking terbaru.'],
  customers: ['Customers', 'Kelola master data customer.'],
  orders: ['Orders', 'Kelola order dan status pemrosesan.'],
  tracks: ['Tracking', 'Catat pergerakan dan status order.'],
}

const title = computed(() => pageCopy[route.name]?.[0] || 'OrderOps')
const subtitle = computed(() => pageCopy[route.name]?.[1] || '')

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.shell {
  display: grid;
  grid-template-columns: 248px 1fr;
}

.sidebar {
  min-height: 100vh;
  padding: 22px 16px;
  background: #101828;
  color: #fff;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 28px;
}

.brand-mark {
  display: grid;
  width: 42px;
  height: 42px;
  place-items: center;
  border-radius: 8px;
  background: #17a2a4;
  font-weight: 800;
}

.brand strong,
.brand span {
  display: block;
}

.brand span {
  color: #98a2b3;
  font-size: 12px;
}

.nav {
  display: grid;
  gap: 6px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: 42px;
  padding: 0 12px;
  border-radius: 8px;
  color: #d0d5dd;
  font-weight: 600;
}

.nav-item.router-link-active {
  background: #1d2939;
  color: #fff;
}

.main {
  min-width: 0;
}

.topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 22px max(24px, calc((100vw - 248px - 1180px) / 2 + 24px));
  background: #fff;
  border-bottom: 1px solid #e4e7ec;
}

.topbar h1 {
  margin: 0;
  font-size: 24px;
}

.topbar p {
  margin: 4px 0 0;
  color: #667085;
}

.profile {
  display: flex;
  align-items: center;
  gap: 12px;
  white-space: nowrap;
}

.main-content {
  padding: 24px 0 40px;
}

@media (max-width: 860px) {
  .shell {
    display: block;
  }

  .sidebar {
    min-height: auto;
  }

  .nav {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .topbar {
    align-items: flex-start;
    flex-direction: column;
    padding: 18px 16px;
  }
}
</style>
