$ErrorActionPreference = 'Continue'
Set-Location 'd:\laragon\www\ITHelpDesk\backend'

# Jalankan API dengan logging JSON (APP_ENV != development) di latar,
# tanpa DB (health/live tidak menyentuh DB), lalu rekam satu baris log request.
$env:APP_ENV = 'staging'
$env:APP_ORIGIN = 'https://example.test'
$env:HTTP_ADDR = '127.0.0.1:18098'
$env:DATABASE_URL = 'postgres://u:p@127.0.0.1:5432/none?sslmode=disable'
$env:STORAGE_ROOT = Join-Path $env:TEMP 'randesk-log-sample'
$env:LOG_LEVEL = 'info'
$env:SESSION_TTL_HOURS = '8'

$logFile = Join-Path $env:TEMP 'randesk-task004-log.txt'
if (Test-Path $logFile) { Remove-Item $logFile -Force }

$proc = Start-Process -FilePath 'go' -ArgumentList 'run','./cmd/api' `
    -RedirectStandardOutput $logFile -RedirectStandardError (Join-Path $env:TEMP 'randesk-task004-err.txt') `
    -PassThru -WindowStyle Hidden

# tunggu server siap (poll health/live maks ~20 detik)
$ok = $false
for ($i = 0; $i -lt 40; $i++) {
    Start-Sleep -Milliseconds 500
    try {
        $resp = Invoke-WebRequest -Uri 'http://127.0.0.1:18098/health/live' -UseBasicParsing -TimeoutSec 2
        if ($resp.StatusCode -eq 200) { $ok = $true; break }
    } catch { }
}

if ($ok) {
    Write-Output '=== health/live -> 200 (X-Request-ID header) ==='
    try {
        $r = Invoke-WebRequest -Uri 'http://127.0.0.1:18098/health/live' -UseBasicParsing -TimeoutSec 2
        Write-Output ('X-Request-ID: ' + $r.Headers['X-Request-ID'])
    } catch { }
} else {
    Write-Output '=== server tidak siap dalam 20 detik ==='
}

# hentikan proses (dan anak proses go run) secara paksa
try { taskkill /PID $proc.Id /T /F | Out-Null } catch { }

Start-Sleep -Milliseconds 300
Write-Output '=== contoh baris log JSON (middleware Logging) ==='
if (Test-Path $logFile) {
    Get-Content $logFile | Where-Object { $_ -match '"msg":"http request"' } | Select-Object -First 3
} else {
    Write-Output '(log file tidak ditemukan)'
}
