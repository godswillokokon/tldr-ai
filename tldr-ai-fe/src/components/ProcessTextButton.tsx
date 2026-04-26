import React from 'react';
import {
  ActivityIndicator,
  StyleSheet,
  Text,
  TouchableOpacity,
} from 'react-native';

import {appTheme as t} from '../theme/appTheme';

type Props = {
  loading: boolean;
  disabled: boolean;
  onPress: () => void;
};

export function ProcessTextButton({loading, disabled, onPress}: Props) {
  return (
    <TouchableOpacity
      style={[styles.button, disabled && styles.buttonDisabled]}
      onPress={onPress}
      disabled={disabled}
      activeOpacity={0.85}>
      {loading ? (
        <ActivityIndicator color={t.ctaText} />
      ) : (
        <Text style={styles.buttonText}>Generate summary</Text>
      )}
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  button: {
    backgroundColor: t.ctaBg,
    height: 52,
    borderRadius: 8,
    alignItems: 'center',
    justifyContent: 'center',
  },
  buttonDisabled: {
    opacity: 0.38,
  },
  buttonText: {
    color: t.ctaText,
    fontWeight: '700',
    fontSize: 16,
    letterSpacing: 0.2,
  },
});
