/** Successful response from POST /api/processText */
export type ProcessResponse = {
  summary: string;
  actionItems: string[];
  model?: string;
};
