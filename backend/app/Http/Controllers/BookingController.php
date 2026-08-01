<?php
namespace App\Http\Controllers;

use App\Models\Seat;
use App\Models\Booking;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Str;

class BookingController extends Controller
{
    public function lockSeat(Request $request, $seatId)
    {
        $userId = auth()->id();

        try {
            $result = DB::transaction(function () use ($seatId, $userId) {
                Seat::where('status', 'locked')
                    ->where('locked_until', '<', now())
                    ->lockForUpdate()
                    ->update(['status' => 'available', 'locked_by' => null, 'locked_until' => null]);

                $seat = Seat::where('id', $seatId)
                    ->lockForUpdate()
                    ->first();

                if (!$seat) {
                    return response()->json(['message' => 'Kursi tidak ditemukan'], 404);
                }

                if ($seat->status !== 'available') {
                    return response()->json(['message' => 'Kursi tidak tersedia. Sudah dikunci atau dipesan oleh pengguna lain.'], 409);
                }

                $lockedUntil = now()->addMinutes(5);
                $seat->update([
                    'status' => 'locked',
                    'locked_by' => $userId,
                    'locked_until' => $lockedUntil,
                ]);

                return response()->json([
                    'message' => 'Kursi berhasil dikunci selama 5 menit',
                    'seat' => [
                        'id' => $seat->id,
                        'seat_number' => $seat->seat_number,
                        'status' => 'locked',
                        'locked_until' => $lockedUntil->toISOString(),
                    ],
                    'schedule_id' => $seat->schedule_id,
                ]);
            });

            return $result;
        } catch (\Exception $e) {
            return response()->json(['message' => 'Gagal mengunci kursi: ' . $e->getMessage()], 500);
        }
    }

    public function confirmBooking(Request $request, $seatId)
    {
        $userId = auth()->id();

        try {
            $result = DB::transaction(function () use ($seatId, $userId) {
                $seat = Seat::where('id', $seatId)
                    ->lockForUpdate()
                    ->first();

                if (!$seat) {
                    return response()->json(['message' => 'Kursi tidak ditemukan'], 404);
                }

                if ($seat->status !== 'locked') {
                    return response()->json(['message' => 'Kursi tidak dalam status terkunci'], 400);
                }

                if ($seat->locked_by !== $userId) {
                    return response()->json(['message' => 'Kursi dikunci oleh pengguna lain'], 403);
                }

                if ($seat->locked_until && $seat->locked_until->isPast()) {
                    $seat->update(['status' => 'available', 'locked_by' => null, 'locked_until' => null]);
                    return response()->json(['message' => 'Waktu kunci kursi telah habis. Silakan coba lagi.'], 410);
                }

                $seat->update([
                    'status' => 'booked',
                    'locked_by' => null,
                    'locked_until' => null,
                ]);

                $booking = Booking::create([
                    'user_id' => $userId,
                    'seat_id' => $seat->id,
                    'schedule_id' => $seat->schedule_id,
                    'booking_code' => 'BK-' . strtoupper(Str::random(8)),
                    'status' => 'confirmed',
                ]);

                $booking->load(['seat', 'schedule.busRoute.operator']);

                return response()->json([
                    'message' => 'Pemesanan berhasil dikonfirmasi!',
                    'booking' => [
                        'id' => $booking->id,
                        'booking_code' => $booking->booking_code,
                        'status' => $booking->status,
                        'seat_number' => $booking->seat->seat_number,
                        'operator' => $booking->schedule->busRoute->operator->name,
                        'origin' => $booking->schedule->busRoute->origin,
                        'destination' => $booking->schedule->busRoute->destination,
                        'departure_date' => $booking->schedule->departure_date,
                        'departure_time' => $booking->schedule->departure_time,
                        'price' => $booking->schedule->busRoute->price,
                        'created_at' => $booking->created_at,
                    ],
                ]);
            });

            return $result;
        } catch (\Exception $e) {
            return response()->json(['message' => 'Gagal mengkonfirmasi pesanan: ' . $e->getMessage()], 500);
        }
    }
}
