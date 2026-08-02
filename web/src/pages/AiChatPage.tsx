import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { AnimatePresence, motion } from 'framer-motion'
import {
  Bot,
  History,
  Loader2,
  MessageSquarePlus,
  PanelLeftClose,
  Terminal,
  X,
  Zap,
} from 'lucide-react'
import { getAiStatus, streamAiChat, type AiChatMessage, type AiResourceRef } from '@/api/ai'
import { AiComposer } from '@/components/AiComposer'
import { AiSessionRow } from '@/components/AiSessionRow'
import {
  clearAiSessions,
  createAiSession,
  deleteAiSession,
  duplicateAiSession,
  getAiSession,
  groupAiSessions,
  listAiSessions,
  renameAiSession,
  saveAiSession,
  type AiSession,
} from '@/lib/aiSessions'
import {
  allAiSkills,
  DEFAULT_AI_AGENT,
  deleteCustomSkill,
  upsertCustomSkill,
  type AiSkillDef,
} from '@/lib/aiSkills'
import { getFleetSummary } from '@/api/cluster'
import {
  buildFleetFocusPrompt,
  buildFleetInspectPrompt,
  buildFleetTourStepPrompt,
  buildFleetTourSummaryPrompt,
  buildInvestigatePrompt,
  fleetIssueScore,
  parseFleetFocus,
  parseFleetInspectSearch,
  parseFleetTourSearch,
  parseInvestigateSearch,
  resourceRefLabel,
} from '@/lib/aiInvestigate'
import { shouldSkipEnterAnim } from '@/lib/motionPrefs'
import { useCluster } from '@/store/cluster'
import { useNamespace } from '@/store/namespace'
import { cn } from '@/lib/utils'

function useIsDesktopHistory() {
  const [desktop, setDesktop] = useState(() =>
    typeof window !== 'undefined' ? window.matchMedia('(min-width: 1024px)').matches : true,
  )
  useEffect(() => {
    const mq = window.matchMedia('(min-width: 1024px)')
    const onChange = () => setDesktop(mq.matches)
    onChange()
    mq.addEventListener('change', onChange)
    return () => mq.removeEventListener('change', onChange)
  }, [])
  return desktop
}

