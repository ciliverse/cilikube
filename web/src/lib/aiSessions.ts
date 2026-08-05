import type { AiChatMessage } from '@/api/ai'

const STORAGE_KEY = 'cilikube_ai_sessions_v1'
const MAX_SESSIONS = 40
const MAX_TITLE = 64

/** Language-neutral marker for an empty chat; UI maps it via i18n (`ai.newChat`). */
export const EMPTY_CHAT_TITLE_KEY = '__new_chat__'

const LEGACY_EMPTY_TITLES = new Set(['新对话', 'New chat', '未命名', 'Untitled', EMPTY_CHAT_TITLE_KEY])

export function isEmptyChatTitle(title?: string | null): boolean {
  const t = (title || '').trim()
  return !t || LEGACY_EMPTY_TITLES.has(t)
}

/** Resolve a stored title for the current UI language. */
export function displayChatTitle(title: string | undefined, emptyLabel: string): string {
  return isEmptyChatTitle(title) ? emptyLabel : (title || emptyLabel)
}

export type AiSession = {
  id: string
  title: string
  /** When true, auto-title from first message will not overwrite. */
  titleCustom?: boolean
  createdAt: number
  updatedAt: number
  messages: AiChatMessage[]
  clusterId?: string
  clusterName?: string
}

export type AiSessionGroup = {
  key: string
  label: string
  sessions: AiSession[]
}

function uid() {
  return `ops_${Date.now().toString(36)}_${Math.random().toString(36).slice(2, 8)}`
}

function readAll(): AiSession[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return []
    const parsed = JSON.parse(raw) as AiSession[]
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

function writeAll(sessions: AiSession[]) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(sessions.slice(0, MAX_SESSIONS)))
  } catch {
    /* quota / private mode */
  }
}

export function sanitizeSessionTitle(raw: string): string {
  const t = raw.trim().replace(/\s+/g, ' ')
  if (!t) return ''
  return t.length > MAX_TITLE ? `${t.slice(0, MAX_TITLE)}…` : t
}

export function titleFromMessages(
  messages: AiChatMessage[],
  emptyTitle = EMPTY_CHAT_TITLE_KEY,
): string {
  const first = messages.find((m) => m.role === 'user' && m.content.trim())
  if (!first) return EMPTY_CHAT_TITLE_KEY
  return sanitizeSessionTitle(first.content) || emptyTitle || EMPTY_CHAT_TITLE_KEY
}

export function listAiSessions(): AiSession[] {
  return readAll().sort((a, b) => b.updatedAt - a.updatedAt)
}

export function getAiSession(id: string): AiSession | undefined {
  return readAll().find((s) => s.id === id)
}

export function createAiSession(meta?: {
  clusterId?: string
  clusterName?: string
  emptyTitle?: string
}): AiSession {
  const session: AiSession = {
    id: uid(),
    title: EMPTY_CHAT_TITLE_KEY,
    titleCustom: false,
    createdAt: Date.now(),
    updatedAt: Date.now(),
    messages: [],
    clusterId: meta?.clusterId,
    clusterName: meta?.clusterName,
  }
  writeAll([session, ...readAll()])
  return session
}

export function saveAiSession(session: AiSession, emptyTitle = EMPTY_CHAT_TITLE_KEY) {
  const prev = readAll().find((s) => s.id === session.id)
  const titleCustom = Boolean(session.titleCustom ?? prev?.titleCustom)
  const all = readAll().filter((s) => s.id !== session.id)
  let title: string
  if (titleCustom) {
    const cleaned = sanitizeSessionTitle(session.title)
    title = isEmptyChatTitle(cleaned) ? EMPTY_CHAT_TITLE_KEY : cleaned || prev?.title || EMPTY_CHAT_TITLE_KEY
  } else if (session.messages.length) {
    title = titleFromMessages(session.messages, emptyTitle)
  } else if (isEmptyChatTitle(session.title) || isEmptyChatTitle(prev?.title)) {
    title = EMPTY_CHAT_TITLE_KEY
  } else {
    title = session.title || prev?.title || EMPTY_CHAT_TITLE_KEY
  }
  const next: AiSession = {
    ...session,
    titleCustom,
    title,
    updatedAt: Date.now(),
  }
  writeAll([next, ...all])
  return next
}

