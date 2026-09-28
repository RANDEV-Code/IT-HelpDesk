$ErrorActionPreference = 'Continue'  # psql menulis error SQL ke stderr; penilaian lewat output
$psql = 'C:\Program Files\PostgreSQL\17\bin\psql.exe'
$migrate = 'C:\Users\ricoa\go\bin\migrate.exe'
$env:PGPASSWORD = 'postgres'
Set-Location 'd:\laragon\www\ITHelpDesk'

Write-Host '=== 1. Buat role (dev-bootstrap.sql) ==='
& $psql -U postgres -h 127.0.0.1 -v ON_ERROR_STOP=1 -f backend/scripts/dev-bootstrap.sql
if ($LASTEXITCODE -ne 0) { throw 'bootstrap role gagal' }

Write-Host '=== 2. Buat database randesk_dev dan randesk_test (owner randesk_migrate) ==='
foreach ($db in 'randesk_dev', 'randesk_test') {
  $exists = & $psql -U postgres -h 127.0.0.1 -tAc "SELECT 1 FROM pg_database WHERE datname='$db'"
  if ($exists -ne '1') {
    & $psql -U postgres -h 127.0.0.1 -c "CREATE DATABASE $db OWNER randesk_migrate"
    Write-Host "DB $db dibuat"
  } else { Write-Host "DB $db sudah ada" }
}

$env:PGPASSWORD = 'randesk_migrate_dev'
$devUrl = 'postgres://randesk_migrate:randesk_migrate_dev@127.0.0.1:5432/randesk_dev?sslmode=disable'
$testUrl = 'postgres://randesk_migrate:randesk_migrate_dev@127.0.0.1:5432/randesk_test?sslmode=disable'

Write-Host '=== 3. migrate up (randesk_dev, dari DB kosong) ==='
& $migrate -path backend/migrations -database $devUrl up
if ($LASTEXITCODE -ne 0) { throw 'migrate up dev gagal' }
& $migrate -path backend/migrations -database $devUrl version

Write-Host '=== 4. Verifikasi jumlah tabel + indeks + constraint ==='
Write-Host ("tabel (harus 11 = 10 SCHEMA + schema_migrations): " + (& $psql -U randesk_migrate -h 127.0.0.1 -d randesk_dev -tAc "SELECT count(*) FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE'"))
& $psql -U randesk_migrate -h 127.0.0.1 -d randesk_dev -tAc "SELECT string_agg(table_name, ', ' ORDER BY table_name) FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE'"
Write-Host ("jumlah indeks: " + (& $psql -U randesk_migrate -h 127.0.0.1 -d randesk_dev -tAc "SELECT count(*) FROM pg_indexes WHERE schemaname='public'"))
Write-Host ("jumlah CHECK constraint: " + (& $psql -U randesk_migrate -h 127.0.0.1 -d randesk_dev -tAc "SELECT count(*) FROM pg_constraint WHERE connamespace='public'::regnamespace AND contype='c'"))

Write-Host '=== 5. Uji CHECK constraint (4 INSERT harus DITOLAK, 1 kontrol positif harus BERHASIL) ==='
$out = & $psql -U randesk_migrate -h 127.0.0.1 -d randesk_dev -f backend/scripts/check-constraints-task003.sql 2>&1 | Out-String
Write-Host $out
$expected = 'tickets_lifecycle_ck', 'tickets_version_check', 'users_email_normalized_ck', 'tickets_description_check'
$fail = $false
foreach ($e in $expected) {
  if ($out -notmatch [regex]::Escape($e)) { Write-Host "!!! constraint $e TIDAK muncul di output"; $fail = $true }
  else { Write-Host "OK ditolak oleh: $e" }
}
if ($out -match 'INSERT 0 1') { Write-Host 'OK kontrol positif: INSERT 0 1' } else { Write-Host '!!! kontrol positif tidak berhasil'; $fail = $true }
if ($fail) { throw 'verifikasi CHECK constraint gagal' }

Write-Host '=== 6. Privilege identity/sequence (catatan SCHEMA bagian 5) ==='
& $psql -U randesk_migrate -h 127.0.0.1 -d randesk_dev -tAc "SELECT table_name || '.' || column_name || ' identity=' || identity_generation FROM information_schema.columns WHERE table_schema='public' AND is_identity='YES'"

Write-Host '=== 7. migrate down 1 lalu up kembali (harus bersih) ==='
& $migrate -path backend/migrations -database $devUrl down 1
if ($LASTEXITCODE -ne 0) { throw 'migrate down gagal' }
$sisa = & $psql -U randesk_migrate -h 127.0.0.1 -d randesk_dev -tAc "SELECT string_agg(table_name, ',') FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE'"
Write-Host "tabel tersisa setelah down: '$sisa' (harus hanya schema_migrations - tabel bookkeeping golang-migrate)"
if ($sisa -ne 'schema_migrations') { throw 'down migration tidak bersih' }
& $migrate -path backend/migrations -database $devUrl up
if ($LASTEXITCODE -ne 0) { throw 're-up gagal' }
Write-Host 're-up OK'

Write-Host '=== 8. migrate up randesk_test ==='
& $migrate -path backend/migrations -database $testUrl up
if ($LASTEXITCODE -ne 0) { throw 'migrate up test gagal' }
& $migrate -path backend/migrations -database $testUrl version

Write-Host '=== SEMUA LANGKAH TASK-003 LULUS ==='
