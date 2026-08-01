<template>
  <button 
    class="seat-btn"
    :class="[
      statusClass,
      { 'is-selected': isSelected }
    ]"
    :disabled="status === 'booked' || (status === 'locked' && !isLockedByMe)"
    @click="$emit('select', seat)"
  >
    <span class="seat-num">{{ seat.seat_number }}</span>
    <span class="seat-status-label">{{ statusText }}</span>
  </button>
</template>

<script setup lang="ts">
import type { SeatItem } from '~/stores/booking'

const props = defineProps<{
  seat: SeatItem
  isSelected: boolean
  currentUserId?: number | null
}>()

defineEmits(['select'])

const isLockedByMe = computed(() => {
  // If user locked it in current session or status is locked
  return props.isSelected
})

const statusClass = computed(() => {
  if (props.isSelected) return 'status-locked-user pulse-lock'
  if (props.seat.status === 'booked') return 'status-booked'
  if (props.seat.status === 'locked') return 'status-locked-other'
  return 'status-available'
})

const statusText = computed(() => {
  if (props.isSelected) return 'Dikunci Anda'
  if (props.seat.status === 'booked') return 'Terisi'
  if (props.seat.status === 'locked') return 'Dikunci User Lain'
  return 'Tersedia'
})
</script>

<style scoped>
.seat-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 64px;
  border-radius: 12px;
  border: 1.5px solid transparent;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.seat-btn:hover:not(:disabled) {
  transform: translateY(-3px) scale(1.03);
}

.seat-num {
  font-size: 1rem;
  font-weight: 800;
}

.seat-status-label {
  font-size: 0.65rem;
  font-weight: 600;
  opacity: 0.85;
}

/* Status variants */
.status-available {
  background: var(--seat-available-bg);
  border-color: var(--seat-available-border);
  color: var(--seat-available-text);
}
.status-available:hover {
  background: #047857;
}

.status-locked-user {
  background: var(--seat-locked-user-bg);
  border-color: var(--seat-locked-user-border);
  color: var(--seat-locked-user-text);
}

.status-locked-other {
  background: var(--seat-locked-other-bg);
  border-color: var(--seat-locked-other-border);
  color: var(--seat-locked-other-text);
  cursor: not-allowed;
  opacity: 0.75;
}

.status-booked {
  background: var(--seat-booked-bg);
  border-color: var(--seat-booked-border);
  color: var(--seat-booked-text);
  cursor: not-allowed;
  opacity: 0.5;
  text-decoration: line-through;
}
</style>
