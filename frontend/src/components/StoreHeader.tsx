import React from 'react'
import { Link } from 'react-router-dom'
import { BusinessHoursData } from '../types/store'

interface Store {
  id: string
  name: string
  address: string
  latitude: number
  longitude: number
  categories: string[]
  business_hours: BusinessHoursData | string
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

interface StoreHeaderProps {
  store: Store
  canEdit?: boolean
  onEdit?: () => void
  onDelete?: () => void
}

const StoreHeader: React.FC<StoreHeaderProps> = ({ store, canEdit, onEdit, onDelete }) => {
  return (
    <div className="store-header">
      <div className="store-header-content">
        <div className="store-title-section">
          <h1 className="store-name">{store.name}</h1>
          <div className="store-meta">
            <span className="store-address">{store.address}</span>
          </div>
          
          {store.categories && store.categories.length > 0 && (
            <div className="categories">
              {store.categories.map((category, index) => (
                <span key={index} className="category-tag">
                  {category}
                </span>
              ))}
            </div>
          )}
          
          {store.tags && store.tags.length > 0 && (
            <div className="tags">
              {store.tags.map((tag, index) => (
                <span key={index} className="tag">
                  {tag}
                </span>
              ))}
            </div>
          )}
        </div>
        
        <div className="store-actions">
          <Link to="/stores" className="btn btn-secondary">
            店舗一覧に戻る
          </Link>
          <Link to={`/map?store=${store.id}`} className="btn btn-secondary">
            地図で見る
          </Link>
          {canEdit && (
            <>
              <button 
                onClick={onEdit}
                className="btn btn-secondary"
              >
                編集
              </button>
              <button 
                onClick={onDelete}
                className="btn btn-danger"
              >
                削除
              </button>
            </>
          )}
        </div>
      </div>
    </div>
  )
}

export default StoreHeader