export function AiChatPage() {
  const { t } = useTranslation()
  const { clusterId, activeCluster, setClusterId, switchCluster, switching } = useCluster()
  const { namespace, setNamespace } = useNamespace()
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const skipMotion = shouldSkipEnterAnim()
  const statusQ = useQuery({
    queryKey: ['ai-status'],
    queryFn: getAiStatus,
    refetchInterval: 30_000,
  })

  const [sessions, setSessions] = useState<AiSession[]>(() => listAiSessions())
  const [activeId, setActiveId] = useState<string>(() => listAiSessions()[0]?.id || '')
  const [messages, setMessages] = useState<AiChatMessage[]>(() => listAiSessions()[0]?.messages || [])
  const [input, setInput] = useState('')
  const [busy, setBusy] = useState(false)
  const [err, setErr] = useState('')
  const [skills, setSkills] = useState<AiSkillDef[]>(() => allAiSkills())
  const isDesktop = useIsDesktopHistory()
  const [historyOpen, setHistoryOpen] = useState(() => {
    try {
      const saved = localStorage.getItem('cilikube_ai_rail')
      if (saved === '0') return false
      if (saved === '1') return true
    } catch {
      /* ignore */
    }
    return typeof window !== 'undefined' ? window.matchMedia('(min-width: 1024px)').matches : true
  })

  const toggleHistory = () => {
    setHistoryOpen((v) => {
      const next = !v
      try {
        localStorage.setItem('cilikube_ai_rail', next ? '1' : '0')
      } catch {
        /* ignore */
      }
      return next
    })
  }

  const closeHistory = () => {
    setHistoryOpen(false)
    try {
      localStorage.setItem('cilikube_ai_rail', '0')
    } catch {
      /* ignore */
    }
  }

  const bottomRef = useRef<HTMLDivElement>(null)
  const abortRef = useRef<AbortController | null>(null)
  /** Skip one persist after hydrating another session (avoids writing A’s messages into B). */
  const skipPersistRef = useRef(false)
  /** Consume resource-page / inspect deep-link (prompt + optional ns/cluster + auto). */
  const investigateJobRef = useRef<{
    prompt: string
    namespaceOverride?: string
    titleHint: string
    auto: boolean
    /** Switch to this cluster before sending (fleet inspect). */
    clusterId?: string
  } | null>(null)
  /** Fleet-wide serial inspect (`?tour=1`). */
  const fleetTourPendingRef = useRef<{ auto: boolean } | null>(null)
  const fleetTourCancelRef = useRef(false)
  const fleetTourRunningRef = useRef(false)

  const ready = Boolean(statusQ.data?.ready)
  const clusterLabel = activeCluster?.name || clusterId || '—'
  const isEmpty = messages.length === 0

  const refreshSessions = useCallback(() => setSessions(listAiSessions()), [])

  const hydrateSession = useCallback((id: string, nextMessages: AiChatMessage[]) => {
    skipPersistRef.current = true
    setActiveId(id)
    setMessages(nextMessages)
    setErr('')
    setInput('')
  }, [])

  const ensureSession = useCallback((): string => {
    if (activeId && getAiSession(activeId)) return activeId
    const s = createAiSession({
      clusterId: clusterId || undefined,
      clusterName: activeCluster?.name,
    })
    hydrateSession(s.id, [])
    refreshSessions()
    return s.id
  }, [activeId, clusterId, activeCluster?.name, hydrateSession, refreshSessions])

  useEffect(() => {
    if (!activeId) {
      const s = createAiSession({
        clusterId: clusterId || undefined,
        clusterName: activeCluster?.name,
      })
      hydrateSession(s.id, [])
      refreshSessions()
    }
  }, [activeId, clusterId, activeCluster?.name, hydrateSession, refreshSessions])

  // Resource page → /ai?investigate=1&kind=&name=&namespace=
  // Fleet card → /ai?inspect=1&cluster=&name=&focus=
  // Fleet tour → /ai?tour=1
  useEffect(() => {
    const tour = parseFleetTourSearch(searchParams)
    if (tour) {
      fleetTourPendingRef.current = { auto: tour.auto }
      navigate('/ai', { replace: true })
      return
    }
    const fleet = parseFleetInspectSearch(searchParams)
    if (fleet) {
      const auto = searchParams.get('auto') !== '0'
      const focus = parseFleetFocus(searchParams)
      investigateJobRef.current = {
        prompt: focus
          ? buildFleetFocusPrompt(fleet.clusterName, focus)
          : buildFleetInspectPrompt(fleet.clusterName),
        titleHint: focus
          ? `${fleet.clusterName} · ${focus === 'unhealthy' ? '异常 Pod' : 'Warning'}`
          : fleet.clusterName,
        auto,
        clusterId: fleet.clusterId,
      }
      if (
        fleet.clusterId &&
        fleet.clusterId !== clusterId &&
        activeCluster?.name !== fleet.clusterId
      ) {
        setClusterId(fleet.clusterId)
      }
      navigate('/ai', { replace: true })
      return
    }
    const target = parseInvestigateSearch(searchParams)
    if (!target) return
    const auto = searchParams.get('auto') !== '0'
    investigateJobRef.current = {
      prompt: buildInvestigatePrompt(target),
      namespaceOverride: target.namespace,
      titleHint: resourceRefLabel(target),
      auto,
    }
    navigate('/ai', { replace: true })
  }, [searchParams, navigate, clusterId, activeCluster?.name, setClusterId])

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: skipMotion ? 'auto' : 'smooth' })
  }, [messages, busy, skipMotion])

  useEffect(() => {
    if (skipPersistRef.current) {
      skipPersistRef.current = false
      return
    }
    if (!activeId || busy) return
    const cur = getAiSession(activeId)
    if (!cur) return
    saveAiSession({
      ...cur,
      messages,
      clusterId: clusterId || cur.clusterId,
      clusterName: activeCluster?.name || cur.clusterName,
    })
    refreshSessions()
  }, [messages, activeId, busy, clusterId, activeCluster?.name, refreshSessions])

  const evidence = useMemo(() => {
    const refs: AiResourceRef[] = []
    const seen = new Set<string>()
    for (const m of messages) {
      for (const r of m.resources || []) {
        const k = r.href + (r.console || '')
        if (seen.has(k)) continue
        seen.add(k)
        refs.push(r)
      }
    }
    return refs
  }, [messages])

  const stopGeneration = useCallback(() => {
    fleetTourCancelRef.current = true
    abortRef.current?.abort()
    abortRef.current = null
    setBusy(false)
  }, [])

  const startNew = () => {
    stopGeneration()
    const s = createAiSession({
      clusterId: clusterId || undefined,
      clusterName: activeCluster?.name,
    })
    hydrateSession(s.id, [])
    if (!isDesktop) closeHistory()
    refreshSessions()
  }

  const selectSession = (id: string) => {
    if (id === activeId) return
    const s = getAiSession(id)
    if (!s) return
    stopGeneration()
    hydrateSession(id, s.messages || [])
    if (!isDesktop) closeHistory()
  }

  const sessionGroups = useMemo(() => groupAiSessions(sessions), [sessions])

  const removeSession = (id: string) => {
    const cur = getAiSession(id)
    const label = cur?.title || '该对话'
    if (!window.confirm(`删除「${label}」？此操作不可恢复。`)) return
    stopGeneration()
    deleteAiSession(id)
    const rest = listAiSessions()
    refreshSessions()
    if (id === activeId) {
      if (rest[0]) {
        hydrateSession(rest[0].id, rest[0].messages || [])
      } else {
        startNew()
      }
    }
  }

  const renameSession = (id: string, title: string) => {
    renameAiSession(id, title)
    refreshSessions()
  }

  const duplicateSession = (id: string) => {
    const copy = duplicateAiSession(id)
    if (!copy) return
    stopGeneration()
    hydrateSession(copy.id, copy.messages || [])
    refreshSessions()
    if (!isDesktop) closeHistory()
  }

  const clearAllSessions = () => {
    if (!sessions.length) return
    if (!window.confirm(`清空全部 ${sessions.length} 条对话历史？此操作不可恢复。`)) return
    stopGeneration()
    clearAiSessions()
    refreshSessions()
    startNew()
  }

  const stop = () => {
    stopGeneration()
  }

  const send = async (
    text?: string,
    skillId?: string,
    opts?: {
      baseMessages?: AiChatMessage[]
      namespaceOverride?: string
      sessionId?: string
      /** Skip in-flight guard (e.g. resource-page investigate after stop). */
      force?: boolean
    },
  ): Promise<AiChatMessage[] | undefined> => {
    const content = (text ?? input).trim()
    if (!content) return undefined
    if (busy && !opts?.force) {
      setErr('正在生成回复，请先点「停止」再发送')
      return undefined
    }
    if (!ready) {
      setErr(statusQ.isError ? t('ai.connectFailed') : t('ai.unavailable'))
      return undefined
    }

    if (!opts?.sessionId) ensureSession()
    setErr('')
    setInput('')

    const prior = opts?.baseMessages ?? messages
    const history = [...prior, { role: 'user' as const, content }]
    setMessages(history)
    setBusy(true)

    let assistant: AiChatMessage = { role: 'assistant', content: '', tools: [], resources: [] }
    setMessages([...history, assistant])

    abortRef.current?.abort()
    const ac = new AbortController()
    abortRef.current = ac

    const ns =
      opts?.namespaceOverride !== undefined
        ? opts.namespaceOverride
        : namespace === 'all'
          ? ''
          : namespace

    const patchAssistant = (fn: (m: AiChatMessage) => AiChatMessage) => {
      assistant = fn(assistant)
      setMessages([...history, assistant])
    }

    try {
      await streamAiChat(
        history.map((m) => ({ role: m.role, content: m.content })),
        {
          namespace: ns,
          skillId,
          signal: ac.signal,
        },
        {
          onToolCall: (name) => {
            patchAssistant((m) => ({
              ...m,
              tools: [...(m.tools || []), { name, ok: undefined }],
            }))
          },
          onToolResult: (name, ok, toolContent) => {
            patchAssistant((m) => {
              const tools = [...(m.tools || [])]
              const idx = tools.map((t) => t.name).lastIndexOf(name)
              if (idx >= 0) tools[idx] = { name, ok, content: toolContent }
              else tools.push({ name, ok, content: toolContent })
              return { ...m, tools }
            })
          },
          onResources: (items) => {
            patchAssistant((m) => ({ ...m, resources: mergeRefs(m.resources || [], items) }))
          },
          onMessage: (msg) => patchAssistant((m) => ({ ...m, content: msg })),
          onError: (message) => {
            setErr(message)
            patchAssistant((m) => ({
              ...m,
              content: m.content || `出错了：${message}`,
            }))
          },
        },
      )
    } catch (e: any) {
      if (e?.name !== 'AbortError') setErr(e?.message || '请求失败')
      else return undefined
    } finally {
      setBusy(false)
    }
    return [...history, assistant]
  }

  // Run pending investigate/inspect job after AI is ready (and cluster switch settles).
  useEffect(() => {
    const job = investigateJobRef.current
    if (!job) return
    if (statusQ.isLoading || switching) return
    if (
      job.clusterId &&
      job.clusterId !== clusterId &&
      activeCluster?.name !== job.clusterId &&
      activeCluster?.id !== job.clusterId
    ) {
      return
    }

    investigateJobRef.current = null
    if (job.namespaceOverride) setNamespace(job.namespaceOverride)

    abortRef.current?.abort()
    abortRef.current = null
    setBusy(false)
    const s = createAiSession({
      clusterId: clusterId || undefined,
      clusterName: activeCluster?.name,
    })
    hydrateSession(s.id, [])
    renameAiSession(s.id, `${job.clusterId ? '巡检' : '调查'} · ${job.titleHint}`)
    refreshSessions()
    if (!isDesktop) closeHistory()

    if (job.auto) {
      if (!ready) {
        setInput(job.prompt)
        setErr(
          statusQ.isError
            ? '无法连接 AI 服务，请确认后端已启动'
            : 'AI 暂不可用，调查 Prompt 已填入输入框，配置好后可直接发送',
        )
        return
      }
      void send(job.prompt, undefined, {
        baseMessages: [],
        namespaceOverride: job.namespaceOverride || '',
        sessionId: s.id,
        force: true,
      })
      return
    }
    setInput(job.prompt)
    // job consumed via ref; deps only gate "when AI/cluster settles"
    // eslint-disable-next-line react-hooks/exhaustive-deps -- run once per queued investigate job
  }, [
    ready,
    statusQ.isLoading,
    statusQ.isError,
    switching,
    clusterId,
    activeCluster?.name,
    activeCluster?.id,
    isDesktop,
    setNamespace,
    hydrateSession,
    refreshSessions,
  ])

  // Fleet tour: serial inspect across reachable clusters, then summary.
  useEffect(() => {
    const pending = fleetTourPendingRef.current
    if (!pending || fleetTourRunningRef.current) return
    if (statusQ.isLoading) return

    fleetTourPendingRef.current = null
    fleetTourRunningRef.current = true
    fleetTourCancelRef.current = false

    void (async () => {
      try {
        const summary = await getFleetSummary()
        const reachable = [...(summary.clusters || [])]
          .filter((c) => c.reachable && c.id)
          .sort((a, b) => fleetIssueScore(b) - fleetIssueScore(a))
        const targets = reachable.slice(0, 8).map((c) => ({ id: c.id, name: c.name }))

        if (!targets.length) {
          setErr('没有可达集群可巡检，请先在集群总览确认连通性')
          return
        }

        abortRef.current?.abort()
        abortRef.current = null
        setBusy(false)

        const s = createAiSession({
          clusterId: targets[0].id,
          clusterName: targets[0].name,
        })
        hydrateSession(s.id, [])
        renameAiSession(
          s.id,
          `舰队健康 · ${targets.length} 集群${reachable.length > targets.length ? '+' : ''}`,
        )
        refreshSessions()
        if (!isDesktop) closeHistory()

        if (!pending.auto || !ready) {
          setInput(buildFleetTourSummaryPrompt(targets.map((c) => c.name)))
          setErr(
            ready
              ? '舰队巡检未自动开始，可编辑后发送；或返回集群总览再点「舰队巡检」'
              : statusQ.isError
                ? '无法连接 AI 服务，请确认后端已启动'
                : 'AI 暂不可用，汇总 Prompt 已填入输入框',
          )
          return
        }

        let base: AiChatMessage[] = []
        for (let i = 0; i < targets.length; i++) {
          if (fleetTourCancelRef.current) break
          const c = targets[i]
          await switchCluster(c.id)
          if (fleetTourCancelRef.current) break
          const next = await send(buildFleetTourStepPrompt(c.name, i + 1, targets.length), undefined, {
            baseMessages: base,
            namespaceOverride: '',
            sessionId: s.id,
            force: true,
          })
          if (!next) break
          base = next
        }

        if (!fleetTourCancelRef.current && base.length > 0) {
          await send(buildFleetTourSummaryPrompt(targets.map((c) => c.name)), undefined, {
            baseMessages: base,
            namespaceOverride: '',
            sessionId: s.id,
            force: true,
          })
        }
      } catch (e: any) {
        if (!fleetTourCancelRef.current) {
          setErr(e?.message || '舰队巡检失败')
        }
      } finally {
        fleetTourRunningRef.current = false
      }
    })()
    // eslint-disable-next-line react-hooks/exhaustive-deps -- one-shot tour kickoff
  }, [ready, statusQ.isLoading, statusQ.isError, isDesktop, hydrateSession, refreshSessions, switchCluster])

  const composer = (
    <AiComposer
      value={input}
      onChange={setInput}
      onSend={(text, skillId) => {
        void send(typeof text === 'string' ? text : input, skillId)
      }}
      onStop={stop}
      busy={busy}
      ready={ready}
      namespaceLabel={namespace === 'all' ? '全部命名空间' : namespace}
      clusterLabel={clusterLabel}
      err={err}
      evidence={evidence}
      skills={skills}
      agent={DEFAULT_AI_AGENT}
      landing={isEmpty}
      onSaveCustomSkill={(input) => {
        const saved = upsertCustomSkill(input)
        if (saved) setSkills(allAiSkills())
        return saved
      }}
      onDeleteCustomSkill={(id) => {
        if (deleteCustomSkill(id)) setSkills(allAiSkills())
      }}
    />
  )

  const historyRail = (onClose?: () => void) => (
    <div className="ai-ops-rail">
      <div className="ai-ops-rail-head">
        <div>
          <div className="ai-ops-kicker">History</div>
          <div className="ai-ops-rail-title">对话</div>
        </div>
        <div className="ai-ops-rail-actions">
          <button type="button" className="ai-ops-icon-btn" onClick={startNew} title="新对话">
            <MessageSquarePlus className="h-4 w-4" />
          </button>
          {onClose ? (
            <button type="button" className="ai-ops-icon-btn" onClick={onClose} title="关闭历史" aria-label="关闭历史">
              <X className="h-4 w-4" />
            </button>
          ) : null}
        </div>
      </div>

      <button type="button" className="ai-ops-new" onClick={startNew}>
        新对话
      </button>

      <div className="ai-ops-session-list">
        {sessions.length === 0 ? (
          <p className="ai-ops-rail-empty">提问后会出现在这里</p>
        ) : (
          <div className="ai-ops-timeline">
            {sessionGroups.map((group) => (
              <div key={group.key} className="ai-ops-session-group">
                <div className="ai-ops-session-group-label">{group.label}</div>
                {group.sessions.map((s) => (
                  <AiSessionRow
                    key={s.id}
                    session={s}
                    active={s.id === activeId}
                    onSelect={() => selectSession(s.id)}
                    onRename={(title) => renameSession(s.id, title)}
                    onDuplicate={() => duplicateSession(s.id)}
                    onDelete={() => removeSession(s.id)}
                  />
                ))}
              </div>
            ))}
          </div>
        )}
      </div>

      {sessions.length > 0 ? (
        <button type="button" className="ai-ops-clear-all" onClick={clearAllSessions}>
          清空历史
        </button>
      ) : null}
    </div>
  )

  return (
    <div className={cn('ai-ops', isEmpty && 'is-landing', historyOpen && 'rail-open')}>
      <div className="ai-ops-atmosphere" aria-hidden />

      <aside
        className={cn('ai-ops-aside', historyOpen ? 'is-open' : 'is-closed')}
        aria-hidden={!historyOpen || !isDesktop}
      >
        {historyRail(closeHistory)}
      </aside>

      <AnimatePresence>
        {historyOpen && !isDesktop ? (
          <>
            <motion.button
              type="button"
              className="ai-ops-backdrop"
              aria-label="关闭历史"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={closeHistory}
            />
            <motion.aside
              className="ai-ops-drawer"
              initial={{ x: '-100%' }}
              animate={{ x: 0 }}
              exit={{ x: '-100%' }}
              transition={{ type: 'tween', duration: 0.2 }}
            >
              {historyRail(closeHistory)}
            </motion.aside>
          </>
        ) : null}
      </AnimatePresence>

      <section className="ai-ops-stage">
        <div className="ai-ops-toolbar">
          <div className="ai-ops-toolbar-left">
            <button
              type="button"
              className="ai-ops-icon-btn ai-ops-history-btn"
              onClick={toggleHistory}
              aria-label={historyOpen ? t('ai.closeHistory') : t('ai.openHistory')}
              aria-expanded={historyOpen}
            >
              {historyOpen && isDesktop ? <PanelLeftClose className="h-4 w-4" /> : <History className="h-4 w-4" />}
            </button>
            <span className={cn('ai-ops-status', ready ? 'is-on' : 'is-off')}>
              <i />
              {ready ? t('ai.online') : t('ai.offline')}
            </span>
            <span className="ai-ops-cluster" title={clusterLabel}>
              {clusterLabel}
            </span>
          </div>
          <Link to="/overview" className="ai-ops-console-btn">
            <Terminal className="h-3.5 w-3.5" />
            {t('ai.console')}
          </Link>
        </div>

        {!ready ? (
          <div className="ai-ops-banner">{t('ai.unavailable')}</div>
        ) : null}

        {isEmpty ? (
          <div className="ai-ops-landing">
            <motion.div
              className="ai-ops-stage-core"
              initial={skipMotion ? false : { opacity: 0, y: 12 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.4 }}
            >
              <BrandTitle skipMotion={skipMotion} />
              <h2 className="ai-ops-hero-title">
                {t('ai.heroTitle')}
                <em>{t('ai.heroTitleEm')}</em>
              </h2>
              <p className="ai-ops-hero-sub">{t('ai.heroSub')}</p>
              <div className="ai-ops-landing-composer">{composer}</div>
            </motion.div>
          </div>
        ) : (
          <>
            <div className="ai-ops-messages">
              <div className="ai-ops-messages-inner">
                {messages.map((m, i) => (
                  <MessageBlock
                    key={`${activeId}-${i}`}
                    message={m}
                    pending={busy && i === messages.length - 1 && m.role === 'assistant' && !m.content}
                  />
                ))}
                <div ref={bottomRef} />
              </div>
            </div>
            <div className="ai-ops-dock">{composer}</div>
          </>
        )}
      </section>
    </div>
  )
}

