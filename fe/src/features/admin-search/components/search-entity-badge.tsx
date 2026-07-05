import { Badge } from '@/components/ui/badge'

import { entityTypeLabel } from '../lib/search-routes'
import type { SearchEntityType } from '../types'

type SearchEntityBadgeProps = {
  entityType: SearchEntityType
}

export function SearchEntityBadge({ entityType }: SearchEntityBadgeProps) {
  return (
    <Badge variant="secondary" className="shrink-0 capitalize">
      {entityTypeLabel(entityType)}
    </Badge>
  )
}
