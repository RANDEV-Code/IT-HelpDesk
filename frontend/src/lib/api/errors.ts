// Pemetaan error API_SPEC §3 ke tipe UI (field errors, conflict, unauthorized,
// rate limit). Klien memanggil fromResponse() untuk non-2xx.

import type { ApiErrorEnvelope } from './types';

/** Kode error yang dikenal. STRING bebas tetap diizinkan (future-proof). */
export const ErrorCodes = {
  MALFORMED_REQUEST: 'MALFORMED_REQUEST',
  UNAUTHENTICATED: 'UNAUTHENTICATED',
  FORBIDDEN: 'FORBIDDEN',
  CSRF_FAILED: 'CSRF_FAILED',
  PASSWORD_CHANGE_REQUIRED: 'PASSWORD_CHANGE_REQUIRED',
  NOT_FOUND: 'NOT_FOUND',
  VERSION_CONFLICT: 'VERSION_CONFLICT',
  ALREADY_ASSIGNED: 'ALREADY_ASSIGNED',
  INVALID_STATE: 'INVALID_STATE',
  ACTIVE_TICKETS_EXIST: 'ACTIVE_TICKETS_EXIST',
  LAST_ADMIN: 'LAST_ADMIN',
  EMAIL_EXISTS: 'EMAIL_EXISTS',
  NAME_EXISTS: 'NAME_EXISTS',
  PAYLOAD_TOO_LARGE: 'PAYLOAD_TOO_LARGE',
  UNSUPPORTED_MEDIA_TYPE: 'UNSUPPORTED_MEDIA_TYPE',
  VALIDATION_ERROR: 'VALIDATION_ERROR',
  MASTER_INACTIVE: 'MASTER_INACTIVE',
  NO_CHANGE: 'NO_CHANGE',
  ATTACHMENT_LIMIT: 'ATTACHMENT_LIMIT',
  RATE_LIMITED: 'RATE_LIMITED',
  INTERNAL_ERROR: 'INTERNAL_ERROR',
  SERVICE_UNAVAILABLE: 'SERVICE_UNAVAILABLE',
  // kode sintetis sisi-klien (bukan dari server)
  NETWORK_ERROR: 'NETWORK_ERROR',
  TIMEOUT: 'TIMEOUT',
} as const;

export class ApiError extends Error {
  readonly code: string;
  readonly status: number;
  readonly fields?: Record<string, string>;
  readonly currentVersion?: number;
  readonly requestId?: string;
  readonly retryAfter?: number; // detik (RATE_LIMITED)

  constructor(init: {
    code: string;
    status: number;
    message: string;
    fields?: Record<string, string>;
    currentVersion?: number;
    requestId?: string;
    retryAfter?: number;
  }) {
    super(init.message);
    this.name = 'ApiError';
    this.code = init.code;
    this.status = init.status;
    this.fields = init.fields;
    this.currentVersion = init.currentVersion;
    this.requestId = init.requestId;
    this.retryAfter = init.retryAfter;
  }

  get isValidation(): boolean {
    return this.code === ErrorCodes.VALIDATION_ERROR;
  }
  get isConflict(): boolean {
    return this.status === 409;
  }
  get isUnauthorized(): boolean {
    return this.code === ErrorCodes.UNAUTHENTICATED;
  }
  get isForbidden(): boolean {
    return this.status === 403;
  }
  get isNotFound(): boolean {
    return this.code === ErrorCodes.NOT_FOUND;
  }
  get isRateLimited(): boolean {
    return this.code === ErrorCodes.RATE_LIMITED;
  }
}

/** Membangun ApiError dari Response non-2xx, mem-parse envelope error + header. */
export async function fromResponse(res: Response): Promise<ApiError> {
  const requestId = res.headers.get('X-Request-ID') ?? undefined;
  const retryAfterHeader = res.headers.get('Retry-After');
  const retryAfterNum = retryAfterHeader ? Number(retryAfterHeader) : undefined;

  let body: Partial<ApiErrorEnvelope['error']> | undefined;
  let metaRequestId: string | undefined;
  try {
    const parsed = (await res.json()) as Partial<ApiErrorEnvelope>;
    if (parsed && typeof parsed === 'object' && parsed.error) {
      body = parsed.error;
      metaRequestId = parsed.meta?.request_id;
    }
  } catch {
    // body bukan JSON — fallback ke generik
  }

  return new ApiError({
    code: body?.code ?? codeForStatus(res.status),
    status: res.status,
    message: body?.message ?? defaultMessageFor(res.status),
    fields: body?.fields,
    currentVersion: body?.current_version,
    requestId: metaRequestId ?? requestId,
    retryAfter: Number.isFinite(retryAfterNum) ? retryAfterNum : undefined,
  });
}

function codeForStatus(status: number): string {
  switch (status) {
    case 400:
      return ErrorCodes.MALFORMED_REQUEST;
    case 401:
      return ErrorCodes.UNAUTHENTICATED;
    case 403:
      return ErrorCodes.FORBIDDEN;
    case 404:
      return ErrorCodes.NOT_FOUND;
    case 409:
      return ErrorCodes.INVALID_STATE;
    case 413:
      return ErrorCodes.PAYLOAD_TOO_LARGE;
    case 415:
      return ErrorCodes.UNSUPPORTED_MEDIA_TYPE;
    case 422:
      return ErrorCodes.VALIDATION_ERROR;
    case 429:
      return ErrorCodes.RATE_LIMITED;
    case 500:
      return ErrorCodes.INTERNAL_ERROR;
    case 503:
      return ErrorCodes.SERVICE_UNAVAILABLE;
    default:
      if (status >= 500) return ErrorCodes.INTERNAL_ERROR;
      return 'UNEXPECTED_ERROR';
  }
}

function defaultMessageFor(status: number): string {
  if (status >= 500) return 'Terjadi kesalahan pada server.';
  if (status === 429) return 'Terlalu banyak permintaan; coba lagi nanti.';
  if (status === 401) return 'Sesi tidak valid atau telah berakhir.';
  return 'Permintaan tidak dapat diproses.';
}
