import { useState, useCallback } from 'react'
import axios, { AxiosResponse } from 'axios'
import toast from 'react-hot-toast'
import { APIResponse } from '@/types/store'

// Custom hook for API calls with error handling and loading state
export const useApi = () => {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const request = useCallback(async <T>(
    apiCall: () => Promise<AxiosResponse<APIResponse<T>>>
  ): Promise<T | null> => {
    try {
      setLoading(true)
      setError(null)
      
      const response = await apiCall()
      
      // Handle API response format
      if (response.data.success && response.data.data) {
        return response.data.data
      } else if (response.data.error) {
        const errorMessage = response.data.error.message || 'Unknown error'
        setError(errorMessage)
        toast.error(errorMessage)
        return null
      } else {
        // Handle legacy response format (direct data)
        return response.data as unknown as T
      }
    } catch (err) {
      let errorMessage = 'Unknown error occurred'
      
      if (axios.isAxiosError(err)) {
        if (err.response?.data?.error) {
          errorMessage = err.response.data.error.message || err.response.data.error
        } else if (err.response?.data?.message) {
          errorMessage = err.response.data.message
        } else if (err.message) {
          errorMessage = err.message
        }
      } else if (err instanceof Error) {
        errorMessage = err.message
      }
      
      setError(errorMessage)
      toast.error(errorMessage)
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const clearError = useCallback(() => {
    setError(null)
  }, [])

  return { request, loading, error, clearError }
}

// Specialized hooks for common API operations
export const useStores = () => {
  const { request, loading, error } = useApi()

  const fetchStores = useCallback(async (filters?: Record<string, any>) => {
    const params = new URLSearchParams()
    
    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          if (Array.isArray(value)) {
            params.append(key, value.join(','))
          } else {
            params.append(key, String(value))
          }
        }
      })
    }

    const url = `/api/v1/stores${params.toString() ? `?${params.toString()}` : ''}`
    
    return request(() => axios.get(url))
  }, [request])

  const fetchStore = useCallback(async (id: string) => {
    return request(() => axios.get(`/api/v1/stores/${id}`))
  }, [request])

  const createStore = useCallback(async (storeData: any) => {
    return request(() => axios.post('/api/v1/stores', storeData))
  }, [request])

  const updateStore = useCallback(async (id: string, storeData: any) => {
    return request(() => axios.put(`/api/v1/stores/${id}`, storeData))
  }, [request])

  const deleteStore = useCallback(async (id: string) => {
    return request(() => axios.delete(`/api/v1/stores/${id}`))
  }, [request])

  return {
    fetchStores,
    fetchStore,
    createStore,
    updateStore,
    deleteStore,
    loading,
    error,
  }
}

export const useReviews = () => {
  const { request, loading, error } = useApi()

  const fetchReviews = useCallback(async (storeId: string) => {
    return request(() => axios.get(`/api/v1/stores/${storeId}/reviews`))
  }, [request])

  const fetchUserReviews = useCallback(async () => {
    return request(() => axios.get('/api/v1/users/me/reviews'))
  }, [request])

  const createReview = useCallback(async (reviewData: any) => {
    return request(() => axios.post('/api/v1/reviews', reviewData))
  }, [request])

  const updateReview = useCallback(async (id: string, reviewData: any) => {
    return request(() => axios.put(`/api/v1/reviews/${id}`, reviewData))
  }, [request])

  const deleteReview = useCallback(async (id: string) => {
    return request(() => axios.delete(`/api/v1/reviews/${id}`))
  }, [request])

  return {
    fetchReviews,
    fetchUserReviews,
    createReview,
    updateReview,
    deleteReview,
    loading,
    error,
  }
}

export const useCategories = () => {
  const { request, loading, error } = useApi()

  const fetchCategories = useCallback(async () => {
    return request(() => axios.get('/api/v1/stores/categories'))
  }, [request])

  return { fetchCategories, loading, error }
}

export const useTags = () => {
  const { request, loading, error } = useApi()

  const fetchTags = useCallback(async () => {
    return request(() => axios.get('/api/v1/stores/tags'))
  }, [request])

  return { fetchTags, loading, error }
}

