import { defineStore } from 'pinia'

export interface SeatItem {
  id: number
  seat_number: string
  status: 'available' | 'locked' | 'booked'
  locked_until: string | null
}

export interface ScheduleItem {
  id: number
  operator: string
  operator_code: string
  origin: string
  destination: string
  departure_date: string
  departure_time: string
  price: number
  available_seats?: number
  total_seats?: number
}

export const useBookingStore = defineStore('booking', () => {
  const searchOrigin = ref('Jakarta')
  const searchDestination = ref('Bandung')
  const searchDate = ref(new Date(Date.now() + 86400000).toISOString().split('T')[0])
  
  const schedules = ref<ScheduleItem[]>([])
  const selectedSchedule = ref<ScheduleItem | null>(null)
  const seats = ref<SeatItem[]>([])
  
  const activeLockedSeat = ref<SeatItem | null>(null)
  const lockExpiresAt = ref<number | null>(null)
  const lockRemainingSeconds = ref(0)
  
  let timerInterval: any = null

  const startLockTimer = (lockedUntilIso: string, seat: SeatItem) => {
    activeLockedSeat.value = seat
    const expiryMs = new Date(lockedUntilIso).getTime()
    lockExpiresAt.value = expiryMs

    if (timerInterval) clearInterval(timerInterval)

    const updateTimer = () => {
      const nowMs = Date.now()
      const diffSec = Math.max(0, Math.floor((expiryMs - nowMs) / 1000))
      lockRemainingSeconds.value = diffSec

      if (diffSec <= 0) {
        clearInterval(timerInterval)
        activeLockedSeat.value = null
        lockExpiresAt.value = null
      }
    }

    updateTimer()
    timerInterval = setInterval(updateTimer, 1000)
  }

  const clearLockTimer = () => {
    if (timerInterval) clearInterval(timerInterval)
    activeLockedSeat.value = null
    lockExpiresAt.value = null
    lockRemainingSeconds.value = 0
  }

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
    clearLockTimer,
  }
})
