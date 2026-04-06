// frontend/src/components/panels/StandDetailPanel.tsx
import { useState } from 'react'
import { reportMissing } from '../../api/actions'
import type { StandFeature, QueryParams } from '../../types'
import styles from './StandDetailPanel.module.css'

const TOASTER_TYPES = ['Wheel Only']

interface Props {
  feature: StandFeature
  queryParams: QueryParams
  onClose: () => void
  onOpenGuide: () => void
}

export function StandDetailPanel({ feature, queryParams, onClose, onOpenGuide }: Props) {
  const p = feature.properties
  const [reportState, setReportState] = useState<'idle' | 'loading' | 'done' | 'error'>('idle')
  const isToaster = TOASTER_TYPES.includes(p.type)
  const googleMapsUrl = `https://www.google.com/maps/dir/?api=1&destination=${feature.geometry.coordinates[1]},${feature.geometry.coordinates[0]}`

  async function handleReportMissing() {
    setReportState('loading')
    try {
      await reportMissing(p.id)
      setReportState('done')
      setTimeout(onClose, 1500)
    } catch {
      setReportState('error')
    }
  }

  return (
    <div className={styles.panel}>
      <div className={styles.hero}>
        <button className={styles.closeBtn} onClick={onClose} aria-label="Close">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
            <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>

        <div className={styles.typeChip}>{p.type}</div>
        <h2 className={styles.name}>{p.name || 'Unnamed stand'}</h2>

        <div className={styles.meta}>
          {p.numberOfStands != null && (
            <span className={styles.metaItem}>
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
                <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z"/><circle cx="12" cy="9" r="2.5"/>
              </svg>
              {p.numberOfStands} {p.numberOfStands === 1 ? 'space' : 'spaces'}
            </span>
          )}
          <span className={p.verified ? styles.verifiedPill : styles.unverifiedPill}>
            {p.verified ? '✓ Verified' : '⚠ Unverified'}
          </span>
        </div>
      </div>

      <div className={styles.body}>
        {isToaster && (
          <div className={styles.toasterWarning}>
            <strong>Wheel-only stand</strong> — you can secure your wheel but not your frame. Consider a Sheffield or hoop stand nearby if possible.
          </div>
        )}

        {p.numberOfStands != null && (
          <div className={styles.stat}>
            <span className={styles.statLabel}>Rack spaces</span>
            <span className={styles.statValue}>{p.numberOfStands}</span>
          </div>
        )}

        <div className={styles.stat}>
          <span className={styles.statLabel}>Source</span>
          <span className={styles.statValueSm}>{p.source || 'Unknown'}</span>
        </div>

        {queryParams.showIDs && (
          <div className={styles.stat}>
            <span className={styles.statLabel}>Stand ID</span>
            <span className={styles.statValueMono}>{p.id}</span>
          </div>
        )}

        {p.publicImageURL && (
          <a className={styles.imageLink} href={p.publicImageURL} target="_blank" rel="noopener noreferrer">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
              <rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/>
            </svg>
            View photo of this stand
          </a>
        )}
      </div>

      <div className={styles.actions}>
        <a className={styles.btnPrimary} href={googleMapsUrl} target="_blank" rel="noopener noreferrer">
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
            <polygon points="3 11 22 2 13 21 11 13 3 11"/>
          </svg>
          Get Directions
        </a>

        {reportState === 'done' ? (
          <div className={styles.reportDone}>✓ Stand reported as missing</div>
        ) : (
          <button
            className={styles.btnDanger}
            onClick={handleReportMissing}
            disabled={reportState === 'loading'}
          >
            {reportState === 'loading' ? 'Reporting…' : reportState === 'error' ? 'Failed — try again' : 'Report Missing'}
          </button>
        )}

        <button className={styles.guideLink} onClick={onOpenGuide}>
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
            <rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/>
          </svg>
          How to lock your bike securely →
        </button>
      </div>
    </div>
  )
}
