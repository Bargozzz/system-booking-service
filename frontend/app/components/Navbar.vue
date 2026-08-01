<template>
  <header class="navbar">
    <div class="navbar-container">
      <NuxtLink to="/" class="brand">
        <span class="brand-icon">🚌</span>
        <span class="brand-name">MiniBooking</span>
      </NuxtLink>

      <div class="user-section">
        <template v-if="authStore.isAuthenticated">
          <div class="user-info">
            <span class="user-icon">👤</span>
            <div class="user-details">
              <span class="user-name">{{ authStore.user?.name }}</span>
              <span class="user-email">{{ authStore.user?.email }}</span>
            </div>
          </div>
          <button @click="handleLogout" class="btn-logout">Logout</button>
        </template>

        <template v-else>
          <div class="demo-logins">
            <span class="demo-label">Demo Login:</span>
            <button 
              v-for="email in ['user1@example.com', 'user2@example.com', 'user3@example.com']" 
              :key="email"
              @click="quickLogin(email)" 
              class="btn-demo"
              :disabled="authStore.loading"
            >
              {{ email.split('@')[0] }}
            </button>
          </div>
          <NuxtLink to="/login" class="btn-login">Login</NuxtLink>
        </template>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
const authStore = useAuthStore()
const { fetchApi } = useApi()

onMounted(() => {
  authStore.restoreSession()
})

const quickLogin = async (email: string) => {
  authStore.loading = true
  authStore.error = null
  const { data, error } = await fetchApi('/login', {
    method: 'POST',
    body: { email, password: 'password' },
  })

  authStore.loading = false
  if (error) {
    alert(error.message || 'Login gagal')
  } else if (data && data.token) {
    authStore.setAuth(data.token, data.user)
  }
}

const handleLogout = async () => {
  await fetchApi('/logout', { method: 'POST' })
  authStore.logout()
  useRouter().push('/')
}
</script>

<style scoped>
.navbar {
  background: rgba(15, 23, 42, 0.85);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--color-border);
  position: sticky;
  top: 0;
  z-index: 100;
}

.navbar-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 14px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
  color: var(--color-text-main);
  font-size: 1.25rem;
  font-weight: 800;
  letter-spacing: -0.5px;
}

.brand-icon {
  font-size: 1.5rem;
}

.user-section {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  background: rgba(30, 41, 59, 0.8);
  padding: 6px 14px;
  border-radius: 12px;
  border: 1px solid var(--color-border);
}

.user-details {
  display: flex;
  flex-direction: column;
}

.user-name {
  font-size: 0.85rem;
  font-weight: 700;
  color: var(--color-text-main);
}

.user-email {
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.demo-logins {
  display: flex;
  align-items: center;
  gap: 8px;
}

.demo-label {
  font-size: 0.8rem;
  color: var(--color-text-muted);
}

.btn-demo {
  background: rgba(56, 189, 248, 0.1);
  border: 1px solid rgba(56, 189, 248, 0.3);
  color: var(--color-primary);
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 0.75rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-demo:hover {
  background: rgba(56, 189, 248, 0.25);
}

.btn-login, .btn-logout {
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  text-decoration: none;
  transition: all 0.2s ease;
}

.btn-login {
  background: var(--color-primary);
  color: #0f172a;
}

.btn-logout {
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #f87171;
}

.btn-logout:hover {
  background: rgba(239, 68, 68, 0.3);
}
</style>
