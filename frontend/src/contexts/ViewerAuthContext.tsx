import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import axios from 'axios'
import { API_BASE_URL } from '@/config'

interface ViewerAuthContextType {
  isAuthenticated: boolean
  login: (password: string) => Promise<void>
  logout: () => void
  validateSession: () => Promise<boolean>
}

const ViewerAuthContext = createContext<ViewerAuthContextType | undefined>(undefined)

export const useViewerAuth = () => {
  const context = useContext(ViewerAuthContext)
  if (context === undefined) {
    throw new Error('useViewerAuth must be used within a ViewerAuthProvider')
  }
  return context
}

interface ViewerAuthProviderProps {
  children: ReactNode
}

export const ViewerAuthProvider: React.FC<ViewerAuthProviderProps> = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false)
  const [isLoading, setIsLoading] = useState<boolean>(true)

  useEffect(() => {
    validateSession()
  }, [])

  const login = async (password: string): Promise<void> => {
    try {
      const response = await axios.post(`${API_BASE_URL}/api/v1/viewer/auth`, {
        password
      })

      if (response.data.token) {
        localStorage.setItem('viewer_token', response.data.token)
        localStorage.setItem('viewer_expires_at', response.data.expires_at)
        setIsAuthenticated(true)
      }
    } catch (error) {
      throw new Error('Authentication failed')
    }
  }

  const logout = (): void => {
    localStorage.removeItem('viewer_token')
    localStorage.removeItem('viewer_expires_at')
    setIsAuthenticated(false)
  }

  const validateSession = async (): Promise<boolean> => {
    try {
      const token = localStorage.getItem('viewer_token')
      const expiresAt = localStorage.getItem('viewer_expires_at')

      if (!token || !expiresAt) {
        setIsAuthenticated(false)
        setIsLoading(false)
        return false
      }

      // Check if token is expired
      const expireTime = new Date(expiresAt)
      if (expireTime <= new Date()) {
        logout()
        setIsLoading(false)
        return false
      }

      // Validate with server
      const response = await axios.get(`${API_BASE_URL}/api/v1/viewer/validate`, {
        headers: {
          'X-Viewer-Token': token
        }
      })

      if (response.data.valid) {
        setIsAuthenticated(true)
        setIsLoading(false)
        return true
      } else {
        logout()
        setIsLoading(false)
        return false
      }
    } catch (error) {
      logout()
      setIsLoading(false)
      return false
    }
  }

  const value = {
    isAuthenticated,
    login,
    logout,
    validateSession
  }

  if (isLoading) {
    return <div>Loading...</div>
  }

  return (
    <ViewerAuthContext.Provider value={value}>
      {children}
    </ViewerAuthContext.Provider>
  )
}