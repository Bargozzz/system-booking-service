<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../api/client'

const route = useRoute()
const router = useRouter()

const seatId = Number(route.query.seatId)
const seatNumber = route.query.seatNumber
const scheduleId = route.query.scheduleId
const lockExpiresAt = Number(route.query.lockExpiresAt)

const now = ref(Date.now())
const confirming = ref(false)
const cancelling = ref(false)
const error = ref('')

let clockTimer = null

const secondsLeft = computed(() => Math.max(0, Math.floor((lockExpiresAt - now.value) / 1000)))
const isExpired = computed(() => secondsLeft.value <= 0)
const busy = computed(() => confirming.value || cancelling.value)

function formatCountdown(sec) {
  const m = Math.floor(sec / 60)
  const s = sec % 60
  return `${m}:${s.toString().padStart(2, '0')}`
}

async function handleConfirm() {
  confirming.value = true
  error.value = ''
  try {
    const { data } = await api.confirmBooking(seatId)
    router.push({
      name: 'Confirmation',
      query: { bookingId: data.booking_id, seatNumber, scheduleId, status: 'success' },
    })
  } catch (err) {
    if (err.response?.status === 409) {
      router.push({ name: 'Confirmation', query: { seatNumber, scheduleId, status: 'expired' } })
    } else if (err.response?.status === 401) {
      error.value = 'Sesi Anda berakhir. Silakan login kembali.'
    } else {
      error.value = 'Booking gagal dikonfirmasi. Silakan coba lagi.'
    }
  } finally {
    confirming.value = false
  }
}

// "Batal" here genuinely releases the hold via the API — not just a
// client-side navigation — so the seat becomes available to other users
// immediately instead of sitting reserved for up to 5 minutes for nothing.
async function handleCancel() {
  cancelling.value = true
  error.value = ''
  try {
    if (!isExpired.value) {
      await api.unlockSeat(seatId)
    }
  } catch (err) {
    // If it's already expired/released server-side, that's fine — the
    // outcome the user wants (seat freed up) is already true.
  } finally {
    cancelling.value = false
    router.push({ name: 'SeatSelection', params: { scheduleId } })
  }
}
</script>

<template>
  <div class="container narrow">
    <p class="eyebrow">Ringkasan booking</p>
    <h1>Periksa sebelum konfirmasi</h1>

    <div class="summary-card card">
      <div class="stub-row">
        <span class="stub-label">Jadwal</span>
        <span class="stub-value">#{{ scheduleId }}</span>
      </div>
      <div class="stub-divider"></div>
      <div class="stub-row">
        <span class="stub-label">Kursi</span>
        <span class="stub-value seat-value">{{ seatNumber }}</span>
      </div>
      <div class="stub-divider"></div>
      <div class="stub-row">
        <span class="stub-label">Waktu tahan tersisa</span>
        <span class="stub-value countdown" :class="{ danger: secondsLeft < 60 }">
          {{ isExpired ? 'Kedaluwarsa' : formatCountdown(secondsLeft) }}
        </span>
      </div>
    </div>

    <p v-if="error" class="error-banner" style="margin-top:1rem">{{ error }}</p>

    <div v-if="isExpired" class="error-banner" style="margin-top:1rem">
      Waktu tahan kursi Anda sudah habis. Silakan kembali dan pilih kursi lagi.
    </div>

    <div class="actions">
      <button class="btn btn-secondary" :disabled="busy" @click="handleCancel">
        {{ cancelling ? 'Membatalkan…' : '← Batal, pilih kursi lain' }}
      </button>
      <button class="btn btn-primary" :disabled="isExpired || busy" @click="handleConfirm">
        {{ confirming ? 'Mengonfirmasi…' : 'Konfirmasi booking' }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.narrow {
  max-width: 560px;
}

.summary-card {
  margin-top: 1.75rem;
  padding: 1.75rem;
  position: relative;
}

.summary-card::before,
.summary-card::after {
  content: '';
  position: absolute;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--color-bg);
  top: 50%;
  transform: translateY(-50%);
}

.summary-card::before { left: -10px; }
.summary-card::after { right: -10px; }

.stub-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.6rem 0;
}

.stub-label {
  color: var(--color-text-dim);
  font-size: 0.85rem;
}

.stub-value {
  font-family: var(--font-display);
  font-weight: 600;
}

.seat-value {
  font-size: 1.4rem;
  color: var(--color-amber);
}

.countdown {
  font-size: 1.1rem;
  color: var(--color-teal);
}

.countdown.danger {
  color: var(--color-danger);
}

.stub-divider {
  border-top: 1px dashed var(--color-border);
}

.actions {
  display: flex;
  justify-content: space-between;
  margin-top: 2rem;
  gap: 1rem;
  flex-wrap: wrap;
}
</style>
