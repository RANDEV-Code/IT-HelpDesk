import { useQuery } from '@tanstack/react-query';

import { ErrorState, LoadingState } from '../components/StateViews';
import { formatWita } from '../lib/date';

// Demo TASK-005: memanggil /health/live lewat proxy Vite ke backend Go,
// membuktikan satu-origin SPA<->API (ARCHITECTURE §7). Juga menampilkan contoh
// timestamp UTC yang dirender WITA (BR-13).
async function fetchLive(): Promise<{ status: string }> {
  const res = await fetch('/health/live');
  if (!res.ok) {
    throw new Error(`health/live -> ${res.status}`);
  }
  return (await res.json()) as { status: string };
}

function StatusBadge({ label, tone }: { label: string; tone: 'ok' | 'bad' }) {
  const cls =
    tone === 'ok'
      ? 'bg-status-resolved-bg text-status-resolved'
      : 'bg-status-open-bg text-danger';
  return (
    <span className={`inline-flex items-center rounded-control px-3 py-1 text-label font-semibold ${cls}`}>
      {label}
    </span>
  );
}

export function DashboardPage() {
  const { isPending, isError, data, refetch } = useQuery({
    queryKey: ['health', 'live'],
    queryFn: fetchLive,
  });

  const contohUtc = '2026-09-26T04:00:00Z';

  return (
    <section className="flex flex-col gap-6">
      <header>
        <h1 className="text-page-title-sm font-bold text-ink md:text-page-title">Dashboard</h1>
        <p className="mt-1 text-label text-muted">
          Kerangka aplikasi (TASK-005). Data layanan akan muncul setelah backend terhubung.
        </p>
      </header>

      <div className="rounded-card border border-line bg-surface p-6">
        <h2 className="text-section font-semibold text-ink">Status koneksi backend (via proxy /api, /health)</h2>
        <div className="mt-3">
          {isPending && <LoadingState message="Memeriksa /health/live…" />}
          {isError && (
            <ErrorState
              message="Backend tidak terjangkau. Jalankan `go run ./cmd/api` di port 8080."
              onRetry={() => void refetch()}
            />
          )}
          {!isPending && !isError && (
            <StatusBadge label={`Backend hidup: ${data?.status ?? 'ok'}`} tone="ok" />
          )}
        </div>
      </div>

      <div className="rounded-card border border-line bg-surface p-6">
        <h2 className="text-section font-semibold text-ink">Contoh format waktu WITA (BR-13)</h2>
        <p className="mt-2 text-body text-muted">
          Timestamp UTC <span className="font-mono tabular-nums">{contohUtc}</span> ditampilkan sebagai{' '}
          <span className="font-semibold text-ink">{formatWita(contohUtc)}</span>.
        </p>
      </div>
    </section>
  );
}
