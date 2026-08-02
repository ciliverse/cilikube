import { useEffect, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { fetchShowcaseInfo } from '@/api/showcase'
import { APP_REPO_URL } from '@/lib/version'
import { cn } from '@/lib/utils'

const DISMISS_KEY = 'cilikube_star_cta_dismissed_until'
const DISMISS_DAYS = 14

function dismissedUntil(): number {
  const raw = localStorage.getItem(DISMISS_KEY)
  const n = raw ? Number(raw) : 0
  return Number.isFinite(n) ? n : 0
}

function isDismissed() {
  return Date.now() < dismissedUntil()
}

function dismissForDays(days: number) {
  localStorage.setItem(DISMISS_KEY, String(Date.now() + days * 86_400_000))
}

/** Rotating support line for login brand / document-adjacent copy. */
export function StarSupportRotator({ className }: { className?: string }) {
  const { t, i18n } = useTranslation()
  const lines = useMemo(
    () => [t('cta.starLine1'), t('cta.starLine2'), t('cta.starLine3')],
    [t, i18n.language],
  )
  const [idx, setIdx] = useState(0)
  const [visible, setVisible] = useState(true)

  useEffect(() => {
    const id = window.setInterval(() => {
      setVisible(false)
      window.setTimeout(() => {
        setIdx((i) => (i + 1) % lines.length)
        setVisible(true)
      }, 280)
    }, 4200)
    return () => window.clearInterval(id)
  }, [lines.length])

  return (
    <a
      href={APP_REPO_URL}
      target="_blank"
      rel="noreferrer noopener"
      className={cn(
        'star-support-rotator inline-flex max-w-md items-center gap-2 text-[13px] leading-snug text-cyan transition-opacity duration-300 hover:underline',
        visible ? 'opacity-100' : 'opacity-0',
        className,
      )}
    >
      <span className="hud-label shrink-0 text-cyan/80">{t('cta.badge')}</span>
      <span>{lines[idx]}</span>
    </a>
  )
}

/** Dismissible floating tip — prefer showcase; still polite on self-host. */
export function StarSupportFloat() {
  const { t } = useTranslation()
  const [show, setShow] = useState(false)
  const [showcase, setShowcase] = useState(false)

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      if (isDismissed()) return
      try {
        const info = await fetchShowcaseInfo()
        if (cancelled) return
        setShowcase(Boolean(info.showcase))
        // Showcase: show after a short beat. Self-host: only if never dismissed this session window.
        window.setTimeout(() => {
          if (!cancelled && !isDismissed()) setShow(true)
        }, info.showcase ? 1200 : 4000)
      } catch {
        if (!cancelled && !isDismissed()) {
          window.setTimeout(() => {
            if (!cancelled && !isDismissed()) setShow(true)
          }, 5000)
        }
      }
    })()
    return () => {
      cancelled = true
    }
  }, [])

  if (!show) return null

  return (
    <div
      className={cn(
        'star-support-float pointer-events-auto fixed z-40 w-[min(22rem,calc(100vw-1.5rem))]',
        'bottom-[max(1rem,env(safe-area-inset-bottom))] right-[max(0.75rem,env(safe-area-inset-right))]',
        'rounded border border-cyan/35 bg-panel-solid/95 px-3.5 py-3 shadow-[0_12px_40px_rgba(0,0,0,0.35)] backdrop-blur-md',
      )}
      role="dialog"
      aria-label={t('cta.badge')}
    >
      <div className="mb-1.5 flex items-center justify-between gap-2">
        <span className="hud-label text-cyan">{t('cta.badge')}</span>
        <button
          type="button"
          className="text-[10px] tracking-[0.12em] text-text-dim uppercase hover:text-text"
          onClick={() => {
            dismissForDays(DISMISS_DAYS)
            setShow(false)
          }}
        >
          {t('common.close')}
        </button>
      </div>
      <p className="text-[13px] leading-relaxed text-text">
        {showcase ? t('cta.starMessageShowcase') : t('cta.starMessage')}
      </p>
      <div className="mt-2.5 flex flex-wrap items-center gap-3">
        <a
          href={APP_REPO_URL}
          target="_blank"
          rel="noreferrer noopener"
          className="inline-flex items-center gap-1.5 rounded border border-cyan/40 bg-cyan/10 px-2.5 py-1 text-[11px] font-semibold tracking-[0.12em] text-cyan uppercase hover:bg-cyan/15"
          onClick={() => {
            dismissForDays(DISMISS_DAYS)
            setShow(false)
          }}
        >
          ★ {t('cta.starAction')}
        </a>
        <button
          type="button"
          className="text-[11px] text-text-dim hover:text-text hover:underline"
          onClick={() => {
            dismissForDays(DISMISS_DAYS)
            setShow(false)
          }}
        >
          {t('cta.later')}
        </button>
      </div>
    </div>
  )
}
