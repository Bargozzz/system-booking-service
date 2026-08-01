<?php
namespace App\Models;
use Illuminate\Database\Eloquent\Model;

class Booking extends Model
{
    protected $fillable = ['user_id', 'seat_id', 'schedule_id', 'booking_code', 'status'];

    public function user()
    {
        return $this->belongsTo(User::class);
    }

    public function seat()
    {
        return $this->belongsTo(Seat::class);
    }

    public function schedule()
    {
        return $this->belongsTo(Schedule::class);
    }
}
