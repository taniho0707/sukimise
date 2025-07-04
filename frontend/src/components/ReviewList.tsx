import React, { useState, useMemo } from 'react'

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
  user?: {
    id: string
    username: string
  }
}

interface ReviewListProps {
  reviews?: Review[]
  currentUserId?: string
  onEditReview?: (review: Review) => void
  onDeleteReview?: (reviewId: string) => void
  showActions?: boolean
  readOnly?: boolean
}

const ReviewList: React.FC<ReviewListProps> = ({ 
  reviews = [], 
  currentUserId, 
  onEditReview, 
  onDeleteReview,
  showActions = true,
  readOnly = false
}) => {
  const renderStars = (rating: number) => {
    return '★'.repeat(rating) + '☆'.repeat(5 - rating)
  }

  const formatDate = (dateString: string | null) => {
    if (!dateString) return null
    return new Date(dateString).toLocaleDateString('ja-JP')
  }

  const renderTextWithLineBreaks = (text: string) => {
    return text.split('\n').map((line, index) => (
      <React.Fragment key={index}>
        {line}
        {index < text.split('\n').length - 1 && <br />}
      </React.Fragment>
    ))
  }

  const canEditReview = (review: Review) => {
    return !readOnly && showActions && currentUserId && review.user_id === currentUserId
  }

  // ページネーション関連の状態
  const [currentPage, setCurrentPage] = useState(1)
  const [pageInput, setPageInput] = useState('')
  const reviewsPerPage = 10

  // ページネーション計算
  const totalPages = Math.ceil((reviews?.length || 0) / reviewsPerPage)
  const paginatedReviews = useMemo(() => {
    if (!reviews) return []
    const startIndex = (currentPage - 1) * reviewsPerPage
    const endIndex = startIndex + reviewsPerPage
    return reviews.slice(startIndex, endIndex)
  }, [reviews, currentPage])

  const handlePageChange = (page: number) => {
    if (page >= 1 && page <= totalPages) {
      setCurrentPage(page)
    }
  }

  const handlePageInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPageInput(e.target.value)
  }

  const handlePageJump = () => {
    const page = parseInt(pageInput)
    if (!isNaN(page) && page >= 1 && page <= totalPages) {
      setCurrentPage(page)
      setPageInput('')
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      handlePageJump()
    }
  }

  return (
    <>
      {!reviews || reviews.length === 0 ? (
        <div className="no-reviews">
          <p>まだレビューがありません。</p>
        </div>
      ) : (
        <>
          <div className="reviews-list">
            {paginatedReviews.map((review) => (
            <div key={review.id} className="review-card">
              <div className="review-header">
                <div className="review-meta">
                  <div className="reviewer">
                    {review.user?.username || 'Unknown User'}
                  </div>
                  {review.visit_date && (
                    <div className="visit-date">
                      訪問日: {formatDate(review.visit_date)}
                    </div>
                  )}
                  {!review.is_visited && (
                    <span className="not-visited">未来店</span>
                  )}
                </div>
                
                <div className="review-header-actions">
                  <div className="star-rating">
                    {Array.from({length: 5}, (_, i) => (
                      <span key={i} className={`star ${i < review.rating ? 'filled' : ''}`}>
                        ★
                      </span>
                    ))}
                    <span className="rating-text">({review.rating}/5)</span>
                  </div>
                  {canEditReview(review) && (
                    <div className="review-actions-inline">
                      <button 
                        onClick={() => onEditReview?.(review)}
                        className="btn-icon"
                        title="編集"
                      >
                        ✏️
                      </button>
                      <button 
                        onClick={() => onDeleteReview?.(review.id)}
                        className="btn-icon"
                        title="削除"
                      >
                        🗑️
                      </button>
                    </div>
                  )}
                </div>
              </div>

              <div className="review-content">
                {review.comment && (
                  <p>{renderTextWithLineBreaks(review.comment)}</p>
                )}

                {(review.payment_amount || review.food_notes) && (
                  <div className="review-payment-info">
                    {review.payment_amount && (
                      <div className="payment-amount">
                        <span className="label">支払金額:</span>
                        <span className="amount">¥{review.payment_amount.toLocaleString()}</span>
                      </div>
                    )}
                    {review.food_notes && (
                      <div className="food-notes">
                        <span className="label">料理メモ:</span>
                        <span className="notes">{renderTextWithLineBreaks(review.food_notes)}</span>
                      </div>
                    )}
                  </div>
                )}

                {review.photos && review.photos.length > 0 && (
                  <div className="review-images">
                    <div className="review-image-gallery">
                      {review.photos.map((photo, index) => (
                        <img 
                          key={index} 
                          src={`/api/v1/uploads/${photo}`} 
                          alt={`Review photo ${index + 1}`}
                          style={{ maxWidth: '200px', height: 'auto', margin: '4px' }}
                        />
                      ))}
                    </div>
                  </div>
                )}
              </div>
              
              <div className="review-footer">
                <span>投稿日: {formatDate(review.created_at)}</span>
              </div>
            </div>
          ))}
          </div>
          
          {/* ページネーション */}
          {totalPages > 1 && (
            <div className="pagination">
              <div className="pagination-controls">
                <button 
                  onClick={() => handlePageChange(currentPage - 1)}
                  disabled={currentPage === 1}
                  className="pagination-btn"
                >
                  前へ
                </button>
                
                <div className="pagination-info">
                  <span>{currentPage} / {totalPages} ページ</span>
                  <span className="reviews-count">（全 {reviews?.length || 0} 件）</span>
                </div>
                
                <button 
                  onClick={() => handlePageChange(currentPage + 1)}
                  disabled={currentPage === totalPages}
                  className="pagination-btn"
                >
                  次へ
                </button>
              </div>
              
              <div className="page-jumper">
                <input 
                  type="number" 
                  value={pageInput}
                  onChange={handlePageInputChange}
                  onKeyPress={handleKeyPress}
                  placeholder="ページ番号"
                  min="1"
                  max={totalPages}
                  className="page-input"
                />
                <button 
                  onClick={handlePageJump}
                  className="jump-btn"
                >
                  移動
                </button>
              </div>
            </div>
          )}
        </>
      )}
    </>
  )
}

export default ReviewList