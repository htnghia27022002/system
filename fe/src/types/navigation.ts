import type { ComponentType } from 'react'

export type BreadcrumbItem = {
  title: string
  href?: string
}

export type NavItem = {
  title: string
  href: string
  icon?: ComponentType<{ className?: string }>
  permission?: string
  items?: NavItem[]
}

export type RouteBreadcrumbHandle = {
  breadcrumbKey?: string
  parentBreadcrumbKey?: string
  parentHref?: string
}
