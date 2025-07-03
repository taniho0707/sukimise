import axios from 'axios'
import { Store, Review } from '@/types/store'

/**
 * 店舗を取得する
 */
export const fetchStore = async (id: string): Promise<Store> => {
  const response = await axios.get(`/api/v1/stores/${id}`)
  return response.data.data
}

/**
 * 店舗一覧を取得する
 */
export const fetchStores = async (params?: Record<string, any>): Promise<Store[]> => {
  const response = await axios.get('/api/v1/stores', { params })
  return response.data.data.stores
}

/**
 * 店舗を作成する
 */
export const createStore = async (storeData: Partial<Store>): Promise<Store> => {
  const response = await axios.post('/api/v1/stores', storeData)
  return response.data.data
}

/**
 * 店舗を更新する
 */
export const updateStore = async (id: string, storeData: Partial<Store>): Promise<Store> => {
  const response = await axios.put(`/api/v1/stores/${id}`, storeData)
  return response.data.data
}

/**
 * 店舗を削除する
 */
export const deleteStore = async (id: string): Promise<void> => {
  await axios.delete(`/api/v1/stores/${id}`)
}

/**
 * 店舗のレビューを取得する
 */
export const fetchStoreReviews = async (storeId: string): Promise<Review[]> => {
  const response = await axios.get(`/api/v1/stores/${storeId}/reviews`)
  
  // APIレスポンス形式の処理
  // Backend returns: { "reviews": [...] }
  if (response.data.reviews && Array.isArray(response.data.reviews)) {
    return response.data.reviews
  } else if (response.data.success && response.data.data) {
    return Array.isArray(response.data.data) ? response.data.data : []
  } else if (Array.isArray(response.data)) {
    return response.data
  } else if (response.data.data && Array.isArray(response.data.data)) {
    return response.data.data
  }
  
  // デフォルトで空配列を返す
  return []
}

/**
 * レビューを作成する
 */
export const createReview = async (reviewData: Partial<Review>): Promise<Review> => {
  const response = await axios.post('/api/v1/reviews', reviewData)
  return response.data.data
}

/**
 * レビューを更新する
 */
export const updateReview = async (id: string, reviewData: Partial<Review>): Promise<Review> => {
  const response = await axios.put(`/api/v1/reviews/${id}`, reviewData)
  return response.data.data
}

/**
 * レビューを削除する
 */
export const deleteReview = async (id: string): Promise<void> => {
  await axios.delete(`/api/v1/reviews/${id}`)
}

/**
 * カテゴリ一覧を取得する
 */
export const fetchCategories = async (): Promise<string[]> => {
  const response = await axios.get('/api/v1/stores/categories')
  return response.data.data.categories
}

/**
 * タグ一覧を取得する
 */
export const fetchTags = async (): Promise<string[]> => {
  const response = await axios.get('/api/v1/stores/tags')
  return response.data.data.tags
}

/**
 * 画像をアップロードする
 */
export const uploadImage = async (file: File): Promise<string> => {
  const formData = new FormData()
  formData.append('image', file)
  
  const response = await axios.post('/api/v1/upload/image', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
  
  return response.data.data.filename
}

/**
 * 画像を削除する
 */
export const deleteImage = async (filename: string): Promise<void> => {
  await axios.delete(`/api/v1/upload/${filename}`)
}