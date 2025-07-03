import { Routes, Route } from 'react-router-dom'
import { AuthProvider } from '@/contexts/AuthContext'
import Layout from '@/components/Layout'
import ProtectedRoute from '@/components/ProtectedRoute'
import Home from '@/pages/Home'
import Login from '@/pages/Login'
import StoreList from '@/pages/StoreList'
import StoreDetail from '@/pages/StoreDetail'
import StoreForm from '@/pages/StoreForm'
import MapView from '@/pages/MapView'
import ApiTest from '@/pages/ApiTest'
import Debug from '@/pages/Debug'

function App() {
  return (
    <AuthProvider>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/api-test" element={<ApiTest />} />
        <Route path="/debug" element={<Debug />} />
        <Route path="/" element={<ProtectedRoute><Layout /></ProtectedRoute>}>
          <Route index element={<Home />} />
          <Route path="stores" element={<StoreList />} />
          <Route path="stores/:id" element={<StoreDetail />} />
          <Route path="stores/new" element={<StoreForm />} />
          <Route path="stores/:id/edit" element={<StoreForm />} />
          <Route path="map" element={<MapView />} />
        </Route>
      </Routes>
    </AuthProvider>
  )
}

export default App