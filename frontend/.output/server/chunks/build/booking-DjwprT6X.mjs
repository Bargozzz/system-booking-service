import { defineStore } from 'pinia';
import { ref } from 'vue';

const intervalError = "[nuxt] `setInterval` should not be used on the server. Consider wrapping it with an `onNuxtReady`, `onBeforeMount` or `onMounted` lifecycle hook, or ensure you only call it in the browser by checking `false`.";
const setInterval = () => {
  console.error(intervalError);
};
const useBookingStore = defineStore("booking", () => {
  const searchOrigin = ref("Jakarta");
  const searchDestination = ref("Bandung");
  const searchDate = ref(new Date(Date.now() + 864e5).toISOString().split("T")[0]);
  const schedules = ref([]);
  const selectedSchedule = ref(null);
  const seats = ref([]);
  const activeLockedSeat = ref(null);
  const lockExpiresAt = ref(null);
  const lockRemainingSeconds = ref(0);
  let timerInterval = null;
  const startLockTimer = (lockedUntilIso, seat) => {
    activeLockedSeat.value = seat;
    const expiryMs = new Date(lockedUntilIso).getTime();
    lockExpiresAt.value = expiryMs;
    if (timerInterval) clearInterval(timerInterval);
    const updateTimer = () => {
      const nowMs = Date.now();
      const diffSec = Math.max(0, Math.floor((expiryMs - nowMs) / 1e3));
      lockRemainingSeconds.value = diffSec;
      if (diffSec <= 0) {
        clearInterval(timerInterval);
        activeLockedSeat.value = null;
        lockExpiresAt.value = null;
      }
    };
    updateTimer();
    timerInterval = setInterval();
  };
  const clearLockTimer = () => {
    if (timerInterval) clearInterval(timerInterval);
    activeLockedSeat.value = null;
    lockExpiresAt.value = null;
    lockRemainingSeconds.value = 0;
  };
  return {
    searchOrigin,
    searchDestination,
    searchDate,
    schedules,
    selectedSchedule,
    seats,
    activeLockedSeat,
    lockExpiresAt,
    lockRemainingSeconds,
    startLockTimer,
    clearLockTimer
  };
});

export { useBookingStore as u };
//# sourceMappingURL=booking-DjwprT6X.mjs.map
