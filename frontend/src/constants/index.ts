// Application constants for the Sukimise frontend

// API Configuration
export const API_CONFIG = {
  BASE_URL: import.meta.env.VITE_API_BASE_URL || '',
  TIMEOUT: 30000, // 30 seconds
  RETRY_ATTEMPTS: 3,
} as const

// User Roles
export const USER_ROLES = {
  ADMIN: 'admin',
  EDITOR: 'editor',
  VIEWER: 'viewer',
} as const

// Business Days
export const BUSINESS_DAYS = {
  MONDAY: 'monday',
  TUESDAY: 'tuesday',
  WEDNESDAY: 'wednesday',
  THURSDAY: 'thursday',
  FRIDAY: 'friday',
  SATURDAY: 'saturday',
  SUNDAY: 'sunday',
} as const

export const BUSINESS_DAY_LABELS = {
  [BUSINESS_DAYS.MONDAY]: '月曜日',
  [BUSINESS_DAYS.TUESDAY]: '火曜日',
  [BUSINESS_DAYS.WEDNESDAY]: '水曜日',
  [BUSINESS_DAYS.THURSDAY]: '木曜日',
  [BUSINESS_DAYS.FRIDAY]: '金曜日',
  [BUSINESS_DAYS.SATURDAY]: '土曜日',
  [BUSINESS_DAYS.SUNDAY]: '日曜日',
} as const

// Pagination
export const PAGINATION = {
  DEFAULT_PAGE: 1,
  DEFAULT_LIMIT: 20,
  MAX_LIMIT: 100,
  PAGE_SIZE_OPTIONS: [10, 20, 50, 100],
} as const

// File Upload
export const FILE_UPLOAD = {
  MAX_FILE_SIZE: 10 * 1024 * 1024, // 10MB
  ALLOWED_IMAGE_TYPES: [
    'image/jpeg',
    'image/png',
    'image/gif',
    'image/webp',
  ],
  MAX_IMAGES_PER_STORE: 20,
  MAX_IMAGES_PER_REVIEW: 10,
} as const

// Map Configuration
export const MAP_CONFIG = {
  DEFAULT_CENTER: {
    latitude: 35.6762,
    longitude: 139.6503, // Tokyo Station
  },
  DEFAULT_ZOOM: 13,
  MIN_ZOOM: 3,
  MAX_ZOOM: 18,
  TILE_URL: 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
  ATTRIBUTION: '© OpenStreetMap contributors',
} as const

// Form Validation
export const VALIDATION = {
  STORE_NAME: {
    MIN_LENGTH: 2,
    MAX_LENGTH: 100,
  },
  ADDRESS: {
    MIN_LENGTH: 5,
    MAX_LENGTH: 200,
  },
  COMMENT: {
    MAX_LENGTH: 1000,
  },
  TAGS: {
    MAX_COUNT: 20,
  },
  CATEGORIES: {
    MAX_COUNT: 10,
  },
  SNS_URLS: {
    MAX_COUNT: 5,
  },
  RATING: {
    MIN: 1,
    MAX: 5,
  },
  PAYMENT_AMOUNT: {
    MIN: 0,
    MAX: 1000000, // 1 million yen
  },
} as const

// Filter Operators
export const FILTER_OPERATORS = {
  AND: 'AND',
  OR: 'OR',
} as const

// Toast Configuration
export const TOAST_CONFIG = {
  DURATION: 4000, // 4 seconds
  POSITION: 'top-right',
} as const

// LocalStorage Keys
export const STORAGE_KEYS = {
  AUTH_TOKEN: 'authToken',
  USER_INFO: 'userInfo',
  USER_PREFERENCES: 'userPreferences',
  RECENT_SEARCHES: 'recentSearches',
  FAVORITE_STORES: 'favoriteStores',
  FILTER_STATE: 'filterState',
} as const

// Route Paths
export const ROUTES = {
  HOME: '/',
  LOGIN: '/login',
  STORES: '/stores',
  STORE_NEW: '/stores/new',
  STORE_DETAIL: '/stores/:id',
  STORE_EDIT: '/stores/:id/edit',
  MAP: '/map',
  PROFILE: '/profile',
  ADMIN: '/admin',
} as const

// Error Messages
export const ERROR_MESSAGES = {
  NETWORK_ERROR: 'ネットワークエラーが発生しました',
  UNAUTHORIZED: '認証が必要です',
  FORBIDDEN: 'アクセス権限がありません',
  NOT_FOUND: 'リソースが見つかりません',
  VALIDATION_ERROR: '入力内容に不備があります',
  SERVER_ERROR: 'サーバーエラーが発生しました',
  UNKNOWN_ERROR: '予期しないエラーが発生しました',
} as const

// Success Messages
export const SUCCESS_MESSAGES = {
  LOGIN_SUCCESS: 'ログインしました',
  LOGOUT_SUCCESS: 'ログアウトしました',
  STORE_CREATED: '店舗を登録しました',
  STORE_UPDATED: '店舗情報を更新しました',
  STORE_DELETED: '店舗を削除しました',
  REVIEW_CREATED: 'レビューを投稿しました',
  REVIEW_UPDATED: 'レビューを更新しました',
  REVIEW_DELETED: 'レビューを削除しました',
  IMAGE_UPLOADED: '画像をアップロードしました',
} as const

// Time Slots for Business Hours (30-minute intervals)
export const TIME_SLOTS = [
  '00:00', '00:30', '01:00', '01:30', '02:00', '02:30',
  '03:00', '03:30', '04:00', '04:30', '05:00', '05:30',
  '06:00', '06:30', '07:00', '07:30', '08:00', '08:30',
  '09:00', '09:30', '10:00', '10:30', '11:00', '11:30',
  '12:00', '12:30', '13:00', '13:30', '14:00', '14:30',
  '15:00', '15:30', '16:00', '16:30', '17:00', '17:30',
  '18:00', '18:30', '19:00', '19:30', '20:00', '20:30',
  '21:00', '21:30', '22:00', '22:30', '23:00', '23:30',
] as const

// HTTP Status Codes
export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  NO_CONTENT: 204,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  CONFLICT: 409,
  UNPROCESSABLE_ENTITY: 422,
  INTERNAL_SERVER_ERROR: 500,
  BAD_GATEWAY: 502,
  SERVICE_UNAVAILABLE: 503,
} as const

// Debounce Delays
export const DEBOUNCE_DELAYS = {
  SEARCH: 300, // milliseconds
  FILTER: 500,
  RESIZE: 250,
  SCROLL: 100,
} as const

// Animation Durations
export const ANIMATION_DURATIONS = {
  FAST: 150,
  NORMAL: 300,
  SLOW: 500,
} as const

// Breakpoints for responsive design
export const BREAKPOINTS = {
  MOBILE: 480,
  TABLET: 768,
  DESKTOP: 1024,
  LARGE_DESKTOP: 1200,
} as const

// Z-Index Values
export const Z_INDEX = {
  DROPDOWN: 1000,
  STICKY: 1020,
  FIXED: 1030,
  MODAL_BACKDROP: 1040,
  MODAL: 1050,
  POPOVER: 1060,
  TOOLTIP: 1070,
  TOAST: 1080,
} as const

// Theme Colors (if implementing theme switching)
export const THEME_COLORS = {
  PRIMARY: '#2563eb',
  SECONDARY: '#64748b',
  SUCCESS: '#16a34a',
  WARNING: '#d97706',
  DANGER: '#dc2626',
  INFO: '#0891b2',
  LIGHT: '#f8fafc',
  DARK: '#1e293b',
} as const