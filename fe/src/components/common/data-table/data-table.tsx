import {
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
  type ColumnDef,
  type SortingState,
  type ColumnFiltersState,
} from '@tanstack/react-table'
import { SearchIcon } from 'lucide-react'
import { useState, type ReactNode } from 'react'

import { Input } from '@/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { cn } from '@/lib/utils'

import { DataTableEmpty } from './data-table-empty'
import { DataTablePagination } from './data-table-pagination'

type DataTableProps<TData, TValue = unknown> = {
  columns: ColumnDef<TData, TValue>[]
  data: TData[]
  getRowId?: (row: TData) => string

  /** Show the filter input. filterKey must be an accessorKey on ColumnDef. */
  filterKey?: string
  filterPlaceholder?: string

  /** Completely custom empty state */
  emptyMessage?: string

  /** Per-row mobile card renderer — required for mobile layout */
  renderMobileCard: (row: TData) => ReactNode

  /** Toolbar slot: actions placed to the right of the filter input */
  toolbar?: ReactNode

  /** Default page size (default 10) */
  defaultPageSize?: number

  className?: string
}

export function DataTable<TData, TValue = unknown>({
  columns,
  data,
  getRowId,
  filterKey,
  filterPlaceholder = 'Search…',
  emptyMessage = 'No results.',
  renderMobileCard,
  toolbar,
  defaultPageSize = 10,
  className,
}: DataTableProps<TData, TValue>) {
  const [sorting, setSorting] = useState<SortingState>([])
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([])

  const table = useReactTable({
    data,
    columns,
    getRowId,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    onSortingChange: setSorting,
    onColumnFiltersChange: (updater) => {
      setColumnFilters(updater)
      table.setPageIndex(0)
    },
    state: { sorting, columnFilters },
    initialState: {
      pagination: { pageSize: defaultPageSize, pageIndex: 0 },
    },
  })

  const rows = table.getRowModel().rows
  const hasFilter = Boolean(filterKey)
  const hasToolbar = hasFilter || Boolean(toolbar)

  return (
    <div className={cn('flex flex-col gap-0', className)}>
      {hasToolbar && (
        <div className="flex flex-col gap-2 pb-3 sm:flex-row sm:items-center sm:justify-between">
          {hasFilter && (
            <div className="relative w-full sm:max-w-xs">
              <SearchIcon className="absolute left-2.5 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder={filterPlaceholder}
                value={
                  (table
                    .getColumn(filterKey!)
                    ?.getFilterValue() as string) ?? ''
                }
                onChange={(e) => {
                  table.getColumn(filterKey!)?.setFilterValue(e.target.value)
                }}
                className="pl-8"
              />
            </div>
          )}
          {toolbar && (
            <div className="flex shrink-0 items-center gap-2">{toolbar}</div>
          )}
        </div>
      )}

      {/* Mobile card list — visible below md */}
      <div className="flex flex-col gap-3 md:hidden">
        {rows.length === 0 ? (
          <DataTableEmpty message={emptyMessage} />
        ) : (
          rows.map((row) => (
            <div key={row.id}>{renderMobileCard(row.original)}</div>
          ))
        )}
      </div>

      {/* Desktop table — visible from md */}
      <div className="hidden min-w-0 rounded-xl border md:flex md:flex-col">
        <div className="relative w-full overflow-x-auto">
          <Table>
            <TableHeader className="bg-muted/50">
              {table.getHeaderGroups().map((headerGroup) => (
                <TableRow
                  key={headerGroup.id}
                  className="hover:bg-transparent"
                >
                  {headerGroup.headers.map((header) => (
                    <TableHead
                      key={header.id}
                      className="text-xs font-semibold uppercase tracking-wide text-muted-foreground"
                      style={{
                        width: header.getSize() !== 150 ? header.getSize() : undefined,
                      }}
                    >
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext(),
                          )}
                    </TableHead>
                  ))}
                </TableRow>
              ))}
            </TableHeader>
            <TableBody>
              {rows.length === 0 ? (
                <TableRow className="hover:bg-transparent">
                  <TableCell colSpan={columns.length}>
                    <DataTableEmpty message={emptyMessage} />
                  </TableCell>
                </TableRow>
              ) : (
                rows.map((row) => (
                  <TableRow
                    key={row.id}
                    data-state={row.getIsSelected() ? 'selected' : undefined}
                  >
                    {row.getVisibleCells().map((cell) => (
                      <TableCell key={cell.id}>
                        {flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext(),
                        )}
                      </TableCell>
                    ))}
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </div>

        {table.getPageCount() > 1 ||
        table.getState().pagination.pageSize < data.length ? (
          <DataTablePagination table={table} />
        ) : null}
      </div>

      {/* Mobile pagination */}
      {rows.length > 0 && (
        <div className="mt-2 md:hidden">
          {table.getPageCount() > 1 ||
          table.getState().pagination.pageSize < data.length ? (
            <DataTablePagination table={table} />
          ) : null}
        </div>
      )}
    </div>
  )
}
