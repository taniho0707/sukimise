import React, { useState, useEffect } from 'react'

interface BusinessHoursData {
  closedDays: string[]
  openTime: string
  closeTime: string
  lastOrderTime: string
}


interface BusinessHoursInputProps {
  value: string
  onChange: (value: string) => void
  className?: string
}

const DAYS_OF_WEEK = [
  { value: 'monday', label: '月', fullLabel: '月曜日' },
  { value: 'tuesday', label: '火', fullLabel: '火曜日' },
  { value: 'wednesday', label: '水', fullLabel: '水曜日' },
  { value: 'thursday', label: '木', fullLabel: '木曜日' },
  { value: 'friday', label: '金', fullLabel: '金曜日' },
  { value: 'saturday', label: '土', fullLabel: '土曜日' },
  { value: 'sunday', label: '日', fullLabel: '日曜日' },
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

const BusinessHoursInput: React.FC<BusinessHoursInputProps> = ({ value, onChange, className }) => {
  const [businessHours, setBusinessHours] = useState<BusinessHoursData>({
    closedDays: [],
    openTime: '11:00',
    closeTime: '22:00',
    lastOrderTime: '21:30',
  })
  

  // テキスト値から構造化データを解析
  const parseBusinessHours = (text: string): BusinessHoursData => {
    const defaultData: BusinessHoursData = {
      closedDays: [],
      openTime: '11:00',
      closeTime: '22:00',
      lastOrderTime: '21:30',
    }

    if (!text) return defaultData

    // 定休日の抽出
    const closedDayMatch = text.match(/定休日[：:]\s*(.+?)(?:\n|$)/i)
    let closedDays: string[] = []
    if (closedDayMatch) {
      const closedDayText = closedDayMatch[1]
      DAYS_OF_WEEK.forEach(day => {
        // フルネーム（月曜日）で先にチェック、その後単文字（月）をチェック
        if (closedDayText.includes(day.fullLabel)) {
          closedDays.push(day.value)
        } else if (closedDayText.includes(day.label) && day.label !== '日') {
          // 「日」以外の単文字をチェック（「日」は「月曜日」等に含まれるため除外）
          closedDays.push(day.value)
        } else if (day.label === '日' && (closedDayText.includes('日曜') || closedDayText === '日')) {
          // 「日」は「日曜」または単独の「日」の場合のみマッチ
          closedDays.push(day.value)
        }
      })
    }

    // 営業時間の抽出
    const timeMatch = text.match(/(\d{1,2}):(\d{2})[-～〜](\d{1,2}):(\d{2})/);
    let openTime = '11:00'
    let closeTime = '22:00'
    if (timeMatch) {
      openTime = `${timeMatch[1].padStart(2, '0')}:${timeMatch[2]}`
      closeTime = `${timeMatch[3].padStart(2, '0')}:${timeMatch[4]}`
    }

    // ラストオーダーの抽出
    const lastOrderMatch = text.match(/(?:ラストオーダー|L\.O\.?|LO)[：:\s]*(\d{1,2}):(\d{2})/i)
    let lastOrderTime = '21:30'
    if (lastOrderMatch) {
      lastOrderTime = `${lastOrderMatch[1].padStart(2, '0')}:${lastOrderMatch[2]}`
    }

    return { closedDays, openTime, closeTime, lastOrderTime }
  }

  // 構造化データからテキスト値を生成
  const generateBusinessHoursText = (data: BusinessHoursData): string => {
    const parts = []
    
    // 営業時間
    parts.push(`営業時間: ${data.openTime}-${data.closeTime}`)
    
    // ラストオーダー
    if (data.lastOrderTime && data.lastOrderTime !== data.closeTime) {
      parts.push(`ラストオーダー: ${data.lastOrderTime}`)
    }
    
    // 定休日
    if (data.closedDays.length > 0) {
      const closedDayLabels = data.closedDays
        .map(day => DAYS_OF_WEEK.find(d => d.value === day)?.fullLabel)
        .filter(Boolean)
      parts.push(`定休日: ${closedDayLabels.join('、')}`)
    } else {
      parts.push('定休日: 年中無休')
    }
    
    return parts.join('\n')
  }

  // 初期値の設定
  useEffect(() => {
    const parsed = parseBusinessHours(value)
    setBusinessHours(parsed)
  }, [value])

  // データ更新時のコールバック
  const updateBusinessHours = (newData: Partial<BusinessHoursData>) => {
    const updatedData = { ...businessHours, ...newData }
    setBusinessHours(updatedData)
    const textValue = generateBusinessHoursText(updatedData)
    onChange(textValue)
  }


  const handleClosedDayToggle = (day: string) => {
    const newClosedDays = businessHours.closedDays.includes(day)
      ? businessHours.closedDays.filter(d => d !== day)
      : [...businessHours.closedDays, day]
    
    updateBusinessHours({ closedDays: newClosedDays })
  }

  return (
    <div className={`business-hours-input ${className || ''}`}>
      {/* 定休日選択 */}
      <div className="form-group">
        <label className="form-label">定休日</label>
        <div className="day-selector">
          {DAYS_OF_WEEK.map((day) => (
            <button
              key={day.value}
              type="button"
              onClick={() => handleClosedDayToggle(day.value)}
              className={`day-button ${businessHours.closedDays.includes(day.value) ? 'active' : ''}`}
            >
              {day.label}
            </button>
          ))}
        </div>
      </div>

      {/* 営業開始時間 */}
      <div className="form-group">
        <label className="form-label">営業開始時間</label>
        <div className="time-selector">
          {TIME_OPTIONS.map((time) => (
            <button
              key={time.value}
              type="button"
              onClick={() => updateBusinessHours({ openTime: time.value })}
              className={`time-button ${businessHours.openTime === time.value ? 'active' : ''}`}
            >
              {time.label}
            </button>
          ))}
        </div>
      </div>

      {/* ラストオーダー時間 */}
      <div className="form-group">
        <label className="form-label">ラストオーダー時間</label>
        <div className="time-selector">
          {TIME_OPTIONS.map((time) => (
            <button
              key={time.value}
              type="button"
              onClick={() => updateBusinessHours({ lastOrderTime: time.value })}
              className={`time-button ${businessHours.lastOrderTime === time.value ? 'active' : ''}`}
            >
              {time.label}
            </button>
          ))}
        </div>
      </div>

      {/* 営業終了時間 */}
      <div className="form-group">
        <label className="form-label">営業終了時間</label>
        <div className="time-selector">
          {TIME_OPTIONS.map((time) => (
            <button
              key={time.value}
              type="button"
              onClick={() => updateBusinessHours({ closeTime: time.value })}
              className={`time-button ${businessHours.closeTime === time.value ? 'active' : ''}`}
            >
              {time.label}
            </button>
          ))}
        </div>
      </div>

      {/* プレビュー */}
      <div className="form-group">
        <label className="form-label">プレビュー</label>
        <div className="business-hours-preview">
          {generateBusinessHoursText(businessHours).split('\n').map((line, index) => (
            <div key={index} className="preview-line">{line}</div>
          ))}
        </div>
      </div>

    </div>
  )
}

export default BusinessHoursInput