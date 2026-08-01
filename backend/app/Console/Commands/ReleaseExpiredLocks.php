<?php
namespace App\Console\Commands;

use App\Models\Seat;
use Illuminate\Console\Command;

class ReleaseExpiredLocks extends Command
{
    protected $signature = 'seats:release-expired';
    protected $description = 'Release seats with expired locks';

    public function handle()
    {
        $count = Seat::where('status', 'locked')
            ->where('locked_until', '<', now())
            ->update([
                'status' => 'available',
                'locked_by' => null,
                'locked_until' => null,
            ]);

        $this->info("Released {$count} expired seat locks.");
        return Command::SUCCESS;
    }
}
