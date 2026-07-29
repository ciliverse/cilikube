/** UI font packs — Latin Maple by default; optional Maple Mono CN (subset, on-demand). */

export type FontPack = {
  id: string
  name: string
  display: string
  sans: string
  mono: string
  /** When true, load self-hosted cn-font-split CSS (subset woff2 chunks). */
  cjk: boolean
}

const SYSTEM_CJK =
  '"PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", "Noto Sans SC", "Noto Sans CJK SC", "Source Han Sans SC", "WenQuanYi Micro Hei", system-ui, sans-serif'

export const FONT_PACKS: FontPack[] = [
  {
    id: 'maple',
    name: 'Maple Mono',
    display: `"Maple Mono", ${SYSTEM_CJK}`,
    sans: `"Maple Mono", ${SYSTEM_CJK}`,
    mono: `"Maple Mono", ${SYSTEM_CJK}`,
    cjk: false,
  },
  {
    id: 'maple-cn',
    name: 'Maple Mono CN',
    display: `"Maple Mono CN", "Maple Mono", ${SYSTEM_CJK}`,
    sans: `"Maple Mono CN", "Maple Mono", ${SYSTEM_CJK}`,
    mono: `"Maple Mono CN", "Maple Mono", ${SYSTEM_CJK}`,
    cjk: true,
  },
]

export const DEFAULT_FONT_ID = 'maple'
export const FONT_STORAGE_KEY = 'cilikube_font'

const MAPLE_CN_CSS_ID = 'maple-cn-split'
const MAPLE_CN_CSS_HREF = '/fonts/maple-cn-split/result.css'

const listeners = new Set<() => void>()
let mapleCnLoading: Promise<void> | null = null

/** Self-hosted cn-font-split subsets (same origin — best for CN networks). */
export function ensureMapleCnCss(): Promise<void> {
  if (typeof document === 'undefined') return Promise.resolve()
  if (document.getElementById(MAPLE_CN_CSS_ID)) return Promise.resolve()
  if (mapleCnLoading) return mapleCnLoading
  mapleCnLoading = new Promise((resolve, reject) => {
    const link = document.createElement('link')
    link.id = MAPLE_CN_CSS_ID
    link.rel = 'stylesheet'
    link.href = MAPLE_CN_CSS_HREF
    link.onload = () => resolve()
    link.onerror = () => {
      mapleCnLoading = null
      reject(new Error('Failed to load Maple Mono CN'))
    }
    document.head.appendChild(link)
  })
  return mapleCnLoading
}

export function resolveFont(id?: string | null): FontPack {
  // Migrate removed / legacy ids → latin maple
  if (id === 'jetbrains' || id === 'hud') return FONT_PACKS[0]
  return FONT_PACKS.find((f) => f.id === id) || FONT_PACKS[0]
}

export function applyFont(pack: FontPack): void {
  const root = document.documentElement
  root.dataset.font = pack.id
  if (pack.cjk) {
    root.dataset.fontCjk = '1'
    void ensureMapleCnCss().catch(() => {
      /* keep system CJK fallback */
    })
  } else {
    delete root.dataset.fontCjk
  }
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
