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
