import React from 'react'
import { BusinessHoursData, DaySchedule } from '../types/store'

interface SafeBusinessHoursDisplayProps {
  businessHours: BusinessHoursData
}

const DAYS_OF_WEEK = [
  { key: 'monday', fullLabel: '月曜日' },
  { key: 'tuesday', fullLabel: '火曜日' },
  { key: 'wednesday', fullLabel: '水曜日' },
  { key: 'thursday', fullLabel: '木曜日' },
  { key: 'friday', fullLabel: '金曜日' },
  { key: 'saturday', fullLabel: '土曜日' },
  { key: 'sunday', fullLabel: '日曜日' },
]

const SafeBusinessHoursDisplay: React.FC<SafeBusinessHoursDisplayProps> = ({ businessHours }) => {
  if (!businessHours) {
    return <div>営業時間情報なし</div>
  }

  return (
    <div className="detailed-business-hours">
      {DAYS_OF_WEEK.map(day => {
        const schedule = businessHours[day.key as keyof BusinessHoursData]
        return (
          <div key={day.key} className="day-schedule">
            <strong>{day.fullLabel}:</strong>
            {schedule?.is_closed ? (
              <span className="closed"> 定休日</span>
            ) : (
              <span className="open">
                {schedule?.time_slots?.map((slot, index) => (
                  <span key={index} className="time-slot">
                    {index > 0 && ', '}
                    {slot.open_time}-{slot.close_time}
                    {slot.last_order_time && slot.last_order_time !== slot.close_time && (
                      <span className="last-order">(L.O.{slot.last_order_time})</span>
                    )}
                  </span>
                )) || <span> 営業時間未設定</span>}
              </span>
            )}
          </div>
        )
      })}
    </div>
  )
}

export default SafeBusinessHoursDisplay