import { useEffect, useRef, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import 'xterm/css/xterm.css'
import { useTranslation } from 'react-i18next'
import { getClusterId, getToken } from '@/lib/api'
import { Button, PageHeader } from '@/components/ui'
import { useAuth } from '@/store/auth'
import { getStoredThemeId, resolveTheme, subscribeTheme, toXtermTheme } from '@/theme/themes'

function shellUrl(kind: 'kubectl' | 'node', nodeName: string) {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const qs = new URLSearchParams()
  const token = getToken()
  const clusterId = getClusterId()
  if (token) qs.set('token', token)
  if (clusterId) qs.set('clusterId', clusterId)
  const path =
    kind === 'node'
      ? `/api/v1/nodes/${encodeURIComponent(nodeName)}/shell`
      : '/api/v1/shell/kubectl'
  return `${protocol}//${window.location.host}${path}?${qs.toString()}`
}

/** Admin-only cluster shell: kubectl on the API host, or a short-lived node shell pod. */
export function ClusterShellPage() {
  const { t } = useTranslation()
  const { isAdmin } = useAuth()
  const [params] = useSearchParams()
  const initialNode = params.get('node') || ''
  const [nodeName, setNodeName] = useState(initialNode)
  const [status, setStatus] = useState<'idle' | 'connecting' | 'open' | 'closed' | 'error'>('idle')
  const hostRef = useRef<HTMLDivElement>(null)
  const termRef = useRef<Terminal | null>(null)
  const fitRef = useRef<FitAddon | null>(null)
  const wsRef = useRef<WebSocket | null>(null)
  const kindRef = useRef<'kubectl' | 'node' | null>(null)

  useEffect(() => {
    if (!hostRef.current || termRef.current) return
    const term = new Terminal({
      cursorBlink: true,
      convertEol: true,
      fontFamily:
        getComputedStyle(document.documentElement).getPropertyValue('--font-mono').trim() ||
        'Maple Mono, Menlo, Consolas, monospace',
      fontSize: 13,
      theme: toXtermTheme(resolveTheme(getStoredThemeId()).terminal),
    })
    const fit = new FitAddon()
    term.loadAddon(fit)
    term.open(hostRef.current)
    fit.fit()
    term.writeln(t('shell.hint'))
    termRef.current = term
    fitRef.current = fit
    const onResize = () => {
      try {
        fit.fit()
      } catch {
        /* ignore */
      }
    }
    window.addEventListener('resize', onResize)
    const dataDisp = term.onData((data) => {
      if (wsRef.current?.readyState === WebSocket.OPEN) wsRef.current.send(data)
    })
    const resizeDisp = term.onResize(({ cols, rows }) => {
      if (kindRef.current !== 'kubectl' || wsRef.current?.readyState !== WebSocket.OPEN) return
      wsRef.current.send(JSON.stringify({ type: 'resize', cols, rows }))
    })
    const unsubTheme = subscribeTheme(() => {
      term.options.theme = toXtermTheme(resolveTheme(getStoredThemeId()).terminal)
    })
    return () => {
      unsubTheme()
      window.removeEventListener('resize', onResize)
      dataDisp.dispose()
      resizeDisp.dispose()
      wsRef.current?.close(1000, 'unmount')
      term.dispose()
      termRef.current = null
    }
  }, [t])

  const disconnect = () => {
    wsRef.current?.close(1000, 'user')
    wsRef.current = null
    setStatus('closed')
  }

  const connect = (kind: 'kubectl' | 'node') => {
    if (kind === 'node' && !nodeName.trim()) return
    disconnect()
    setStatus('connecting')
    termRef.current?.clear()
    termRef.current?.writeln(
      kind === 'node' ? t('shell.connectingNode', { name: nodeName.trim() }) : t('shell.connectingKubectl'),
    )
    kindRef.current = kind
    const ws = new WebSocket(shellUrl(kind, nodeName.trim()))
    ws.binaryType = 'arraybuffer'
    wsRef.current = ws
    ws.onopen = () => {
      setStatus('open')
      try {
        fitRef.current?.fit()
      } catch {
        /* ignore */
      }
    }
    ws.onerror = () => setStatus('error')
    ws.onclose = () => setStatus((prev) => (prev === 'error' ? prev : 'closed'))
    ws.onmessage = (ev) => {
      if (typeof ev.data === 'string') termRef.current?.write(ev.data)
      else if (ev.data instanceof ArrayBuffer) termRef.current?.write(new Uint8Array(ev.data))
    }
  }

  return (
    <div className="flex h-full min-h-0 flex-col gap-3">
      <PageHeader title={t('shell.title')} subtitle={t('shell.subtitle')} />
      {!isAdmin ? (
        <p className="text-sm text-amber-300">{t('shell.adminRequired')}</p>
      ) : (
        <div className="flex flex-wrap items-end gap-2">
          <label className="flex flex-col gap-1 text-xs">
            <span className="hud-label">{t('shell.node')}</span>
            <input
              className="hud-field w-56 font-mono text-xs"
              value={nodeName}
              onChange={(e) => setNodeName(e.target.value)}
              placeholder={t('shell.nodePlaceholder')}
            />
          </label>
          <Button type="button" variant="outline" onClick={() => connect('kubectl')}>
            {t('shell.connectKubectl')}
          </Button>
          <Button type="button" variant="outline" disabled={!nodeName.trim()} onClick={() => connect('node')}>
            {t('shell.connectNode')}
          </Button>
          <Button type="button" variant="ghost" onClick={disconnect}>
            {t('shell.disconnect')}
          </Button>
          <span className="text-xs text-muted">{status}</span>
        </div>
      )}
      <div ref={hostRef} className="min-h-[420px] flex-1 overflow-hidden rounded border border-line bg-black/40" />
    </div>
  )
}
