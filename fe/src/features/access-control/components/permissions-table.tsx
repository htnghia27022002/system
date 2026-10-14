'use client'

import { useTranslation } from 'react-i18next'

import { PageContainer, PageHeader } from '@/components/common/page-layout'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

import { usePermissionsCatalog } from '../hooks/use-permissions-catalog'

export function PermissionsTable() {
  const { t } = useTranslation('admin')
  const query = usePermissionsCatalog()

  if (query.isLoading) {
    return (
      <PageContainer>
        <Skeleton className="h-10 w-48" />
        <Skeleton className="h-64 w-full" />
      </PageContainer>
    )
  }

  if (query.isError) {
    return (
      <PageContainer>
        <Card>
          <CardContent className="flex flex-col gap-3 py-6">
            <p className="text-sm text-muted-foreground">
              {t('access.permissions.error')}
            </p>
            <Button variant="outline" onClick={() => void query.refetch()}>
              {t('access.actions.retry')}
            </Button>
          </CardContent>
        </Card>
      </PageContainer>
    )
  }

  const permissions = query.data ?? []

  return (
    <PageContainer>
      <PageHeader
        title={t('access.permissions.title')}
        description={t('access.permissions.description')}
      />

      {permissions.length === 0 ? (
        <Card>
          <CardContent className="py-8 text-center text-sm text-muted-foreground">
            {t('access.permissions.empty')}
          </CardContent>
        </Card>
      ) : (
        <div className="rounded-xl border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t('access.permissions.fields.key')}</TableHead>
                <TableHead>{t('access.permissions.fields.name')}</TableHead>
                <TableHead>{t('access.permissions.fields.group')}</TableHead>
                <TableHead>
                  {t('access.permissions.fields.description')}
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {permissions.map((perm) => (
                <TableRow key={perm.key}>
                  <TableCell className="font-mono text-xs">{perm.key}</TableCell>
                  <TableCell className="font-medium">{perm.name}</TableCell>
                  <TableCell>
                    <Badge variant="secondary">{perm.group}</Badge>
                  </TableCell>
                  <TableCell className="text-muted-foreground">
                    {perm.description}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}
    </PageContainer>
  )
}
