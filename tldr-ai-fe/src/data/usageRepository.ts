import {getUsageUrl} from '../config/apiConfig';
import {parseUsageResponse} from '../domain/parseUsageResponse';
import type {UsageSnapshot} from '../domain/usageTypes';

export async function fetchUsage(): Promise<UsageSnapshot> {
  let response: Response;
  try {
    response = await fetch(getUsageUrl());
  } catch {
    throw new Error('Could not reach server for usage');
  }
  const raw = await response.text();
  if (!response.ok) {
    throw new Error(`Usage request failed (${response.status})`);
  }
  return parseUsageResponse(raw);
}
