import type { NearestStand } from '../../hooks/useWatchPosition'
import type { SessionStats } from '../../hooks/useVerifySession'
import styles from './VerifyHUD.module.css'

interface VerifyHUDProps {
  nearest: NearestStand | null
  heading: number | null
  stats: SessionStats
  gpsError: string | null
  onVerify: () => void
  onMissing: () => void
  onSkip: () => void
  onExit: () => void
}

function getArrow(side: string): string {
  switch (side) {
    case 'left': return '\u25C4'
    case 'right': return '\u25BA'
    case 'ahead': return '\u25B2'
    case 'behind': return '\u25BC'
    default: return '\u25CF'
  }
}

function getSideLabel(side: string): string {
  switch (side) {
    case 'left': return 'LEFT'
    case 'right': return 'RIGHT'
    case 'ahead': return 'AHEAD'
    case 'behind': return 'BEHIND'
    default: return ''
  }
}

function getDistanceColor(distanceM: number): string {
  if (distanceM < 5) return styles.arrowGreen
  if (distanceM < 25) return styles.arrowYellow
  return styles.arrowDim
}

function formatDistance(m: number): string {
  if (m < 1000) return `~${Math.round(m)}m`
  return `~${(m / 1000).toFixed(1)}km`
}

export function VerifyHUD({
  nearest,
  heading,
  stats,
  gpsError,
  onVerify,
  onMissing,
  onSkip,
  onExit,
}: VerifyHUDProps) {
  const statsText = [
    stats.verified > 0 && `${stats.verified} verified`,
    stats.missing > 0 && `${stats.missing} missing`,
    stats.skipped > 0 && `${stats.skipped} skipped`,
  ]
    .filter(Boolean)
    .join(', ')

  if (gpsError) {
    return (
      <div className={styles.hud}>
        <div className={styles.topBar}>
          <button className={styles.exitButton} onClick={onExit}>X Exit</button>
          <span className={styles.stats}>{statsText || 'Starting...'}</span>
        </div>
        <div className={styles.gpsError}>{gpsError}</div>
      </div>
    )
  }

  if (!nearest) {
    return (
      <div className={styles.hud}>
        <div className={styles.topBar}>
          <button className={styles.exitButton} onClick={onExit}>X Exit</button>
          <span className={styles.stats}>{statsText || '0 verified'}</span>
        </div>
        <div className={styles.emptyState}>
          <div className={styles.emptyTitle}>No more stands nearby</div>
          <ul className={styles.summaryList}>
            <li>Verified: {stats.verified}</li>
            <li>Missing: {stats.missing}</li>
            <li>Skipped: {stats.skipped}</li>
          </ul>
          <button className={styles.finishButton} onClick={onExit}>Finish Ride</button>
        </div>
      </div>
    )
  }

  const { feature, distanceM, side, compass } = nearest
  const useCompass = heading == null

  return (
    <div className={styles.hud}>
      <div className={styles.topBar}>
        <button className={styles.exitButton} onClick={onExit}>X Exit</button>
        <span className={styles.stats}>{statsText || 'Starting...'}</span>
      </div>

      <div className={styles.directionArea}>
        {useCompass ? (
          <div className={styles.compassText}>
            Stand is to the {compass}
          </div>
        ) : (
          <>
            <div className={`${styles.arrow} ${getDistanceColor(distanceM)}`}>
              {getArrow(side)}
            </div>
            <div className={styles.compassText}>{getSideLabel(side)}</div>
          </>
        )}

        <div className={styles.standInfo}>
          <div className={styles.standName}>{feature.properties.name || 'Unnamed Stand'}</div>
          <div className={styles.standType}>{feature.properties.type}</div>
        </div>

        <div className={styles.distance}>{formatDistance(distanceM)}</div>
      </div>

      <div className={styles.actionButtons}>
        <button className={styles.verifyButton} onClick={onVerify}>VERIFY</button>
        <button className={styles.missingButton} onClick={onMissing}>MISSING</button>
      </div>
      <button className={styles.skipButton} onClick={onSkip}>SKIP</button>
    </div>
  )
}
