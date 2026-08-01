<template>
  <div class="bus-container glass-card">
    <div class="bus-front">
      <span class="driver-wheel">☸️ Supir Bus</span>
    </div>

    <!-- Seat Grid (4 cols, 5 rows) -->
    <div class="seat-grid">
      <div v-for="row in rows" :key="row" class="grid-row">
        <!-- Left 2 seats (Col 1, 2) -->
        <div class="seat-pair">
          <SeatItem
            v-for="col in [1, 2]"
            :key="`${row}${col}`"
            :seat="getSeat(`${row}${col}`)"
            :is-selected="selectedSeatId === getSeat(`${row}${col}`).id"
            @select="onSelectSeat"
          />
        </div>

        <!-- Aisle spacing -->
        <div class="aisle">
          <span class="row-label">{{ row }}</span>
        </div>

        <!-- Right 2 seats (Col 3, 4) -->
        <div class="seat-pair">
          <SeatItem
            v-for="col in [3, 4]"
            :key="`${row}${col}`"
            :seat="getSeat(`${row}${col}`)"
            :is-selected="selectedSeatId === getSeat(`${row}${col}`).id"
            @select="onSelectSeat"
          />
        </div>
      </div>
    </div>

    <!-- Legend -->
    <div class="seat-legend">
      <div class="legend-item"><span class="legend-dot status-available"></span> Tersedia</div>
      <div class="legend-item"><span class="legend-dot status-locked-user"></span> Dikunci Anda</div>
      <div class="legend-item"><span class="legend-dot status-locked-other"></span> Dikunci User Lain</div>
      <div class="legend-item"><span class="legend-dot status-booked"></span> Terisi</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { SeatItem as ISeatItem } from '~/stores/booking'

const props = defineProps<{
  seats: ISeatItem[]
  selectedSeatId: number | null
}>()

const emit = defineEmits(['select-seat'])

const rows = ['A', 'B', 'C', 'D', 'E']

const fallbackSeat = (code: string): ISeatItem => ({
  id: 0,
  seat_number: code,
  status: 'available',
  locked_until: null,
})

const getSeat = (code: string): ISeatItem => {
  return props.seats.find(s => s.seat_number === code) || fallbackSeat(code)
}

const onSelectSeat = (seat: ISeatItem) => {
  if (seat.id !== 0) {
    emit('select-seat', seat)
  }
}
</script>

<style scoped>
.bus-container {
  max-width: 480px;
  margin: 0 auto;
  padding: 24px;
}

.bus-front {
  background: rgba(15, 23, 42, 0.8);
  border: 1px dashed var(--color-border);
  padding: 10px;
  border-radius: 10px;
  text-align: right;
  margin-bottom: 24px;
  font-size: 0.85rem;
  color: var(--color-text-muted);
}

.seat-grid {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.grid-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.seat-pair {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  flex: 1;
}

.aisle {
  width: 36px;
  text-align: center;
}

.row-label {
  font-size: 0.8rem;
  font-weight: 700;
  color: var(--color-text-muted);
}

.seat-legend {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  margin-top: 28px;
  padding-top: 16px;
  border-top: 1px solid var(--color-border);
  font-size: 0.8rem;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.legend-dot {
  width: 14px;
  height: 14px;
  border-radius: 4px;
}

.status-available { background: var(--seat-available-border); }
.status-locked-user { background: var(--seat-locked-user-border); }
.status-locked-other { background: var(--seat-locked-other-border); }
.status-booked { background: var(--seat-booked-border); }
</style>
