import React, { useState, useEffect } from 'react'
import { BusinessHoursData, DaySchedule, TimeSlot } from '../types/store'

interface BusinessHoursInputProps {
  value: BusinessHoursData
  onChange: (value: BusinessHoursData) => void
  className?: string
}

const DAYS_OF_WEEK = [
  { key: 'monday', label: '月', fullLabel: '月曜日' },
  { key: 'tuesday', label: '火', fullLabel: '火曜日' },
  { key: 'wednesday', label: '水', fullLabel: '水曜日' },
  { key: 'thursday', label: '木', fullLabel: '木曜日' },
  { key: 'friday', label: '金', fullLabel: '金曜日' },
  { key: 'saturday', label: '土', fullLabel: '土曜日' },
  { key: 'sunday', label: '日', fullLabel: '日曜日' },
]

// 30分単位の時間オプションを生成
const generateTimeOptions = () => {
  const options: Array<{ value: string; label: string }> = []
  
  for (let hour = 0; hour < 24; hour++) {
    for (let minute of [0, 30]) {
      const timeString = `${hour.toString().padStart(2, '0')}:${minute.toString().padStart(2, '0')}`
      const displayTime = `${hour}:${minute.toString().padStart(2, '0')}`
      options.push({ value: timeString, label: displayTime })
    }
  }
  
  return options
}

const TIME_OPTIONS = generateTimeOptions()

const getDefaultBusinessHours = (): BusinessHoursData => ({
  monday: { is_closed: false, time_slots: [] },
  tuesday: { is_closed: false, time_slots: [] },
  wednesday: { is_closed: false, time_slots: [] },
  thursday: { is_closed: false, time_slots: [] },
  friday: { is_closed: false, time_slots: [] },
  saturday: { is_closed: false, time_slots: [] },
  sunday: { is_closed: false, time_slots: [] },
})

// データの整合性を保証する関数
const sanitizeBusinessHours = (data: BusinessHoursData | null | undefined): BusinessHoursData => {
  if (!data) return getDefaultBusinessHours()
  
  const defaultData = getDefaultBusinessHours()
  
  // 各曜日のデータを確認し、不整合があればデフォルト値で補完
  Object.keys(defaultData).forEach(day => {
    const dayKey = day as keyof BusinessHoursData
    if (!data[dayKey]) {
      data[dayKey] = defaultData[dayKey]
    } else {
      // time_slotsがnullまたはundefinedの場合は空配列に初期化
      if (!Array.isArray(data[dayKey].time_slots)) {
        data[dayKey].time_slots = []
      }
    }
  })
  
  return data
}

