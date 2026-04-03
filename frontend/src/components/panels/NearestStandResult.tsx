import type { NearestResult } from '../../hooks/useGeolocation'
import styles from './NearestStandResult.module.css'

interface Props {
  result: NearestResult
  onSelect: () => void
  onDismiss: () => void
}

function formatDist(m: number): string {
  return m < 1000 ? `${Math.round(m)}m` : `${(m / 1000).toFixed(1)}km`
}

export function NearestStandResult({ result, onSelect, onDismiss }: Props) {
  const { feature, distanceM } = result
  const p = feature.properties
  const [lng, lat] = feature.geometry.coordinates
  const mapsUrl = `https://www.google.com/maps/dir/?api=1&destination=${lat},${lng}`

  return (
    <div className={styles.card}>
      <button className={styles.dismiss} onClick={onDismiss} aria-label="Dismiss">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
          <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
        </svg>
      </button>

      <div className={styles.label}>Nearest stand</div>
      <div className={styles.distance}>{formatDist(distanceM)}</div>

      <div className={styles.name}>{p.name || p.type || 'Bike stand'}</div>
      {p.numberOfStands != null && (
        <div className={styles.spaces}>{p.numberOfStands} {p.numberOfStands === 1 ? 'space' : 'spaces'}</div>
      )}

      <div className={styles.actions}>
        <button className={styles.btnDetails} onClick={onSelect}>Details</button>
        <a className={styles.btnDirections} href={mapsUrl} target="_blank" rel="noopener noreferrer">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
            <polygon points="3 11 22 2 13 21 11 13 3 11"/>
          </svg>
          Directions
        </a>
      </div>
    </div>
  )
}
