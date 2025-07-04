import React, { useEffect, useRef } from 'react'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import './Map.css'

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
}

const Map: React.FC<MapProps> = ({
  stores = [],
  markers = [],
  center = [35.6762, 139.6503], // 東京駅
  zoom = 13,
  height = '400px',
  onStoreClick,
  onMarkerClick,
  selectedStore
}) => {
  const mapRef = useRef<L.Map | null>(null)
  const markersRef = useRef<L.Marker[]>([])
  const mapContainerRef = useRef<HTMLDivElement>(null)

  // 地図の初期化
  useEffect(() => {
    if (!mapContainerRef.current) return

    // 既存の地図があれば削除
    if (mapRef.current) {
      mapRef.current.remove()
    }

    // centerの型を統一
    const mapCenter: [number, number] = Array.isArray(center) 
      ? center 
      : [center.lat, center.lng]
    
    // 地図を作成
    const map = L.map(mapContainerRef.current).setView(mapCenter, zoom)
    
    // OpenStreetMapタイルレイヤーを追加
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map)

    mapRef.current = map

    return () => {
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
      
      // カスタムアイコンを作成
      const icon = L.divIcon({
        className: `custom-marker ${isSelected ? 'selected' : ''}`,
        html: `
          <div class="marker-pin ${isSelected ? 'selected' : ''}">
            <div class="marker-content">
              <span class="marker-text">${store.name.charAt(0)}</span>
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
  }, [stores, selectedStore, onStoreClick])

  // 選択された店舗にフォーカス
  useEffect(() => {
    if (selectedStore && mapRef.current) {
      mapRef.current.setView([selectedStore.latitude, selectedStore.longitude], 16)
    }
  }, [selectedStore])

  return (
    <div className="map-container">
      <div
        ref={mapContainerRef}
        className="map"
        style={{ height }}
      />
    </div>
  )
}

export default Map