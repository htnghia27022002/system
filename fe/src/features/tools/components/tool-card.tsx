'use client'

import Link from 'next/link'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { cn } from '@/lib/utils'

import { toolsFeatureContent } from '../content'
import type { ToolEntry } from '../types'

type ToolCardProps = {
  tool: ToolEntry
  className?: string
}

function isNavigable(tool: ToolEntry): tool is ToolEntry & { href: string } {
  return tool.status === 'available' && typeof tool.href === 'string' && tool.href.length > 0
}

export function ToolCard({ tool, className }: ToolCardProps) {
  const navigable = isNavigable(tool)
  const statusLabel =
    tool.status === 'available' ? toolsFeatureContent.available : toolsFeatureContent.comingSoon

  return (
    <Card
      className={cn(
        'h-full transition-colors focus-within:ring-2 focus-within:ring-ring/50',
        className,
      )}
      data-tool-id={tool.id}
      data-tool-status={tool.status}
    >
      <CardHeader>
        <div className="flex flex-wrap items-start justify-between gap-2">
          <CardTitle className="text-lg">{tool.name}</CardTitle>
          <Badge variant={tool.status === 'available' ? 'default' : 'secondary'}>
            {statusLabel}
          </Badge>
        </div>
        {tool.category ? (
          <p className="text-xs text-muted-foreground">{tool.category}</p>
        ) : null}
        <CardDescription>{tool.description}</CardDescription>
      </CardHeader>
      <CardContent className="flex-1" />
      <CardFooter className="justify-end gap-2">
        {navigable ? (
          <Button asChild size="sm">
            <Link href={tool.href}>{toolsFeatureContent.openTool}</Link>
          </Button>
        ) : (
          <Button size="sm" variant="outline" disabled aria-disabled="true">
            {tool.status === 'comingSoon'
              ? toolsFeatureContent.comingSoon
              : toolsFeatureContent.unavailable}
          </Button>
        )}
      </CardFooter>
    </Card>
  )
}
