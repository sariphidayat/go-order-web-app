<template>
  <AppShell>
    <div class="stats">
      <n-card v-for="item in stats" :key="item.label">
        <span>{{ item.label }}</span>
        <strong>{{ item.value }}</strong>
      </n-card>
    </div>

    <n-card title="Order terbaru">
      <n-data-table :columns="columns" :data="orders" :loading="loading" :pagination="{ pageSize: 6 }" />
    </n-card>
  </AppShell>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { NCard, NDataTable } from 'naive-ui'
import { http } from '../api/http'
import AppShell from '../components/AppShell.vue'

const customers = ref([])
const orders = ref([])
const tracks = ref([])
const loading = ref(false)

const stats = computed(() => [
  { label: 'Customers', value: customers.value.length },
  { label: 'Orders', value: orders.value.length },
  { label: 'Tracking events', value: tracks.value.length },
])

const columns = [
  { title: 'Order', key: 'order_number' },
  { title: 'Customer', key: 'customer.name' },
  { title: 'Status', key: 'status' },
  {
    title: 'Total',
    key: 'total_amount',
    render: (row) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(row.total_amount),
  },
]

onMounted(async () => {
  loading.value = true
  try {
    const [customerRes, orderRes, trackRes] = await Promise.all([
      http.get('/customers'),
      http.get('/orders'),
      http.get('/order-tracks'),
    ])
    customers.value = customerRes.data
    orders.value = orderRes.data
    tracks.value = trackRes.data
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
  margin-bottom: 18px;
}

.stats span,
.stats strong {
  display: block;
}

.stats span {
  color: #667085;
}

.stats strong {
  margin-top: 8px;
  font-size: 30px;
}

@media (max-width: 720px) {
  .stats {
    grid-template-columns: 1fr;
  }
}
</style>
