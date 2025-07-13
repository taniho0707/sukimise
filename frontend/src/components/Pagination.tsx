import React, { useState } from 'react'

interface PaginationProps {
  currentPage: number
  totalPages: number
  onPageChange: (page: number) => void
  className?: string
}

const Pagination: React.FC<PaginationProps> = ({
  currentPage,
  totalPages,
  onPageChange,
  className = ''
}) => {
  const [jumpToPage, setJumpToPage] = useState('')

  // 表示するページ番号の範囲を計算
  const getVisiblePages = () => {
    const visiblePages: number[] = []
    const delta = 2 // 現在のページから前後2つまで表示

    for (let i = Math.max(1, currentPage - delta); i <= Math.min(totalPages, currentPage + delta); i++) {
      visiblePages.push(i)
    }

    return visiblePages
  }

  const handleJumpToPage = (e: React.FormEvent) => {
    e.preventDefault()
    const pageNumber = parseInt(jumpToPage, 10)
    if (pageNumber >= 1 && pageNumber <= totalPages) {
      onPageChange(pageNumber)
      setJumpToPage('')
    }
  }

  if (totalPages <= 1) {
    return null
  }

  const visiblePages = getVisiblePages()

  return (
    <div className={`pagination-container ${className}`}>
      <div className="pagination-info">
        <span>{totalPages}ページ中 {currentPage}ページ目</span>
      </div>
      
      <div className="pagination-controls">
        {/* 前のページボタン */}
        <button
          type="button"
          onClick={() => onPageChange(currentPage - 1)}
          disabled={currentPage <= 1}
          className="pagination-button"
          aria-label="前のページ"
        >
          ←
        </button>

        {/* 最初のページ */}
        {visiblePages[0] > 1 && (
          <>
            <button
              type="button"
              onClick={() => onPageChange(1)}
              className="pagination-button"
            >
              1
            </button>
            {visiblePages[0] > 2 && (
              <span className="pagination-ellipsis">...</span>
            )}
          </>
        )}

        {/* 表示範囲のページ番号 */}
        {visiblePages.map((page) => (
          <button
            key={page}
            type="button"
            onClick={() => onPageChange(page)}
            className={`pagination-button ${page === currentPage ? 'active' : ''}`}
          >
            {page}
          </button>
        ))}

        {/* 最後のページ */}
        {visiblePages[visiblePages.length - 1] < totalPages && (
          <>
            {visiblePages[visiblePages.length - 1] < totalPages - 1 && (
              <span className="pagination-ellipsis">...</span>
            )}
            <button
              type="button"
              onClick={() => onPageChange(totalPages)}
              className="pagination-button"
            >
              {totalPages}
            </button>
          </>
        )}

        {/* 次のページボタン */}
        <button
          type="button"
          onClick={() => onPageChange(currentPage + 1)}
          disabled={currentPage >= totalPages}
          className="pagination-button"
          aria-label="次のページ"
        >
          →
        </button>
      </div>

      {/* ページジャンプ機能 */}
      <div className="page-jump">
        <form onSubmit={handleJumpToPage} className="page-jump-form">
          <input
            type="number"
            min="1"
            max={totalPages}
            value={jumpToPage}
            onChange={(e) => setJumpToPage(e.target.value)}
            placeholder="ページ番号"
            className="page-jump-input"
          />
          <button type="submit" className="page-jump-button">
            移動
          </button>
        </form>
      </div>
    </div>
  )
}

export default Pagination