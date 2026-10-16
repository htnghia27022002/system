'use client'

import type { FormEvent, ReactNode } from 'react'

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Drawer,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
} from '@/components/ui/drawer'
import { useIsMobile } from '@/hooks/use-mobile'
import { cn } from '@/lib/utils'

type FormOverlayProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  title: string
  description?: string
  children: ReactNode
  footer: ReactNode
  onSubmit: (event: FormEvent<HTMLFormElement>) => void
  className?: string
}

/**
 * Official shadcn composition:
 * DialogContent → Header + form + Footer (desktop)
 * DrawerContent → Header + form + Footer (mobile)
 *
 * No inner scroll panes. Dialog sizes to content; the page/dialog
 * only scrolls when the whole panel exceeds the viewport.
 */
export function FormOverlay({
  open,
  onOpenChange,
  title,
  description,
  children,
  footer,
  onSubmit,
  className,
}: FormOverlayProps) {
  const isMobile = useIsMobile()

  if (isMobile) {
    return (
      <Drawer open={open} onOpenChange={onOpenChange}>
        <DrawerContent>
          <DrawerHeader className="text-left">
            <DrawerTitle>{title}</DrawerTitle>
            {description ? (
              <DrawerDescription>{description}</DrawerDescription>
            ) : null}
          </DrawerHeader>
          <form onSubmit={onSubmit} noValidate>
            <div className="grid gap-6 px-4">{children}</div>
            <DrawerFooter>{footer}</DrawerFooter>
          </form>
        </DrawerContent>
      </Drawer>
    )
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent
        className={cn(
          'max-h-[calc(100dvh-2rem)] grid-rows-[auto_minmax(0,1fr)] sm:max-w-xl',
          className,
        )}
      >
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          {description ? (
            <DialogDescription>{description}</DialogDescription>
          ) : null}
        </DialogHeader>
        <form
          className="grid min-h-0 grid-rows-[minmax(0,1fr)_auto] gap-6"
          onSubmit={onSubmit}
          noValidate
        >
          {/* Body scrolls when taller than the viewport; header and footer stay put. */}
          <div className="-mx-4 grid min-h-0 content-start gap-6 overflow-y-auto px-4">
            {children}
          </div>
          <DialogFooter>{footer}</DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
