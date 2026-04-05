# Bike Locking Guide — Design Spec

**Date:** 2026-04-05
**Status:** Approved

---

## Overview

An animated, step-by-step guide teaching users how to securely lock their bike. Presented as a full-screen modal with slide-in step transitions. Shown automatically on first visit, and accessible at any time via the app header and from within a stand's detail panel.

---

## Entry Points

1. **First visit** — Auto-opens on mount if `localStorage.getItem('guide_seen')` is falsy. Sets the flag immediately on open so it doesn't re-appear.
2. **Header button** — A small circular "?" button added to the right of the existing "Find Nearest" button in `AppHeader`. Matches the existing button style.
3. **Stand detail panel** — A subtle text link at the bottom of `StandDetailPanel`: "How to lock your bike securely →"

---

## State Management

`App.tsx` owns `guideOpen: boolean` state.

- On mount: check `localStorage.getItem('guide_seen')`. If absent, set `guideOpen = true` and write `localStorage.setItem('guide_seen', 'true')`.
- `onOpenGuide` sets `guideOpen = true`.
- `onCloseGuide` sets `guideOpen = false`.
- Props flow: `App` → `AppShell` → `LockingGuide` (rendered inside AppShell), and `onOpenGuide` is passed to `AppHeader` and `StandDetailPanel`.

---

## Step Content

6 steps, each with a title, short body text, and an inline SVG illustration.

| # | Title | Body |
|---|-------|------|
| 1 | Check the stand | Look for signs of tampering — loose bolts, damaged metalwork, anything that looks interfered with. If it looks dodgy, find another stand. |
| 2 | Choose the right lock | Use a Sold Secure rated U-lock or D-lock as your primary. Avoid cable locks as your only lock — they can be cut in seconds. |
| 3 | Lock your frame | Pass the lock through your frame (not just the wheel) and around the stand. A wheel-only lock leaves your frame behind. |
| 4 | Add a secondary lock | Use a cable or chain to also secure your rear wheel. Two locks = double the deterrent. |
| 5 | Fill the lock | Leave as little space as possible inside the lock. A tight fit makes it much harder to lever open. |
| 6 | Pick your spot | Lock in a visible, well-lit location. Thieves prefer cover — busy, open areas are safer. |

Illustrations: simple inline SVGs using the app's navy/sky palette, consistent with existing SVG icon style.

---

## Components

### New files

- `frontend/src/components/guide/LockingGuide.tsx` — modal component
- `frontend/src/components/guide/LockingGuide.module.css` — styles
- `frontend/src/components/guide/steps.ts` — step content array (title, body, illustration component)

### Modified files

| File | Change |
|------|--------|
| `App.tsx` | Add `guideOpen` state + localStorage check on mount; pass `onOpenGuide`/`onCloseGuide` down |
| `AppShell.tsx` | Accept and forward guide props; render `<LockingGuide>` |
| `AppHeader.tsx` | Add "?" circular button, calls `onOpenGuide` |
| `StandDetailPanel.tsx` | Add "How to lock your bike securely →" text link at bottom |

---

## Modal UI

**Layout (top to bottom):**
- Thin sky-blue progress bar (fills proportionally as steps advance)
- "Skip" button — top right, muted; hidden on last step
- Illustration area — ~200px tall inline SVG
- Step counter — "Step 2 of 6", muted small text
- Title — Barlow Condensed, large, bold
- Body text — DM Sans, regular
- Prev / Next navigation buttons at bottom; Next becomes "Done" on final step

**Visual tokens:** uses existing CSS variables — `--navy`, `--sky`, `--navy-deep`, `--radius-lg`, `--shadow-lg`, `--font-display`, `--font-body`.

**Backdrop:** `rgba(9, 21, 56, 0.85)` full-screen overlay, fades in on open.

**Card:** centred, `max-width: 480px`, `--radius-lg` corners, `--shadow-lg`. On mobile: `width: 90vw`, fills more of the screen.

---

## Animations

No external library. Pure CSS transitions:

- **Backdrop:** `opacity` fade on open/close.
- **Step transition:** outgoing step slides left (`translateX(-100%)`), incoming step slides in from right (`translateX(100%) → 0`). Implemented by toggling a CSS class on the step wrapper when the active index changes. Direction reverses for Prev.
- Transition duration: ~250ms, `ease-in-out`.

---

## Accessibility

- Modal traps focus while open.
- `aria-modal="true"`, `role="dialog"`, `aria-label="How to lock your bike"`.
- Skip and close buttons are keyboard accessible.
- Progress bar has `aria-valuenow` / `aria-valuemax`.
