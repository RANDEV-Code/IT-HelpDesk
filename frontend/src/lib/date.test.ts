import { describe, expect, it } from 'vitest';

import { formatWita, toUtcIso } from './date';

describe('formatWita (BR-13: UTC -> Asia/Makassar, UTC+8)', () => {
  it('selalu melabeli hasil dengan WITA', () => {
    expect(formatWita('2026-09-26T04:00:00Z')).toMatch(/WITA$/);
  });

  it('menggeser hari saat melewati tengah malam WITA (UTC+8)', () => {
    // 2026-09-26T16:00Z = 2026-09-27 00:00 WITA
    expect(formatWita('2026-09-26T16:00:00Z')).toContain('27');
    // 2026-09-26T15:59Z = 2026-09-26 23:59 WITA
    expect(formatWita('2026-09-26T15:59:00Z')).toContain('26');
  });

  it('menolak timestamp tidak valid', () => {
    expect(() => formatWita('bukan-tanggal')).toThrow(/tidak valid/);
  });
});

describe('toUtcIso', () => {
  it('menghasilkan string RFC 3339 UTC berakhiran Z', () => {
    expect(toUtcIso('2026-09-26T04:00:00.000Z')).toBe('2026-09-26T04:00:00.000Z');
  });

  it('menolak nilai tak dapat diparse', () => {
    expect(() => toUtcIso('salah')).toThrow(/tidak valid/);
  });
});
