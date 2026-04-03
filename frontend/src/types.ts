export type StandType =
  | 'Sheffield Stand'
  | 'Hoop'
  | 'Hoops'
  | 'Stainless Steel Curved'
  | 'Railing'
  | 'Wheel Only'
  | 'Pride'
  | string

export interface StandProperties {
  id: string
  name: string
  type: StandType
  numberOfStands: number | null
  notes: string
  imageId: string
  source: string
  checked: boolean
  verified: boolean
  publicImageURL: string
}

export interface StandFeature {
  type: 'Feature'
  geometry: {
    type: 'Point'
    coordinates: [number, number] // [lng, lat]
  }
  properties: StandProperties
}

export interface StandCollection {
  type: 'FeatureCollection'
  features: StandFeature[]
}

export interface QueryParams {
  showIDs: boolean
  verified: boolean    // ?checked=unchecked shows unverified
  pendingReview: boolean
}
