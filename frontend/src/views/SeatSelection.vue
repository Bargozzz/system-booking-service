<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api/client'
import { authState } from '../store/auth'

const props = defineProps({ scheduleId: { type: [String, Number], required: true } })
const router = useRouter()

const seats = ref([])
const loading = ref(true)
const error = ref('')
const locking = ref(false)
const lockedSeat = ref(null) // the seat this browser session just locked (pending confirmation)
const lockExpiresAt = ref(null)
const now = ref(Date.now())

let pollTimer = null
let clockTimer = null

async function loadSeats() {
  try {
    const { data } = await api.viewSeats(props.scheduleId)
    seats.value = data
  } catch (err) {
    if (err.response?.status === 404) {
      error.value = 'Jadwal ini tidak ditemukan.'
    } else {
      error.value = 'Peta kursi gagal dimuat. Silakan coba lagi.'
    }
  } finally {
    loading.value = false
  }
}

// A seat's real state, cross-referenced against who is logged in — this is
// what lets us tell "booked by me" apart from "booked by someone else",
// since the backend keeps locked_by populated as the booking owner even
// after status flips to 'booked'.
function seatState(seat) {
  const isMine = authState.userId != null && seat.locked_by === authState.userId
  if (seat.status === 'booked') return isMine ? 'booked-mine' : 'booked-other'
  if (seat.status === 'locked') return isMine ? 'holding' : 'locked-other'
  return 'available'
}

async function handleSeatClick(seat) {
  if (seatState(seat) !== 'available' || locking.value) return
  locking.value = true
  error.value = ''
  try {
    const { data } = await api.lockSeat(seat.id)
    lockedSeat.value = seat
    lockExpiresAt.value = new Date(data.lock_expires_at).getTime()
    await loadSeats()
  } catch (err) {
    if (err.response?.status === 409) {
      error.value = err.response.data?.error || 'Kursi ini baru saja diambil orang lain. Silakan pilih kursi lain.'
      await loadSeats()
    } else if (err.response?.status === 401) {
      error.value = 'Sesi Anda berakhir. Silakan login kembali.'
    } else {
      error.value = 'Kursi gagal dikunci. Silakan coba lagi.'
    }
  } finally {
    locking.value = false
  }
}

function goToSummary() {
  router.push({
    name: 'BookingSummary',
    query: {
      seatId: lockedSeat.value.id,
      seatNumber: lockedSeat.value.seat_number,
      scheduleId: props.scheduleId,
      lockExpiresAt: lockExpiresAt.value,
    },
  })
}

const cancelling = ref(false)

async function handleCancel() {
  if (!lockedSeat.value || cancelling.value) return
  cancelling.value = true
  error.value = ''
  try {
    await api.unlockSeat(lockedSeat.value.id)
  } catch (err) {
    // Even if the unlock call fails (e.g. it already expired on its own),
    // we still want to clear the local hold state below — there's nothing
    // useful the user can do about a stale hold that's already gone.
  } finally {
    lockedSeat.value = null
    lockExpiresAt.value = null
    cancelling.value = false
    await loadSeats()
  }
}

const secondsLeft = computed(() => {
  if (!lockExpiresAt.value) return 0
  return Math.max(0, Math.floor((lockExpiresAt.value - now.value) / 1000))
})

const lockExpired = computed(() => lockedSeat.value && secondsLeft.value <= 0)

function formatCountdown(sec) {
  const m = Math.floor(sec / 60)
  const s = sec % 60
  return `${m}:${s.toString().padStart(2, '0')}`
}

// Vertical cabin layout: rows run front-to-back (top-to-bottom), each row
// holds a left-side seat and a right-side seat with an aisle between them —
// mirrors how the cabin actually looks from above, so "kiri"/"kanan" reads
// intuitively instead of an arbitrary horizontal strip.
const cabinRows = computed(() => {
  const byNumber = {}
  for (const seat of seats.value) {
    const letter = seat.seat_number.charAt(0)
    const number = seat.seat_number.slice(1)
    if (!byNumber[number]) byNumber[number] = {}
    byNumber[number][letter] = seat
  }
  return Object.keys(byNumber)
    .sort((a, b) => Number(a) - Number(b))
    .map((number) => ({ number, left: byNumber[number].A, right: byNumber[number].B }))
})

onMounted(() => {
  loadSeats()
  pollTimer = setInterval(loadSeats, 5000)
  clockTimer = setInterval(() => { now.value = Date.now() }, 1000)
})

onUnmounted(() => {
  clearInterval(pollTimer)
  clearInterval(clockTimer)
})
</script>

