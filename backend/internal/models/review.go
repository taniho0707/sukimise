package models

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID            uuid.UUID   `json:"id" db:"id"`
	StoreID       uuid.UUID   `json:"store_id" db:"store_id"`
	UserID        uuid.UUID   `json:"user_id" db:"user_id"`
	Rating        int         `json:"rating" db:"rating"`
	Comment       *string     `json:"comment" db:"comment"`
	Photos        StringArray `json:"photos" db:"photos"`
	VisitDate     *time.Time  `json:"visit_date" db:"visit_date"`
	IsVisited     bool        `json:"is_visited" db:"is_visited"`
	PaymentAmount *int        `json:"payment_amount" db:"payment_amount"` // 支払金額（円）
	FoodNotes     *string     `json:"food_notes" db:"food_notes"`         // 料理についてのメモ
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
	User          *User       `json:"user,omitempty"`                     // ユーザー情報（JOINで取得）
}

type MenuItem struct {
	ID       uuid.UUID `json:"id" db:"id"`
	ReviewID uuid.UUID `json:"review_id" db:"review_id"`
	Name     string    `json:"name" db:"name"`
	Comment  string    `json:"comment" db:"comment"`
}