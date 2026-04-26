import {getProcessTextUrl} from '../config/apiConfig';
import {parseProcessTextResponse} from '../domain/parseProcessResponse';
import type {ProcessResponse} from '../domain/types';

function networkErrorMessage(err: unknown): string {
  if (err instanceof TypeError && err.message.includes('fetch')) {
    return 'Network error: could not reach the server. Is the backend running?';
  }
  if (err instanceof Error) {
    return err.message;
  }
  return 'Network request failed';
}

/**
 * Single gateway for the process-text HTTP contract (Open/Closed: swap implementation for tests).
 */
export async function fetchProcessText(text: string): Promise<ProcessResponse> {
  let response: Response;
  try {
    response = await fetch(getProcessTextUrl(), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({text}),
    });
  } catch (err) {
    throw new Error(networkErrorMessage(err));
  }

  const raw = await response.text();
  return parseProcessTextResponse(response.status, response.ok, raw);
}
