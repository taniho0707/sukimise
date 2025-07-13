import React, { useState, useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import axios from 'axios'
import toast from 'react-hot-toast'
import Map from '../components/Map'
import SafeBusinessHoursDisplay from '../components/SafeBusinessHoursDisplay'
import './MapView.css'

interface Store {
  id: string
  name: string
  address: string
  latitude: number
  longitude: number
  categories: string[]
  business_hours?: string | object
  tags: string[]
}

const MapView: React.FC = () => {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [stores, setStores] = useState<Store[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedStore, setSelectedStore] = useState<Store | null>(null)
  const [showSidebar, setShowSidebar] = useState(false)
  const [mapCenter, setMapCenter] = useState<[number, number]>([35.6762, 139.6503]) // デフォルト: 東京駅
  const [mapZoom, setMapZoom] = useState(13)

  useEffect(() => {
    const fetchStores = async () => {
      try {
        setLoading(true)
        const response = await axios.get('/api/v1/stores')
        
        let storesData: Store[] = []
        // Handle the new API response format
        if (response.data.success && response.data.data && response.data.data.stores) {
          storesData = response.data.data.stores
        } else if (response.data.stores) {
          // Fallback for old format
          storesData = response.data.stores
        } else {
          storesData = []
        }
        
        setStores(storesData)
        
        // URLパラメータから店舗IDを取得して、その店舗を中心に表示
        const storeId = searchParams.get('store')
        if (storeId && storesData.length > 0) {
          const targetStore = storesData.find(store => store.id === storeId)
          if (targetStore) {
            setMapCenter([targetStore.latitude, targetStore.longitude])
            setMapZoom(16) // より詳細なズームレベル
            setSelectedStore(targetStore)
            setShowSidebar(true)
            console.log(`Centering map on store: ${targetStore.name} at ${targetStore.latitude}, ${targetStore.longitude}`)
          } else {
            console.warn(`Store with ID ${storeId} not found`)
          }
        }
      } catch (error) {
        console.error('Stores fetch error:', error)
        toast.error('店舗データの取得に失敗しました')
      } finally {
        setLoading(false)
      }
    }

    fetchStores()
  }, [searchParams])

  const handleStoreClick = (store: Store) => {
    setSelectedStore(store)
    setShowSidebar(true)
  }

  const handleStoreDetailClick = (storeId: string) => {
    navigate(`/stores/${storeId}`)
  }

  const handleCloseSidebar = () => {
    setShowSidebar(false)
    setSelectedStore(null)
  }

  if (loading) {
    return (
      <div className="map-view-loading">
        <div className="loading-spinner">地図を読み込み中...</div>
      </div>
    )
  }

  return (
    <div className="map-view">
      <div className="map-view-header">
        <h1>店舗マップ</h1>
        <div className="map-controls">
          <span className="store-count">
            {stores.length}件の店舗
          </span>
          <button
            onClick={() => navigate('/stores')}
            className="btn btn-secondary"
          >
            リスト表示
          </button>
        </div>
      </div>

      <div className="map-view-content">
        <div className="map-wrapper">
          <Map
            stores={stores}
            center={mapCenter}
            zoom={mapZoom}
            height="calc(100vh - 140px)"
            onStoreClick={handleStoreClick}
            selectedStore={selectedStore}
          />
        </div>

        {/* サイドバー */}
        {showSidebar && selectedStore && (
          <>
            <div 
              className="sidebar-overlay"
              onClick={handleCloseSidebar}
            />
            <div className="store-sidebar">
              <div className="sidebar-header">
                <h2>{selectedStore.name}</h2>
                <button
                  onClick={handleCloseSidebar}
                  className="close-btn"
                >
                  ×
                </button>
              </div>

              <div className="sidebar-content">
                <div className="store-info">
                  <div className="info-item">
                    <strong>住所</strong>
                    <p>{selectedStore.address}</p>
                  </div>

                  {selectedStore.categories.length > 0 && (
                    <div className="info-item">
                      <strong>カテゴリ</strong>
                      <div className="categories">
                        {selectedStore.categories.map((category, index) => (
                          <span key={index} className="category-tag">
                            {category}
                          </span>
                        ))}
                      </div>
                    </div>
                  )}


                  {selectedStore.business_hours && (
                    <div className="info-item">
                      <strong>営業時間</strong>
                      <div className="business-hours">
                        <SafeBusinessHoursDisplay businessHours={selectedStore.business_hours} />
                      </div>
                    </div>
                  )}

                  {selectedStore.tags.length > 0 && (
                    <div className="info-item">
                      <strong>タグ</strong>
                      <div className="tags">
                        {selectedStore.tags.map((tag, index) => (
                          <span key={index} className="tag">
                            #{tag}
                          </span>
                        ))}
                      </div>
                    </div>
                  )}

                  <div className="info-item">
                    <strong>位置情報</strong>
                    <p className="coordinates">
                      緯度: {selectedStore.latitude}<br />
                      経度: {selectedStore.longitude}
                    </p>
                  </div>
                </div>

                <div className="sidebar-actions">
                  <button
                    onClick={() => handleStoreDetailClick(selectedStore.id)}
                    className="btn btn-primary btn-full"
                  >
                    詳細を見る
                  </button>
                  <button
                    onClick={() => {
                      const url = `https://www.google.com/maps?q=${selectedStore.latitude},${selectedStore.longitude}`
                      window.open(url, '_blank')
                    }}
                    className="btn btn-secondary btn-full"
                  >
                    Googleマップで開く
                  </button>
                </div>
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  )
}

export default MapView