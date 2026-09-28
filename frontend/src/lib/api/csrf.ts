// Penyimpanan CSRF token HANYA di memori modul (SECURITY §4 / ARCHITECTURE §7).
// Token TIDAK pernah masuk localStorage/sessionStorage/cookie yang dibuat JS.

let csrfToken: string | null = null;

export const csrf = {
  /** Set token dari respons /auth/login atau /auth/me. */
  set(token: string | null): void {
    csrfToken = token;
  },
  /** Ambil token untuk header mutasi (null bila belum bersesi). */
  get(): string | null {
    return csrfToken;
  },
  /** Kosongkan saat logout / sesi berakhir. */
  clear(): void {
    csrfToken = null;
  },
} as const;
