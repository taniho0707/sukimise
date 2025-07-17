import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '@/contexts/AuthContext'
import { useStores, useReviews } from '@/hooks/useApi'
import './Home.css'

const Home: React.FC = () => {
  const { user } = useAuth()
  const { fetchStores } = useStores()
  const { fetchUserReviews } = useReviews()
  
  const [stats, setStats] = useState({
    totalStores: 0,
    userReviews: 0,
    monthlyVisits: 0,
    loading: true
  })

  useEffect(() => {
    const loadStats = async () => {
      try {
        // 登録店舗数を取得（認証不要）
        const storesData = await fetchStores()
        console.log('Stores API response:', storesData)
        
        // storesDataがオブジェクトでstoresプロパティを持つ場合の処理
        let stores = []
        if (storesData && typeof storesData === 'object') {
          if (Array.isArray(storesData)) {
            stores = storesData
          } else if ((storesData as any).stores && Array.isArray((storesData as any).stores)) {
            stores = (storesData as any).stores
          }
        }
        const totalStores = stores.length
        console.log('Total stores:', totalStores)

        let userReviews = 0
        let monthlyVisits = 0

        // ユーザーが認証済みの場合のみレビューを取得
        if (user) {
          try {
            const reviewsData = await fetchUserReviews()
            console.log('Reviews API response:', reviewsData)
            
            // reviewsDataがオブジェクトでreviewsプロパティを持つ場合の処理
            let reviews = []
            if (reviewsData && typeof reviewsData === 'object') {
              if (Array.isArray(reviewsData)) {
                reviews = reviewsData
              } else if ((reviewsData as any).reviews && Array.isArray((reviewsData as any).reviews)) {
                reviews = (reviewsData as any).reviews
              }
            }
            userReviews = reviews.length
            console.log('User reviews count:', userReviews)

            // 今月の訪問店舗数を計算
            const currentMonth = new Date().getMonth()
            const currentYear = new Date().getFullYear()
            monthlyVisits = reviews.filter((review: any) => {
              if (!review.visit_date || !review.is_visited) return false
              const visitDate = new Date(review.visit_date)
              return visitDate.getMonth() === currentMonth && visitDate.getFullYear() === currentYear
            }).length
            console.log('Monthly visits:', monthlyVisits)
          } catch (reviewError) {
            console.warn('レビューデータの取得に失敗しました（認証が必要）:', reviewError)
          }
        }

        setStats({
          totalStores,
          userReviews,
          monthlyVisits,
          loading: false
        })
      } catch (error) {
        console.error('統計データの取得に失敗しました:', error)
        setStats(prev => ({ ...prev, loading: false }))
      }
    }

    loadStats()
  }, [user, fetchStores, fetchUserReviews])

  return (
    <div className="home-page">
      <div className="hero-section">
        <h1 className="hero-title">お気に入りの店舗を記録・共有</h1>
        <p className="hero-subtitle">
          {user?.username}さん、お疲れ様です！<br />
          今日はどちらのお店に行かれましたか？
        </p>
        <div className="hero-actions">
          <Link to="/stores/new" className="btn btn-primary">
            新しい店舗を登録
          </Link>
          <Link to="/stores" className="btn btn-secondary">
            店舗一覧を見る
          </Link>
        </div>
      </div>

      <div className="feature-grid">
        <div className="feature-card">
          <h3>店舗管理</h3>
          <p>お気に入りの店舗情報を詳細に記録。住所、営業時間、価格帯など、必要な情報をすべて管理できます。</p>
          <Link to="/stores" className="feature-link">店舗一覧へ →</Link>
        </div>

        <div className="feature-card">
          <h3>地図表示</h3>
          <p>登録した店舗を地図上で確認。近くのお店を探したり、エリア別に店舗を検索することができます。</p>
          <Link to="/map" className="feature-link">地図を見る →</Link>
        </div>

        <div className="feature-card">
          <h3>レビュー機能</h3>
          <p>個人的な評価やメモを記録。他の編集者のレビューも確認して、店舗選びの参考にできます。</p>
          <Link to="/stores" className="feature-link">レビューを見る →</Link>
        </div>
      </div>

      <div className="stats-section">
        <div className="stat-item">
          <div className="stat-number">{stats.loading ? '-' : stats.totalStores}</div>
          <div className="stat-label">登録店舗数</div>
        </div>
        <div className="stat-item">
          <div className="stat-number">{stats.loading ? '-' : stats.userReviews}</div>
          <div className="stat-label">あなたのレビュー</div>
        </div>
        <div className="stat-item">
          <div className="stat-number">{stats.loading ? '-' : stats.monthlyVisits}</div>
          <div className="stat-label">今月の訪問店舗</div>
        </div>
      </div>
    </div>
  )
}

export default Home