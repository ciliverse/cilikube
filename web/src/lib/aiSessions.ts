import type { AiChatMessage } from '@/api/ai'

const STORAGE_KEY = 'cilikube_ai_sessions_v1'
const MAX_SESSIONS = 40
const MAX_TITLE = 64

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

export function titleFromMessages(messages: AiChatMessage[]): string {
  const first = messages.find((m) => m.role === 'user' && m.content.trim())
  if (!first) return '新对话'
  return sanitizeSessionTitle(first.content) || '新对话'
}

export function listAiSessions(): AiSession[] {
  return readAll().sort((a, b) => b.updatedAt - a.updatedAt)
}

export function getAiSession(id: string): AiSession | undefined {
  return readAll().find((s) => s.id === id)
}

export function createAiSession(meta?: { clusterId?: string; clusterName?: string }): AiSession {
  const session: AiSession = {
    id: uid(),
    title: '新对话',
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

export function saveAiSession(session: AiSession) {
  const prev = readAll().find((s) => s.id === session.id)
  const titleCustom = Boolean(session.titleCustom ?? prev?.titleCustom)
  const all = readAll().filter((s) => s.id !== session.id)
  const next: AiSession = {
    ...session,
    titleCustom,
    title: titleCustom
      ? sanitizeSessionTitle(session.title) || prev?.title || '新对话'
      : session.messages.length
        ? titleFromMessages(session.messages)
        : session.title || prev?.title || '新对话',
    updatedAt: Date.now(),
  }
  writeAll([next, ...all])
  return next
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

export function duplicateAiSession(id: string): AiSession | undefined {
  const cur = getAiSession(id)
  if (!cur) return undefined
  const copy: AiSession = {
    ...cur,
    id: uid(),
    title: sanitizeSessionTitle(`${cur.title.replace(/（副本）$/, '')}（副本）`) || '新对话（副本）',
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
export function groupAiSessions(sessions: AiSession[]): AiSessionGroup[] {
  const now = new Date()
  const today = startOfDay(now)
  const yesterday = today - 86_400_000
  const week = today - 7 * 86_400_000
  const month = today - 30 * 86_400_000

  const buckets: { key: string; label: string; sessions: AiSession[] }[] = [
    { key: 'today', label: '今天', sessions: [] },
    { key: 'yesterday', label: '昨天', sessions: [] },
    { key: 'week', label: '近 7 天', sessions: [] },
    { key: 'month', label: '近 30 天', sessions: [] },
    { key: 'older', label: '更早', sessions: [] },
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
