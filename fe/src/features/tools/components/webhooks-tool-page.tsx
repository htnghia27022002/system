'use client'

import Link from 'next/link'
import { ArrowLeftIcon } from 'lucide-react'

import { PublicSiteHeader } from '@/components/common/public-site-header'
import { Button } from '@/components/ui/button'
import { getToolById } from '../catalog'
import { toolsFeatureContent } from '../content'

/**
 * Phase 1 Webhooks tool shell — identity + placeholder body + back to hub.
 * Full webhook product logic is out of scope for this page.
 */
export function WebhooksToolPage() {
  const tool = getToolById('webhooks')
  const copy = toolsFeatureContent.webhooks
  const title = tool?.name ?? copy.title
  const subtitle = tool?.description ?? copy.subtitle

  return (
    <div className="min-h-[100dvh] bg-background text-foreground">
      <PublicSiteHeader variant="tools" />
      <main className="mx-auto w-full max-w-6xl px-4 py-10 sm:px-6 sm:py-14 lg:px-8">
        <Button variant="ghost" size="sm" asChild className="-ml-2 mb-6">
          <Link href="/tools">
            <ArrowLeftIcon className="size-4" />
            {toolsFeatureContent.backToTools}
          </Link>
        </Button>

        <header className="max-w-2xl">
          <p className="text-sm font-medium text-primary">{toolsFeatureContent.navTools}</p>
          <h1 className="mt-2 font-heading text-3xl font-semibold tracking-tight sm:text-4xl">
            {title}
          </h1>
          <p className="mt-3 text-base text-muted-foreground sm:text-lg">{subtitle}</p>
        </header>

        <section
          className="mt-12 rounded-2xl border border-dashed border-border/80 bg-muted/20 px-6 py-12 sm:px-10"
          aria-labelledby="webhooks-placeholder-title"
        >
          <h2
            id="webhooks-placeholder-title"
            className="font-heading text-xl font-semibold tracking-tight"
          >
            {copy.placeholderTitle}
          </h2>
          <p className="mt-3 max-w-prose text-sm leading-relaxed text-muted-foreground sm:text-base">
            {copy.placeholderBody}
          </p>
        </section>
      </main>
    </div>
  )
}
