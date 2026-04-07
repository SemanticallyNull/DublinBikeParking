import { useRef, useCallback } from 'react'

type WebkitWindow = Window & { webkitAudioContext?: typeof AudioContext }

export function useVerifyAudio() {
  const ctxRef = useRef<AudioContext | null>(null)
  const approachedIds = useRef<Set<string>>(new Set())
  const closeIds = useRef<Set<string>>(new Set())

  const init = useCallback(() => {
    if (!ctxRef.current) {
      const Ctor =
        typeof AudioContext !== 'undefined'
          ? AudioContext
          : (window as WebkitWindow).webkitAudioContext
      if (!Ctor) return
      ctxRef.current = new Ctor()
    }
    const ctx = ctxRef.current
    if (ctx.state === 'suspended') {
      void ctx.resume()
    }
    // Unlock audio on iOS (Safari/Firefox/Chrome on iOS all use WebKit) by
    // playing a silent buffer inside the user-gesture handler. Without this,
    // subsequent oscillator nodes are silent on iOS browsers.
    try {
      const buffer = ctx.createBuffer(1, 1, 22050)
      const source = ctx.createBufferSource()
      source.buffer = buffer
      source.connect(ctx.destination)
      source.start(0)
    } catch {
      // ignore
    }
  }, [])

  const playTone = useCallback((freq: number, durationMs: number) => {
    const ctx = ctxRef.current
    if (!ctx) return
    if (ctx.state === 'suspended') void ctx.resume()

    const now = ctx.currentTime
    const duration = durationMs / 1000

    const osc = ctx.createOscillator()
    const gain = ctx.createGain()
    osc.type = 'sine'
    osc.frequency.setValueAtTime(freq, now)

    // iOS WebKit requires explicit gain scheduling via setValueAtTime /
    // ramps to produce audible output reliably. Setting gain.value alone
    // is silent on iOS Firefox/Safari. Use a short attack/release to also
    // avoid clicks.
    const peak = 0.3
    gain.gain.setValueAtTime(0.0001, now)
    gain.gain.exponentialRampToValueAtTime(peak, now + 0.01)
    gain.gain.setValueAtTime(peak, now + Math.max(duration - 0.02, 0.01))
    gain.gain.exponentialRampToValueAtTime(0.0001, now + duration)

    osc.connect(gain)
    gain.connect(ctx.destination)
    osc.start(now)
    osc.stop(now + duration + 0.02)
  }, [])

  const playApproaching = useCallback((standId: string) => {
    if (approachedIds.current.has(standId)) return
    approachedIds.current.add(standId)
    playTone(440, 200)
  }, [playTone])

  const playClose = useCallback((standId: string) => {
    if (closeIds.current.has(standId)) return
    closeIds.current.add(standId)
    playTone(880, 150)
  }, [playTone])

  return { init, playApproaching, playClose }
}
