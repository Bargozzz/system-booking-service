<script setup>
import { useRouter } from 'vue-router'
import { authState, clearSession } from './store/auth'

const router = useRouter()

function handleLogout() {
  clearSession()
  router.push('/search')
}
</script>

<template>
  <div class="app-shell">
    <header class="topbar">
      <RouterLink to="/search" class="brand">
        <span class="brand-mark">✈</span>
        <span class="brand-text">MINI<span class="brand-accent">BOOKING</span></span>
      </RouterLink>
      <div class="topbar-right">
        <span v-if="authState.isAuthenticated" class="user-chip">
          {{ authState.username }}
          <button class="link-btn" @click="handleLogout">log out</button>
        </span>
        <RouterLink v-else to="/login" class="btn btn-secondary login-link">Log in</RouterLink>
      </div>
    </header>
    <main>
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
.app-shell {
  min-height: 100%;
  display: flex;
  flex-direction: column;
}

.topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.1rem 1.75rem;
  border-bottom: 1px solid var(--color-border);
  background: rgba(15, 20, 32, 0.85);
  backdrop-filter: blur(6px);
  position: sticky;
  top: 0;
  z-index: 10;
}

.brand {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  text-decoration: none;
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 1.05rem;
  letter-spacing: 0.06em;
  color: var(--color-text);
}

.brand-mark {
  color: var(--color-amber);
  transform: rotate(45deg);
  display: inline-block;
}

.brand-accent {
  color: var(--color-amber);
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-chip {
  font-family: var(--font-display);
  font-size: 0.85rem;
  color: var(--color-text-dim);
  display: flex;
  align-items: center;
  gap: 0.7rem;
}

.link-btn {
  background: none;
  border: none;
  color: var(--color-teal);
  font-family: var(--font-display);
  font-size: 0.8rem;
  padding: 0;
  text-decoration: underline;
  text-underline-offset: 3px;
}

.login-link {
  font-family: var(--font-display);
  font-size: 0.85rem;
  padding: 0.5rem 1.1rem;
}
</style>
