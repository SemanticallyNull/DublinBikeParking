# Bike Locking Guide Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add an animated, 6-step full-screen guide modal teaching users how to lock their bike securely, accessible from the header, on first visit, and from stand detail panels.

**Architecture:** A new `LockingGuide` modal component owns all rendering and animation logic. `App.tsx` owns `guideOpen` state and the localStorage first-visit check. Guide open/close callbacks flow down through `AppShell` to `AppHeader` and `StandDetailPanel`.

**Tech Stack:** React 18, TypeScript, CSS Modules, inline SVG illustrations, CSS keyframe animations (no external animation library).

> **Note:** This project has no frontend test infrastructure (no vitest/jest). Skip unit test steps; use `npm run build` and browser verification instead.

---

## File Map

| Action | Path | Responsibility |
|--------|------|----------------|
| Create | `frontend/src/components/guide/steps.tsx` | Step content array (title, body, SVG illustration component) |
| Create | `frontend/src/components/guide/LockingGuide.tsx` | Modal component — layout, animation, focus trap, progress bar |
| Create | `frontend/src/components/guide/LockingGuide.module.css` | All styles for the modal |
| Modify | `frontend/src/App.tsx` | Add `guideOpen` state + first-visit localStorage check |
| Modify | `frontend/src/components/layout/AppShell.tsx` | Accept guide props, render `<LockingGuide>`, forward `onOpenGuide` |
| Modify | `frontend/src/components/layout/AppHeader.tsx` | Add "?" guide button |
| Modify | `frontend/src/components/panels/StandDetailPanel.tsx` | Add "How to lock your bike →" link + accept `onOpenGuide` prop |

---

## Task 1: Create step content

**Files:**
- Create: `frontend/src/components/guide/steps.tsx`

- [ ] **Step 1: Create `steps.tsx` with the step data array**

