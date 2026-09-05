'use client'

import Link from 'next/link'
import { useTranslation } from 'react-i18next'

import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar'
import { useCurrentPath } from '@/hooks/use-current-path'
import type { NavItem } from '@/types/navigation'

type NavMainProps = {
  items?: NavItem[]
}

export function NavMain({ items = [] }: NavMainProps) {
  const { t } = useTranslation('admin')
  const { isCurrentPath } = useCurrentPath()

  const platformItems = items.filter((item) => !item.items?.length)
  const sectionGroups = items.filter(
    (item) => item.items && item.items.length > 0,
  )

  const renderLeaf = (item: NavItem) => (
    <SidebarMenuItem key={item.href}>
      <SidebarMenuButton
        asChild
        isActive={isCurrentPath(item.href)}
        tooltip={item.title}
      >
        <Link href={item.href}>
          {item.icon ? <item.icon /> : null}
          <span>{item.title}</span>
        </Link>
      </SidebarMenuButton>
    </SidebarMenuItem>
  )

  return (
    <>
      {platformItems.length > 0 ? (
        <SidebarGroup>
          <SidebarGroupLabel>{t('nav.platform')}</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>{platformItems.map(renderLeaf)}</SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      ) : null}

      {sectionGroups.map((group) => (
        <SidebarGroup key={group.title}>
          <SidebarGroupLabel>{group.title}</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {group.items?.map(renderLeaf)}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      ))}
    </>
  )
}