<template>
  <div class="container">
    <p class="eyebrow">Jadwal #{{ scheduleId }}</p>
    <h1>Pilih kursi Anda</h1>

    <p v-if="error" class="error-banner" style="margin-top:1rem">{{ error }}</p>

    <div v-if="loading" class="seat-loading pulse">Memuat peta kursi…</div>

    <template v-else>
      <div class="legend">
        <span class="legend-item"><span class="dot available"></span>Kosong</span>
        <span class="legend-item"><span class="dot locked-other"></span>Sedang dikunci orang lain</span>
        <span class="legend-item"><span class="dot booked-other"></span>Sudah dipesan orang lain</span>
        <span class="legend-item"><span class="dot booked-mine"></span>Sudah dipesan oleh Anda</span>
        <span class="legend-item"><span class="dot holding"></span>Sedang Anda tahan</span>
      </div>

      <div class="cabin">
        <div class="cabin-nose">
          <span class="nose-label">Depan kabin</span>
        </div>

        <div class="cabin-body">
          <div class="col-heading">
            <span>Kiri (A)</span>
            <span class="aisle-heading">Lorong</span>
            <span>Kanan (B)</span>
          </div>

          <div v-for="row in cabinRows" :key="row.number" class="cabin-row">
            <button
              v-if="row.left"
              class="seat"
              :class="seatState(row.left)"
              :disabled="seatState(row.left) !== 'available' || locking"
              @click="handleSeatClick(row.left)"
            >
              {{ row.left.seat_number }}
            </button>
            <div v-else class="seat-empty"></div>

            <span class="row-number">{{ row.number }}</span>

            <button
              v-if="row.right"
              class="seat"
              :class="seatState(row.right)"
              :disabled="seatState(row.right) !== 'available' || locking"
              @click="handleSeatClick(row.right)"
            >
              {{ row.right.seat_number }}
            </button>
            <div v-else class="seat-empty"></div>
          </div>
        </div>
      </div>

      <div v-if="lockedSeat && !lockExpired" class="lock-banner card">
        <div>
          <p class="eyebrow">Kursi {{ lockedSeat.seat_number }} ditahan untuk Anda</p>
          <p class="countdown">{{ formatCountdown(secondsLeft) }} tersisa untuk konfirmasi</p>
        </div>
        <div class="lock-banner-actions">
          <button class="btn btn-secondary" :disabled="cancelling" @click="handleCancel">
            {{ cancelling ? 'Membatalkan…' : 'Batal' }}
          </button>
          <button class="btn btn-primary" @click="goToSummary">Lanjut &amp; konfirmasi →</button>
        </div>
      </div>
      <div v-else-if="lockedSeat && lockExpired" class="error-banner">
        Waktu tahan kursi {{ lockedSeat.seat_number }} sudah habis. Silakan pilih kursi lagi.
      </div>
    </template>
  </div>
</template>

<style scoped>
.seat-loading {
  margin-top: 2rem;
  color: var(--color-text-dim);
}

.legend {
  display: flex;
  gap: 1.25rem;
  margin: 1.5rem 0 1rem;
  flex-wrap: wrap;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.8rem;
  color: var(--color-text-dim);
}

.dot {
  width: 10px;
  height: 10px;
  border-radius: 3px;
  display: inline-block;
}

.dot.available { background: var(--color-surface-raised); border: 1px solid var(--color-teal); }
.dot.locked-other { background: rgba(169, 118, 44, 0.4); }
.dot.booked-other { background: var(--color-border); }
.dot.booked-mine { background: #2fae6b; }
.dot.holding { background: var(--color-amber); }

/* Cabin: a tall, narrow silhouette (rounded top = front of vehicle) instead
   of a generic flat card, so the vertical seat plan reads as an actual
   cabin viewed from above rather than an arbitrary grid. */
.cabin {
  max-width: 340px;
  margin: 0 auto;
}

.cabin-nose {
  height: 40px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-bottom: none;
  border-radius: 60px 60px 0 0;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding-bottom: 0.4rem;
}

.nose-label {
  font-family: var(--font-display);
  font-size: 0.68rem;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--color-text-dim);
}

.cabin-body {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 0 0 16px 16px;
  padding: 1.25rem 1.5rem 1.75rem;
  display: flex;
  flex-direction: column;
  gap: 0.65rem;
}

.col-heading {
  display: grid;
  grid-template-columns: 56px 40px 56px;
  justify-content: center;
  column-gap: 1.5rem;
  font-size: 0.7rem;
  color: var(--color-text-dim);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  margin-bottom: 0.4rem;
  text-align: center;
}

.aisle-heading {
  opacity: 0.5;
}

.cabin-row {
  display: grid;
  grid-template-columns: 56px 40px 56px;
  align-items: center;
  justify-content: center;
  column-gap: 1.5rem;
}

.row-number {
  text-align: center;
  font-family: var(--font-display);
  font-size: 0.75rem;
  color: var(--color-text-dim);
}

.seat, .seat-empty {
  width: 56px;
  height: 48px;
}

.seat {
  border-radius: 8px;
  font-family: var(--font-display);
  font-size: 0.8rem;
  font-weight: 600;
  border: 1px solid var(--color-border);
  background: var(--color-surface-raised);
  color: var(--color-text-dim);
  transition: all 0.15s ease;
}

.seat.available {
  color: var(--color-text);
  border-color: var(--color-teal);
}

.seat.available:hover:not(:disabled) {
  background: rgba(79, 209, 197, 0.15);
  transform: scale(1.04);
}

.seat.locked-other {
  background: rgba(169, 118, 44, 0.18);
  color: var(--color-amber-dim);
  cursor: not-allowed;
}

.seat.holding {
  background: var(--color-amber);
  color: #1a1305;
  border-color: var(--color-amber);
}

.seat.booked-other {
  background: var(--color-border);
  color: #5a6180;
  cursor: not-allowed;
  text-decoration: line-through;
}

/* Distinct from booked-other on purpose: this is YOUR confirmed seat. */
.seat.booked-mine {
  background: rgba(47, 174, 107, 0.18);
  border-color: #2fae6b;
  color: #2fae6b;
  cursor: not-allowed;
  font-weight: 700;
}

.lock-banner {
  margin-top: 1.5rem;
  padding: 1.25rem 1.5rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-color: var(--color-amber);
  flex-wrap: wrap;
  gap: 1rem;
  max-width: 500px;
  margin-left: auto;
  margin-right: auto;
}

.countdown {
  font-family: var(--font-display);
  font-size: 1.3rem;
  margin: 0.2rem 0 0;
}

.lock-banner-actions {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}
</style>
