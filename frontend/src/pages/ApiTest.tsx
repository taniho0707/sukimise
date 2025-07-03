import React, { useState } from 'react'
import axios from 'axios'

const ApiTest: React.FC = () => {
  const [apiStatus, setApiStatus] = useState<string>('未確認')
  const [loading, setLoading] = useState(false)

  const testApiConnection = async () => {
    setLoading(true)
    try {
      const response = await axios.get('/api/v1/stores')
      setApiStatus(`成功: ${response.status}`)
    } catch (error: any) {
      if (error.response) {
        setApiStatus(`エラー: ${error.response.status} - ${error.response.statusText}`)
      } else if (error.request) {
        setApiStatus('ネットワークエラー: バックエンドに接続できません')
      } else {
        setApiStatus(`エラー: ${error.message}`)
      }
    } finally {
      setLoading(false)
    }
  }

  const testLogin = async () => {
    setLoading(true)
    try {
      const response = await axios.post('/api/v1/auth/login', {
        username: 'admin',
        password: 'admin123'
      })
      setApiStatus(`ログイン成功: ${response.status}`)
      console.log('Login response:', response.data)
    } catch (error: any) {
      if (error.response) {
        setApiStatus(`ログインエラー: ${error.response.status} - ${error.response.data?.error || error.response.statusText}`)
      } else {
        setApiStatus(`ログインエラー: ${error.message}`)
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ padding: '20px', maxWidth: '600px', margin: '0 auto' }}>
      <h1>API接続テスト</h1>
      
      <div style={{ marginBottom: '20px', padding: '15px', backgroundColor: '#f8f9fa', border: '1px solid #dee2e6', borderRadius: '5px' }}>
        <h3>API状態: {apiStatus}</h3>
      </div>

      <div style={{ display: 'flex', gap: '10px', marginBottom: '20px' }}>
        <button 
          onClick={testApiConnection}
          disabled={loading}
          style={{
            padding: '10px 20px',
            backgroundColor: '#007bff',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: loading ? 'not-allowed' : 'pointer'
          }}
        >
          {loading ? '確認中...' : '店舗API確認'}
        </button>
        
        <button 
          onClick={testLogin}
          disabled={loading}
          style={{
            padding: '10px 20px',
            backgroundColor: '#28a745',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: loading ? 'not-allowed' : 'pointer'
          }}
        >
          {loading ? 'ログイン中...' : 'ログインテスト'}
        </button>
      </div>

      <div style={{ fontSize: '14px', color: '#6c757d' }}>
        <h4>テスト情報:</h4>
        <ul>
          <li>フロントエンド: http://localhost:3000</li>
          <li>バックエンドAPI: http://localhost:8080</li>
          <li>テストユーザー: admin / admin123</li>
        </ul>
      </div>
    </div>
  )
}

export default ApiTest