import { redirect } from 'next/navigation'

/**
 * Owner Webhooks UI lives under admin RBAC.
 * Public capture remains `/tools/webhooks/{uuid}` (nginx / Next rewrite → BE).
 */
export default function WebhooksOwnerRedirectPage() {
  redirect('/admin/tools/webhooks')
}
