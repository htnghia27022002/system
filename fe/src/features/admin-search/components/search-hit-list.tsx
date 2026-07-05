'use client'

import Link from 'next/link'
import { ArrowRightIcon } from 'lucide-react'

import { Card, CardContent } from '@/components/ui/card'

import { useSearchNavigation } from '../hooks/use-search-navigation'
import type { SearchHit } from '../types'
import { SearchEntityBadge } from './search-entity-badge'

type SearchHitListProps = {
  hits: SearchHit[]
  emptyMessage: string
}

export function SearchHitList({ hits, emptyMessage }: SearchHitListProps) {
  const { resolveHitHref } = useSearchNavigation()

  if (hits.length === 0) {
    return (
      <p className="rounded-lg border border-dashed px-4 py-10 text-center text-sm text-muted-foreground">
        {emptyMessage}
      </p>
    )
  }

  return (
    <ul className="space-y-3">
      {hits.map((hit) => {
        const href = resolveHitHref(hit)
        const content = (
          <Card className="transition-colors hover:bg-muted/40">
            <CardContent className="flex items-start gap-3 p-4">
              <div className="min-w-0 flex-1 space-y-2">
                <div className="flex flex-wrap items-center gap-2">
                  <SearchEntityBadge entityType={hit.entityType} />
                  <h3 className="truncate font-medium">{hit.title}</h3>
                </div>
                {hit.snippet ? (
                  <p className="line-clamp-2 text-sm text-muted-foreground">{hit.snippet}</p>
                ) : null}
                {hit.metadata?.email ? (
                  <p className="text-xs text-muted-foreground">{hit.metadata.email}</p>
                ) : null}
                {hit.metadata?.key ? (
                  <p className="font-mono text-xs text-muted-foreground">{hit.metadata.key}</p>
                ) : null}
              </div>
              {href ? (
                <ArrowRightIcon className="mt-1 size-4 shrink-0 text-muted-foreground" aria-hidden />
              ) : null}
            </CardContent>
          </Card>
        )

        return (
          <li key={`${hit.entityType}:${hit.entityId}`}>
            {href ? (
              <Link href={href} className="block rounded-xl focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">
                {content}
              </Link>
            ) : (
              content
            )}
          </li>
        )
      })}
    </ul>
  )
}
