import { hexToRgba, type Theme, type ThemeColors, type TerminalTheme } from './themes'

const STORAGE_KEY = 'cilikube_custom_themes'

let cache: Theme[] = read()

function read(): Theme[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return []
    const parsed = JSON.parse(raw) as Theme[]
    return Array.isArray(parsed) ? parsed.filter((t) => t?.id?.startsWith('custom-')) : []
  } catch {
    return []
  }
}

function write(next: Theme[]) {
  cache = next
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(next))
  } catch {
    /* ignore quota */
  }
}

/** Stable snapshot for useSyncExternalStore. */
export function listCustomThemes(): Theme[] {
  return cache
}

export function findCustomTheme(id: string): Theme | undefined {
  return cache.find((t) => t.id === id)
}

export function isCustomThemeId(id: string): boolean {
  return id.startsWith('custom-')
}

export function terminalFromColors(colors: ThemeColors, mode: Theme['mode']): TerminalTheme {
  return {
    bg: mode === 'light' ? '#fafafa' : colors.bg,
    fg: colors.text,
    cursor: colors.primary,
    selection: hexToRgba(colors.primary, mode === 'light' ? 0.22 : 0.28),
  }
}

/** Copy a theme into an editable pack. The source object is not changed. */
export function duplicateTheme(source: Theme): Theme {
  const colors = { ...source.colors }
  return {
    id: `custom-${Date.now().toString(36)}`,
    name: `${source.name} copy`,
    mode: source.mode,
    colors,
    terminal: terminalFromColors(colors, source.mode),
  }
}

export function saveCustomTheme(theme: Theme, notify: () => void): Theme {
  const nextTheme: Theme = {
    ...theme,
    id: isCustomThemeId(theme.id) ? theme.id : `custom-${Date.now().toString(36)}`,
    terminal: terminalFromColors(theme.colors, theme.mode),
  }
  const exists = cache.some((t) => t.id === nextTheme.id)
  const next = exists ? cache.map((t) => (t.id === nextTheme.id ? nextTheme : t)) : [...cache, nextTheme]
  write(next)
  notify()
  return nextTheme
}

export function deleteCustomTheme(id: string, notify: () => void): void {
  write(cache.filter((t) => t.id !== id))
  notify()
}
