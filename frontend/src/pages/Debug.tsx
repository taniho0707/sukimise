import React, { useState } from 'react'

const Debug: React.FC = () => {
  const [testResults, setTestResults] = useState<string[]>([])

  const addResult = (result: string) => {
    setTestResults(prev => [...prev, `${new Date().toLocaleTimeString()}: ${result}`])
  }

  const testFetch = async () => {
    try {
      addResult('フェッチテスト開始...')
      
      // 直接fetch APIを使用してテスト
      const response = await fetch('/api/v1/stores')
      addResult(`レスポンス状態: ${response.status} ${response.statusText}`)
      
      const data = await response.json()
      addResult(`レスポンスデータ: ${JSON.stringify(data)}`)
      
    } catch (error: any) {
      addResult(`エラー: ${error.message}`)
    }
  }

  const testLogin = async () => {
    try {
      addResult('ログインテスト開始...')
      
      const response = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          username: 'admin',
          password: 'admin123'
        })
      })
      
      addResult(`ログインレスポンス: ${response.status} ${response.statusText}`)
      
      if (response.ok) {
        const data = await response.json()
        addResult(`ログイン成功: ユーザー ${data.user.username}`)
      } else {
        const errorData = await response.text()
        addResult(`ログイン失敗: ${errorData}`)
      }
      
    } catch (error: any) {
      addResult(`ログインエラー: ${error.message}`)
    }
  }

  const clearResults = () => {
    setTestResults([])
  }

  return (
    <div style={{ padding: '20px', fontFamily: 'Arial, sans-serif' }}>
      <h1>🔍 Sukimise デバッグページ</h1>
      
      <div style={{ marginBottom: '20px' }}>
        <button 
          onClick={testFetch}
          style={{
            padding: '10px 20px',
            marginRight: '10px',
            backgroundColor: '#007bff',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer'
          }}
        >
          API接続テスト
        </button>
        
        <button 
          onClick={testLogin}
          style={{
            padding: '10px 20px',
            marginRight: '10px',
            backgroundColor: '#28a745',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer'
          }}
        >
          ログインテスト
        </button>
        
        <button 
          onClick={clearResults}
          style={{
            padding: '10px 20px',
            backgroundColor: '#6c757d',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer'
          }}
        >
          結果クリア
        </button>
      </div>

      <div style={{
        backgroundColor: '#f8f9fa',
        border: '1px solid #dee2e6',
        borderRadius: '4px',
        padding: '15px',
        minHeight: '200px',
        maxHeight: '400px',
        overflowY: 'auto'
      }}>
        <h3>テスト結果:</h3>
        {testResults.length === 0 ? (
          <p style={{ color: '#6c757d' }}>まだテストが実行されていません</p>
        ) : (
          <div style={{ fontFamily: 'monospace', fontSize: '14px' }}>
            {testResults.map((result, index) => (
              <div key={index} style={{ marginBottom: '5px' }}>
                {result}
              </div>
            ))}
          </div>
        )}
      </div>

      <div style={{ marginTop: '20px', fontSize: '14px', color: '#6c757d' }}>
        <h4>現在の状況:</h4>
        <ul>
          <li>フロントエンド: {window.location.origin}</li>
          <li>APIプロキシ: /api → backend:8080</li>
          <li>React Router: 動作中</li>
        </ul>
      </div>
    </div>
  )
}

export default Debug