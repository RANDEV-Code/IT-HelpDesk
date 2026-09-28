$ErrorActionPreference = 'Continue'
Set-Location 'd:\laragon\www\ITHelpDesk\backend'

# Tambah toolchain MinGW-w64 (WinLibs POSIX/UCRT) ke PATH untuk CGO/race.
$mingwBin = Get-ChildItem -Path "$env:LOCALAPPDATA\Microsoft\WinGet\Packages" -Recurse -Filter 'gcc.exe' -ErrorAction SilentlyContinue | Select-Object -First 1 -ExpandProperty DirectoryName
if (-not $mingwBin) { Write-Output 'GCC tidak ditemukan; batalkan'; exit 1 }
$env:PATH = "$mingwBin;$env:PATH"
$env:CGO_ENABLED = '1'
$env:CC = 'gcc'
Write-Output "gcc: $mingwBin"
& gcc --version | Select-Object -First 1

Write-Output '=== go env (CGO) ==='
go env CGO_ENABLED CC

Write-Output '=== go test -race ./... (unit + integrasi terhadap randesk_test) ==='
$env:TEST_DATABASE_URL = 'postgres://randesk_migrate:randesk_migrate_dev@127.0.0.1:5432/randesk_test?sslmode=disable'
go test -race ./... 2>&1
Write-Output "exit=$LASTEXITCODE"