function BrandTitle({ skipMotion }: { skipMotion: boolean }) {
  const brand = 'CILIKUBE'
  return (
    <div className="ai-ops-brand">
      <div className="ai-ops-brand-row" aria-label="CiliKube AI">
        <span className="ai-ops-brand-main">
          {brand.split('').map((ch, i) => (
            <motion.span
              key={`${ch}-${i}`}
              className="ai-ops-brand-letter"
              initial={skipMotion ? false : { opacity: 0, y: 18, filter: 'blur(6px)' }}
              animate={{ opacity: 1, y: 0, filter: 'blur(0px)' }}
              transition={{
                delay: skipMotion ? 0 : 0.04 * i,
                duration: 0.4,
                ease: [0.22, 1, 0.36, 1],
              }}
            >
              {ch}
            </motion.span>
          ))}
        </span>
        <motion.span
          className="ai-ops-brand-ai"
          initial={skipMotion ? false : { opacity: 0, scale: 0.88, filter: 'blur(8px)' }}
          animate={{ opacity: 1, scale: 1, filter: 'blur(0px)' }}
          transition={{ delay: skipMotion ? 0 : 0.34, duration: 0.5, ease: [0.22, 1, 0.36, 1] }}
        >
          <span className="ai-ops-brand-ai-corners" aria-hidden>
            <i />
            <i />
            <i />
            <i />
          </span>
          <span className="ai-ops-brand-ai-grid" aria-hidden />
          <span className="ai-ops-brand-ai-core">
            <span className="ai-ops-brand-ai-glyph">A</span>
            <span className="ai-ops-brand-ai-glyph">I</span>
          </span>
          <span className="ai-ops-brand-ai-beam" aria-hidden />
          <span className="ai-ops-brand-ai-sheen" aria-hidden />
        </motion.span>
      </div>
      <motion.div
        className="ai-ops-brand-rule"
        initial={skipMotion ? false : { scaleX: 0, opacity: 0 }}
        animate={{ scaleX: 1, opacity: 1 }}
        transition={{ delay: skipMotion ? 0 : 0.52, duration: 0.55, ease: [0.22, 1, 0.36, 1] }}
      >
        <span className="ai-ops-brand-rule-packet" aria-hidden />
      </motion.div>
    </div>
  )
}

