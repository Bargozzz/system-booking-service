<?php
namespace App\Models;
use Illuminate\Database\Eloquent\Model;

class Seat extends Model
{
    protected $fillable = ['schedule_id', 'seat_number', 'status', 'locked_by', 'locked_until'];

    protected function casts(): array
    {
        return ['locked_until' => 'datetime'];
    }

    public function schedule()
    {
        return $this->belongsTo(Schedule::class);
    }

    public function lockedByUser()
    {
        return $this->belongsTo(User::class, 'locked_by');
    }

    public function booking()
    {
        return $this->hasOne(Booking::class);
    }

    public function isLockExpired(): bool
    {
        return $this->status === 'locked' && $this->locked_until && $this->locked_until->isPast();
    }

    public function getEffectiveStatusAttribute(): string
    {
        if ($this->isLockExpired()) {
            return 'available';
        }
        return $this->status;
    }
}
