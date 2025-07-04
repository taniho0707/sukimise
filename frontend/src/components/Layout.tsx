import React, { useEffect } from 'react'
import { Outlet, Link, useNavigate } from 'react-router-dom'
import { useAuth } from '@/contexts/AuthContext'
import './Layout.css'

const Layout: React.FC = () => {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  useEffect(() => {
    if (!user) {
      navigate('/login')
    }
  }, [user, navigate])

  if (!user) {
    return null
  }

  return (
    <div className="layout">
      <header className="header">
        <div className="container">
          <div className="header-content">
            <Link to="/" className="logo">
              Sukimise
            </Link>
            <nav className="nav">
              <Link to="/" className="nav-link">
                ホーム
              </Link>
              <Link to="/stores" className="nav-link">
                店舗一覧
              </Link>
              <Link to="/map" className="nav-link">
                地図
              </Link>
              <Link to="/stores/new" className="nav-link">
                店舗登録
              </Link>
              {user.role === 'admin' && (
                <Link to="/admin/viewer-settings" className="nav-link">
                  閲覧者設定
                </Link>
              )}
            </nav>
            <div className="user-menu">
              <span className="user-name">{user.username} ({user.role})</span>
              <button onClick={handleLogout} className="btn btn-secondary">
                ログアウト
              </button>
            </div>
          </div>
        </div>
      </header>
      <main className="main">
        <div className="container">
          <Outlet />
        </div>
      </main>
    </div>
  )
}

export default Layout