/** Normalize legacy localized empty titles ("新对话" / "New chat") to the sentinel. */
export function migrateEmptyChatTitles(): AiSession[] {
  const all = readAll()
  let changed = false
  const next = all.map((s) => {
    if (!s.titleCustom && isEmptyChatTitle(s.title) && s.title !== EMPTY_CHAT_TITLE_KEY) {
      changed = true
      return { ...s, title: EMPTY_CHAT_TITLE_KEY }
    }
    return s
  })
  if (changed) writeAll(next)
  return next.sort((a, b) => b.updatedAt - a.updatedAt)
}

export function renameAiSession(id: string, title: string): AiSession | undefined {
  const cur = getAiSession(id)
  if (!cur) return undefined
  const nextTitle = sanitizeSessionTitle(title)
  if (!nextTitle) return cur
  return saveAiSession({
    ...cur,
    title: nextTitle,
    titleCustom: true,
  })
}

export function duplicateAiSession(
  id: string,
  opts?: { emptyTitle?: string; copySuffix?: string },
): AiSession | undefined {
  const cur = getAiSession(id)
  if (!cur) return undefined
  const emptyLabel = opts?.emptyTitle || 'New chat'
  const suffix = opts?.copySuffix || ' (copy)'
  const shown = displayChatTitle(cur.title, emptyLabel)
  const base = shown.replace(/\s*(\(copy\)|（副本）)$/u, '')
  const copy: AiSession = {
    ...cur,
    id: uid(),
    title: sanitizeSessionTitle(`${base}${suffix}`) || `${emptyLabel}${suffix}`,
    titleCustom: true,
    createdAt: Date.now(),
    updatedAt: Date.now(),
    messages: structuredClone(cur.messages),
  }
  writeAll([copy, ...readAll()])
  return copy
}

export function deleteAiSession(id: string) {
  writeAll(readAll().filter((s) => s.id !== id))
}

export function clearAiSessions() {
  writeAll([])
}

function startOfDay(d: Date) {
  return new Date(d.getFullYear(), d.getMonth(), d.getDate()).getTime()
}

/** Cursor-style recency buckets for the history rail. */
export function groupAiSessions(
  sessions: AiSession[],
  labels?: Partial<Record<'today' | 'yesterday' | 'week' | 'month' | 'older', string>>,
): AiSessionGroup[] {
  const now = new Date()
  const today = startOfDay(now)
  const yesterday = today - 86_400_000
  const week = today - 7 * 86_400_000
  const month = today - 30 * 86_400_000

  const buckets: { key: string; label: string; sessions: AiSession[] }[] = [
    { key: 'today', label: labels?.today || 'Today', sessions: [] },
    { key: 'yesterday', label: labels?.yesterday || 'Yesterday', sessions: [] },
    { key: 'week', label: labels?.week || 'Previous 7 days', sessions: [] },
    { key: 'month', label: labels?.month || 'Previous 30 days', sessions: [] },
    { key: 'older', label: labels?.older || 'Older', sessions: [] },
  ]

  for (const s of sessions) {
    const t = s.updatedAt
    if (t >= today) buckets[0].sessions.push(s)
    else if (t >= yesterday) buckets[1].sessions.push(s)
    else if (t >= week) buckets[2].sessions.push(s)
    else if (t >= month) buckets[3].sessions.push(s)
    else buckets[4].sessions.push(s)
  }

  return buckets.filter((b) => b.sessions.length > 0)
}

export function formatSessionTime(ts: number): string {
  const d = new Date(ts)
  const now = new Date()
  const sameDay =
    d.getFullYear() === now.getFullYear() &&
    d.getMonth() === now.getMonth() &&
    d.getDate() === now.getDate()
  if (sameDay) {
    return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }
  return d.toLocaleDateString([], { month: 'short', day: 'numeric' })
}
