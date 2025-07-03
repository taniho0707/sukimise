import React, { useState, useEffect } from 'react'
import axios from 'axios'

interface TagCategoryInputProps {
  label: string
  value: string
  onChange: (value: string) => void
  apiEndpoint: string
  placeholder: string
  className?: string
}

const TagCategoryInput: React.FC<TagCategoryInputProps> = ({
  label,
  value,
  onChange,
  apiEndpoint,
  placeholder,
  className
}) => {
  const [suggestions, setSuggestions] = useState<string[]>([])
  const [selectedItems, setSelectedItems] = useState<string[]>([])
  const [inputValue, setInputValue] = useState('')
  const [originalSuggestions, setOriginalSuggestions] = useState<string[]>([]) // 元々の候補を保持

  useEffect(() => {
    // valueからselectedItemsを初期化
    const items = value ? value.split(',').map(item => item.trim()).filter(item => item) : []
    setSelectedItems(items)
  }, [value])

  useEffect(() => {
    const fetchSuggestions = async () => {
      try {
        console.log('Fetching suggestions from:', apiEndpoint)
        const response = await axios.get(apiEndpoint)
        const data = response.data
        console.log('API response:', data)
        
        if (apiEndpoint.includes('categories')) {
          const categories = data.data?.categories || data.categories || []
          console.log('Setting categories:', categories)
          setSuggestions(categories)
          setOriginalSuggestions(categories)
        } else if (apiEndpoint.includes('tags')) {
          const tags = data.data?.tags || data.tags || []
          console.log('Setting tags:', tags)
          setSuggestions(tags)
          setOriginalSuggestions(tags)
        }
      } catch (error) {
        console.error('Failed to fetch suggestions:', error)
      }
    }
    
    fetchSuggestions()
  }, [apiEndpoint])

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value
    setInputValue(newValue)
  }

  const handleInputKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && inputValue.trim()) {
      e.preventDefault()
      handleAddNewItem(inputValue.trim())
      setInputValue('')
    }
  }

  const handleAddNewItem = (newItem: string) => {
    // 既に存在するかチェック
    if (selectedItems.includes(newItem) || suggestions.includes(newItem)) {
      // 既に存在する場合はトグル
      handleToggleItem(newItem)
    } else {
      // 新しいアイテムの場合は候補に追加して選択
      setSuggestions(prev => [...prev, newItem])
      const newItems = [...selectedItems, newItem]
      setSelectedItems(newItems)
      onChange(newItems.join(', '))
    }
  }

  const handleAddButtonClick = () => {
    if (inputValue.trim()) {
      handleAddNewItem(inputValue.trim())
      setInputValue('')
    }
  }

  const handleToggleItem = (item: string) => {
    if (selectedItems.includes(item)) {
      // 既に選択されている場合は削除
      const newItems = selectedItems.filter(selected => selected !== item)
      setSelectedItems(newItems)
      onChange(newItems.join(', '))
    } else {
      // 選択されていない場合は追加
      const newItems = [...selectedItems, item]
      setSelectedItems(newItems)
      onChange(newItems.join(', '))
    }
  }

  return (
    <div className={`tag-category-input ${className || ''}`}>
      <label htmlFor={label} className="form-label">
        {label}
      </label>
      
      {/* カスタム入力フィールド（新しいアイテムを追加用） */}
      <div className="form-group" style={{ marginBottom: '12px' }}>
        <div style={{ display: 'flex', gap: '8px' }}>
          <input
            type="text"
            id={label}
            value={inputValue}
            onChange={handleInputChange}
            onKeyPress={handleInputKeyPress}
            className="form-input"
            placeholder={placeholder}
            style={{ fontSize: '14px', flex: 1 }}
          />
          <button
            type="button"
            onClick={handleAddButtonClick}
            disabled={!inputValue.trim()}
            style={{
              padding: '8px 16px',
              backgroundColor: inputValue.trim() ? '#0066cc' : '#e9ecef',
              color: inputValue.trim() ? 'white' : '#6c757d',
              border: 'none',
              borderRadius: '4px',
              cursor: inputValue.trim() ? 'pointer' : 'not-allowed',
              fontSize: '14px',
              fontWeight: '500',
              transition: 'all 0.2s ease'
            }}
          >
            追加
          </button>
        </div>
        <small className="form-help">新しい項目を入力してEnterまたは「追加」ボタンで追加</small>
      </div>

      {/* 選択済みアイテムの表示 */}
      {selectedItems.length > 0 && (
        <div style={{ marginBottom: '12px' }}>
          <div style={{ fontSize: '14px', fontWeight: '500', marginBottom: '6px', color: '#333' }}>
            選択中:
          </div>
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: '6px' }}>
            {selectedItems.map((item, index) => (
              <span
                key={index}
                style={{
                  padding: '4px 8px',
                  backgroundColor: '#e6f3ff',
                  border: '1px solid #0066cc',
                  borderRadius: '16px',
                  fontSize: '13px',
                  color: '#0066cc',
                  display: 'flex',
                  alignItems: 'center',
                  gap: '4px'
                }}
              >
                {item}
                <button
                  type="button"
                  onClick={() => handleToggleItem(item)}
                  style={{
                    background: 'none',
                    border: 'none',
                    color: '#0066cc',
                    cursor: 'pointer',
                    padding: '0',
                    fontSize: '16px',
                    lineHeight: '1'
                  }}
                >
                  ×
                </button>
              </span>
            ))}
          </div>
        </div>
      )}

      {/* 候補一覧（常時表示） */}
      {suggestions.length > 0 && (
        <div>
          <div style={{ fontSize: '14px', fontWeight: '500', marginBottom: '8px', color: '#333' }}>
            候補から選択:
          </div>
          <div style={{ 
            display: 'flex', 
            flexWrap: 'wrap', 
            gap: '6px',
            padding: '12px',
            backgroundColor: '#f8f9fa',
            border: '1px solid #e9ecef',
            borderRadius: '4px',
            maxHeight: '200px',
            overflowY: 'auto'
          }}>
            {suggestions.map((suggestion, index) => {
              const isSelected = selectedItems.includes(suggestion)
              const isNewItem = !originalSuggestions.includes(suggestion) // 新しく追加されたアイテムかチェック
              return (
                <button
                  key={index}
                  type="button"
                  onClick={() => handleToggleItem(suggestion)}
                  style={{
                    padding: '6px 12px',
                    border: '1px solid',
                    borderColor: isSelected ? '#0066cc' : '#dee2e6',
                    borderRadius: '16px',
                    backgroundColor: isSelected ? '#e6f3ff' : 'white',
                    color: isSelected ? '#0066cc' : '#495057',
                    cursor: 'pointer',
                    fontSize: '13px',
                    transition: 'all 0.2s ease',
                    fontWeight: isSelected ? '500' : 'normal',
                    position: 'relative'
                  }}
                  onMouseEnter={(e) => {
                    if (!isSelected) {
                      e.currentTarget.style.backgroundColor = '#f8f9fa'
                      e.currentTarget.style.borderColor = '#adb5bd'
                    }
                  }}
                  onMouseLeave={(e) => {
                    if (!isSelected) {
                      e.currentTarget.style.backgroundColor = 'white'
                      e.currentTarget.style.borderColor = '#dee2e6'
                    }
                  }}
                >
                  {suggestion}
                  {isNewItem && (
                    <span style={{
                      marginLeft: '4px',
                      fontSize: '10px',
                      backgroundColor: '#28a745',
                      color: 'white',
                      padding: '1px 4px',
                      borderRadius: '8px',
                      fontWeight: 'bold'
                    }}>
                      新
                    </span>
                  )}
                </button>
              )
            })}
          </div>
        </div>
      )}
      
      <small className="form-help" style={{ marginTop: '8px', display: 'block' }}>
        候補をクリックして選択/解除できます。新しい項目は上のテキストボックスから追加できます。
      </small>
    </div>
  )
}

export default TagCategoryInput