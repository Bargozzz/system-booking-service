param(
    [string]$BaseUrl = "http://localhost:8000/api",
    [int]$ConcurrentRequests = 50
)

Write-Host "=== Mini Booking Service - Load Test ===" -ForegroundColor Cyan
Write-Host "Target: $BaseUrl"
Write-Host "Concurrent Requests: $ConcurrentRequests"
Write-Host ""

Write-Host "1. Logging in..." -ForegroundColor Yellow
$loginBody = @{ email = "user1@example.com"; password = "password" } | ConvertTo-Json
$loginResponse = Invoke-RestMethod -Uri "$BaseUrl/login" -Method Post -Body $loginBody -ContentType "application/json"
$token = $loginResponse.token
Write-Host "   Token acquired." -ForegroundColor Green

Write-Host "2. Finding available seat..." -ForegroundColor Yellow
$scheduleId = 1
$seatsResponse = Invoke-RestMethod -Uri "$BaseUrl/schedules/$scheduleId/seats" -Method Get
$availableSeat = ($seatsResponse.seats | Where-Object { $_.status -eq "available" } | Select-Object -First 1)
if (-not $availableSeat) {
    Write-Host "   No available seats found! Run: php artisan migrate:fresh --seed" -ForegroundColor Red
    exit 1
}
$seatId = $availableSeat.id
Write-Host "   Target Seat ID: $seatId ($($availableSeat.seat_number))" -ForegroundColor Green

Write-Host "3. Sending $ConcurrentRequests concurrent lock requests to seat $seatId..." -ForegroundColor Yellow

$jobs = @()
for ($i = 1; $i -le $ConcurrentRequests; $i++) {
    $jobs += Start-Job -ScriptBlock {
        param($url, $token, $seatId)
        try {
            $headers = @{ Authorization = "Bearer $token"; "Content-Type" = "application/json" }
            $response = Invoke-WebRequest -Uri "$url/seats/$seatId/lock" -Method Post -Headers $headers -UseBasicParsing
            return @{ StatusCode = $response.StatusCode; Body = $response.Content }
        } catch {
            $statusCode = $_.Exception.Response.StatusCode.value__
            return @{ StatusCode = $statusCode; Body = $_.ErrorDetails.Message }
        }
    } -ArgumentList $BaseUrl, $token, $seatId
}

$results = $jobs | Wait-Job | Receive-Job
$jobs | Remove-Job

$successCount = ($results | Where-Object { $_.StatusCode -eq 200 }).Count
$conflictCount = ($results | Where-Object { $_.StatusCode -eq 409 }).Count
$otherCount = $ConcurrentRequests - $successCount - $conflictCount

Write-Host ""
Write-Host "=== RESULTS ===" -ForegroundColor Cyan
Write-Host "  Success (200): $successCount" -ForegroundColor Green
Write-Host "  Conflict (409): $conflictCount" -ForegroundColor Yellow
Write-Host "  Other Errors:  $otherCount" -ForegroundColor Red
Write-Host ""

if ($successCount -eq 1) {
    Write-Host "PASS: Exactly 1 request succeeded. Race condition handling is working!" -ForegroundColor Green
} elseif ($successCount -eq 0) {
    Write-Host "WARNING: No requests succeeded. Check if the server is running." -ForegroundColor Yellow
} else {
    Write-Host "FAIL: $successCount requests succeeded. Double-booking detected!" -ForegroundColor Red
}
