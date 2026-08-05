import { useEffect, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { fetchShowcaseInfo } from '@/api/showcase'
import { APP_REPO_URL } from '@/lib/version'
import { cn } from '@/lib/utils'

const LEGACY_DISMISS_KEY = 'cilikube_star_cta_dismissed_until'

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
        'star-support-rotator inline-flex max-w-md items-center text-[13px] leading-snug text-cyan transition-opacity duration-300 hover:underline',
        visible ? 'opacity-100' : 'opacity-0',
        className,
      )}
    >
      <span>{lines[idx]}</span>
    </a>
  )
}

/** Floating tip — close only hides for this page view; every fresh open shows again. */
export function StarSupportFloat() {
  const { t } = useTranslation()
  const [show, setShow] = useState(false)
  const [showcase, setShowcase] = useState(false)

  useEffect(() => {
    // Drop the old 14-day localStorage dismiss so prior closes don't stick.
    try {
      localStorage.removeItem(LEGACY_DISMISS_KEY)
    } catch {
      /* ignore */
    }

    let cancelled = false
    ;(async () => {
      try {
        const info = await fetchShowcaseInfo()
        if (cancelled) return
        // Local / self-hosted: never show the floating star CTA.
        if (!info.showcase) {
          setShowcase(false)
          setShow(false)
          return
        }
        setShowcase(true)
        window.setTimeout(() => {
          if (!cancelled) setShow(true)
        }, 1200)
      } catch {
        // Unknown mode — prefer quiet local UX over a surprise promo card.
        if (!cancelled) {
          setShowcase(false)
          setShow(false)
        }
      }
    })()
    return () => {
      cancelled = true
    }
  }, [])

  if (!show || !showcase) return null

  return (
    <div
      className={cn(
        'star-support-float pointer-events-auto fixed z-[var(--z-float-cta)] w-[min(22rem,calc(100vw-1.5rem))]',
        'bottom-[max(1rem,env(safe-area-inset-bottom))] right-[max(0.75rem,env(safe-area-inset-right))]',
        'rounded border border-cyan/35 bg-panel-solid/95 px-3.5 py-3 shadow-[0_12px_40px_rgba(0,0,0,0.35)] backdrop-blur-md',
      )}
      role="dialog"
      aria-label={t('cta.starAction')}
    >
      <div className="mb-1.5 flex items-center justify-end gap-2">
        <button
          type="button"
          className="text-[10px] tracking-[0.12em] text-text-dim uppercase hover:text-text"
          onClick={() => setShow(false)}
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
          onClick={() => setShow(false)}
        >
          ★ {t('cta.starAction')}
        </a>
        <button
          type="button"
          className="text-[11px] text-text-dim hover:text-text hover:underline"
          onClick={() => setShow(false)}
        >
          {t('cta.later')}
        </button>
      </div>
    </div>
  )
}
