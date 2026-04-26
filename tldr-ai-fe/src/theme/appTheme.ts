/**
 * App chrome: dark canvas, warm neutrals, high-contrast primary actions.
 */
export const appTheme = {
  /** App canvas — near-black */
  bg: '#0B0C0E',
  /** Cards / inputs */
  surface: '#13141A',
  surfaceHover: '#181A22',
  border: '#2A2D36',
  borderStrong: '#3D424D',
  /** Primary copy — warm off-white */
  text: '#F4F2ED',
  textSecondary: '#B8B5AE',
  textMuted: '#7A7770',
  /** Primary CTA */
  ctaBg: '#F4F2ED',
  ctaText: '#0B0C0E',
  ctaPressed: '#DEDAD3',
  /** Accent line / progress */
  accentLine: '#E8E4DC',
  accentWarm: '#C9B8A2',
  danger: '#E85D5D',
  dangerMuted: '#F0A0A0',
} as const;

export type AppTheme = typeof appTheme;
