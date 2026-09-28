// Utilitas tanggal — BR-13: data disimpan UTC, ditampilkan Asia/Makassar (WITA).
// Menggunakan Intl API bawaan (tanpa dependensi eksternal).

const WITA_TIME_ZONE = 'Asia/Makassar';

const dateFormatter = new Intl.DateTimeFormat('id-ID', {
  timeZone: WITA_TIME_ZONE,
  day: '2-digit',
  month: 'short',
  year: 'numeric',
});

const timeFormatter = new Intl.DateTimeFormat('id-ID', {
  timeZone: WITA_TIME_ZONE,
  hour: '2-digit',
  minute: '2-digit',
});

/**
 * Format timestamp ISO (UTC) menjadi "26 Sep 2026 12:05 WITA".
 * Melempar Error bila input bukan timestamp valid.
 */
export function formatWita(iso: string): string {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) {
    throw new Error(`timestamp tidak valid: ${iso}`);
  }
  return `${dateFormatter.format(d)} ${timeFormatter.format(d)} WITA`;
}

/** Hanya tanggalnya: "26 Sep 2026" (WITA). */
export function formatWitaDate(iso: string): string {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) {
    throw new Error(`timestamp tidak valid: ${iso}`);
  }
  return dateFormatter.format(d);
}

/** Konversi Date/ISO menjadi string RFC 3339 UTC untuk dikirim ke API. */
export function toUtcIso(value: string | Date): string {
  const d = value instanceof Date ? value : new Date(value);
  if (Number.isNaN(d.getTime())) {
    throw new Error(`timestamp tidak valid: ${String(value)}`);
  }
  return d.toISOString();
}
