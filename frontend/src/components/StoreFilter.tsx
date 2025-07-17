import React, { useState, useEffect } from 'react'
import axios from 'axios'
import { API_BASE_URL } from '@/config'

interface FilterState {
  name: string
  categories: string[]
  categoriesOperator: string
  tags: string[]
  tagsOperator: string
  priceMin: number
  priceMax: number
  businessDay: string
  businessTime: string
}

interface StoreFilterProps {
  initialFilters: FilterState
  onFilterChange: (filters: FilterState) => void
  onSearch?: () => void
  onReset?: () => void
}


const BUSINESS_DAYS = [
  { value: '', label: '指定なし' },
  { value: 'monday', label: '月曜日' },
  { value: 'tuesday', label: '火曜日' },
  { value: 'wednesday', label: '水曜日' },
  { value: 'thursday', label: '木曜日' },
  { value: 'friday', label: '金曜日' },
  { value: 'saturday', label: '土曜日' },
  { value: 'sunday', label: '日曜日' },
]

// 30分単位の営業時間選択肢を生成
const generateTimeOptions = () => {
  const options = [{ value: '', label: '指定なし' }]
  
  for (let hour = 0; hour < 24; hour++) {
    for (let minute of [0, 30]) {
      const timeString = `${hour.toString().padStart(2, '0')}:${minute.toString().padStart(2, '0')}`
      const displayTime = `${hour}:${minute.toString().padStart(2, '0')}`
      options.push({ value: timeString, label: displayTime })
    }
  }
  
  return options
}

const BUSINESS_TIMES = generateTimeOptions()

const StoreFilter: React.FC<StoreFilterProps> = ({ initialFilters, onFilterChange, onSearch, onReset }) => {
  const [filters, setFilters] = useState<FilterState>(initialFilters)
  const [availableCategories, setAvailableCategories] = useState<string[]>([])
  const [availableTags, setAvailableTags] = useState<string[]>([])

  // initialFiltersが変更された場合にfiltersを更新
  useEffect(() => {
    setFilters(initialFilters)
  }, [initialFilters])


  useEffect(() => {
    const fetchOptions = async () => {
      try {
        console.log('Fetching categories and tags...')
        
        const [categoriesRes, tagsRes] = await Promise.all([
          axios.get(`${API_BASE_URL}/api/v1/stores/categories`),
          axios.get(`${API_BASE_URL}/api/v1/stores/tags`),
        ])
        
        console.log('Categories response:', categoriesRes.data)
        console.log('Tags response:', tagsRes.data)
        
        // Handle new API response format
        const categories = categoriesRes.data.success && categoriesRes.data.data 
          ? categoriesRes.data.data.categories 
          : categoriesRes.data.categories || []
        const tags = tagsRes.data.success && tagsRes.data.data 
          ? tagsRes.data.data.tags 
          : tagsRes.data.tags || []
          
        setAvailableCategories(categories)
        setAvailableTags(tags)

        // 初回読み込み時のデフォルト値設定は行わない（指定なしのまま）
      } catch (error: any) {
        console.error('Failed to fetch filter options:', error)
        if (error.response) {
          console.error('Response status:', error.response.status)
          console.error('Response data:', error.response.data)
        }
      }
    }

    fetchOptions()
  }, [])

  const handleFilterChange = (newFilters: Partial<FilterState>) => {
    const updatedFilters = { ...filters, ...newFilters }
    setFilters(updatedFilters)
    onFilterChange(updatedFilters)
  }

  const handleCategoryToggle = (category: string) => {
    const newCategories = filters.categories.includes(category)
      ? filters.categories.filter(c => c !== category)
      : [...filters.categories, category]
    
    handleFilterChange({ categories: newCategories })
  }

  const handleTagToggle = (tag: string) => {
    const newTags = filters.tags.includes(tag)
      ? filters.tags.filter(t => t !== tag)
      : [...filters.tags, tag]
    
    handleFilterChange({ tags: newTags })
  }

  const handleReset = () => {
    onReset?.()
  }


  return (
    <div className="search-section">
      <div className="search-form">
        <div className="search-header">
          <h3>店舗検索・フィルタ</h3>
        </div>

        {/* 基本検索 */}
        <div className="search-row">
          <div className="form-group">
            <label className="form-label">店舗名</label>
            <input
              type="text"
              value={filters.name}
              onChange={(e) => handleFilterChange({ name: e.target.value })}
              onKeyPress={(e) => e.key === 'Enter' && onSearch?.()}
              placeholder="店舗名で検索..."
              className="form-input"
            />
          </div>
        </div>

        {/* 詳細フィルタ */}
        <div className="filter-details">
            {/* カテゴリ選択 */}
            <div className="form-group">
              <div className="filter-header">
                <label className="form-label">カテゴリ</label>
                <select
                  value={filters.categoriesOperator}
                  onChange={(e) => handleFilterChange({ categoriesOperator: e.target.value })}
                  className="operator-select"
                >
                  <option value="OR">いずれかに該当 (OR)</option>
                  <option value="AND">すべてに該当 (AND)</option>
                </select>
              </div>
              <div className="checkbox-grid">
                {availableCategories.map((category) => (
                  <label key={category} className="checkbox-item">
                    <input
                      type="checkbox"
                      checked={filters.categories.includes(category)}
                      onChange={() => handleCategoryToggle(category)}
                    />
                    <span>{category}</span>
                  </label>
                ))}
              </div>
            </div>

            {/* タグ選択 */}
            <div className="form-group">
              <div className="filter-header">
                <label className="form-label">タグ</label>
                <select
                  value={filters.tagsOperator}
                  onChange={(e) => handleFilterChange({ tagsOperator: e.target.value })}
                  className="operator-select"
                >
                  <option value="AND">すべてに該当 (AND)</option>
                  <option value="OR">いずれかに該当 (OR)</option>
                </select>
              </div>
              <div className="checkbox-grid tag-grid">
                {availableTags.map((tag) => (
                  <label key={tag} className="checkbox-item">
                    <input
                      type="checkbox"
                      checked={filters.tags.includes(tag)}
                      onChange={() => handleTagToggle(tag)}
                    />
                    <span>#{tag}</span>
                  </label>
                ))}
              </div>
            </div>

            {/* 営業時間 */}
            <div className="search-row">

              <div className="form-group">
                <label className="form-label">営業日</label>
                <div className="day-selector">
                  {BUSINESS_DAYS.map((day) => (
                    <button
                      key={day.value}
                      type="button"
                      onClick={() => handleFilterChange({ businessDay: day.value })}
                      className={`day-button ${filters.businessDay === day.value ? 'active' : ''}`}
                    >
                      {day.label}
                    </button>
                  ))}
                </div>
              </div>

              <div className="form-group">
                <label className="form-label">営業時間</label>
                <div className="time-selector">
                  {BUSINESS_TIMES.map((time) => (
                    <button
                      key={time.value}
                      type="button"
                      onClick={() => handleFilterChange({ businessTime: time.value })}
                      className={`time-button ${filters.businessTime === time.value ? 'active' : ''}`}
                    >
                      {time.label}
                    </button>
                  ))}
                </div>
              </div>
            </div>
          </div>

        {/* アクションボタン */}
        <div className="search-actions">
          <button onClick={onSearch} className="btn btn-primary">
            検索
          </button>
          <button onClick={handleReset} className="btn btn-secondary">
            リセット
          </button>
        </div>
      </div>
    </div>
  )
}

export default StoreFilter