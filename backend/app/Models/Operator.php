<?php
namespace App\Models;
use Illuminate\Database\Eloquent\Model;

class Operator extends Model
{
    protected $fillable = ['name', 'code'];
    
    public function busRoutes()
    {
        return $this->hasMany(BusRoute::class);
    }
}
