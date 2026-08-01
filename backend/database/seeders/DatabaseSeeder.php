<?php
namespace Database\Seeders;

use App\Models\User;
use App\Models\Operator;
use App\Models\BusRoute;
use App\Models\Schedule;
use App\Models\Seat;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\Hash;

class DatabaseSeeder extends Seeder
{
    public function run(): void
    {
        // Users
        $users = [
            ['name' => 'User Satu', 'email' => 'user1@example.com', 'password' => Hash::make('password')],
            ['name' => 'User Dua', 'email' => 'user2@example.com', 'password' => Hash::make('password')],
            ['name' => 'User Tiga', 'email' => 'user3@example.com', 'password' => Hash::make('password')],
        ];
        foreach ($users as $u) {
            User::create($u);
        }

        // Operators
        $op1 = Operator::create(['name' => 'Tiketux Express', 'code' => 'TKX']);
        $op2 = Operator::create(['name' => 'Nusantara Bus', 'code' => 'NTB']);

        // Routes
        $routes = [
            ['operator_id' => $op1->id, 'origin' => 'Jakarta', 'destination' => 'Bandung', 'price' => 150000],
            ['operator_id' => $op1->id, 'origin' => 'Bandung', 'destination' => 'Jakarta', 'price' => 150000],
            ['operator_id' => $op2->id, 'origin' => 'Surabaya', 'destination' => 'Malang', 'price' => 75000],
            ['operator_id' => $op2->id, 'origin' => 'Malang', 'destination' => 'Surabaya', 'price' => 75000],
        ];
        foreach ($routes as $r) {
            $route = BusRoute::create($r);

            // Create 2 schedules per route with future dates
            $dates = [
                now()->addDays(1)->format('Y-m-d'),
                now()->addDays(3)->format('Y-m-d'),
            ];
            $times = ['08:00:00', '14:00:00'];

            foreach ($dates as $i => $date) {
                $schedule = Schedule::create([
                    'bus_route_id' => $route->id,
                    'departure_date' => $date,
                    'departure_time' => $times[$i],
                ]);

                // Create 20 seats per schedule (A1-A4, B1-B4, C1-C4, D1-D4, E1-E4)
                $rows = ['A', 'B', 'C', 'D', 'E'];
                $cols = [1, 2, 3, 4];
                foreach ($rows as $row) {
                    foreach ($cols as $col) {
                        Seat::create([
                            'schedule_id' => $schedule->id,
                            'seat_number' => $row . $col,
                            'status' => 'available',
                        ]);
                    }
                }
            }
        }
    }
}
