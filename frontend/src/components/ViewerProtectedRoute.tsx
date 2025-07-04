import React from 'react'
import { Navigate } from 'react-router-dom'
import { useViewerAuth } from '@/contexts/ViewerAuthContext'

interface ViewerProtectedRouteProps {
  children: React.ReactNode
}

const ViewerProtectedRoute: React.FC<ViewerProtectedRouteProps> = ({ children }) => {
  const { isAuthenticated } = useViewerAuth()

  if (!isAuthenticated) {
    return <Navigate to="/viewer-login" replace />
  }

  return <>{children}</>
}

export default ViewerProtectedRoute