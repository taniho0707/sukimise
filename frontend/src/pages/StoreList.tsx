import React, { useState, useEffect } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import axios from 'axios'
import toast from 'react-hot-toast'
import StoreFilter from '../components/StoreFilter'
import Pagination from '../components/Pagination'
import './StoreList.css'

interface Store {
  id: string
  name: string
  address: string
  categories: string[]
  tags: string[]
  created_at: string
}

interface FilterState {
  name: string
  categories: string[]
  categoriesOperator: string
  tags: string[]
  tagsOperator: string
  priceMin: number
  priceMax: number
  businessDay: string
  businessTime: string
}

interface PaginationMeta {
  total: number
  page: number
  limit: number
  offset: number
  total_pages: number
}

const StoreList: React.FC = () => {
  const [stores, setStores] = useState<Store[]>([])
  const [loading, setLoading] = useState(true)
  const [pagination, setPagination] = useState<PaginationMeta>({
    total: 0,
    page: 1,
    limit: 20,
    offset: 0,
    total_pages: 1
  })
  const [searchParams, setSearchParams] = useSearchParams()
  // 現在時刻から30分以上後の最短時間を計算

  const [filters, setFilters] = useState<FilterState>({
    name: '',
    categories: [],
    categoriesOperator: 'OR', // デフォルトはOR
    tags: [],
    tagsOperator: 'AND', // デフォルトはAND
    priceMin: 0,
    priceMax: 10000,
    businessDay: '', // デフォルトは指定なし
    businessTime: '', // デフォルトは指定なし
  })

  const fetchStores = async (filterParams?: FilterState, page?: number) => {
    try {
      setLoading(true)
      const params = new URLSearchParams()
      const currentFilters = filterParams || filters
      const currentPage = page || pagination.page
      
      console.log('Fetching stores with filters:', currentFilters, 'page:', currentPage)
      
      if (currentFilters.name) params.append('name', currentFilters.name)
      if (currentFilters.categories.length > 0) {
        params.append('categories', currentFilters.categories.join(','))
        params.append('categories_operator', currentFilters.categoriesOperator)
      }
      if (currentFilters.tags.length > 0) {
        params.append('tags', currentFilters.tags.join(','))
        params.append('tags_operator', currentFilters.tagsOperator)
      }
      if (currentFilters.businessDay) params.append('business_day', currentFilters.businessDay)
      if (currentFilters.businessTime) params.append('business_time', currentFilters.businessTime)

      // Add pagination parameters
      params.append('limit', '20')
      params.append('offset', ((currentPage - 1) * 20).toString())

      const url = `/api/v1/stores?${params.toString()}`
      console.log('Request URL:', url)
      
      const response = await axios.get(url)
      console.log('Response:', response.data)
      
      // Handle the new API response format
      if (response.data.success && response.data.data && response.data.data.stores) {
        setStores(response.data.data.stores)
        
        // Update pagination metadata
        if (response.data.meta) {
          setPagination({
            total: response.data.meta.total || 0,
            page: response.data.meta.page || currentPage,
            limit: response.data.meta.limit || 20,
            offset: response.data.meta.offset || 0,
            total_pages: response.data.meta.total_pages || 1
          })
        }
      } else if (response.data.stores) {
        // Fallback for old format
        setStores(response.data.stores)
        setPagination({
          total: response.data.stores.length,
          page: currentPage,
          limit: 20,
          offset: 0,
          total_pages: 1
        })
      } else {
        setStores([])
        setPagination({
          total: 0,
          page: currentPage,
          limit: 20,
          offset: 0,
          total_pages: 1
        })
      }

      // Update URL parameters
      const urlParams = new URLSearchParams(searchParams)
      urlParams.set('page', currentPage.toString())
      setSearchParams(urlParams)
    } catch (error) {
      toast.error('店舗の取得に失敗しました')
      console.error('Error fetching stores:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    // URLパラメータからページ番号を取得
    const pageParam = searchParams.get('page')
    const initialPage = pageParam ? parseInt(pageParam, 10) : 1
    
    setPagination(prev => ({ ...prev, page: initialPage }))
    fetchStores(undefined, initialPage)
  }, [])

  const handleFilterChange = (newFilters: FilterState) => {
    setFilters(newFilters)
    // 自動検索を無効化 - ユーザーが検索ボタンを押すまで検索しない
  }

  const handleSearch = () => {
    // 検索時は1ページ目に戻る
    setPagination(prev => ({ ...prev, page: 1 }))
    fetchStores(filters, 1)
  }

  const handleReset = () => {
    const resetFilters = {
      name: '',
      categories: [],
      categoriesOperator: 'OR',
      tags: [],
      tagsOperator: 'AND',
      priceMin: 0,
      priceMax: 10000,
      businessDay: '', // デフォルトは指定なし
      businessTime: '', // デフォルトは指定なし
    }
    setFilters(resetFilters)
    setPagination(prev => ({ ...prev, page: 1 }))
    fetchStores(resetFilters, 1) // リセット後は自動的に検索実行
  }

  const handlePageChange = (page: number) => {
    setPagination(prev => ({ ...prev, page }))
    fetchStores(filters, page)
  }

  const handleExportCSV = async () => {
    try {
      const params = new URLSearchParams()
      
      if (filters.name) params.append('name', filters.name)
      if (filters.categories.length > 0) {
        params.append('categories', filters.categories.join(','))
        params.append('categories_operator', filters.categoriesOperator)
      }
      if (filters.tags.length > 0) {
        params.append('tags', filters.tags.join(','))
        params.append('tags_operator', filters.tagsOperator)
      }
      if (filters.businessDay) params.append('business_day', filters.businessDay)
      if (filters.businessTime) params.append('business_time', filters.businessTime)

      const url = `/api/v1/stores/export/csv?${params.toString()}`
      
      const response = await axios.get(url, {
        responseType: 'blob', // Important for downloading files
      })

      // Create blob link to download
      const blob = new Blob([response.data], { type: 'text/csv' })
      const link = document.createElement('a')
      link.href = window.URL.createObjectURL(blob)
      
      // Get filename from response headers or use default
      const contentDisposition = response.headers['content-disposition']
      let filename = 'sukimise_stores.csv'
      if (contentDisposition) {
        const filenameMatch = contentDisposition.match(/filename=(.+)/)
        if (filenameMatch) {
          filename = filenameMatch[1]
        }
      }
      
      link.download = filename
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      
      toast.success('CSVファイルのダウンロードを開始しました')
    } catch (error) {
      toast.error('CSVエクスポートに失敗しました')
      console.error('Error exporting CSV:', error)
    }
  }

  if (loading) {
    return <div className="loading">店舗情報を読み込み中...</div>
  }

  return (
    <div className="store-list-page">
      <div className="page-header">
        <h1>店舗一覧</h1>
        <div className="header-actions">
          <button onClick={handleExportCSV} className="btn btn-outline">
            CSVエクスポート
          </button>
          <Link to="/map" className="btn btn-secondary">
            地図表示
          </Link>
          <Link to="/stores/new" className="btn btn-primary">
            新しい店舗を登録
          </Link>
        </div>
      </div>

      <StoreFilter 
        initialFilters={filters}
        onFilterChange={handleFilterChange}
        onSearch={handleSearch}
        onReset={handleReset}
      />

      <div className="stores-grid">
        {stores.length === 0 ? (
          <div className="no-stores">
            <p>店舗が見つかりませんでした。</p>
            <Link to="/stores/new" className="btn btn-primary">
              最初の店舗を登録する
            </Link>
          </div>
        ) : (
          stores.map((store) => (
            <div key={store.id} className="store-card">
              <div className="store-card-header">
                <h3 className="store-name">
                  <Link to={`/stores/${store.id}`}>{store.name}</Link>
                </h3>
              </div>
              
              <p className="store-address">{store.address}</p>
              
              {store.categories && store.categories.length > 0 && (
                <div className="store-categories">
                  {store.categories.map((category, index) => (
                    <span key={index} className="category-tag">
                      {category}
                    </span>
                  ))}
                </div>
              )}
              
              {store.tags && store.tags.length > 0 && (
                <div className="store-tags">
                  {store.tags.map((tag, index) => (
                    <span key={index} className="tag">
                      #{tag}
                    </span>
                  ))}
                </div>
              )}
              
              <div className="store-card-footer">
                <span className="created-date">
                  {new Date(store.created_at).toLocaleDateString('ja-JP')}
                </span>
                <Link to={`/stores/${store.id}`} className="view-link">
                  詳細を見る →
                </Link>
              </div>
            </div>
          ))
        )}
      </div>

      {/* Pagination */}
      {stores.length > 0 && (
        <Pagination
          currentPage={pagination.page}
          totalPages={pagination.total_pages}
          onPageChange={handlePageChange}
          className="store-list-pagination"
        />
      )}
    </div>
  )
}

export default StoreList