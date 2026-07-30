import { useEffect, useId, useState, type FormEvent } from 'react'
import { X } from 'lucide-react'
import type { AiSkillDef, CustomSkillInput } from '@/lib/aiSkills'

type Props = {
  open: boolean
  initial?: AiSkillDef | null
  onClose: () => void
  onSave: (input: CustomSkillInput & { id?: string }) => void
}

export function AiSkillEditor({ open, initial, onClose, onSave }: Props) {
  const titleId = useId()
  const [label, setLabel] = useState('')
  const [blurb, setBlurb] = useState('')
  const [prompt, setPrompt] = useState('')

  useEffect(() => {
    if (!open) return
    setLabel(initial?.label || '')
    setBlurb(initial?.custom ? initial.blurb || '' : '')
    setPrompt(initial?.prompt || '')
  }, [open, initial])

  useEffect(() => {
    if (!open) return
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [open, onClose])

  if (!open) return null

  const submit = (e: FormEvent) => {
    e.preventDefault()
    onSave({
      id: initial?.custom ? initial.id : undefined,
      label,
      blurb,
      prompt,
    })
  }

  const canSave = label.trim().length > 0 && prompt.trim().length > 0

  return (
    <div className="ai-ops-skill-editor-backdrop" role="presentation" onClick={onClose}>
      <form
        className="ai-ops-skill-editor"
        role="dialog"
        aria-modal="true"
        aria-labelledby={titleId}
        onClick={(e) => e.stopPropagation()}
        onSubmit={submit}
      >
        <div className="ai-ops-skill-editor-head">
          <div>
            <div className="ai-ops-kicker">Skill</div>
            <h2 id={titleId} className="ai-ops-skill-editor-title">
              {initial?.custom ? '编辑自定义 Skill' : '创建自定义 Skill'}
            </h2>
          </div>
          <button type="button" className="ai-ops-icon-btn" onClick={onClose} aria-label="关闭">
            <X className="h-4 w-4" />
          </button>
        </div>

        <label className="ai-ops-skill-editor-field">
          <span>名称</span>
          <input
            autoFocus
            value={label}
            maxLength={24}
            placeholder="例如：夜班点检"
            onChange={(e) => setLabel(e.target.value)}
          />
        </label>

        <label className="ai-ops-skill-editor-field">
          <span>简介（可选）</span>
          <input
            value={blurb}
            maxLength={48}
            placeholder="一句话说明这个 Skill 做什么"
            onChange={(e) => setBlurb(e.target.value)}
          />
        </label>

        <label className="ai-ops-skill-editor-field">
          <span>Prompt</span>
          <textarea
            value={prompt}
            maxLength={2000}
            rows={6}
            placeholder="写下你希望调查员执行的完整指令…"
            onChange={(e) => setPrompt(e.target.value)}
          />
        </label>

        <div className="ai-ops-skill-editor-foot">
          <p className="ai-ops-skill-editor-note">保存在本机浏览器，可随时删除</p>
          <div className="ai-ops-skill-editor-actions">
            <button type="button" className="ai-ops-skill-editor-cancel" onClick={onClose}>
              取消
            </button>
            <button type="submit" className="ai-ops-send" disabled={!canSave}>
              保存
            </button>
          </div>
        </div>
      </form>
    </div>
  )
}
