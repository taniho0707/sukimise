import React from 'react'
import './ParkingDisplay.css'

interface ParkingInfo {
  hasPrivateParking: boolean
  hasCoinParkingService: boolean
  hasNearbyCoinParking: boolean
  memo: string
}

interface ParkingDisplayProps {
  parkingInfo: string
  className?: string
}

const ParkingDisplay: React.FC<ParkingDisplayProps> = ({ parkingInfo, className }) => {
  // 駐車場情報を解析する関数
  const parseParkingInfo = (text: string): ParkingInfo | null => {
    if (!text) return null

    // 構造化データ（JSON）を試す
    try {
      const parsed = JSON.parse(text)
      if (parsed && typeof parsed === 'object') {
        return {
          hasPrivateParking: parsed.hasPrivateParking || false,
          hasCoinParkingService: parsed.hasCoinParkingService || false,
          hasNearbyCoinParking: parsed.hasNearbyCoinParking || false,
          memo: parsed.memo || ''
        }
      }
    } catch {
      // JSON以外の場合は従来のテキスト表示にフォールバック
      return null
    }

    return null
  }

  const parkingData = parseParkingInfo(parkingInfo)

  // 構造化データの場合
  if (parkingData) {
    const hasAnyParking = parkingData.hasPrivateParking || 
                         parkingData.hasCoinParkingService || 
                         parkingData.hasNearbyCoinParking

    return (
      <div className={`parking-display ${className || ''}`}>
        {hasAnyParking ? (
          <div className="parking-features">
            {parkingData.hasPrivateParking && (
              <div className="parking-feature available">
                <span className="feature-icon">🅿️</span>
                <span className="feature-text">専用駐車場あり</span>
              </div>
            )}
            
            {parkingData.hasCoinParkingService && (
              <div className="parking-feature available">
                <span className="feature-icon">🎫</span>
                <span className="feature-text">コインパーキングサービスあり</span>
              </div>
            )}
            
            {parkingData.hasNearbyCoinParking && (
              <div className="parking-feature available">
                <span className="feature-icon">🗺️</span>
                <span className="feature-text">近隣コインパーキングあり</span>
              </div>
            )}
          </div>
        ) : (
          <div className="parking-feature unavailable">
            <span className="feature-icon">❌</span>
            <span className="feature-text">駐車場情報なし</span>
          </div>
        )}

        {parkingData.memo && (
          <div className="parking-memo">
            <div className="memo-label">詳細情報</div>
            <div className="memo-content">
              {parkingData.memo.split('\n').map((line, index) => (
                <div key={index}>{line}</div>
              ))}
            </div>
          </div>
        )}
      </div>
    )
  }

  // 従来のテキスト形式の場合
  if (parkingInfo) {
    return (
      <div className={`parking-display ${className || ''}`}>
        <div className="parking-text">
          {parkingInfo.split('\n').map((line, index) => (
            <div key={index}>{line}</div>
          ))}
        </div>
      </div>
    )
  }

  // 情報がない場合
  return (
    <div className={`parking-display ${className || ''}`}>
      <div className="parking-feature unavailable">
        <span className="feature-icon">❓</span>
        <span className="feature-text">駐車場情報未設定</span>
      </div>
    </div>
  )
}

export default ParkingDisplay