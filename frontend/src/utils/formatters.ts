/**
 * 日付をフォーマットする
 */
export const formatDate = (dateString: string | null): string | null => {
  if (!dateString) return null
  return new Date(dateString).toLocaleDateString('ja-JP')
}

/**
 * 星の評価を描画する
 */
export const renderStars = (rating: number): string => {
  return '★'.repeat(rating) + '☆'.repeat(5 - rating)
}

/**
 * 金額をフォーマットする
 */
export const formatCurrency = (amount: number | null): string => {
  if (!amount) return ''
  return `¥${amount.toLocaleString()}`
}

/**
 * 画像URLを生成する
 */
export const getImageUrl = (filename: string): string => {
  return `/api/v1/uploads/${filename}`
}

/**
 * 文字列を安全にトリムする
 */
export const safeTrim = (str: string | null | undefined): string => {
  return str?.trim() || ''
}

/**
 * 配列を安全に結合する
 */
export const safeJoin = (arr: string[] | null | undefined, separator: string = ', '): string => {
  return arr?.join(separator) || ''
}