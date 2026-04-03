import { useState, useEffect, useCallback } from 'react'
import { useDropzone } from 'react-dropzone'
import { createStand } from '../../api/actions'
import { checkImageAvailability } from '../../api/stands'
import styles from './AddStandForm.module.css'

const STAND_TYPES = [
  'Sheffield Stand',
  'Hoop',
  'Stainless Steel Curved',
  'Railing',
  'Wheel Only',
  'Pride',
  'Other',
]

interface Props {
  lat: number
  lng: number
  onSuccess: () => void
  onCancel: () => void
}

export function AddStandForm({ lat, lng, onSuccess, onCancel }: Props) {
  const [type, setType] = useState(STAND_TYPES[0])
  const [name, setName] = useState('')
  const [count, setCount] = useState('')
  const [imageId, setImageId] = useState('')
  const [imageAvailable, setImageAvailable] = useState(false)
  const [uploadState, setUploadState] = useState<'idle' | 'uploading' | 'done' | 'error'>('idle')
  const [fileName, setFileName] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    checkImageAvailability().then(setImageAvailable).catch(() => setImageAvailable(false))
  }, [])

  const onDrop = useCallback(async (accepted: File[]) => {
    const file = accepted[0]
    if (!file) return
    setUploadState('uploading')
    setFileName(file.name)
    try {
      const body = new FormData()
      body.append('photo', file)
      const res = await fetch('/api/v0/image', { method: 'POST', body })
      if (!res.ok) throw new Error()
      const id = await res.text()
      setImageId(id)
      setUploadState('done')
    } catch {
      setUploadState('error')
    }
  }, [])

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: { 'image/jpeg': [], 'image/png': [], 'image/gif': [] },
    maxFiles: 1,
    maxSize: 5 * 1024 * 1024,
  })

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setSubmitting(true)
    setError('')
    try {
      await createStand({
        Lat: lat,
        Lng: lng,
        Type: type,
        Name: name || undefined,
        NumberOfStands: count ? parseInt(count, 10) : undefined,
        ImageID: imageId || undefined,
      })
      onSuccess()
    } catch {
      setError('Failed to submit. Please try again.')
      setSubmitting(false)
    }
  }

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      <div className={styles.header}>
        <h3 className={styles.title}>Add Bike Stand</h3>
        <button type="button" className={styles.closeBtn} onClick={onCancel} aria-label="Cancel">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
            <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>

      <div className={styles.coords}>
        {lat.toFixed(5)}, {lng.toFixed(5)}
      </div>

      <div className={styles.field}>
        <label className={styles.label}>Stand type *</label>
        <select className={styles.select} value={type} onChange={e => setType(e.target.value)} required>
          {STAND_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
        </select>
      </div>

      <div className={styles.field}>
        <label className={styles.label}>Location name</label>
        <input
          className={styles.input}
          type="text"
          value={name}
          onChange={e => setName(e.target.value)}
          placeholder="e.g. Outside Tesco on Grafton St"
          maxLength={200}
        />
      </div>

      <div className={styles.field}>
        <label className={styles.label}>Number of spaces</label>
        <input
          className={styles.input}
          type="number"
          value={count}
          onChange={e => setCount(e.target.value)}
          placeholder="e.g. 8"
          min="1"
          max="999"
        />
      </div>

      {imageAvailable && (
        <div className={styles.field}>
          <label className={styles.label}>Photo (optional)</label>
          {uploadState === 'done' ? (
            <div className={styles.uploadDone}>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
                <polyline points="20 6 9 17 4 12"/>
              </svg>
              {fileName}
              <button type="button" className={styles.uploadRemove} onClick={() => { setImageId(''); setUploadState('idle'); setFileName('') }}>Remove</button>
            </div>
          ) : (
            <div {...getRootProps()} className={`${styles.dropzone} ${isDragActive ? styles.dropzoneActive : ''} ${uploadState === 'uploading' ? styles.dropzoneUploading : ''} ${uploadState === 'error' ? styles.dropzoneError : ''}`}>
              <input {...getInputProps()} />
              {uploadState === 'uploading' ? (
                <span>Uploading…</span>
              ) : uploadState === 'error' ? (
                <span>Upload failed — try again</span>
              ) : isDragActive ? (
                <span>Drop photo here</span>
              ) : (
                <span>Drag a photo here, or <u>browse</u></span>
              )}
            </div>
          )}
        </div>
      )}

      {error && <div className={styles.error}>{error}</div>}

      <button type="submit" className={styles.submitBtn} disabled={submitting || uploadState === 'uploading'}>
        {submitting ? 'Submitting…' : 'Submit Stand'}
      </button>
    </form>
  )
}
