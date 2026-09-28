$ErrorActionPreference = 'Continue'
Set-Location 'd:\laragon\www\ITHelpDesk\backend'
$p = Start-Process -FilePath '.\bin\api.exe' -WorkingDirectory (Get-Location) -NoNewWindow -PassThru -RedirectStandardOutput 'bin\ready-out.log' -RedirectStandardError 'bin\ready-err.log'
$ok = $false
foreach ($i in 1..20) {
  Start-Sleep -Milliseconds 500
  try { Invoke-WebRequest -Uri 'http://127.0.0.1:8080/health/live' -UseBasicParsing -TimeoutSec 2 | Out-Null; $ok = $true; break } catch { }
}
if (-not $ok) { Write-Host 'SERVER TIDAK MERESPONS'; Get-Content 'bin\ready-err.log'; exit 1 }
$live = Invoke-WebRequest -Uri 'http://127.0.0.1:8080/health/live' -UseBasicParsing
Write-Host ("live  -> {0} {1}" -f $live.StatusCode, $live.Content)
try {
  $ready = Invoke-WebRequest -Uri 'http://127.0.0.1:8080/health/ready' -UseBasicParsing
  Write-Host ("ready -> {0} {1}" -f $ready.StatusCode, $ready.Content)
} catch {
  $resp = $_.Exception.Response
  $body = (New-Object IO.StreamReader($resp.GetResponseStream())).ReadToEnd()
  Write-Host ("ready -> {0} {1}" -f [int]$resp.StatusCode, $body)
}
taskkill /F /PID $p.Id | Out-Null
Write-Host '--- log api ---'
Get-Content 'bin\ready-out.log'
