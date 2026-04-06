import { useRef, useCallback } from 'react'

export function useVerifyAudio() {
  const ctxRef = useRef<AudioContext | null>(null)
  const lastApproachRef = useRef(0)
  const lastCloseRef = useRef(0)

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

  const playApproaching = useCallback(() => {
    const now = Date.now()
    if (now - lastApproachRef.current < 5000) return
    lastApproachRef.current = now
    playTone(440, 200)
  }, [playTone])

  const playClose = useCallback(() => {
    const now = Date.now()
    if (now - lastCloseRef.current < 3000) return
    lastCloseRef.current = now
    playTone(880, 150)
  }, [playTone])

  return { init, playApproaching, playClose }
}
