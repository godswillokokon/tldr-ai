import React from 'react';
import {StyleSheet, Text, View} from 'react-native';

import type {UsageSnapshot} from '../domain/usageTypes';
import {appTheme as t} from '../theme/appTheme';

type Props = {
  usage: UsageSnapshot | null;
  error: string | null;
};

function fmtUsd(n: number): string {
  return `$${n.toFixed(2)}`;
}

export function UsageBar({usage, error}: Props) {
  if (error) {
    return (
      <View style={styles.wrap}>
        <Text style={styles.label}>Usage</Text>
        <Text style={styles.err}>{error}</Text>
      </View>
    );
  }
  if (!usage) {
    return (
      <View style={styles.wrap}>
        <Text style={styles.label}>Usage</Text>
        <Text style={styles.muted}>Loading…</Text>
      </View>
    );
  }

  const month = usage.billingMonth ? ` · ${usage.billingMonth}` : '';

  if (usage.usdCapActive && usage.budgetUsd > 0) {
    const pct = Math.min(
      100,
      Math.max(0, (usage.spentUsd / usage.budgetUsd) * 100),
    );
    return (
      <View style={styles.wrap}>
        <Text style={styles.label}>Usage{month}</Text>
        <Text style={styles.title}>
          {fmtUsd(usage.spentUsd)} / {fmtUsd(usage.budgetUsd)} ·{' '}
          {fmtUsd(usage.remainingUsd)} left
        </Text>
        <View style={styles.track}>
          <View style={[styles.fill, {width: `${pct}%`}]} />
        </View>
        <Text style={styles.muted}>
          Estimated cost per run: ~{fmtUsd(usage.estimateUsdPerCall)}
        </Text>
        {usage.callCapActive && usage.cap > 0 ? (
          <Text style={styles.muted}>
            Process limit: {usage.cap} calls ({usage.remaining} remaining).
          </Text>
        ) : null}
      </View>
    );
  }

  if (usage.callCapActive && usage.cap > 0) {
    const pct = Math.min(100, Math.max(0, (usage.used / usage.cap) * 100));
    return (
      <View style={styles.wrap}>
        <Text style={styles.label}>Usage</Text>
        <Text style={styles.title}>
          Calls · {usage.used} / {usage.cap} ({usage.remaining} left)
        </Text>
        <View style={styles.track}>
          <View style={[styles.fill, {width: `${pct}%`}]} />
        </View>
      </View>
    );
  }

  return (
    <View style={styles.wrap}>
      <Text style={styles.label}>Usage{month}</Text>
      <Text style={styles.title}>
        No cap · ~{fmtUsd(usage.spentUsd)} est. this month
      </Text>
      <Text style={styles.muted}>
        Set USAGE_BUDGET_USD on the server for a monthly ceiling.
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  wrap: {
    backgroundColor: t.surface,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: t.border,
    padding: 14,
    gap: 6,
    opacity: 0.94,
  },
  label: {
    fontSize: 11,
    fontWeight: '700',
    letterSpacing: 1.6,
    color: t.textMuted,
    textTransform: 'uppercase',
  },
  title: {
    fontSize: 14,
    fontWeight: '500',
    color: t.text,
    lineHeight: 20,
  },
  muted: {
    fontSize: 12,
    color: t.textSecondary,
    lineHeight: 17,
  },
  err: {
    fontSize: 13,
    color: t.dangerMuted,
    lineHeight: 18,
  },
  track: {
    height: 6,
    borderRadius: 3,
    backgroundColor: t.border,
    overflow: 'hidden',
  },
  fill: {
    height: '100%',
    backgroundColor: t.accentLine,
    borderRadius: 3,
  },
});
