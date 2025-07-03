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
	BusinessHours string           `json:"business_hours" db:"business_hours"`
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