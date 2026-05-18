<template>
  <AppShell>
    <n-card>
      <template #header>
        <div class="card-head">
          <h2 class="section-title">Data customer</h2>
          <n-button type="primary" @click="openCreate">Tambah</n-button>
        </div>
      </template>
      <n-data-table :columns="columns" :data="customers" :loading="loading" :pagination="{ pageSize: 8 }" />
    </n-card>

    <n-drawer v-model:show="showForm" width="420">
      <n-drawer-content :title="editingId ? 'Edit customer' : 'Tambah customer'">
        <n-form :model="form" :rules="rules" ref="formRef">
          <n-form-item label="Nama" path="name">
            <n-input v-model:value="form.name" />
          </n-form-item>
          <n-form-item label="Email" path="email">
            <n-input v-model:value="form.email" />
          </n-form-item>
          <n-form-item label="Telepon" path="phone">
            <n-input v-model:value="form.phone" />
          </n-form-item>
          <n-form-item label="Alamat" path="address">
            <n-input v-model:value="form.address" type="textarea" />
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
import { h, onMounted, reactive, ref } from 'vue'
import { NButton, NCard, NDataTable, NDrawer, NDrawerContent, NForm, NFormItem, NInput, useMessage } from 'naive-ui'
import { http } from '../api/http'
import AppShell from '../components/AppShell.vue'

const message = useMessage()
const customers = ref([])
const loading = ref(false)
const saving = ref(false)
const showForm = ref(false)
const editingId = ref(null)
const formRef = ref(null)

const form = reactive({ name: '', email: '', phone: '', address: '' })
const rules = {
  name: [{ required: true, message: 'Nama wajib diisi' }],
  email: [{ required: true, type: 'email', message: 'Email valid wajib diisi' }],
}

const columns = [
  { title: 'Nama', key: 'name' },
  { title: 'Email', key: 'email' },
  { title: 'Telepon', key: 'phone' },
  { title: 'Alamat', key: 'address' },
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
  Object.assign(form, { name: '', email: '', phone: '', address: '' })
  editingId.value = null
}

function openCreate() {
  resetForm()
  showForm.value = true
}

function openEdit(row) {
  Object.assign(form, row)
  editingId.value = row.id
  showForm.value = true
}

async function load() {
  loading.value = true
  try {
    customers.value = (await http.get('/customers')).data
  } finally {
    loading.value = false
  }
}

async function save() {
  await formRef.value?.validate()
  saving.value = true
  try {
    if (editingId.value) {
      await http.put(`/customers/${editingId.value}`, form)
    } else {
      await http.post('/customers', form)
    }
    message.success('Customer tersimpan')
    showForm.value = false
    await load()
  } catch (error) {
    message.error(error.response?.data?.error || 'Gagal menyimpan customer')
  } finally {
    saving.value = false
  }
}

async function remove(id) {
  await http.delete(`/customers/${id}`)
  message.success('Customer dihapus')
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
</style>
