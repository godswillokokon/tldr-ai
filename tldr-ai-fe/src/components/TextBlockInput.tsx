import React from 'react';
import {StyleSheet, TextInput} from 'react-native';

import {appTheme as t} from '../theme/appTheme';

type Props = {
  value: string;
  onChangeText: (value: string) => void;
};

/**
 * Multiline input — no custom keyboard accessory. On iOS, drag the keyboard down
 * (ScrollView `keyboardDismissMode="interactive"`) or tap away to dismiss; Return inserts newlines.
 */
export function TextBlockInput({value, onChangeText}: Props) {
  return (
    <TextInput
      style={styles.input}
      multiline
      placeholder="Paste or type your text (at least 20 characters)…"
      placeholderTextColor={t.textMuted}
      value={value}
      onChangeText={onChangeText}
      textAlignVertical="top"
      selectionColor={t.accentWarm}
      cursorColor={t.accentLine}
      returnKeyType="default"
      blurOnSubmit={false}
    />
  );
}

const styles = StyleSheet.create({
  input: {
    minHeight: 168,
    borderWidth: 1,
    borderColor: t.border,
    borderRadius: 8,
    padding: 14,
    backgroundColor: t.surface,
    fontSize: 15,
    lineHeight: 22,
    color: t.text,
  },
});
