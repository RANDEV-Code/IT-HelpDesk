$ErrorActionPreference = 'Continue'
Set-Location 'd:\laragon\www\ITHelpDesk\backend'

Write-Host '=== 1. Build exe ==='
go build -o bin/api.exe ./cmd/api
if ($LASTEXITCODE -ne 0) { Write-Host 'BUILD GAGAL'; exit 1 }
Write-Host 'build OK'

Write-Host '=== 2. Start server (background, log ke file) ==='
$stdout = 'd:\laragon\www\ITHelpDesk\backend\bin\smoke-out.log'
$stderr = 'd:\laragon\www\ITHelpDesk\backend\bin\smoke-err.log'
$p = Start-Process -FilePath '.\bin\api.exe' -WorkingDirectory 'd:\laragon\www\ITHelpDesk\backend' -RedirectStandardOutput $stdout -RedirectStandardError $stderr -PassThru -NoNewWindow
Write-Host ("PID = {0}" -f $p.Id)

# tunggu server siap menjawab (maks 10 detik)
$ok = $false
foreach ($i in 1..20) {
  Start-Sleep -Milliseconds 500
  try { Invoke-WebRequest -Uri 'http://127.0.0.1:8080/health/live' -UseBasicParsing -TimeoutSec 2 | Out-Null; $ok = $true; break } catch { }
}
if (-not $ok) { Write-Host 'SERVER TIDAK MERESPONS'; Get-Content $stderr; exit 1 }

Write-Host '=== 3. GET /health/live (harapkan 200) ==='
$live = Invoke-WebRequest -Uri 'http://127.0.0.1:8080/health/live' -UseBasicParsing
Write-Host ("live  -> {0} {1}" -f $live.StatusCode, $live.Content)

Write-Host '=== 4. GET /health/ready (harapkan 503, DB mati) ==='
try {
  $ready = Invoke-WebRequest -Uri 'http://127.0.0.1:8080/health/ready' -UseBasicParsing
  Write-Host ("ready -> {0} {1}" -f $ready.StatusCode, $ready.Content)
} catch {
  $resp = $_.Exception.Response
  $body = (New-Object IO.StreamReader($resp.GetResponseStream())).ReadToEnd()
  Write-Host ("ready -> {0} {1}" -f [int]$resp.StatusCode, $body)
}

Write-Host '=== 5. Hentikan server ==='
# Catatan: graceful shutdown (SIGTERM/Ctrl+C <=15 detik) dibuktikan oleh tes Go
# TestGracefulShutdown di backend/cmd/api/shutdown_test.go — taskkill tanpa /F
# tidak terkirim ke proses tanpa window, jadi di sini cukup hentikan paksa.
taskkill /F /PID $p.Id | Out-Null
$p.WaitForExit(10000) | Out-Null
Write-Host ("server dihentikan; HasExited={0}" -f $p.HasExited)

Write-Host '--- log stdout ---'
Get-Content $stdout
Write-Host '--- log stderr ---'
Get-Content $stderr -ErrorAction SilentlyContinue

Write-Host '=== 6. Production tanpa konfigurasi aman harus menolak start ==='
$env:APP_ENV = 'production'
$env:APP_ORIGIN = 'http://helpdesk.example.com'
$env:HTTP_ADDR = '127.0.0.1:8081'
$env:DATABASE_URL = 'postgres://u:p@localhost:5432/d'
$env:STORAGE_ROOT = './storage'
$env:LOG_LEVEL = 'info'
& .\bin\api.exe
Write-Host ("exit code = {0} (harapkan 1)" -f $LASTEXITCODE)
Remove-Item Env:APP_ENV, Env:APP_ORIGIN, Env:HTTP_ADDR, Env:DATABASE_URL, Env:STORAGE_ROOT, Env:LOG_LEVEL
Write-Host '=== SMOKE TEST SELESAI ==='
