<script setup>
import { ref } from "vue";
import { useRouter, useRoute } from "vue-router";
import { api } from "../api/client";
import { setSession } from "../store/auth";

const router = useRouter();
const route = useRoute();

const username = ref("alice");
const password = ref("password123");
const loading = ref(false);
const error = ref("");

async function handleLogin() {
  error.value = "";
  loading.value = true;
  try {
    const { data } = await api.login(username.value, password.value);
    setSession(username.value, data.user_id, data.access_token, data.refresh_token);
    router.push(route.query.redirect || "/search");
  } catch (err) {
    if (err.response?.status === 401) {
      error.value = "Username or password is incorrect.";
    } else {
      error.value = "Something went wrong logging in. Please try again.";
    }
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="container login-container">
    <div class="card login-card">
      <p class="eyebrow">Boarding pass access</p>
      <h1>Log in</h1>
      <p class="subtitle">Use one of the seeded demo accounts to lock and confirm seats.</p>

      <form @submit.prevent="handleLogin" class="login-form">
        <label class="field">
          <span class="field-label">Username</span>
          <input v-model="username" class="input" type="text" autocomplete="username" required />
        </label>
        <label class="field">
          <span class="field-label">Password</span>
          <input v-model="password" class="input" type="password" autocomplete="current-password" required />
        </label>

        <p v-if="error" class="error-banner">{{ error }}</p>

        <button class="btn btn-primary" type="submit" :disabled="loading">
          {{ loading ? "Signing in…" : "Log in" }}
        </button>
      </form>

      <p class="demo-hint">Demo accounts: <code>alice</code> / <code>bob</code>, password <code>password123</code></p>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  padding-top: 4rem;
}

.login-card {
  width: 100%;
  max-width: 400px;
  padding: 2rem;
}

.subtitle {
  color: var(--color-text-dim);
  font-size: 0.9rem;
  margin: 0.5rem 0 1.5rem;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.field-label {
  font-size: 0.8rem;
  color: var(--color-text-dim);
}

.demo-hint {
  margin-top: 1.5rem;
  font-size: 0.8rem;
  color: var(--color-text-dim);
  border-top: 1px dashed var(--color-border);
  padding-top: 1rem;
}

.demo-hint code {
  color: var(--color-teal);
}
</style>
