// Formatting utilities for consistent display across the application

// Date formatting
export const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString('ja-JP', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

export const formatDateTime = (dateString: string): string => {
  return new Date(dateString).toLocaleString('ja-JP', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

export const formatDateTimeShort = (dateString: string): string => {
  return new Date(dateString).toLocaleString('ja-JP', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Currency formatting
export const formatCurrency = (amount: number): string => {
  return amount.toLocaleString('ja-JP', {
    style: 'currency',
    currency: 'JPY',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  })
}

export const formatPrice = (amount: number): string => {
  return `${amount.toLocaleString()}円`
}

// Text formatting
export const truncateText = (text: string, maxLength: number): string => {
  if (text.length <= maxLength) {
    return text
  }
  return text.substring(0, maxLength) + '...'
}

export const capitalizeFirst = (text: string): string => {
  if (!text) return text
  return text.charAt(0).toUpperCase() + text.slice(1).toLowerCase()
}

// Business hours formatting
export const formatBusinessHours = (hours: string): string[] => {
  if (!hours) return []
  return hours.split('\n').filter(line => line.trim() !== '')
}

// Coordinates formatting
export const formatCoordinates = (lat: number, lng: number): string => {
  return `${lat.toFixed(6)}, ${lng.toFixed(6)}`
}

export const formatLatitude = (lat: number): string => {
  return `緯度: ${lat.toFixed(6)}`
}

export const formatLongitude = (lng: number): string => {
  return `経度: ${lng.toFixed(6)}`
}

// Rating formatting
export const formatRating = (rating: number): string => {
  return `${rating.toFixed(1)}`
}

export const getRatingStars = (rating: number): string => {
  const fullStars = Math.floor(rating)
  const hasHalfStar = rating % 1 >= 0.5
  const emptyStars = 5 - fullStars - (hasHalfStar ? 1 : 0)
  
  return '★'.repeat(fullStars) + 
         (hasHalfStar ? '☆' : '') + 
         '☆'.repeat(emptyStars)
}

// Array formatting
export const formatList = (items: string[], maxItems: number = 3): string => {
  if (items.length === 0) return ''
  
  if (items.length <= maxItems) {
    return items.join(', ')
  }
  
  const visibleItems = items.slice(0, maxItems)
  const remainingCount = items.length - maxItems
  return `${visibleItems.join(', ')} 他${remainingCount}件`
}

// URL formatting
export const ensureHttps = (url: string): string => {
  if (!url) return url
  if (url.startsWith('http://') || url.startsWith('https://')) {
    return url
  }
  return `https://${url}`
}

export const formatWebsiteUrl = (url: string): string => {
  if (!url) return ''
  const formatted = ensureHttps(url)
  // Remove protocol for display
  return formatted.replace(/^https?:\/\//, '')
}

// File size formatting
export const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes'
  
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// Distance formatting
export const formatDistance = (distance: number): string => {
  if (distance < 1000) {
    return `${Math.round(distance)}m`
  } else {
    return `${(distance / 1000).toFixed(1)}km`
  }
}

// Time formatting
export const formatTime = (time: string): string => {
  // Ensure time is in HH:MM format
  const [hours, minutes] = time.split(':')
  return `${hours.padStart(2, '0')}:${minutes.padStart(2, '0')}`
}

export const formatBusinessDay = (day: string): string => {
  const dayNames: Record<string, string> = {
    monday: '月',
    tuesday: '火',
    wednesday: '水',
    thursday: '木',
    friday: '金',
    saturday: '土',
    sunday: '日'
  }
  
  return dayNames[day] || day
}

// Search query formatting
export const formatSearchQuery = (query: string): string => {
  return query.trim().toLowerCase()
}

// Error message formatting
export const formatErrorMessage = (error: any): string => {
  if (typeof error === 'string') {
    return error
  }
  
  if (error?.message) {
    return error.message
  }
  
  if (error?.response?.data?.error) {
    return error.response.data.error
  }
  
  if (error?.response?.data?.message) {
    return error.response.data.message
  }
  
  return 'エラーが発生しました'
}

// Input value formatting
export const formatInputValue = (value: any): string => {
  if (value === null || value === undefined) {
    return ''
  }
  return String(value)
}

// Phone number formatting (Japanese format)
export const formatPhoneNumber = (phone: string): string => {
  // Remove all non-digit characters
  const digits = phone.replace(/\D/g, '')
  
  // Format based on length
  if (digits.length === 10) {
    return digits.replace(/(\d{2})(\d{4})(\d{4})/, '$1-$2-$3')
  } else if (digits.length === 11) {
    return digits.replace(/(\d{3})(\d{4})(\d{4})/, '$1-$2-$3')
  }
  
  return phone
}

// Tag formatting
export const formatTag = (tag: string): string => {
  return tag.startsWith('#') ? tag : `#${tag}`
}

export const formatTags = (tags: string[]): string[] => {
  return tags.map(formatTag)
}