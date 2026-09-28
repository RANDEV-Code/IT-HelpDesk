/** Kondisi sistem: loading / empty / error (DESIGN §7, §8). */

export function LoadingState({ message = 'Memuat data…' }: { message?: string }) {
  return (
    <div role="status" aria-live="polite" className="flex items-center gap-3 p-4 text-muted">
      <span
        aria-hidden="true"
        className="inline-block size-5 animate-spin rounded-full border-2 border-line border-t-primary"
      />
      <span className="text-body">{message}</span>
    </div>
  );
}

export function EmptyState({
  title,
  description,
}: {
  title: string;
  description?: string;
}) {
  return (
    <div className="flex flex-col items-start gap-1 rounded-card border border-line bg-surface p-6">
      <p className="text-section font-semibold text-ink">{title}</p>
      {description && <p className="text-label text-muted">{description}</p>}
    </div>
  );
}

export function ErrorState({
  message = 'Terjadi kesalahan. Coba lagi.',
  onRetry,
}: {
  message?: string;
  onRetry?: () => void;
}) {
  return (
    <div role="alert" className="flex flex-col items-start gap-3 rounded-card border border-danger bg-surface p-6">
      <p className="text-section font-semibold text-danger">{message}</p>
      {onRetry && (
        <button
          type="button"
          onClick={onRetry}
          className="min-h-11 rounded-control border border-control bg-surface px-4 text-label font-semibold text-ink"
        >
          Coba lagi
        </button>
      )}
    </div>
  );
}
