import React, { useEffect, useRef, useState, useId } from 'react'
import L from 'leaflet'
import axios from 'axios'
import 'leaflet/dist/leaflet.css'
import './Map.css'
import { API_BASE_URL } from '@/config'

// Leafletのデフォルトアイコンの問題を修正
delete (L.Icon.Default.prototype as any)._getIconUrl
L.Icon.Default.mergeOptions({
  iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png',
  iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png',
  shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png',
})

interface Store {
  id: string
  name: string
  address: string
  latitude: number
  longitude: number
  categories: string[]
  tags: string[]
}

interface CategoryCustomization {
  id: string
  category_name: string
  icon?: string
  color?: string
  created_at: string
  updated_at: string
}

interface Marker {
  id: string
  position: { lat: number; lng: number }
  title: string
}

interface MapProps {
  stores?: Store[]
  markers?: Marker[]
  center?: [number, number] | { lat: number; lng: number }
  zoom?: number
  height?: string
  onStoreClick?: (store: Store) => void
  onMarkerClick?: (marker: Marker) => void
  selectedStore?: Store | null
  onMapCenterChange?: (center: [number, number]) => void
  onMapZoomChange?: (zoom: number) => void
}

const Map: React.FC<MapProps> = ({
  stores = [],
  markers = [],
  center = [35.6762, 139.6503], // 東京駅
  zoom = 13,
  height = '400px',
  onStoreClick,
  onMarkerClick,
  selectedStore,
  onMapCenterChange,
  onMapZoomChange
}) => {
  const mapId = useId()
  const mapRef = useRef<L.Map | null>(null)
  const markersRef = useRef<L.Marker[]>([])
  const mapContainerRef = useRef<HTMLDivElement>(null)
  const [categoryCustomizations, setCategoryCustomizations] = useState<CategoryCustomization[]>([])

  // Fetch category customizations
  useEffect(() => {
    const fetchCategoryCustomizations = async () => {
      try {
        const response = await axios.get(`${API_BASE_URL}/api/v1/category-customizations`)
        const responseData = response.data
        
        let customizationsData = []
        if (responseData.success && responseData.data && responseData.data.category_customizations) {
          customizationsData = responseData.data.category_customizations
        } else if (Array.isArray(responseData.data)) {
          customizationsData = responseData.data
        } else {
          customizationsData = []
        }
        
        setCategoryCustomizations(customizationsData)
      } catch (error) {
        console.error('Error fetching category customizations:', error)
        // Continue with empty customizations
        setCategoryCustomizations([])
      }
    }

    fetchCategoryCustomizations()
  }, [])

  // Get the best category customization for a store
  const getStoreCustomization = (store: Store): { icon: string; color: string } => {
    // Find the first category that has a customization
    for (const category of store.categories || []) {
      const customization = categoryCustomizations.find(cc => cc.category_name === category)
      if (customization && customization.icon && customization.color) {
        return {
          icon: customization.icon,
          color: customization.color
        }
      }
    }
    
    // Fallback to store name first character and default color
    return {
      icon: store.name.charAt(0),
      color: '#007bff'
    }
  }

  // 地図の初期化
  useEffect(() => {
    let resizeObserver: ResizeObserver | null = null
    
    const initializeMap = () => {
      if (!mapContainerRef.current) return

      // 既存の地図があれば削除
      if (mapRef.current) {
        mapRef.current.remove()
        mapRef.current = null
      }

      // centerの型を統一
      const mapCenter: [number, number] = Array.isArray(center) 
        ? center 
        : [center.lat, center.lng]
      
      try {
        // 地図を作成
        const map = L.map(mapContainerRef.current, {
          preferCanvas: true,
          fadeAnimation: false,
          zoomAnimation: true,
          markerZoomAnimation: true
        }).setView(mapCenter, zoom)
    
        // OpenStreetMapタイルレイヤーを追加（最適化設定）
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
          attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
          maxZoom: 19,
          minZoom: 1,
          tileSize: 256,
          opacity: 1,
          keepBuffer: 2,
          updateWhenIdle: false,
          updateWhenZooming: true,
          crossOrigin: true
        }).addTo(map)

        // 初期化後に地図サイズを再計算
        setTimeout(() => {
          map.invalidateSize()
        }, 100)

        // 地図の中心とズーム変更を親コンポーネントに通知
        map.on('moveend', () => {
          const center = map.getCenter()
          const currentZoom = map.getZoom()
          if (onMapCenterChange) {
            onMapCenterChange([center.lat, center.lng])
          }
          if (onMapZoomChange) {
            onMapZoomChange(currentZoom)
          }
        })

        // ResizeObserverでコンテナサイズ変更を監視
        if (window.ResizeObserver && mapContainerRef.current) {
          resizeObserver = new ResizeObserver(() => {
            if (map) {
              // デバウンスして地図サイズを再計算
              setTimeout(() => {
                map.invalidateSize()
              }, 100)
            }
          })
          resizeObserver.observe(mapContainerRef.current)
        }

        mapRef.current = map
      } catch (error) {
        console.error('Error initializing map:', error)
      }
    }

    // DOM要素が確実に存在することを確認してから初期化
    if (mapContainerRef.current) {
      const container = mapContainerRef.current
      if (container.offsetWidth > 0 && container.offsetHeight > 0) {
        initializeMap()
      } else {
        // コンテナのサイズが0の場合は段階的に待機
        let retryCount = 0
        const maxRetries = 10
        const retryDelay = 50
        
        const retryInit = () => {
          if (retryCount >= maxRetries) {
            console.warn('Map container size check timed out, initializing anyway')
            initializeMap()
            return
          }
          
          if (container.offsetWidth > 0 && container.offsetHeight > 0) {
            initializeMap()
          } else {
            retryCount++
            setTimeout(retryInit, retryDelay * retryCount)
          }
        }
        
        retryInit()
      }
    }

    return () => {
      if (resizeObserver) {
        resizeObserver.disconnect()
      }
      if (mapRef.current) {
        mapRef.current.remove()
        mapRef.current = null
      }
    }
  }, [center, zoom])

  // マーカーの更新
  useEffect(() => {
    if (!mapRef.current) return

    // 既存のマーカーを削除
    markersRef.current.forEach(marker => {
      mapRef.current?.removeLayer(marker)
    })
    markersRef.current = []

    // storesからマーカーを追加
    if (stores && Array.isArray(stores) && stores.length > 0) {
      stores.forEach(store => {
        const isSelected = selectedStore?.id === store.id
        const customization = getStoreCustomization(store)
      
      // カスタムアイコンを作成
      const icon = L.divIcon({
        className: `custom-marker ${isSelected ? 'selected' : ''}`,
        html: `
          <div class="marker-pin ${isSelected ? 'selected' : ''}" style="background-color: ${customization.color};">
            <div class="marker-content">
              <span class="marker-text">${customization.icon}</span>
            </div>
          </div>
        `,
        iconSize: [30, 30],
        iconAnchor: [15, 30],
        popupAnchor: [0, -30]
      })

      const marker = L.marker([store.latitude, store.longitude], { icon })
        .addTo(mapRef.current!)

      // ポップアップを作成
      const popupContent = `
        <div class="store-popup">
          <h3>${store.name}</h3>
          <p><strong>住所:</strong> ${store.address}</p>
          ${store.categories.length > 0 ? `<p><strong>カテゴリ:</strong> ${store.categories.join(', ')}</p>` : ''}
          ${store.tags && store.tags.length > 0 ? `<p><strong>タグ:</strong> ${store.tags.join(', ')}</p>` : ''}
          <button class="popup-detail-btn" data-store-id="${store.id}">詳細を見る</button>
        </div>
      `
      
      marker.bindPopup(popupContent)

      // マーカークリック時のイベント
      marker.on('click', () => {
        if (onStoreClick) {
          onStoreClick(store)
        }
      })

        markersRef.current.push(marker)
      })
    }

    // markersからマーカーを追加（storesがない場合）
    if (markers && Array.isArray(markers) && markers.length > 0 && stores.length === 0) {
      markers.forEach(markerData => {
        const icon = L.divIcon({
          className: 'custom-marker',
          html: `
            <div class="marker-pin">
              <div class="marker-content">
                <span class="marker-text">${markerData.title.charAt(0)}</span>
              </div>
            </div>
          `,
          iconSize: [30, 30],
          iconAnchor: [15, 30],
          popupAnchor: [0, -30]
        })

        const marker = L.marker([markerData.position.lat, markerData.position.lng], { icon })
          .addTo(mapRef.current!)

        // ポップアップを作成
        const popupContent = `
          <div class="store-popup">
            <h3>${markerData.title}</h3>
          </div>
        `
        
        marker.bindPopup(popupContent)

        // マーカークリック時のイベント
        marker.on('click', () => {
          if (onMarkerClick) {
            onMarkerClick(markerData)
          }
        })

        markersRef.current.push(marker)
      })
    }

    // ポップアップ内のボタンクリックイベント
    const handlePopupButtonClick = (e: Event) => {
      const target = e.target as HTMLElement
      if (target.classList.contains('popup-detail-btn')) {
        const storeId = target.getAttribute('data-store-id')
        const store = stores.find(s => s.id === storeId)
        if (store && onStoreClick) {
          onStoreClick(store)
        }
      }
    }

    document.addEventListener('click', handlePopupButtonClick)

    return () => {
      document.removeEventListener('click', handlePopupButtonClick)
    }
  }, [stores, selectedStore, onStoreClick, categoryCustomizations])

  // 選択された店舗にフォーカス
  useEffect(() => {
    if (selectedStore && mapRef.current) {
      mapRef.current.setView([selectedStore.latitude, selectedStore.longitude], 16)
    }
  }, [selectedStore])

  return (
    <div className="map-container">
      <div
        id={mapId}
        ref={mapContainerRef}
        className="map"
        style={{ height, width: '100%', minHeight: '300px' }}
      />
    </div>
  )
}

export default Map