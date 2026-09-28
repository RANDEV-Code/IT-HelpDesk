/** Halaman sementara untuk route yang belum memiliki fitur (TASK-005). */
export function PlaceholderPage({ title, note }: { title: string; note?: string }) {
  return (
    <section>
      <h1 className="text-page-title-sm font-bold text-ink md:text-page-title">{title}</h1>
      <p className="mt-3 rounded-card border border-line bg-surface p-6 text-body text-muted">
        {note ?? 'Bagian ini akan diisi pada task berikutnya.'}
      </p>
    </section>
  );
}
