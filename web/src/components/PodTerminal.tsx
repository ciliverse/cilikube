import { useEffect, useRef, useState } from 'react'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import 'xterm/css/xterm.css'
import { buildPodWsUrl } from '@/api/resources'
import { Badge, Button, HudSelect } from '@/components/ui'
import {
  getStoredThemeId,
  resolveTheme,
  subscribeTheme,
  toXtermTheme,
} from '@/theme/themes'

type Props = {
  namespace: string
  podName: string
  containers: Array<{ name: string; image: string }>
  /** exec opens a shell; attach attaches to the running process */
  mode?: 'exec' | 'attach'
}

export function PodTerminal({ namespace, podName, containers, mode = 'exec' }: Props) {
  const [container, setContainer] = useState(containers[0]?.name || '')
  const [shell, setShell] = useState('/bin/sh')
  const [status, setStatus] = useState<'idle' | 'connecting' | 'open' | 'closed' | 'error'>('idle')
  const hostRef = useRef<HTMLDivElement>(null)
  const termRef = useRef<Terminal | null>(null)
  const fitRef = useRef<FitAddon | null>(null)
  const wsRef = useRef<WebSocket | null>(null)

  useEffect(() => {
    if (!container && containers[0]?.name) setContainer(containers[0].name)
  }, [containers, container])

  useEffect(() => {
    if (!hostRef.current || termRef.current) return
    const initial = resolveTheme(getStoredThemeId())
    const term = new Terminal({
      cursorBlink: true,
      convertEol: true,
      fontFamily:
        getComputedStyle(document.documentElement).getPropertyValue('--font-mono').trim() ||
        'Maple Mono, Menlo, Consolas, monospace',
      fontSize: 13,
      theme: toXtermTheme(initial.terminal),
    })
    const fit = new FitAddon()
    term.loadAddon(fit)
    term.open(hostRef.current)
    fit.fit()
    term.writeln(
      mode === 'attach'
        ? 'Select container and click Attach.'
        : 'Select container and click Connect.',
    )
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
      if (wsRef.current?.readyState === WebSocket.OPEN) {
        wsRef.current.send(data)
      }
    })

    const unsubTheme = subscribeTheme(() => {
      const next = resolveTheme(getStoredThemeId())
      term.options.theme = toXtermTheme(next.terminal)
      try {
        term.refresh(0, Math.max(0, term.rows - 1))
      } catch {
        /* ignore */
      }
    })

    return () => {
      unsubTheme()
      window.removeEventListener('resize', onResize)
      dataDisp.dispose()
      wsRef.current?.close(1000, 'unmount')
      term.dispose()
      termRef.current = null
      fitRef.current = null
    }
  }, [mode])

  const disconnect = () => {
    wsRef.current?.close(1000, 'user')
    wsRef.current = null
    setStatus('closed')
  }

  const connect = () => {
    if (!container || !namespace || !podName) return
    disconnect()
    setStatus('connecting')
    termRef.current?.clear()
    termRef.current?.writeln(
      mode === 'attach'
        ? `Attaching to ${podName}/${container}…`
        : `Connecting to ${podName}/${container}…`,
    )

    const url =
      mode === 'attach'
        ? buildPodWsUrl('attach', namespace, podName, {
            container,
            tty: true,
            stdin: true,
            stdout: true,
            stderr: true,
          })
        : buildPodWsUrl('exec', namespace, podName, {
            container,
            command: shell || '/bin/sh',
            tty: true,
            stdin: true,
            stdout: true,
            stderr: true,
          })
    const ws = new WebSocket(url)
    ws.binaryType = 'arraybuffer'
    wsRef.current = ws

    ws.onopen = () => {
      setStatus('open')
      termRef.current?.writeln('\x1b[32mConnected.\x1b[0m')
      try {
        fitRef.current?.fit()
      } catch {
        /* ignore */
      }
      termRef.current?.focus()
    }
    ws.onmessage = async (event) => {
      const term = termRef.current
      if (!term) return
      if (typeof event.data === 'string') {
        term.write(event.data)
        return
      }
      if (event.data instanceof ArrayBuffer) {
        term.write(new Uint8Array(event.data))
        return
      }
      if (event.data instanceof Blob) {
        term.write(new Uint8Array(await event.data.arrayBuffer()))
      }
    }
    ws.onerror = () => {
      setStatus('error')
      termRef.current?.writeln('\r\n\x1b[31mWebSocket error\x1b[0m')
    }
    ws.onclose = (event) => {
      setStatus(event.code === 1000 ? 'closed' : 'error')
      termRef.current?.writeln(`\r\n\x1b[33mClosed: ${event.reason || event.code}\x1b[0m`)
      wsRef.current = null
    }
  }

  const statusTone =
    status === 'open' ? 'ok' : status === 'connecting' ? 'warn' : status === 'error' ? 'danger' : 'neutral'

  return (
    <div className="flex h-[min(62vh,70dvh)] min-h-[240px] flex-col sm:min-h-[360px]">
      <div className="flex flex-wrap items-center gap-3 border-b border-line px-4 py-3">
        <label className="flex items-center gap-2 text-xs">
          <span className="hud-label">Container</span>
          <HudSelect
            aria-label="Container"
            className="w-auto min-w-[140px]"
            value={container}
            disabled={status === 'open' || status === 'connecting'}
            onChange={setContainer}
            searchableWhen={99}
            options={containers.map((c) => ({ value: c.name, label: c.name }))}
          />
        </label>
        {mode === 'exec' ? (
          <label className="flex items-center gap-2 text-xs">
            <span className="hud-label">Shell</span>
            <input
              className="hud-field w-32"
              value={shell}
              disabled={status === 'open' || status === 'connecting'}
              onChange={(e) => setShell(e.target.value)}
            />
          </label>
        ) : null}
        {status === 'open' ? (
          <Button variant="danger" className="px-3 py-1.5 text-xs" type="button" onClick={disconnect}>
            Disconnect
          </Button>
        ) : (
          <Button
            variant="primary"
            className="px-3 py-1.5 text-xs"
            type="button"
            disabled={!container || status === 'connecting'}
            onClick={connect}
          >
            {mode === 'attach' ? 'Attach' : 'Connect'}
          </Button>
        )}
        <div className="ml-auto">
          <Badge tone={statusTone}>{status}</Badge>
        </div>
      </div>
      <div ref={hostRef} className="term-surface min-h-0 flex-1 px-2 py-2" />
    </div>
  )
}
