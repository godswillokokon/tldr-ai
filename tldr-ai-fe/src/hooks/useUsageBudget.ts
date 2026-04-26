import {useCallback, useEffect, useState} from 'react';

import {fetchUsage} from '../data/usageRepository';
import type {UsageSnapshot} from '../domain/usageTypes';

const POLL_MS = 8000;

export function useUsageBudget() {
  const [usage, setUsage] = useState<UsageSnapshot | null>(null);
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    try {
      const next = await fetchUsage();
      setUsage(next);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Usage unavailable');
    }
  }, []);

  useEffect(() => {
    refresh().catch(() => {});
    const id = setInterval(() => {
      refresh().catch(() => {});
    }, POLL_MS);
    return () => clearInterval(id);
  }, [refresh]);

  return {usage, error, refresh};
}
