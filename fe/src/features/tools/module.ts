import { Link2Icon, WrenchIcon } from 'lucide-react'

import type { AppModule } from '@/types/modules'

import { WebhooksPermissions } from './permissions'

const toolsGroup = { key: 'tools', titleKey: 'nav.tools', icon: WrenchIcon }

/** Tools module descriptor — registered in src/config/modules.ts. */
export const toolsModule: AppModule = {
  id: 'tools',
  backendModule: 'webhooks',
  adminNav: [
    {
      titleKey: 'nav.webhooks',
      href: '/admin/tools/webhooks',
      icon: Link2Icon,
      permission: WebhooksPermissions.view,
      group: toolsGroup,
    },
  ],
  breadcrumbs: [
    {
      path: '/admin/tools/webhooks',
      labelKey: 'nav.webhooks',
      parent: { labelKey: 'nav.tools' },
    },
  ],
  permissions: [
    {
      key: WebhooksPermissions.view,
      name: 'View webhooks',
      group: 'webhooks',
      description: 'Open the webhooks inbox and list captured requests',
    },
    {
      key: WebhooksPermissions.modify,
      name: 'Modify webhooks',
      group: 'webhooks',
      description: 'Regenerate the public URL and soft-delete captured requests',
    },
  ],
  publicLinks: [{ href: '/tools', labelKey: 'nav.tools' }],
  sitemap: ['/tools'],
}
