<?php
namespace Tests\Feature;

use App\Models\User;
use App\Models\Operator;
use App\Models\BusRoute;
use App\Models\Schedule;
use App\Models\Seat;
use App\Models\Booking;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;
use Carbon\Carbon;

class BookingTest extends TestCase
{
    use RefreshDatabase;

    protected function setUp(): void
    {
        parent::setUp();
        $this->seed();
    }

    protected function getToken($email = 'user1@example.com', $password = 'password')
    {
        $response = $this->postJson('/api/login', [
            'email' => $email,
            'password' => $password,
        ]);
        return $response->json('token');
    }

    /** @test */
    public function user_can_login_with_valid_credentials()
    {
        $response = $this->postJson('/api/login', [
            'email' => 'user1@example.com',
            'password' => 'password',
        ]);

        $response->assertStatus(200)
            ->assertJsonStructure(['token', 'token_type', 'expires_in', 'user']);
    }

    /** @test */
    public function login_fails_with_invalid_credentials()
    {
        $response = $this->postJson('/api/login', [
            'email' => 'user1@example.com',
            'password' => 'wrongpassword',
        ]);

        $response->assertStatus(401);
    }

    /** @test */
    public function can_search_schedules()
    {
        $schedule = Schedule::first();
        $route = $schedule->busRoute;

        $response = $this->getJson('/api/schedules/search?origin=' . $route->origin . '&destination=' . $route->destination . '&date=' . $schedule->departure_date);

        $response->assertStatus(200)
            ->assertJsonStructure(['data']);
    }

    /** @test */
    public function can_view_seat_availability()
    {
        $schedule = Schedule::first();

        $response = $this->getJson('/api/schedules/' . $schedule->id . '/seats');

        $response->assertStatus(200)
            ->assertJsonStructure(['schedule', 'seats']);
    }

    /** @test */
    public function authenticated_user_can_lock_available_seat()
    {
        $token = $this->getToken();
        $seat = Seat::where('status', 'available')->first();

        $response = $this->withHeader('Authorization', 'Bearer ' . $token)
            ->postJson('/api/seats/' . $seat->id . '/lock');

        $response->assertStatus(200)
            ->assertJson(['message' => 'Kursi berhasil dikunci selama 5 menit']);

        $this->assertDatabaseHas('seats', ['id' => $seat->id, 'status' => 'locked']);
    }

    /** @test */
    public function lock_fails_when_seat_already_locked()
    {
        $token1 = $this->getToken('user1@example.com');
        $token2 = $this->getToken('user2@example.com');
        $seat = Seat::where('status', 'available')->first();

        $this->withHeader('Authorization', 'Bearer ' . $token1)
            ->postJson('/api/seats/' . $seat->id . '/lock')
            ->assertStatus(200);

        $response = $this->withHeader('Authorization', 'Bearer ' . $token2)
            ->postJson('/api/seats/' . $seat->id . '/lock');

        $response->assertStatus(409);
    }

    /** @test */
    public function lock_fails_without_authentication()
    {
        $seat = Seat::where('status', 'available')->first();

        $response = $this->postJson('/api/seats/' . $seat->id . '/lock');

        $response->assertStatus(401);
    }

    /** @test */
    public function user_can_confirm_locked_seat()
    {
        $token = $this->getToken();
        $seat = Seat::where('status', 'available')->first();

        $this->withHeader('Authorization', 'Bearer ' . $token)
            ->postJson('/api/seats/' . $seat->id . '/lock');

        $response = $this->withHeader('Authorization', 'Bearer ' . $token)
            ->postJson('/api/seats/' . $seat->id . '/confirm');

        $response->assertStatus(200)
            ->assertJsonStructure(['booking' => ['booking_code']]);

        $this->assertDatabaseHas('seats', ['id' => $seat->id, 'status' => 'booked']);
        $this->assertDatabaseHas('bookings', ['seat_id' => $seat->id, 'status' => 'confirmed']);
    }

    /** @test */
    public function confirm_fails_after_lock_expires()
    {
        $token = $this->getToken();
        $seat = Seat::where('status', 'available')->first();

        $this->withHeader('Authorization', 'Bearer ' . $token)
            ->postJson('/api/seats/' . $seat->id . '/lock');

        Carbon::setTestNow(now()->addMinutes(6));

        $response = $this->withHeader('Authorization', 'Bearer ' . $token)
            ->postJson('/api/seats/' . $seat->id . '/confirm');

        $response->assertStatus(410);

        Carbon::setTestNow();
    }

    /** @test */
    public function expired_lock_releases_seat_as_available()
    {
        $token = $this->getToken();
        $seat = Seat::where('status', 'available')->first();

        $this->withHeader('Authorization', 'Bearer ' . $token)
            ->postJson('/api/seats/' . $seat->id . '/lock');

        Carbon::setTestNow(now()->addMinutes(6));

        $schedule = $seat->schedule;
        $response = $this->getJson('/api/schedules/' . $schedule->id . '/seats');

        $response->assertStatus(200);
        
        $seatData = collect($response->json('seats'))->firstWhere('id', $seat->id);
        $this->assertEquals('available', $seatData['status']);

        Carbon::setTestNow();
    }
}
