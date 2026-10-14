'use client'

import { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { PageContainer, PageHeader, Panel } from '@/components/common/page-layout'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { FieldGroup, FieldLegend, FieldSet } from '@/components/ui/field'
import { Skeleton } from '@/components/ui/skeleton'
import { Spinner } from '@/components/ui/spinner'
import { usePermissions } from '@/features/access-control'

import { useSettings, useUpdateSettings } from '../hooks/use-system'
import type { AppSetting, SettingValues } from '../types'
import { SettingField } from './setting-field'

export function SettingsPage() {
  const { t } = useTranslation('admin')
  const query = useSettings()
  const update = useUpdateSettings()
  const { canModify } = usePermissions()
  const editable = canModify('settings')
  const [draft, setDraft] = useState<SettingValues>({})

  const groups = useMemo(() => {
    const byGroup = new Map<string, AppSetting[]>()
    for (const setting of query.data ?? []) {
      byGroup.set(setting.group, [...(byGroup.get(setting.group) ?? []), setting])
    }
    return [...byGroup.entries()]
  }, [query.data])

  const changed = useMemo(() => {
    const out: SettingValues = {}
    for (const setting of query.data ?? []) {
      if (setting.key in draft && draft[setting.key] !== setting.value) out[setting.key] = draft[setting.key]
    }
    return out
  }, [draft, query.data])
  const dirty = Object.keys(changed).length > 0

  const header = (
    <PageHeader
      title={t('system.settings.title')}
      description={t('system.settings.description')}
      actions={
        editable ? (
          <div className="flex w-full flex-col gap-2 md:w-auto md:flex-row">
            <Button variant="outline" disabled={!dirty || update.isPending} onClick={() => setDraft({})}>
              {t('access.actions.cancel')}
            </Button>
            <Button
              disabled={!dirty || update.isPending}
              onClick={() => update.mutate(changed, { onSuccess: () => setDraft({}) })}
            >
              {update.isPending ? <Spinner className="size-4" /> : null}
              {t('access.actions.save')}
            </Button>
          </div>
        ) : null
      }
    />
  )

  if (query.isLoading) {
    return (
      <PageContainer>
        {header}
        <Skeleton className="h-40 w-full max-w-3xl" />
        <Skeleton className="h-40 w-full max-w-3xl" />
      </PageContainer>
    )
  }

  if (query.isError) {
    return (
      <PageContainer>
        {header}
        <Card className="max-w-3xl">
          <CardContent className="flex flex-col gap-3 py-6">
            <p className="text-sm text-muted-foreground">{t('system.settings.error')}</p>
            <Button variant="outline" className="w-fit" onClick={() => void query.refetch()}>
              {t('access.actions.retry')}
            </Button>
          </CardContent>
        </Card>
      </PageContainer>
    )
  }

  return (
    <PageContainer>
      {header}
      {groups.length === 0 ? (
        <p className="text-sm text-muted-foreground">{t('system.settings.empty')}</p>
      ) : null}
      {groups.map(([group, settings]) => (
        <Panel key={group} className="max-w-3xl">
          <FieldSet>
            <FieldLegend>{t(`system.settings.groups.${group}`, { defaultValue: group })}</FieldLegend>
            <FieldGroup>
              {settings.map((setting) => (
                <SettingField
                  key={setting.key}
                  setting={setting}
                  value={setting.key in draft ? draft[setting.key] : setting.value}
                  disabled={!editable || update.isPending}
                  onChange={(value) => setDraft((prev) => ({ ...prev, [setting.key]: value }))}
                />
              ))}
            </FieldGroup>
          </FieldSet>
        </Panel>
      ))}
    </PageContainer>
  )
}
