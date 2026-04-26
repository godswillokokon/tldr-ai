import type {ProcessResponse} from './types';

type ErrorBody = {
  error?: string;
};

/**
 * Parses and validates the HTTP response for POST /api/processText (pure — easy to test).
 */
export function parseProcessTextResponse(
  status: number,
  ok: boolean,
  rawBody: string,
): ProcessResponse {
  let data: unknown;
  try {
    data = rawBody.trim() === '' ? {} : JSON.parse(rawBody);
  } catch {
    if (!ok) {
      throw new Error(`Request failed (${status})`);
    }
    throw new Error('Invalid JSON response from server');
  }

  const obj = data as ProcessResponse & ErrorBody;

  if (!ok) {
    const msg =
      typeof obj?.error === 'string' && obj.error.length > 0
        ? obj.error
        : `Request failed (${status})`;
    throw new Error(msg);
  }

  if (
    obj == null ||
    typeof obj !== 'object' ||
    typeof obj.summary !== 'string' ||
    !Array.isArray(obj.actionItems)
  ) {
    throw new Error('Invalid response shape from server');
  }

  if (obj.actionItems.length !== 3) {
    throw new Error('Invalid response: expected exactly 3 action items');
  }

  for (const item of obj.actionItems) {
    if (typeof item !== 'string' || item.trim() === '') {
      throw new Error('Invalid response: action items must be non-empty strings');
    }
  }

  const out: ProcessResponse = {
    summary: obj.summary,
    actionItems: [...obj.actionItems],
  };
  if (typeof obj.model === 'string' && obj.model.length > 0) {
    out.model = obj.model;
  }
  return out;
}
