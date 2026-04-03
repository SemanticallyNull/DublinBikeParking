import { useEffect, useRef, useState } from 'react'
import styles from './BottomSheet.module.css'

interface Props {
  open: boolean
  onClose: () => void
  children: React.ReactNode
}

export function BottomSheet({ open, onClose, children }: Props) {
  const [expanded, setExpanded] = useState(false)
  const startY = useRef<number | null>(null)

  // Reset expanded state when sheet closes
  useEffect(() => { if (!open) setExpanded(false) }, [open])

  function onTouchStart(e: React.TouchEvent) {
    startY.current = e.touches[0].clientY
  }

  function onTouchEnd(e: React.TouchEvent) {
    if (startY.current === null) return
    const dy = e.changedTouches[0].clientY - startY.current
    if (dy < -40) setExpanded(true)   // swipe up
    if (dy > 40)  expanded ? setExpanded(false) : onClose()
    startY.current = null
  }

  if (!open) return null

  return (
    <>
      <div className={styles.backdrop} onClick={onClose} />
      <div
        className={`${styles.sheet} ${expanded ? styles.expanded : ''}`}
        onTouchStart={onTouchStart}
        onTouchEnd={onTouchEnd}
      >
        <div className={styles.handle} onClick={() => setExpanded(e => !e)} />
        <div className={styles.content}>
          {children}
        </div>
      </div>
    </>
  )
}
