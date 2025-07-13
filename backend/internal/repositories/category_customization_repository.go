package repositories

import (
	"database/sql"
	"sukimise/internal/models"
	"time"

	"github.com/google/uuid"
)

type CategoryCustomizationRepository struct {
	db *sql.DB
}

func NewCategoryCustomizationRepository(db *sql.DB) CategoryCustomizationRepositoryInterface {
	return &CategoryCustomizationRepository{db: db}
}

func (r *CategoryCustomizationRepository) Create(categoryCustomization *models.CategoryCustomization) error {
	query := `
		INSERT INTO category_customizations (id, category_name, icon, color, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	
	categoryCustomization.ID = uuid.New()
	categoryCustomization.CreatedAt = time.Now()
	categoryCustomization.UpdatedAt = time.Now()
	
	_, err := r.db.Exec(query,
		categoryCustomization.ID,
		categoryCustomization.CategoryName,
		categoryCustomization.Icon,
		categoryCustomization.Color,
		categoryCustomization.CreatedAt,
		categoryCustomization.UpdatedAt,
	)
	
	return err
}

func (r *CategoryCustomizationRepository) GetByCategoryName(categoryName string) (*models.CategoryCustomization, error) {
	query := `
		SELECT id, category_name, icon, color, created_at, updated_at
		FROM category_customizations
		WHERE category_name = $1
	`
	
	var categoryCustomization models.CategoryCustomization
	err := r.db.QueryRow(query, categoryName).Scan(
		&categoryCustomization.ID,
		&categoryCustomization.CategoryName,
		&categoryCustomization.Icon,
		&categoryCustomization.Color,
		&categoryCustomization.CreatedAt,
		&categoryCustomization.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}
	
	return &categoryCustomization, nil
}

func (r *CategoryCustomizationRepository) GetAll() ([]*models.CategoryCustomization, error) {
	query := `
		SELECT id, category_name, icon, color, created_at, updated_at
		FROM category_customizations
		ORDER BY category_name ASC
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var customizations []*models.CategoryCustomization
	for rows.Next() {
		var customization models.CategoryCustomization
		err := rows.Scan(
			&customization.ID,
			&customization.CategoryName,
			&customization.Icon,
			&customization.Color,
			&customization.CreatedAt,
			&customization.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		customizations = append(customizations, &customization)
	}
	
	return customizations, nil
}

func (r *CategoryCustomizationRepository) Update(categoryCustomization *models.CategoryCustomization) error {
	query := `
		UPDATE category_customizations
		SET icon = $1, color = $2, updated_at = $3
		WHERE category_name = $4
	`
	
	categoryCustomization.UpdatedAt = time.Now()
	
	_, err := r.db.Exec(query,
		categoryCustomization.Icon,
		categoryCustomization.Color,
		categoryCustomization.UpdatedAt,
		categoryCustomization.CategoryName,
	)
	
	return err
}

func (r *CategoryCustomizationRepository) Delete(categoryName string) error {
	query := `DELETE FROM category_customizations WHERE category_name = $1`
	_, err := r.db.Exec(query, categoryName)
	return err
}