import React from 'react';
import {StyleSheet, Text, View} from 'react-native';

import type {ProcessResponse} from '../domain/types';
import {appTheme as t} from '../theme/appTheme';

type Props = {
  result: ProcessResponse;
};

export function ProcessResultCard({result}: Props) {
  return (
    <View style={styles.card}>
      <Text style={styles.sectionTitle}>Summary</Text>
      <Text style={styles.bodyText}>{result.summary}</Text>

      <Text style={styles.sectionTitle}>Next steps</Text>
      {result.actionItems.map((item, index) => (
        <View key={`${item}-${index}`} style={styles.actionRow}>
          <View style={styles.actionBadge}>
            <Text style={styles.actionIndex}>{index + 1}</Text>
          </View>
          <Text style={styles.actionText}>{item}</Text>
        </View>
      ))}

      {result.model ? (
        <View style={styles.metaWrap}>
          <Text style={styles.model}>Model: {result.model}</Text>
        </View>
      ) : null}
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    marginTop: 8,
    backgroundColor: t.surface,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: t.border,
    padding: 16,
    gap: 12,
  },
  sectionTitle: {
    fontSize: 12,
    fontWeight: '700',
    letterSpacing: 1.4,
    color: t.accentWarm,
    textTransform: 'uppercase',
  },
  bodyText: {
    fontSize: 15,
    color: t.text,
    lineHeight: 24,
  },
  actionRow: {
    flexDirection: 'row',
    gap: 12,
    alignItems: 'flex-start',
  },
  actionBadge: {
    marginTop: 2,
    minWidth: 22,
    height: 22,
    borderRadius: 11,
    borderWidth: 1,
    borderColor: t.borderStrong,
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: t.bg,
  },
  actionIndex: {
    fontSize: 12,
    fontWeight: '700',
    color: t.textSecondary,
  },
  actionText: {
    flex: 1,
    fontSize: 15,
    color: t.text,
    lineHeight: 24,
  },
  metaWrap: {
    marginTop: 4,
    paddingTop: 10,
    borderTopWidth: 1,
    borderTopColor: t.border,
  },
  model: {
    fontSize: 12,
    color: t.textMuted,
  },
});
