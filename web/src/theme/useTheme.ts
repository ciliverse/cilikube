import { useSyncExternalStore } from 'react'
import { listCustomThemes } from './customThemes'
import {
  BUILTIN_THEMES,
  getStoredThemeId,
  resolveTheme,
  setThemeId,
  subscribeTheme,
  type Theme,
} from './themes'

export function useTheme(): {
  theme: Theme
  themeId: string
  themes: Theme[]
  builtins: Theme[]
  customs: Theme[]
  setTheme: (id: string) => void
} {
  const themeId = useSyncExternalStore(subscribeTheme, getStoredThemeId, () => 'paper')
  const customs = useSyncExternalStore(subscribeTheme, listCustomThemes, () => [])
  return {
    theme: resolveTheme(themeId),
    themeId,
    themes: [...BUILTIN_THEMES, ...customs],
    builtins: BUILTIN_THEMES,
    customs,
    setTheme: setThemeId,
  }
}
