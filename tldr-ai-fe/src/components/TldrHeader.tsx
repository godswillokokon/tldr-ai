import React from 'react';
import {StyleSheet, Text, View} from 'react-native';

import {appTheme as t} from '../theme/appTheme';

export function TldrHeader() {
  return (
    <View style={styles.wrap}>
      <View style={styles.accentRule} />
      <Text style={styles.eyebrow}>TLDR AI</Text>
      <Text style={styles.title}>
        From long text,{' '}
        <Text style={styles.titleHuman}>clear next moves.</Text>
      </Text>
      <Text style={styles.subtitle}>
        Paste or type a note, email, or transcript — we return a tight summary and three next steps.
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  wrap: {
    gap: 12,
    paddingBottom: 6,
  },
  accentRule: {
    width: 48,
    height: 3,
    borderRadius: 2,
    backgroundColor: t.accentWarm,
    marginBottom: 2,
  },
  eyebrow: {
    fontSize: 11,
    fontWeight: '700',
    letterSpacing: 2.2,
    color: t.textMuted,
  },
  title: {
    fontSize: 26,
    fontWeight: '700',
    color: t.text,
    lineHeight: 32,
    letterSpacing: -0.3,
  },
  titleHuman: {
    color: t.accentWarm,
    fontWeight: '600',
    fontStyle: 'italic',
  },
  subtitle: {
    color: t.textSecondary,
    fontSize: 15,
    lineHeight: 22,
  },
});