function MessageBlock({ message, pending }: { message: AiChatMessage; pending?: boolean }) {
  const user = message.role === 'user'
  return (
    <div className={cn('ai-ops-msg', user ? 'is-user' : 'is-bot')}>
      {!user ? (
        <div className="ai-ops-avatar">
          <Bot className="h-3.5 w-3.5" />
        </div>
      ) : (
        <div className="ai-ops-avatar is-user" aria-hidden>
          你
        </div>
      )}
      <div className="ai-ops-bubble">
        {message.tools?.length ? (
          <div className="ai-ops-tools">
            {message.tools.map((t, ti) => (
              <div key={ti} className="ai-ops-tool">
                {t.ok === undefined ? (
                  <Loader2 className="h-3 w-3 animate-spin" />
                ) : (
                  <span className={t.ok ? 'ok' : 'err'}>{t.ok ? 'OK' : 'ERR'}</span>
                )}
                <code>{t.name}</code>
              </div>
            ))}
          </div>
        ) : null}
        <div className="ai-ops-text">
          {message.content || (pending ? <span className="ai-ops-typing">正在查看集群…</span> : '')}
        </div>
        {message.resources?.length ? (
          <div className="ai-ops-cards">
            {message.resources.map((r) => (
              <ResourceCard key={r.href + (r.console || '')} refItem={r} />
            ))}
          </div>
        ) : null}
      </div>
    </div>
  )
}

function ResourceCard({ refItem }: { refItem: AiResourceRef }) {
  const isConsole = Boolean(refItem.console)
  return (
    <Link to={refItem.href} className="ai-ops-card">
      {isConsole ? <Terminal className="h-3.5 w-3.5 shrink-0" /> : <Zap className="h-3.5 w-3.5 shrink-0" />}
      <span className="truncate">{refItem.label || refItem.name}</span>
    </Link>
  )
}

function mergeRefs(a: AiResourceRef[], b: AiResourceRef[]): AiResourceRef[] {
  const seen = new Set(a.map((r) => r.href + (r.console || '')))
  const out = [...a]
  for (const r of b) {
    const k = r.href + (r.console || '')
    if (!r.href || seen.has(k)) continue
    seen.add(k)
    out.push(r)
  }
  return out
}
