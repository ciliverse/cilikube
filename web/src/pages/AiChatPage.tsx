import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { Link } from 'react-router-dom'
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

const PROBES = [
  { id: '01', label: '集群脉诊', hint: '节点 · 负载 · 命名空间', q: '集群现在怎么样？节点和关键负载健康吗？' },
  { id: '02', label: '故障扫描', hint: 'Failed · Pending', q: '有哪些 Failed 或 Pending 的 Pod？' },
  { id: '03', label: '部署盘点', hint: 'default Deployments', q: '列出 default 命名空间的 Deployments' },
  { id: '04', label: '日志取样', hint: '最近输出', q: '随便挑一个 Pod，看看最近日志' },
]

export function AiChatPage() {
  const { clusterId, activeCluster } = useCluster()
  const { namespace } = useNamespace()
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

  const send = async (text?: string) => {
    const content = (text ?? input).trim()
    if (!content) return
    if (busy) {
      setErr('正在生成回复，请先点「停止」再发送')
      return
    }
    if (!ready) {
      setErr(statusQ.isError ? '无法连接 AI 服务，请确认后端已启动' : 'AI 暂不可用，请到控制台 Settings → AI 检查配置')
      return
    }

    ensureSession()
    setErr('')
    setInput('')

    const history = [...messages, { role: 'user' as const, content }]
    setMessages(history)
    setBusy(true)

    const assistant: AiChatMessage = { role: 'assistant', content: '', tools: [], resources: [] }
    setMessages([...history, assistant])

    abortRef.current?.abort()
    const ac = new AbortController()
    abortRef.current = ac

    const patchAssistant = (fn: (m: AiChatMessage) => AiChatMessage) => {
      setMessages((prev) => {
        const next = [...prev]
        const last = next[next.length - 1]
        if (last?.role === 'assistant') next[next.length - 1] = fn(last)
        return next
      })
    }

    try {
      await streamAiChat(
        history.map((m) => ({ role: m.role, content: m.content })),
        {
          namespace: namespace === 'all' ? '' : namespace,
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
    } finally {
      setBusy(false)
    }
  }

  const composer = (
    <AiComposer
      value={input}
      onChange={setInput}
      onSend={(text) => {
        void send(typeof text === 'string' ? text : input)
      }}
      onStop={stop}
      busy={busy}
      ready={ready}
      namespaceLabel={namespace === 'all' ? '全部命名空间' : namespace}
      clusterLabel={clusterLabel}
      err={err}
      evidence={evidence}
      probes={PROBES}
      landing={isEmpty}
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
              aria-label={historyOpen ? '收起历史' : '打开历史'}
              aria-expanded={historyOpen}
            >
              {historyOpen && isDesktop ? <PanelLeftClose className="h-4 w-4" /> : <History className="h-4 w-4" />}
            </button>
            <span className={cn('ai-ops-status', ready ? 'is-on' : 'is-off')}>
              <i />
              {ready ? '在线' : '离线'}
            </span>
            <span className="ai-ops-cluster" title={clusterLabel}>
              {clusterLabel}
            </span>
          </div>
          <Link to="/overview" className="ai-ops-console-btn">
            <Terminal className="h-3.5 w-3.5" />
            控制台
          </Link>
        </div>

        {!ready ? (
          <div className="ai-ops-banner">AI 暂不可用，请到控制台 Settings → AI 检查配置。</div>
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
                先问清楚
                <em>再动手改</em>
              </h2>
              <p className="ai-ops-hero-sub">查状态、找故障、一点进控制台</p>
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
          initial={skipMotion ? false : { opacity: 0, scale: 0.85, x: -6 }}
          animate={{ opacity: 1, scale: 1, x: 0 }}
          transition={{ delay: skipMotion ? 0 : 0.38, duration: 0.45, ease: [0.22, 1, 0.36, 1] }}
        >
          AI
          <span className="ai-ops-brand-ai-scan" aria-hidden />
        </motion.span>
      </div>
      <motion.div
        className="ai-ops-brand-rule"
        initial={skipMotion ? false : { scaleX: 0, opacity: 0 }}
        animate={{ scaleX: 1, opacity: 1 }}
        transition={{ delay: skipMotion ? 0 : 0.5, duration: 0.55, ease: [0.22, 1, 0.36, 1] }}
      />
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
