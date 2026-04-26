/**
 * @format
 */

import {fetchProcessText} from '../src/data/processTextRepository';

describe('fetchProcessText', () => {
  const originalFetch = global.fetch;

  afterEach(() => {
    global.fetch = originalFetch;
    jest.resetAllMocks();
  });

  it('maps fetch TypeError to a clear network message', async () => {
    global.fetch = jest.fn().mockRejectedValue(new TypeError('Failed to fetch'));

    await expect(fetchProcessText('x'.repeat(25))).rejects.toThrow(
      'Network error',
    );
  });

  it('returns parsed body on success', async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: true,
      status: 200,
      text: async () =>
        JSON.stringify({
          summary: 'Done.',
          actionItems: ['one', 'two', 'three'],
        }),
    });

    const r = await fetchProcessText('x'.repeat(25));
    expect(r.summary).toBe('Done.');
    expect(r.actionItems).toEqual(['one', 'two', 'three']);
  });
});
