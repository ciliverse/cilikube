import { FileCode2 } from 'lucide-react'
import { useAuth } from '@/store/auth'
import { Button } from '@/components/ui'

export function YamlActionButton({
  resource,
  onClick,
}: {
  resource: string
  onClick: () => void
}) {
  const { checkPermission } = useAuth()
  if (!checkPermission(resource, 'read')) return null
  return (
    <Button variant="ghost" className="px-2 py-1" type="button" title="YAML" onClick={onClick}>
      <FileCode2 className="h-3.5 w-3.5" />
    </Button>
  )
}
