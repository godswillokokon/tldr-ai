/**
 * @format
 */

import {parseProcessTextResponse} from '../src/domain/parseProcessResponse';

describe('parseProcessTextResponse', () => {
  it('returns ProcessResponse for valid 200 JSON', () => {
    const body = JSON.stringify({
      summary: 'Short.',
      actionItems: ['a', 'b', 'c'],
      model: 'claude:test',
    });
    const r = parseProcessTextResponse(200, true, body);
    expect(r.summary).toBe('Short.');
    expect(r.actionItems).toEqual(['a', 'b', 'c']);
    expect(r.model).toBe('claude:test');
  });

  it('throws on non-JSON body when ok', () => {
    expect(() => parseProcessTextResponse(200, true, 'not json')).toThrow(
      'Invalid JSON response from server',
    );
  });

  it('throws with status when error body is not JSON', () => {
    expect(() => parseProcessTextResponse(502, false, '<html>')).toThrow(
      'Request failed (502)',
    );
  });

  it('uses error field when present for failed response', () => {
    const body = JSON.stringify({error: 'too many requests'});
    expect(() => parseProcessTextResponse(429, false, body)).toThrow(
      'too many requests',
    );
  });

  it('rejects wrong action item count', () => {
    const body = JSON.stringify({
      summary: 'S',
      actionItems: ['a', 'b'],
    });
    expect(() => parseProcessTextResponse(200, true, body)).toThrow(
      'expected exactly 3 action items',
    );
  });

  it('rejects empty action item', () => {
    const body = JSON.stringify({
      summary: 'S',
      actionItems: ['a', '  ', 'c'],
    });
    expect(() => parseProcessTextResponse(200, true, body)).toThrow(
      'non-empty strings',
    );
  });
});
