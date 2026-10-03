import { useParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { apiGet } from '@/lib/api'
import { useCluster } from '@/store/cluster'
import { PageHeader } from '@/components/ui'

export type PluginManifest = {
  id: string
  name: string
  title?: string
  titleZh?: string
  resources?: string[]
}

export function PluginFramePage() {
  const { id = '' } = useParams()
  const { clusterId } = useCluster()
  const q = useQuery({
    queryKey: ['plugins'],
    queryFn: () => apiGet<PluginManifest[]>('/api/v1/plugins'),
  })
  const plugin = (q.data || []).find((item) => item.id === id)
  const title = plugin?.title || plugin?.name || id
  const src = `/plugin-assets/${encodeURIComponent(id)}/index.html?clusterId=${encodeURIComponent(clusterId)}`
  return (
    <div className="flex h-full min-h-0 flex-col gap-3">
      <PageHeader title={title} subtitle={id} />
      <iframe title={title} src={src} className="min-h-[480px] flex-1 rounded border border-line bg-panel" />
    </div>
  )
}
