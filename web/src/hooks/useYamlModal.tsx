import { useState, type ReactNode } from 'react'
import { ResourceYamlModal } from '@/components/ResourceYamlModal'
import { YamlActionButton } from '@/components/YamlActionButton'

export function useYamlModal(resource: string, namespaced = true) {
  const [target, setTarget] = useState<any | null>(null)

  const yamlButton = (item: any) => (
    <YamlActionButton resource={resource} onClick={() => setTarget(item)} />
  )

  const yamlModal: ReactNode = (
    <ResourceYamlModal
      open={Boolean(target)}
      item={target}
      resource={resource}
      namespaced={namespaced}
      onClose={() => setTarget(null)}
    />
  )

  return { yamlButton, yamlModal }
}
