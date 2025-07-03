package repositories

import (
	"database/sql"
	"sukimise/internal/models"

	"github.com/google/uuid"
)

type ReviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(review *models.Review) error {
	query := `
		INSERT INTO reviews (id, store_id, user_id, rating, comment, photos, visit_date, is_visited, payment_amount, food_notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
	`
	review.ID = uuid.New()
	_, err := r.db.Exec(query,
		review.ID, review.StoreID, review.UserID, review.Rating, review.Comment,
		review.Photos, review.VisitDate, review.IsVisited, review.PaymentAmount, review.FoodNotes,
	)
	return err
}

func (r *ReviewRepository) GetByID(id uuid.UUID) (*models.Review, error) {
	var review models.Review
	query := `
		SELECT id, store_id, user_id, rating, comment, photos, visit_date, is_visited, payment_amount, food_notes, created_at, updated_at
		FROM reviews WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&review.ID, &review.StoreID, &review.UserID, &review.Rating, &review.Comment,
		&review.Photos, &review.VisitDate, &review.IsVisited, &review.PaymentAmount, &review.FoodNotes, &review.CreatedAt, &review.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepository) GetByStoreID(storeID uuid.UUID) ([]*models.Review, error) {
	query := `
		SELECT 
			r.id, r.store_id, r.user_id, r.rating, r.comment, r.photos, r.visit_date, r.is_visited, r.payment_amount, r.food_notes, r.created_at, r.updated_at,
			u.id, u.username
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.store_id = $1 
		ORDER BY r.created_at DESC
	`
	rows, err := r.db.Query(query, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*models.Review
	for rows.Next() {
		var review models.Review
		var user models.User
		var userID *uuid.UUID
		var username *string
		
		err := rows.Scan(
			&review.ID, &review.StoreID, &review.UserID, &review.Rating, &review.Comment,
			&review.Photos, &review.VisitDate, &review.IsVisited, &review.PaymentAmount, &review.FoodNotes, &review.CreatedAt, &review.UpdatedAt,
			&userID, &username,
		)
		if err != nil {
			return nil, err
		}
		
		// ユーザー情報が存在する場合のみ設定
		if userID != nil && username != nil {
			user.ID = *userID
			user.Username = *username
			review.User = &user
		}
		
		reviews = append(reviews, &review)
	}

	return reviews, nil
}

func (r *ReviewRepository) GetByUserID(userID uuid.UUID) ([]*models.Review, error) {
	query := `
		SELECT 
			r.id, r.store_id, r.user_id, r.rating, r.comment, r.photos, r.visit_date, r.is_visited, r.payment_amount, r.food_notes, r.created_at, r.updated_at,
			u.id, u.username
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.user_id = $1 
		ORDER BY r.created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*models.Review
	for rows.Next() {
		var review models.Review
		var user models.User
		var userUUID *uuid.UUID
		var username *string
		
		err := rows.Scan(
			&review.ID, &review.StoreID, &review.UserID, &review.Rating, &review.Comment,
			&review.Photos, &review.VisitDate, &review.IsVisited, &review.PaymentAmount, &review.FoodNotes, &review.CreatedAt, &review.UpdatedAt,
			&userUUID, &username,
		)
		if err != nil {
			return nil, err
		}
		
		// ユーザー情報が存在する場合のみ設定
		if userUUID != nil && username != nil {
			user.ID = *userUUID
			user.Username = *username
			review.User = &user
		}
		
		reviews = append(reviews, &review)
	}

	return reviews, nil
}

func (r *ReviewRepository) GetByStoreAndUser(storeID, userID uuid.UUID) (*models.Review, error) {
	var review models.Review
	query := `
		SELECT id, store_id, user_id, rating, comment, photos, visit_date, is_visited, payment_amount, food_notes, created_at, updated_at
		FROM reviews WHERE store_id = $1 AND user_id = $2
	`
	err := r.db.QueryRow(query, storeID, userID).Scan(
		&review.ID, &review.StoreID, &review.UserID, &review.Rating, &review.Comment,
		&review.Photos, &review.VisitDate, &review.IsVisited, &review.PaymentAmount, &review.FoodNotes, &review.CreatedAt, &review.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepository) Update(review *models.Review) error {
	query := `
		UPDATE reviews SET 
			rating = $2, comment = $3, photos = $4, visit_date = $5, is_visited = $6, payment_amount = $7, food_notes = $8, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(query,
		review.ID, review.Rating, review.Comment, review.Photos, review.VisitDate, review.IsVisited, review.PaymentAmount, review.FoodNotes,
	)
	return err
}

func (r *ReviewRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM reviews WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *ReviewRepository) CreateMenuItem(menuItem *models.MenuItem) error {
	query := `
		INSERT INTO menu_items (id, review_id, name, comment)
		VALUES ($1, $2, $3, $4)
	`
	menuItem.ID = uuid.New()
	_, err := r.db.Exec(query, menuItem.ID, menuItem.ReviewID, menuItem.Name, menuItem.Comment)
	return err
}

func (r *ReviewRepository) GetMenuItemsByReviewID(reviewID uuid.UUID) ([]*models.MenuItem, error) {
	query := `
		SELECT id, review_id, name, comment
		FROM menu_items WHERE review_id = $1
	`
	rows, err := r.db.Query(query, reviewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menuItems []*models.MenuItem
	for rows.Next() {
		var menuItem models.MenuItem
		err := rows.Scan(&menuItem.ID, &menuItem.ReviewID, &menuItem.Name, &menuItem.Comment)
		if err != nil {
			return nil, err
		}
		menuItems = append(menuItems, &menuItem)
	}

	return menuItems, nil
}

// GetAveragePaymentAmount returns the average payment amount from the latest 3 reviews for a store
func (r *ReviewRepository) GetAveragePaymentAmount(storeID uuid.UUID) (*float64, error) {
	query := `
		SELECT AVG(payment_amount) as avg_amount 
		FROM (
			SELECT payment_amount 
			FROM reviews 
			WHERE store_id = $1 AND payment_amount IS NOT NULL 
			ORDER BY created_at DESC 
			LIMIT 3
		) as recent_reviews
	`
	var avgAmount sql.NullFloat64
	err := r.db.QueryRow(query, storeID).Scan(&avgAmount)
	if err != nil {
		return nil, err
	}
	
	if avgAmount.Valid {
		return &avgAmount.Float64, nil
	}
	return nil, nil
}