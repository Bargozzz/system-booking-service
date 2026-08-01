<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const isSuccess = computed(() => route.query.status === 'success')
const bookingId = route.query.bookingId
const seatNumber = route.query.seatNumber
const scheduleId = route.query.scheduleId

function backToSearch() {
  router.push('/search')
}

function tryAgain() {
  router.push({ name: 'SeatSelection', params: { scheduleId } })
}
</script>

<template>
  <div class="container narrow">
    <div v-if="isSuccess" class="result-card card success">
      <span class="result-icon">✓</span>
      <p class="eyebrow">Booking confirmed</p>
      <h1>You're all set</h1>
      <p class="detail">Seat <strong>{{ seatNumber }}</strong> on schedule #{{ scheduleId }} is booked.</p>
      <p class="booking-ref">Booking reference: <span>#{{ bookingId }}</span></p>
      <button class="btn btn-primary" @click="backToSearch">Book another trip</button>
    </div>

    <div v-else class="result-card card failure">
      <span class="result-icon">✕</span>
      <p class="eyebrow">Hold expired</p>
      <h1>We couldn't complete that booking</h1>
      <p class="detail">
        Your hold on seat <strong>{{ seatNumber }}</strong> expired before confirmation went through,
        or someone else booked it in the meantime.
      </p>
      <div class="actions">
        <button class="btn btn-secondary" @click="backToSearch">Back to search</button>
        <button class="btn btn-primary" @click="tryAgain">Pick another seat</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.narrow {
  max-width: 480px;
}

.result-card {
  margin-top: 3rem;
  padding: 2.5rem;
  text-align: center;
}

.result-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  font-size: 1.6rem;
  margin-bottom: 1rem;
}

.success .result-icon {
  background: rgba(79, 209, 197, 0.15);
  color: var(--color-teal);
}

.failure .result-icon {
  background: rgba(226, 87, 76, 0.15);
  color: var(--color-danger);
}

.detail {
  color: var(--color-text-dim);
  margin: 1rem 0;
  line-height: 1.5;
}

.booking-ref {
  font-family: var(--font-display);
  color: var(--color-amber);
  margin-bottom: 1.5rem;
}

.actions {
  display: flex;
  justify-content: center;
  gap: 1rem;
  margin-top: 1rem;
}
</style>
