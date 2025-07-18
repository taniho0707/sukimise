import React, { createContext, useContext, useState, useEffect, useCallback, ReactNode } from 'react'
import axios from 'axios'
import { API_BASE_URL } from '@/config'
import { fetchCSRFToken, getCSRFConfig } from '@/utils/csrf'

interface User {
  id: string
  username: string
  email: string
  role: string
}

interface AuthContextType {
  user: User | null
  token: string | null
  login: (username: string, password: string) => Promise<void>
  logout: () => void
  loading: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

interface AuthProviderProps {
  children: ReactNode
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const initializeAuth = async () => {
      // CSRFトークンを最初に取得
      await fetchCSRFToken(API_BASE_URL)
      
      const savedToken = localStorage.getItem('token')
      const savedUser = localStorage.getItem('user')
      
      if (savedToken && savedUser) {
        setToken(savedToken)
        setUser(JSON.parse(savedUser))
        
        // Axiosのデフォルトヘッダーを設定
        axios.defaults.headers.common['Authorization'] = `Bearer ${savedToken}`
      }
      
      setLoading(false)
    }
    
    initializeAuth()
  }, [])

  const logout = useCallback(() => {
    setUser(null)
    setToken(null)
    
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    
    delete axios.defaults.headers.common['Authorization']
  }, [])

  // Axiosインターセプターを設定してトークンとCSRFトークンを自動的に追加
  useEffect(() => {
    const requestInterceptor = axios.interceptors.request.use(
      (config) => {
        const currentToken = localStorage.getItem('token')
        if (currentToken) {
          config.headers.Authorization = `Bearer ${currentToken}`
        }
        
        // POSTリクエストなどにCSRFトークンを追加
        if (config.method !== 'get' && config.method !== 'head' && config.method !== 'options') {
          const csrfConfig = getCSRFConfig()
          Object.assign(config.headers, csrfConfig.headers)
        }
        
        return config
      },
      (error) => {
        return Promise.reject(error)
      }
    )

    const responseInterceptor = axios.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          // トークンが無効な場合はログアウト
          logout()
        }
        return Promise.reject(error)
      }
    )

    return () => {
      axios.interceptors.request.eject(requestInterceptor)
      axios.interceptors.response.eject(responseInterceptor)
    }
  }, [logout])

  const login = async (username: string, password: string) => {
    try {
      // CSRFトークンを含めてログインリクエストを送信
      const csrfConfig = getCSRFConfig()
      const response = await axios.post(`${API_BASE_URL}/api/v1/auth/login`, {
        username,
        password,
      }, csrfConfig)

      const { access_token, user: userData } = response.data
      
      setToken(access_token)
      setUser(userData)
      
      localStorage.setItem('token', access_token)
      localStorage.setItem('user', JSON.stringify(userData))
      
      axios.defaults.headers.common['Authorization'] = `Bearer ${access_token}`
    } catch (error) {
      console.error('Login failed:', error)
      throw error
    }
  }

  const value = {
    user,
    token,
    login,
    logout,
    loading,
  }

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  )
}