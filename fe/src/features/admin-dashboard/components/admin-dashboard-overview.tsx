import { PageContainer } from '@/components/common/page-layout'
import { PlaceholderPattern } from '@/components/ui/placeholder-pattern'

export function AdminDashboardOverview() {
  return (
    <PageContainer>
      <div className="grid auto-rows-min gap-4 md:grid-cols-3">
        <div className="relative aspect-video overflow-hidden rounded-xl border border-border">
          <PlaceholderPattern className="absolute inset-0 size-full stroke-foreground/15" />
        </div>
        <div className="relative aspect-video overflow-hidden rounded-xl border border-border">
          <PlaceholderPattern className="absolute inset-0 size-full stroke-foreground/15" />
        </div>
        <div className="relative aspect-video overflow-hidden rounded-xl border border-border">
          <PlaceholderPattern className="absolute inset-0 size-full stroke-foreground/15" />
        </div>
      </div>
      <div className="relative min-h-screen flex-1 overflow-hidden rounded-xl border border-sidebar-border/70 md:min-h-min dark:border-sidebar-border">
        <PlaceholderPattern className="absolute inset-0 size-full stroke-foreground/15" />
      </div>
    </PageContainer>
  )
}
