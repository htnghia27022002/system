'use client'

import { SearchIcon, XIcon } from 'lucide-react'
import {
  useEffect,
  useRef,
  useState,
  type FormEvent,
  type KeyboardEvent,
} from 'react'

import { useDebouncedValue } from '@/hooks/use-debounced-value'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

type ServerListSearchProps = {
  value: string
  placeholder: string
  debounceMs?: number
  onSearch: (value: string) => void
  onClear: () => void
}

export function ServerListSearch({
  value,
  placeholder,
  debounceMs = 300,
  onSearch,
  onClear,
}: ServerListSearchProps) {
  const [draft, setDraft] = useState(value)
  const debouncedDraft = useDebouncedValue(draft, debounceMs)
  const onSearchRef = useRef(onSearch)
  const onClearRef = useRef(onClear)

  useEffect(() => {
    onSearchRef.current = onSearch
  }, [onSearch])

  useEffect(() => {
    onClearRef.current = onClear
  }, [onClear])

  useEffect(() => {
    setDraft(value)
  }, [value])

  useEffect(() => {
    const debouncedTrimmed = debouncedDraft.trim()
    const draftTrimmed = draft.trim()
    const valueTrimmed = value.trim()

    // Ignore stale debounced values while the user is still typing or clearing.
    if (debouncedTrimmed !== draftTrimmed) {
      return
    }

    if (debouncedTrimmed === valueTrimmed) {
      return
    }

    onSearchRef.current(debouncedTrimmed)
  }, [debouncedDraft, draft, value])

  function submitSearch() {
    onSearchRef.current(draft.trim())
  }

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    submitSearch()
  }

  function handleKeyUp(event: KeyboardEvent<HTMLInputElement>) {
    if (event.key === 'Enter') {
      submitSearch()
    }
  }

  function handleClear() {
    setDraft('')
    onClearRef.current()
  }

  const showClear = Boolean(draft.trim() || value.trim())

  return (
    <form
      onSubmit={handleSubmit}
      className="relative w-full sm:max-w-xs"
      role="search"
    >
      <SearchIcon className="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-muted-foreground" />
      <Input
        value={draft}
        onChange={(event) => setDraft(event.target.value)}
        onKeyUp={handleKeyUp}
        placeholder={placeholder}
        className="pr-9 pl-8"
        aria-label={placeholder}
      />
      {showClear ? (
        <Button
          type="button"
          variant="ghost"
          size="icon"
          className="absolute top-1/2 right-0.5 size-7 -translate-y-1/2"
          onClick={handleClear}
          aria-label="Clear search"
        >
          <XIcon className="size-3.5" />
        </Button>
      ) : null}
    </form>
  )
}
