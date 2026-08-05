import { useTranslation } from 'react-i18next'
import { HudSelect } from '@/components/ui'
import { useCluster } from '@/store/cluster'
import { ALL_NAMESPACES, useNamespace } from '@/store/namespace'
import { cn } from '@/lib/utils'

type Layout = 'inline' | 'stack' | 'grid'

type Props = {
  /** inline = topbar; stack = drawer; grid = main content strip below lg */
  layout: Layout
  /** Show CLUSTER / NAMESPACE text labels (topbar hides them until xl via CSS). */
  showLabels?: boolean
  className?: string
  /** Extra class on each select wrapper */
  selectClassName?: string
}

/**
 * Single source of truth for cluster + namespace pickers.
 * Breakpoint contract (AppShell):
 * - < lg: main strip (`grid`) + mobile drawer (`stack`)
 * - >= lg: topbar (`inline`) only — main strip hidden
 */
export function ClusterNamespaceControls({
  layout,
  showLabels = true,
  className,
  selectClassName,
}: Props) {
  const { t } = useTranslation()
  const { clusters, clusterId, setClusterId, switching } = useCluster()
  const { namespace, setNamespace, namespaces } = useNamespace()

  const clusterLabel = t('nav.cluster')
  const nsLabel = t('common.namespace')

  const clusterOptions = clusters.map((c) => ({
    value: c.id || c.name,
    label: c.name || c.id,
  }))
  const nsOptions = [
    { value: ALL_NAMESPACES, label: t('common.allNamespaces') },
    ...namespaces.map((n) => ({ value: n, label: n })),
  ]

  const clusterSelect = (
    <HudSelect
      aria-label={clusterLabel}
      className={cn(
        layout === 'inline' &&
          'app-topbar-select w-auto min-w-[7.5rem] max-w-[9.5rem] xl:min-w-[8.5rem] xl:max-w-[12rem]',
        layout !== 'inline' && 'w-full',
        selectClassName,
      )}
      value={clusterId}
      onChange={setClusterId}
      disabled={switching || !clusters.length}
      options={clusterOptions}
    />
  )

  const nsSelect = (
    <HudSelect
      aria-label={nsLabel}
      className={cn(
        layout === 'inline' &&
          'app-topbar-select w-auto min-w-[7rem] max-w-[9rem] xl:min-w-[7.5rem] xl:max-w-[11rem]',
        layout !== 'inline' && 'w-full',
        selectClassName,
      )}
      value={namespace}
      onChange={setNamespace}
      searchableWhen={0}
      options={nsOptions}
    />
  )

  if (layout === 'inline') {
    return (
      <div className={cn('app-topbar-context', className)}>
        <label className="app-topbar-field flex items-center gap-2">
          {showLabels ? <span className="hud-label app-topbar-context-label">{clusterLabel}</span> : null}
          {clusterSelect}
        </label>
        <label className="app-topbar-field flex items-center gap-2">
          {showLabels ? <span className="hud-label app-topbar-context-label">{nsLabel}</span> : null}
          {nsSelect}
        </label>
      </div>
    )
  }

  if (layout === 'grid') {
    return (
      <div className={cn('grid grid-cols-2 gap-2', className)}>
        <div className="min-w-0">
          {showLabels ? <div className="hud-label mb-1">{clusterLabel}</div> : null}
          {clusterSelect}
        </div>
        <div className="min-w-0">
          {showLabels ? <div className="hud-label mb-1">{nsLabel}</div> : null}
          {nsSelect}
        </div>
      </div>
    )
  }

  return (
    <div className={cn('space-y-3', className)}>
      <div>
        {showLabels ? <div className="hud-label mb-2">{clusterLabel}</div> : null}
        {clusterSelect}
      </div>
      <div>
        {showLabels ? <div className="hud-label mb-2">{nsLabel}</div> : null}
        {nsSelect}
      </div>
    </div>
  )
}
