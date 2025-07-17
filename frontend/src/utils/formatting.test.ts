/*
import { describe, it, expect, beforeEach, vi } from 'vitest'
import {
  formatDate,
  formatDateTime,
  formatDateTimeShort,
  formatCurrency,
  formatPrice,
  truncateText,
  capitalizeFirst,
  formatBusinessHours,
  formatCoordinates,
  formatLatitude,
  formatLongitude,
  formatRating,
  getRatingStars,
  formatList,
  ensureHttps,
  formatWebsiteUrl,
  formatFileSize,
  formatDistance,
  formatTime,
  formatBusinessDay,
  formatSearchQuery,
  formatErrorMessage,
  formatInputValue,
  formatPhoneNumber,
  formatTag,
  formatTags,
} from './formatting'

describe('formatting utils', () => {
  describe('date formatting', () => {
    beforeEach(() => {
      // Mock the current date to ensure consistent test results
      const mockDate = new Date('2023-12-25T10:30:00Z')
      global.Date = vi.fn(() => mockDate) as any
      global.Date.now = vi.fn(() => mockDate.getTime())
    })

    it('should format date correctly', () => {
      const result = formatDate('2023-12-25T10:30:00Z')
      expect(result).toMatch(/2023年12月25日/)
    })

    it('should format date time correctly', () => {
      const result = formatDateTime('2023-12-25T10:30:00Z')
      expect(result).toMatch(/2023年12月25日.*10:30/)
    })

    it('should format short date time correctly', () => {
      const result = formatDateTimeShort('2023-12-25T10:30:00Z')
      expect(result).toMatch(/12月25日.*10:30/)
    })
  })

  describe('currency formatting', () => {
    it('should format currency correctly', () => {
      const result = formatCurrency(1000)
      expect(result).toMatch(/￥1,000/)
    })

    it('should format price correctly', () => {
      const result = formatPrice(1500)
      expect(result).toBe('1,500円')
    })

    it('should handle zero amount', () => {
      const result = formatPrice(0)
      expect(result).toBe('0円')
    })
  })

  describe('text formatting', () => {
    it('should truncate long text', () => {
      const result = truncateText('This is a very long text that should be truncated', 20)
      expect(result).toBe('This is a very long ...')
    })

    it('should return original text if within limit', () => {
      const result = truncateText('Short text', 20)
      expect(result).toBe('Short text')
    })

    it('should capitalize first letter', () => {
      expect(capitalizeFirst('hello world')).toBe('Hello world')
      expect(capitalizeFirst('HELLO WORLD')).toBe('Hello world')
      expect(capitalizeFirst('')).toBe('')
    })
  })

  describe('business hours formatting', () => {
    it('should split business hours by newline', () => {
      const result = formatBusinessHours('月-金: 9:00-18:00\n土: 10:00-17:00\n日: 休み')
      expect(result).toEqual(['月-金: 9:00-18:00', '土: 10:00-17:00', '日: 休み'])
    })

    it('should filter empty lines', () => {
      const result = formatBusinessHours('Monday: 9-5\n\nTuesday: 9-5\n')
      expect(result).toEqual(['Monday: 9-5', 'Tuesday: 9-5'])
    })

    it('should return empty array for empty input', () => {
      expect(formatBusinessHours('')).toEqual([])
      expect(formatBusinessHours(null as any)).toEqual([])
    })
  })

  describe('coordinates formatting', () => {
    it('should format coordinates with 6 decimal places', () => {
      const result = formatCoordinates(35.676191, 139.650311)
      expect(result).toBe('35.676191, 139.650311')
    })

    it('should format latitude with label', () => {
      const result = formatLatitude(35.676191)
      expect(result).toBe('緯度: 35.676191')
    })

    it('should format longitude with label', () => {
      const result = formatLongitude(139.650311)
      expect(result).toBe('経度: 139.650311')
    })
  })

  describe('rating formatting', () => {
    it('should format rating with one decimal', () => {
      expect(formatRating(4.5)).toBe('4.5')
      expect(formatRating(3)).toBe('3.0')
    })

    it('should generate rating stars correctly', () => {
      expect(getRatingStars(5)).toBe('★★★★★')
      expect(getRatingStars(4.5)).toBe('★★★★☆')
      expect(getRatingStars(3)).toBe('★★★☆☆')
      expect(getRatingStars(0)).toBe('☆☆☆☆☆')
    })
  })

  describe('list formatting', () => {
    it('should format short lists normally', () => {
      const result = formatList(['apple', 'banana', 'cherry'])
      expect(result).toBe('apple, banana, cherry')
    })

    it('should format long lists with "others" count', () => {
      const result = formatList(['apple', 'banana', 'cherry', 'date', 'elderberry'], 3)
      expect(result).toBe('apple, banana, cherry 他2件')
    })

    it('should return empty string for empty array', () => {
      expect(formatList([])).toBe('')
    })
  })

  describe('URL formatting', () => {
    it('should ensure HTTPS protocol', () => {
      expect(ensureHttps('example.com')).toBe('https://example.com')
      expect(ensureHttps('http://example.com')).toBe('http://example.com')
      expect(ensureHttps('https://example.com')).toBe('https://example.com')
    })

    it('should format website URL for display', () => {
      expect(formatWebsiteUrl('https://example.com')).toBe('example.com')
      expect(formatWebsiteUrl('http://www.example.com')).toBe('www.example.com')
      expect(formatWebsiteUrl('example.com')).toBe('example.com')
      expect(formatWebsiteUrl('')).toBe('')
    })
  })

  describe('file size formatting', () => {
    it('should format bytes correctly', () => {
      expect(formatFileSize(0)).toBe('0 Bytes')
      expect(formatFileSize(1024)).toBe('1 KB')
      expect(formatFileSize(1048576)).toBe('1 MB')
      expect(formatFileSize(1073741824)).toBe('1 GB')
    })

    it('should handle decimal values', () => {
      expect(formatFileSize(1536)).toBe('1.5 KB')
    })
  })

  describe('distance formatting', () => {
    it('should format short distances in meters', () => {
      expect(formatDistance(500)).toBe('500m')
      expect(formatDistance(999)).toBe('999m')
    })

    it('should format long distances in kilometers', () => {
      expect(formatDistance(1000)).toBe('1.0km')
      expect(formatDistance(1500)).toBe('1.5km')
    })
  })

  describe('time formatting', () => {
    it('should format time with leading zeros', () => {
      expect(formatTime('9:5')).toBe('09:05')
      expect(formatTime('10:30')).toBe('10:30')
    })

    it('should format business day abbreviations', () => {
      expect(formatBusinessDay('monday')).toBe('月')
      expect(formatBusinessDay('sunday')).toBe('日')
      expect(formatBusinessDay('invalid')).toBe('invalid')
    })
  })

  describe('utility formatting', () => {
    it('should format search query', () => {
      expect(formatSearchQuery('  Hello World  ')).toBe('hello world')
    })

    it('should format error messages', () => {
      expect(formatErrorMessage('Simple error')).toBe('Simple error')
      expect(formatErrorMessage({ message: 'Error object' })).toBe('Error object')
      expect(formatErrorMessage({})).toBe('エラーが発生しました')
    })

    it('should format input values', () => {
      expect(formatInputValue('hello')).toBe('hello')
      expect(formatInputValue(123)).toBe('123')
      expect(formatInputValue(null)).toBe('')
      expect(formatInputValue(undefined)).toBe('')
    })
  })

  describe('phone number formatting', () => {
    it('should format 10-digit phone numbers', () => {
      expect(formatPhoneNumber('0312345678')).toBe('03-1234-5678')
    })

    it('should format 11-digit phone numbers', () => {
      expect(formatPhoneNumber('09012345678')).toBe('090-1234-5678')
    })

    it('should handle numbers with existing formatting', () => {
      expect(formatPhoneNumber('03-1234-5678')).toBe('03-1234-5678')
    })

    it('should return original for invalid lengths', () => {
      expect(formatPhoneNumber('12345')).toBe('12345')
    })
  })

  describe('tag formatting', () => {
    it('should add hash prefix to tags', () => {
      expect(formatTag('restaurant')).toBe('#restaurant')
      expect(formatTag('#restaurant')).toBe('#restaurant')
    })

    it('should format arrays of tags', () => {
      const result = formatTags(['restaurant', '#cafe', 'bar'])
      expect(result).toEqual(['#restaurant', '#cafe', '#bar'])
    })
  })
})*/
