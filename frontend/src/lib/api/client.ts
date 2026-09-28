// Klien HTTP tunggal frontend (TASK-006). Satu jalur untuk mem-parse envelope,
// memetakan error, mengirim X-CSRF-Token pada mutasi, dan timeout sisi-klien.
// Kebijakan retry diekspor untuk dipakai QueryClient (read terbatas, mutasi nol).

import { csrf } from './csrf';
import { ApiError, ErrorCodes, fromResponse } from './errors';
import type { ApiEnvelope, ApiMeta } from './types';

const BASE = '/api/v1';
const DEFAULT_TIMEOUT_MS = 15_000;

export type HttpMethod = 'GET' | 'HEAD' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';

export interface RequestOptions {
  method?: HttpMethod;
  body?: unknown;
  query?: Record<string, string | number | boolean | undefined | null>;
  signal?: AbortSignal;
  timeoutMs?: number;
}

export function isMutation(method: HttpMethod): boolean {
  return method !== 'GET' && method !== 'HEAD';
}

function buildUrl(path: string, query?: RequestOptions['query']): string {
  const url = `${BASE}${path}`;
  if (!query) return url;
  const params = new URLSearchParams();
  for (const [k, v] of Object.entries(query)) {
    if (v !== undefined && v !== null) params.append(k, String(v));
  }
  const qs = params.toString();
  return qs ? `${url}?${qs}` : url;
}

/** Permintaan JSON generik. Mengembalikan envelope { data, meta }. */
export async function request<T>(path: string, opts: RequestOptions = {}): Promise<ApiEnvelope<T>> {
  const method = opts.method ?? 'GET';
  const mutating = isMutation(method);

  const headers: Record<string, string> = { Accept: 'application/json' };
  if (opts.body !== undefined) headers['Content-Type'] = 'application/json';
  if (mutating) {
    const token = csrf.get();
    if (token) headers['X-CSRF-Token'] = token;
  }

  const controller = new AbortController();
  const timeoutMs = opts.timeoutMs ?? DEFAULT_TIMEOUT_MS;
  const timer = setTimeout(() => controller.abort(), timeoutMs);
  const onOuterAbort = () => controller.abort();
  opts.signal?.addEventListener('abort', onOuterAbort);

  let res: Response;
  try {
    res = await fetch(buildUrl(path, opts.query), {
      method,
      headers,
      credentials: 'same-origin',
      body: opts.body !== undefined ? JSON.stringify(opts.body) : undefined,
      signal: controller.signal,
    });
  } catch (err) {
    const aborted = err instanceof DOMException ? err.name === 'AbortError' : controller.signal.aborted;
    throw new ApiError({
      code: aborted ? ErrorCodes.TIMEOUT : ErrorCodes.NETWORK_ERROR,
      status: 0,
      message: aborted ? 'Waktu permintaan melebihi batas.' : 'Koneksi ke server gagal.',
    });
  } finally {
    clearTimeout(timer);
    opts.signal?.removeEventListener('abort', onOuterAbort);
  }

  if (!res.ok) {
    throw await fromResponse(res);
  }

  const requestId = res.headers.get('X-Request-ID') ?? '';
  if (res.status === 204) {
    return { data: undefined as T, meta: { request_id: requestId } };
  }

  const json = (await res.json()) as Partial<ApiEnvelope<T>>;
  const meta: ApiMeta = json?.meta ?? { request_id: requestId };
  return { data: json.data as T, meta };
}

export const api = {
  get: <T>(path: string, query?: RequestOptions['query'], signal?: AbortSignal) =>
    request<T>(path, { method: 'GET', query, signal }),
  post: <T>(path: string, body?: unknown) => request<T>(path, { method: 'POST', body }),
  put: <T>(path: string, body?: unknown) => request<T>(path, { method: 'PUT', body }),
  patch: <T>(path: string, body?: unknown) => request<T>(path, { method: 'PATCH', body }),
  del: <T>(path: string, body?: unknown) => request<T>(path, { method: 'DELETE', body }),
};

/**
 * Kebijakan retry TanStack Query (API_SPEC §10):
 * - READ: boleh diulang terbatas untuk kegagalan transien (jaringan/timeout,
 *   5xx, 408, 429). Maks 2 percobaan ulang.
 * - Jangan pernah retry error 4xx non-transien.
 */
export function queryRetryPolicy(failureCount: number, error: unknown): boolean {
  if (failureCount >= 2) return false;
  if (!(error instanceof ApiError)) return true; // error tak dikenal: biarkan 1x retry
  if (error.code === ErrorCodes.NETWORK_ERROR || error.code === ErrorCodes.TIMEOUT) return true;
  if (error.status === 408 || error.status === 429) return true;
  if (error.status >= 500) return true;
  return false;
}

/** Mutasi TIDAK di-retry otomatis (create/comment/upload mungkin sudah commit). */
export const mutationRetryPolicy = false;

/** Konvensi query key menyertakan user id (ARCHITECTURE §7). */
export function queryKeys(userId: string | undefined, ...parts: Array<string | number>) {
  return ['u', userId ?? 'anon', ...parts] as const;
}
