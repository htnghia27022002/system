'use client'

import type { ReactNode } from 'react'

import { AppLogo } from '@/components/common/app-logo'
import { ThemeToggle } from '@/components/common/theme-toggle'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'

type AuthLayoutProps = {
  title: string
  description: string
  children: ReactNode
}

export function AuthLayout({ title, description, children }: AuthLayoutProps) {
  return (
    <div className="relative flex min-h-svh flex-col items-center justify-center bg-muted/30 p-6 md:p-10">
      <div className="absolute right-4 top-4">
        <ThemeToggle className="border-transparent bg-transparent text-muted-foreground hover:bg-muted hover:text-foreground" />
      </div>
      <div className="w-full max-w-md">
        <Card className="border-border/60 shadow-sm">
          <CardHeader className="items-center space-y-4 pb-2 text-center">
            <AppLogo />
            <div className="space-y-1">
              <CardTitle className="text-xl font-semibold">{title}</CardTitle>
              <CardDescription>{description}</CardDescription>
            </div>
          </CardHeader>
          <CardContent className="pt-2">{children}</CardContent>
        </Card>
      </div>
    </div>
  )
}
