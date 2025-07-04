import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import axios from 'axios'
import toast from 'react-hot-toast'
import StoreFilter from '../components/StoreFilter'
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
  businessDay: string
  businessTime: string
}

const StoreList: React.FC = () => {
  const [stores, setStores] = useState<Store[]>([])
  const [loading, setLoading] = useState(true)
  // 現在時刻から30分以上後の最短時間を計算
  const getDefaultBusinessDateTime = () => {
    const now = new Date()
    const currentMinutes = now.getMinutes()
    const currentHour = now.getHours()
    
    // 30分単位に丸める + 30分追加
    let targetMinutes = currentMinutes < 30 ? 30 : 0
    let targetHour = currentMinutes < 30 ? currentHour : currentHour + 1
    
    // さらに30分追加
    if (targetMinutes === 30) {
      targetMinutes = 0
      targetHour += 1
    } else {
      targetMinutes = 30
    }
    
    // 24時間を超えた場合は翌日
    if (targetHour >= 24) {
      targetHour = targetHour % 24
      now.setDate(now.getDate() + 1)
    }
    
    const dayNames = ['sunday', 'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday']
    const targetDay = dayNames[now.getDay()]
    const targetTime = `${targetHour.toString().padStart(2, '0')}:${targetMinutes.toString().padStart(2, '0')}`
    
    return { day: targetDay, time: targetTime }
  }

  const defaultDateTime = getDefaultBusinessDateTime()

  const [filters, setFilters] = useState<FilterState>({
    name: '',
    categories: [],
    categoriesOperator: 'OR', // デフォルトはOR
    tags: [],
    tagsOperator: 'AND', // デフォルトはAND
    businessDay: '', // デフォルトは指定なし
    businessTime: '', // デフォルトは指定なし
  })

  const fetchStores = async (filterParams?: FilterState) => {
    try {
      setLoading(true)
      const params = new URLSearchParams()
      const currentFilters = filterParams || filters
      
      console.log('Fetching stores with filters:', currentFilters)
      
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

      const url = `/api/v1/stores?${params.toString()}`
      console.log('Request URL:', url)
      
      const response = await axios.get(url)
      console.log('Response:', response.data)
      
      // Handle the new API response format
      if (response.data.success && response.data.data && response.data.data.stores) {
        setStores(response.data.data.stores)
      } else if (response.data.stores) {
        // Fallback for old format
        setStores(response.data.stores)
      } else {
        setStores([])
      }
    } catch (error) {
      toast.error('店舗の取得に失敗しました')
      console.error('Error fetching stores:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchStores()
  }, [])

  const handleFilterChange = (newFilters: FilterState) => {
    setFilters(newFilters)
    // 自動検索を無効化 - ユーザーが検索ボタンを押すまで検索しない
  }

  const handleSearch = () => {
    fetchStores(filters)
  }

  const handleReset = () => {
    const resetFilters = {
      name: '',
      categories: [],
      categoriesOperator: 'OR',
      tags: [],
      tagsOperator: 'AND',
      businessDay: '', // デフォルトは指定なし
      businessTime: '', // デフォルトは指定なし
    }
    setFilters(resetFilters)
    fetchStores(resetFilters) // リセット後は自動的に検索実行
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
    </div>
  )
}

export default StoreList