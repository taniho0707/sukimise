import React, { useState } from 'react'
import { getImageUrl } from '../utils/imageUrl'
import './ImageGallery.css'

interface ImageGalleryProps {
  images: string[]
  alt?: string
  className?: string
}

const ImageGallery: React.FC<ImageGalleryProps> = ({
  images,
  alt = '画像',
  className = ''
}) => {
  const [selectedImage, setSelectedImage] = useState<string | null>(null)

  if (!images || images.length === 0) {
    return null
  }

  const openModal = (imageUrl: string) => {
    setSelectedImage(imageUrl)
  }

  const closeModal = () => {
    setSelectedImage(null)
  }

  const nextImage = () => {
    if (!selectedImage) return
    const currentIndex = images.indexOf(selectedImage)
    const nextIndex = (currentIndex + 1) % images.length
    setSelectedImage(images[nextIndex])
  }

  const prevImage = () => {
    if (!selectedImage) return
    const currentIndex = images.indexOf(selectedImage)
    const prevIndex = (currentIndex - 1 + images.length) % images.length
    setSelectedImage(images[prevIndex])
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Escape') {
      closeModal()
    } else if (e.key === 'ArrowRight') {
      nextImage()
    } else if (e.key === 'ArrowLeft') {
      prevImage()
    }
  }

  return (
    <>
      <div className={`image-gallery ${className}`}>
        {images.length === 1 ? (
          <div className="single-image">
            <img
              src={getImageUrl(images[0])}
              alt={alt}
              className="gallery-image"
              onClick={() => openModal(images[0])}
            />
          </div>
        ) : (
          <div className="image-grid">
            {images.slice(0, 4).map((imageUrl, index) => (
              <div
                key={index}
                className={`image-item ${index === 0 ? 'main-image' : ''}`}
                onClick={() => openModal(imageUrl)}
              >
                <img
                  src={getImageUrl(imageUrl)}
                  alt={`${alt} ${index + 1}`}
                  className="gallery-image"
                />
                {index === 3 && images.length > 4 && (
                  <div className="more-images-overlay">
                    <span>+{images.length - 4}</span>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>

      {/* モーダル */}
      {selectedImage && (
        <div
          className="image-modal"
          onClick={closeModal}
          onKeyDown={handleKeyDown}
          tabIndex={0}
        >
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <button className="modal-close" onClick={closeModal}>
              ×
            </button>
            
            {images.length > 1 && (
              <>
                <button className="modal-nav prev" onClick={prevImage}>
                  ‹
                </button>
                <button className="modal-nav next" onClick={nextImage}>
                  ›
                </button>
              </>
            )}
            
            <img
              src={getImageUrl(selectedImage)}
              alt={alt}
              className="modal-image"
            />
            
            {images.length > 1 && (
              <div className="image-counter">
                {images.indexOf(selectedImage) + 1} / {images.length}
              </div>
            )}
          </div>
        </div>
      )}
    </>
  )
}

export default ImageGallery