$ErrorActionPreference = 'Continue'
Set-Location 'd:\laragon\www\ITHelpDesk\backend'

Write-Output '=== gofmt -l (harus kosong) ==='
gofmt -l .

Write-Output '=== go vet ./... ==='
go vet ./... 2>&1

Write-Output '=== go test ./... (testdb SKIP tanpa TEST_DATABASE_URL) ==='
$env:TEST_DATABASE_URL = $null
go test ./... 2>&1

Write-Output '=== tes integrasi NYATA terhadap randesk_test ==='
$env:TEST_DATABASE_URL = 'postgres://randesk_migrate:randesk_migrate_dev@127.0.0.1:5432/randesk_test?sslmode=disable'
go test ./internal/platform/testdb/ -run TestIntegrationMigrateAndCRUD -v 2>&1
