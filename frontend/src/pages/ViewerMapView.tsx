import React, { useState, useEffect } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import axios from 'axios'
import toast from 'react-hot-toast'
import Map from '../components/Map'
import StoreFilter from '../components/StoreFilter'
import { API_BASE_URL } from '@/config'
import './MapView.css'

interface Store {
  id: string
  name: string
  address: string
  latitude: number
  longitude: number
  categories: string[]
  tags: string[]
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

interface Marker {
  id: string
  position: { lat: number; lng: number }
  title: string
  store: Store
}

const ViewerMapView: React.FC = () => {
  const [searchParams] = useSearchParams()
  const [stores, setStores] = useState<Store[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedStore, setSelectedStore] = useState<Store | null>(null)
  const [showSidebar, setShowSidebar] = useState(false)
  // Load saved map center and display limit from localStorage (viewer version)
  const loadSavedSettings = () => {
    const savedCenter = localStorage.getItem('viewerMapCenter')
    const savedZoom = localStorage.getItem('viewerMapZoom')
    const savedLimit = localStorage.getItem('viewerDisplayLimit')
    
    return {
      center: savedCenter ? JSON.parse(savedCenter) : [35.6762, 139.6503],
      zoom: savedZoom ? parseInt(savedZoom) : 13,
      limit: savedLimit ? parseInt(savedLimit) : 30
    }
  }

  const savedSettings = loadSavedSettings()
  const [mapCenter, setMapCenter] = useState<[number, number]>(savedSettings.center)
  const [mapZoom, setMapZoom] = useState(savedSettings.zoom)
  const [displayLimit, setDisplayLimit] = useState(savedSettings.limit)
  
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

  const fetchStoresByProximity = async () => {
    try {
      setLoading(true)
      
      // Save current settings to localStorage
      localStorage.setItem('viewerMapCenter', JSON.stringify(mapCenter))
      localStorage.setItem('viewerMapZoom', mapZoom.toString())
      localStorage.setItem('viewerDisplayLimit', displayLimit.toString())
      
      // 地図中心から近い順で店舗を取得
      const params = new URLSearchParams({
        latitude: mapCenter[0].toString(),
        longitude: mapCenter[1].toString(),
        order_by_proximity: 'true',
        limit: displayLimit.toString()
      })
      
      const response = await axios.get(`${API_BASE_URL}/api/v1/stores?${params.toString()}`)
      const responseData = response.data
      
      // レスポンス構造を確認
      let storesData = []
      if (responseData.success && responseData.data && responseData.data.stores) {
        storesData = responseData.data.stores
      } else if (Array.isArray(responseData.data)) {
        storesData = responseData.data
      } else if (Array.isArray(responseData)) {
        storesData = responseData
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
      console.error('Error fetching stores:', error)
      toast.error('店舗の取得に失敗しました')
    } finally {
      setLoading(false)
    }
  }

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
      let storesData = []
      if (responseData.success && responseData.data && responseData.data.stores) {
        storesData = responseData.data.stores
      } else if (Array.isArray(responseData.data)) {
        storesData = responseData.data
      } else if (Array.isArray(responseData)) {
        storesData = responseData
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
      console.error('Error fetching stores:', error)
      toast.error('店舗の取得に失敗しました')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchStores(filterState)
  }, [searchParams])

  const handleFilterChange = (newFilters: FilterState) => {
    setFilterState(newFilters)
    fetchStores(newFilters)
  }

  const markers: Marker[] = stores.map(store => ({
    id: store.id,
    position: { lat: store.latitude, lng: store.longitude },
    title: store.name,
    store
  }))

  const handleMarkerClick = (marker: Marker) => {
    setSelectedStore(marker.store)
    setShowSidebar(true)
  }

  const handleCloseSidebar = () => {
    setShowSidebar(false)
    setSelectedStore(null)
  }

  return (
    <div className="map-view-page">
      <div className="page-header">
        <h1>地図から探す</h1>
      </div>

      <div className="map-view-content">
        <div className="map-sidebar">
          <StoreFilter
            initialFilters={filterState}
            onFilterChange={handleFilterChange}
          />
          
          {showSidebar && selectedStore && (
            <div className="selected-store-info">
              <div className="selected-store-header">
                <h3>選択中の店舗</h3>
                <button
                  onClick={handleCloseSidebar}
                  className="close-btn"
                >
                  ×
                </button>
              </div>
              <div className="store-card">
                <h4>
                  <Link to={`/viewer/stores/${selectedStore.id}`}>
                    {selectedStore.name}
                  </Link>
                </h4>
                <p>{selectedStore.address}</p>
                <div className="store-categories">
                  {selectedStore.categories && selectedStore.categories.map((category, index) => (
                    <span key={index} className="category-tag">{category}</span>
                  ))}
                </div>
                <div className="store-tags">
                  {selectedStore.tags && selectedStore.tags.map((tag, index) => (
                    <span key={index} className="tag">{tag}</span>
                  ))}
                </div>
              </div>
            </div>
          )}
        </div>
        
        <div className="map-container">
          {loading ? (
            <div className="map-loading">地図を読み込み中...</div>
          ) : (
            <Map
              markers={markers}
              center={mapCenter}
              zoom={mapZoom}
              onMarkerClick={handleMarkerClick}
            />
          )}
        </div>
      </div>
    </div>
  )
}

export default ViewerMapView