<template>
  <div class="seat-selection-page">
    <!-- Header Back Navigation -->
    <div class="top-bar">
      <NuxtLink to="/" class="back-link">
        ← Kembali ke Cari Jadwal
      </NuxtLink>
      <div v-if="schedule" class="schedule-summary">
        <span class="operator-tag">{{ schedule.operator }}</span>
        <span class="route-tag">{{ schedule.origin }} ➔ {{ schedule.destination }}</span>
        <span class="date-tag">📅 {{ schedule.departure_date }} ({{ schedule.departure_time.substring(0, 5) }})</span>
      </div>
    </div>

    <!-- Live Lock Banner (If active locked seat) -->
    <div v-if="bookingStore.activeLockedSeat" class="lock-banner-wrapper">
      <CountdownTimer :seconds="bookingStore.lockRemainingSeconds" />
    </div>

    <div class="selection-grid">
      <!-- Interactive Bus Seat Grid -->
      <div class="grid-column">
        <SeatGrid 
          :seats="bookingStore.seats" 
          :selected-seat-id="selectedSeat?.id || null"
          @select-seat="handleSelectSeat"
        />
      </div>

      <!-- Action & Confirmation Sidebar -->
      <div class="sidebar-column">
        <div class="booking-card glass-card">
          <h3>Ringkasan Pemesanan</h3>

          <div v-if="selectedSeat" class="selected-seat-box">
            <div class="box-header">
              <span>Nomor Kursi</span>
              <span class="seat-pill">{{ selectedSeat.seat_number }}</span>
            </div>

            <div class="detail-row">
              <span>Harga Tiket</span>
              <span class="price">Rp {{ schedule?.price.toLocaleString('id-ID') }}</span>
            </div>

            <div class="detail-row">
              <span>Status Kunci</span>
              <span v-if="isLockedByMe" class="badge badge-warning">Terkunci (5 Menit)</span>
              <span v-else class="badge badge-success">Siap Dikunci</span>
            </div>
          </div>

          <div v-else class="no-seat-box">
            <span class="icon">💺</span>
            <p>Silakan klik salah satu kursi hijau yang tersedia pada layout bus di samping.</p>
          </div>

          <!-- Message Banner -->
          <div v-if="message" :class="['msg-banner', messageType]">
            {{ message }}
          </div>

          <!-- Step 1: Lock Seat Button -->
          <button 
            v-if="selectedSeat && !isLockedByMe"
            @click="lockSelectedSeat" 
            class="btn-primary btn-block"
            :disabled="actionLoading"
          >
            <span v-if="actionLoading">Mengunci...</span>
            <span v-else>🔒 Kunci Kursi (5 Menit)</span>
          </button>

          <!-- Step 2: Confirm Booking Button -->
          <button 
            v-if="selectedSeat && isLockedByMe"
            @click="confirmSelectedBooking" 
            class="btn-primary btn-block btn-confirm"
            :disabled="actionLoading"
          >
            <span v-if="actionLoading">Memproses...</span>
            <span v-else>✅ Konfirmasi & Bayar Pesanan</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Booking Success Modal -->
    <div v-if="confirmedBooking" class="modal-overlay">
      <div class="modal-card glass-card">
        <div class="success-icon">🎉</div>
        <h2>Pemesanan Tiket Berhasil!</h2>
        <p class="modal-sub">Kode Booking Resmi Anda:</p>

        <div class="ticket-code">
          {{ confirmedBooking.booking_code }}
        </div>

        <div class="ticket-details">
          <div class="t-row"><span>Operator:</span> <strong>{{ confirmedBooking.operator }}</strong></div>
          <div class="t-row"><span>Rute:</span> <strong>{{ confirmedBooking.origin }} ➔ {{ confirmedBooking.destination }}</strong></div>
          <div class="t-row"><span>Kursi:</span> <strong>{{ confirmedBooking.seat_number }}</strong></div>
          <div class="t-row"><span>Waktu:</span> <strong>{{ confirmedBooking.departure_date }} @ {{ confirmedBooking.departure_time }}</strong></div>
          <div class="t-row"><span>Total:</span> <strong>Rp {{ confirmedBooking.price.toLocaleString('id-ID') }}</strong></div>
        </div>

        <button @click="closeSuccessModal" class="btn-primary btn-block">
          Selesai & Kembali ke Beranda
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { SeatItem } from '~/stores/booking'

const route = useRoute()
const router = useRouter()
const { fetchApi } = useApi()

const authStore = useAuthStore()
const bookingStore = useBookingStore()

const scheduleId = route.params.id
const schedule = computed(() => bookingStore.selectedSchedule)
const selectedSeat = ref<SeatItem | null>(null)

const actionLoading = ref(false)
const message = ref<string | null>(null)
const messageType = ref<'success' | 'error'>('success')

const confirmedBooking = ref<any | null>(null)

const isLockedByMe = computed(() => {
  if (!selectedSeat.value || !bookingStore.activeLockedSeat) return false
  return selectedSeat.value.id === bookingStore.activeLockedSeat.id && bookingStore.lockRemainingSeconds > 0
})

const fetchSeats = async () => {
  const { data } = await fetchApi(`/schedules/${scheduleId}/seats`)
  if (data) {
    if (data.schedule) bookingStore.selectedSchedule = data.schedule
    if (data.seats) bookingStore.seats = data.seats
  }
}

const handleSelectSeat = (seat: SeatItem) => {
  if (seat.status === 'booked' || (seat.status === 'locked' && seat.id !== bookingStore.activeLockedSeat?.id)) {
    return
  }
  selectedSeat.value = seat
  message.value = null
}

