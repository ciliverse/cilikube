import { useEffect, useMemo, useRef, useState, type FormEvent, type KeyboardEvent } from 'react'
import { useTranslation } from 'react-i18next'
import { ArrowUp, Bot, CornerDownLeft, Plus, Slash, Square, Terminal, Trash2, Zap } from 'lucide-react'
import { Link } from 'react-router-dom'
import type { AiResourceRef } from '@/api/ai'
import { AiSkillEditor } from '@/components/AiSkillEditor'
import {
  defaultAiAgent,
  filterSkillsByQuery,
  parseSlashToken,
  skillGroupLabel,
  skillsByGroup,
  type AiAgentMeta,
  type AiSkillDef,
  type AiSkillGroup,
  type CustomSkillInput,
} from '@/lib/aiSkills'
import { cn } from '@/lib/utils'

type Props = {
  value: string
  onChange: (value: string) => void
  onSend: (text?: string, skillId?: string) => void
  onStop?: () => void
  busy?: boolean
  ready?: boolean
  namespaceLabel: string
  clusterLabel: string
  err?: string
  evidence?: AiResourceRef[]
  skills?: AiSkillDef[]
  agent?: AiAgentMeta
  landing?: boolean
  onSaveCustomSkill?: (input: CustomSkillInput & { id?: string }) => AiSkillDef | null
  onDeleteCustomSkill?: (id: string) => void
}

