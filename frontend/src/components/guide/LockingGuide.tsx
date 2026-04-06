// frontend/src/components/guide/LockingGuide.tsx
import { useEffect, useRef, useState } from 'react'
import { GUIDE_STEPS } from './steps'
import styles from './LockingGuide.module.css'

interface Props {
  open: boolean
  onClose: () => void
}

export function LockingGuide({ open, onClose }: Props) {
  const [step, setStep] = useState(0)
  const [direction, setDirection] = useState<'forward' | 'back'>('forward')
  const cardRef = useRef<HTMLDivElement>(null)
  const total = GUIDE_STEPS.length

  // Focus trap + keyboard handler
  useEffect(() => {
    if (!open) return
    const prev = document.activeElement as HTMLElement | null
    cardRef.current?.focus()

    function onKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') { onClose(); return }
      if (e.key !== 'Tab') return
      const focusable = cardRef.current?.querySelectorAll<HTMLElement>(
        'button, [href], input, [tabindex]:not([tabindex="-1"])'
      )
      if (!focusable?.length) return
      const first = focusable[0]
      const last = focusable[focusable.length - 1]
      if (e.shiftKey && document.activeElement === first) {
        e.preventDefault()
        last.focus()
      } else if (!e.shiftKey && document.activeElement === last) {
        e.preventDefault()
        first.focus()
      }
    }

    document.addEventListener('keydown', onKeyDown)
    return () => {
      document.removeEventListener('keydown', onKeyDown)
      prev?.focus()
    }
  }, [open, onClose])

  // Reset to step 0 whenever the guide opens
  useEffect(() => {
    if (open) { setStep(0); setDirection('forward') }
  }, [open])

  if (!open) return null

  const current = GUIDE_STEPS[step]
  const Illustration = current.illustration
  const isLast = step === total - 1
  const isFirst = step === 0
  const progressPct = ((step + 1) / total) * 100

  function goNext() {
    if (isLast) { onClose(); return }
    setDirection('forward')
    setStep(s => s + 1)
  }

  function goPrev() {
    if (isFirst) return
    setDirection('back')
    setStep(s => s - 1)
  }

  return (
    <div className={styles.backdrop} role="dialog" aria-modal="true" aria-label="How to lock your bike">
      <div className={styles.card} ref={cardRef} tabIndex={-1}>
        {/* Progress bar */}
        <div className={styles.progressBar} role="progressbar" aria-valuenow={step + 1} aria-valuemin={1} aria-valuemax={total}>
          <div className={styles.progressFill} style={{ width: `${progressPct}%` }} />
        </div>

        {/* Top row */}
        <div className={styles.topRow}>
          <span className={styles.stepCounter}>Step {step + 1} of {total}</span>
          {!isLast && (
            <button className={styles.skipBtn} onClick={onClose}>Skip</button>
          )}
        </div>

        {/* Animated step content */}
        <div
          key={step}
          className={`${styles.stepContent} ${direction === 'forward' ? styles.stepContentForward : styles.stepContentBack}`}
        >
          <div className={styles.illustration}>
            <Illustration />
          </div>
          <h2 className={styles.stepTitle}>{current.title}</h2>
          <p className={styles.stepBody}>{current.body}</p>
        </div>

        {/* Navigation */}
        <div className={styles.nav}>
          {!isFirst && (
            <button className={styles.navBtnPrev} onClick={goPrev}>Back</button>
          )}
          <button className={styles.navBtnNext} onClick={goNext}>
            {isLast ? 'Done' : 'Next'}
          </button>
        </div>
      </div>
    </div>
  )
}
