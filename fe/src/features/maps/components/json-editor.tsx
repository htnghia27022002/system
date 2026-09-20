'use client'

import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

import { formatJsonBody, isValidJsonBody } from '../schemas/source-form-schema'

type JsonEditorProps = {
  id?: string
  value: string
  onChange: (value: string) => void
  invalid?: boolean
  className?: string
}

const TOKEN =
  /("(?:\\.|[^"\\])*")\s*:|("(?:\\.|[^"\\])*")|\b(?:true|false|null)\b|-?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?|[{}[\],:]/g

function highlightJson(source: string) {
  const parts: Array<{ text: string; kind: string }> = []
  let cursor = 0

  for (const match of source.matchAll(TOKEN)) {
    const index = match.index ?? 0
    if (index > cursor) {
      parts.push({ text: source.slice(cursor, index), kind: 'plain' })
    }
    const [token, key] = match
    if (key) parts.push({ text: token, kind: 'key' })
    else if (token.startsWith('"')) parts.push({ text: token, kind: 'string' })
    else if (token === 'true' || token === 'false' || token === 'null') {
      parts.push({ text: token, kind: 'atom' })
    } else if (/^-?\d/.test(token)) parts.push({ text: token, kind: 'number' })
    else parts.push({ text: token, kind: 'punct' })
    cursor = index + token.length
  }

  if (cursor < source.length) {
    parts.push({ text: source.slice(cursor), kind: 'plain' })
  }

  return parts.length > 0 ? parts : [{ text: source, kind: 'plain' }]
}

const TOKEN_CLASS: Record<string, string> = {
  key: 'text-sky-700 dark:text-sky-300',
  string: 'text-emerald-700 dark:text-emerald-300',
  atom: 'text-violet-700 dark:text-violet-300',
  number: 'text-amber-700 dark:text-amber-300',
  punct: 'text-muted-foreground',
  plain: 'text-foreground',
}

export function JsonEditor({
  id,
  value,
  onChange,
  invalid,
  className,
}: JsonEditorProps) {
  const { t } = useTranslation('admin')
  const display = value.length > 0 ? value : ' '
  const lines = display.split('\n')
  const rowCount = Math.max(8, lines.length)
  const tokens = useMemo(() => highlightJson(display), [display])
  const valid = isValidJsonBody(value)

  return (
    <div className={cn('space-y-2', className)}>
      <div className="flex items-center justify-end">
        <Button
          type="button"
          variant="outline"
          size="sm"
          disabled={!valid || !value.trim()}
          onClick={() => onChange(formatJsonBody(value))}
        >
          {t('maps.sources.fields.formatJson')}
        </Button>
      </div>
      <div
        className={cn(
          'max-h-64 overflow-auto rounded-lg border bg-muted/30 font-mono text-xs leading-5',
          invalid || !valid
            ? 'border-destructive'
            : 'border-input focus-within:border-ring focus-within:ring-ring/50 focus-within:ring-[3px]',
        )}
      >
        <div className="grid grid-cols-[2.5rem_minmax(0,1fr)]">
          <div
            aria-hidden="true"
            className="select-none border-r border-border bg-muted/50 py-2.5 text-right text-muted-foreground"
          >
            {Array.from({ length: rowCount }, (_, index) => (
              <div key={index} className="px-2">
                {index + 1}
              </div>
            ))}
          </div>
          <div className="relative">
            <pre
              aria-hidden="true"
              className="pointer-events-none absolute inset-0 whitespace-pre-wrap break-words px-3 py-2.5"
            >
              {tokens.map((part, index) => (
                <span key={index} className={TOKEN_CLASS[part.kind]}>
                  {part.text}
                </span>
              ))}
            </pre>
            <textarea
              id={id}
              value={value}
              rows={rowCount}
              spellCheck={false}
              aria-invalid={invalid || !valid}
              onChange={(event) => onChange(event.target.value)}
              className="relative z-10 w-full resize-none bg-transparent px-3 py-2.5 text-transparent caret-foreground outline-none"
            />
          </div>
        </div>
      </div>
      {!valid ? (
        <p className="text-sm text-destructive">
          {t('maps.sources.fields.bodyInvalid')}
        </p>
      ) : null}
    </div>
  )
}
