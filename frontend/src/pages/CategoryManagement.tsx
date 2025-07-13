import React, { useState, useEffect } from 'react'
import axios from 'axios'
import toast from 'react-hot-toast'
import { API_BASE_URL } from '@/config'
import './CategoryManagement.css'

interface CategoryCustomization {
  id: string
  category_name: string
  icon?: string
  color?: string
  created_at: string
  updated_at: string
}

interface CategoryCustomizationRequest {
  category_name: string
  icon?: string
  color?: string
}

const CategoryManagement: React.FC = () => {
  const [customizations, setCustomizations] = useState<CategoryCustomization[]>([])
  const [loading, setLoading] = useState(true)
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [editingCategory, setEditingCategory] = useState<CategoryCustomization | null>(null)
  const [formData, setFormData] = useState<CategoryCustomizationRequest>({
    category_name: '',
    icon: '',
    color: '#FF5733'
  })

  useEffect(() => {
    const initializeData = async () => {
      await fetchCategoryCustomizations()
      // Automatically sync categories on load to ensure consistency
      try {
        await axios.post(`${API_BASE_URL}/api/v1/admin/category-customizations/sync`)
        // Refetch after sync to show any newly added categories
        await fetchCategoryCustomizations()
      } catch (error) {
        console.error('Auto-sync failed:', error)
        // Continue with existing data if sync fails
      }
    }
    
    initializeData()
  }, [])

  const fetchCategoryCustomizations = async () => {
    try {
      setLoading(true)
      const response = await axios.get(`${API_BASE_URL}/api/v1/category-customizations`)
      const responseData = response.data
      
      let customizationsData = []
      if (responseData.success && responseData.data && responseData.data.category_customizations) {
        customizationsData = responseData.data.category_customizations
      } else if (Array.isArray(responseData.data)) {
        customizationsData = responseData.data
      } else {
        customizationsData = []
      }
      
      setCustomizations(customizationsData)
    } catch (error) {
      console.error('Error fetching category customizations:', error)
      toast.error('カテゴリカスタマイズ情報の取得に失敗しました')
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!formData.category_name.trim()) {
      toast.error('カテゴリ名を入力してください')
      return
    }

    try {
      if (editingCategory) {
        // Update existing customization
        await axios.put(
          `${API_BASE_URL}/api/v1/admin/category-customizations/${encodeURIComponent(editingCategory.category_name)}`,
          {
            category_name: formData.category_name,
            icon: formData.icon || null,
            color: formData.color || null
          }
        )
        toast.success('カテゴリカスタマイズを更新しました')
      } else {
        // Create new customization
        await axios.post(`${API_BASE_URL}/api/v1/admin/category-customizations`, {
          category_name: formData.category_name,
          icon: formData.icon || null,
          color: formData.color || null
        })
        toast.success('カテゴリカスタマイズを作成しました')
      }
      
      resetForm()
      await fetchCategoryCustomizations()
    } catch (error: any) {
      console.error('Error saving category customization:', error)
      const errorMessage = error.response?.data?.message || 'カテゴリカスタマイズの保存に失敗しました'
      toast.error(errorMessage)
    }
  }

  const handleEdit = (customization: CategoryCustomization) => {
    setEditingCategory(customization)
    setFormData({
      category_name: customization.category_name,
      icon: customization.icon || '',
      color: customization.color || '#FF5733'
    })
    setShowCreateForm(true)
  }

  const handleDelete = async (categoryName: string) => {
    if (!confirm(`カテゴリ「${categoryName}」のカスタマイズを削除しますか？`)) {
      return
    }

    try {
      await axios.delete(`${API_BASE_URL}/api/v1/admin/category-customizations/${encodeURIComponent(categoryName)}`)
      toast.success('カテゴリカスタマイズを削除しました')
      await fetchCategoryCustomizations()
    } catch (error) {
      console.error('Error deleting category customization:', error)
      toast.error('カテゴリカスタマイズの削除に失敗しました')
    }
  }

  const resetForm = () => {
    setFormData({
      category_name: '',
      icon: '',
      color: '#FF5733'
    })
    setEditingCategory(null)
    setShowCreateForm(false)
  }

  const handleSync = async () => {
    try {
      const response = await axios.post(`${API_BASE_URL}/api/v1/admin/category-customizations/sync`)
      if (response.data.success) {
        const syncedCategories = response.data.data.synchronized_categories
        toast.success(`${syncedCategories.length}個のカテゴリと同期しました`)
        await fetchCategoryCustomizations()
      }
    } catch (error: any) {
      console.error('Error syncing categories:', error)
      const errorMessage = error.response?.data?.message || 'カテゴリの同期に失敗しました'
      toast.error(errorMessage)
    }
  }

  const validateIcon = (icon: string): boolean => {
    if (!icon) return true // Empty is allowed
    const chars = Array.from(icon)
    return chars.length === 1
  }

  const validateColor = (color: string): boolean => {
    if (!color) return true // Empty is allowed
    return /^#[0-9A-Fa-f]{6}$/.test(color)
  }

  if (loading) {
    return <div className="loading">読み込み中...</div>
  }

  return (
    <div className="category-management">
      <div className="category-management-header">
        <h1>カテゴリの編集</h1>
        <div className="header-buttons">
          <button
            onClick={handleSync}
            className="btn btn-secondary"
            title="店舗検索で使用されているカテゴリと同期します"
          >
            店舗カテゴリと同期
          </button>
          <button
            onClick={() => setShowCreateForm(true)}
            className="btn btn-primary"
          >
            新しいカテゴリカスタマイズを追加
          </button>
        </div>
      </div>

      {showCreateForm && (
        <div className="modal-overlay">
          <div className="modal">
            <div className="modal-header">
              <h2>{editingCategory ? 'カテゴリカスタマイズを編集' : '新しいカテゴリカスタマイズを追加'}</h2>
              <button onClick={resetForm} className="close-btn">×</button>
            </div>
            
            <form onSubmit={handleSubmit} className="category-form">
              <div className="form-group">
                <label>カテゴリ名 *</label>
                <input
                  type="text"
                  value={formData.category_name}
                  onChange={(e) => setFormData({ ...formData, category_name: e.target.value })}
                  placeholder="レストラン"
                  required
                  disabled={!!editingCategory} // Disable editing category name for existing categories
                />
              </div>

              <div className="form-group">
                <label>アイコン (絵文字1文字または文字1文字)</label>
                <input
                  type="text"
                  value={formData.icon}
                  onChange={(e) => {
                    const value = e.target.value
                    if (validateIcon(value)) {
                      setFormData({ ...formData, icon: value })
                    }
                  }}
                  placeholder="🍽️"
                  maxLength={2} // Allow for emojis that might be 2 UTF-16 code units
                />
                {formData.icon && !validateIcon(formData.icon) && (
                  <small className="error">アイコンは1文字または絵文字1つにしてください</small>
                )}
              </div>

              <div className="form-group">
                <label>色 (HEXカラーコード)</label>
                <div className="color-input-group">
                  <input
                    type="color"
                    value={formData.color}
                    onChange={(e) => setFormData({ ...formData, color: e.target.value })}
                  />
                  <input
                    type="text"
                    value={formData.color}
                    onChange={(e) => {
                      const value = e.target.value
                      if (validateColor(value) || value.length <= 7) {
                        setFormData({ ...formData, color: value })
                      }
                    }}
                    placeholder="#FF5733"
                    pattern="^#[0-9A-Fa-f]{6}$"
                  />
                </div>
                {formData.color && !validateColor(formData.color) && (
                  <small className="error">正しいHEXカラーコード形式で入力してください（例：#FF5733）</small>
                )}
              </div>

              {formData.icon && formData.color && (
                <div className="preview">
                  <label>プレビュー:</label>
                  <div 
                    className="icon-preview"
                    style={{ 
                      backgroundColor: formData.color,
                      color: 'white'
                    }}
                  >
                    {formData.icon}
                  </div>
                </div>
              )}

              <div className="form-actions">
                <button type="button" onClick={resetForm} className="btn btn-secondary">
                  キャンセル
                </button>
                <button 
                  type="submit" 
                  className="btn btn-primary"
                  disabled={!formData.category_name.trim() || (formData.icon && !validateIcon(formData.icon)) || (formData.color && !validateColor(formData.color))}
                >
                  {editingCategory ? '更新' : '作成'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      <div className="customizations-list">
        <h2>登録済みカテゴリカスタマイズ ({customizations.length}件)</h2>
        
        {customizations.length === 0 ? (
          <div className="empty-state">
            <p>まだカテゴリカスタマイズが登録されていません。</p>
          </div>
        ) : (
          <div className="customizations-grid">
            {customizations.map((customization) => (
              <div key={customization.id} className="customization-card">
                <div className="card-header">
                  <div className="category-info">
                    {customization.icon && customization.color && (
                      <div 
                        className="category-icon"
                        style={{ 
                          backgroundColor: customization.color,
                          color: 'white'
                        }}
                      >
                        {customization.icon}
                      </div>
                    )}
                    <div>
                      <h3>{customization.category_name}</h3>
                      <div className="category-details">
                        {customization.icon && <span>アイコン: {customization.icon}</span>}
                        {customization.color && <span>色: {customization.color}</span>}
                      </div>
                    </div>
                  </div>
                  <div className="card-actions">
                    <button
                      onClick={() => handleEdit(customization)}
                      className="btn btn-sm btn-secondary"
                    >
                      編集
                    </button>
                    <button
                      onClick={() => handleDelete(customization.category_name)}
                      className="btn btn-sm btn-danger"
                    >
                      削除
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

export default CategoryManagement