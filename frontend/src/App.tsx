import { Routes, Route } from 'react-router-dom'
import { AuthProvider } from '@/contexts/AuthContext'
import { ViewerAuthProvider } from '@/contexts/ViewerAuthContext'
import Layout from '@/components/Layout'
import ViewerLayout from '@/components/ViewerLayout'
import ProtectedRoute from '@/components/ProtectedRoute'
import ViewerProtectedRoute from '@/components/ViewerProtectedRoute'
import Home from '@/pages/Home'
import Login from '@/pages/Login'
import ViewerLogin from '@/pages/ViewerLogin'
import ViewerHome from '@/pages/ViewerHome'
import AdminViewerSettings from '@/pages/AdminViewerSettings'
import StoreList from '@/pages/StoreList'
import ViewerStoreList from '@/pages/ViewerStoreList'
import StoreDetail from '@/pages/StoreDetail'
import ViewerStoreDetail from '@/pages/ViewerStoreDetail'
import StoreForm from '@/pages/StoreForm'
import MapView from '@/pages/MapView'
import ViewerMapView from '@/pages/ViewerMapView'
import CategoryManagement from '@/pages/CategoryManagement'
import ApiTest from '@/pages/ApiTest'
import Debug from '@/pages/Debug'

function App() {
  return (
    <AuthProvider>
      <ViewerAuthProvider>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/viewer-login" element={<ViewerLogin />} />
          <Route path="/api-test" element={<ApiTest />} />
          <Route path="/debug" element={<Debug />} />
          
          {/* 編集者・管理者用ルート */}
          <Route path="/" element={<ProtectedRoute><Layout /></ProtectedRoute>}>
            <Route index element={<Home />} />
            <Route path="stores" element={<StoreList />} />
            <Route path="stores/:id" element={<StoreDetail />} />
            <Route path="stores/new" element={<StoreForm />} />
            <Route path="stores/:id/edit" element={<StoreForm />} />
            <Route path="map" element={<MapView />} />
            <Route path="admin/viewer-settings" element={<AdminViewerSettings />} />
            <Route path="admin/category-management" element={<CategoryManagement />} />
          </Route>
          
          {/* 閲覧者専用ルート */}
          <Route path="/viewer" element={<ViewerProtectedRoute><ViewerLayout /></ViewerProtectedRoute>}>
            <Route index element={<ViewerHome />} />
            <Route path="stores" element={<ViewerStoreList />} />
            <Route path="stores/:id" element={<ViewerStoreDetail />} />
            <Route path="map" element={<ViewerMapView />} />
          </Route>
        </Routes>
      </ViewerAuthProvider>
    </AuthProvider>
  )
}

export default App