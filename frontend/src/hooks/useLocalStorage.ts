import { useState, useCallback, useEffect } from 'react'

type SetValue<T> = T | ((val: T) => T)

// Custom hook for localStorage with TypeScript support
export const useLocalStorage = <T>(key: string, initialValue: T) => {
  // State to store our value
  const [storedValue, setStoredValue] = useState<T>(() => {
    if (typeof window === 'undefined') {
      return initialValue
    }
    
    try {
      const item = window.localStorage.getItem(key)
      return item ? JSON.parse(item) : initialValue
    } catch (error) {
      console.error(`Error reading localStorage key "${key}":`, error)
      return initialValue
    }
  })

  // Return a wrapped version of useState's setter function that persists the new value to localStorage
  const setValue = useCallback((value: SetValue<T>) => {
    try {
      // Allow value to be a function so we have the same API as useState
      const valueToStore = value instanceof Function ? value(storedValue) : value
      
      // Save state
      setStoredValue(valueToStore)
      
      // Save to local storage
      if (typeof window !== 'undefined') {
        window.localStorage.setItem(key, JSON.stringify(valueToStore))
      }
    } catch (error) {
      console.error(`Error setting localStorage key "${key}":`, error)
    }
  }, [key, storedValue])

  // Remove item from localStorage
  const removeValue = useCallback(() => {
    try {
      setStoredValue(initialValue)
      if (typeof window !== 'undefined') {
        window.localStorage.removeItem(key)
      }
    } catch (error) {
      console.error(`Error removing localStorage key "${key}":`, error)
    }
  }, [key, initialValue])

  // Listen for changes to the localStorage key in other tabs
  useEffect(() => {
    if (typeof window === 'undefined') {
      return
    }

    const handleStorageChange = (e: StorageEvent) => {
      if (e.key === key && e.newValue !== null) {
        try {
          setStoredValue(JSON.parse(e.newValue))
        } catch (error) {
          console.error(`Error parsing localStorage value for key "${key}":`, error)
        }
      }
    }

    window.addEventListener('storage', handleStorageChange)
    
    return () => {
      window.removeEventListener('storage', handleStorageChange)
    }
  }, [key])

  return [storedValue, setValue, removeValue] as const
}

// Specialized hooks for common localStorage use cases
export const useAuthToken = () => {
  return useLocalStorage<string | null>('authToken', null)
}

export const useUserPreferences = () => {
  return useLocalStorage('userPreferences', {
    theme: 'light',
    language: 'ja',
    mapZoom: 13,
    mapCenter: { lat: 35.6762, lng: 139.6503 }, // Tokyo
  })
}

export const useRecentSearches = () => {
  return useLocalStorage<string[]>('recentSearches', [])
}

export const useFavoriteStores = () => {
  return useLocalStorage<string[]>('favoriteStores', [])
}