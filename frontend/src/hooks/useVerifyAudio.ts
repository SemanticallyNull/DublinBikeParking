import { useRef, useCallback, useMemo } from 'react'

// Generate a tiny WAV file (PCM 16-bit mono) as a Blob URL containing a sine
// tone. Using HTMLAudioElement with a pre-generated WAV is the most reliable
// way to play sound on iOS browsers (Safari, Firefox, Chrome — all WebKit).
// The Web Audio API is unreliable on Firefox iOS in particular, where
// oscillator nodes frequently produce no audible output even after the
// context is resumed inside a user gesture.
function makeToneWavUrl(freq: number, durationMs: number, volume = 0.3): string {
  const sampleRate = 22050
  const numSamples = Math.floor((durationMs / 1000) * sampleRate)
  const bytesPerSample = 2
  const blockAlign = bytesPerSample
  const byteRate = sampleRate * blockAlign
  const dataSize = numSamples * bytesPerSample
  const buffer = new ArrayBuffer(44 + dataSize)
  const view = new DataView(buffer)

  const writeStr = (offset: number, s: string) => {
    for (let i = 0; i < s.length; i++) view.setUint8(offset + i, s.charCodeAt(i))
  }

  // RIFF header
  writeStr(0, 'RIFF')
  view.setUint32(4, 36 + dataSize, true)
  writeStr(8, 'WAVE')
  // fmt chunk
  writeStr(12, 'fmt ')
  view.setUint32(16, 16, true) // chunk size
  view.setUint16(20, 1, true) // PCM format
  view.setUint16(22, 1, true) // mono
  view.setUint32(24, sampleRate, true)
  view.setUint32(28, byteRate, true)
  view.setUint16(32, blockAlign, true)
  view.setUint16(34, 16, true) // bits per sample
  // data chunk
  writeStr(36, 'data')
  view.setUint32(40, dataSize, true)

  // PCM samples with a short linear attack/release envelope to avoid clicks.
  const attack = Math.min(0.01 * sampleRate, numSamples / 2)
  const release = Math.min(0.02 * sampleRate, numSamples / 2)
  for (let i = 0; i < numSamples; i++) {
    let env = 1
    if (i < attack) env = i / attack
    else if (i > numSamples - release) env = (numSamples - i) / release
    const sample = Math.sin((2 * Math.PI * freq * i) / sampleRate) * volume * env
    view.setInt16(44 + i * 2, Math.max(-1, Math.min(1, sample)) * 0x7fff, true)
  }

  const blob = new Blob([buffer], { type: 'audio/wav' })
  return URL.createObjectURL(blob)
}

export function useVerifyAudio() {
  const approachedIds = useRef<Set<string>>(new Set())
  const closeIds = useRef<Set<string>>(new Set())

  const { approachUrl, closeUrl } = useMemo(
    () => ({
      approachUrl: makeToneWavUrl(440, 200),
      closeUrl: makeToneWavUrl(880, 150),
    }),
    []
  )

  const approachAudioRef = useRef<HTMLAudioElement | null>(null)
  const closeAudioRef = useRef<HTMLAudioElement | null>(null)

  const getApproach = useCallback(() => {
    if (!approachAudioRef.current) {
      const a = new Audio(approachUrl)
      a.preload = 'auto'
      approachAudioRef.current = a
    }
    return approachAudioRef.current
  }, [approachUrl])

  const getClose = useCallback(() => {
    if (!closeAudioRef.current) {
      const a = new Audio(closeUrl)
      a.preload = 'auto'
      closeAudioRef.current = a
    }
    return closeAudioRef.current
  }, [closeUrl])

  // Must be called from inside a user-gesture handler to "unlock" audio on
  // iOS. We briefly play each element muted, which primes it so subsequent
  // play() calls (from non-gesture contexts like GPS updates) succeed.
  const init = useCallback(() => {
    const unlock = (a: HTMLAudioElement) => {
      try {
        a.muted = true
        const p = a.play()
        if (p && typeof p.then === 'function') {
          p.then(() => {
            a.pause()
            a.currentTime = 0
            a.muted = false
          }).catch(() => {
            a.muted = false
          })
        } else {
          a.pause()
          a.currentTime = 0
          a.muted = false
        }
      } catch {
        a.muted = false
      }
    }
    unlock(getApproach())
    unlock(getClose())
  }, [getApproach, getClose])

  const play = useCallback((a: HTMLAudioElement) => {
    try {
      a.currentTime = 0
      const p = a.play()
      if (p && typeof p.catch === 'function') {
        p.catch(() => {
          // Autoplay blocked — ignore
        })
      }
    } catch {
      // ignore
    }
  }, [])

  const playApproaching = useCallback(
    (standId: string) => {
      if (approachedIds.current.has(standId)) return
      approachedIds.current.add(standId)
      play(getApproach())
    },
    [play, getApproach]
  )

  const playClose = useCallback(
    (standId: string) => {
      if (closeIds.current.has(standId)) return
      closeIds.current.add(standId)
      play(getClose())
    },
    [play, getClose]
  )

  const playTest = useCallback(() => {
    // Runs from a click handler, so this is a user gesture — unlock first.
    init()
    play(getApproach())
    setTimeout(() => play(getClose()), 300)
  }, [init, play, getApproach, getClose])

  return { init, playApproaching, playClose, playTest }
}
