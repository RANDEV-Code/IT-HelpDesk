import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { api, isMutation, mutationRetryPolicy, queryRetryPolicy, request } from './client';
import { csrf } from './csrf';
import { ApiError, ErrorCodes } from './errors';

interface FetchCall {
  url: string;
  init: RequestInit;
}

let lastCall: FetchCall | undefined;
function mockFetch(response: Response) {
  lastCall = undefined;
  const fn = vi.fn(async (url: string, init: RequestInit) => {
    lastCall = { url, init };
    return response;
  });
  vi.stubGlobal('fetch', fn);
  return fn;
}

function jsonRes(status: number, body: unknown, headers: Record<string, string> = {}) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json', ...headers },
  });
}

function headersOf(init: RequestInit): Record<string, string> {
  return (init.headers ?? {}) as Record<string, string>;
}

// Penanda localStorage: bila ada kode menulis token ke storage, tes gagal.
let storageWrites: string[];
beforeEach(() => {
  storageWrites = [];
  vi.stubGlobal('localStorage', {
    getItem: () => null,
    setItem: (k: string) => storageWrites.push(k),
    removeItem: () => {},
    clear: () => {},
    key: () => null,
    length: 0,
  } as unknown as Storage);
});

afterEach(() => {
  vi.unstubAllGlobals();
  csrf.clear();
  lastCall = undefined;
});

describe('client.request', () => {
  it('GET sukses mengembalikan data + meta, tanpa CSRF header', async () => {
    mockFetch(jsonRes(200, { data: { id: '1' }, meta: { request_id: 'r1' } }, { 'X-Request-ID': 'r1' }));
    csrf.set('tok'); // token ada, tapi GET tidak boleh mengirimnya

    const env = await request<{ id: string }>('/tickets/1');
    expect(env.data).toEqual({ id: '1' });
    expect(env.meta.request_id).toBe('r1');
    expect(lastCall?.url).toBe('/api/v1/tickets/1');
    expect(lastCall?.init.method).toBe('GET');
    expect(lastCall?.init.credentials).toBe('same-origin');
    expect(headersOf(lastCall!.init)['X-CSRF-Token']).toBeUndefined();
  });

  it('POST mutasi mengirim X-CSRF-Token saat token tersedia', async () => {
    mockFetch(jsonRes(201, { data: { ok: true }, meta: { request_id: 'r2' } }));
    csrf.set('abc123');

    await api.post('/tickets', { title: 'Contoh' });
    const h = headersOf(lastCall!.init);
    expect(lastCall?.init.method).toBe('POST');
    expect(h['X-CSRF-Token']).toBe('abc123');
    expect(h['Content-Type']).toBe('application/json');
    expect(lastCall?.init.body).toBe(JSON.stringify({ title: 'Contoh' }));
  });

  it('PUT/PATCH/DELETE juga mengirim CSRF (mutasi)', async () => {
    for (const method of ['PUT', 'PATCH', 'DELETE'] as const) {
      csrf.clear();
      csrf.set('t-' + method);
      mockFetch(jsonRes(200, { data: null, meta: { request_id: 'r' } }));
      await request('/x', { method, body: {} });
      expect(headersOf(lastCall!.init)['X-CSRF-Token']).toBe('t-' + method);
      expect(isMutation(method)).toBe(true);
    }
  });

  it('non-2xx melempar ApiError terpetakan', async () => {
    mockFetch(jsonRes(422, { error: { code: 'VALIDATION_ERROR', message: 'x', fields: { a: 'b' } }, meta: { request_id: 'r' } }));
    await expect(request('/y')).rejects.toMatchObject({ code: ErrorCodes.VALIDATION_ERROR, status: 422 });
  });

  it('reject fetch -> ApiError NETWORK_ERROR (status 0)', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => {
      throw new TypeError('failed to fetch');
    }));
    const p = request('/z').catch((e: unknown) => e as ApiError);
    const err = await p;
    expect(err).toBeInstanceOf(ApiError);
    expect((err as ApiError).code).toBe(ErrorCodes.NETWORK_ERROR);
    expect((err as ApiError).status).toBe(0);
  });

  it('204 No Content mengembalikan data undefined', async () => {
    mockFetch(new Response(null, { status: 204 }));
    const env = await request<undefined>('/logout', { method: 'POST', body: {} });
    expect(env.data).toBeUndefined();
  });

  it('CSRF token TIDAK pernah ditulis ke localStorage', async () => {
    mockFetch(jsonRes(200, { data: {}, meta: { request_id: 'r' } }));
    csrf.set('secret-token');
    await api.post('/tickets', { title: 'q' });
    expect(storageWrites).toEqual([]);
  });
});

describe('retry policy (API_SPEC §10)', () => {
  const apiErr = (status: number, code = 'X') => new ApiError({ code, status, message: 'm' });

  it('read: retry transien (5xx, 429, 408, network, timeout) hingga maks 2', () => {
    expect(queryRetryPolicy(0, apiErr(500))).toBe(true);
    expect(queryRetryPolicy(0, apiErr(429, ErrorCodes.RATE_LIMITED))).toBe(true);
    expect(queryRetryPolicy(0, apiErr(0, ErrorCodes.NETWORK_ERROR))).toBe(true);
    expect(queryRetryPolicy(1, apiErr(503))).toBe(true);
    expect(queryRetryPolicy(2, apiErr(500))).toBe(false); // batas 2 percobaan ulang
  });

  it('read: TIDAK retry 4xx non-transien', () => {
    expect(queryRetryPolicy(0, apiErr(404, ErrorCodes.NOT_FOUND))).toBe(false);
    expect(queryRetryPolicy(0, apiErr(422, ErrorCodes.VALIDATION_ERROR))).toBe(false);
    expect(queryRetryPolicy(0, apiErr(403, ErrorCodes.FORBIDDEN))).toBe(false);
  });

  it('read: error tak dikenal diberi kelonggaran 1x', () => {
    expect(queryRetryPolicy(0, new Error('unexpected'))).toBe(true);
    expect(queryRetryPolicy(2, new Error('unexpected'))).toBe(false);
  });

  it('mutasi: tidak pernah retry', () => {
    expect(mutationRetryPolicy).toBe(false);
  });
});
