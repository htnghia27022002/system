import { useCallback, useEffect, useRef, useState } from 'react'

const STORAGE_WIDTH_KEY = 'admin_sidebar_width'
const STORAGE_COLLAPSED_KEY = 'admin_sidebar_collapsed'

export const SIDEBAR_COLLAPSED_WIDTH = 56
export const SIDEBAR_MIN_WIDTH = 188
export const SIDEBAR_MAX_WIDTH = 420
export const SIDEBAR_DEFAULT_WIDTH = 240

function getStoredWidth(): number {
  try {
    const raw = localStorage.getItem(STORAGE_WIDTH_KEY)
    if (!raw) return SIDEBAR_DEFAULT_WIDTH
    const n = Number(raw)
    if (!Number.isFinite(n)) return SIDEBAR_DEFAULT_WIDTH
    return Math.min(SIDEBAR_MAX_WIDTH, Math.max(SIDEBAR_MIN_WIDTH, n))
  } catch {
    return SIDEBAR_DEFAULT_WIDTH
  }
}

function getStoredCollapsed(): boolean {
  try {
    return localStorage.getItem(STORAGE_COLLAPSED_KEY) === 'true'
  } catch {
    return false
  }
}

export type AdminSidebarState = {
  width: number
  isCollapsed: boolean
  isMobileOpen: boolean
  isResizing: boolean
  toggleCollapsed: () => void
  openMobile: () => void
  closeMobile: () => void
  handleResizeStart: (e: React.MouseEvent) => void
}

export function useAdminSidebar(): AdminSidebarState {
  const [width, setWidth] = useState(() => getStoredWidth())
  const [isCollapsed, setIsCollapsed] = useState(() => getStoredCollapsed())
  const [isMobileOpen, setIsMobileOpen] = useState(false)
  const [isResizing, setIsResizing] = useState(false)

  const widthRef = useRef(width)
  widthRef.current = width

  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_WIDTH_KEY, String(width))
    } catch {}
  }, [width])

  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_COLLAPSED_KEY, String(isCollapsed))
    } catch {}
  }, [isCollapsed])

  const toggleCollapsed = useCallback(() => {
    setIsCollapsed((prev) => !prev)
  }, [])

  const openMobile = useCallback(() => setIsMobileOpen(true), [])
  const closeMobile = useCallback(() => setIsMobileOpen(false), [])

  const handleResizeStart = useCallback((e: React.MouseEvent) => {
    e.preventDefault()
    const startX = e.clientX
    const startWidth = widthRef.current

    setIsResizing(true)
    document.body.style.cursor = 'col-resize'
    document.body.style.userSelect = 'none'

    const onMouseMove = (moveEvent: MouseEvent) => {
      const delta = moveEvent.clientX - startX
      const newWidth = Math.min(
        SIDEBAR_MAX_WIDTH,
        Math.max(SIDEBAR_MIN_WIDTH, startWidth + delta),
      )
      setWidth(newWidth)
    }

    const onMouseUp = () => {
      setIsResizing(false)
      document.body.style.cursor = ''
      document.body.style.userSelect = ''
      document.removeEventListener('mousemove', onMouseMove)
      document.removeEventListener('mouseup', onMouseUp)
    }

    document.addEventListener('mousemove', onMouseMove)
    document.addEventListener('mouseup', onMouseUp)
  }, [])

  return {
    width: isCollapsed ? SIDEBAR_COLLAPSED_WIDTH : width,
    isCollapsed,
    isMobileOpen,
    isResizing,
    toggleCollapsed,
    openMobile,
    closeMobile,
    handleResizeStart,
  }
}
