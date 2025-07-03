/**
 * メールアドレスの形式をチェックする
 */
export const isValidEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

/**
 * URLの形式をチェックする
 */
export const isValidUrl = (url: string): boolean => {
  if (!url) return true // 空の場合は有効とする
  try {
    new URL(url)
    return true
  } catch {
    return false
  }
}

/**
 * 緯度の有効性をチェックする
 */
export const isValidLatitude = (lat: number): boolean => {
  return lat >= -90 && lat <= 90
}

/**
 * 経度の有効性をチェックする
 */
export const isValidLongitude = (lng: number): boolean => {
  return lng >= -180 && lng <= 180
}

/**
 * 評価の有効性をチェックする
 */
export const isValidRating = (rating: number): boolean => {
  return rating >= 1 && rating <= 5 && Number.isInteger(rating)
}

/**
 * ファイルサイズをチェックする
 */
export const isValidFileSize = (file: File, maxSizeMB: number = 5): boolean => {
  const maxSizeBytes = maxSizeMB * 1024 * 1024
  return file.size <= maxSizeBytes
}

/**
 * 画像ファイルの形式をチェックする
 */
export const isValidImageFile = (file: File): boolean => {
  const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp']
  return allowedTypes.includes(file.type)
}

/**
 * 文字列の長さをチェックする
 */
export const isValidLength = (str: string, minLength: number = 0, maxLength: number = Infinity): boolean => {
  return str.length >= minLength && str.length <= maxLength
}

/**
 * 必須フィールドのチェック
 */
export const isRequired = (value: any): boolean => {
  if (typeof value === 'string') {
    return value.trim().length > 0
  }
  return value !== null && value !== undefined && value !== ''
}