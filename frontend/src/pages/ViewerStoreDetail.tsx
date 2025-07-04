import React, { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import axios from 'axios'
import toast from 'react-hot-toast'
import Map from '../components/Map'
import ImageGallery from '../components/ImageGallery'
import StoreHeader from '../components/StoreHeader'
import StoreInfo from '../components/StoreInfo'
import ReviewList from '../components/ReviewList'
import { API_BASE_URL } from '@/config'
import './StoreDetail.css'

interface Store {
  id: string
  name: string
  address: string
  latitude: number
  longitude: number
  categories: string[]
  business_hours: string
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

interface Review {
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
  user: {
    id: string
    username: string
  }
}

const ViewerStoreDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const [store, setStore] = useState<Store | null>(null)
  const [reviews, setReviews] = useState<Review[]>([])
  const [loading, setLoading] = useState(true)
  const [reviewsLoading, setReviewsLoading] = useState(true)

  useEffect(() => {
    if (id) {
      fetchStore()
      fetchReviews()
    }
  }, [id])

  const fetchStore = async () => {
    try {
      setLoading(true)
      const response = await axios.get(`${API_BASE_URL}/api/v1/stores/${id}`)
      const responseData = response.data
      
      // レスポンス構造を確認
      let store = null
      if (responseData.success && responseData.data) {
        store = responseData.data
      } else if (responseData.data) {
        store = responseData.data
      } else {
        store = responseData
      }
      
      setStore(store)
    } catch (error) {
      console.error('Error fetching store:', error)
      toast.error('店舗情報の取得に失敗しました')
    } finally {
      setLoading(false)
    }
  }

  const fetchReviews = async () => {
    try {
      setReviewsLoading(true)
      const response = await axios.get(`${API_BASE_URL}/api/v1/stores/${id}/reviews`)
      const responseData = response.data
      
      // レスポンス構造を確認
      let reviews = []
      if (responseData.success && responseData.data && responseData.data.reviews) {
        reviews = responseData.data.reviews
      } else if (Array.isArray(responseData.data)) {
        reviews = responseData.data
      } else if (Array.isArray(responseData)) {
        reviews = responseData
      }
      
      setReviews(reviews)
    } catch (error) {
      console.error('Error fetching reviews:', error)
      toast.error('レビュー情報の取得に失敗しました')
    } finally {
      setReviewsLoading(false)
    }
  }

  if (loading) {
    return <div className="loading">読み込み中...</div>
  }

  if (!store) {
    return (
      <div className="error">
        <h2>店舗が見つかりません</h2>
        <Link to="/viewer/stores" className="btn btn-primary">
          店舗一覧に戻る
        </Link>
      </div>
    )
  }

  return (
    <div className="store-detail-page">
      <div className="store-detail-header">
        <Link to="/viewer/stores" className="back-link">
          ← 店舗一覧に戻る
        </Link>
      </div>

      <StoreHeader store={store} />
      
      <div className="store-detail-content">
        <div className="store-detail-main">
          <StoreInfo store={store} />
          
          {store.photos && store.photos.length > 0 && (
            <div className="store-photos-section">
              <h3>写真</h3>
              <ImageGallery images={store.photos} />
            </div>
          )}
          
          <div className="reviews-section">
            <h3>レビュー ({reviews.length}件)</h3>
            {reviewsLoading ? (
              <div className="loading">レビューを読み込み中...</div>
            ) : (
              <ReviewList 
                reviews={reviews} 
                showActions={false}
                readOnly={true}
              />
            )}
          </div>
        </div>
        
        <div className="store-detail-sidebar">
          <div className="map-section">
            <h3>地図</h3>
            <Map
              center={{ lat: store.latitude, lng: store.longitude }}
              markers={[{
                id: store.id,
                position: { lat: store.latitude, lng: store.longitude },
                title: store.name
              }]}
              zoom={16}
            />
          </div>
        </div>
      </div>
    </div>
  )
}

export default ViewerStoreDetail