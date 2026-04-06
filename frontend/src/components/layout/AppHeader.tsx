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
