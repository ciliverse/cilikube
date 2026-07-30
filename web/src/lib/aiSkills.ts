/** Lightweight Agent + Skill catalog for the AI landing page (not a skill marketplace). */

export type AiSkillGroup = 'combo' | 'inspect' | 'incident' | 'navigate' | 'custom'

export type AiAgentMeta = {
  id: string
  name: string
  blurb: string
}

export type AiSkillDef = {
  id: string
  code: string
  label: string
  blurb: string
  prompt: string
  group: AiSkillGroup
  /** User-authored skill persisted in localStorage. */
  custom?: boolean
  /** Display-only hint of tools the skill tends to exercise. */
  toolsHint?: string[]
}

export type CustomSkillInput = {
  label: string
  prompt: string
  blurb?: string
}

const CUSTOM_STORAGE_KEY = 'cilikube_ai_custom_skills_v1'
const MAX_CUSTOM_SKILLS = 30
const MAX_LABEL = 24
const MAX_BLURB = 48
const MAX_PROMPT = 2000

function uid() {
  return `skill_custom_${Date.now().toString(36)}_${Math.random().toString(36).slice(2, 7)}`
}

export const DEFAULT_AI_AGENT: AiAgentMeta = {
  id: 'cluster_investigator',
  name: '集群调查员',
  blurb: '只读查证 · 问清楚再进控制台',
}

export const AI_SKILLS: AiSkillDef[] = [
  {
    id: 'skill_inspect_combo',
    code: 'S1',
    label: '智能点检',
    blurb: '概览 + 异常 Pod + 近期事件',
    prompt:
      '请做一次智能点检：先给集群概览，再列出 Failed 或 Pending 的 Pod，并汇总最近 Warning / Error 事件。用中文简洁结论。',
    group: 'combo',
    toolsHint: ['get_cluster_overview', 'list_resources'],
  },
  {
    id: 'skill_triage_combo',
    code: 'S2',
    label: '快速调查',
    blurb: '异常 Pod → 抽样日志',
    prompt:
      '快速调查：找出 Failed 或 Pending 的 Pod，挑一个有代表性的看最近日志，并告诉我下一步该进控制台看哪里。',
    group: 'combo',
    toolsHint: ['list_resources', 'get_pod_logs'],
  },
  {
    id: 'skill_workload_snap',
    code: 'S3',
    label: '工作负载快照',
    blurb: 'Deployments + Services',
    prompt: '给当前关注范围做工作负载快照：列出 Deployments 和 Services（优先 default 命名空间）。',
    group: 'combo',
    toolsHint: ['list_resources'],
  },
  {
    id: 'skill_cluster_pulse',
    code: '01',
    label: '集群脉诊',
    blurb: '节点 · 负载健康',
    prompt: '集群现在怎么样？节点和关键负载健康吗？',
    group: 'inspect',
    toolsHint: ['get_cluster_overview'],
  },
  {
    id: 'skill_fault_scan',
    code: '02',
    label: '故障扫描',
    blurb: 'Failed · Pending',
    prompt: '有哪些 Failed 或 Pending 的 Pod？',
    group: 'incident',
    toolsHint: ['list_resources'],
  },
  {
    id: 'skill_deploy_list',
    code: '03',
    label: '部署盘点',
    blurb: 'default Deployments',
    prompt: '列出 default 命名空间的 Deployments',
    group: 'navigate',
    toolsHint: ['list_resources'],
  },
  {
    id: 'skill_log_sample',
    code: '04',
    label: '日志取样',
    blurb: '抽样最近输出',
    prompt: '随便挑一个 Pod，看看最近日志',
    group: 'inspect',
    toolsHint: ['list_resources', 'get_pod_logs'],
  },
  {
    id: 'skill_events',
    code: '05',
    label: '事件热线',
    blurb: 'Warning · Error',
    prompt: '最近有哪些 Warning 或 Error 事件？',
    group: 'incident',
    toolsHint: ['list_resources'],
  },
  {
    id: 'skill_crashloop',
    code: '06',
    label: '崩坏回路',
    blurb: 'CrashLoop · ImagePull',
    prompt: '有没有 CrashLoopBackOff 或 ImagePullBackOff 的 Pod？',
    group: 'incident',
    toolsHint: ['list_resources'],
  },
  {
    id: 'skill_node_pressure',
    code: '07',
    label: '节点压力',
    blurb: 'CPU · 内存紧张度',
    prompt: '哪个节点资源最紧张？帮我看看节点状态',
    group: 'inspect',
    toolsHint: ['list_resources', 'get_cluster_overview'],
  },
  {
    id: 'skill_svc_map',
    code: '08',
    label: '服务地图',
    blurb: 'Services 列表',
    prompt: '列出 default 命名空间里的 Service',
    group: 'navigate',
    toolsHint: ['list_resources'],
  },
  {
    id: 'skill_restarts',
    code: '09',
    label: '重启异常',
    blurb: '高 restartCount',
    prompt: '有没有重启次数特别高的 Pod？',
    group: 'incident',
    toolsHint: ['list_resources'],
  },
  {
    id: 'skill_namespaces',
    code: '10',
    label: '命名空间',
    blurb: 'Namespaces 盘点',
    prompt: '当前集群有哪些命名空间？各自大概有多少工作负载？',
    group: 'navigate',
    toolsHint: ['list_resources', 'get_cluster_overview'],
  },
]

