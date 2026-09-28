import { Link } from 'react-router-dom';

/**
 * "Tiket tidak ditemukan" — dipakai untuk URL objek yang tidak dapat dibaca
 * actor tanpa mengonfirmasi keberadaannya (DESIGN §3).
 */
export function NotFoundPage() {
  return (
    <section>
      <h1 className="text-page-title-sm font-bold text-ink md:text-page-title">
        Data tidak ditemukan
      </h1>
      <p className="mt-3 text-body text-muted">
        Halaman atau tiket yang Anda cari tidak tersedia.
      </p>
      <Link
        to="/dashboard"
        className="mt-4 inline-flex min-h-11 items-center rounded-control bg-primary px-4 text-label font-semibold text-white"
      >
        Kembali ke dashboard
      </Link>
    </section>
  );
}
