import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useViewerAuth } from '@/contexts/ViewerAuthContext'
import toast from 'react-hot-toast'
import './Login.css'

const ViewerLogin: React.FC = () => {
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const { login } = useViewerAuth()
  const navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      await login(password)
      toast.success('ログインしました')
      navigate('/viewer')
    } catch (error) {
      toast.error('パスワードが間違っています')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="login-page">
      <div className="login-container">
        <div className="login-card">
          <h1 className="login-title">Sukimise</h1>
          <p className="login-subtitle">閲覧者向けログイン</p>
          
          <form onSubmit={handleSubmit} className="login-form">
            <div className="form-group">
              <label htmlFor="password" className="form-label">
                パスワード
              </label>
              <input
                type="password"
                id="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="form-input"
                required
                disabled={loading}
                placeholder="管理者から提供されたパスワードを入力"
              />
            </div>
            
            <button
              type="submit"
              className="btn btn-primary login-button"
              disabled={loading}
            >
              {loading ? 'ログイン中...' : 'ログイン'}
            </button>
          </form>
          
          <div className="login-note">
            <p>このページは閲覧者向けです。</p>
            <p>編集者の方は<a href="/login">こちら</a>からログインしてください。</p>
          </div>
        </div>
      </div>
    </div>
  )
}

export default ViewerLogin