import styles from './FloatingButton.module.css'

interface Props {
  loading?: boolean
  onClick: () => void
}

export function FloatingButton({ loading, onClick }: Props) {
  return (
    <button
      className={styles.fab}
      onClick={onClick}
      disabled={loading}
      aria-label="Find nearest bike stand"
      title="Find nearest bike stand"
    >
      {loading ? (
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" className={styles.spin}>
          <path d="M12 2a10 10 0 0 1 10 10"/>
        </svg>
      ) : (
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
          <circle cx="12" cy="12" r="9"/><circle cx="12" cy="12" r="3"/>
          <line x1="12" y1="2" x2="12" y2="5"/><line x1="12" y1="19" x2="12" y2="22"/>
          <line x1="2" y1="12" x2="5" y2="12"/><line x1="19" y1="12" x2="22" y2="12"/>
        </svg>
      )}
    </button>
  )
}
