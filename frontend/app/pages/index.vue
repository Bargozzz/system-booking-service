<template>
  <div class="search-page">
    <!-- Hero Banner -->
    <section class="hero-section">
      <div class="hero-content">
        <h1>Pesan Tiket Bus Real-Time</h1>
        <p>Sistem reservasi instan dengan fitur penguncian kursi transaksional 5 menit.</p>
      </div>

      <!-- Search Form Card -->
      <div class="search-card glass-card">
        <form @submit.prevent="onSearch" class="search-form">
          <div class="form-group">
            <label>Kota Asal</label>
            <select v-model="bookingStore.searchOrigin" class="input-field">
              <option value="Jakarta">Jakarta</option>
              <option value="Bandung">Bandung</option>
              <option value="Surabaya">Surabaya</option>
              <option value="Malang">Malang</option>
            </select>
          </div>

          <div class="form-group">
            <label>Kota Tujuan</label>
            <select v-model="bookingStore.searchDestination" class="input-field">
              <option value="Bandung">Bandung</option>
              <option value="Jakarta">Jakarta</option>
              <option value="Malang">Malang</option>
              <option value="Surabaya">Surabaya</option>
            </select>
          </div>

          <div class="form-group">
            <label>Tanggal Keberangkatan</label>
            <input 
              v-model="bookingStore.searchDate" 
              type="date" 
              class="input-field" 
              required
            />
          </div>

          <button type="submit" class="btn-primary btn-search" :disabled="loading">
            <span v-if="loading">Mencari...</span>
            <span v-else>🔍 Cari Jadwal</span>
          </button>
        </form>
      </div>
    </section>

    <!-- Schedule Results Section -->
    <section class="results-section">
      <div v-if="searched" class="section-header">
        <h2>Hasil Pencarian ({{ bookingStore.schedules.length }})</h2>
        <p>Menampilkan jadwal dari {{ bookingStore.searchOrigin }} ke {{ bookingStore.searchDestination }}</p>
      </div>

      <div v-if="loading" class="loading-state glass-card">
        <p>Memuat data jadwal bus...</p>
      </div>

      <div v-else-if="searched && bookingStore.schedules.length === 0" class="empty-state glass-card">
        <span class="empty-icon">🚫</span>
        <h3>Tidak Ada Jadwal Ditemukan</h3>
        <p>Coba pilih kombinasi kota asal/tujuan atau tanggal keberangkatan yang lain.</p>
      </div>

      <div v-else class="schedule-grid">
        <div 
          v-for="schedule in bookingStore.schedules" 
          :key="schedule.id"
          class="schedule-card glass-card"
        >
          <div class="operator-info">
            <div class="operator-badge">
              <span class="operator-code">{{ schedule.operator_code }}</span>
              <span class="operator-name">{{ schedule.operator }}</span>
            </div>
            <div class="seats-badge badge" :class="schedule.available_seats ? 'badge-success' : 'badge-warning'">
              {{ schedule.available_seats }} dari {{ schedule.total_seats }} kursi tersedia
            </div>
          </div>

          <div class="route-details">
            <div class="time-block">
              <span class="time">{{ schedule.departure_time.substring(0, 5) }}</span>
              <span class="city">{{ schedule.origin }}</span>
            </div>
            <div class="route-line">
              <span class="duration">Langsung</span>
              <div class="line"></div>
            </div>
            <div class="time-block">
              <span class="time">Tiba</span>
              <span class="city">{{ schedule.destination }}</span>
            </div>
          </div>

          <div class="card-footer">
            <div class="price-tag">
              <span class="price-label">Harga per kursi</span>
              <span class="price-val">Rp {{ schedule.price.toLocaleString('id-ID') }}</span>
            </div>

            <button 
              @click="selectSchedule(schedule)" 
              class="btn-primary"
              :disabled="!schedule.available_seats"
            >
              Pilih Kursi ➔
            </button>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
const bookingStore = useBookingStore()
const { fetchApi } = useApi()
const router = useRouter()

const loading = ref(false)
const searched = ref(false)

const onSearch = async () => {
  loading.value = true
  searched.value = true

  const { data, error } = await fetchApi(
    `/schedules/search?origin=${bookingStore.searchOrigin}&destination=${bookingStore.searchDestination}&date=${bookingStore.searchDate}`
  )

  loading.value = false
  if (data && data.data) {
    bookingStore.schedules = data.data
  }
}

const selectSchedule = (schedule: any) => {
  bookingStore.selectedSchedule = schedule
  router.push(`/schedules/${schedule.id}`)
}

onMounted(() => {
  onSearch()
})
</script>

<style scoped>
.search-page {
  max-width: 1100px;
  margin: 0 auto;
  padding: 32px 24px;
}

.hero-section {
  text-align: center;
  margin-bottom: 40px;
}

.hero-content h1 {
  font-size: 2.2rem;
  font-weight: 800;
  background: linear-gradient(135deg, #f8fafc 0%, #38bdf8 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.hero-content p {
  color: var(--color-text-muted);
  font-size: 1rem;
  margin-top: 8px;
}

.search-card {
  margin-top: 28px;
  padding: 24px;
}

.search-form {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr auto;
  gap: 16px;
  align-items: end;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  text-align: left;
}

.form-group label {
  font-size: 0.85rem;
  font-weight: 600;
}

.btn-search {
  height: 46px;
}

.results-section {
  margin-top: 32px;
}

.section-header {
  margin-bottom: 20px;
}

.section-header h2 {
  font-size: 1.4rem;
  font-weight: 700;
}

.section-header p {
  font-size: 0.85rem;
  color: var(--color-text-muted);
}

.schedule-grid {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.schedule-card {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.operator-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.operator-badge {
  display: flex;
  align-items: center;
  gap: 8px;
}

.operator-code {
  background: rgba(56, 189, 248, 0.15);
  color: var(--color-primary);
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 0.8rem;
  font-weight: 800;
}

.operator-name {
  font-weight: 700;
  font-size: 1.05rem;
}

.route-details {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 12px 0;
  border-top: 1px solid var(--color-border);
  border-bottom: 1px solid var(--color-border);
}

.time-block {
  display: flex;
  flex-direction: column;
}

.time {
  font-size: 1.3rem;
  font-weight: 800;
}

.city {
  font-size: 0.85rem;
  color: var(--color-text-muted);
}

.route-line {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.duration {
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.line {
  width: 100%;
  height: 2px;
  background: var(--color-border);
  position: relative;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.price-tag {
  display: flex;
  flex-direction: column;
}

.price-label {
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.price-val {
  font-size: 1.25rem;
  font-weight: 800;
  color: var(--color-primary);
}

.empty-state {
  text-align: center;
  padding: 48px;
}

.empty-icon {
  font-size: 3rem;
}

@media (max-width: 768px) {
  .search-form {
    grid-template-columns: 1fr;
  }
}
</style>
