import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import axios from 'axios'
import toast from 'react-hot-toast'
import { useAuth } from '../contexts/AuthContext'
import { BusinessHoursData } from '../types/store'
import ImageUpload from '../components/ImageUpload'
import BusinessHoursInput from '../components/BusinessHoursInput'
import TagCategoryInput from '../components/TagCategoryInput'
import ParkingInput from '../components/ParkingInput'
import './StoreForm.css'
import './StoreList.css'

const storeSchema = z.object({
  name: z.string().min(1, '店舗名は必須です'),
  address: z.string().min(1, '住所は必須です'),
  latitude: z.preprocess((val) => {
    if (val === '' || val === null || val === undefined || val === 'NaN') return 0;
    const num = parseFloat(val as string);
    return isNaN(num) ? 0 : num;
  }, z.number().min(-90).max(90)),
  longitude: z.preprocess((val) => {
    if (val === '' || val === null || val === undefined || val === 'NaN') return 0;
    const num = parseFloat(val as string);
    return isNaN(num) ? 0 : num;
  }, z.number().min(-180).max(180)),
  categories: z.string(),
  parkingInfo: z.string(),
  websiteUrl: z.string().url('正しいURLを入力してください').or(z.literal('')),
  googleMapUrl: z.string().url('正しいURLを入力してください').or(z.literal('')),
  snsUrls: z.string(),
  tags: z.string(),
})

type StoreFormData = z.infer<typeof storeSchema>

interface Store {
  id?: string
  name: string
  address: string
  latitude: number
  longitude: number
  categories: string[]
  business_hours: BusinessHoursData
  parking_info: string
  website_url: string
  google_map_url: string
  sns_urls: string[]
  tags: string[]
  photos: string[]
}

const getDefaultBusinessHours = (): BusinessHoursData => ({
  monday: { is_closed: false, time_slots: [] },
  tuesday: { is_closed: false, time_slots: [] },
  wednesday: { is_closed: false, time_slots: [] },
  thursday: { is_closed: false, time_slots: [] },
  friday: { is_closed: false, time_slots: [] },
  saturday: { is_closed: false, time_slots: [] },
  sunday: { is_closed: false, time_slots: [] },
})

