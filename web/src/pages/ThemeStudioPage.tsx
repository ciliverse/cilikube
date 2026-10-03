import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Button, PageHeader } from '@/components/ui'
import { useTheme } from '@/theme/useTheme'
import { switchTheme } from '@/theme/switchTheme'
import {
  deleteCustomTheme,
  duplicateTheme,
  isCustomThemeId,
  saveCustomTheme,
} from '@/theme/customThemes'
import { notifyThemeListeners, type Theme, type ThemeColors } from '@/theme/themes'

const FIELDS: Array<{ key: keyof ThemeColors; labelKey: string }> = [
  { key: 'bg', labelKey: 'themeStudio.bg' },
  { key: 'bgPanelSolid', labelKey: 'themeStudio.panel' },
  { key: 'primary', labelKey: 'themeStudio.primary' },
  { key: 'primaryDim', labelKey: 'themeStudio.primaryDim' },
  { key: 'secondary', labelKey: 'themeStudio.secondary' },
  { key: 'text', labelKey: 'themeStudio.text' },
  { key: 'textDim', labelKey: 'themeStudio.textDim' },
  { key: 'green', labelKey: 'themeStudio.ok' },
  { key: 'red', labelKey: 'themeStudio.danger' },
]

export function ThemeStudioPage() {
  const { t } = useTranslation()
  const { builtins, customs, themeId } = useTheme()
  const [draft, setDraft] = useState<Theme | null>(null)
  const [msg, setMsg] = useState('')

  const startCopy = (source: Theme) => {
    setMsg('')
    setDraft(duplicateTheme(source))
  }

  const edit = (theme: Theme) => {
    setMsg('')
    setDraft({ ...theme, colors: { ...theme.colors } })
  }

  const save = () => {
    if (!draft) return
    const saved = saveCustomTheme(draft, notifyThemeListeners)
    switchTheme(saved.id)
    setDraft(saved)
    setMsg(t('themeStudio.saved'))
  }

  const remove = (id: string) => {
    deleteCustomTheme(id, notifyThemeListeners)
    if (themeId === id) switchTheme('paper')
    if (draft?.id === id) setDraft(null)
    setMsg(t('themeStudio.deleted'))
  }

  return (
    <div className="flex w-full min-w-0 flex-col gap-4">
      <PageHeader title={t('themeStudio.title')} subtitle={t('themeStudio.subtitle')} />
      {msg ? <p className="text-sm text-cyan">{msg}</p> : null}
      <div className="grid gap-4 lg:grid-cols-[16rem_1fr]">
        <div className="space-y-3">
          <section>
            <h2 className="mb-2 text-xs tracking-wide text-text-dim uppercase">{t('themeStudio.builtin')}</h2>
            <ul className="space-y-1">
              {builtins.map((theme) => (
                <li key={theme.id} className="flex items-center justify-between gap-2 rounded border border-line px-2 py-1.5">
                  <button type="button" className="text-left text-sm" onClick={() => switchTheme(theme.id)}>
                    {theme.name}
                  </button>
                  <Button type="button" variant="ghost" className="px-2 py-1 text-xs" onClick={() => startCopy(theme)}>
                    {t('themeStudio.duplicate')}
                  </Button>
                </li>
              ))}
            </ul>
          </section>
          <section>
            <h2 className="mb-2 text-xs tracking-wide text-text-dim uppercase">{t('themeStudio.custom')}</h2>
            {customs.length === 0 ? <p className="text-xs text-text-dim">{t('themeStudio.empty')}</p> : null}
            <ul className="space-y-1">
              {customs.map((theme) => (
                <li key={theme.id} className="flex items-center justify-between gap-2 rounded border border-line px-2 py-1.5">
                  <button type="button" className="text-left text-sm" onClick={() => edit(theme)}>
                    {theme.name}
                  </button>
                  <Button type="button" variant="ghost" className="px-2 py-1 text-xs" onClick={() => remove(theme.id)}>
                    {t('themeStudio.delete')}
                  </Button>
                </li>
              ))}
            </ul>
          </section>
        </div>
        {draft ? (
          <form
            className="space-y-3 rounded border border-line p-4"
            onSubmit={(e) => {
              e.preventDefault()
              save()
            }}
          >
            <label className="block space-y-1 text-sm">
              <span className="hud-label">{t('themeStudio.name')}</span>
              <input
                className="hud-field"
                value={draft.name}
                onChange={(e) => setDraft({ ...draft, name: e.target.value })}
              />
            </label>
            <label className="block space-y-1 text-sm">
              <span className="hud-label">{t('themeStudio.mode')}</span>
              <select
                className="hud-field"
                value={draft.mode}
                onChange={(e) => setDraft({ ...draft, mode: e.target.value as Theme['mode'] })}
              >
                <option value="dark">dark</option>
                <option value="light">light</option>
              </select>
            </label>
            <div className="grid gap-3 sm:grid-cols-2">
              {FIELDS.map((field) => (
                <label key={field.key} className="flex items-center justify-between gap-2 text-sm">
                  <span>{t(field.labelKey)}</span>
                  <input
                    type="color"
                    value={toHex(draft.colors[field.key] || '#000000')}
                    onChange={(e) =>
                      setDraft({
                        ...draft,
                        colors: { ...draft.colors, [field.key]: e.target.value },
                      })
                    }
                  />
                </label>
              ))}
            </div>
            <Button type="submit">{t('themeStudio.save')}</Button>
            {!isCustomThemeId(draft.id) ? null : (
              <p className="text-xs text-text-dim">{draft.id}</p>
            )}
          </form>
        ) : (
          <p className="text-sm text-text-dim">{t('themeStudio.pick')}</p>
        )}
      </div>
    </div>
  )
}

function toHex(value: string): string {
  if (value.startsWith('#') && (value.length === 7 || value.length === 4)) {
    if (value.length === 4) {
      return `#${value[1]}${value[1]}${value[2]}${value[2]}${value[3]}${value[3]}`
    }
    return value
  }
  return '#000000'
}
