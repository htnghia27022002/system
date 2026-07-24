'use client'

import Link from 'next/link'

import { PublicSiteHeader } from '@/components/common/public-site-header'
import { Button } from '@/components/ui/button'

import { getAllTools } from './catalog'
import { ToolCard } from './components/tool-card'
import { toolsFeatureContent } from './content'
import type { ToolEntry } from './types'

function groupByCategory(tools: ToolEntry[]): Array<{ category: string; tools: ToolEntry[] }> {
  const map = new Map<string, ToolEntry[]>()
  for (const tool of tools) {
    const key = tool.category?.trim() || 'General'
    const list = map.get(key) ?? []
    list.push(tool)
    map.set(key, list)
  }
  return Array.from(map.entries()).map(([category, groupTools]) => ({
    category,
    tools: groupTools,
  }))
}

export function ToolsHubPage() {
  const tools = getAllTools()
  const groups = groupByCategory(tools)
  const isEmpty = tools.length === 0

  return (
    <div className="min-h-[100dvh] bg-background text-foreground">
      <PublicSiteHeader variant="tools" />
      <main className="mx-auto w-full max-w-6xl px-4 py-12 sm:px-6 sm:py-16 lg:px-8">
        <div className="max-w-2xl">
          <h1 className="font-heading text-3xl font-semibold tracking-tight sm:text-4xl">
            {toolsFeatureContent.hubTitle}
          </h1>
          <p className="mt-3 text-base text-muted-foreground sm:text-lg">
            {toolsFeatureContent.hubSubtitle}
          </p>
        </div>

        {isEmpty ? (
          <div
            className="mt-12 rounded-xl border border-dashed border-border/80 bg-muted/20 px-6 py-16 text-center"
            role="status"
          >
            <p className="font-medium text-foreground">{toolsFeatureContent.emptyTitle}</p>
            <p className="mt-2 text-sm text-muted-foreground">{toolsFeatureContent.emptyBody}</p>
            <Button asChild className="mt-6" variant="outline">
              <Link href="/">{toolsFeatureContent.emptyCta}</Link>
            </Button>
          </div>
        ) : (
          <div className="mt-12 space-y-12">
            {groups.map((group) => (
              <section key={group.category} aria-labelledby={`tools-group-${group.category}`}>
                <h2
                  id={`tools-group-${group.category}`}
                  className="mb-4 text-sm font-medium text-muted-foreground"
                >
                  {group.category}
                </h2>
                <ul className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                  {group.tools.map((tool) => (
                    <li key={tool.id}>
                      <ToolCard tool={tool} />
                    </li>
                  ))}
                </ul>
              </section>
            ))}
          </div>
        )}
      </main>
    </div>
  )
}
