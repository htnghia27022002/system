import type { ReactNode } from 'react'

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { cn } from '@/lib/utils'

/** DESIGN.md store-utility-card + caption-strong field labels */
export const authFieldInputClassName =
  'h-11 rounded-[var(--ds-radius-md)] text-[17px] md:text-[17px]'

type AuthFormShellProps = {
  title: string
  description: string
  children: ReactNode
  footer: ReactNode
}

export function AuthFormShell({
  title,
  description,
  children,
  footer,
}: AuthFormShellProps) {
  return (
    <Card
      className="w-full max-w-md gap-0 border border-border bg-card py-0 shadow-none ring-0"
      aria-labelledby="auth-form-title"
    >
      <CardHeader className="gap-2 px-6 pt-6">
        <CardTitle
          id="auth-form-title"
          className="font-[family-name:var(--ds-font-display)] text-[32px] font-semibold leading-[1.1] tracking-normal sm:text-[40px]"
        >
          {title}
        </CardTitle>
        <CardDescription className="text-[17px] leading-[1.47] tracking-[-0.374px] text-muted-foreground">
          {description}
        </CardDescription>
      </CardHeader>
      <CardContent className="px-6 pb-6">{children}</CardContent>
      <CardFooter className="flex-col gap-0 border-t border-[var(--ds-divider-soft)] bg-transparent px-6 py-6">
        {footer}
      </CardFooter>
    </Card>
  )
}

type AuthFormFieldProps = {
  id: string
  label: string
  error?: string
  children: ReactNode
}

export function AuthFormField({ id, label, error, children }: AuthFormFieldProps) {
  return (
    <div className="space-y-2">
      <Label
        htmlFor={id}
        className="text-[14px] font-semibold leading-[1.29] tracking-[-0.224px] text-foreground"
      >
        {label}
      </Label>
      {children}
      {error ? (
        <p id={`${id}-error`} className="text-sm text-destructive" role="alert">
          {error}
        </p>
      ) : null}
    </div>
  )
}

export function authFieldDescribedBy(id: string, hasError: boolean) {
  return hasError ? `${id}-error` : undefined
}

export function AuthFormFooterText({
  className,
  children,
}: {
  className?: string
  children: ReactNode
}) {
  return (
    <p
      className={cn(
        'w-full text-center text-[17px] leading-[1.47] tracking-[-0.374px] text-muted-foreground',
        className,
      )}
    >
      {children}
    </p>
  )
}
