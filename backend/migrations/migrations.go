// Package migrations menyematkan seluruh file .sql agar dapat dijalankan
// dari binari tes/aplikasi tanpa bergantung pada lokasi filesystem saat runtime
// (dipakai helper testdb dan, nanti, perintah maintenance).
package migrations

import "embed"

// FS berisi file migrasi (000001_baseline.up.sql dst) pada direktori ini.
//
//go:embed *.sql
var FS embed.FS
