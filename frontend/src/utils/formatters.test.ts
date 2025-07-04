import { describe, it, expect } from 'vitest'
import {
  formatDate,
  renderStars,
  formatCurrency,
  getImageUrl,
  safeTrim,
  safeJoin,
} from './formatters'

describe('formatters', () => {
  describe('formatDate', () => {
    it('should format valid date string', () => {
      const result = formatDate('2023-12-25T10:30:00Z')
      expect(result).toBe('2023/12/25')
    })

    it('should return null for null input', () => {
      const result = formatDate(null)
      expect(result).toBeNull()
    })

    it('should return null for empty string', () => {
      const result = formatDate('')
      expect(result).toBeNull()
    })
  })

  describe('renderStars', () => {
    it('should render 5 filled stars for rating 5', () => {
      const result = renderStars(5)
      expect(result).toBe('★★★★★')
    })

    it('should render 3 filled and 2 empty stars for rating 3', () => {
      const result = renderStars(3)
      expect(result).toBe('★★★☆☆')
    })

    it('should render 0 filled and 5 empty stars for rating 0', () => {
      const result = renderStars(0)
      expect(result).toBe('☆☆☆☆☆')
    })

    it('should render 1 filled and 4 empty stars for rating 1', () => {
      const result = renderStars(1)
      expect(result).toBe('★☆☆☆☆')
    })
  })

  describe('formatCurrency', () => {
    it('should format positive amount with yen symbol', () => {
      const result = formatCurrency(1000)
      expect(result).toBe('¥1,000')
    })

    it('should format large amount with commas', () => {
      const result = formatCurrency(1234567)
      expect(result).toBe('¥1,234,567')
    })

    it('should return empty string for null', () => {
      const result = formatCurrency(null)
      expect(result).toBe('')
    })

    it('should handle zero amount', () => {
      const result = formatCurrency(0)
      expect(result).toBe('')
    })
  })

  describe('getImageUrl', () => {
    it('should generate correct image URL', () => {
      const result = getImageUrl('image.jpg')
      expect(result).toBe('/api/v1/uploads/image.jpg')
    })

    it('should handle filename with special characters', () => {
      const result = getImageUrl('my-image_123.png')
      expect(result).toBe('/api/v1/uploads/my-image_123.png')
    })
  })

  describe('safeTrim', () => {
    it('should trim whitespace from string', () => {
      const result = safeTrim('  hello world  ')
      expect(result).toBe('hello world')
    })

    it('should return empty string for null', () => {
      const result = safeTrim(null)
      expect(result).toBe('')
    })

    it('should return empty string for undefined', () => {
      const result = safeTrim(undefined)
      expect(result).toBe('')
    })

    it('should handle empty string', () => {
      const result = safeTrim('')
      expect(result).toBe('')
    })

    it('should handle string with only whitespace', () => {
      const result = safeTrim('   ')
      expect(result).toBe('')
    })
  })

  describe('safeJoin', () => {
    it('should join array with default separator', () => {
      const result = safeJoin(['apple', 'banana', 'cherry'])
      expect(result).toBe('apple, banana, cherry')
    })

    it('should join array with custom separator', () => {
      const result = safeJoin(['apple', 'banana', 'cherry'], ' | ')
      expect(result).toBe('apple | banana | cherry')
    })

    it('should return empty string for null array', () => {
      const result = safeJoin(null)
      expect(result).toBe('')
    })

    it('should return empty string for undefined array', () => {
      const result = safeJoin(undefined)
      expect(result).toBe('')
    })

    it('should handle empty array', () => {
      const result = safeJoin([])
      expect(result).toBe('')
    })

    it('should handle single item array', () => {
      const result = safeJoin(['apple'])
      expect(result).toBe('apple')
    })
  })
})