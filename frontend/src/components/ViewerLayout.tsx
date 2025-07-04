import React from 'react'
import { Outlet, Link, useNavigate } from 'react-router-dom'
import { useViewerAuth } from '@/contexts/ViewerAuthContext'
import './Layout.css'

const ViewerLayout: React.FC = () => {
  const { logout } = useViewerAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/viewer-login')
  }

  return (
    <div className="layout">
      <header className="header">
        <div className="container">
          <div className="header-content">
            <Link to="/viewer" className="logo">
              Sukimise (閲覧者)
            </Link>
            <nav className="nav">
              <Link to="/viewer" className="nav-link">
                ホーム
              </Link>
              <Link to="/viewer/stores" className="nav-link">
                店舗一覧
              </Link>
              <Link to="/viewer/map" className="nav-link">
                地図
              </Link>
            </nav>
            <div className="user-menu">
              <span className="user-name">閲覧者</span>
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

export default ViewerLayout