/** Stable composer — must live outside the page so IME / Chinese input is not remounted each keystroke. */
export function AiComposer({
  value,
  onChange,
  onSend,
  onStop,
  busy,
  ready,
  namespaceLabel,
  clusterLabel,
  err,
  evidence,
  skills = [],
  agent,
  landing,
  onSaveCustomSkill,
  onDeleteCustomSkill,
}: Props) {
  const { t, i18n } = useTranslation()
  const lang = i18n.language
  const resolvedAgent = agent || defaultAiAgent(lang)
  const textareaRef = useRef<HTMLTextAreaElement>(null)
  const composingRef = useRef(false)
  const menuRef = useRef<HTMLDivElement>(null)
  const maxHeight = landing ? 200 : 180
  const groups = useMemo(() => skillsByGroup(skills, lang), [skills, lang])
  const [activeGroup, setActiveGroup] = useState<AiSkillGroup>('combo')
  const [slashOpen, setSlashOpen] = useState(false)
  const [slashQuery, setSlashQuery] = useState('')
  const [slashStart, setSlashStart] = useState(0)
  const [highlight, setHighlight] = useState(0)
  const [menuForced, setMenuForced] = useState(false)
  const [editorOpen, setEditorOpen] = useState(false)
  const [editing, setEditing] = useState<AiSkillDef | null>(null)

  useEffect(() => {
    if (!groups.length) return
    if (!groups.some((g) => g.group === activeGroup)) {
      setActiveGroup(groups[0].group)
    }
  }, [groups, activeGroup])

  const activeSkills = useMemo(() => {
    const hit = groups.find((g) => g.group === activeGroup)
    return hit?.skills || groups[0]?.skills || []
  }, [groups, activeGroup])

  const slashMatches = useMemo(
    () => filterSkillsByQuery(skills, slashQuery, lang).slice(0, 8),
    [skills, slashQuery, lang],
  )

  const showSlashMenu = (slashOpen || menuForced) && !busy && skills.length > 0
  const canManageCustom = Boolean(onSaveCustomSkill)

  useEffect(() => {
    setHighlight(0)
    const list = menuRef.current?.querySelector('.ai-ops-slash-list')
    if (list) list.scrollTop = 0
  }, [slashQuery, showSlashMenu])

  useEffect(() => {
    if (!showSlashMenu) return
    const active = menuRef.current?.querySelector<HTMLElement>('.ai-ops-slash-item.is-active')
    active?.scrollIntoView({ block: 'nearest' })
  }, [highlight, showSlashMenu])

  const syncSlashFromValue = (next: string, cursor: number) => {
    const token = parseSlashToken(next, cursor)
    if (!token) {
      setSlashOpen(false)
      setMenuForced(false)
      setSlashQuery('')
      return
    }
    setSlashOpen(true)
    setMenuForced(false)
    setSlashStart(token.start)
    setSlashQuery(token.query)
  }

  const resize = () => {
    const el = textareaRef.current
    if (!el || composingRef.current) return
    el.style.height = 'auto'
    el.style.height = `${Math.min(el.scrollHeight, maxHeight)}px`
  }

  const clearSlashToken = () => {
    const el = textareaRef.current
    const cursor = el?.selectionStart ?? value.length
    const token = parseSlashToken(value, cursor)
    const start = token?.start ?? slashStart
    const before = value.slice(0, start).replace(/\s+$/, '')
    onChange(before ? `${before} ` : '')
    setSlashOpen(false)
    setMenuForced(false)
    setSlashQuery('')
  }

  const submit = (e?: FormEvent) => {
    e?.preventDefault()
    if (showSlashMenu && slashMatches.length > 0) return
    setSlashOpen(false)
    setMenuForced(false)
    onSend(value)
  }

  const applySkill = (skill: AiSkillDef) => {
    const el = textareaRef.current
    const cursor = el?.selectionStart ?? value.length
    const token = parseSlashToken(value, cursor)
    const start = token?.start ?? slashStart
    const before = value.slice(0, start).replace(/\s+$/, '')
    onChange(before ? `${before} ` : '')
    setSlashOpen(false)
    setMenuForced(false)
    setSlashQuery('')
    requestAnimationFrame(() => {
      const box = textareaRef.current
      if (box) {
        box.style.height = 'auto'
        box.style.height = `${Math.min(box.scrollHeight, maxHeight)}px`
        box.focus()
      }
    })
    onSend(skill.prompt, skill.id)
  }

  const openCreateEditor = (seed?: AiSkillDef | null) => {
    setSlashOpen(false)
    setMenuForced(false)
    if (showSlashMenu) clearSlashToken()
    setEditing(seed || null)
    setEditorOpen(true)
  }

  const openSkillMenu = () => {
    if (busy || !skills.length) return
    const el = textareaRef.current
    const cursor = el?.selectionStart ?? value.length
    const needsSlash = !parseSlashToken(value, cursor)
    if (needsSlash) {
      const insertAt = cursor
      const next = `${value.slice(0, insertAt)}/${value.slice(insertAt)}`
      onChange(next)
      requestAnimationFrame(() => {
        const box = textareaRef.current
        if (!box) return
        const pos = insertAt + 1
        box.focus()
        box.setSelectionRange(pos, pos)
        syncSlashFromValue(next, pos)
        resize()
      })
      return
    }
    setMenuForced(true)
    setSlashOpen(true)
  }

  const onKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.nativeEvent.isComposing || e.key === 'Process' || composingRef.current) return

    if (showSlashMenu) {
      if (slashMatches.length > 0 && e.key === 'ArrowDown') {
        e.preventDefault()
        setHighlight((i) => (i + 1) % slashMatches.length)
        return
      }
      if (slashMatches.length > 0 && e.key === 'ArrowUp') {
        e.preventDefault()
        setHighlight((i) => (i - 1 + slashMatches.length) % slashMatches.length)
        return
      }
      if (slashMatches.length > 0 && (e.key === 'Enter' || e.key === 'Tab')) {
        e.preventDefault()
        const skill = slashMatches[highlight]
        if (skill) applySkill(skill)
        return
      }
      if (e.key === 'Escape') {
        e.preventDefault()
        setSlashOpen(false)
        setMenuForced(false)
        return
      }
    }

    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      submit()
    }
  }

  return (
    <div className={cn('ai-ops-composer-wrap', landing && 'is-landing')}>
      {err ? <div className="ai-ops-err">{err}</div> : null}
      {evidence && evidence.length > 0 ? (
        <div className="ai-ops-evidence-strip">
          <span className="ai-ops-evidence-label">{t('ai.clues')}</span>
          <div className="ai-ops-evidence-scroll">
            {evidence.map((r) => (
              <Link key={r.href + (r.console || '')} to={r.href} className="ai-ops-card">
                {r.console ? (
                  <Terminal className="h-3.5 w-3.5 shrink-0" />
                ) : (
                  <Zap className="h-3.5 w-3.5 shrink-0" />
                )}
                <span className="truncate">{r.label || r.name}</span>
              </Link>
            ))}
          </div>
        </div>
      ) : null}

      <form className={cn('ai-ops-composer', landing && skills.length && 'has-probes')} onSubmit={submit}>
        {landing && skills.length ? (
          <div className="ai-ops-skills">
            <div className="ai-ops-skills-top">
              <div className="ai-ops-skills-tabs" role="tablist" aria-label={t('ai.skillGroups')}>
                {groups.map((g) => (
                  <button
                    key={g.group}
                    type="button"
                    role="tab"
                    aria-selected={g.group === activeGroup}
                    className={cn('ai-ops-skills-tab', g.group === activeGroup && 'is-active')}
                    onClick={() => setActiveGroup(g.group)}
                  >
                    {g.label}
                    <span className="ai-ops-skills-tab-count">{g.skills.length}</span>
                  </button>
                ))}
              </div>
              <div className="ai-ops-skills-top-actions">
                {canManageCustom ? (
                  <button
                    type="button"
                    className="ai-ops-skill-create"
                    disabled={busy}
                    onClick={() => openCreateEditor(null)}
                  >
                    <Plus className="h-3.5 w-3.5" />
                    {t('ai.custom')}
                  </button>
                ) : null}
                <span className="ai-ops-agent-pill" title={resolvedAgent.blurb}>
                  <Bot className="h-3.5 w-3.5" />
                  {resolvedAgent.name}
                </span>
              </div>
            </div>

            <div className="ai-ops-skills-panel" role="tabpanel">
              {activeGroup === 'custom' && activeSkills.length === 0 ? (
                <div className="ai-ops-skills-empty">
                  <p>{t('ai.noCustomSkills')}</p>
                  {canManageCustom ? (
                    <button type="button" className="ai-ops-skill-create is-ghost" onClick={() => openCreateEditor(null)}>
                      <Plus className="h-3.5 w-3.5" />
                      {t('ai.writeCustomSkill')}
                    </button>
                  ) : null}
                </div>
              ) : (
                <div className="ai-ops-skills-grid">
                  {activeSkills.map((s) => (
                    <div key={s.id} className={cn('ai-ops-skill-wrap', s.custom && 'is-custom')}>
                      <button
                        type="button"
                        className="ai-ops-skill"
                        disabled={busy}
                        title={s.blurb}
                        onClick={() => applySkill(s)}
                        onContextMenu={
                          s.custom && canManageCustom
                            ? (e) => {
                                e.preventDefault()
                                openCreateEditor(s)
                              }
                            : undefined
                        }
                      >
                        <span className="ai-ops-skill-code">{s.code}</span>
                        <span className="ai-ops-skill-label">{s.label}</span>
                      </button>
                      {s.custom && onDeleteCustomSkill ? (
                        <button
                          type="button"
                          className="ai-ops-skill-del"
                          title={t('ai.delete')}
                          aria-label={t('ai.deleteSkill', { label: s.label })}
                          onClick={(e) => {
                            e.preventDefault()
                            e.stopPropagation()
                            onDeleteCustomSkill(s.id)
                          }}
                        >
                          <Trash2 className="h-3 w-3" />
                        </button>
                      ) : null}
                    </div>
                  ))}
                </div>
              )}
              <p className="ai-ops-skills-hint">
                {t('ai.slashHint')}
                {canManageCustom ? t('ai.slashHintEdit') : null}
              </p>
            </div>
          </div>
        ) : null}

        <div className="ai-ops-composer-main">
          {showSlashMenu ? (
            <div ref={menuRef} className="ai-ops-slash" role="listbox" aria-label={t('ai.skillList')}>
              <div className="ai-ops-slash-head">
                <span>
                  Skills
                  {slashQuery ? <em> · /{slashQuery}</em> : null}
                </span>
                <span className="ai-ops-slash-keys">{t('ai.slashKeys')}</span>
              </div>
              <div className="ai-ops-slash-list">
                {slashMatches.length === 0 ? (
                  <div className="ai-ops-slash-empty">{t('ai.noSkillMatch')}</div>
                ) : (
                  slashMatches.map((s, i) => (
                    <button
                      key={s.id}
                      type="button"
                      role="option"
                      aria-selected={i === highlight}
                      className={cn('ai-ops-slash-item', i === highlight && 'is-active')}
                      onMouseEnter={() => setHighlight(i)}
                      onClick={() => applySkill(s)}
                    >
                      <span className="ai-ops-slash-code">{s.code}</span>
                      <span className="ai-ops-slash-body">
                        <span className="ai-ops-slash-title">
                          {s.label}
                          <span className="ai-ops-slash-group">{skillGroupLabel(s.group, lang)}</span>
                        </span>
                        <span className="ai-ops-slash-blurb">{s.blurb}</span>
                      </span>
                      {i === highlight ? <CornerDownLeft className="h-3.5 w-3.5 ai-ops-slash-enter" /> : null}
                    </button>
                  ))
                )}
              </div>
              {canManageCustom ? (
                <button type="button" className="ai-ops-slash-create" onClick={() => openCreateEditor(null)}>
                  <Plus className="h-3.5 w-3.5" />
                  {t('ai.createCustomSkill')}
                </button>
              ) : null}
            </div>
          ) : null}

          <textarea
            ref={textareaRef}
            className="ai-ops-input"
            rows={landing ? 3 : 2}
            placeholder={
              ready
                ? landing
                  ? t('ai.placeholderReady')
                  : t('ai.placeholderContinue')
                : t('ai.placeholderOffline')
            }
            value={value}
            onCompositionStart={() => {
              composingRef.current = true
            }}
            onCompositionEnd={(e) => {
              composingRef.current = false
              onChange(e.currentTarget.value)
              syncSlashFromValue(e.currentTarget.value, e.currentTarget.selectionStart ?? 0)
              requestAnimationFrame(resize)
            }}
            onChange={(e) => {
              const next = e.target.value
              onChange(next)
              syncSlashFromValue(next, e.target.selectionStart ?? next.length)
              resize()
            }}
            onClick={(e) => {
              syncSlashFromValue(e.currentTarget.value, e.currentTarget.selectionStart ?? 0)
            }}
            onKeyUp={(e) => {
              if (e.key === 'ArrowLeft' || e.key === 'ArrowRight' || e.key === 'Home' || e.key === 'End') {
                syncSlashFromValue(e.currentTarget.value, e.currentTarget.selectionStart ?? 0)
              }
            }}
            onKeyDown={onKeyDown}
          />
          <div className="ai-ops-composer-bar">
            <div className="ai-ops-composer-meta">
              <span className="ai-ops-agent-inline" title={resolvedAgent.blurb}>
                <Bot className="h-3 w-3" />
                {resolvedAgent.name}
              </span>
              <button
                type="button"
                className="ai-ops-slash-trigger"
                disabled={busy || !skills.length}
                onClick={openSkillMenu}
                title={t('ai.openSkillMenu')}
              >
                <Slash className="h-3 w-3" />
                Skill
              </button>
              {canManageCustom ? (
                <button
                  type="button"
                  className="ai-ops-slash-trigger is-muted"
                  disabled={busy}
                  onClick={() => openCreateEditor(null)}
                  title={t('ai.createCustomSkill')}
                >
                  <Plus className="h-3 w-3" />
                  {t('ai.custom')}
                </button>
              ) : null}
              <span className="ai-ops-hint">
                {namespaceLabel}
                <span className="ai-ops-hint-sep">·</span>
                {clusterLabel}
              </span>
            </div>
            {busy ? (
              <button type="button" className="ai-ops-send is-stop" onClick={onStop}>
                <Square className="h-3.5 w-3.5" />
                {t('ai.stop')}
              </button>
            ) : (
              <button
                type="submit"
                className="ai-ops-send"
                disabled={!ready || !value.trim() || (showSlashMenu && slashMatches.length > 0)}
              >
                <ArrowUp className="h-4 w-4" />
                <span>{t('ai.send')}</span>
              </button>
            )}
          </div>
        </div>
      </form>

      {canManageCustom ? (
        <AiSkillEditor
          open={editorOpen}
          initial={editing}
          onClose={() => {
            setEditorOpen(false)
            setEditing(null)
          }}
          onSave={(input) => {
            const saved = onSaveCustomSkill?.(input)
            if (!saved) return
            setEditorOpen(false)
            setEditing(null)
            setActiveGroup('custom')
          }}
        />
      ) : null}
    </div>
  )
}
