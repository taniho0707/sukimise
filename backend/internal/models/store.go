package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Store struct {
	ID            uuid.UUID        `json:"id" db:"id"`
	Name          string           `json:"name" db:"name"`
	Address       string           `json:"address" db:"address"`
	Latitude      float64          `json:"latitude" db:"latitude"`
	Longitude     float64          `json:"longitude" db:"longitude"`
	Categories    StringArray      `json:"categories" db:"categories"`
	BusinessHours BusinessHoursData `json:"business_hours" db:"business_hours"`
	ParkingInfo   string           `json:"parking_info" db:"parking_info"`
	WebsiteURL    string           `json:"website_url" db:"website_url"`
	GoogleMapURL  string           `json:"google_map_url" db:"google_map_url"`
	SnsUrls       StringArray      `json:"sns_urls" db:"sns_urls"`
	Tags          StringArray      `json:"tags" db:"tags"`
	Photos        StringArray      `json:"photos" db:"photos"`
	CreatedBy     uuid.UUID        `json:"created_by" db:"created_by"`
	CreatedAt     time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at" db:"updated_at"`
}

// BusinessHoursData represents detailed business hours structure
type BusinessHoursData struct {
	Monday    DaySchedule `json:"monday"`
	Tuesday   DaySchedule `json:"tuesday"`
	Wednesday DaySchedule `json:"wednesday"`
	Thursday  DaySchedule `json:"thursday"`
	Friday    DaySchedule `json:"friday"`
	Saturday  DaySchedule `json:"saturday"`
	Sunday    DaySchedule `json:"sunday"`
}

// DaySchedule represents a day's schedule with up to 3 time slots
type DaySchedule struct {
	IsClosed   bool       `json:"is_closed"`
	TimeSlots  []TimeSlot `json:"time_slots"` // Maximum 3 slots
}

// TimeSlot represents one time period in a day
type TimeSlot struct {
	OpenTime      string `json:"open_time"`       // HH:MM format
	CloseTime     string `json:"close_time"`      // HH:MM format
	LastOrderTime string `json:"last_order_time"` // HH:MM format
}

// Value implements driver.Valuer for BusinessHoursData
func (b BusinessHoursData) Value() (driver.Value, error) {
	return json.Marshal(b)
}

// Scan implements sql.Scanner for BusinessHoursData
func (b *BusinessHoursData) Scan(value interface{}) error {
	if value == nil {
		*b = GetDefaultBusinessHours()
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("cannot scan business_hours: unsupported type")
	}

	return json.Unmarshal(bytes, b)
}

// GetDefaultBusinessHours returns default business hours structure
func GetDefaultBusinessHours() BusinessHoursData {
	return BusinessHoursData{
		Monday:    DaySchedule{IsClosed: false, TimeSlots: []TimeSlot{}},
		Tuesday:   DaySchedule{IsClosed: false, TimeSlots: []TimeSlot{}},
		Wednesday: DaySchedule{IsClosed: false, TimeSlots: []TimeSlot{}},
		Thursday:  DaySchedule{IsClosed: false, TimeSlots: []TimeSlot{}},
		Friday:    DaySchedule{IsClosed: false, TimeSlots: []TimeSlot{}},
		Saturday:  DaySchedule{IsClosed: false, TimeSlots: []TimeSlot{}},
		Sunday:    DaySchedule{IsClosed: false, TimeSlots: []TimeSlot{}},
	}
}



type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, a)
}