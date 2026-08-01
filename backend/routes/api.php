<?php
use Illuminate\Support\Facades\Route;
use App\Http\Controllers\AuthController;
use App\Http\Controllers\ScheduleController;
use App\Http\Controllers\BookingController;

Route::post('/login', [AuthController::class, 'login']);

Route::middleware('throttle:search')->group(function () {
    Route::get('/schedules/search', [ScheduleController::class, 'search']);
});

Route::get('/schedules/{id}/seats', [ScheduleController::class, 'seats']);

Route::middleware('auth:api')->group(function () {
    Route::get('/me', [AuthController::class, 'me']);
    Route::post('/logout', [AuthController::class, 'logout']);
    Route::post('/seats/{id}/lock', [BookingController::class, 'lockSeat']);
    Route::post('/seats/{id}/confirm', [BookingController::class, 'confirmBooking']);
});
