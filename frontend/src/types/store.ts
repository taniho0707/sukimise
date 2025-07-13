export interface Store {
  id: string
  name: string
  address: string
  latitude: number
  longitude: number
  categories: string[]
  business_hours: BusinessHoursData
  parking_info: string
  website_url: string
  google_map_url: string
  sns_urls: string[]
  tags: string[]
  photos: string[]
  created_by: string
  created_at: string
  updated_at: string
}

export interface Review {
  id: string
  store_id: string
  user_id: string
  rating: number
  comment: string | null
  photos: string[]
  visit_date: string | null
  is_visited: boolean
  payment_amount: number | null
  food_notes: string | null
  created_at: string
  updated_at: string
  user?: {
    id: string
    username: string
  }
}

export interface User {
  id: string
  username: string
  email: string
  role: 'admin' | 'editor' | 'viewer'
  created_at: string
  updated_at: string
}

export interface FilterState {
  name: string
  categories: string[]
  categoriesOperator: 'AND' | 'OR'
  tags: string[]
  tagsOperator: 'AND' | 'OR'
  businessDay: string
  businessTime: string
}

export interface APIResponse<T = any> {
  success: boolean
  data?: T
  error?: {
    code: string
    message: string
    details?: string
  }
  meta?: {
    total?: number
    page?: number
    limit?: number
    offset?: number
  }
}

export interface LocationCoordinates {
  latitude: number
  longitude: number
}

export interface TimeSlot {
  open_time: string
  close_time: string
  last_order_time: string
}

export interface DaySchedule {
  is_closed: boolean
  time_slots: TimeSlot[]
}

export interface BusinessHoursData {
  monday: DaySchedule
  tuesday: DaySchedule
  wednesday: DaySchedule
  thursday: DaySchedule
  friday: DaySchedule
  saturday: DaySchedule
  sunday: DaySchedule
}

export interface BusinessHoursInput {
  day: string
  time: string
}

export interface ValidationError {
  field: string
  message: string
}

export interface StoreFormData {
  name: string
  address: string
  latitude: number
  longitude: number
  categories: string[]
  business_hours: BusinessHoursData
  parking_info: string
  website_url: string
  google_map_url: string
  sns_urls: string[]
  tags: string[]
  photos: string[]
}

export interface ReviewFormData {
  store_id: string
  rating: number
  comment: string
  photos: string[]
  visit_date: Date | null
  is_visited: boolean
  payment_amount: number | undefined
  food_notes: string
}

export interface AuthContextType {
  user: User | null
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  loading: boolean
  hasRole?: (role: string) => boolean
  isAdmin?: () => boolean
  isEditor?: () => boolean
  canEdit?: (createdBy: string) => boolean
}

export const USER_ROLES = {
  ADMIN: 'admin' as const,
  EDITOR: 'editor' as const,
  VIEWER: 'viewer' as const,
} as const

export const DEFAULT_FILTER_STATE: FilterState = {
  name: '',
  categories: [],
  categoriesOperator: 'OR',
  tags: [],
  tagsOperator: 'AND',
  businessDay: '',
  businessTime: '',
} as const