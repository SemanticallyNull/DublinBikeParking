import { useRef, useCallback } from 'react'

export function useVerifyAudio() {
  const ctxRef = useRef<AudioContext | null>(null)
  const approachedIds = useRef<Set<string>>(new Set())
  const closeIds = useRef<Set<string>>(new Set())

  const init = useCallback(() => {
    if (!ctxRef.current) {
      ctxRef.current = new AudioContext()
    }
    if (ctxRef.current.state === 'suspended') {
      ctxRef.current.resume()
    }
  }, [])

  const playTone = useCallback((freq: number, durationMs: number) => {
    const ctx = ctxRef.current
    if (!ctx) return
    if (ctx.state === 'suspended') ctx.resume()

    const osc = ctx.createOscillator()
    const gain = ctx.createGain()
    osc.type = 'sine'
    osc.frequency.value = freq
    gain.gain.value = 0.3
    osc.connect(gain)
    gain.connect(ctx.destination)
    osc.start()
    osc.stop(ctx.currentTime + durationMs / 1000)
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
