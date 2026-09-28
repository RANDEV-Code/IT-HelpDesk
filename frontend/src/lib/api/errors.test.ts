import { describe, expect, it } from 'vitest';

import { ApiError, ErrorCodes, fromResponse } from './errors';

function jsonResponse(status: number, body: unknown, headers: Record<string, string> = {}): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json', ...headers },
  });
}

describe('fromResponse (API_SPEC §3)', () => {
  it('422 VALIDATION_ERROR terpetakan ke map field yang bisa dirender form', async () => {
    const res = jsonResponse(
      422,
      {
        error: { code: 'VALIDATION_ERROR', message: 'Periksa data.', fields: { title: 'minimal 5 karakter' } },
        meta: { request_id: 'req-422' },
      },
      { 'X-Request-ID': 'req-422' },
    );
    const err = await fromResponse(res);
    expect(err).toBeInstanceOf(ApiError);
    expect(err.code).toBe(ErrorCodes.VALIDATION_ERROR);
    expect(err.status).toBe(422);
    expect(err.isValidation).toBe(true);
    expect(err.fields).toEqual({ title: 'minimal 5 karakter' });
    expect(err.requestId).toBe('req-422');
  });

  it('409 VERSION_CONFLICT terpetakan dengan current_version', async () => {
    const res = jsonResponse(409, {
      error: { code: 'VERSION_CONFLICT', message: 'Data berubah.', current_version: 7 },
      meta: { request_id: 'req-409' },
    });
    const err = await fromResponse(res);
    expect(err.code).toBe(ErrorCodes.VERSION_CONFLICT);
    expect(err.isConflict).toBe(true);
    expect(err.currentVersion).toBe(7);
  });

  it('401 UNAUTHENTICATED memicu penanda logout', async () => {
    const res = jsonResponse(401, { error: { code: 'UNAUTHENTICATED', message: 'Sesi habis.' }, meta: { request_id: 'r' } });
    const err = await fromResponse(res);
    expect(err.isUnauthorized).toBe(true);
  });

  it('429 memuat Retry-After detik', async () => {
    const res = new Response(
      JSON.stringify({ error: { code: 'RATE_LIMITED', message: 'Terlalu sering.' }, meta: { request_id: 'r' } }),
      { status: 429, headers: { 'Content-Type': 'application/json', 'Retry-After': '30' } },
    );
    const err = await fromResponse(res);
    expect(err.isRateLimited).toBe(true);
    expect(err.retryAfter).toBe(30);
  });

  it('body non-JSON tetap menghasilkan error generik ber-status', async () => {
    const res = new Response('gateway bad', { status: 502 });
    const err = await fromResponse(res);
    expect(err.status).toBe(502);
    expect(err.code).toBe(ErrorCodes.INTERNAL_ERROR);
    expect(err.message).toContain('server');
  });
});
