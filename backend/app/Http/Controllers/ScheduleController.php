<?php
namespace App\Http\Controllers;

use App\Models\Schedule;
use App\Models\Seat;
use Illuminate\Http\Request;

class ScheduleController extends Controller
{
    public function search(Request $request)
    {
        $request->validate([
            'origin' => 'required|string',
            'destination' => 'required|string',
            'date' => 'required|date',
        ]);

        $schedules = Schedule::whereHas('busRoute', function ($q) use ($request) {
            $q->where('origin', $request->origin)
              ->where('destination', $request->destination);
        })
        ->where('departure_date', $request->date)
        ->with(['busRoute.operator'])
        ->get()
        ->map(function ($schedule) {
            $availableSeats = $schedule->seats()
                ->where(function ($q) {
                    $q->where('status', 'available')
                      ->orWhere(function ($q2) {
                          $q2->where('status', 'locked')
                             ->where('locked_until', '<', now());
                      });
                })->count();

            return [
                'id' => $schedule->id,
                'operator' => $schedule->busRoute->operator->name,
                'operator_code' => $schedule->busRoute->operator->code,
                'origin' => $schedule->busRoute->origin,
                'destination' => $schedule->busRoute->destination,
                'departure_date' => $schedule->departure_date,
                'departure_time' => $schedule->departure_time,
                'price' => $schedule->busRoute->price,
                'available_seats' => $availableSeats,
                'total_seats' => 20,
            ];
        });

        return response()->json(['data' => $schedules]);
    }

    public function seats($scheduleId)
    {
        $schedule = Schedule::with(['busRoute.operator'])->findOrFail($scheduleId);

        Seat::where('schedule_id', $scheduleId)
            ->where('status', 'locked')
            ->where('locked_until', '<', now())
            ->update(['status' => 'available', 'locked_by' => null, 'locked_until' => null]);

        $seats = Seat::where('schedule_id', $scheduleId)
            ->orderBy('seat_number')
            ->get()
            ->map(function ($seat) {
                return [
                    'id' => $seat->id,
                    'seat_number' => $seat->seat_number,
                    'status' => $seat->status,
                    'locked_until' => $seat->status === 'locked' ? $seat->locked_until : null,
                ];
            });

        return response()->json([
            'schedule' => [
                'id' => $schedule->id,
                'operator' => $schedule->busRoute->operator->name,
                'operator_code' => $schedule->busRoute->operator->code,
                'origin' => $schedule->busRoute->origin,
                'destination' => $schedule->busRoute->destination,
                'departure_date' => $schedule->departure_date,
                'departure_time' => $schedule->departure_time,
                'price' => $schedule->busRoute->price,
            ],
            'seats' => $seats,
        ]);
    }
}