const lockSelectedSeat = async () => {
  if (!authStore.isAuthenticated) {
    alert('Silakan login terlebih dahulu untuk mengunci kursi!')
    router.push('/login')
    return
  }

  if (!selectedSeat.value) return

  actionLoading.value = true
  message.value = null

  const { data, error, status } = await fetchApi(`/seats/${selectedSeat.value.id}/lock`, {
    method: 'POST',
  })

  actionLoading.value = false

  if (error) {
    messageType.value = 'error'
    if (status === 409) {
      message.value = 'Gagal! Kursi telah dikunci atau dipesan oleh pengguna lain.'
    } else {
      message.value = error.message || 'Gagal mengunci kursi'
    }
    fetchSeats()
  } else if (data && data.seat) {
    messageType.value = 'success'
    message.value = data.message || 'Kursi berhasil dikunci selama 5 menit!'
    selectedSeat.value = data.seat
    bookingStore.startLockTimer(data.seat.locked_until, data.seat)
    fetchSeats()
  }
}

const confirmSelectedBooking = async () => {
  if (!selectedSeat.value) return

  actionLoading.value = true
  message.value = null

  const { data, error, status } = await fetchApi(`/seats/${selectedSeat.value.id}/confirm`, {
    method: 'POST',
  })

  actionLoading.value = false

  if (error) {
    messageType.value = 'error'
    if (status === 410) {
      message.value = 'Waktu kunci kursi telah habis. Silakan pilih kembali.'
      bookingStore.clearLockTimer()
    } else {
      message.value = error.message || 'Gagal mengonfirmasi pesanan'
    }
    fetchSeats()
  } else if (data && data.booking) {
    confirmedBooking.value = data.booking
    bookingStore.clearLockTimer()
    fetchSeats()
  }
}

const closeSuccessModal = () => {
  confirmedBooking.value = null
  router.push('/')
}

let pollInterval: any = null

onMounted(() => {
  fetchSeats()
  pollInterval = setInterval(fetchSeats, 5000)
})

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval)
})
</script>

<style scoped>
.seat-selection-page {
  max-width: 1100px;
  margin: 0 auto;
  padding: 32px 24px;
}

.top-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.back-link {
  color: var(--color-primary);
  text-decoration: none;
  font-weight: 600;
  font-size: 0.9rem;
}

.schedule-summary {
  display: flex;
  gap: 10px;
  font-size: 0.85rem;
}

.operator-tag {
  background: rgba(56, 189, 248, 0.15);
  color: var(--color-primary);
  padding: 4px 10px;
  border-radius: 8px;
  font-weight: 700;
}

.route-tag, .date-tag {
  background: rgba(30, 41, 59, 0.8);
  border: 1px solid var(--color-border);
  padding: 4px 10px;
  border-radius: 8px;
}

.lock-banner-wrapper {
  margin-bottom: 24px;
}

.selection-grid {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 32px;
  align-items: start;
}

.booking-card {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.booking-card h3 {
  font-size: 1.2rem;
  font-weight: 700;
}

.selected-seat-box {
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid var(--color-border);
  border-radius: 12px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.box-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
}

.seat-pill {
  background: var(--color-primary);
  color: #0f172a;
  padding: 4px 12px;
  border-radius: 8px;
  font-weight: 800;
  font-size: 1.1rem;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.9rem;
}

.price {
  font-weight: 800;
  color: var(--color-primary);
}

.no-seat-box {
  text-align: center;
  padding: 24px;
  border: 1px dashed var(--color-border);
  border-radius: 12px;
  color: var(--color-text-muted);
  font-size: 0.85rem;
}

.no-seat-box .icon {
  font-size: 2rem;
  display: block;
  margin-bottom: 8px;
}

.btn-block {
  width: 100%;
}

.btn-confirm {
  background: linear-gradient(135deg, #10b981 0%, #34d399 100%);
  color: #064e3b;
  box-shadow: 0 4px 15px rgba(16, 185, 129, 0.3);
}

.msg-banner {
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 0.85rem;
  text-align: center;
}

.msg-banner.success {
  background: rgba(16, 185, 129, 0.15);
  border: 1px solid rgba(16, 185, 129, 0.4);
  color: #34d399;
}

.msg-banner.error {
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid rgba(239, 68, 68, 0.4);
  color: #fca5a5;
}

/* Modal Overlay */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(15, 23, 42, 0.85);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
  padding: 24px;
}

.modal-card {
  max-width: 440px;
  width: 100%;
  text-align: center;
}

.success-icon {
  font-size: 3.5rem;
  margin-bottom: 12px;
}

.modal-sub {
  font-size: 0.85rem;
  color: var(--color-text-muted);
  margin-top: 4px;
}

.ticket-code {
  background: rgba(56, 189, 248, 0.15);
  border: 2px dashed var(--color-primary);
  color: var(--color-primary);
  font-size: 1.8rem;
  font-weight: 800;
  letter-spacing: 2px;
  padding: 12px;
  border-radius: 12px;
  margin: 20px 0;
}

.ticket-details {
  background: rgba(15, 23, 42, 0.5);
  padding: 16px;
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  font-size: 0.85rem;
  margin-bottom: 24px;
  text-align: left;
}

.t-row {
  display: flex;
  justify-content: space-between;
}

@media (max-width: 868px) {
  .selection-grid {
    grid-template-columns: 1fr;
  }
}
</style>
