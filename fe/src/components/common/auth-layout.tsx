'use client'

import type { ReactNode } from 'react'

import { AppLogo } from '@/components/common/app-logo'
import { ThemeToggle } from '@/components/common/theme-toggle'

type AuthLayoutProps = {
  title: string
  description: string
  children: ReactNode
}

export function AuthLayout({ title, description, children }: AuthLayoutProps) {
  return (
    <div className="relative flex min-h-svh flex-col items-center justify-center gap-6 bg-background p-6 md:p-10">
      <div className="absolute right-4 top-4">
        <ThemeToggle className="border-transparent bg-transparent text-muted-foreground hover:bg-muted hover:text-foreground" />
      </div>
      <div className="w-full max-w-sm">
        <div className="flex flex-col gap-8">
          <div className="flex flex-col items-center gap-4">
            <AppLogo />
            <div className="space-y-2 text-center">
              <h1 className="text-xl font-medium">{title}</h1>
              <p className="text-center text-sm text-muted-foreground">
                {description}
              </p>
            </div>
          </div>
          {children}
        </div>
      </div>
    </div>
  )
}
