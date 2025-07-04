import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import axios from 'axios'
import { API_BASE_URL } from '@/config'
import './Home.css'

interface Stats {
  total_stores: number
  total_reviews: number
  recent_stores: Array<{
    id: string
    name: string
    address: string
    categories: string[]
  }>
}

const ViewerHome: React.FC = () => {
  const [stats, setStats] = useState<Stats | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchStats()
  }, [])

  const fetchStats = async () => {
    try {
      // 統計情報を取得（既存のAPIを活用）
      const storesResponse = await axios.get(`${API_BASE_URL}/api/v1/stores?limit=5`)
      const responseData = storesResponse.data
      
      // レスポンス構造を確認
      let stores = []
      if (responseData.success && responseData.data && responseData.data.stores) {
        stores = responseData.data.stores
      } else if (Array.isArray(responseData.data)) {
        stores = responseData.data
      } else if (Array.isArray(responseData)) {
        stores = responseData
      }
      
      // 簡易統計を作成
      const statsData: Stats = {
        total_stores: responseData.meta?.total || stores.length,
        total_reviews: 0, // レビュー総数は個別に計算が必要
        recent_stores: stores.slice(0, 5)
      }
      
      setStats(statsData)
    } catch (error) {
      console.error('Error fetching stats:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="loading">読み込み中...</div>
  }

  return (
    <div className="home-page">
      <div className="hero-section">
        <h1>Sukimise へようこそ</h1>
        <p className="hero-subtitle">
          お気に入りの店舗を検索・閲覧できます
        </p>
        <div className="hero-actions">
          <Link to="/viewer/stores" className="btn btn-primary btn-large">
            店舗一覧を見る
          </Link>
          <Link to="/viewer/map" className="btn btn-secondary btn-large">
            地図から探す
          </Link>
        </div>
      </div>

      {stats && (
        <div className="stats-section">
          <div className="stats-grid">
            <div className="stat-card">
              <h3>{stats.total_stores}</h3>
              <p>登録店舗数</p>
            </div>
          </div>
        </div>
      )}

      {stats && stats.recent_stores.length > 0 && (
        <div className="recent-stores-section">
          <h2>新着店舗</h2>
          <div className="store-grid">
            {stats.recent_stores.map((store) => (
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
              </div>
            ))}
          </div>
          <div className="section-footer">
            <Link to="/viewer/stores" className="view-all-link">
              すべての店舗を見る →
            </Link>
          </div>
        </div>
      )}

      <div className="features-section">
        <h2>機能</h2>
        <div className="features-grid">
          <div className="feature-card">
            <h3>🔍 店舗検索</h3>
            <p>名前、カテゴリ、タグから店舗を検索できます</p>
          </div>
          <div className="feature-card">
            <h3>🗺️ 地図表示</h3>
            <p>地図上で店舗の位置を確認できます</p>
          </div>
          <div className="feature-card">
            <h3>📄 詳細情報</h3>
            <p>営業時間、写真、レビューなど詳細情報を閲覧できます</p>
          </div>
          <div className="feature-card">
            <h3>📊 CSVエクスポート</h3>
            <p>店舗データをCSVファイルでダウンロードできます</p>
          </div>
        </div>
      </div>
    </div>
  )
}

export default ViewerHome