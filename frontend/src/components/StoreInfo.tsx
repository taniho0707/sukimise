import React from 'react'
import { Link } from 'react-router-dom'
import ImageGallery from './ImageGallery'
import ParkingDisplay from './ParkingDisplay'
import SafeBusinessHoursDisplay from './SafeBusinessHoursDisplay'

interface Store {
  id: string
  name: string
  address: string
  latitude: number
  longitude: number
  categories: string[]
  business_hours: string
  parking_info: string
  website_url: string
  google_map_url: string
  sns_urls: string[]
  tags: string[]
  photos: string[]
  created_by: string
  created_at: string
  updated_at: string
}

interface StoreInfoProps {
  store: Store
  isViewer?: boolean
}

const StoreInfo: React.FC<StoreInfoProps> = ({ store, isViewer = false }) => {
  return (
    <>
      {/* 画像ギャラリー */}
      {store.photos && store.photos.length > 0 && (
        <div className="store-images-section">
          <div className="store-images">
            <ImageGallery images={store.photos} />
          </div>
        </div>
      )}

      {/* 店舗情報グリッド */}
      <div className="store-info-grid">
        <div className="store-info-main">
          {/* 営業時間 */}
          {store.business_hours && (
            <div className="info-section">
              <h3>営業時間</h3>
              <div className="business-hours">
                <SafeBusinessHoursDisplay businessHours={store.business_hours} />
              </div>
            </div>
          )}

          {/* 駐車場情報 */}
          {store.parking_info && (
            <div className="info-section">
              <h3>駐車場情報</h3>
              <ParkingDisplay parkingInfo={store.parking_info} />
            </div>
          )}
        </div>

        <div className="store-info-side">
          {/* Web情報 */}
          <div className="info-section">
            <h3>Web情報</h3>
            <div className="web-links">
              {store.website_url && (
                <a 
                  href={store.website_url} 
                  target="_blank" 
                  rel="noopener noreferrer"
                  className="web-link"
                >
                  🌐 ホームページ
                </a>
              )}
              {store.google_map_url && (
                <a 
                  href={store.google_map_url} 
                  target="_blank" 
                  rel="noopener noreferrer"
                  className="web-link"
                >
                  📍 GoogleMap
                </a>
              )}
              {store.sns_urls && store.sns_urls.length > 0 && (
                store.sns_urls.map((url, index) => (
                  <a 
                    key={index}
                    href={url} 
                    target="_blank" 
                    rel="noopener noreferrer"
                    className="web-link"
                  >
                    🔗 SNS {index + 1}
                  </a>
                ))
              )}
            </div>
          </div>

          {/* 位置情報 */}
          <div className="info-section">
            <h3>位置情報</h3>
            <div className="location-info">
              <div className="coordinates">
                <div>緯度: {store.latitude}</div>
                <div>経度: {store.longitude}</div>
              </div>
              <Link 
                to={isViewer ? `/viewer/map?store=${store.id}` : `/map?store=${store.id}`} 
                className="btn btn-secondary map-btn"
              >
                地図で表示
              </Link>
            </div>
          </div>

          {/* メタ情報 */}
          <div className="info-section">
            <h3>情報</h3>
            <div className="meta-info">
              <div>登録日: {new Date(store.created_at).toLocaleDateString('ja-JP')}</div>
              <div>更新日: {new Date(store.updated_at).toLocaleDateString('ja-JP')}</div>
            </div>
          </div>
        </div>
      </div>
    </>
  )
}

export default StoreInfo