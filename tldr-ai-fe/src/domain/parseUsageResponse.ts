import type {UsageSnapshot} from './usageTypes';

type Raw = Record<string, unknown>;

function num(v: unknown, fallback = 0): number {
  return typeof v === 'number' && Number.isFinite(v) ? v : fallback;
}

function bool(v: unknown, fallback = false): boolean {
  return typeof v === 'boolean' ? v : fallback;
}

function str(v: unknown, fallback = ''): string {
  return typeof v === 'string' ? v : fallback;
}

export function parseUsageResponse(rawBody: string): UsageSnapshot {
  let data: unknown;
  try {
    data = rawBody.trim() === '' ? {} : JSON.parse(rawBody);
  } catch {
    throw new Error('Invalid JSON from usage endpoint');
  }
  const o = data as Raw;

  return {
    used: num(o.used),
    cap: num(o.cap),
    remaining: num(o.remaining),
    unlimited: bool(o.unlimited),
    spentUsd: num(o.spentUsd),
    budgetUsd: num(o.budgetUsd),
    remainingUsd: num(o.remainingUsd),
    usdCapActive: bool(o.usdCapActive),
    callCapActive: bool(o.callCapActive),
    estimateUsdPerCall: num(o.estimateUsdPerCall, 0.02),
    billingMonth: str(o.billingMonth),
  };
}
