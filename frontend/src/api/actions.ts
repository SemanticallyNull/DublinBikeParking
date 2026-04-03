export interface CreateStandPayload {
  Lat: number
  Lng: number
  Type: string
  Name?: string
  NumberOfStands?: number | null
  ImageID?: string
  Notes?: string
}

export async function createStand(payload: CreateStandPayload): Promise<void> {
  const res = await fetch('/api/v0/stand', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })
  if (!res.ok) throw new Error(`Failed to create stand: ${res.status}`)
}

export async function reportMissing(id: string): Promise<void> {
  const res = await fetch(`/api/v0/stand/${id}/missing`)
  if (!res.ok) throw new Error(`Failed to report missing: ${res.status}`)
}
