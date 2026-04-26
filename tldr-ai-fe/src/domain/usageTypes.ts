/**
 * GET /api/usage response (subset validated in parseUsageResponse).
 */
export type UsageSnapshot = {
  used: number;
  cap: number;
  remaining: number;
  unlimited: boolean;
  spentUsd: number;
  budgetUsd: number;
  remainingUsd: number;
  usdCapActive: boolean;
  callCapActive: boolean;
  estimateUsdPerCall: number;
  billingMonth: string;
};
