<template>
  <main class="login-page">
    <section class="login-panel">
      <div>
        <h1>OrderOps</h1>
        <p>Masuk dengan password dan kode 6 digit dari authenticator app.</p>
      </div>

      <n-alert type="info" :show-icon="false">
        Development admin: <strong>admin@example.com</strong> / <strong>admin12345</strong>.
        Secret TOTP default: <strong>JBSWY3DPEHPK3PXP</strong>.
      </n-alert>

      <n-form ref="formRef" :model="form" :rules="rules" @submit.prevent="submit">
        <n-form-item label="Email" path="email">
          <n-input v-model:value="form.email" placeholder="admin@example.com" />
        </n-form-item>
        <n-form-item label="Password" path="password">
          <n-input v-model:value="form.password" type="password" show-password-on="click" />
        </n-form-item>
        <n-form-item label="Authenticator code" path="otp">
          <n-input v-model:value="form.otp" maxlength="6" placeholder="123456" />
        </n-form-item>
        <n-button :disabled="loading" attr-type="submit" type="primary" block :loading="loading">Login</n-button>
      </n-form>
    </section>
  </main>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { NAlert, NButton, NForm, NFormItem, NInput, useMessage } from 'naive-ui'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const message = useMessage()
const auth = useAuthStore()
const formRef = ref(null)
const loading = ref(false)

const form = reactive({
  email: 'admin@example.com',
  password: 'admin12345',
  otp: '',
})

const rules = {
  email: [{ required: true, message: 'Email wajib diisi' }],
  password: [{ required: true, message: 'Password wajib diisi' }],
  otp: [{ required: true, len: 6, message: 'Kode OTP harus 6 digit' }],
}

async function submit() {
  await formRef.value?.validate()
  loading.value = true
  try {
    await auth.login(form)
    message.success('Login berhasil')
    router.push('/')
  } catch (error) {
    await new Promise((resolve) => setTimeout(resolve, 2000)) // Simulate delay
    message.error(error.response?.data?.error || 'Login gagal')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  display: grid;
  min-height: 100vh;
  place-items: center;
  padding: 24px;
  background:
    linear-gradient(rgba(16, 24, 40, 0.72), rgba(16, 24, 40, 0.72)),
    url("https://images.unsplash.com/photo-1553413077-190dd305871c?auto=format&fit=crop&w=1600&q=80");
  background-position: center;
  background-size: cover;
}

.login-panel {
  display: grid;
  width: min(100%, 430px);
  gap: 18px;
  padding: 28px;
  border-radius: 8px;
  background: #fff;
  box-shadow: 0 24px 80px rgba(16, 24, 40, 0.25);
}

h1 {
  margin: 0;
  font-size: 32px;
}

p {
  margin: 8px 0 0;
  color: #667085;
}
</style>
