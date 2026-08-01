<?php
namespace App\Models;
use Illuminate\Database\Eloquent\Model;

class BusRoute extends Model
{
    protected $fillable = ['operator_id', 'origin', 'destination', 'price'];

    public function operator()
    {
        return $this->belongsTo(Operator::class);
    }

    public function schedules()
    {
        return $this->hasMany(Schedule::class);
    }
}
