import React, {useMemo} from 'react';
import {Platform, ScrollView, StatusBar, StyleSheet} from 'react-native';
import {
  SafeAreaProvider,
  SafeAreaView,
  useSafeAreaInsets,
} from 'react-native-safe-area-context';

import {ProcessResultCard} from './src/components/ProcessResultCard';
import {ProcessTextButton} from './src/components/ProcessTextButton';
import {TextBlockInput} from './src/components/TextBlockInput';
import {TldrHeader} from './src/components/TldrHeader';
import {UsageBar} from './src/components/UsageBar';
import {useKeyboardBottomInset} from './src/hooks/useKeyboardBottomInset';
import {useProcessText} from './src/hooks/useProcessText';
import {useUsageBudget} from './src/hooks/useUsageBudget';
import {appTheme as t} from './src/theme/appTheme';

const SCROLL_BASE_PADDING_BOTTOM = 32;

function AppInner() {
  const insets = useSafeAreaInsets();
  const keyboardBottom = useKeyboardBottomInset();
  const {usage, error: usageError, refresh: refreshUsage} = useUsageBudget();
  const {text, setText, loading, result, canSubmit, submit} = useProcessText({
    onSuccess: refreshUsage,
  });

  const scrollContentStyle = useMemo(
    () => [
      styles.container,
      {
        paddingBottom:
          SCROLL_BASE_PADDING_BOTTOM +
          keyboardBottom +
          Math.max(insets.bottom, 8),
      },
    ],
    [insets.bottom, keyboardBottom],
  );

  return (
    <SafeAreaView style={styles.safeArea} edges={['top', 'left', 'right']}>
      <StatusBar barStyle="light-content" backgroundColor={t.bg} />
      <ScrollView
        style={styles.scroll}
        contentContainerStyle={scrollContentStyle}
        indicatorStyle="white"
        keyboardShouldPersistTaps="handled"
        keyboardDismissMode={Platform.OS === 'ios' ? 'interactive' : 'on-drag'}
        showsVerticalScrollIndicator>
        <TldrHeader />

        <TextBlockInput value={text} onChangeText={setText} />

        <ProcessTextButton
          loading={loading}
          disabled={!canSubmit}
          onPress={submit}
        />

        <UsageBar usage={usage} error={usageError} />

        {result ? <ProcessResultCard result={result} /> : null}
      </ScrollView>
    </SafeAreaView>
  );
}

function App() {
  return (
    <SafeAreaProvider>
      <AppInner />
    </SafeAreaProvider>
  );
}

const styles = StyleSheet.create({
  safeArea: {
    flex: 1,
    backgroundColor: t.bg,
  },
  scroll: {
    flex: 1,
  },
  container: {
    flexGrow: 1,
    paddingHorizontal: 20,
    paddingTop: 18,
    paddingBottom: 8,
    gap: 16,
  },
});

export default App;
