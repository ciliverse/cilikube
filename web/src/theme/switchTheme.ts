import { setThemeId } from './themes'

/** Apply a theme with a short CSS transition (same as CiliTerm). */
export function switchTheme(id: string, ms = 260): void {
  const root = document.documentElement
  root.classList.add('theme-switching')
  setThemeId(id)
  window.setTimeout(() => root.classList.remove('theme-switching'), ms)
}
