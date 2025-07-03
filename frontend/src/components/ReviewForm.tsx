import React, { useState, useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import axios from 'axios'
import toast from 'react-hot-toast'
import { useAuth } from '../contexts/AuthContext'
import ImageUpload from './ImageUpload'
import './ReviewForm.css'

const reviewSchema = z.object({
  rating: z.number().min(1, '評価は必須です').max(5, '評価は1〜5の範囲で入力してください'),
  comment: z.string().optional(),
  visitDate: z.string().optional(),
  isVisited: z.boolean(),
  paymentAmount: z.number().min(0, '支払金額は0円以上で入力してください').optional(),
  foodNotes: z.string().optional(),
})

type ReviewFormData = z.infer<typeof reviewSchema>

interface ReviewFormProps {
  storeId: string
  onSuccess: () => void
  onCancel: () => void
  existingReview?: {
    id: string
    rating: number
    comment: string
    visit_date?: string
    is_visited: boolean
    photos?: string[]
    payment_amount?: number
    food_notes?: string
  }
  userLatestRating?: number
}

const ReviewForm: React.FC<ReviewFormProps> = ({ 
  storeId, 
  onSuccess, 
  onCancel, 
  existingReview,
  userLatestRating 
}) => {
  const { user } = useAuth()
  const [loading, setLoading] = useState(false)
  const [isVisited, setIsVisited] = useState(true)
  const [images, setImages] = useState<string[]>([])
  const isEdit = !!existingReview

  // デフォルト評価値を決定
  const getDefaultRating = () => {
    if (existingReview) return existingReview.rating
    if (userLatestRating) return userLatestRating
    return 5
  }

  // デフォルト来店日を今日の日付に設定
  const getCurrentDate = () => {
    return new Date().toISOString().split('T')[0]
  }

  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
    setValue
  } = useForm<ReviewFormData>({
    resolver: zodResolver(reviewSchema),
    defaultValues: {
      rating: getDefaultRating(),
      comment: existingReview?.comment || '',
      visitDate: existingReview?.visit_date ? 
        new Date(existingReview.visit_date).toISOString().split('T')[0] : getCurrentDate(),
      isVisited: existingReview?.is_visited ?? true,
      paymentAmount: existingReview?.payment_amount || undefined,
      foodNotes: existingReview?.food_notes || '',
    }
  })

  // 初期化時に来店状態を設定
  useEffect(() => {
    if (existingReview) {
      setIsVisited(existingReview.is_visited)
      setImages(existingReview.photos || [])
    } else {
      setIsVisited(true)
      setValue('isVisited', true)
      setValue('visitDate', getCurrentDate())
    }
  }, [existingReview, setValue])

  const toggleVisited = () => {
    const newVisitedState = !isVisited
    setIsVisited(newVisitedState)
    setValue('isVisited', newVisitedState)
    
    if (newVisitedState) {
      // 来店状態にした場合、来店日をデフォルトに設定
      setValue('visitDate', getCurrentDate())
    }
  }

  const onSubmit = async (data: ReviewFormData) => {
    if (!user) {
      toast.error('ログインが必要です')
      return
    }

    setLoading(true)
    try {
      const reviewData = {
        store_id: storeId,
        rating: data.rating,
        comment: data.comment || '',
        visit_date: data.visitDate ? new Date(data.visitDate).toISOString() : null,
        is_visited: data.isVisited,
        photos: images,
        payment_amount: (data.paymentAmount !== undefined && data.paymentAmount !== null && !isNaN(data.paymentAmount)) ? data.paymentAmount : null,
        food_notes: data.foodNotes || '',
      }

      if (isEdit && existingReview) {
        await axios.put(`/api/v1/reviews/${existingReview.id}`, reviewData)
        toast.success('レビューを更新しました')
      } else {
        await axios.post('/api/v1/reviews', reviewData)
        toast.success('レビューを投稿しました')
      }
      
      onSuccess()
    } catch (error: any) {
      console.error('Review save error:', error)
      if (error.response?.data?.error) {
        toast.error(error.response.data.error)
      } else {
        toast.error(isEdit ? 'レビューの更新に失敗しました' : 'レビューの投稿に失敗しました')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="review-form-overlay">
      <div className="review-form-modal">
        <div className="review-form-header">
          <h3>{isEdit ? 'レビューを編集' : 'レビューを投稿'}</h3>
          <button 
            type="button" 
            onClick={onCancel}
            className="close-btn"
            disabled={loading}
          >
            ×
          </button>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="review-form">
          {/* 評価 */}
          <div className="form-group">
            <label className="form-label required">評価</label>
            <div className="rating-input">
              {[1, 2, 3, 4, 5].map((star) => (
                <button
                  key={star}
                  type="button"
                  className={`star-btn ${watch('rating') >= star ? 'active' : ''}`}
                  onClick={() => setValue('rating', star)}
                  disabled={loading}
                >
                  ★
                </button>
              ))}
              <span className="rating-text">
                {watch('rating')}/5
              </span>
            </div>
            {errors.rating && (
              <span className="error-message">{errors.rating.message}</span>
            )}
          </div>

          {/* 来店ボタン */}
          <div className="form-group">
            <label className="form-label">来店</label>
            <button
              type="button"
              onClick={toggleVisited}
              className={`visit-toggle-btn ${isVisited ? 'active' : ''}`}
              disabled={loading}
            >
              {isVisited ? '✓ 来店済み' : '未来店'}
            </button>
          </div>

          {/* 来店日 */}
          {isVisited && (
            <div className="form-group">
              <label htmlFor="visitDate" className="form-label">
                来店日
              </label>
              <input
                type="date"
                id="visitDate"
                {...register('visitDate')}
                className="form-input"
                disabled={loading}
              />
            </div>
          )}

          {/* コメント */}
          <div className="form-group">
            <label htmlFor="comment" className="form-label">
              コメント
            </label>
            <textarea
              id="comment"
              {...register('comment')}
              className="form-textarea"
              placeholder="料理の感想、サービス、雰囲気など..."
              rows={4}
              disabled={loading}
            />
          </div>

          {/* 支払金額 */}
          {isVisited && (
            <div className="form-group">
              <label htmlFor="paymentAmount" className="form-label">
                支払金額（1人分）
              </label>
              <input
                type="number"
                id="paymentAmount"
                {...register('paymentAmount', { 
                  setValueAs: (value) => value === '' || value === null ? undefined : Number(value)
                })}
                className="form-input"
                placeholder="例: 1500"
                min="0"
                step="1"
                disabled={loading}
              />
              <small className="form-help">
                1人で支払った金額を円単位で入力してください（任意）
              </small>
              {errors.paymentAmount && (
                <span className="error-message">{errors.paymentAmount.message}</span>
              )}
            </div>
          )}

          {/* 料理メモ */}
          {isVisited && (
            <div className="form-group">
              <label htmlFor="foodNotes" className="form-label">
                注文した料理
              </label>
              <textarea
                id="foodNotes"
                {...register('foodNotes')}
                className="form-textarea"
                placeholder="例: ハンバーグセット、ドリンクバー、デザート"
                rows={2}
                disabled={loading}
              />
              <small className="form-help">
                注文した料理についてのメモ（任意）
              </small>
            </div>
          )}

          {/* 画像 */}
          <div className="form-group">
            <ImageUpload
              images={images}
              onChange={setImages}
              maxImages={5}
              disabled={loading}
            />
          </div>

          {/* 送信ボタン */}
          <div className="form-actions">
            <button
              type="button"
              onClick={onCancel}
              className="btn btn-secondary"
              disabled={loading}
            >
              キャンセル
            </button>
            <button
              type="submit"
              className="btn btn-primary"
              disabled={loading}
            >
              {loading ? 
                (isEdit ? '更新中...' : '投稿中...') : 
                (isEdit ? '更新' : '投稿')
              }
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default ReviewForm