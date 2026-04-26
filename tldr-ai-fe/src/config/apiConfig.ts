import {Platform} from 'react-native';

export const API_PORT = 8080;

export const PROCESS_TEXT_PATH = '/api/processText';
export const USAGE_PATH = '/api/usage';

/**
 * Base URL for the Go server on the dev machine.
 * iOS Simulator: localhost. Android Emulator: host loopback alias.
 */
export function getApiBaseUrl(): string {
  if (Platform.OS === 'android') {
    return `http://10.0.2.2:${API_PORT}`;
  }
  return `http://localhost:${API_PORT}`;
}

export function getProcessTextUrl(): string {
  return `${getApiBaseUrl()}${PROCESS_TEXT_PATH}`;
}

export function getUsageUrl(): string {
  return `${getApiBaseUrl()}${USAGE_PATH}`;
}
