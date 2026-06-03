import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Label } from '@/components/ui/label'
import { cn } from '@/lib/utils'

import type { Permission } from '../types'

/**
 * Preferred action order for column display.
 * Actions not in this list are appended at the end.
 */
const ACTION_ORDER = ['read', 'create', 'update', 'delete']

type CheckState = boolean | 'indeterminate'

function resolveState(selected: number, total: number): CheckState {
  if (total === 0) return false
  if (selected === total) return true
  if (selected > 0) return 'indeterminate'
  return false
}

type PermissionPickerProps = {
  permissions: Permission[]
  selectedKeys: string[]
  onChange: (keys: string[]) => void
  error?: string
}

export function PermissionPicker({
  permissions,
  selectedKeys,
  onChange,
  error,
}: PermissionPickerProps) {
  const { t } = useTranslation('admin')
  const selected = useMemo(() => new Set(selectedKeys), [selectedKeys])

  // ---------- derive matrix structure ----------
  const { groups, actions, lookup } = useMemo(() => {
    const groupsSet: string[] = []
    const actionsSet: string[] = []
    const lookup = new Map<string, Permission>()

    for (const p of permissions) {
      const [grp, act] = p.key.split(':')
      if (!groupsSet.includes(grp)) groupsSet.push(grp)
      if (!actionsSet.includes(act)) actionsSet.push(act)
      lookup.set(p.key, p)
    }

    // sort actions by preferred order
    actionsSet.sort((a, b) => {
      const ai = ACTION_ORDER.indexOf(a)
      const bi = ACTION_ORDER.indexOf(b)
      if (ai === -1 && bi === -1) return a.localeCompare(b)
      if (ai === -1) return 1
      if (bi === -1) return -1
      return ai - bi
    })

    return { groups: groupsSet, actions: actionsSet, lookup }
  }, [permissions])

  // ---------- helpers ----------
  const cell = (group: string, action: string) =>
    lookup.get(`${group}:${action}`) ?? null

  const toggleKeys = (keys: string[], force: boolean) => {
    const next = new Set(selected)
    for (const k of keys) {
      force ? next.add(k) : next.delete(k)
    }
    onChange([...next])
  }

  // row (group) state
  const rowKeys = (group: string) =>
    actions
      .map((a) => cell(group, a))
      .filter(Boolean)
      .map((p) => p!.key)

  const rowState = (group: string): CheckState => {
    const keys = rowKeys(group)
    const sel = keys.filter((k) => selected.has(k)).length
    return resolveState(sel, keys.length)
  }

  const toggleRow = (group: string, state: CheckState) =>
    toggleKeys(rowKeys(group), state !== true)

  // column (action) state
  const colKeys = (action: string) =>
    groups
      .map((g) => cell(g, action))
      .filter(Boolean)
      .map((p) => p!.key)

  const colState = (action: string): CheckState => {
    const keys = colKeys(action)
    const sel = keys.filter((k) => selected.has(k)).length
    return resolveState(sel, keys.length)
  }

  const toggleCol = (action: string, state: CheckState) =>
    toggleKeys(colKeys(action), state !== true)

  // global state
  const allKeys = permissions.map((p) => p.key)
  const globalState: CheckState = resolveState(
    allKeys.filter((k) => selected.has(k)).length,
    allKeys.length,
  )

  const toggleAll = (state: CheckState) =>
    toggleKeys(allKeys, state !== true)

  // per-group selected count
  const groupCount = (group: string) =>
    rowKeys(group).filter((k) => selected.has(k)).length

  return (
    <div className="flex flex-col gap-3">
      {/* Header row */}
      <div className="flex items-center justify-between gap-2">
        <div className="flex flex-col gap-0.5">
          <Label className="text-sm font-medium">
            {t('access.roles.fields.permissions')}
          </Label>
          <span className="text-xs text-muted-foreground">
            {t('access.roles.selectedCount', {
              count: selected.size,
              total: permissions.length,
            })}
          </span>
        </div>

        <div className="flex items-center gap-1">
          <Button
            type="button"
            variant="ghost"
            size="sm"
            className="h-7 text-xs"
            onClick={() => onChange([...allKeys])}
            disabled={selected.size === allKeys.length}
          >
            {t('access.roles.selectAll', { defaultValue: 'Select all' })}
          </Button>
          <span className="text-muted-foreground/40">·</span>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            className="h-7 text-xs"
            onClick={() => onChange([])}
            disabled={selected.size === 0}
          >
            {t('access.roles.clearAll', { defaultValue: 'Clear' })}
          </Button>
        </div>
      </div>

      {/* Matrix table */}
      <div className="overflow-x-auto rounded-lg border">
        <table className="w-full border-collapse">
          {/* ---- thead: action column headers ---- */}
          <thead>
            <tr className="border-b bg-muted/50">
              {/* top-left: ALL checkbox */}
              <th
                scope="col"
                className="sticky left-0 z-10 min-w-[160px] bg-muted/50 px-4 py-3 text-align"
              >
                <div className="flex items-center gap-2.5">
                  <Checkbox
                    id="perm-all"
                    checked={globalState}
                    onCheckedChange={() => toggleAll(globalState)}
                    aria-label="Toggle all permissions"
                  />
                  <label
                    htmlFor="perm-all"
                    className="cursor-pointer text-xs font-semibold uppercase tracking-wider text-muted-foreground"
                  >
                    {t('access.roles.groupLabel', { defaultValue: 'Group' })}
                  </label>
                </div>
              </th>

              {/* one th per action */}
              {actions.map((action) => {
                const state = colState(action)
                return (
                  <th
                    key={action}
                    scope="col"
                    className="w-36 px-4 py-3"
                  >
                    <div className="flex flex-col items-center gap-2">
                      <Checkbox
                        id={`col-${action}`}
                        checked={state}
                        onCheckedChange={() => toggleCol(action, state)}
                        aria-label={`Toggle all ${action} permissions`}
                      />
                      <label
                        htmlFor={`col-${action}`}
                        className="cursor-pointer text-xs font-semibold uppercase tracking-wider text-muted-foreground"
                      >
                        {action}
                      </label>
                    </div>
                  </th>
                )
              })}

              {/* count column */}
              <th
                scope="col"
                className="w-20 px-4 py-3 text-center text-xs font-semibold uppercase tracking-wider text-muted-foreground"
              >
                {t('access.roles.countHeader', { defaultValue: 'Count' })}
              </th>
            </tr>
          </thead>

          {/* ---- tbody: one row per group ---- */}
          <tbody>
            {groups.map((group, i) => {
              const rState = rowState(group)
              const count = groupCount(group)
              const total = rowKeys(group).length

              return (
                <tr
                  key={group}
                  className={cn(
                    'border-b transition-colors last:border-b-0',
                    i % 2 === 0 ? 'bg-background' : 'bg-muted/20',
                  )}
                >
                  {/* group name + row checkbox — sticky */}
                  <th
                    scope="row"
                    className={cn(
                      'sticky left-0 z-10 px-4 py-3',
                      i % 2 === 0 ? 'bg-background' : 'bg-muted/20',
                    )}
                  >
                    <div className="flex items-center gap-2.5">
                      <Checkbox
                        id={`row-${group}`}
                        checked={rState}
                        onCheckedChange={() => toggleRow(group, rState)}
                        aria-label={`Toggle all ${group} permissions`}
                      />
                      <label
                        htmlFor={`row-${group}`}
                        className="cursor-pointer text-sm font-medium capitalize text-foreground/80"
                      >
                        {group}
                      </label>
                    </div>
                  </th>

                  {/* one cell per action */}
                  {actions.map((action) => {
                    const perm = cell(group, action)
                    if (!perm) {
                      return (
                        <td
                          key={action}
                          className="px-4 py-3 text-base text-muted-foreground/25"
                          aria-label="Not available"
                        >
                          <div className="flex justify-center">—</div>
                        </td>
                      )
                    }
                    const checked = selected.has(perm.key)
                    return (
                      <td
                        key={action}
                        className={cn(
                          'px-4 py-3 transition-colors',
                          checked && 'bg-primary/5',
                        )}
                      >
                        <div className="flex justify-center">
                          <Checkbox
                            id={`cell-${perm.key}`}
                            checked={checked}
                            onCheckedChange={(v) => {
                              const next = new Set(selected)
                              v === true ? next.add(perm.key) : next.delete(perm.key)
                              onChange([...next])
                            }}
                            aria-label={perm.name}
                            title={`${perm.name} — ${perm.description}`}
                          />
                        </div>
                      </td>
                    )
                  })}

                  {/* row count */}
                  <td className="px-4 py-3 text-right">
                    <span
                      className={cn(
                        'text-sm tabular-nums',
                        count === total
                          ? 'font-semibold text-emerald-600 dark:text-emerald-400'
                          : count > 0
                            ? 'font-semibold text-amber-600 dark:text-amber-400'
                            : 'text-muted-foreground',
                      )}
                    >
                      {count}/{total}
                    </span>
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>

      {error ? (
        <p className="text-sm text-destructive" role="alert">
          {error}
        </p>
      ) : null}
    </div>
  )
}
