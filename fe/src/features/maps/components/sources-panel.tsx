'use client'

import { useState } from 'react'
import { useTranslation } from 'react-i18next'

import { PermissionGate, PermissionKeys } from '@/features/access-control'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { StatusBadge } from '@/components/ui/status-badge'
import { CircleOffIcon, PowerIcon, Trash2Icon } from 'lucide-react'

import {
  useMapSourceMutations,
  useMapSources,
} from '../hooks/use-map-sources'
import type { SourceRecord } from '../types'
import { SourceFormDialog } from './source-form-dialog'

export function SourcesPanel() {
  const { t } = useTranslation('admin')
  const sourcesQuery = useMapSources(1)
  const { createSource, updateSource, deleteSource } = useMapSourceMutations()
  const [dialogOpen, setDialogOpen] = useState(false)
  const [mode, setMode] = useState<'create' | 'edit'>('create')
  const [editing, setEditing] = useState<SourceRecord | undefined>()
  const [deleteTarget, setDeleteTarget] = useState<SourceRecord | null>(null)

  const sources = sourcesQuery.data?.items ?? []

  const openCreate = () => {
    setMode('create')
    setEditing(undefined)
    setDialogOpen(true)
  }

  const openEdit = (source: SourceRecord) => {
    setMode('edit')
    setEditing(source)
    setDialogOpen(true)
  }

  return (
    <div className="space-y-4 px-4 pb-6">
      <PermissionGate
        permission={PermissionKeys.maps.modify}
        fallback={
          <p className="text-sm text-muted-foreground">
            {t('maps.sources.readOnlyHint')}
          </p>
        }
      >
        <Button className="w-full" onClick={openCreate}>
          {t('maps.sources.createAction')}
        </Button>
      </PermissionGate>

      {sourcesQuery.isLoading ? (
        <Card>
          <CardContent className="p-4 text-sm text-muted-foreground">
            {t('maps.sources.description')}
          </CardContent>
        </Card>
      ) : null}

      {sourcesQuery.isError ? (
        <Card>
          <CardContent className="flex flex-col gap-3 p-4">
            <p className="text-sm text-muted-foreground">
              {t('maps.sources.error')}
            </p>
            <Button
              variant="outline"
              onClick={() => void sourcesQuery.refetch()}
            >
              {t('maps.retry')}
            </Button>
          </CardContent>
        </Card>
      ) : null}

      {!sourcesQuery.isLoading && !sourcesQuery.isError && sources.length === 0 ? (
        <Card>
          <CardContent className="space-y-1 p-4">
            <p className="text-sm font-medium">{t('maps.sources.emptyTitle')}</p>
            <p className="text-sm text-muted-foreground">
              {t('maps.sources.emptyDescription')}
            </p>
          </CardContent>
        </Card>
      ) : null}

      <ul className="space-y-3">
        {sources.map((source) => (
          <li key={source.id}>
            <Card>
              <CardContent className="flex flex-col gap-3 p-4 sm:flex-row sm:items-start sm:justify-between">
                <div className="min-w-0 space-y-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <PermissionGate
                      permission={PermissionKeys.maps.modify}
                      fallback={
                        <span className="text-sm font-medium">{source.name}</span>
                      }
                    >
                      <button
                        type="button"
                        className="text-left text-sm font-medium text-primary underline-offset-4 hover:underline focus-visible:underline"
                        onClick={() => openEdit(source)}
                      >
                        {source.name}
                      </button>
                    </PermissionGate>
                    <StatusBadge
                      variant={source.enabled ? 'success' : 'neutral'}
                      label={
                        source.enabled
                          ? t('maps.sources.enabled')
                          : t('maps.sources.disabled')
                      }
                    />
                  </div>
                  <p className="text-sm text-muted-foreground">
                    {source.httpMethod} · {source.url}
                  </p>
                </div>
                <PermissionGate permission={PermissionKeys.maps.modify}>
                  <div className="flex shrink-0 items-center gap-1 self-start">
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      className="size-11 shrink-0 sm:size-8"
                      disabled={updateSource.isPending}
                      aria-label={
                        source.enabled
                          ? t('maps.sources.disable')
                          : t('maps.sources.enable')
                      }
                      title={
                        source.enabled
                          ? t('maps.sources.disable')
                          : t('maps.sources.enable')
                      }
                      onClick={() =>
                        updateSource.mutate({
                          id: source.id,
                          input: { enabled: !source.enabled },
                        })
                      }
                    >
                      {source.enabled ? (
                        <CircleOffIcon className="size-4" />
                      ) : (
                        <PowerIcon className="size-4" />
                      )}
                    </Button>
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      className="size-11 shrink-0 text-destructive hover:bg-destructive/10 hover:text-destructive sm:size-8"
                      disabled={deleteSource.isPending}
                      aria-label={t('access.actions.delete')}
                      title={t('access.actions.delete')}
                      onClick={() => setDeleteTarget(source)}
                    >
                      <Trash2Icon className="size-4" />
                    </Button>
                  </div>
                </PermissionGate>
              </CardContent>
            </Card>
          </li>
        ))}
      </ul>

      <SourceFormDialog
        open={dialogOpen}
        mode={mode}
        source={editing}
        isPending={createSource.isPending || updateSource.isPending}
        onOpenChange={setDialogOpen}
        onCreate={(input) =>
          createSource.mutate(input, { onSuccess: () => setDialogOpen(false) })
        }
        onUpdate={(input) => {
          if (!editing) return
          updateSource.mutate(
            { id: editing.id, input },
            { onSuccess: () => setDialogOpen(false) },
          )
        }}
      />

      <AlertDialog
        open={Boolean(deleteTarget)}
        onOpenChange={(open) => {
          if (!open) setDeleteTarget(null)
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t('maps.sources.deleteTitle')}</AlertDialogTitle>
            <AlertDialogDescription>
              {t('maps.sources.deleteDescription', {
                name: deleteTarget?.name ?? '',
              })}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{t('access.actions.cancel')}</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => {
                if (!deleteTarget) return
                deleteSource.mutate(deleteTarget.id, {
                  onSuccess: () => setDeleteTarget(null),
                })
              }}
            >
              {t('access.actions.delete')}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}
