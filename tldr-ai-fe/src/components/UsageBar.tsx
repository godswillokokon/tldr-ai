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
        <Text style={styles.label}>This month</Text>
        <Text style={styles.err}>{error}</Text>
      </View>
    );
  }
  if (!usage) {
    return (
      <View style={styles.wrap}>
        <Text style={styles.label}>This month</Text>
        <Text style={styles.muted}>Loading…</Text>
      </View>
    );
  }

  const monthBit = usage.billingMonth ? ` · ${usage.billingMonth}` : '';

  if (usage.usdCapActive && usage.budgetUsd > 0) {
    const pct = Math.min(
      100,
      Math.max(0, (usage.spentUsd / usage.budgetUsd) * 100),
    );
    return (
      <View style={styles.wrap}>
        <Text style={styles.label}>This month{monthBit}</Text>
        <Text style={styles.hero}>{fmtUsd(usage.spentUsd)}</Text>
        <Text style={styles.sub}>Updates after each successful summarize</Text>
        <Text style={styles.meta}>
          {fmtUsd(usage.budgetUsd)} cap · {fmtUsd(usage.remainingUsd)} left
        </Text>
        <View style={styles.track}>
          <View style={[styles.fill, {width: `${pct}%`}]} />
        </View>
        {usage.callCapActive && usage.cap > 0 ? (
          <Text style={styles.meta}>
            Runs {usage.used} / {usage.cap} · {usage.remaining} left
          </Text>
        ) : null}
      </View>
    );
  }

  if (usage.callCapActive && usage.cap > 0) {
    const pct = Math.min(100, Math.max(0, (usage.used / usage.cap) * 100));
    return (
      <View style={styles.wrap}>
        <Text style={styles.label}>This month{monthBit}</Text>
        <Text style={styles.hero}>{fmtUsd(usage.spentUsd)}</Text>
        <Text style={styles.sub}>Updates after each successful summarize</Text>
        <Text style={styles.meta}>
          Runs {usage.used} / {usage.cap} · {usage.remaining} left
        </Text>
        <View style={styles.track}>
          <View style={[styles.fill, {width: `${pct}%`}]} />
        </View>
      </View>
    );
  }

  return (
    <View style={styles.wrap}>
      <Text style={styles.label}>This month{monthBit}</Text>
      <Text style={styles.hero}>{fmtUsd(usage.spentUsd)}</Text>
      <Text style={styles.sub}>Updates after each successful summarize</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  wrap: {
    backgroundColor: t.surface,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: t.border,
    paddingVertical: 10,
    paddingHorizontal: 12,
    gap: 1,
    opacity: 0.94,
  },
  label: {
    fontSize: 10,
    fontWeight: '700',
    letterSpacing: 1.2,
    color: t.textMuted,
    textTransform: 'uppercase',
    marginBottom: 1,
  },
  hero: {
    fontSize: 20,
    fontWeight: '700',
    color: t.text,
    letterSpacing: -0.3,
    lineHeight: 24,
  },
  sub: {
    fontSize: 11,
    color: t.textSecondary,
    lineHeight: 15,
    marginTop: 1,
  },
  meta: {
    fontSize: 11,
    color: t.textSecondary,
    lineHeight: 15,
    marginTop: 3,
  },
  muted: {
    fontSize: 12,
    color: t.textSecondary,
    lineHeight: 16,
  },
  err: {
    fontSize: 12,
    color: t.dangerMuted,
    lineHeight: 16,
    marginTop: 2,
  },
  track: {
    height: 5,
    borderRadius: 3,
    backgroundColor: t.border,
    overflow: 'hidden',
    marginTop: 6,
  },
  fill: {
    height: '100%',
    backgroundColor: t.accentLine,
    borderRadius: 3,
  },
});
