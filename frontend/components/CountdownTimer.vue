<template>
  <div v-if="seconds > 0" class="timer-card pulse-lock">
    <span class="timer-icon">⏱️</span>
    <div class="timer-content">
      <span class="timer-title">Kursi Dikunci Sementara</span>
      <span class="timer-clock">{{ formattedTime }}</span>
    </div>
  </div>
  <div v-else-if="expired" class="timer-card timer-expired">
    <span class="timer-icon">⚠️</span>
    <span class="timer-title">Waktu kunci kursi habis! Silakan pilih kembali.</span>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  seconds: number
}>()

const expired = computed(() => props.seconds <= 0)

const formattedTime = computed(() => {
  const m = Math.floor(props.seconds / 60)
  const s = props.seconds % 60
  return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
})
</script>

<style scoped>
.timer-card {
  display: flex;
  align-items: center;
  gap: 12px;
  background: rgba(245, 158, 11, 0.15);
  border: 1px solid rgba(245, 158, 11, 0.5);
  padding: 12px 20px;
  border-radius: 12px;
  color: #fcd34d;
}

.timer-expired {
  background: rgba(239, 68, 68, 0.15);
  border-color: rgba(239, 68, 68, 0.5);
  color: #fca5a5;
}

.timer-icon {
  font-size: 1.4rem;
}

.timer-content {
  display: flex;
  flex-direction: column;
}

.timer-title {
  font-size: 0.8rem;
  opacity: 0.9;
}

.timer-clock {
  font-size: 1.3rem;
  font-weight: 800;
  font-variant-numeric: tabular-nums;
  letter-spacing: 1px;
}
</style>