```tsx
// frontend/src/components/guide/steps.tsx

export interface GuideStep {
  title: string
  body: string
  illustration: () => JSX.Element
}

const IllustrationCheckStand = () => (
  <svg viewBox="0 0 120 120" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
    {/* Sheffield stand arch */}
    <path d="M20 95 L20 42 Q20 18 45 18 L75 18 Q100 18 100 42 L100 95"
      stroke="white" strokeWidth="5" strokeLinecap="round" strokeLinejoin="round"/>
    {/* Ground rail */}
    <line x1="8" y1="95" x2="112" y2="95" stroke="#6C96BB" strokeWidth="5" strokeLinecap="round"/>
    {/* Magnifying glass */}
    <circle cx="72" cy="62" r="20" stroke="#5A94C8" strokeWidth="4"/>
    <line x1="86" y1="76" x2="100" y2="90" stroke="#5A94C8" strokeWidth="4" strokeLinecap="round"/>
    {/* Warning marks on stand */}
    <circle cx="45" cy="50" r="4" fill="#DC2626"/>
    <circle cx="45" cy="64" r="4" fill="#DC2626"/>
  </svg>
)

const IllustrationChooseLock = () => (
  <svg viewBox="0 0 120 120" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
    {/* U-lock shackle */}
    <path d="M22 80 L22 48 Q22 26 38 26 Q54 26 54 48 L54 80"
      stroke="#4ADE80" strokeWidth="6" strokeLinecap="round"/>
    {/* U-lock body */}
    <rect x="14" y="76" width="48" height="22" rx="6"
      fill="rgba(74,222,128,0.12)" stroke="#4ADE80" strokeWidth="4"/>
    {/* Good label */}
    <text x="38" y="115" textAnchor="middle" fontSize="10" fill="#4ADE80"
      fontFamily="system-ui, sans-serif" fontWeight="700">U-LOCK</text>
    {/* Cable lock squiggle */}
    <path d="M76 30 Q82 45 78 58 Q74 70 82 80 Q90 90 86 100"
      stroke="#DC2626" strokeWidth="5" strokeLinecap="round"/>
    {/* X mark */}
    <line x1="78" y1="110" x2="90" y2="118" stroke="#DC2626" strokeWidth="3" strokeLinecap="round"/>
    <line x1="90" y1="110" x2="78" y2="118" stroke="#DC2626" strokeWidth="3" strokeLinecap="round"/>
    <text x="84" y="115" textAnchor="middle" fontSize="10" fill="#DC2626"
      fontFamily="system-ui, sans-serif" fontWeight="700" dy="10">CABLE</text>
  </svg>
)

const IllustrationLockFrame = () => (
  <svg viewBox="0 0 120 120" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
    {/* Bike frame (simplified triangle) */}
    <path d="M20 90 L60 30 L90 90 Z" stroke="white" strokeWidth="4" strokeLinejoin="round"/>
    {/* Rear wheel hint */}
    <circle cx="20" cy="90" r="16" stroke="#6C96BB" strokeWidth="3"/>
    {/* Front wheel hint */}
    <circle cx="90" cy="90" r="16" stroke="#6C96BB" strokeWidth="3"/>
    {/* Sheffield stand */}
    <path d="M50 105 L50 70 Q50 58 60 58 Q70 58 70 70 L70 105"
      stroke="#5A94C8" strokeWidth="5" strokeLinecap="round"/>
    {/* U-lock through frame + stand */}
    <path d="M44 96 L44 78 Q44 68 54 68 L66 68 Q76 68 76 78 L76 96"
      stroke="#4ADE80" strokeWidth="5" strokeLinecap="round"/>
    <rect x="38" y="92" width="44" height="16" rx="5"
      fill="rgba(74,222,128,0.15)" stroke="#4ADE80" strokeWidth="3"/>
  </svg>
)

const IllustrationSecondaryLock = () => (
  <svg viewBox="0 0 120 120" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
    {/* Rear wheel */}
    <circle cx="60" cy="70" r="32" stroke="white" strokeWidth="4"/>
    <circle cx="60" cy="70" r="6" fill="#6C96BB"/>
    {/* Spokes */}
    <line x1="60" y1="38" x2="60" y2="102" stroke="#6C96BB" strokeWidth="2" opacity="0.5"/>
    <line x1="28" y1="70" x2="92" y2="70" stroke="#6C96BB" strokeWidth="2" opacity="0.5"/>
    {/* Sheffield stand */}
    <path d="M36 108 L36 80 Q36 68 48 68 L72 68 Q84 68 84 80 L84 108"
      stroke="#5A94C8" strokeWidth="5" strokeLinecap="round"/>
    {/* Cable lock loop around wheel + stand */}
    <path d="M44 100 Q34 90 38 70 Q42 50 60 46 Q78 42 86 62 Q92 78 80 92"
      stroke="#5A94C8" strokeWidth="4" strokeLinecap="round" strokeDasharray="6 3"/>
    {/* Lock clasp */}
    <rect x="74" y="90" width="14" height="11" rx="3"
      fill="rgba(90,148,200,0.2)" stroke="#5A94C8" strokeWidth="3"/>
  </svg>
)

const IllustrationFillLock = () => (
  <svg viewBox="0 0 120 120" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
    {/* U-lock shackle */}
    <path d="M30 85 L30 48 Q30 28 50 28 Q70 28 70 48 L70 85"
      stroke="white" strokeWidth="6" strokeLinecap="round"/>
    {/* Lock body */}
    <rect x="22" y="80" width="56" height="26" rx="7"
      fill="rgba(90,148,200,0.12)" stroke="white" strokeWidth="4"/>
    {/* Keyhole */}
    <circle cx="50" cy="90" r="5" fill="#6C96BB"/>
    <rect x="47" y="90" width="6" height="8" rx="2" fill="#6C96BB"/>
    {/* Stand bar inside lock — showing tight fit */}
    <rect x="42" y="45" width="16" height="40" rx="4"
      fill="rgba(74,222,128,0.2)" stroke="#4ADE80" strokeWidth="3"/>
    {/* Arrow indicators showing tight fit */}
    <line x1="34" y1="55" x2="42" y2="55" stroke="#4ADE80" strokeWidth="2.5" strokeLinecap="round"/>
    <line x1="58" y1="55" x2="66" y2="55" stroke="#4ADE80" strokeWidth="2.5" strokeLinecap="round"/>
    <line x1="39" y1="52" x2="42" y2="55" stroke="#4ADE80" strokeWidth="2.5" strokeLinecap="round"/>
    <line x1="39" y1="58" x2="42" y2="55" stroke="#4ADE80" strokeWidth="2.5" strokeLinecap="round"/>
    <line x1="63" y1="52" x2="60" y2="55" stroke="#4ADE80" strokeWidth="2.5" strokeLinecap="round"/>
    <line x1="63" y1="58" x2="60" y2="55" stroke="#4ADE80" strokeWidth="2.5" strokeLinecap="round"/>
  </svg>
)

const IllustrationPickSpot = () => (
  <svg viewBox="0 0 120 120" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
    {/* Sun / visibility rays */}
    <circle cx="60" cy="38" r="14" fill="rgba(90,148,200,0.2)" stroke="#5A94C8" strokeWidth="3"/>
    <line x1="60" y1="16" x2="60" y2="10" stroke="#5A94C8" strokeWidth="3" strokeLinecap="round"/>
    <line x1="60" y1="60" x2="60" y2="66" stroke="#5A94C8" strokeWidth="3" strokeLinecap="round"/>
    <line x1="38" y1="38" x2="32" y2="38" stroke="#5A94C8" strokeWidth="3" strokeLinecap="round"/>
    <line x1="82" y1="38" x2="88" y2="38" stroke="#5A94C8" strokeWidth="3" strokeLinecap="round"/>
    <line x1="44" y1="24" x2="40" y2="20" stroke="#5A94C8" strokeWidth="3" strokeLinecap="round"/>
    <line x1="76" y1="52" x2="80" y2="56" stroke="#5A94C8" strokeWidth="3" strokeLinecap="round"/>
    <line x1="76" y1="24" x2="80" y2="20" stroke="#5A94C8" strokeWidth="3" strokeLinecap="round"/>
    <line x1="44" y1="52" x2="40" y2="56" stroke="#5A94C8" strokeWidth="3" strokeLinecap="round"/>
    {/* Map pin */}
    <path d="M60 115 C60 115 38 88 38 74 C38 61.3 48.1 52 60 52 C71.9 52 82 61.3 82 74 C82 88 60 115 60 115Z"
      fill="rgba(90,148,200,0.15)" stroke="white" strokeWidth="4"/>
    <circle cx="60" cy="74" r="8" fill="white" fillOpacity="0.3" stroke="white" strokeWidth="3"/>
  </svg>
)

export const GUIDE_STEPS: GuideStep[] = [
  {
    title: 'Check the stand',
    body: 'Look for signs of tampering — loose bolts, damaged metalwork, anything that looks interfered with. If it looks dodgy, find another stand.',
    illustration: IllustrationCheckStand,
  },
  {
    title: 'Choose the right lock',
    body: 'Use a Sold Secure rated U-lock or D-lock as your primary. Avoid cable locks as your only lock — they can be cut in seconds.',
    illustration: IllustrationChooseLock,
  },
  {
    title: 'Lock your frame',
    body: 'Pass the lock through your frame (not just the wheel) and around the stand. A wheel-only lock leaves your frame behind.',
    illustration: IllustrationLockFrame,
  },
  {
    title: 'Add a secondary lock',
    body: 'Use a cable or chain to also secure your rear wheel. Two locks double the deterrent and the time cost for a thief.',
    illustration: IllustrationSecondaryLock,
  },
  {
    title: 'Fill the lock',
    body: 'Leave as little space as possible inside your U-lock. A tight fit makes it much harder to lever open.',
    illustration: IllustrationFillLock,
  },
  {
    title: 'Pick your spot',
    body: 'Lock in a visible, well-lit location. Thieves prefer cover — busy, open areas are much safer.',
    illustration: IllustrationPickSpot,
  },
]
```

