import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import axios from 'axios'
import toast from 'react-hot-toast'
import StoreFilter from '../components/StoreFilter'
import { API_BASE_URL } from '@/config'
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

const ViewerStoreList: React.FC = () => {
  const [stores, setStores] = useState<Store[]>([])
  const [loading, setLoading] = useState(true)
  
  // 現在時刻から30分以上後の最短時間を計算
  const getCurrentPlus30MinTime = () => {
    const now = new Date()
    now.setMinutes(now.getMinutes() + 30)
    const hours = String(now.getHours()).padStart(2, '0')
    const minutes = Math.ceil(now.getMinutes() / 30) * 30
    const formattedMinutes = String(minutes === 60 ? 0 : minutes).padStart(2, '0')
    const formattedHours = minutes === 60 ? String(Number(hours) + 1).padStart(2, '0') : hours
    return `${formattedHours}:${formattedMinutes}`
  }

  const [filterState, setFilterState] = useState<FilterState>({
    name: '',
    categories: [],
    categoriesOperator: 'OR',
    tags: [],
    tagsOperator: 'OR',
    priceMin: 0,
    priceMax: 10000,
    businessDay: '',
    businessTime: getCurrentPlus30MinTime()
  })

  const fetchStores = async (filters: FilterState) => {
    try {
      setLoading(true)
      const params = new URLSearchParams()
      
      if (filters.name.trim()) {
        params.append('name', filters.name.trim())
      }
      
      if (filters.categories.length > 0) {
        params.append('categories', filters.categories.join(','))
        params.append('categories_operator', filters.categoriesOperator)
      }
      
      if (filters.tags.length > 0) {
        params.append('tags', filters.tags.join(','))
        params.append('tags_operator', filters.tagsOperator)
      }
      
      if (filters.businessDay && filters.businessTime) {
        params.append('business_day', filters.businessDay)
        params.append('business_time', filters.businessTime)
      }

      const response = await axios.get(`${API_BASE_URL}/api/v1/stores?${params.toString()}`)
      const responseData = response.data
      
      // レスポンス構造を確認
      let stores = []
      if (responseData.success && responseData.data && responseData.data.stores) {
        stores = responseData.data.stores
      } else if (Array.isArray(responseData.data)) {
        stores = responseData.data
      } else if (Array.isArray(responseData)) {
        stores = responseData
      }
      
      setStores(stores)
    } catch (error) {
      console.error('Error fetching stores:', error)
      toast.error('店舗の取得に失敗しました')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchStores(filterState)
  }, [])

  const handleFilterChange = (newFilters: FilterState) => {
    setFilterState(newFilters)
    fetchStores(newFilters)
  }

  const handleExportCSV = async () => {
    try {
      const params = new URLSearchParams()
      
      if (filterState.name.trim()) {
        params.append('name', filterState.name.trim())
      }
      
      if (filterState.categories.length > 0) {
        params.append('categories', filterState.categories.join(','))
        params.append('categories_operator', filterState.categoriesOperator)
      }
      
      if (filterState.tags.length > 0) {
        params.append('tags', filterState.tags.join(','))
        params.append('tags_operator', filterState.tagsOperator)
      }
      
      if (filterState.businessDay && filterState.businessTime) {
        params.append('business_day', filterState.businessDay)
        params.append('business_time', filterState.businessTime)
      }

      const response = await axios.get(`${API_BASE_URL}/api/v1/stores/export/csv?${params.toString()}`, {
        responseType: 'blob'
      })

      const blob = new Blob([response.data], { type: 'text/csv' })
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      
      const timestamp = new Date().toISOString().slice(0, 19).replace(/:/g, '-')
      link.download = `stores_${timestamp}.csv`
      
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
      
      toast.success('CSVファイルをダウンロードしました')
    } catch (error) {
      console.error('Error exporting CSV:', error)
      toast.error('CSVエクスポートに失敗しました')
    }
  }

  return (
    <div className="store-list-page">
      <div className="page-header">
        <h1>店舗一覧</h1>
        <div className="page-actions">
          <button onClick={handleExportCSV} className="btn btn-secondary">
            CSVエクスポート
          </button>
        </div>
      </div>

      <StoreFilter
        initialFilters={filterState}
        onFilterChange={handleFilterChange}
      />

      {loading ? (
        <div className="loading">読み込み中...</div>
      ) : (
        <div className="store-grid">
          {stores.length === 0 ? (
            <div className="no-stores">
              <p>条件に一致する店舗が見つかりませんでした。</p>
            </div>
          ) : (
            stores.map((store) => (
              <div key={store.id} className="store-card">
                <h3 className="store-name">
                  <Link to={`/viewer/stores/${store.id}`}>
                    {store.name}
                  </Link>
                </h3>
                <p className="store-address">{store.address}</p>
                <div className="store-categories">
                  {store.categories && store.categories.map((category, index) => (
                    <span key={index} className="category-tag">{category}</span>
                  ))}
                </div>
                <div className="store-tags">
                  {store.tags && store.tags.map((tag, index) => (
                    <span key={index} className="tag">{tag}</span>
                  ))}
                </div>
              </div>
            ))
          )}
        </div>
      )}
    </div>
  )
}

export default ViewerStoreList