import { useEffect, useRef, useState } from 'react'
import { buildPodWsUrl } from '@/api/resources'
import { Badge, Button, HudSelect } from '@/components/ui'

type Props = {
  namespace: string
  podName: string
  containers: Array<{ name: string; image: string }>
}

export function PodLogs({ namespace, podName, containers }: Props) {
  const [container, setContainer] = useState(containers[0]?.name || '')
  const [tailLines, setTailLines] = useState(200)
  const [timestamps, setTimestamps] = useState(true)
  const [follow, setFollow] = useState(true)
  const [logs, setLogs] = useState('')
  const [status, setStatus] = useState<'idle' | 'connecting' | 'open' | 'closed' | 'error'>(
    'idle',
  )
  const [bump, setBump] = useState(0)
  const preRef = useRef<HTMLPreElement>(null)
  const wsRef = useRef<WebSocket | null>(null)
  const autoScrollRef = useRef(true)

  useEffect(() => {
    if (!container && containers[0]?.name) setContainer(containers[0].name)
  }, [containers, container])

  useEffect(() => {
    if (!container || !namespace || !podName) return

    wsRef.current?.close(1000, 'reconnect')
    setLogs('')
    setStatus('connecting')

    const url = buildPodWsUrl('logs', namespace, podName, {
      container,
      tailLines,
      timestamps,
      follow,
    })
    const ws = new WebSocket(url)
    wsRef.current = ws

    ws.onopen = () => setStatus('open')
    ws.onmessage = (event) => {
      if (typeof event.data !== 'string') return
      setLogs((prev) => (prev ? `${prev}\n${event.data}` : event.data))
    }
    ws.onerror = () => {
      setStatus('error')
      setLogs((prev) => `${prev}\n# WebSocket error`)
    }
    ws.onclose = (event) => {
      setStatus(event.code === 1000 ? 'closed' : 'error')
      if (event.code !== 1000) {
        setLogs((prev) => `${prev}\n# disconnected: ${event.reason || event.code}`)
      }
    }

    return () => {
      ws.close(1000, 'cleanup')
    }
  }, [container, namespace, podName, tailLines, timestamps, follow, bump])

  useEffect(() => {
    if (!autoScrollRef.current || !preRef.current) return
    preRef.current.scrollTop = preRef.current.scrollHeight
  }, [logs])

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
            onChange={setContainer}
            searchableWhen={99}
            options={containers.map((c) => ({ value: c.name, label: c.name }))}
          />
        </label>
        <label className="flex items-center gap-2 text-xs">
          <span className="hud-label">Tail</span>
          <input
            type="number"
            min={10}
            max={10000}
            value={tailLines}
            onChange={(e) => setTailLines(Number(e.target.value) || 100)}
            className="hud-field w-24"
          />
        </label>
        <label className="flex items-center gap-2 text-xs text-text-dim">
          <input
            type="checkbox"
            checked={timestamps}
            onChange={(e) => setTimestamps(e.target.checked)}
          />
          Timestamps
        </label>
        <label className="flex items-center gap-2 text-xs text-text-dim">
          <input type="checkbox" checked={follow} onChange={(e) => setFollow(e.target.checked)} />
          Follow
        </label>
        <Button
          variant="outline"
          className="px-3 py-1.5 text-xs"
          type="button"
          onClick={() => setBump((n) => n + 1)}
        >
          Refresh
        </Button>
        <Button
          variant="ghost"
          className="px-3 py-1.5 text-xs"
          type="button"
          onClick={() => setLogs('')}
        >
          Clear
        </Button>
        <div className="ml-auto">
          <Badge tone={statusTone}>{status}</Badge>
        </div>
      </div>
      <pre
        ref={preRef}
        onScroll={(e) => {
          const el = e.currentTarget
          autoScrollRef.current = el.scrollTop + el.clientHeight >= el.scrollHeight - 24
        }}
        className="m-0 min-h-0 flex-1 overflow-auto bg-[#040a0e] px-4 py-3 font-mono text-[12px] leading-5 text-text whitespace-pre-wrap break-all"
      >
        {logs ||
          (status === 'connecting'
            ? '# connecting log stream…'
            : '# select a container to stream logs')}
      </pre>
    </div>
  )
}