- [ ] **Step 2: Verify TypeScript compiles**

```bash
cd frontend && npx tsc --noEmit
```

Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/guide/steps.tsx
git commit -m "feat: add bike locking guide step content"
```

---

## Task 2: Create LockingGuide modal component

**Files:**
- Create: `frontend/src/components/guide/LockingGuide.tsx`
- Create: `frontend/src/components/guide/LockingGuide.module.css`

- [ ] **Step 1: Create `LockingGuide.module.css`**

```css
/* frontend/src/components/guide/LockingGuide.module.css */

.backdrop {
  position: fixed;
  inset: 0;
  background: rgba(9, 21, 56, 0.85);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: fadeIn 0.2s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to   { opacity: 1; }
}

.card {
  background: var(--surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  width: min(480px, 90vw);
  max-height: 90dvh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;
}

/* Progress bar */
.progressBar {
  height: 3px;
  background: var(--border);
  flex-shrink: 0;
}

.progressFill {
  height: 100%;
  background: var(--gold);
  transition: width 0.3s ease;
}

/* Top row: step counter + skip */
.topRow {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 20px 0;
  flex-shrink: 0;
}

.stepCounter {
  font-size: 12px;
  font-weight: 600;
  color: var(--muted);
  letter-spacing: 0.04em;
}

.skipBtn {
  background: none;
  border: none;
  color: var(--muted);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  transition: color 0.15s, background 0.15s;
}

.skipBtn:hover {
  color: var(--text);
  background: var(--border);
}

/* Animated step content wrapper */
.stepContent {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 28px 20px;
  overflow: hidden;
  min-height: 0;
}

.stepContentForward {
  animation: slideInRight 0.25s ease-in-out;
}

.stepContentBack {
  animation: slideInLeft 0.25s ease-in-out;
}

@keyframes slideInRight {
  from { transform: translateX(60px); opacity: 0; }
  to   { transform: translateX(0);    opacity: 1; }
}

@keyframes slideInLeft {
  from { transform: translateX(-60px); opacity: 0; }
  to   { transform: translateX(0);     opacity: 1; }
}

/* Illustration */
.illustration {
  width: 160px;
  height: 160px;
  flex-shrink: 0;
  background: var(--navy);
  border-radius: var(--radius);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  margin-bottom: 20px;
}

.illustration svg {
  width: 100%;
  height: 100%;
}

/* Step text */
.stepTitle {
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 26px;
  letter-spacing: -0.01em;
  color: var(--navy);
  text-align: center;
  margin-bottom: 10px;
  line-height: 1.1;
}

.stepBody {
  font-size: 14px;
  line-height: 1.6;
  color: var(--muted);
  text-align: center;
  max-width: 340px;
}

/* Navigation buttons */
.nav {
  display: flex;
  gap: 10px;
  padding: 0 20px 20px;
  flex-shrink: 0;
}

.navBtnPrev {
  flex: 1;
  background: none;
  border: 1px solid var(--border);
  color: var(--muted);
  font-family: var(--font-display);
  font-size: 14px;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  padding: 12px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
}

.navBtnPrev:hover {
  background: var(--border);
  color: var(--text);
}

.navBtnNext {
  flex: 2;
  background: var(--gold);
  border: none;
  color: var(--white);
  font-family: var(--font-display);
  font-size: 14px;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  padding: 12px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background 0.15s, box-shadow 0.15s;
  box-shadow: 0 2px 12px var(--gold-glow);
}

.navBtnNext:hover {
  background: var(--gold-hover);
  box-shadow: 0 4px 16px var(--gold-glow);
}

/* Mobile adjustments */
@media (max-width: 480px) {
  .illustration {
    width: 130px;
    height: 130px;
  }

  .stepTitle {
    font-size: 22px;
  }

  .stepBody {
    font-size: 13px;
  }
}
```

- [ ] **Step 2: Create `LockingGuide.tsx`**

```tsx
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
    if (open) setStep(0)
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
        <div className={styles.progressBar} role="progressbar" aria-valuenow={step + 1} aria-valuemax={total}>
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
```

- [ ] **Step 3: Verify TypeScript compiles**

```bash
cd frontend && npx tsc --noEmit
```

Expected: no errors.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/guide/LockingGuide.tsx frontend/src/components/guide/LockingGuide.module.css
git commit -m "feat: add LockingGuide modal component"
```

---

## Task 3: Wire up App.tsx

**Files:**
- Modify: `frontend/src/App.tsx`

- [ ] **Step 1: Update `App.tsx`**

Add `guideOpen` state and first-visit localStorage check. Replace the existing file content with:

```tsx
// frontend/src/App.tsx
import { useEffect, useState } from 'react'
import { AppShell } from './components/layout/AppShell'
import { useStands } from './hooks/useStands'
import { useQueryParams } from './hooks/useQueryParams'
import { useGeolocation } from './hooks/useGeolocation'
import type { StandFeature } from './types'

const GUIDE_SEEN_KEY = 'guide_seen'

export default function App() {
  const params = useQueryParams()
  const { features, loading, reload } = useStands({ verified: params.verified, pendingReview: params.pendingReview })
  const [selectedId, setSelectedId] = useState<string | null>(null)
  const [placementMode, setPlacementMode] = useState(false)
  const [addPin, setAddPin] = useState<{ lat: number; lng: number } | null>(null)
  const [guideOpen, setGuideOpen] = useState(false)

  const geo = useGeolocation(features)

  // Auto-open on first visit
  useEffect(() => {
    if (!localStorage.getItem(GUIDE_SEEN_KEY)) {
      localStorage.setItem(GUIDE_SEEN_KEY, 'true')
      setGuideOpen(true)
    }
  }, [])

  const selectedStand = selectedId
    ? (features.find(f => f.properties.id === selectedId) ?? null)
    : null

  function handleSelect(feature: StandFeature) {
    setSelectedId(feature.properties.id)
    setPlacementMode(false)
  }

  function handleClose() {
    setSelectedId(null)
  }

  function handleAddStand() {
    const next = !placementMode
    setPlacementMode(next)
    setSelectedId(null)
    if (!next) setAddPin(null)
  }

  function handleMapClick(lat: number, lng: number) {
    if (placementMode) {
      setAddPin({ lat, lng })
      setPlacementMode(false)
    }
  }

  function handleAddSuccess() {
    setAddPin(null)
    reload()
  }

  function handleAddCancel() {
    setAddPin(null)
  }

  return (
    <AppShell
      features={features}
      loading={loading}
      selectedStand={selectedStand}
      placementMode={placementMode}
      queryParams={params}
      userPosition={geo.position}
      geoLoading={geo.loading}
      nearestResult={geo.nearest}
      addPin={addPin}
      guideOpen={guideOpen}
      onSelect={handleSelect}
      onClose={handleClose}
      onAddStand={handleAddStand}
      onLocate={geo.locate}
      onDismissNearest={geo.clear}
      onMapClick={handleMapClick}
      onAddSuccess={handleAddSuccess}
      onAddCancel={handleAddCancel}
      onOpenGuide={() => setGuideOpen(true)}
      onCloseGuide={() => setGuideOpen(false)}
    />
  )
}
```

- [ ] **Step 2: Verify TypeScript compiles**

```bash
cd frontend && npx tsc --noEmit
```

Expected: type errors for `guideOpen`/`onOpenGuide`/`onCloseGuide` on AppShell (not yet added) — that's fine, proceed.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/App.tsx
git commit -m "feat: add guideOpen state and first-visit check to App"
```

---

## Task 4: Update AppShell

**Files:**
- Modify: `frontend/src/components/layout/AppShell.tsx`

- [ ] **Step 1: Update `AppShell.tsx`**

Add guide props to the interface, pass `onOpenGuide` down to `AppHeader` and `StandDetailPanel`, and render `<LockingGuide>`. Replace existing file content with:

```tsx
// frontend/src/components/layout/AppShell.tsx
import { useEffect } from 'react'
import { AppHeader } from './AppHeader'
import { MapView } from '../map/MapView'
import { StandDetailPanel } from '../panels/StandDetailPanel'
import { AddStandForm } from '../panels/AddStandForm'
import { BottomSheet } from '../ui/BottomSheet'
import { FloatingButton } from '../ui/FloatingButton'
import { NearestStandResult } from '../panels/NearestStandResult'
import { LockingGuide } from '../guide/LockingGuide'
import type { StandFeature, QueryParams } from '../../types'
import type { NearestResult } from '../../hooks/useGeolocation'
import styles from './AppShell.module.css'

interface Props {
  features: StandFeature[]
  loading: boolean
  selectedStand: StandFeature | null
  placementMode: boolean
  queryParams: QueryParams
  userPosition: GeolocationPosition | null
  geoLoading: boolean
  nearestResult: NearestResult | null
  addPin: { lat: number; lng: number } | null
  guideOpen: boolean
  onSelect: (feature: StandFeature) => void
  onClose: () => void
  onAddStand: () => void
  onLocate: () => void
  onDismissNearest: () => void
  onMapClick: (lat: number, lng: number) => void
  onAddSuccess: () => void
  onAddCancel: () => void
  onOpenGuide: () => void
  onCloseGuide: () => void
}

export function AppShell({
  features,
  loading,
  selectedStand,
  placementMode,
  queryParams,
  userPosition,
  geoLoading,
  nearestResult,
  addPin,
  guideOpen,
  onSelect,
  onClose,
  onAddStand,
  onLocate,
  onDismissNearest,
  onMapClick,
  onAddSuccess,
  onAddCancel,
  onOpenGuide,
  onCloseGuide,
}: Props) {
  useEffect(() => {
    if (selectedStand || addPin) {
      document.body.classList.add('panel-open')
    } else {
      document.body.classList.remove('panel-open')
    }
    return () => document.body.classList.remove('panel-open')
  }, [selectedStand, addPin])

  return (
    <div className={styles.shell}>
      {placementMode && (
        <div className={styles.placementBanner}>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
            <circle cx="12" cy="12" r="9"/><circle cx="12" cy="12" r="3"/>
            <line x1="12" y1="2" x2="12" y2="5"/><line x1="12" y1="19" x2="12" y2="22"/>
            <line x1="2" y1="12" x2="5" y2="12"/><line x1="19" y1="12" x2="22" y2="12"/>
          </svg>
          Tap the map where the stand is
          <button className={styles.placementCancel} onClick={onAddStand}>Cancel</button>
        </div>
      )}

      <AppHeader
        standCount={features.length}
        placementMode={placementMode}
        onAddStand={onAddStand}
        onFindNearest={onLocate}
        onOpenGuide={onOpenGuide}
      />

      <div className={styles.body}>
        {loading && (
          <div className={styles.loadingPill}>Loading stands…</div>
        )}

        <div className={styles.mapArea}>
          <MapView
            features={features}
            selectedId={selectedStand?.properties.id ?? null}
            placementMode={placementMode}
            userPosition={userPosition}
            addPin={addPin}
            onSelect={onSelect}
            onMapClick={onMapClick}
          />

          {nearestResult && !selectedStand && (
            <NearestStandResult
              result={nearestResult}
              onSelect={() => onSelect(nearestResult.feature)}
              onDismiss={onDismissNearest}
            />
          )}

          <FloatingButton loading={geoLoading} onClick={onLocate} />
        </div>

        {(selectedStand || addPin) && (
          <aside className={styles.sidePanel}>
            {addPin
              ? <AddStandForm lat={addPin.lat} lng={addPin.lng} onSuccess={onAddSuccess} onCancel={onAddCancel} />
              : <StandDetailPanel feature={selectedStand!} queryParams={queryParams} onClose={onClose} onOpenGuide={onOpenGuide} />
            }
          </aside>
        )}
      </div>

      <div className={styles.mobileOnly}>
        <BottomSheet open={!!(selectedStand || addPin)} onClose={addPin ? onAddCancel : onClose}>
          {addPin
            ? <div className="bottom-sheet-content"><AddStandForm lat={addPin.lat} lng={addPin.lng} onSuccess={onAddSuccess} onCancel={onAddCancel} /></div>
            : selectedStand
              ? <div className="bottom-sheet-content"><StandDetailPanel feature={selectedStand} queryParams={queryParams} onClose={onClose} onOpenGuide={onOpenGuide} /></div>
              : null
          }
        </BottomSheet>
      </div>

      <LockingGuide open={guideOpen} onClose={onCloseGuide} />
    </div>
  )
}
```

- [ ] **Step 2: Verify TypeScript compiles**

```bash
cd frontend && npx tsc --noEmit
```

Expected: type errors for `onOpenGuide` on `AppHeader` and `StandDetailPanel` (not yet added) — proceed.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/layout/AppShell.tsx
git commit -m "feat: wire LockingGuide into AppShell"
```

---

## Task 5: Add guide button to AppHeader

**Files:**
- Modify: `frontend/src/components/layout/AppHeader.tsx`
- Modify: `frontend/src/components/layout/AppHeader.module.css`

- [ ] **Step 1: Add `onOpenGuide` prop and "?" button to `AppHeader.tsx`**

Replace existing file content with:

```tsx
// frontend/src/components/layout/AppHeader.tsx
import styles from './AppHeader.module.css'

interface Props {
  standCount: number
  placementMode: boolean
  onAddStand: () => void
  onFindNearest: () => void
  onOpenGuide: () => void
}

const LogoSvg = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 953.24 935.45">
    <path fill="white" d="M63.19,730.15V123.82h257.05c27.89,0,53.65,5.7,77.29,17.08,23.62,11.4,43.98,26.77,61.06,46.12,17.08,19.36,30.32,40.99,39.71,64.9,9.39,23.91,14.09,48.4,14.09,73.44,0,33.6-7.69,65.9-23.06,96.93-15.37,31.04-37.15,56.23-65.33,75.58-28.18,19.36-61.06,29.04-98.64,29.04h-144.32v203.25H63.19ZM181.04,423.57h136.64c14.8,0,27.89-3.98,39.28-11.96,11.38-7.97,20.35-19.5,26.9-34.59,6.54-15.08,9.82-32.3,9.82-51.67,0-21.06-3.84-38.86-11.53-53.37-7.69-14.52-17.65-25.62-29.89-33.31-12.25-7.69-25.49-11.53-39.71-11.53h-131.51v196.42Z"/>
    <g><circle cx="405.26" cy="800.75" r="119.19" style={{fill:'none',stroke:'white',strokeMiterlimit:10,strokeWidth:'31px'}}/><circle cx="818.54" cy="800.75" r="119.19" style={{fill:'none',stroke:'white',strokeMiterlimit:10,strokeWidth:'31px'}}/><path d="M405.26,800.75c17.83-54.63,29.44-99.47,38.88-129.76,27.97-89.71,35.19-104.05,43.21-115.7,9.88-14.37,27.28-16.14,44.01-14.37" style={{fill:'none',stroke:'white',strokeLinecap:'round',strokeMiterlimit:10,strokeWidth:'31px'}}/><polygon points="454.25 638.88 643.09 800.75 818.54 798.23 688.82 633.71 454.25 638.88" style={{fill:'none',stroke:'white',strokeLinecap:'round',strokeLinejoin:'round',strokeWidth:'31px'}}/><path d="M701.41,596.37c-17.82,64.73-35.65,129.47-53.47,194.2" style={{fill:'none',stroke:'white',strokeLinecap:'round',strokeLinejoin:'round',strokeWidth:'31px'}}/><line x1="659.85" y1="589.17" x2="733.53" y2="589.17" style={{fill:'none',stroke:'white',strokeLinecap:'round',strokeMiterlimit:10,strokeWidth:'31px'}}/></g>
  </svg>
)

export function AppHeader({ standCount, placementMode, onAddStand, onFindNearest, onOpenGuide }: Props) {
  return (
    <header className={styles.header}>
      <div className={styles.logo}>
        <div className={styles.logoMark}><LogoSvg /></div>
        <div className={styles.logoText}>
          <span className={styles.logoName}>Bike Parking</span>
          <span className={styles.logoCity}>Dublin</span>
        </div>
      </div>

      <div className={styles.gap} />

      {standCount > 0 && (
        <span className={styles.count}>{standCount} stands</span>
      )}

      <button
        className={`${styles.addBtn} ${placementMode ? styles.addBtnActive : ''}`}
        onClick={onAddStand}
        title="Add a stand"
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
          <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
        <span>Add Stand</span>
      </button>

      <button className={styles.findBtn} onClick={onFindNearest} title="Find nearest stand">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
          <circle cx="12" cy="12" r="9"/><circle cx="12" cy="12" r="3"/>
          <line x1="12" y1="2" x2="12" y2="5"/><line x1="12" y1="19" x2="12" y2="22"/>
          <line x1="2" y1="12" x2="5" y2="12"/><line x1="19" y1="12" x2="22" y2="12"/>
        </svg>
        <span>Find Nearest</span>
      </button>

      <button className={styles.guideBtn} onClick={onOpenGuide} title="How to lock your bike">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
          <circle cx="12" cy="12" r="10"/>
          <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/>
          <line x1="12" y1="17" x2="12.01" y2="17"/>
        </svg>
      </button>
    </header>
  )
}
```

- [ ] **Step 2: Add `.guideBtn` style to `AppHeader.module.css`**

Append to the end of the file:

```css
.guideBtn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background: rgba(108, 150, 187, 0.12);
  color: var(--sky-light);
  border-radius: 50%;
  border: 1px solid rgba(108, 150, 187, 0.2);
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
  flex-shrink: 0;
}

.guideBtn:hover {
  background: rgba(108, 150, 187, 0.22);
  color: var(--white);
}
```

- [ ] **Step 3: Verify TypeScript compiles**

```bash
cd frontend && npx tsc --noEmit
```

Expected: only remaining error is `onOpenGuide` missing on `StandDetailPanel` — proceed.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/layout/AppHeader.tsx frontend/src/components/layout/AppHeader.module.css
git commit -m "feat: add guide button to AppHeader"
```

---

## Task 6: Add guide link to StandDetailPanel

**Files:**
- Modify: `frontend/src/components/panels/StandDetailPanel.tsx`
- Modify: `frontend/src/components/panels/StandDetailPanel.module.css`

- [ ] **Step 1: Add `onOpenGuide` prop and link to `StandDetailPanel.tsx`**

Add `onOpenGuide: () => void` to the `Props` interface, destructure it, and add the link at the bottom of the `.actions` div. The diff from current content:

```tsx
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
```

- [ ] **Step 2: Add `.guideLink` to `StandDetailPanel.module.css`**

Append to the end of the file:

```css
.guideLink {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
  background: none;
  border: none;
  color: var(--sky);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  padding: 8px;
  border-radius: var(--radius-sm);
  transition: color 0.15s, background 0.15s;
  width: 100%;
  text-align: center;
}

.guideLink:hover {
  color: var(--gold);
  background: rgba(90, 148, 200, 0.07);
}
```

- [ ] **Step 3: Verify TypeScript compiles clean**

```bash
cd frontend && npx tsc --noEmit
```

Expected: no errors.

- [ ] **Step 4: Build to confirm no bundle errors**

```bash
cd frontend && npm run build
```

Expected: build succeeds with output to `../static`.

- [ ] **Step 5: Commit**

```bash
git add frontend/src/components/panels/StandDetailPanel.tsx frontend/src/components/panels/StandDetailPanel.module.css
git commit -m "feat: add lock guide link to StandDetailPanel"
```

---

## Task 7: Manual verification

- [ ] **Step 1: Start dev server**

```bash
cd frontend && npm run dev
```

Open `http://localhost:5173` in a browser.

- [ ] **Step 2: Verify first-visit auto-open**

Clear `localStorage` in DevTools (`Application → Local Storage → Clear All`), reload. Guide should open automatically on the first render.

- [ ] **Step 3: Verify step navigation**

- Click Next through all 6 steps — steps should slide in from the right
- Click Back — steps should slide in from the left
- Progress bar should fill as steps advance
- "Skip" button should be visible on steps 1–5, hidden on step 6
- "Done" should appear on step 6 and close the guide

- [ ] **Step 4: Verify entry points**

- "?" button in the header should open the guide
- Clicking a stand and opening its detail panel should show "How to lock your bike securely →" at the bottom — clicking it opens the guide

- [ ] **Step 5: Verify keyboard accessibility**

- Open guide, press Escape → should close
- Tab should cycle through focusable elements within the modal only (focus trap)

- [ ] **Step 6: Verify no re-show on next visit**

Close the guide, reload the page. Guide should NOT auto-open again (`guide_seen` is set in localStorage).

- [ ] **Step 7: Final commit**

```bash
git add -A
git commit -m "feat: complete bike locking guide implementation"
```
