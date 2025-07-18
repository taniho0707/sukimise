// 画像URLを絶対URLに変換する関数
export const getImageUrl = (imageUrl: string): string => {
  // 既に完全なURLの場合はそのまま返す
  if (imageUrl.startsWith('http://') || imageUrl.startsWith('https://')) {
    return imageUrl
  }
  
  // 相対パスの場合はベースURLを追加
  const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080' // バックエンドのURL
  
  // /uploads/で始まる場合
  if (imageUrl.startsWith('/uploads/')) {
    const finalUrl = `${baseUrl}${imageUrl}`
    console.log('Image URL converted:', imageUrl, '->', finalUrl)
    return finalUrl
  }
  
  // uploads/で始まる場合
  if (imageUrl.startsWith('uploads/')) {
    const finalUrl = `${baseUrl}/${imageUrl}`
    console.log('Image URL converted:', imageUrl, '->', finalUrl)
    return finalUrl
  }
  
  // その他の相対パス
  if (imageUrl.startsWith('/')) {
    const finalUrl = `${baseUrl}${imageUrl}`
    console.log('Image URL converted:', imageUrl, '->', finalUrl)
    return finalUrl
  }
  
  const finalUrl = `${baseUrl}/uploads/${imageUrl}`
  console.log('Image URL converted:', imageUrl, '->', finalUrl)
  return finalUrl
}

// 複数の画像URLを変換する関数
export const getImageUrls = (imageUrls: string[]): string[] => {
  return imageUrls.map(getImageUrl)
}