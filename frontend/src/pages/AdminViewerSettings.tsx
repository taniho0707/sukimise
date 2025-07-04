import React, { useState, useEffect } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import toast from 'react-hot-toast'
import axios from 'axios'
import { API_BASE_URL } from '@/config'
import './AdminViewerSettings.css'

interface ViewerSettings {
  id: string
  session_duration_days: number
  created_at: string
  updated_at: string
}

interface ViewerLoginHistory {
  id: string
  ip_address: string
  user_agent: string
  login_time: string
  expires_at: string
}

const AdminViewerSettings: React.FC = () => {
  const { user } = useAuth()
  const [settings, setSettings] = useState<ViewerSettings | null>(null)
  const [history, setHistory] = useState<ViewerLoginHistory[]>([])
  const [loading, setLoading] = useState(false)
  const [historyLoading, setHistoryLoading] = useState(false)
  const [password, setPassword] = useState('')
  const [sessionDuration, setSessionDuration] = useState(7)
  const [showPassword, setShowPassword] = useState(false)

  useEffect(() => {
    loadSettings()
    loadHistory()
  }, [])

  const loadSettings = async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/api/v1/admin/viewer-settings`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        }
      })
      setSettings(response.data)
      setSessionDuration(response.data.session_duration_days)
    } catch (error) {
      toast.error('設定の読み込みに失敗しました')
    }
  }

  const loadHistory = async () => {
    try {
      setHistoryLoading(true)
      const response = await axios.get(`${API_BASE_URL}/api/v1/admin/viewer-history`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        }
      })
      setHistory(response.data.history || [])
    } catch (error) {
      toast.error('履歴の読み込みに失敗しました')
    } finally {
      setHistoryLoading(false)
    }
  }

  const handleUpdateSettings = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!password) {
      toast.error('パスワードを入力してください')
      return
    }

    setLoading(true)
    try {
      await axios.put(`${API_BASE_URL}/api/v1/admin/viewer-settings`, {
        password,
        session_duration_days: sessionDuration
      }, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        }
      })
      toast.success('設定を更新しました')
      setPassword('')
      loadSettings()
    } catch (error) {
      toast.error('設定の更新に失敗しました')
    } finally {
      setLoading(false)
    }
  }

  const handleCleanupSessions = async () => {
    try {
      await axios.post(`${API_BASE_URL}/api/v1/admin/viewer-cleanup`, {}, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        }
      })
      toast.success('期限切れセッションを削除しました')
      loadHistory()
    } catch (error) {
      toast.error('セッションの削除に失敗しました')
    }
  }

  const formatDateTime = (dateString: string) => {
    return new Date(dateString).toLocaleString('ja-JP')
  }

  if (user?.role !== 'admin') {
    return (
      <div className="admin-settings-page">
        <div className="access-denied">
          <h2>アクセスが拒否されました</h2>
          <p>この機能は管理者のみ利用できます。</p>
        </div>
      </div>
    )
  }

  return (
    <div className="admin-settings-page">
      <div className="admin-settings-container">
        <h1>閲覧者認証設定</h1>

        <div className="settings-section">
          <h2>パスワード設定</h2>
          <form onSubmit={handleUpdateSettings} className="settings-form">
            <div className="form-group">
              <label htmlFor="password">新しいパスワード</label>
              <div className="password-input-container">
                <input
                  type={showPassword ? 'text' : 'password'}
                  id="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="form-input"
                  required
                  disabled={loading}
                />
                <button
                  type="button"
                  className="password-toggle"
                  onClick={() => setShowPassword(!showPassword)}
                >
                  {showPassword ? '隠す' : '表示'}
                </button>
              </div>
            </div>

            <div className="form-group">
              <label htmlFor="sessionDuration">セッション有効期間（日）</label>
              <input
                type="number"
                id="sessionDuration"
                value={sessionDuration}
                onChange={(e) => setSessionDuration(parseInt(e.target.value))}
                className="form-input"
                min="1"
                max="365"
                required
                disabled={loading}
              />
            </div>

            <button
              type="submit"
              className="btn btn-primary"
              disabled={loading}
            >
              {loading ? '更新中...' : '設定を更新'}
            </button>
          </form>
        </div>

        <div className="settings-section">
          <div className="section-header">
            <h2>ログイン履歴</h2>
            <button
              onClick={handleCleanupSessions}
              className="btn btn-secondary"
            >
              期限切れセッション削除
            </button>
          </div>

          {historyLoading ? (
            <div className="loading">履歴を読み込み中...</div>
          ) : (
            <div className="history-table">
              {history.length === 0 ? (
                <p>ログイン履歴がありません</p>
              ) : (
                <table>
                  <thead>
                    <tr>
                      <th>ログイン時間</th>
                      <th>IPアドレス</th>
                      <th>ユーザーエージェント</th>
                      <th>有効期限</th>
                    </tr>
                  </thead>
                  <tbody>
                    {history.map((item) => (
                      <tr key={item.id}>
                        <td>{formatDateTime(item.login_time)}</td>
                        <td>{item.ip_address}</td>
                        <td className="user-agent">{item.user_agent}</td>
                        <td>{formatDateTime(item.expires_at)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </div>
          )}
        </div>

        {settings && (
          <div className="settings-section">
            <h2>現在の設定</h2>
            <div className="current-settings">
              <p><strong>セッション有効期間:</strong> {settings.session_duration_days}日</p>
              <p><strong>最終更新:</strong> {formatDateTime(settings.updated_at)}</p>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

export default AdminViewerSettings