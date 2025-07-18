// CSRF トークン管理ユーティリティ

/**
 * クッキーからCSRFトークンを取得
 */
export const getCSRFToken = (): string | null => {
  const cookies = document.cookie.split(';')
  for (const cookie of cookies) {
    const [name, value] = cookie.trim().split('=')
    if (name === 'csrf_token') {
      return decodeURIComponent(value)
    }
  }
  return null
}

/**
 * CSRFトークンをHTTPヘッダーに設定するためのヘッダーオブジェクトを返す
 */
export const getCSRFHeaders = (): Record<string, string> => {
  const token = getCSRFToken()
  if (token) {
    return {
      'X-CSRF-Token': token
    }
  }
  return {}
}

/**
 * CSRFトークンを含むAxios設定オブジェクトを返す
 */
export const getCSRFConfig = () => {
  const headers = getCSRFHeaders()
  return {
    headers
  }
}

/**
 * バックエンドからCSRFトークンを取得（初回アクセス時）
 */
export const fetchCSRFToken = async (apiBaseUrl: string): Promise<void> => {
  try {
    // GETリクエストでCSRFトークンを取得（自動的にクッキーに設定される）
    await fetch(`${apiBaseUrl}/health`, {
      method: 'GET',
      credentials: 'include' // クッキーを含める
    })
  } catch (error) {
    console.warn('CSRFトークンの取得に失敗しました:', error)
  }
}