import { useSyncExternalStore } from 'react'
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
  setTheme: (id: string) => void
} {
  const themeId = useSyncExternalStore(subscribeTheme, getStoredThemeId, () => 'tron')
  return {
    theme: resolveTheme(themeId),
    themeId,
    themes: BUILTIN_THEMES,
    setTheme: setThemeId,
  }
}