const StoreForm: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { user } = useAuth()
  const isEdit = !!id
  const [loading, setLoading] = useState(false)
  const [store, setStore] = useState<Store | null>(null)
  const [images, setImages] = useState<string[]>([])
  const [businessHours, setBusinessHours] = useState<BusinessHoursData>(getDefaultBusinessHours())

  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue,
    watch
  } = useForm<StoreFormData>({
    resolver: zodResolver(storeSchema),
    defaultValues: {
      name: '',
      address: '',
      latitude: 0,
      longitude: 0,
      categories: '',
      parkingInfo: '',
      websiteUrl: '',
      googleMapUrl: '',
      snsUrls: '',
      tags: '',
    }
  })

  // 編集時の既存データ取得
  useEffect(() => {
    if (isEdit && id) {
      const fetchStore = async () => {
        try {
          const response = await axios.get(`/api/v1/stores/${id}`)
          const storeData = response.data.data
          setStore(storeData)
          
          // 編集権限の確認
          if (user && storeData.created_by !== user.id && user.role !== 'admin') {
            toast.error('この店舗を編集する権限がありません')
            navigate('/stores')
            return
          }
          
          // フォームに既存データを設定
          setValue('name', storeData.name)
          setValue('address', storeData.address)
          setValue('latitude', storeData.latitude.toString())
          setValue('longitude', storeData.longitude.toString())
          setValue('categories', storeData.categories?.join(', ') || '')
          setValue('parkingInfo', storeData.parking_info || '')
          setValue('websiteUrl', storeData.website_url || '')
          setValue('googleMapUrl', storeData.google_map_url || '')
          setValue('snsUrls', storeData.sns_urls?.join(', ') || '')
          setValue('tags', storeData.tags?.join(', ') || '')
          setImages(storeData.photos || [])
          
          // BusinessHoursDataを直接設定
          if (storeData.business_hours) {
            setBusinessHours(storeData.business_hours)
          }
        } catch (error: any) {
          console.error('Store fetch error:', error)
          if (error.response?.status === 404) {
            toast.error('店舗が見つかりません')
          } else if (error.response?.status === 403) {
            toast.error('この店舗を編集する権限がありません')
          } else {
            toast.error('店舗データの取得に失敗しました')
          }
          navigate('/stores')
        }
      }
      fetchStore()
    }
  }, [isEdit, id, setValue, navigate, user])

  const onSubmit = async (data: StoreFormData) => {
    if (!user) {
      toast.error('ログインが必要です')
      navigate('/login')
      return
    }

    setLoading(true)
    try {
      const storeData = {
        name: data.name,
        address: data.address,
        latitude: data.latitude,
        longitude: data.longitude,
        categories: data.categories.split(',').map(c => c.trim()).filter(c => c),
        business_hours: businessHours,
        parking_info: data.parkingInfo,
        website_url: data.websiteUrl,
        google_map_url: data.googleMapUrl,
        sns_urls: data.snsUrls.split(',').map(u => u.trim()).filter(u => u),
        tags: data.tags.split(',').map(t => t.trim()).filter(t => t),
        photos: images,
      }

      if (isEdit && id) {
        await axios.put(`/api/v1/stores/${id}`, storeData)
        toast.success('店舗情報を更新しました')
      } else {
        await axios.post('/api/v1/stores', storeData)
        toast.success('店舗を登録しました')
      }
      
      navigate('/stores')
    } catch (error: any) {
      console.error('Store save error:', error)
      if (error.response?.data?.error) {
        toast.error(error.response.data.error)
      } else {
        toast.error(isEdit ? '店舗の更新に失敗しました' : '店舗の登録に失敗しました')
      }
    } finally {
      setLoading(false)
    }
  }

  const handleCancel = () => {
    navigate('/stores')
  }

  const getCurrentLocation = () => {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          setValue('latitude', position.coords.latitude)
          setValue('longitude', position.coords.longitude)
          toast.success('現在位置を取得しました')
        },
        () => {
          toast.error('位置情報の取得に失敗しました')
        }
      )
    } else {
      toast.error('この環境では位置情報を取得できません')
    }
  }

  const getLocationFromAddress = async () => {
    const address = watch('address')
    if (!address) {
      toast.error('住所を入力してください')
      return
    }

    try {
      // OpenStreetMap Nominatimを使用した簡単なジオコーディング
      const response = await fetch(`https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(address)}&limit=1`)
      const data = await response.json()
      
      if (data && data.length > 0) {
        const lat = parseFloat(data[0].lat)
        const lon = parseFloat(data[0].lon)
        setValue('latitude', lat)
        setValue('longitude', lon)
        toast.success('住所から位置情報を取得しました')
      } else {
        toast.error('住所から位置情報を取得できませんでした')
      }
    } catch (error) {
      console.error('Geocoding error:', error)
      toast.error('位置情報の取得に失敗しました')
    }
  }

  return (
    <div className="store-form-page">
      <div className="store-form-container">
        <div className="store-form-header">
          <h1>{isEdit ? '店舗編集' : '店舗登録'}</h1>
          {isEdit && store && (
            <p className="edit-note">店舗ID: {id}</p>
          )}
        </div>

        <form onSubmit={handleSubmit(onSubmit as any)} className="store-form">
          {/* 基本情報セクション */}
          <section className="form-section">
            <h2>基本情報</h2>
            
            <div className="form-group">
              <label htmlFor="name" className="form-label required">
                店舗名
              </label>
              <input
                type="text"
                id="name"
                {...register('name')}
                className={`form-input ${errors.name ? 'error' : ''}`}
                placeholder="例: すき家 新宿店"
              />
              {errors.name && (
                <span className="error-message">{errors.name.message}</span>
              )}
            </div>

            <div className="form-group">
              <label htmlFor="address" className="form-label required">
                住所
              </label>
              <div className="form-row">
                <input
                  type="text"
                  id="address"
                  {...register('address')}
                  className={`form-input ${errors.address ? 'error' : ''}`}
                  placeholder="例: 東京都新宿区新宿3-1-1"
                  style={{ flex: 1 }}
                />
                <button
                  type="button"
                  onClick={getLocationFromAddress}
                  className="btn btn-secondary"
                  style={{ marginLeft: '8px' }}
                >
                  位置を取得
                </button>
              </div>
              {errors.address && (
                <span className="error-message">{errors.address.message}</span>
              )}
            </div>

            <div className="form-row">
              <div className="form-group">
                <label htmlFor="latitude" className="form-label">
                  緯度（任意）
                </label>
                <input
                  type="number"
                  step="any"
                  id="latitude"
                  {...register('latitude')}
                  className={`form-input ${errors.latitude ? 'error' : ''}`}
                  placeholder="35.6762"
                />
                {errors.latitude && (
                  <span className="error-message">{errors.latitude.message}</span>
                )}
              </div>

              <div className="form-group">
                <label htmlFor="longitude" className="form-label">
                  経度（任意）
                </label>
                <input
                  type="number"
                  step="any"
                  id="longitude"
                  {...register('longitude')}
                  className={`form-input ${errors.longitude ? 'error' : ''}`}
                  placeholder="139.6503"
                />
                {errors.longitude && (
                  <span className="error-message">{errors.longitude.message}</span>
                )}
              </div>

              <div className="form-group">
                <button
                  type="button"
                  onClick={getCurrentLocation}
                  className="btn btn-secondary location-btn"
                >
                  現在位置を取得
                </button>
              </div>
            </div>

            <TagCategoryInput
              label="カテゴリ"
              value={watch('categories') || ''}
              onChange={(value) => setValue('categories', value)}
              apiEndpoint="/api/v1/stores/categories"
              placeholder="例: 和食, ラーメン, 定食"
            />
          </section>

          {/* 営業情報セクション */}
          <section className="form-section">
            <h2>営業情報</h2>

            <div className="form-group">
              <label className="form-label">
                営業時間
              </label>
              <BusinessHoursInput
                value={businessHours}
                onChange={setBusinessHours}
                className="business-hours-input"
              />
            </div>


            <div className="form-group">
              <label className="form-label">
                駐車場情報
              </label>
              <ParkingInput
                value={watch('parkingInfo') || ''}
                onChange={(value) => setValue('parkingInfo', value)}
                className="parking-input"
              />
            </div>
          </section>

          {/* Web情報セクション */}
          <section className="form-section">
            <h2>Web情報</h2>

            <div className="form-group">
              <label htmlFor="websiteUrl" className="form-label">
                ホームページURL
              </label>
              <input
                type="url"
                id="websiteUrl"
                {...register('websiteUrl')}
                className={`form-input ${errors.websiteUrl ? 'error' : ''}`}
                placeholder="https://example.com"
              />
              {errors.websiteUrl && (
                <span className="error-message">{errors.websiteUrl.message}</span>
              )}
            </div>

            <div className="form-group">
              <label htmlFor="googleMapUrl" className="form-label">
                GoogleマップURL
              </label>
              <input
                type="url"
                id="googleMapUrl"
                {...register('googleMapUrl')}
                className={`form-input ${errors.googleMapUrl ? 'error' : ''}`}
                placeholder="https://maps.google.com/..."
              />
              {errors.googleMapUrl && (
                <span className="error-message">{errors.googleMapUrl.message}</span>
              )}
            </div>

            <div className="form-group">
              <label htmlFor="snsUrls" className="form-label">
                SNS URL
              </label>
              <input
                type="text"
                id="snsUrls"
                {...register('snsUrls')}
                className="form-input"
                placeholder="例: https://twitter.com/..., https://instagram.com/..."
              />
              <small className="form-help">複数のSNS URLをカンマ区切りで入力できます</small>
            </div>
          </section>

          {/* 分類・検索セクション */}
          <section className="form-section">
            <h2>分類・検索</h2>

            <TagCategoryInput
              label="タグ"
              value={watch('tags') || ''}
              onChange={(value) => setValue('tags', value)}
              apiEndpoint="/api/v1/stores/tags"
              placeholder="例: デート向け, 家族連れ, 駅近, 深夜営業"
            />
          </section>

          {/* 画像セクション */}
          <section className="form-section">
            <h2>画像</h2>
            <ImageUpload
              images={images}
              onChange={setImages}
              maxImages={10}
              disabled={loading}
            />
          </section>

          {/* 送信ボタン */}
          <div className="form-actions">
            <button
              type="button"
              onClick={handleCancel}
              className="btn btn-secondary"
              disabled={loading}
            >
              キャンセル
            </button>
            <button
              type="submit"
              className="btn btn-primary"
              disabled={loading}
            >
              {loading ? (isEdit ? '更新中...' : '登録中...') : (isEdit ? '更新' : '登録')}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default StoreForm