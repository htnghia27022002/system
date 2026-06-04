'use client'

import Link from 'next/link'
import { ChevronRightIcon } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible'
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from '@/components/ui/sidebar'
import { useCurrentPath } from '@/hooks/use-current-path'
import type { NavItem } from '@/types/navigation'

type NavMainProps = {
  items?: NavItem[]
}

export function NavMain({ items = [] }: NavMainProps) {
  const { t } = useTranslation('admin')
  const { isCurrentPath } = useCurrentPath()

  return (
    <SidebarGroup className="px-2 py-0">
      <SidebarGroupLabel>{t('nav.platform')}</SidebarGroupLabel>
      <SidebarMenu>
        {items.map((item) => {
          const hasChildren = item.items && item.items.length > 0

          if (hasChildren) {
            const hasActiveChild = item.items?.some((subItem) =>
              isCurrentPath(subItem.href),
            )

            return (
              <Collapsible
                key={item.title}
                defaultOpen={hasActiveChild}
                className="group/collapsible"
              >
                <SidebarMenuItem>
                  <CollapsibleTrigger asChild>
                    <SidebarMenuButton
                      tooltip={{ children: item.title }}
                      className="data-[state=open]:bg-muted"
                    >
                      {item.icon ? <item.icon /> : null}
                      <span>{item.title}</span>
                      <ChevronRightIcon className="ml-auto transition-transform group-data-[state=open]/collapsible:rotate-90" />
                    </SidebarMenuButton>
                  </CollapsibleTrigger>
                  <CollapsibleContent>
                    <SidebarMenuSub>
                      {item.items?.map((subItem) => (
                        <SidebarMenuSubItem key={subItem.title}>
                          <SidebarMenuSubButton
                            asChild
                            isActive={isCurrentPath(subItem.href)}
                          >
                            <Link href={subItem.href}>
                              <span>{subItem.title}</span>
                            </Link>
                          </SidebarMenuSubButton>
                        </SidebarMenuSubItem>
                      ))}
                    </SidebarMenuSub>
                  </CollapsibleContent>
                </SidebarMenuItem>
              </Collapsible>
            )
          }

          return (
            <SidebarMenuItem key={item.title}>
              <SidebarMenuButton
                asChild
                isActive={isCurrentPath(item.href)}
                tooltip={{ children: item.title }}
              >
                <Link href={item.href}>
                  {item.icon ? <item.icon /> : null}
                  <span>{item.title}</span>
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          )
        })}
      </SidebarMenu>
    </SidebarGroup>
  )
}
