'use client'

import Link from 'next/link'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { getToolsTeaser, type ToolEntry } from '@/features/tools'

import { landingContent } from '../../content'
import { SectionReveal } from '../section-reveal'
import { SectionBase } from './section-base'

function isNavigable(tool: ToolEntry): tool is ToolEntry & { href: string } {
  return tool.status === 'available' && typeof tool.href === 'string' && tool.href.length > 0
}

export function ToolsSection() {
  const { tools: toolsCopy } = landingContent
  const teaser = getToolsTeaser(3)
  const isEmpty = teaser.length === 0

  return (
    <SectionBase id="tools" className="landing-section-tools">
      <SectionReveal>
        <h2 className="max-w-xl font-heading text-3xl font-semibold tracking-tight sm:text-4xl">
          {toolsCopy.title}
        </h2>
        <p className="mt-3 max-w-xl text-muted-foreground">{toolsCopy.subtitle}</p>
      </SectionReveal>

      {isEmpty ? (
        <SectionReveal delay={0.06}>
          <div
            className="mt-10 rounded-xl border border-dashed border-border bg-muted/30 px-6 py-10"
            role="status"
          >
            <p className="font-medium text-foreground">{toolsCopy.emptyTitle}</p>
            <p className="mt-2 text-sm text-muted-foreground">{toolsCopy.emptyBody}</p>
          </div>
        </SectionReveal>
      ) : (
        <ul className="mt-10 grid gap-4 md:grid-cols-3">
          {teaser.map((tool, index) => {
            const navigable = isNavigable(tool)
            const body = (
              <>
                <div className="flex flex-wrap items-center gap-2">
                  <h3 className="font-heading text-base font-semibold text-foreground">
                    {tool.name}
                  </h3>
                  <Badge variant={tool.status === 'available' ? 'default' : 'secondary'}>
                    {tool.status === 'available' ? 'Available' : 'Coming soon'}
                  </Badge>
                </div>
                <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
                  {tool.description}
                </p>
              </>
            )

            return (
              <li key={tool.id}>
                <SectionReveal delay={0.04 * (index + 1)} className="h-full">
                  <div className="h-full rounded-xl border border-border bg-card p-5 shadow-sm transition-colors hover:border-primary/40">
                    {navigable ? (
                      <Link
                        href={tool.href}
                        className="block h-full rounded-md outline-none focus-visible:ring-2 focus-visible:ring-ring/60"
                      >
                        {body}
                      </Link>
                    ) : (
                      <div className="h-full opacity-90">{body}</div>
                    )}
                  </div>
                </SectionReveal>
              </li>
            )
          })}
        </ul>
      )}

      <SectionReveal delay={0.1}>
        <div className="mt-10">
          <Button size="lg" asChild>
            <Link href={toolsCopy.primaryCtaHref}>{toolsCopy.primaryCta}</Link>
          </Button>
        </div>
      </SectionReveal>
    </SectionBase>
  )
}
