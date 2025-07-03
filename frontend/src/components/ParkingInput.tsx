import React, { useState, useEffect } from 'react'

interface ParkingInfo {
  hasPrivateParking: boolean
  hasCoinParkingService: boolean
  hasNearbyCoinParking: boolean
  memo: string
}

interface ParkingInputProps {
  value: string
  onChange: (value: string) => void
  className?: string
}

const ParkingInput: React.FC<ParkingInputProps> = ({ value, onChange, className }) => {
  const [parkingInfo, setParkingInfo] = useState<ParkingInfo>({
    hasPrivateParking: false,
    hasCoinParkingService: false,
    hasNearbyCoinParking: false,
    memo: ''
  })

  // valueからparkingInfoを解析する関数
  const parseParkingInfo = (text: string): ParkingInfo => {
    const defaultInfo: ParkingInfo = {
      hasPrivateParking: false,
      hasCoinParkingService: false,
      hasNearbyCoinParking: false,
      memo: ''
    }

    if (!text) return defaultInfo

    // 構造化データか判定（JSON形式）
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
      // JSONでない場合は従来の文字列として処理
    }

    // 従来の文字列形式から推測
    const lowerText = text.toLowerCase()
    return {
      hasPrivateParking: lowerText.includes('専用') || lowerText.includes('自社') || lowerText.includes('店舗駐車場'),
      hasCoinParkingService: lowerText.includes('サービス') || lowerText.includes('割引'),
      hasNearbyCoinParking: lowerText.includes('近隣') || lowerText.includes('周辺') || lowerText.includes('コイン'),
      memo: text
    }
  }

  // parkingInfoからvalueを生成する関数
  const generateParkingText = (info: ParkingInfo): string => {
    // 構造化データとしてJSON形式で保存
    return JSON.stringify(info)
  }

  // 初期値の設定
  useEffect(() => {
    const parsed = parseParkingInfo(value)
    setParkingInfo(parsed)
  }, [value])

  // データ更新時のコールバック
  const updateParkingInfo = (newInfo: Partial<ParkingInfo>) => {
    const updatedInfo = { ...parkingInfo, ...newInfo }
    setParkingInfo(updatedInfo)
    const textValue = generateParkingText(updatedInfo)
    onChange(textValue)
  }

  const handleToggle = (field: keyof Omit<ParkingInfo, 'memo'>) => {
    updateParkingInfo({ [field]: !parkingInfo[field] })
  }

  const handleMemoChange = (memo: string) => {
    updateParkingInfo({ memo })
  }

  // 表示用のサマリーを生成
  const generateDisplaySummary = (): string => {
    const parts = []
    
    if (parkingInfo.hasPrivateParking) {
      parts.push('専用駐車場あり')
    }
    
    if (parkingInfo.hasCoinParkingService) {
      parts.push('コインパーキングサービスあり')
    }
    
    if (parkingInfo.hasNearbyCoinParking) {
      parts.push('近隣コインパーキングあり')
    }
    
    if (parts.length === 0) {
      parts.push('駐車場情報なし')
    }
    
    if (parkingInfo.memo) {
      parts.push(`メモ: ${parkingInfo.memo}`)
    }
    
    return parts.join(' / ')
  }

  return (
    <div className={className || ''} style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
      <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
        {/* 専用駐車場 */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
          <label style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            padding: '12px',
            border: '1px solid #e9ecef',
            borderRadius: '8px',
            backgroundColor: '#f8f9fa',
            cursor: 'pointer',
            transition: 'all 0.2s ease'
          }}
          onMouseEnter={(e) => e.currentTarget.style.backgroundColor = '#e9ecef'}
          onMouseLeave={(e) => e.currentTarget.style.backgroundColor = '#f8f9fa'}
          >
            <span style={{ fontWeight: '500', color: '#333', fontSize: '14px' }}>専用駐車場</span>
            <button
              type="button"
              onClick={() => handleToggle('hasPrivateParking')}
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: '8px',
                padding: '6px 12px',
                border: '1px solid #dee2e6',
                borderRadius: '20px',
                backgroundColor: parkingInfo.hasPrivateParking ? '#0066cc' : 'white',
                borderColor: parkingInfo.hasPrivateParking ? '#0066cc' : '#dee2e6',
                color: parkingInfo.hasPrivateParking ? 'white' : 'inherit',
                cursor: 'pointer',
                transition: 'all 0.2s ease',
                fontSize: '13px'
              }}
            >
              <span style={{
                width: '16px',
                height: '16px',
                borderRadius: '50%',
                backgroundColor: parkingInfo.hasPrivateParking ? 'white' : '#dee2e6',
                transition: 'all 0.2s ease'
              }}></span>
              <span style={{ fontWeight: '500', minWidth: '30px' }}>
                {parkingInfo.hasPrivateParking ? 'あり' : 'なし'}
              </span>
            </button>
          </label>
          <small style={{ fontSize: '12px', color: '#6c757d', marginLeft: '12px' }}>店舗専用の駐車場がある場合</small>
        </div>

        {/* コインパーキングサービス */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
          <label style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            padding: '12px',
            border: '1px solid #e9ecef',
            borderRadius: '8px',
            backgroundColor: '#f8f9fa',
            cursor: 'pointer',
            transition: 'all 0.2s ease'
          }}
          onMouseEnter={(e) => e.currentTarget.style.backgroundColor = '#e9ecef'}
          onMouseLeave={(e) => e.currentTarget.style.backgroundColor = '#f8f9fa'}
          >
            <span style={{ fontWeight: '500', color: '#333', fontSize: '14px' }}>コインパーキングサービス</span>
            <button
              type="button"
              onClick={() => handleToggle('hasCoinParkingService')}
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: '8px',
                padding: '6px 12px',
                border: '1px solid #dee2e6',
                borderRadius: '20px',
                backgroundColor: parkingInfo.hasCoinParkingService ? '#0066cc' : 'white',
                borderColor: parkingInfo.hasCoinParkingService ? '#0066cc' : '#dee2e6',
                color: parkingInfo.hasCoinParkingService ? 'white' : 'inherit',
                cursor: 'pointer',
                transition: 'all 0.2s ease',
                fontSize: '13px'
              }}
            >
              <span style={{
                width: '16px',
                height: '16px',
                borderRadius: '50%',
                backgroundColor: parkingInfo.hasCoinParkingService ? 'white' : '#dee2e6',
                transition: 'all 0.2s ease'
              }}></span>
              <span style={{ fontWeight: '500', minWidth: '30px' }}>
                {parkingInfo.hasCoinParkingService ? 'あり' : 'なし'}
              </span>
            </button>
          </label>
          <small style={{ fontSize: '12px', color: '#6c757d', marginLeft: '12px' }}>駐車料金の割引やサービスがある場合</small>
        </div>

        {/* 近隣コインパーキング */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
          <label style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            padding: '12px',
            border: '1px solid #e9ecef',
            borderRadius: '8px',
            backgroundColor: '#f8f9fa',
            cursor: 'pointer',
            transition: 'all 0.2s ease'
          }}
          onMouseEnter={(e) => e.currentTarget.style.backgroundColor = '#e9ecef'}
          onMouseLeave={(e) => e.currentTarget.style.backgroundColor = '#f8f9fa'}
          >
            <span style={{ fontWeight: '500', color: '#333', fontSize: '14px' }}>近隣コインパーキング</span>
            <button
              type="button"
              onClick={() => handleToggle('hasNearbyCoinParking')}
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: '8px',
                padding: '6px 12px',
                border: '1px solid #dee2e6',
                borderRadius: '20px',
                backgroundColor: parkingInfo.hasNearbyCoinParking ? '#0066cc' : 'white',
                borderColor: parkingInfo.hasNearbyCoinParking ? '#0066cc' : '#dee2e6',
                color: parkingInfo.hasNearbyCoinParking ? 'white' : 'inherit',
                cursor: 'pointer',
                transition: 'all 0.2s ease',
                fontSize: '13px'
              }}
            >
              <span style={{
                width: '16px',
                height: '16px',
                borderRadius: '50%',
                backgroundColor: parkingInfo.hasNearbyCoinParking ? 'white' : '#dee2e6',
                transition: 'all 0.2s ease'
              }}></span>
              <span style={{ fontWeight: '500', minWidth: '30px' }}>
                {parkingInfo.hasNearbyCoinParking ? 'あり' : 'なし'}
              </span>
            </button>
          </label>
          <small style={{ fontSize: '12px', color: '#6c757d', marginLeft: '12px' }}>歩いて行ける距離にコインパーキングがある場合</small>
        </div>
      </div>

      {/* メモ欄 */}
      <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
        <label style={{ fontWeight: '500', color: '#333', fontSize: '14px' }}>駐車場メモ</label>
        <textarea
          value={parkingInfo.memo}
          onChange={(e) => handleMemoChange(e.target.value)}
          placeholder="駐車台数、料金、注意事項など詳細情報を入力"
          rows={3}
          style={{
            width: '100%',
            padding: '8px 12px',
            border: '1px solid #dee2e6',
            borderRadius: '4px',
            fontSize: '14px',
            lineHeight: '1.4',
            resize: 'vertical'
          }}
          onFocus={(e) => {
            e.currentTarget.style.outline = 'none'
            e.currentTarget.style.borderColor = '#0066cc'
            e.currentTarget.style.boxShadow = '0 0 0 2px rgba(0, 102, 204, 0.25)'
          }}
          onBlur={(e) => {
            e.currentTarget.style.borderColor = '#dee2e6'
            e.currentTarget.style.boxShadow = 'none'
          }}
        />
      </div>

      {/* プレビュー */}
      <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
        <label style={{ fontWeight: '500', color: '#333', fontSize: '14px' }}>表示プレビュー</label>
        <div style={{
          padding: '12px',
          backgroundColor: '#f8f9fa',
          border: '1px solid #e9ecef',
          borderRadius: '4px',
          fontSize: '14px',
          color: '#495057',
          minHeight: '20px'
        }}>
          {generateDisplaySummary()}
        </div>
      </div>

    </div>
  )
}

export default ParkingInput