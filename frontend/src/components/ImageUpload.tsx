import React, { useState, useRef } from 'react'
import axios from 'axios'
import toast from 'react-hot-toast'
import { getImageUrl } from '../utils/imageUrl'
import './ImageUpload.css'

interface UploadedImage {
  filename: string
  url: string
  size: number
}

interface ImageUploadProps {
  images: string[]
  onChange: (images: string[]) => void
  maxImages?: number
  disabled?: boolean
}

const ImageUpload: React.FC<ImageUploadProps> = ({
  images,
  onChange,
  maxImages = 5,
  disabled = false
}) => {
  const [uploading, setUploading] = useState(false)
  const [dragOver, setDragOver] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFileSelect = async (files: FileList | null) => {
    if (!files || files.length === 0) return
    
    const remainingSlots = maxImages - images.length
    if (remainingSlots <= 0) {
      toast.error(`最大${maxImages}枚まで追加できます`)
      return
    }

    const filesToUpload = Array.from(files).slice(0, remainingSlots)
    
    // ファイルサイズとタイプの検証
    for (const file of filesToUpload) {
      if (file.size > 10 * 1024 * 1024) {
        toast.error(`${file.name}は10MBを超えています`)
        return
      }
      
      if (!file.type.startsWith('image/')) {
        toast.error(`${file.name}は画像ファイルではありません`)
        return
      }
    }

    setUploading(true)

    try {
      const uploadPromises = filesToUpload.map(async (file) => {
        const formData = new FormData()
        formData.append('image', file)

        const response = await axios.post<UploadedImage>('/api/v1/upload/image', formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        })

        return response.data.url
      })

      const uploadedUrls = await Promise.all(uploadPromises)
      const newImages = [...images, ...uploadedUrls]
      onChange(newImages)
      
      toast.success(`${uploadedUrls.length}枚の画像をアップロードしました`)
    } catch (error: any) {
      console.error('Upload error:', error)
      if (error.response?.data?.error) {
        toast.error(error.response.data.error)
      } else {
        toast.error('画像のアップロードに失敗しました')
      }
    } finally {
      setUploading(false)
      // ファイル入力をリセット
      if (fileInputRef.current) {
        fileInputRef.current.value = ''
      }
    }
  }

  const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    handleFileSelect(e.target.files)
  }

  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    setDragOver(false)
    
    if (disabled) return
    
    handleFileSelect(e.dataTransfer.files)
  }

  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    setDragOver(true)
  }

  const handleDragLeave = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    setDragOver(false)
  }

  const handleRemoveImage = async (index: number) => {
    const imageUrl = images[index]
    
    // サーバーから画像を削除
    try {
      const filename = imageUrl.split('/').pop()
      if (filename) {
        await axios.delete(`/api/v1/upload/${filename}`)
      }
    } catch (error) {
      console.error('Delete error:', error)
      // 削除に失敗してもUI上は削除する
    }

    const newImages = images.filter((_, i) => i !== index)
    onChange(newImages)
    toast.success('画像を削除しました')
  }

  const openFileDialog = () => {
    if (!disabled && fileInputRef.current) {
      fileInputRef.current.click()
    }
  }

  return (
    <div className="image-upload">
      <div className="image-upload-header">
        <label className="form-label">画像</label>
        <span className="image-count">
          {images.length}/{maxImages}
        </span>
      </div>

      {/* 画像プレビュー */}
      {images.length > 0 && (
        <div className="image-preview-grid">
          {images.map((imageUrl, index) => (
            <div key={index} className="image-preview-item">
              <img
                src={getImageUrl(imageUrl)}
                alt={`アップロード画像 ${index + 1}`}
                className="preview-image"
              />
              <button
                type="button"
                onClick={() => handleRemoveImage(index)}
                className="remove-image-btn"
                disabled={disabled}
              >
                ×
              </button>
            </div>
          ))}
        </div>
      )}

      {/* アップロードエリア */}
      {images.length < maxImages && (
        <div
          className={`upload-area ${dragOver ? 'drag-over' : ''} ${disabled ? 'disabled' : ''}`}
          onDrop={handleDrop}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onClick={openFileDialog}
        >
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            multiple
            onChange={handleFileInputChange}
            className="file-input"
            disabled={disabled}
          />
          
          {uploading ? (
            <div className="upload-loading">
              <div className="spinner"></div>
              <p>アップロード中...</p>
            </div>
          ) : (
            <div className="upload-content">
              <div className="upload-icon">📷</div>
              <p className="upload-text">
                クリックまたはドラッグ&ドロップで画像を追加
              </p>
              <p className="upload-hint">
                JPEG, PNG, GIF, WebP（最大10MB）
              </p>
            </div>
          )}
        </div>
      )}
    </div>
  )
}

export default ImageUpload