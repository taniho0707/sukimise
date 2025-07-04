import { describe, it, expect } from 'vitest'
import {
  isValidEmail,
  isValidUrl,
  isValidLatitude,
  isValidLongitude,
  isValidRating,
  isValidFileSize,
  isValidImageFile,
  isValidLength,
  isRequired,
} from './validation'

describe('validation', () => {
  describe('isValidEmail', () => {
    it('should validate correct email formats', () => {
      expect(isValidEmail('test@example.com')).toBe(true)
      expect(isValidEmail('user.name@domain.co.jp')).toBe(true)
      expect(isValidEmail('user+tag@example.org')).toBe(true)
      expect(isValidEmail('123@example.com')).toBe(true)
    })

    it('should reject invalid email formats', () => {
      expect(isValidEmail('invalid')).toBe(false)
      expect(isValidEmail('invalid@')).toBe(false)
      expect(isValidEmail('@domain.com')).toBe(false)
      expect(isValidEmail('user@')).toBe(false)
      expect(isValidEmail('user@domain')).toBe(false)
      expect(isValidEmail('')).toBe(false)
    })
  })

  describe('isValidUrl', () => {
    it('should validate correct URL formats', () => {
      expect(isValidUrl('https://example.com')).toBe(true)
      expect(isValidUrl('http://example.com')).toBe(true)
      expect(isValidUrl('https://www.example.com/path?query=value')).toBe(true)
      expect(isValidUrl('https://subdomain.example.co.jp')).toBe(true)
    })

    it('should return true for empty URL', () => {
      expect(isValidUrl('')).toBe(true)
    })

    it('should reject invalid URL formats', () => {
      expect(isValidUrl('not-a-url')).toBe(false)
      expect(isValidUrl('http://')).toBe(false)
      expect(isValidUrl('https://')).toBe(false)
      expect(isValidUrl('ftp://example.com')).toBe(false)
    })
  })

  describe('isValidLatitude', () => {
    it('should validate correct latitude values', () => {
      expect(isValidLatitude(0)).toBe(true)
      expect(isValidLatitude(90)).toBe(true)
      expect(isValidLatitude(-90)).toBe(true)
      expect(isValidLatitude(35.6762)).toBe(true) // Tokyo
      expect(isValidLatitude(-33.8688)).toBe(true) // Sydney
    })

    it('should reject invalid latitude values', () => {
      expect(isValidLatitude(91)).toBe(false)
      expect(isValidLatitude(-91)).toBe(false)
      expect(isValidLatitude(180)).toBe(false)
    })
  })

  describe('isValidLongitude', () => {
    it('should validate correct longitude values', () => {
      expect(isValidLongitude(0)).toBe(true)
      expect(isValidLongitude(180)).toBe(true)
      expect(isValidLongitude(-180)).toBe(true)
      expect(isValidLongitude(139.6503)).toBe(true) // Tokyo
      expect(isValidLongitude(151.2093)).toBe(true) // Sydney
    })

    it('should reject invalid longitude values', () => {
      expect(isValidLongitude(181)).toBe(false)
      expect(isValidLongitude(-181)).toBe(false)
      expect(isValidLongitude(360)).toBe(false)
    })
  })

  describe('isValidRating', () => {
    it('should validate correct rating values', () => {
      expect(isValidRating(1)).toBe(true)
      expect(isValidRating(2)).toBe(true)
      expect(isValidRating(3)).toBe(true)
      expect(isValidRating(4)).toBe(true)
      expect(isValidRating(5)).toBe(true)
    })

    it('should reject invalid rating values', () => {
      expect(isValidRating(0)).toBe(false)
      expect(isValidRating(6)).toBe(false)
      expect(isValidRating(-1)).toBe(false)
      expect(isValidRating(3.5)).toBe(false) // Not an integer
      expect(isValidRating(NaN)).toBe(false)
    })
  })

  describe('isValidFileSize', () => {
    it('should validate files within size limit', () => {
      const smallFile = new File(['content'], 'test.jpg', { type: 'image/jpeg' })
      Object.defineProperty(smallFile, 'size', { value: 1024 * 1024 }) // 1MB
      
      expect(isValidFileSize(smallFile, 5)).toBe(true)
      expect(isValidFileSize(smallFile, 1)).toBe(true)
    })

    it('should reject files exceeding size limit', () => {
      const largeFile = new File(['content'], 'test.jpg', { type: 'image/jpeg' })
      Object.defineProperty(largeFile, 'size', { value: 6 * 1024 * 1024 }) // 6MB
      
      expect(isValidFileSize(largeFile, 5)).toBe(false)
    })

    it('should use default 5MB limit when not specified', () => {
      const file = new File(['content'], 'test.jpg', { type: 'image/jpeg' })
      Object.defineProperty(file, 'size', { value: 4 * 1024 * 1024 }) // 4MB
      
      expect(isValidFileSize(file)).toBe(true)
    })
  })

  describe('isValidImageFile', () => {
    it('should validate allowed image types', () => {
      const jpegFile = new File(['content'], 'test.jpg', { type: 'image/jpeg' })
      const pngFile = new File(['content'], 'test.png', { type: 'image/png' })
      const gifFile = new File(['content'], 'test.gif', { type: 'image/gif' })
      const webpFile = new File(['content'], 'test.webp', { type: 'image/webp' })
      
      expect(isValidImageFile(jpegFile)).toBe(true)
      expect(isValidImageFile(pngFile)).toBe(true)
      expect(isValidImageFile(gifFile)).toBe(true)
      expect(isValidImageFile(webpFile)).toBe(true)
    })

    it('should reject non-image file types', () => {
      const textFile = new File(['content'], 'test.txt', { type: 'text/plain' })
      const pdfFile = new File(['content'], 'test.pdf', { type: 'application/pdf' })
      const bmpFile = new File(['content'], 'test.bmp', { type: 'image/bmp' })
      
      expect(isValidImageFile(textFile)).toBe(false)
      expect(isValidImageFile(pdfFile)).toBe(false)
      expect(isValidImageFile(bmpFile)).toBe(false)
    })
  })

  describe('isValidLength', () => {
    it('should validate strings within length limits', () => {
      expect(isValidLength('hello', 3, 10)).toBe(true)
      expect(isValidLength('hello', 5, 5)).toBe(true)
      expect(isValidLength('', 0, 10)).toBe(true)
    })

    it('should reject strings outside length limits', () => {
      expect(isValidLength('hi', 3, 10)).toBe(false)
      expect(isValidLength('very long string', 3, 10)).toBe(false)
    })

    it('should use default limits when not specified', () => {
      expect(isValidLength('any string')).toBe(true)
      expect(isValidLength('', 1)).toBe(false)
    })
  })

  describe('isRequired', () => {
    it('should validate non-empty required values', () => {
      expect(isRequired('hello')).toBe(true)
      expect(isRequired('  text  ')).toBe(true)
      expect(isRequired(123)).toBe(true)
      expect(isRequired(0)).toBe(true)
      expect(isRequired(false)).toBe(true)
      expect(isRequired([])).toBe(true)
      expect(isRequired({})).toBe(true)
    })

    it('should reject empty required values', () => {
      expect(isRequired('')).toBe(false)
      expect(isRequired('   ')).toBe(false)
      expect(isRequired(null)).toBe(false)
      expect(isRequired(undefined)).toBe(false)
    })
  })
})