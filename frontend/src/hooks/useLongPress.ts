import { useRef, useCallback } from 'react'

interface Options {
  delay?: number       // ms before long-press fires (default 600)
  moveThreshold?: number // px movement before cancelling (default 10)
}

export function useLongPress(callback: (e: TouchEvent) => void, { delay = 600, moveThreshold = 10 }: Options = {}) {
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const startPos = useRef<{ x: number; y: number } | null>(null)

  const cancel = useCallback(() => {
    if (timerRef.current !== null) {
      clearTimeout(timerRef.current)
      timerRef.current = null
    }
    startPos.current = null
  }, [])

  const onTouchStart = useCallback((e: TouchEvent) => {
    const touch = e.touches[0]
    startPos.current = { x: touch.clientX, y: touch.clientY }
    timerRef.current = setTimeout(() => {
      timerRef.current = null
      callback(e)
    }, delay)
  }, [callback, delay])

  const onTouchMove = useCallback((e: TouchEvent) => {
    if (!startPos.current) return
    const touch = e.touches[0]
    const dx = Math.abs(touch.clientX - startPos.current.x)
    const dy = Math.abs(touch.clientY - startPos.current.y)
    if (dx > moveThreshold || dy > moveThreshold) cancel()
  }, [cancel, moveThreshold])

  return { onTouchStart, onTouchMove, onTouchEnd: cancel }
}
