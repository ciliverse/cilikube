import { useRef, type FormEvent, type KeyboardEvent } from 'react'
import { ArrowUp, Square } from 'lucide-react'
import type { AiResourceRef } from '@/api/ai'
import { cn } from '@/lib/utils'
import { Link } from 'react-router-dom'
import { Terminal, Zap } from 'lucide-react'

export type AiProbe = { id: string; label: string; q: string }

type Props = {
  value: string
  onChange: (value: string) => void
  onSend: (text?: string) => void
  onStop?: () => void
  busy?: boolean
  ready?: boolean
  namespaceLabel: string
  clusterLabel: string
  err?: string
  evidence?: AiResourceRef[]
  probes?: AiProbe[]
  landing?: boolean
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
  probes,
  landing,
}: Props) {
  const textareaRef = useRef<HTMLTextAreaElement>(null)
  const composingRef = useRef(false)
  const maxHeight = landing ? 220 : 180

  const resize = () => {
    const el = textareaRef.current
    if (!el || composingRef.current) return
    el.style.height = 'auto'
    el.style.height = `${Math.min(el.scrollHeight, maxHeight)}px`
  }

  const submit = (e?: FormEvent) => {
    e?.preventDefault()
    onSend(value)
  }

  const sendProbe = (q: string) => {
    // Always surface the prompt first so the click never feels dead.
    onChange(q)
    requestAnimationFrame(() => {
      const el = textareaRef.current
      if (el) {
        el.style.height = 'auto'
        el.style.height = `${Math.min(el.scrollHeight, maxHeight)}px`
      }
    })
    onSend(q)
  }

  const onKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    // Don't send while Chinese/Japanese IME is composing
    if (e.nativeEvent.isComposing || e.key === 'Process' || composingRef.current) return
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
          <span className="ai-ops-evidence-label">线索</span>
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

      <form className={cn('ai-ops-composer', landing && probes?.length && 'has-probes')} onSubmit={submit}>
        {landing && probes?.length ? (
          <div className="ai-ops-composer-probes">
            <div className="ai-ops-composer-probes-label">快速开始</div>
            <div className="ai-ops-composer-chips">
              {probes.map((p) => (
                <button
                  key={p.id}
                  type="button"
                  className="ai-ops-chip-btn"
                  disabled={busy}
                  onClick={() => sendProbe(p.q)}
                >
                  <span className="ai-ops-chip-id">{p.id}</span>
                  {p.label}
                </button>
              ))}
            </div>
          </div>
        ) : null}

        <div className="ai-ops-composer-main">
          <textarea
            ref={textareaRef}
            className="ai-ops-input"
            rows={landing ? 4 : 2}
            placeholder={ready ? '集群里发生了什么？直接问…' : 'AI 暂不可用，仍可先输入草稿…'}
            value={value}
            onCompositionStart={() => {
              composingRef.current = true
            }}
            onCompositionEnd={(e) => {
              composingRef.current = false
              onChange(e.currentTarget.value)
              requestAnimationFrame(resize)
            }}
            onChange={(e) => {
              onChange(e.target.value)
              resize()
            }}
            onKeyDown={onKeyDown}
          />
          <div className="ai-ops-composer-bar">
            <span className="ai-ops-hint">
              {namespaceLabel}
              <span className="ai-ops-hint-sep">·</span>
              {clusterLabel}
            </span>
            {busy ? (
              <button type="button" className="ai-ops-send is-stop" onClick={onStop}>
                <Square className="h-3.5 w-3.5" />
                停止
              </button>
            ) : (
              <button type="submit" className="ai-ops-send" disabled={!ready || !value.trim()}>
                <ArrowUp className="h-4 w-4" />
                <span>发送</span>
              </button>
            )}
          </div>
        </div>
      </form>
    </div>
  )
}
