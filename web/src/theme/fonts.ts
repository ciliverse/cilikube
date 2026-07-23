/** UI font packs — HUD (Orbitron) or Maple Mono. */

export type FontPack = {
  id: string
  name: string
  display: string
  sans: string
  mono: string
}

export const FONT_PACKS: FontPack[] = [
  {
    id: 'hud',
    name: 'HUD (Orbitron)',
    display: '"Orbitron", "Maple Mono", ui-monospace, monospace',
    sans: '"Maple Mono", ui-monospace, Menlo, Consolas, monospace',
    mono: '"Maple Mono", ui-monospace, Menlo, Consolas, monospace',
  },
  {
    id: 'maple',
    name: 'Maple Mono',
    display: '"Maple Mono", ui-monospace, monospace',
    sans: '"Maple Mono", ui-monospace, Menlo, Consolas, monospace',
    mono: '"Maple Mono", ui-monospace, Menlo, Consolas, monospace',
  },
]

export const DEFAULT_FONT_ID = 'hud'
export const FONT_STORAGE_KEY = 'cilikube_font'

const listeners = new Set<() => void>()

export function resolveFont(id?: string | null): FontPack {
  // Removed JetBrains pack — migrate old preference
  if (id === 'jetbrains') return FONT_PACKS[0]
  return FONT_PACKS.find((f) => f.id === id) || FONT_PACKS[0]
}

export function applyFont(pack: FontPack): void {
  const root = document.documentElement
  root.dataset.font = pack.id
  root.style.setProperty('--font-display', pack.display)
  root.style.setProperty('--font-sans', pack.sans)
  root.style.setProperty('--font-mono', pack.mono)
}

export function getStoredFontId(): string {
  try {
    const raw = localStorage.getItem(FONT_STORAGE_KEY)
    return resolveFont(raw).id
  } catch {
    return DEFAULT_FONT_ID
  }
}

export function setFontId(id: string): void {
  const pack = resolveFont(id)
  try {
    localStorage.setItem(FONT_STORAGE_KEY, pack.id)
  } catch {
    /* ignore */
  }
  applyFont(pack)
  listeners.forEach((l) => l())
}

export function subscribeFont(listener: () => void): () => void {
  listeners.add(listener)
  return () => listeners.delete(listener)
}

export function initFont(): FontPack {
  const pack = resolveFont(getStoredFontId())
  applyFont(pack)
  return pack
}
