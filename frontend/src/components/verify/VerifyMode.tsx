import { useState, useEffect, useRef } from 'react'
import { useStands } from '../../hooks/useStands'
import { useWatchPosition } from '../../hooks/useWatchPosition'
import { useVerifyAudio } from '../../hooks/useVerifyAudio'
import { useVerifySession } from '../../hooks/useVerifySession'
import { VerifyHUD } from './VerifyHUD'
import styles from './VerifyMode.module.css'

export function VerifyMode() {
  const [password, setPassword] = useState('')
  const [started, setStarted] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const wakeLockRef = useRef<WakeLockSentinel | null>(null)

  // Fetch only unverified stands
  const { features, loading } = useStands({ verified: false })

  const audio = useVerifyAudio()
  const session = useVerifySession(password)
  const watch = useWatchPosition(features, session.allExcluded, started)

  // Audio triggers based on distance
  useEffect(() => {
    if (!started || !watch.nearest) return
    const d = watch.nearest.distanceM
    if (d < 30) {
      audio.playClose()
    } else if (d < 100) {
      audio.playApproaching()
    }
  }, [started, watch.nearest, audio])

  // Wake Lock
  useEffect(() => {
    if (!started) return
    async function requestWakeLock() {
      try {
        if ('wakeLock' in navigator) {
          wakeLockRef.current = await navigator.wakeLock.request('screen')
        }
      } catch {
        // Wake Lock not available or denied
      }
    }

    requestWakeLock()

    return () => {
      wakeLockRef.current?.release()
      wakeLockRef.current = null
    }
  }, [started])

  function handleStart() {
    if (!password.trim()) {
      setError('Please enter a password')
      return
    }
    setError(null)
    audio.init()
    setStarted(true)
  }

  function handleExit() {
    setStarted(false)
    wakeLockRef.current?.release()
    wakeLockRef.current = null
  }

  // Handle 401 errors from API calls — shown if password is wrong
  useEffect(() => {
    if (!started) return
    // We'll detect auth errors when verifying the first stand
  }, [started])

  if (!started) {
    const unverifiedCount = features.filter(f => !f.properties.verified).length

    return (
      <div className={styles.entry}>
        <h1 className={styles.title}>Rider Verification Mode</h1>
        <p className={styles.subtitle}>
          Cycle around Dublin and verify bike stands exist
        </p>

        <input
          className={styles.passwordInput}
          type="password"
          placeholder="Enter verification password"
          value={password}
          onChange={e => setPassword(e.target.value)}
          onKeyDown={e => e.key === 'Enter' && handleStart()}
        />

        {error && <div className={styles.error}>{error}</div>}

        <button
          className={styles.startButton}
          onClick={handleStart}
          disabled={loading}
        >
          {loading ? 'Loading stands...' : 'Start Verification Ride'}
        </button>

        {!loading && (
          <div className={styles.standCount}>
            {unverifiedCount} unverified stands to check
          </div>
        )}
      </div>
    )
  }

  return (
    <VerifyHUD
      nearest={watch.nearest}
      heading={watch.heading}
      stats={session.sessionStats}
      gpsError={watch.error}
      onVerify={() => {
        if (watch.nearest) session.verify(watch.nearest.feature.properties.id)
      }}
      onMissing={() => {
        if (watch.nearest) session.markMissing(watch.nearest.feature.properties.id)
      }}
      onSkip={() => {
        if (watch.nearest) session.skip(watch.nearest.feature.properties.id)
      }}
      onExit={handleExit}
    />
  )
}
