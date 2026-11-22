import { useTranslation } from 'react-i18next'
import { podProblem } from '@/lib/podReason'

const KNOWN = new Set([
  'CrashLoopBackOff',
  'ImagePullBackOff',
  'ErrImagePull',
  'CreateContainerConfigError',
  'CreateContainerError',
  'RunContainerError',
  'InvalidImageName',
  'OOMKilled',
  'Error',
  'ContainerCannotRun',
  'StartError',
])

export function PodReasonNote({
  pod,
  verbose = false,
}: {
  pod: Parameters<typeof podProblem>[0]
  verbose?: boolean
}) {
  const { t } = useTranslation()
  const { reason, message } = podProblem(pod)
  if (!reason) return null
  const explain = KNOWN.has(reason)
    ? t(`podReason.${reason}`)
    : t('podReason.generic', { reason })
  const title = message ? `${explain} — ${message}` : explain
  if (!verbose) {
    return (
      <div className="mt-1 max-w-[18rem] truncate font-mono text-[11px] text-warn" title={title}>
        {reason}
      </div>
    )
  }
  return (
    <p className="mt-2 max-w-xl text-xs leading-relaxed text-text-dim" title={title}>
      <span className="font-mono text-warn">{reason}</span>
      {' · '}
      {explain}
    </p>
  )
}
