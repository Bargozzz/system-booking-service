<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api/client'

const router = useRouter()

const origin = ref('')
const destination = ref('')
const date = ref('')
const schedules = ref([])
const loading = ref(false)
const error = ref('')
const hasSearched = ref(false)

async function search() {
  loading.value = true
  error.value = ''
  hasSearched.value = true
  try {
    const params = {}
    if (origin.value) params.origin = origin.value
    if (destination.value) params.destination = destination.value
    if (date.value) params.date = date.value
    const { data } = await api.searchSchedules(params)
    schedules.value = data
  } catch (err) {
    if (err.response?.status === 429) {
      error.value = 'Too many searches at once — please wait a moment and try again.'
    } else {
      error.value = 'Could not load schedules right now. Please try again.'
    }
  } finally {
    loading.value = false
  }
}

function formatTime(iso) {
  return new Date(iso).toLocaleString(undefined, {
    weekday: 'short', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit',
  })
}

function formatPrice(price) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(price)
}

function selectSchedule(schedule) {
  router.push({ name: 'SeatSelection', params: { scheduleId: schedule.id } })
}

onMounted(search)
</script>

<template>
  <div class="container">
    <p class="eyebrow">Departures</p>
    <h1>Find a schedule</h1>

    <form @submit.prevent="search" class="search-form card">
      <label class="field">
        <span class="field-label">Origin</span>
        <input v-model="origin" class="input" type="text" placeholder="e.g. Jakarta" />
      </label>
      <label class="field">
        <span class="field-label">Destination</span>
        <input v-model="destination" class="input" type="text" placeholder="e.g. Bandung" />
      </label>
      <label class="field">
        <span class="field-label">Date</span>
        <input v-model="date" class="input" type="date" />
      </label>
      <button class="btn btn-primary search-btn" type="submit" :disabled="loading">
        {{ loading ? 'Searching…' : 'Search' }}
      </button>
    </form>

    <p v-if="error" class="error-banner" style="margin-top:1.5rem">{{ error }}</p>

    <div class="board" v-if="!loading">
      <div v-if="schedules.length === 0 && hasSearched" class="empty-state card">
        No schedules match that search. Try a different route or date.
      </div>

      <button
        v-for="schedule in schedules"
        :key="schedule.id"
        class="schedule-row card"
        @click="selectSchedule(schedule)"
      >
        <div class="route-col">
          <span class="route-cities">{{ schedule.origin }} <span class="arrow">→</span> {{ schedule.destination }}</span>
          <span class="operator">{{ schedule.operator_name }}</span>
        </div>
        <div class="time-col">
          <span class="depart-label">Departs</span>
          <span class="depart-time">{{ formatTime(schedule.departure_time) }}</span>
        </div>
        <div class="price-col">
          <span class="price">{{ formatPrice(schedule.price) }}</span>
          <span class="select-cta">Select seats →</span>
        </div>
      </button>
    </div>

    <div v-else class="board">
      <div class="schedule-row card pulse" v-for="n in 3" :key="n" style="height:76px"></div>
    </div>
  </div>
</template>

<style scoped>
.search-form {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr auto;
  gap: 1rem;
  padding: 1.5rem;
  margin-top: 1.5rem;
  align-items: end;
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

.search-btn {
  height: fit-content;
}

.board {
  margin-top: 2rem;
  display: flex;
  flex-direction: column;
  gap: 0.9rem;
}

.empty-state {
  padding: 2rem;
  text-align: center;
  color: var(--color-text-dim);
}

.schedule-row {
  display: grid;
  grid-template-columns: 1.4fr 1fr 1fr;
  align-items: center;
  padding: 1.1rem 1.5rem;
  text-align: left;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  color: var(--color-text);
  transition: border-color 0.15s ease, transform 0.12s ease;
}

.schedule-row:hover {
  border-color: var(--color-amber);
  transform: translateY(-1px);
}

.route-col {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.route-cities {
  font-family: var(--font-display);
  font-weight: 600;
  font-size: 1rem;
}

.arrow {
  color: var(--color-amber);
}

.operator {
  font-size: 0.8rem;
  color: var(--color-text-dim);
}

.time-col {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.depart-label {
  font-size: 0.72rem;
  color: var(--color-text-dim);
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.depart-time {
  font-family: var(--font-display);
}

.price-col {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 0.3rem;
}

.price {
  font-family: var(--font-display);
  font-weight: 600;
  color: var(--color-teal);
}

.select-cta {
  font-size: 0.78rem;
  color: var(--color-amber);
}

@media (max-width: 720px) {
  .search-form {
    grid-template-columns: 1fr;
  }
  .schedule-row {
    grid-template-columns: 1fr;
    gap: 0.5rem;
  }
  .price-col {
    align-items: flex-start;
  }
}
</style>
