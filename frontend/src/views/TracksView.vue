<template>
  <AppShell>
    <n-card>
      <template #header>
        <div class="card-head">
          <h2 class="section-title">Order tracking</h2>
          <n-button type="primary" @click="openCreate">Tambah</n-button>
        </div>
      </template>
      <n-data-table :columns="columns" :data="tracks" :loading="loading" :pagination="{ pageSize: 8 }" />
    </n-card>

    <n-drawer v-model:show="showForm" width="460">
      <n-drawer-content :title="editingId ? 'Edit tracking' : 'Tambah tracking'">
        <n-form :model="form" :rules="rules" ref="formRef">
          <n-form-item label="Order" path="order_id">
            <n-select v-model:value="form.order_id" :options="orderOptions" />
          </n-form-item>
          <n-form-item label="Status" path="status">
            <n-select v-model:value="form.status" :options="statusOptions" />
          </n-form-item>
          <n-form-item label="Lokasi" path="location">
            <n-input v-model:value="form.location" />
          </n-form-item>
          <n-form-item label="Waktu tracking" path="tracked_at">
            <n-date-picker v-model:value="form.tracked_at" type="datetime" class="full" />
          </n-form-item>
          <n-form-item label="Deskripsi" path="description">
            <n-input v-model:value="form.description" type="textarea" />
          </n-form-item>
        </n-form>
        <template #footer>
          <n-button @click="showForm = false">Batal</n-button>
          <n-button type="primary" :loading="saving" @click="save">Simpan</n-button>
        </template>
      </n-drawer-content>
    </n-drawer>
  </AppShell>
</template>

<script setup>
import { computed, h, onMounted, reactive, ref } from 'vue'
import { NButton, NCard, NDataTable, NDatePicker, NDrawer, NDrawerContent, NForm, NFormItem, NInput, NSelect, useMessage } from 'naive-ui'
import { http } from '../api/http'
import AppShell from '../components/AppShell.vue'

const message = useMessage()
const orders = ref([])
const tracks = ref([])
const loading = ref(false)
const saving = ref(false)
const showForm = ref(false)
const editingId = ref(null)
const formRef = ref(null)

const statusOptions = ['pending', 'processing', 'packed', 'shipped', 'in_transit', 'delivered', 'cancelled'].map((value) => ({ label: value, value }))
const orderOptions = computed(() => orders.value.map((order) => ({ label: order.order_number, value: order.id })))
const form = reactive({ order_id: null, status: 'processing', location: '', tracked_at: Date.now(), description: '' })
const rules = {
  order_id: [{ required: true, type: 'number', message: 'Order wajib dipilih' }],
  status: [{ required: true, message: 'Status wajib diisi' }],
}

const columns = [
  { title: 'Order', key: 'order.order_number' },
  { title: 'Status', key: 'status' },
  { title: 'Lokasi', key: 'location' },
  {
    title: 'Waktu',
    key: 'tracked_at',
    render: (row) => new Date(row.tracked_at).toLocaleString('id-ID'),
  },
  { title: 'Deskripsi', key: 'description' },
  {
    title: '',
    key: 'actions',
    render(row) {
      return h('div', { class: 'table-actions' }, [
        h(NButton, { size: 'small', secondary: true, onClick: () => openEdit(row) }, { default: () => 'Edit' }),
        h(NButton, { size: 'small', type: 'error', secondary: true, onClick: () => remove(row.id) }, { default: () => 'Hapus' }),
      ])
    },
  },
]

function resetForm() {
  Object.assign(form, { order_id: null, status: 'processing', location: '', tracked_at: Date.now(), description: '' })
  editingId.value = null
}

function openCreate() {
  resetForm()
  showForm.value = true
}

function openEdit(row) {
  Object.assign(form, {
    order_id: row.order_id,
    status: row.status,
    location: row.location,
    tracked_at: new Date(row.tracked_at).getTime(),
    description: row.description,
  })
  editingId.value = row.id
  showForm.value = true
}

async function load() {
  loading.value = true
  try {
    const [orderRes, trackRes] = await Promise.all([http.get('/orders'), http.get('/order-tracks')])
    orders.value = orderRes.data
    tracks.value = trackRes.data
  } finally {
    loading.value = false
  }
}

async function save() {
  await formRef.value?.validate()
  saving.value = true
  const payload = { ...form, tracked_at: new Date(form.tracked_at).toISOString() }
  try {
    if (editingId.value) {
      await http.put(`/order-tracks/${editingId.value}`, payload)
    } else {
      await http.post('/order-tracks', payload)
    }
    message.success('Tracking tersimpan')
    showForm.value = false
    await load()
  } catch (error) {
    message.error(error.response?.data?.error || 'Gagal menyimpan tracking')
  } finally {
    saving.value = false
  }
}

async function remove(id) {
  await http.delete(`/order-tracks/${id}`)
  message.success('Tracking dihapus')
  await load()
}

onMounted(load)
</script>

<style scoped>
.card-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.full {
  width: 100%;
}
</style>