export const AI_SKILL_GROUP_LABEL: Record<AiSkillGroup, string> = {
  combo: '精选',
  inspect: '巡检',
  incident: '排障',
  navigate: '导航',
  custom: '自定义',
}

function readCustomSkills(): AiSkillDef[] {
  try {
    const raw = localStorage.getItem(CUSTOM_STORAGE_KEY)
    if (!raw) return []
    const parsed = JSON.parse(raw) as AiSkillDef[]
    if (!Array.isArray(parsed)) return []
    return parsed
      .filter((s) => s && typeof s.id === 'string' && typeof s.label === 'string' && typeof s.prompt === 'string')
      .map((s, i) => ({
        id: s.id,
        code: s.code || `C${i + 1}`,
        label: String(s.label).slice(0, MAX_LABEL),
        blurb: String(s.blurb || '自定义 Prompt').slice(0, MAX_BLURB),
        prompt: String(s.prompt).slice(0, MAX_PROMPT),
        group: 'custom' as const,
        custom: true,
      }))
  } catch {
    return []
  }
}

function writeCustomSkills(skills: AiSkillDef[]) {
  try {
    localStorage.setItem(CUSTOM_STORAGE_KEY, JSON.stringify(skills.slice(0, MAX_CUSTOM_SKILLS)))
  } catch {
    /* quota / private mode */
  }
}

export function listCustomSkills(): AiSkillDef[] {
  return readCustomSkills()
}

export function allAiSkills(): AiSkillDef[] {
  return [...AI_SKILLS, ...listCustomSkills()]
}

export function sanitizeCustomSkillInput(input: CustomSkillInput): CustomSkillInput | null {
  const label = input.label.trim().replace(/\s+/g, ' ').slice(0, MAX_LABEL)
  const prompt = input.prompt.trim().slice(0, MAX_PROMPT)
  const blurb = (input.blurb || '').trim().replace(/\s+/g, ' ').slice(0, MAX_BLURB)
  if (!label || !prompt) return null
  return { label, prompt, blurb: blurb || '自定义 Prompt' }
}

export function upsertCustomSkill(input: CustomSkillInput & { id?: string }): AiSkillDef | null {
  const clean = sanitizeCustomSkillInput(input)
  if (!clean) return null
  const existing = readCustomSkills()
  if (input.id) {
    const idx = existing.findIndex((s) => s.id === input.id)
    if (idx >= 0) {
      const next: AiSkillDef = {
        ...existing[idx],
        label: clean.label,
        blurb: clean.blurb || existing[idx].blurb,
        prompt: clean.prompt,
        group: 'custom',
        custom: true,
      }
      const list = [...existing]
      list[idx] = next
      writeCustomSkills(list)
      return next
    }
  }
  const skill: AiSkillDef = {
    id: uid(),
    code: `C${existing.length + 1}`,
    label: clean.label,
    blurb: clean.blurb || '自定义 Prompt',
    prompt: clean.prompt,
    group: 'custom',
    custom: true,
  }
  writeCustomSkills([skill, ...existing])
  return skill
}

export function deleteCustomSkill(id: string): boolean {
  const existing = readCustomSkills()
  const next = existing.filter((s) => s.id !== id)
  if (next.length === existing.length) return false
  // Re-number codes for stable display
  writeCustomSkills(next.map((s, i) => ({ ...s, code: `C${i + 1}` })))
  return true
}

export function skillsByGroup(skills: AiSkillDef[] = AI_SKILLS): { group: AiSkillGroup; label: string; skills: AiSkillDef[] }[] {
  const order: AiSkillGroup[] = ['combo', 'inspect', 'incident', 'navigate', 'custom']
  return order
    .map((group) => ({
      group,
      label: AI_SKILL_GROUP_LABEL[group],
      skills: skills.filter((s) => s.group === group),
    }))
    .filter((g) => g.group === 'custom' || g.skills.length > 0)
}

/** Slash token just before the caret, e.g. `/点检` or `/fault`. */
export function parseSlashToken(
  value: string,
  cursor: number,
): { start: number; query: string } | null {
  const before = value.slice(0, Math.max(0, cursor))
  const m = before.match(/(?:^|[\s\n])(\/[^\s]*)$/)
  if (!m) return null
  const token = m[1]
  const start = before.length - token.length
  return { start, query: token.slice(1) }
}

export function filterSkillsByQuery(skills: AiSkillDef[], query: string): AiSkillDef[] {
  const q = query.trim().toLowerCase()
  if (!q) return skills
  return skills.filter((s) => {
    const hay = [
      s.label,
      s.blurb,
      s.code,
      s.id,
      s.id.replace(/^skill_/, ''),
      AI_SKILL_GROUP_LABEL[s.group],
      ...(s.toolsHint || []),
    ]
      .join(' ')
      .toLowerCase()
    return hay.includes(q) || s.label.includes(query.trim())
  })
}
