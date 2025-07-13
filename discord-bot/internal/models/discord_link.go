package models

import (
	"time"

	"github.com/google/uuid"
)

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
	IsClosed  bool       `json:"is_closed"`
	TimeSlots []TimeSlot `json:"time_slots"`
}

// TimeSlot represents one time period in a day
type TimeSlot struct {
	OpenTime      string `json:"open_time"`
	CloseTime     string `json:"close_time"`
	LastOrderTime string `json:"last_order_time"`
}

type DiscordLink struct {
	ID           uuid.UUID `json:"id"`
	DiscordID    string    `json:"discord_id"`
	UserID       uuid.UUID `json:"user_id"`
	Username     string    `json:"username"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenExpiry  time.Time `json:"token_expiry"`
	LinkedAt     time.Time `json:"linked_at"`
	LastUsedAt   time.Time `json:"last_used_at"`
}

type SukimiseUser struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         SukimiseUser `json:"user"`
}

type StoreCreateRequest struct {
	Name            string            `json:"name"`
	Address         string            `json:"address"`
	Latitude        float64           `json:"latitude"`
	Longitude       float64           `json:"longitude"`
	Categories      []string          `json:"categories"`
	BusinessHours   BusinessHoursData `json:"business_hours"`
	ParkingInfo     string            `json:"parking_info"`
	WebsiteURL      string            `json:"website_url"`
	GoogleMapURL    string            `json:"google_map_url"`
	SNSUrls         []string          `json:"sns_urls"`
	Tags            []string          `json:"tags"`
	Photos          []string          `json:"photos"`
}

type StoreCreateResponse struct {
	ID            uuid.UUID         `json:"id"`
	Name          string            `json:"name"`
	Address       string            `json:"address"`
	Latitude      float64           `json:"latitude"`
	Longitude     float64           `json:"longitude"`
	Categories    []string          `json:"categories"`
	BusinessHours BusinessHoursData `json:"business_hours"`
	ParkingInfo   string            `json:"parking_info"`
	WebsiteURL    string            `json:"website_url"`
	GoogleMapURL  string            `json:"google_map_url"`
	SNSUrls       []string          `json:"sns_urls"`
	Tags          []string          `json:"tags"`
	Photos        []string          `json:"photos"`
	CreatedBy     uuid.UUID         `json:"created_by"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`
	Message       string            `json:"message"`
}

type APIResponse struct {
	Success bool                 `json:"success"`
	Data    StoreCreateResponse  `json:"data"`
}