import { useEffect, useLayoutEffect, useRef, useState, type KeyboardEvent, type MouseEvent } from 'react'
import { createPortal } from 'react-dom'
import { Copy, MoreHorizontal, Pencil, Trash2 } from 'lucide-react'
import { formatSessionTime, type AiSession } from '@/lib/aiSessions'
import { cn } from '@/lib/utils'

type Props = {
  session: AiSession
  active: boolean
  onSelect: () => void
  onRename: (title: string) => void
  onDuplicate: () => void
  onDelete: () => void
}

export function AiSessionRow({ session, active, onSelect, onRename, onDuplicate, onDelete }: Props) {
  const [menuOpen, setMenuOpen] = useState(false)
  const [renaming, setRenaming] = useState(false)
  const [draft, setDraft] = useState(session.title)
  const [menuPos, setMenuPos] = useState({ top: 0, left: 0 })
  const inputRef = useRef<HTMLInputElement>(null)
  const moreBtnRef = useRef<HTMLButtonElement>(null)
  const menuRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!renaming) setDraft(session.title)
  }, [session.title, renaming])

  useEffect(() => {
    if (renaming) {
      inputRef.current?.focus()
      inputRef.current?.select()
    }
  }, [renaming])

  useLayoutEffect(() => {
    if (!menuOpen || !moreBtnRef.current) return
    const rect = moreBtnRef.current.getBoundingClientRect()
    const menuW = 140
    const menuH = 120
    const left = Math.min(rect.right - menuW, window.innerWidth - menuW - 8)
    const top =
      rect.bottom + 6 + menuH > window.innerHeight ? Math.max(8, rect.top - menuH - 6) : rect.bottom + 6
    setMenuPos({ top, left: Math.max(8, left) })
  }, [menuOpen])

  useEffect(() => {
    if (!menuOpen) return
    const onDoc = (e: Event) => {
      const t = e.target as Node
      if (menuRef.current?.contains(t) || moreBtnRef.current?.contains(t)) return
      setMenuOpen(false)
    }
    const onKey = (e: globalThis.KeyboardEvent) => {
      if (e.key === 'Escape') setMenuOpen(false)
    }
    const onScroll = () => setMenuOpen(false)
    document.addEventListener('mousedown', onDoc)
    document.addEventListener('keydown', onKey)
    window.addEventListener('resize', onScroll)
    document.addEventListener('scroll', onScroll, true)
    return () => {
      document.removeEventListener('mousedown', onDoc)
      document.removeEventListener('keydown', onKey)
      window.removeEventListener('resize', onScroll)
      document.removeEventListener('scroll', onScroll, true)
    }
  }, [menuOpen])

  const commitRename = () => {
    const next = draft.trim()
    setRenaming(false)
    if (!next || next === session.title) {
      setDraft(session.title)
      return
    }
    onRename(next)
  }

  const startRename = (e?: MouseEvent) => {
    e?.stopPropagation()
    setMenuOpen(false)
    setDraft(session.title)
    setRenaming(true)
  }

  const onRowKey = (e: KeyboardEvent) => {
    if (renaming) return
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault()
      onSelect()
    }
    if (e.key === 'F2') {
      e.preventDefault()
      startRename()
    }
  }

  return (
    <div
      role="button"
      tabIndex={0}
      className={cn('ai-ops-session', active && 'is-active', renaming && 'is-renaming')}
      onClick={() => {
        if (!renaming) onSelect()
      }}
      onDoubleClick={(e) => {
        e.preventDefault()
        startRename(e)
      }}
      onKeyDown={onRowKey}
    >
      <span className="ai-ops-session-dot" aria-hidden />
      <div className="ai-ops-session-body">
        <div className="ai-ops-session-meta">
          <span className="ai-ops-session-time">{formatSessionTime(session.updatedAt)}</span>
          <div className="ai-ops-session-actions">
            <button
              ref={moreBtnRef}
              type="button"
              className="ai-ops-session-more"
              title="更多"
              aria-label="更多操作"
              aria-expanded={menuOpen}
              onClick={(e) => {
                e.stopPropagation()
                setMenuOpen((v) => !v)
              }}
            >
              <MoreHorizontal className="h-3.5 w-3.5" />
            </button>
            {menuOpen
              ? createPortal(
                  <div
                    ref={menuRef}
                    className="ai-ops-session-menu"
                    role="menu"
                    style={{ top: menuPos.top, left: menuPos.left }}
                  >
                    <button
                      type="button"
                      role="menuitem"
                      onClick={(e) => {
                        e.stopPropagation()
                        startRename(e)
                      }}
                    >
                      <Pencil className="h-3.5 w-3.5" />
                      重命名
                    </button>
                    <button
                      type="button"
                      role="menuitem"
                      onClick={(e) => {
                        e.stopPropagation()
                        setMenuOpen(false)
                        onDuplicate()
                      }}
                    >
                      <Copy className="h-3.5 w-3.5" />
                      创建副本
                    </button>
                    <button
                      type="button"
                      role="menuitem"
                      className="is-danger"
                      onClick={(e) => {
                        e.stopPropagation()
                        setMenuOpen(false)
                        onDelete()
                      }}
                    >
                      <Trash2 className="h-3.5 w-3.5" />
                      删除
                    </button>
                  </div>,
                  document.body,
                )
              : null}
          </div>
        </div>
        {renaming ? (
          <input
            ref={inputRef}
            className="ai-ops-session-rename"
            value={draft}
            maxLength={64}
            aria-label="重命名对话"
            onClick={(e) => e.stopPropagation()}
            onChange={(e) => setDraft(e.target.value)}
            onBlur={commitRename}
            onKeyDown={(e) => {
              e.stopPropagation()
              if (e.key === 'Enter') {
                e.preventDefault()
                commitRename()
              }
              if (e.key === 'Escape') {
                e.preventDefault()
                setDraft(session.title)
                setRenaming(false)
              }
            }}
          />
        ) : (
          <div className="ai-ops-session-title" title={session.title || '未命名'}>
            {session.title || '未命名'}
          </div>
        )}
      </div>
    </div>
  )
}
