<template>
  <AppShell>
    <n-card>
      <template #header>
        <div class="card-head">
          <h2 class="section-title">Data order</h2>
          <n-button type="primary" @click="openCreate">Tambah</n-button>
        </div>
      </template>
      <n-data-table :columns="columns" :data="orders" :loading="loading" :pagination="{ pageSize: 8 }" />
    </n-card>

    <n-drawer v-model:show="showForm" width="460">
      <n-drawer-content :title="editingId ? 'Edit order' : 'Tambah order'">
        <n-form :model="form" :rules="rules" ref="formRef">
          <n-form-item label="Nomor order" path="order_number">
            <n-input v-model:value="form.order_number" placeholder="Kosongkan untuk auto-generate" />
          </n-form-item>
          <n-form-item label="Customer" path="customer_id">
            <n-select v-model:value="form.customer_id" :options="customerOptions" />
          </n-form-item>
          <n-form-item label="Status" path="status">
            <n-select v-model:value="form.status" :options="statusOptions" />
          </n-form-item>
          <n-form-item label="Total" path="total_amount">
            <n-input-number v-model:value="form.total_amount" :min="0" class="full" />
          </n-form-item>
          <n-form-item label="Catatan" path="notes">
            <n-input v-model:value="form.notes" type="textarea" />
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
import { NButton, NCard, NDataTable, NDrawer, NDrawerContent, NForm, NFormItem, NInput, NInputNumber, NSelect, useMessage } from 'naive-ui'
import { http } from '../api/http'
import AppShell from '../components/AppShell.vue'

const message = useMessage()
const customers = ref([])
const orders = ref([])
const loading = ref(false)
const saving = ref(false)
const showForm = ref(false)
const editingId = ref(null)
const formRef = ref(null)

const statusOptions = ['pending', 'processing', 'shipped', 'delivered', 'cancelled'].map((value) => ({ label: value, value }))
const customerOptions = computed(() => customers.value.map((customer) => ({ label: customer.name, value: customer.id })))
const form = reactive({ order_number: '', customer_id: null, status: 'pending', total_amount: 0, notes: '' })
const rules = {
  customer_id: [{ required: true, type: 'number', message: 'Customer wajib dipilih' }],
  status: [{ required: true, message: 'Status wajib diisi' }],
}

const columns = [
  { title: 'Nomor', key: 'order_number' },
  { title: 'Customer', key: 'customer.name' },
  { title: 'Status', key: 'status' },
  {
    title: 'Total',
    key: 'total_amount',
    render: (row) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(row.total_amount),
  },
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
  Object.assign(form, { order_number: '', customer_id: null, status: 'pending', total_amount: 0, notes: '' })
  editingId.value = null
}

function openCreate() {
  resetForm()
  showForm.value = true
}

function openEdit(row) {
  Object.assign(form, {
    order_number: row.order_number,
    customer_id: row.customer_id,
    status: row.status,
    total_amount: row.total_amount,
    notes: row.notes,
  })
  editingId.value = row.id
  showForm.value = true
}

async function load() {
  loading.value = true
  try {
    const [customerRes, orderRes] = await Promise.all([http.get('/customers'), http.get('/orders')])
    customers.value = customerRes.data
    orders.value = orderRes.data
  } finally {
    loading.value = false
  }
}

async function save() {
  await formRef.value?.validate()
  saving.value = true
  try {
    if (editingId.value) {
      await http.put(`/orders/${editingId.value}`, form)
    } else {
      await http.post('/orders', form)
    }
    message.success('Order tersimpan')
    showForm.value = false
    await load()
  } catch (error) {
    message.error(error.response?.data?.error || 'Gagal menyimpan order')
  } finally {
    saving.value = false
  }
}

async function remove(id) {
  await http.delete(`/orders/${id}`)
  message.success('Order dihapus')
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
