import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import axios from 'axios'
import toast from 'react-hot-toast'
import { useAuth } from '@/contexts/AuthContext'
import ReviewForm from '../components/ReviewForm'
import StoreHeader from '../components/StoreHeader'
import StoreInfo from '../components/StoreInfo'
import ReviewList from '../components/ReviewList'
import { Store, Review } from '../types/store'
import { fetchStore, fetchStoreReviews, deleteStore, deleteReview } from '../utils/api'
import './StoreDetail.css'

const StoreDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { user } = useAuth()
  
  const [store, setStore] = useState<Store | null>(null)
  const [reviews, setReviews] = useState<Review[]>([])
  const [loading, setLoading] = useState(true)
  const [reviewsLoading, setReviewsLoading] = useState(false)
  const [showReviewForm, setShowReviewForm] = useState(false)
  const [editingReview, setEditingReview] = useState<Review | null>(null)
  const [userLatestRating, setUserLatestRating] = useState<number | undefined>()
  useEffect(() => {
    if (!id) {
      navigate('/stores')
      return
    }

    const loadStoreData = async () => {
      try {
        setLoading(true)
        const storeData = await fetchStore(id)
        setStore(storeData)
      } catch (error) {
        console.error('Store fetch error:', error)
        toast.error('店舗情報の取得に失敗しました')
        navigate('/stores')
      } finally {
        setLoading(false)
      }
    }

    const loadReviews = async () => {
      try {
        setReviewsLoading(true)
        const reviewsData = await fetchStoreReviews(id)
        setReviews(reviewsData)
        
        // 現在のユーザーのレビューを特定
        if (user) {
          // const userReviews = reviewsData.filter((review: Review) => review.user_id === user.id)
          // setCurrentUserReviews(userReviews) - removed variable
        }
      } catch (error) {
        console.error('Reviews fetch error:', error)
        // レビュー取得エラーは致命的ではないため、エラートーストは表示しない
      } finally {
        setReviewsLoading(false)
      }
    }

    loadStoreData()
    loadReviews()
  }, [id, navigate, user])

  // ユーザーの最新レビュー評価を取得する関数
  useEffect(() => {
    const fetchUserLatestRating = async () => {
      if (!user) return
      
      try {
        // ユーザーの全レビューを取得
        const response = await axios.get('/api/v1/users/me/reviews')
        
        // APIレスポンス形式の処理
        let userReviews = []
        if (response.data.success && response.data.data && response.data.data.reviews) {
          userReviews = response.data.data.reviews
        } else if (response.data.reviews) {
          userReviews = response.data.reviews
        } else if (Array.isArray(response.data)) {
          userReviews = response.data
        }
        
        if (userReviews.length > 0) {
          // 最新のレビューの評価を取得
          const latestReview = userReviews.sort((a: Review, b: Review) => 
            new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
          )[0]
          setUserLatestRating(latestReview.rating)
        }
      } catch (error) {
        console.error('User reviews fetch error:', error)
        setUserLatestRating(undefined)
      }
    }

    if (user) {
      fetchUserLatestRating()
    }
  }, [user])

  const handleEdit = () => {
    if (store) {
      navigate(`/stores/${store.id}/edit`)
    }
  }

  const handleDelete = async () => {
    if (!store || !window.confirm('この店舗を削除しますか？この操作は取り消せません。')) {
      return
    }

    try {
      await deleteStore(store.id)
      toast.success('店舗を削除しました')
      navigate('/stores')
    } catch (error: any) {
      console.error('Delete error:', error)
      toast.error('店舗の削除に失敗しました')
    }
  }

  const handleReviewSubmit = () => {
    // レビューが追加/更新された後にリロード
    const reloadReviews = async () => {
      if (id) {
        try {
          const reviewsData = await fetchStoreReviews(id)
          setReviews(reviewsData)
          
          if (user) {
            // const userReviews = reviewsData.filter((review: Review) => review.user_id === user.id)
            // setCurrentUserReviews(userReviews) - removed variable
          }
        } catch (error) {
          console.error('Reviews reload error:', error)
          // レビューリロード時のエラーは表示しない（レビュー自体の投稿は成功しているため）
        }
      }
    }

    reloadReviews()
    setShowReviewForm(false)
    setEditingReview(null)
  }

  const handleEditReview = (review: Review) => {
    setEditingReview(review)
    setShowReviewForm(true)
  }

  const handleDeleteReview = async (reviewId: string) => {
    if (!window.confirm('このレビューを削除しますか？')) {
      return
    }

    try {
      await deleteReview(reviewId)
      toast.success('レビューを削除しました')
      
      // レビューリストから削除
      setReviews(prev => prev.filter(review => review.id !== reviewId))
      // setCurrentUserReviews(prev => prev.filter(review => review.id !== reviewId)) - removed variable
    } catch (error: any) {
      console.error('Review delete error:', error)
      toast.error('レビューの削除に失敗しました')
    }
  }

  const canEdit = () => {
    return user && store && (store.created_by === user.id || user.role === 'admin')
  }

  if (loading) {
    return <div className="loading">読み込み中...</div>
  }

  if (!store) {
    return <div className="error">店舗が見つかりません</div>
  }

  return (
    <div className="store-detail-page">
      <div className="store-detail-container">
        {/* 店舗ヘッダー */}
        <StoreHeader
          store={store}
          canEdit={canEdit() || false}
          onEdit={handleEdit}
          onDelete={handleDelete}
        />

        {/* 店舗情報 */}
        <StoreInfo store={store} />

        {/* レビューセクション */}
        <div className="reviews-section">
          <div className="reviews-header">
            <h2>レビュー</h2>
            {user && (
              <button
                onClick={() => {
                  setEditingReview(null)
                  setShowReviewForm(true)
                }}
                className="btn btn-primary"
              >
                レビューを投稿
              </button>
            )}
          </div>

          {/* レビューフォーム */}
          {showReviewForm && (
            <div className="review-form-section">
              <ReviewForm
                storeId={store.id}
                existingReview={editingReview ? {
                  ...editingReview,
                  comment: editingReview.comment || '',
                  visit_date: editingReview.visit_date || undefined,
                  payment_amount: editingReview.payment_amount || undefined,
                  food_notes: editingReview.food_notes || undefined
                } : undefined}
                onSuccess={handleReviewSubmit}
                onCancel={() => {
                  setShowReviewForm(false)
                  setEditingReview(null)
                }}
                userLatestRating={userLatestRating}
              />
            </div>
          )}

          {/* レビューリスト */}
          {reviewsLoading ? (
            <div className="loading">レビューを読み込み中...</div>
          ) : (
            <ReviewList
              reviews={reviews}
              currentUserId={user?.id}
              onEditReview={handleEditReview}
              onDeleteReview={handleDeleteReview}
            />
          )}
        </div>
      </div>
    </div>
  )
}

export default StoreDetail