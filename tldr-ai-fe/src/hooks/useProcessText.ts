import {useCallback, useMemo, useState} from 'react';
import {Alert} from 'react-native';

import {fetchProcessText} from '../data/processTextRepository';
import {MIN_INPUT_LENGTH} from '../domain/constants';
import type {ProcessResponse} from '../domain/types';

type ProcessFn = (text: string) => Promise<ProcessResponse>;

export type UseProcessTextOptions = {
  /** Injected for tests / alternate backends (Dependency Inversion). */
  processFn?: ProcessFn;
  /** Called after a successful summarize (e.g. refresh usage UI). */
  onSuccess?: () => void;
};

/**
 * Encapsulates submit state and validation for the summarization flow (SRP).
 */
export function useProcessText(options: UseProcessTextOptions = {}) {
  const run = options.processFn ?? fetchProcessText;
  const onSuccess = options.onSuccess;

  const [text, setText] = useState('');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<ProcessResponse | null>(null);

  const trimmed = text.trim();
  const canSubmit = useMemo(
    () => trimmed.length >= MIN_INPUT_LENGTH && !loading,
    [trimmed, loading],
  );

  const submit = useCallback(async () => {
    if (!canSubmit) {
      Alert.alert(
        'Input too short',
        `Please enter at least ${MIN_INPUT_LENGTH} characters of text.`,
      );
      return;
    }
    setLoading(true);
    try {
      const response = await run(trimmed);
      setResult(response);
      onSuccess?.();
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unknown error';
      Alert.alert('Request failed', message);
    } finally {
      setLoading(false);
    }
  }, [canSubmit, onSuccess, run, trimmed]);

  return {
    text,
    setText,
    loading,
    result,
    canSubmit,
    submit,
  };
}
