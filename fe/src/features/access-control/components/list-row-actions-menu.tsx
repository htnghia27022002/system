import { MoreHorizontalIcon } from 'lucide-react'
import type { ReactNode } from 'react'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

type ListRowActionsMenuProps = {
  children: ReactNode
}

export function ListRowActionsMenu({ children }: ListRowActionsMenuProps) {
  const { t } = useTranslation('admin')

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          className="size-11 shrink-0 sm:size-8"
        >
          <MoreHorizontalIcon className="size-4" />
          <span className="sr-only">{t('access.actions.openMenu')}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">{children}</DropdownMenuContent>
    </DropdownMenu>
  )
}
