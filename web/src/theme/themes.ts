/** CiliKube theme packs — same idea as CiliTerm: semantic colors → CSS variables. */

export type ThemeColors = {
  bg: string
  bgPanel: string
  bgPanelSolid: string
  primary: string
  primaryDim: string
  secondary: string
  red: string
  green: string
  text: string
  textDim: string
}

export type Theme = {
  id: string
  name: string
  /** light | dark — for system/OS hints and terminal defaults */
  mode: 'dark' | 'light'
  colors: ThemeColors
}

export const BUILTIN_THEMES: Theme[] = [
  {
    id: 'tron',
    name: 'TRON',
    mode: 'dark',
    colors: {
      bg: '#020508',
      bgPanel: 'rgba(6, 16, 22, 0.9)',
      bgPanelSolid: '#061018',
      primary: '#2ef0ff',
      primaryDim: '#178a9c',
      secondary: '#ffb000',
      red: '#ff3d52',
      green: '#3dffb0',
      text: '#e8fbff',
      textDim: '#7aafbb',
    },
  },
  {
    id: 'paper',
    name: 'Paper',
    mode: 'light',
    colors: {
      bg: '#ebe6da',
      bgPanel: 'rgba(255, 253, 248, 0.97)',
      bgPanelSolid: '#fffdf8',
      primary: '#065a74',
      primaryDim: '#2f6678',
      secondary: '#a03d0a',
      red: '#a82222',
      green: '#0f6b42',
      text: '#0b1218',
      textDim: '#3a4a54',
    },
  },
  {
    id: 'matrix',
    name: 'Matrix',
    mode: 'dark',
    colors: {
      bg: '#000400',
      bgPanel: 'rgba(2, 14, 6, 0.92)',
      bgPanelSolid: '#03140a',
      primary: '#2dff7a',
      primaryDim: '#1a7a3f',
      secondary: '#c8ff2e',
      red: '#ff4444',
      green: '#7affc0',
      text: '#d4ffe0',
      textDim: '#5a9a68',
    },
  },
  {
    id: 'amber',
    name: 'Amber',
    mode: 'dark',
    colors: {
      bg: '#080400',
      bgPanel: 'rgba(22, 14, 2, 0.92)',
      bgPanelSolid: '#120c02',
      primary: '#ffb400',
      primaryDim: '#9a6800',
      secondary: '#ff6a1a',
      red: '#ff4a2e',
      green: '#d4ef40',
      text: '#ffe29a',
      textDim: '#b88a30',
    },
  },
  {
    id: 'nord',
    name: 'Nord',
    mode: 'dark',
    colors: {
      bg: '#242933',
      bgPanel: 'rgba(46, 52, 64, 0.94)',
      bgPanelSolid: '#3b4252',
      primary: '#8fdbef',
      primaryDim: '#5e81ac',
      secondary: '#ebcb8b',
      red: '#bf616a',
      green: '#a3be8c',
      text: '#eceff4',
      textDim: '#9aa3b5',
    },
  },
]

export const DEFAULT_THEME_ID = 'tron'
export const THEME_STORAGE_KEY = 'cilikube_theme'

export function hexToRgba(hex: string, alpha: number): string {
  let m = hex.replace('#', '').trim()
  if (m.length === 3) m = [...m].map((c) => c + c).join('')
  const r = parseInt(m.slice(0, 2), 16)
  const g = parseInt(m.slice(2, 4), 16)
  const b = parseInt(m.slice(4, 6), 16)
  return `rgba(${r}, ${g}, ${b}, ${alpha})`
}

export function resolveTheme(id?: string | null): Theme {
  return BUILTIN_THEMES.find((t) => t.id === id) || BUILTIN_THEMES[0]
}

/** Apply theme colors to :root (Tailwind @theme + custom HUD classes). */
export function applyTheme(theme: Theme): void {
  const c = theme.colors
  const root = document.documentElement
  const set = (k: string, v: string) => root.style.setProperty(k, v)

  root.dataset.theme = theme.id
  root.dataset.themeMode = theme.mode

  set('--color-bg', c.bg)
  set('--color-panel', c.bgPanel)
  set('--color-panel-solid', c.bgPanelSolid)
  const light = theme.mode === 'light'
  set('--color-line', hexToRgba(c.primary, light ? 0.38 : 0.26))
  set('--color-cyan', c.primary)
  set('--color-cyan-dim', c.primaryDim)
  set('--color-cyan-faint', hexToRgba(c.primary, light ? 0.2 : 0.18))
  set('--color-orange', c.secondary)
  set('--color-text', c.text)
  set('--color-text-dim', c.textDim)
  set('--color-ok', c.green)
  set('--color-warn', c.secondary)
  set('--color-danger', c.red)
  set('--color-accent', c.primary)
  set('--color-accent-bright', c.primary)
  set('--color-signal', c.secondary)
  set('--color-ink', c.text)
  set('--color-ink-soft', c.textDim)
  set('--color-mist', hexToRgba(c.primary, light ? 0.11 : 0.09))
  set('--color-ambient', hexToRgba(c.primary, light ? 0.14 : 0.1))
  set('--color-grid-line', hexToRgba(c.primary, light ? 0.12 : 0.08))
  set('--color-selection', hexToRgba(c.primary, light ? 0.26 : 0.32))
  set('--color-glow', hexToRgba(c.primary, light ? 0.14 : 0.4))
  set('--color-scroll-thumb', hexToRgba(c.primary, light ? 0.5 : 0.45))
  set('--color-scroll-track', hexToRgba(c.bgPanelSolid, light ? 0.85 : 0.55))
}

const listeners = new Set<() => void>()

export function getStoredThemeId(): string {
  try {
    return localStorage.getItem(THEME_STORAGE_KEY) || DEFAULT_THEME_ID
  } catch {
    return DEFAULT_THEME_ID
  }
}

export function setThemeId(id: string): void {
  const theme = resolveTheme(id)
  try {
    localStorage.setItem(THEME_STORAGE_KEY, theme.id)
  } catch {
    /* ignore */
  }
  applyTheme(theme)
  listeners.forEach((l) => l())
}

export function subscribeTheme(listener: () => void): () => void {
  listeners.add(listener)
  return () => listeners.delete(listener)
}

export function initTheme(): Theme {
  const theme = resolveTheme(getStoredThemeId())
  applyTheme(theme)
  return theme
}