const BusinessHoursInput: React.FC<BusinessHoursInputProps> = ({ value, onChange, className }) => {
  // 共通の営業時間を抽出
  const extractCommonHours = (data: BusinessHoursData) => {
    const openDays = DAYS_OF_WEEK.filter(day => {
      const schedule = data[day.key as keyof BusinessHoursData]
      return !schedule?.is_closed && schedule?.time_slots && schedule.time_slots.length > 0
    })
    
    const closedDays = DAYS_OF_WEEK.filter(day => {
      const schedule = data[day.key as keyof BusinessHoursData]
      return schedule?.is_closed
    }).map(day => day.key)
    
    if (openDays.length > 0) {
      const firstOpenDay = data[openDays[0].key as keyof BusinessHoursData]
      const firstSlot = firstOpenDay?.time_slots?.[0]
      if (firstSlot) {
        return {
          open_time: firstSlot.open_time,
          close_time: firstSlot.close_time,
          last_order_time: firstSlot.last_order_time,
          closedDays
        }
      }
    }
    
    return {
      open_time: '11:00',
      close_time: '22:00',
      last_order_time: '21:30',
      closedDays
    }
  }

  // 初期化時に営業時間スロットを設定
  const initializeBusinessHours = () => {
    if (value) {
      return sanitizeBusinessHours(value)
    }
    
    const initialHours = getDefaultBusinessHours()
    const defaultTimeSlot = {
      open_time: '11:00',
      close_time: '22:00',
      last_order_time: '21:30'
    }
    
    // 全ての曜日に初期時間スロットを設定
    Object.keys(initialHours).forEach(dayKey => {
      const schedule = initialHours[dayKey as keyof BusinessHoursData]
      if (!schedule.is_closed) {
        schedule.time_slots = [defaultTimeSlot]
      }
    })
    
    return initialHours
  }

  // 初期化時の共通時間を設定
  const initializeCommonHours = () => {
    if (value) {
      return extractCommonHours(sanitizeBusinessHours(value))
    }
    return {
      open_time: '11:00',
      close_time: '22:00',
      last_order_time: '21:30',
      closedDays: [] as string[]
    }
  }

  // State定義
  const [businessHours, setBusinessHours] = useState<BusinessHoursData>(
    initializeBusinessHours()
  )
  
  const [isDetailMode, setIsDetailMode] = useState(false)
  const [commonHours, setCommonHours] = useState(initializeCommonHours())

  // 初期値の設定
  useEffect(() => {
    const initializedHours = initializeBusinessHours()
    setBusinessHours(initializedHours)
    setCommonHours(initializeCommonHours())
  }, [value])

  // 営業時間データを更新
  const updateBusinessHours = (newData: BusinessHoursData) => {
    setBusinessHours(newData)
    onChange(newData)
  }

  // 定休日の切り替え
  const toggleClosedDay = (dayKey: string) => {
    const newData = { ...businessHours }
    const daySchedule = newData[dayKey as keyof BusinessHoursData]
    
    if (daySchedule.is_closed) {
      // 定休日を営業日に変更
      daySchedule.is_closed = false
      daySchedule.time_slots = [{
        open_time: commonHours.open_time,
        close_time: commonHours.close_time,
        last_order_time: commonHours.last_order_time
      }]
    } else {
      // 営業日を定休日に変更
      daySchedule.is_closed = true
      daySchedule.time_slots = []
    }
    
    updateBusinessHours(newData)
    setCommonHours(prev => ({
      ...prev,
      closedDays: prev.closedDays.includes(dayKey) 
        ? prev.closedDays.filter(d => d !== dayKey)
        : [...prev.closedDays, dayKey]
    }))
  }

  // 共通営業時間の更新
  const updateCommonHours = (field: 'open_time' | 'close_time' | 'last_order_time', value: string) => {
    const newCommonHours = { ...commonHours, [field]: value }
    setCommonHours(newCommonHours)
    
    // 営業している日にのみ適用（time_slotsが空の場合は新しく作成）
    const newData = { ...businessHours }
    DAYS_OF_WEEK.forEach(day => {
      const daySchedule = newData[day.key as keyof BusinessHoursData]
      if (!daySchedule?.is_closed) {
        if (!daySchedule.time_slots) {
          daySchedule.time_slots = []
        }
        
        // time_slotsが空の場合は新しいスロットを作成
        if (daySchedule.time_slots.length === 0) {
          daySchedule.time_slots.push({
            open_time: newCommonHours.open_time,
            close_time: newCommonHours.close_time,
            last_order_time: newCommonHours.last_order_time
          })
        } else {
          // 既存の最初のスロットを更新
          daySchedule.time_slots[0] = {
            open_time: newCommonHours.open_time,
            close_time: newCommonHours.close_time,
            last_order_time: newCommonHours.last_order_time
          }
        }
      }
    })
    
    updateBusinessHours(newData)
  }

  // 詳細モードでの時間スロット更新
  const updateTimeSlot = (dayKey: string, slotIndex: number, field: keyof TimeSlot, value: string) => {
    const newData = { ...businessHours }
    const daySchedule = newData[dayKey as keyof BusinessHoursData]
    
    if (daySchedule?.time_slots?.[slotIndex]) {
      daySchedule.time_slots[slotIndex] = {
        ...daySchedule.time_slots[slotIndex],
        [field]: value
      }
    }
    
    updateBusinessHours(newData)
  }

  // 時間スロットの追加
  const addTimeSlot = (dayKey: string) => {
    const newData = { ...businessHours }
    const daySchedule = newData[dayKey as keyof BusinessHoursData]
    
    if (daySchedule?.time_slots && daySchedule.time_slots.length < 3) {
      daySchedule.time_slots.push({
        open_time: '11:00',
        close_time: '22:00',
        last_order_time: '21:30'
      })
    }
    
    updateBusinessHours(newData)
  }

  // 時間スロットの削除
  const removeTimeSlot = (dayKey: string, slotIndex: number) => {
    const newData = { ...businessHours }
    const daySchedule = newData[dayKey as keyof BusinessHoursData]
    
    if (daySchedule?.time_slots) {
      daySchedule.time_slots.splice(slotIndex, 1)
    }
    updateBusinessHours(newData)
  }

  // 営業時間のクリア
  const clearTimeSlots = (dayKey: string) => {
    const newData = { ...businessHours }
    const daySchedule = newData[dayKey as keyof BusinessHoursData]
    
    if (daySchedule) {
      daySchedule.time_slots = []
    }
    updateBusinessHours(newData)
  }

  return (
    <div className={`business-hours-input ${className || ''}`}>
      <div className="form-group">
        <div className="time-input-header">
          <label className="form-label">営業時間設定</label>
          <div className="input-mode-toggle">
            <button
              type="button"
              className={`mode-button ${!isDetailMode ? 'active' : ''}`}
              onClick={() => setIsDetailMode(false)}
            >
              簡単設定
            </button>
            <button
              type="button"
              className={`mode-button ${isDetailMode ? 'active' : ''}`}
              onClick={() => setIsDetailMode(true)}
            >
              詳細設定
            </button>
          </div>
        </div>

        {!isDetailMode ? (
          // 簡単設定モード
          <div className="simple-mode">
            <div className="form-group">
              <label className="form-label">営業時間</label>
              <div className="time-selector">
                <select
                  value={commonHours.open_time}
                  onChange={(e) => updateCommonHours('open_time', e.target.value)}
                >
                  {TIME_OPTIONS.map(option => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </select>
                <span> - </span>
                <select
                  value={commonHours.close_time}
                  onChange={(e) => updateCommonHours('close_time', e.target.value)}
                >
                  {TIME_OPTIONS.map(option => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            <div className="form-group">
              <label className="form-label">ラストオーダー</label>
              <select
                value={commonHours.last_order_time}
                onChange={(e) => updateCommonHours('last_order_time', e.target.value)}
              >
                {TIME_OPTIONS.map(option => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
            </div>

            <div className="form-group">
              <label className="form-label">定休日</label>
              <div className="day-selector">
                {DAYS_OF_WEEK.map(day => (
                  <button
                    key={day.key}
                    type="button"
                    className={`day-button ${businessHours[day.key as keyof BusinessHoursData]?.is_closed ? 'active' : ''}`}
                    onClick={() => toggleClosedDay(day.key)}
                  >
                    {day.label}
                  </button>
                ))}
              </div>
            </div>
          </div>
        ) : (
          // 詳細設定モード
          <div className="detailed-mode">
            {DAYS_OF_WEEK.map(day => {
              const daySchedule = businessHours[day.key as keyof BusinessHoursData]
              return (
                <div key={day.key} className="form-group">
                  <div className="day-header">
                    <label className="form-label">{day.fullLabel}</label>
                    <button
                      type="button"
                      className={`day-button ${daySchedule?.is_closed ? 'active' : ''}`}
                      onClick={() => toggleClosedDay(day.key)}
                    >
                      {daySchedule?.is_closed ? '定休日' : '営業'}
                    </button>
                  </div>
                  
                  {!daySchedule?.is_closed && (
                    <div className="time-slots">
                      {daySchedule.time_slots && daySchedule.time_slots.map((slot, index) => (
                        <div key={index} className="time-slot-row">
                          <div className="time-selector">
                            <select
                              value={slot.open_time}
                              onChange={(e) => updateTimeSlot(day.key, index, 'open_time', e.target.value)}
                            >
                              {TIME_OPTIONS.map(option => (
                                <option key={option.value} value={option.value}>
                                  {option.label}
                                </option>
                              ))}
                            </select>
                            <span> - </span>
                            <select
                              value={slot.close_time}
                              onChange={(e) => updateTimeSlot(day.key, index, 'close_time', e.target.value)}
                            >
                              {TIME_OPTIONS.map(option => (
                                <option key={option.value} value={option.value}>
                                  {option.label}
                                </option>
                              ))}
                            </select>
                            <span> L.O.</span>
                            <select
                              value={slot.last_order_time}
                              onChange={(e) => updateTimeSlot(day.key, index, 'last_order_time', e.target.value)}
                            >
                              {TIME_OPTIONS.map(option => (
                                <option key={option.value} value={option.value}>
                                  {option.label}
                                </option>
                              ))}
                            </select>
                          </div>
                          
                          <div className="slot-actions">
                            {daySchedule?.time_slots && daySchedule.time_slots.length > 1 && (
                              <button
                                type="button"
                                className="btn-remove"
                                onClick={() => removeTimeSlot(day.key, index)}
                              >
                                削除
                              </button>
                            )}
                          </div>
                        </div>
                      ))}
                      
                      <div className="slot-controls">
                        {daySchedule?.time_slots && daySchedule.time_slots.length < 3 && (
                          <button
                            type="button"
                            className="btn-add"
                            onClick={() => addTimeSlot(day.key)}
                          >
                            営業時間追加
                          </button>
                        )}
                        
                        {daySchedule?.time_slots && daySchedule.time_slots.length > 0 && (
                          <button
                            type="button"
                            className="btn-clear"
                            onClick={() => clearTimeSlots(day.key)}
                          >
                            クリア
                          </button>
                        )}
                      </div>
                    </div>
                  )}
                </div>
              )
            })}
          </div>
        )}
      </div>
    </div>
  )
}

export default BusinessHoursInput