<template>
  <div class="login-page">
    <div class="login-card glass-card">
      <div class="header">
        <span class="logo">🚌</span>
        <h2>Login MiniBooking</h2>
        <p>Masuk ke akun untuk mulai mengunci dan memesan kursi bus</p>
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label>Email Address</label>
          <input 
            v-model="email" 
            type="email" 
            class="input-field"
            placeholder="user1@example.com" 
            required 
          />
        </div>

        <div class="form-group">
          <label>Password</label>
          <input 
            v-model="password" 
            type="password" 
            class="input-field"
            placeholder="••••••••" 
            required 
          />
        </div>

        <div v-if="error" class="error-banner">
          {{ error }}
        </div>

        <button type="submit" class="btn-primary btn-block" :disabled="loading">
          <span v-if="loading">Loading...</span>
          <span v-else>Login Sekarang</span>
        </button>

        <div class="quick-logins">
          <span>Akun Demo Uji Coba:</span>
          <div class="demo-buttons">
            <button 
              type="button"
              v-for="u in demoUsers" 
              :key="u.email"
              @click="fillDemo(u.email)"
              class="btn-demo-pill"
            >
              {{ u.name }} ({{ u.email }})
            </button>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
const email = ref('user1@example.com')
const password = ref('password')
const loading = ref(false)
const error = ref<string | null>(null)

const authStore = useAuthStore()
const { fetchApi } = useApi()
const router = useRouter()

const demoUsers = [
  { name: 'User Satu', email: 'user1@example.com' },
  { name: 'User Dua', email: 'user2@example.com' },
  { name: 'User Tiga', email: 'user3@example.com' },
]

const fillDemo = (e: string) => {
  email.value = e
  password.value = 'password'
}

const handleLogin = async () => {
  loading.value = true
  error.value = null

  const { data, error: err } = await fetchApi('/login', {
    method: 'POST',
    body: { email: email.value, password: password.value },
  })

  loading.value = false

  if (err) {
    error.value = err.message || 'Email atau password salah'
  } else if (data && data.token) {
    authStore.setAuth(data.token, data.user)
    router.push('/')
  }
}
</script>

<style scoped>
.login-page {
  min-height: calc(100vh - 80px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.login-card {
  max-width: 420px;
  width: 100%;
}

.header {
  text-align: center;
  margin-bottom: 24px;
}

.logo {
  font-size: 2.5rem;
}

.header h2 {
  font-size: 1.5rem;
  font-weight: 800;
  margin-top: 8px;
}

.header p {
  font-size: 0.85rem;
  color: var(--color-text-muted);
  margin-top: 4px;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 0.85rem;
  font-weight: 600;
}

.btn-block {
  width: 100%;
  margin-top: 8px;
}

.error-banner {
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid rgba(239, 68, 68, 0.4);
  color: #fca5a5;
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 0.85rem;
  text-align: center;
}

.quick-logins {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--color-border);
  font-size: 0.8rem;
  color: var(--color-text-muted);
}

.demo-buttons {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-top: 8px;
}

.btn-demo-pill {
  background: rgba(56, 189, 248, 0.08);
  border: 1px solid rgba(56, 189, 248, 0.2);
  color: var(--color-primary);
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 0.78rem;
  cursor: pointer;
  text-align: left;
  transition: all 0.2s ease;
}

.btn-demo-pill:hover {
  background: rgba(56, 189, 248, 0.2);
}
</style>
