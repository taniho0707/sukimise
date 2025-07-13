package repositories

import (
	"database/sql"
	"fmt"
	"strings"
	"sukimise/internal/models"

	"github.com/google/uuid"
)

type StoreRepository struct {
	db *sql.DB
}

func NewStoreRepository(db *sql.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

func (r *StoreRepository) Create(store *models.Store) error {
	query := `
		INSERT INTO stores (id, name, address, latitude, longitude, categories, business_hours, 
						  parking_info, website_url, google_map_url, sns_urls, 
						  tags, photos, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, NOW(), NOW())
	`
	store.ID = uuid.New()
	_, err := r.db.Exec(query,
		store.ID, store.Name, store.Address, store.Latitude, store.Longitude,
		store.Categories, store.BusinessHours, store.ParkingInfo,
		store.WebsiteURL, store.GoogleMapURL, store.SnsUrls, store.Tags,
		store.Photos, store.CreatedBy,
	)
	return err
}

func (r *StoreRepository) GetByID(id uuid.UUID) (*models.Store, error) {
	var store models.Store
	query := `
		SELECT id, name, address, latitude, longitude, categories, business_hours,
			   parking_info, website_url, google_map_url, sns_urls,
			   tags, photos, created_by, created_at, updated_at
		FROM stores WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&store.ID, &store.Name, &store.Address, &store.Latitude, &store.Longitude,
		&store.Categories, &store.BusinessHours, &store.ParkingInfo,
		&store.WebsiteURL, &store.GoogleMapURL, &store.SnsUrls, &store.Tags,
		&store.Photos, &store.CreatedBy, &store.CreatedAt, &store.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &store, nil
}

type StoreFilter struct {
	Name              string
	Categories        []string
	CategoriesOperator string // "OR" or "AND"
	Tags              []string
	TagsOperator      string // "OR" or "AND"
	Latitude          *float64
	Longitude         *float64
	Radius            *float64
	BusinessDay       string
	BusinessTime      string
	Limit             int
	Offset            int
}

func (r *StoreRepository) GetAll(filter *StoreFilter) ([]*models.Store, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	baseQuery := `
		SELECT id, name, address, latitude, longitude, categories, business_hours,
			   parking_info, website_url, google_map_url, sns_urls,
			   tags, photos, created_by, created_at, updated_at
		FROM stores
	`

	if filter.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argIndex))
		args = append(args, "%"+filter.Name+"%")
		argIndex++
	}

	if len(filter.Categories) > 0 {
		placeholders := make([]string, len(filter.Categories))
		for i, cat := range filter.Categories {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, cat)
			argIndex++
		}
		
		// デフォルトはOR、ANDが指定された場合は&&演算子を使用
		operator := "?|" // OR演算子
		if filter.CategoriesOperator == "AND" {
			operator = "?&" // AND演算子
		}
		conditions = append(conditions, fmt.Sprintf("categories %s ARRAY[%s]::text[]", operator, strings.Join(placeholders, ",")))
	}

	if len(filter.Tags) > 0 {
		placeholders := make([]string, len(filter.Tags))
		for i, tag := range filter.Tags {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, tag)
			argIndex++
		}
		
		// デフォルトはAND、ORが指定された場合は?|演算子を使用
		operator := "?&" // AND演算子
		if filter.TagsOperator == "OR" {
			operator = "?|" // OR演算子
		}
		conditions = append(conditions, fmt.Sprintf("tags %s ARRAY[%s]::text[]", operator, strings.Join(placeholders, ",")))
	}

	if filter.Latitude != nil && filter.Longitude != nil && filter.Radius != nil {
		conditions = append(conditions, fmt.Sprintf(`
			ST_DWithin(
				ST_Point(longitude, latitude)::geography,
				ST_Point($%d, $%d)::geography,
				$%d
			)
		`, argIndex, argIndex+1, argIndex+2))
		args = append(args, *filter.Longitude, *filter.Latitude, *filter.Radius)
		argIndex += 3
	}


	// 営業時間検索（JSON形式対応）
	if filter.BusinessDay != "" && filter.BusinessTime != "" {
		// 両方指定: 指定曜日の指定時間に営業している店
		conditions = append(conditions, fmt.Sprintf(`
			(business_hours IS NOT NULL AND
			 business_hours->$%d->>'is_closed' != 'true' AND
			 EXISTS (
			   SELECT 1 FROM jsonb_array_elements(business_hours->$%d->'time_slots') AS slot
			   WHERE $%d::time >= (slot->>'open_time')::time 
			     AND $%d::time <= COALESCE(
			       NULLIF(slot->>'last_order_time', ''),
			       slot->>'close_time'
			     )::time
			 ))`, argIndex, argIndex+1, argIndex+2, argIndex+3))
		args = append(args, filter.BusinessDay, filter.BusinessDay, filter.BusinessTime, filter.BusinessTime)
		argIndex += 4
	} else if filter.BusinessDay != "" {
		// 営業日のみ指定: その日が休業日ではない店
		conditions = append(conditions, fmt.Sprintf(`
			(business_hours IS NULL OR 
			 business_hours->$%d->>'is_closed' != 'true')`, argIndex))
		args = append(args, filter.BusinessDay)
		argIndex++
	} else if filter.BusinessTime != "" {
		// 営業時間のみ指定: 1週間のうちどこか1日でもその時間に営業している店
		timeConditions := []string{}
		for _, day := range []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"} {
			timeConditions = append(timeConditions, fmt.Sprintf(`
				(business_hours->$%d->>'is_closed' != 'true' AND
				 EXISTS (
				   SELECT 1 FROM jsonb_array_elements(business_hours->$%d->'time_slots') AS slot
				   WHERE $%d::time >= (slot->>'open_time')::time 
				     AND $%d::time <= COALESCE(
				       NULLIF(slot->>'last_order_time', ''),
				       slot->>'close_time'
				     )::time
				 ))`, argIndex, argIndex+1, argIndex+2, argIndex+3))
			args = append(args, day, day, filter.BusinessTime, filter.BusinessTime)
			argIndex += 4
		}
		conditions = append(conditions, "("+strings.Join(timeConditions, " OR ")+")")
	}

	query := baseQuery
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stores []*models.Store
	for rows.Next() {
		var store models.Store
		err := rows.Scan(
			&store.ID, &store.Name, &store.Address, &store.Latitude, &store.Longitude,
			&store.Categories, &store.BusinessHours, &store.ParkingInfo,
			&store.WebsiteURL, &store.GoogleMapURL, &store.SnsUrls, &store.Tags,
			&store.Photos, &store.CreatedBy, &store.CreatedAt, &store.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		stores = append(stores, &store)
	}

	return stores, nil
}

// GetCount returns the total count of stores matching the filter
func (r *StoreRepository) GetCount(filter *StoreFilter) (int, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	baseQuery := `SELECT COUNT(*) FROM stores`

	if filter.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argIndex))
		args = append(args, "%"+filter.Name+"%")
		argIndex++
	}

	if len(filter.Categories) > 0 {
		placeholders := make([]string, len(filter.Categories))
		for i, cat := range filter.Categories {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, cat)
			argIndex++
		}
		
		operator := "?|" // OR演算子
		if filter.CategoriesOperator == "AND" {
			operator = "?&" // AND演算子
		}
		conditions = append(conditions, fmt.Sprintf("categories %s ARRAY[%s]::text[]", operator, strings.Join(placeholders, ",")))
	}

	if len(filter.Tags) > 0 {
		placeholders := make([]string, len(filter.Tags))
		for i, tag := range filter.Tags {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, tag)
			argIndex++
		}
		
		operator := "?&" // AND演算子
		if filter.TagsOperator == "OR" {
			operator = "?|" // OR演算子
		}
		conditions = append(conditions, fmt.Sprintf("tags %s ARRAY[%s]::text[]", operator, strings.Join(placeholders, ",")))
	}

	// 近隣検索（緯度・経度・半径指定）
	if filter.Latitude != nil && filter.Longitude != nil && filter.Radius != nil {
		conditions = append(conditions, fmt.Sprintf(`
			(6371 * acos(cos(radians($%d)) * cos(radians(latitude)) * 
			cos(radians(longitude) - radians($%d)) + sin(radians($%d)) * 
			sin(radians(latitude)))) <= $%d`, argIndex, argIndex+1, argIndex, argIndex+2))
		args = append(args, *filter.Latitude, *filter.Longitude, *filter.Latitude, *filter.Radius)
		argIndex += 3
	}

	// 営業時間検索（JSON形式対応）
	if filter.BusinessDay != "" && filter.BusinessTime != "" {
		// 両方指定: 指定曜日の指定時間に営業している店
		conditions = append(conditions, fmt.Sprintf(`
			(business_hours IS NOT NULL AND
			 business_hours->$%d->>'is_closed' != 'true' AND
			 EXISTS (
			   SELECT 1 FROM jsonb_array_elements(business_hours->$%d->'time_slots') AS slot
			   WHERE $%d::time >= (slot->>'open_time')::time 
			     AND $%d::time <= COALESCE(
			       NULLIF(slot->>'last_order_time', ''),
			       slot->>'close_time'
			     )::time
			 ))`, argIndex, argIndex+1, argIndex+2, argIndex+3))
		args = append(args, filter.BusinessDay, filter.BusinessDay, filter.BusinessTime, filter.BusinessTime)
		argIndex += 4
	} else if filter.BusinessDay != "" {
		// 営業日のみ指定: その日が休業日ではない店
		conditions = append(conditions, fmt.Sprintf(`
			(business_hours IS NULL OR 
			 business_hours->$%d->>'is_closed' != 'true')`, argIndex))
		args = append(args, filter.BusinessDay)
		argIndex++
	} else if filter.BusinessTime != "" {
		// 営業時間のみ指定: 1週間のうちどこか1日でもその時間に営業している店
		timeConditions := []string{}
		for _, day := range []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"} {
			timeConditions = append(timeConditions, fmt.Sprintf(`
				(business_hours->$%d->>'is_closed' != 'true' AND
				 EXISTS (
				   SELECT 1 FROM jsonb_array_elements(business_hours->$%d->'time_slots') AS slot
				   WHERE $%d::time >= (slot->>'open_time')::time 
				     AND $%d::time <= COALESCE(
				       NULLIF(slot->>'last_order_time', ''),
				       slot->>'close_time'
					     )::time
				 ))`, argIndex, argIndex+1, argIndex+2, argIndex+3))
			args = append(args, day, day, filter.BusinessTime, filter.BusinessTime)
			argIndex += 4
		}
		conditions = append(conditions, "("+strings.Join(timeConditions, " OR ")+")")
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int
	err := r.db.QueryRow(baseQuery, args...).Scan(&count)
	return count, err
}

func (r *StoreRepository) Update(store *models.Store) error {
	query := `
		UPDATE stores SET 
			name = $2, address = $3, latitude = $4, longitude = $5, categories = $6,
			business_hours = $7, parking_info = $8, website_url = $9,
			google_map_url = $10, sns_urls = $11, tags = $12, photos = $13, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(query,
		store.ID, store.Name, store.Address, store.Latitude, store.Longitude,
		store.Categories, store.BusinessHours, store.ParkingInfo,
		store.WebsiteURL, store.GoogleMapURL, store.SnsUrls, store.Tags, store.Photos,
	)
	return err
}

func (r *StoreRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM stores WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *StoreRepository) GetAllCategories() ([]string, error) {
	query := `
		SELECT DISTINCT jsonb_array_elements_text(categories) as category 
		FROM stores 
		WHERE categories IS NOT NULL AND jsonb_array_length(categories) > 0
		ORDER BY category
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *StoreRepository) GetAllTags() ([]string, error) {
	query := `
		SELECT DISTINCT jsonb_array_elements_text(tags) as tag 
		FROM stores 
		WHERE tags IS NOT NULL AND jsonb_array_length(tags) > 0
		ORDER BY tag
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

