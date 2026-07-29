import { apiGet } from '@/lib/api'

export type AiStatus = {
  enabled: boolean
  ready: boolean
  provider: string
  model: string
}

export type AiResourceRef = {
  kind: string
  namespace?: string
  name: string
  href: string
  console?: string
  label?: string
}

export type AiChatMessage = {
  role: 'user' | 'assistant'
  content: string
  resources?: AiResourceRef[]
  tools?: { name: string; ok?: boolean; content?: string }[]
}

export type AiStreamHandlers = {
  onToolCall?: (name: string, args: unknown) => void
  onToolResult?: (name: string, ok: boolean, content: string) => void
  onResources?: (items: AiResourceRef[]) => void
  onMessage?: (content: string) => void
  onError?: (message: string) => void
  onDone?: () => void
}

export function getAiStatus() {
  return apiGet<AiStatus>('/api/v1/ai/status')
}

function authToken(): string {
  try {
    return localStorage.getItem('cilikube_token') || ''
  } catch {
    return ''
  }
}

function clusterId(): string {
  try {
    return localStorage.getItem('cilikube_cluster') || ''
  } catch {
    return ''
  }
}

/** SSE chat — uses fetch so Authorization header works. */
export async function streamAiChat(
  messages: { role: string; content: string }[],
  opts: { namespace?: string; signal?: AbortSignal },
  handlers: AiStreamHandlers,
): Promise<void> {
  const qs = new URLSearchParams()
  const cid = clusterId()
  if (cid) qs.set('clusterId', cid)
  const url = `/api/v1/ai/chat${qs.toString() ? `?${qs}` : ''}`

  const res = await fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${authToken()}`,
    },
    body: JSON.stringify({
      messages,
      namespace: opts.namespace || '',
    }),
    signal: opts.signal,
  })

  if (!res.ok) {
    let msg = `AI chat failed (${res.status})`
    try {
      const j = await res.json()
      msg = j?.message || j?.error || msg
    } catch {
      /* ignore */
    }
    handlers.onError?.(msg)
    return
  }

  const reader = res.body?.getReader()
  if (!reader) {
    handlers.onError?.('No response stream')
    return
  }

  const decoder = new TextDecoder()
  let buf = ''
  let eventName = 'message'

  while (true) {
    const { done, value } = await reader.read()
    if (done) break
    buf += decoder.decode(value, { stream: true })
    const chunks = buf.split('\n\n')
    buf = chunks.pop() || ''
    for (const chunk of chunks) {
      const lines = chunk.split('\n')
      let dataLine = ''
      for (const line of lines) {
        if (line.startsWith('event:')) eventName = line.slice(6).trim()
        if (line.startsWith('data:')) dataLine += line.slice(5).trim()
      }
      if (!dataLine) continue
      let data: any = {}
      try {
        data = JSON.parse(dataLine)
      } catch {
        continue
      }
      switch (eventName) {
        case 'tool_call':
          handlers.onToolCall?.(data.name, data.arguments)
          break
        case 'tool_result':
          handlers.onToolResult?.(data.name, Boolean(data.ok), String(data.content || ''))
          break
        case 'resources':
          handlers.onResources?.(data.items || [])
          break
        case 'message':
          handlers.onMessage?.(String(data.content || ''))
          break
        case 'error':
          handlers.onError?.(String(data.message || 'AI error'))
          break
        case 'done':
          handlers.onDone?.()
          break
      }
      eventName = 'message'
    }
  }
  handlers.onDone?.()
}
