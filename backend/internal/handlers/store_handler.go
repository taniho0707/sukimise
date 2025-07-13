package handlers

import (
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sukimise/internal/constants"
	"sukimise/internal/errors"
	"sukimise/internal/models"
	"sukimise/internal/repositories"
	"sukimise/internal/types"
	"sukimise/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StoreRequest represents the request body for store operations
type StoreRequest struct {
	Name          string                   `json:"name"`
	Address       string                   `json:"address"`
	Latitude      float64                  `json:"latitude"`
	Longitude     float64                  `json:"longitude"`
	Categories    []string                 `json:"categories"`
	BusinessHours models.BusinessHoursData `json:"business_hours"`
	ParkingInfo   string                   `json:"parking_info"`
	WebsiteURL    string                   `json:"website_url"`
	GoogleMapURL  string                   `json:"google_map_url"`
	SnsUrls       []string                 `json:"sns_urls"`
	Tags          []string                 `json:"tags"`
	Photos        []string                 `json:"photos"`
}

// ValidateForCreate validates store data for creation
func (r *StoreRequest) ValidateForCreate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.NewValidationError("Store name is required", "")
	}
	if strings.TrimSpace(r.Address) == "" {
		return errors.NewValidationError("Store address is required", "")
	}
	return r.validate()
}

// ValidateForUpdate validates store data for update
func (r *StoreRequest) ValidateForUpdate() error {
	return r.validate()
}

// validate performs common validation for both create and update
func (r *StoreRequest) validate() error {
	if err := utils.ValidateCoordinates(r.Latitude, r.Longitude); err != nil {
		return err
	}
	if err := utils.ValidateURL(r.WebsiteURL); err != nil {
		return err
	}
	if err := utils.ValidateURL(r.GoogleMapURL); err != nil {
		return err
	}
	if err := utils.ValidateStringArray(r.Categories, "categories", 10); err != nil {
		return err
	}
	if err := utils.ValidateStringArray(r.Tags, "tags", 20); err != nil {
		return err
	}
	if err := utils.ValidateStringArray(r.Photos, "photos", 20); err != nil {
		return err
	}
	for _, url := range r.SnsUrls {
		if err := utils.ValidateURL(url); err != nil {
			return err
		}
	}
	return nil
}

// ToModel converts StoreRequest to Store model
func (r *StoreRequest) ToModel(userID uuid.UUID) *models.Store {
	return &models.Store{
		Name:          r.Name,
		Address:       r.Address,
		Latitude:      r.Latitude,
		Longitude:     r.Longitude,
		Categories:    models.StringArray(r.Categories),
		BusinessHours: r.BusinessHours,
		ParkingInfo:   r.ParkingInfo,
		WebsiteURL:    r.WebsiteURL,
		GoogleMapURL:  r.GoogleMapURL,
		SnsUrls:       models.StringArray(r.SnsUrls),
		Tags:          models.StringArray(r.Tags),
		Photos:        models.StringArray(r.Photos),
		CreatedBy:     userID,
	}
}

// UpdateModel updates existing store model with request data
func (r *StoreRequest) UpdateModel(store *models.Store) {
	if r.Name != "" {
		store.Name = r.Name
	}
	if r.Address != "" {
		store.Address = r.Address
	}
	// Always update latitude and longitude (including 0 values)
	store.Latitude = r.Latitude
	store.Longitude = r.Longitude
	if r.Categories != nil {
		store.Categories = models.StringArray(r.Categories)
	}
	// Always update business hours since it's a structured object
	store.BusinessHours = r.BusinessHours
	if r.ParkingInfo != "" {
		store.ParkingInfo = r.ParkingInfo
	}
	if r.WebsiteURL != "" {
		store.WebsiteURL = r.WebsiteURL
	}
	if r.GoogleMapURL != "" {
		store.GoogleMapURL = r.GoogleMapURL
	}
	if r.SnsUrls != nil {
		store.SnsUrls = models.StringArray(r.SnsUrls)
	}
	if r.Tags != nil {
		store.Tags = models.StringArray(r.Tags)
	}
	if r.Photos != nil {
		store.Photos = models.StringArray(r.Photos)
	}
}

func (h *Handler) CreateStore(c *gin.Context) {
	var req StoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.NewValidationError("Invalid request data", err.Error()))
		return
	}

	// Validate request data
	if err := req.ValidateForCreate(); err != nil {
		errors.HandleError(c, err)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		errors.HandleError(c, errors.NewUnauthorizedError("User ID not found in token"))
		return
	}

	store := req.ToModel(userID.(uuid.UUID))

	if err := h.storeService.CreateStore(store); err != nil {
		log.Printf("Failed to create store: %v", err)
		errors.HandleError(c, errors.NewInternalError("Failed to create store"))
		return
	}

	errors.SendCreated(c, store)
}

func (h *Handler) GetStore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		errors.HandleError(c, errors.NewValidationError("Invalid store ID", err.Error()))
		return
	}

	store, err := h.storeService.GetStoreByID(id)
	if err != nil {
		log.Printf("Store not found: %v", err)
		errors.HandleError(c, errors.NewNotFoundError("Store"))
		return
	}

	errors.SendSuccess(c, store)
}

func (h *Handler) GetStores(c *gin.Context) {
	filter := h.parseStoreFilter(c)

	// Validate filter parameters
	if err := h.validateStoreFilter(filter); err != nil {
		errors.HandleError(c, err)
		return
	}

	stores, err := h.storeService.GetStores(filter)
	if err != nil {
		log.Printf("Failed to get stores: %v", err)
		errors.HandleError(c, errors.NewInternalError("Failed to get stores"))
		return
	}

	// Get total count for pagination
	totalCount, err := h.storeService.GetStoresCount(filter)
	if err != nil {
		log.Printf("Failed to get stores count: %v", err)
		errors.HandleError(c, errors.NewInternalError("Failed to get stores count"))
		return
	}

	// Calculate pagination info
	totalPages := (totalCount + filter.Limit - 1) / filter.Limit
	currentPage := (filter.Offset / filter.Limit) + 1

	// Prepare metadata
	meta := &types.MetaInfo{
		Total:       totalCount,
		Limit:       filter.Limit,
		Offset:      filter.Offset,
		Page:        &currentPage,
		TotalPages:  &totalPages,
	}

	errors.SendSuccess(c, map[string]interface{}{
		"stores": stores,
	}, meta)
}


func (h *Handler) UpdateStore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		errors.HandleError(c, errors.NewValidationError("Invalid store ID", err.Error()))
		return
	}

	var req StoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.NewValidationError("Invalid request data", err.Error()))
		return
	}

	// Validate request data
	if err := req.ValidateForUpdate(); err != nil {
		errors.HandleError(c, err)
		return
	}

	existingStore, err := h.storeService.GetStoreByID(id)
	if err != nil {
		log.Printf("Store not found for update: %v", err)
		errors.HandleError(c, errors.NewNotFoundError("Store"))
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		errors.HandleError(c, errors.NewUnauthorizedError("User ID not found in token"))
		return
	}

	userRole, exists := c.Get("role")
	if !exists {
		errors.HandleError(c, errors.NewUnauthorizedError("User role not found in token"))
		return
	}

	// Check permissions
	if existingStore.CreatedBy != userID.(uuid.UUID) && userRole.(string) != constants.RoleAdmin {
		errors.HandleError(c, errors.NewForbiddenError("You can only update stores you created"))
		return
	}

	// Update store with request data
	req.UpdateModel(existingStore)

	if err := h.storeService.UpdateStore(existingStore); err != nil {
		log.Printf("Failed to update store: %v", err)
		errors.HandleError(c, errors.NewInternalError("Failed to update store"))
		return
	}

	errors.SendSuccess(c, existingStore)
}

func (h *Handler) DeleteStore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		errors.HandleError(c, errors.NewValidationError("Invalid store ID", err.Error()))
		return
	}

	existingStore, err := h.storeService.GetStoreByID(id)
	if err != nil {
		log.Printf("Store not found for deletion: %v", err)
		errors.HandleError(c, errors.NewNotFoundError("Store"))
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		errors.HandleError(c, errors.NewUnauthorizedError("User ID not found in token"))
		return
	}

	userRole, exists := c.Get("role")
	if !exists {
		errors.HandleError(c, errors.NewUnauthorizedError("User role not found in token"))
		return
	}

	// Check permissions
	if existingStore.CreatedBy != userID.(uuid.UUID) && userRole.(string) != constants.RoleAdmin {
		errors.HandleError(c, errors.NewForbiddenError("You can only delete stores you created"))
		return
	}

	if err := h.storeService.DeleteStore(id); err != nil {
		log.Printf("Failed to delete store: %v", err)
		errors.HandleError(c, errors.NewInternalError("Failed to delete store"))
		return
	}

	errors.SendSuccess(c, map[string]string{"message": "Store deleted successfully"})
}

func (h *Handler) GetCategories(c *gin.Context) {
	categories, err := h.storeService.GetAllCategories()
	if err != nil {
		log.Printf("Failed to get categories: %v", err)
		errors.HandleError(c, errors.NewInternalError("Failed to get categories"))
		return
	}

	log.Printf("Found %d categories: %v", len(categories), categories)
	errors.SendSuccess(c, map[string][]string{"categories": categories})
}

func (h *Handler) GetTags(c *gin.Context) {
	tags, err := h.storeService.GetAllTags()
	if err != nil {
		log.Printf("Failed to get tags: %v", err)
		errors.HandleError(c, errors.NewInternalError("Failed to get tags"))
		return
	}

	log.Printf("Found %d tags: %v", len(tags), tags)
	errors.SendSuccess(c, map[string][]string{"tags": tags})
}

// ExportStoresCSV handles CSV export of stores with current filter conditions
func (h *Handler) ExportStoresCSV(c *gin.Context) {
	filter := h.parseStoreFilter(c)

	// Validate filter parameters
	if err := h.validateStoreFilter(filter); err != nil {
		errors.HandleError(c, err)
		return
	}

	stores, err := h.storeService.GetStores(filter)
	if err != nil {
		log.Printf("Failed to get stores for CSV export: %v", err)
		errors.HandleError(c, errors.NewInternalError("Failed to get stores"))
		return
	}

	// Set response headers for CSV download
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("sukimise_stores_%s.csv", timestamp)
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// Create CSV writer with proper configuration
	writer := csv.NewWriter(c.Writer)
	writer.UseCRLF = false // Use LF only for better compatibility
	defer writer.Flush()

	// Write BOM for UTF-8 (Excel compatibility)
	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})

	// Write CSV header
	header := []string{
		"ID",
		"店名",
		"住所",
		"緯度",
		"経度",
		"カテゴリ",
		"営業時間",
		"駐車場情報",
		"ウェブサイト",
		"GoogleMap URL",
		"SNS URL",
		"タグ",
		"作成者",
		"作成日時",
		"更新日時",
	}
	if err := writer.Write(header); err != nil {
		log.Printf("Failed to write CSV header: %v", err)
		errors.HandleError(c, errors.NewInternalError("Failed to generate CSV"))
		return
	}

	// Write store data
	for _, store := range stores {
		record := []string{
			store.ID.String(),
			sanitizeCSVField(store.Name),
			sanitizeCSVField(store.Address),
			fmt.Sprintf("%.6f", store.Latitude),
			fmt.Sprintf("%.6f", store.Longitude),
			sanitizeCSVField(strings.Join(store.Categories, "; ")),
			sanitizeCSVField(formatBusinessHoursForCSV(store.BusinessHours)),
			sanitizeCSVField(store.ParkingInfo),
			sanitizeCSVField(store.WebsiteURL),
			sanitizeCSVField(store.GoogleMapURL),
			sanitizeCSVField(strings.Join(store.SnsUrls, "; ")),
			sanitizeCSVField(strings.Join(store.Tags, "; ")),
			store.CreatedBy.String(),
			store.CreatedAt.Format("2006-01-02 15:04:05"),
			store.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(record); err != nil {
			log.Printf("Failed to write CSV record: %v", err)
			errors.HandleError(c, errors.NewInternalError("Failed to generate CSV"))
			return
		}
	}

	log.Printf("CSV export completed: %d stores exported", len(stores))
}

// parseStoreFilter parses query parameters into StoreFilter
func (h *Handler) parseStoreFilter(c *gin.Context) *repositories.StoreFilter {
	filter := &repositories.StoreFilter{
		Limit: constants.DefaultLimit,
	}

	if name := c.Query("name"); name != "" {
		filter.Name = strings.TrimSpace(name)
	}

	if categories := c.Query("categories"); categories != "" {
		filter.Categories = strings.Split(categories, ",")
		for i, cat := range filter.Categories {
			filter.Categories[i] = strings.TrimSpace(cat)
		}
	}

	if categoriesOp := c.Query("categories_operator"); categoriesOp != "" {
		filter.CategoriesOperator = categoriesOp
	}

	if tags := c.Query("tags"); tags != "" {
		filter.Tags = strings.Split(tags, ",")
		for i, tag := range filter.Tags {
			filter.Tags[i] = strings.TrimSpace(tag)
		}
	}

	if tagsOp := c.Query("tags_operator"); tagsOp != "" {
		filter.TagsOperator = tagsOp
	}

	if latStr := c.Query("latitude"); latStr != "" {
		if lat, err := strconv.ParseFloat(latStr, 64); err == nil {
			filter.Latitude = &lat
		}
	}

	if lngStr := c.Query("longitude"); lngStr != "" {
		if lng, err := strconv.ParseFloat(lngStr, 64); err == nil {
			filter.Longitude = &lng
		}
	}

	if radiusStr := c.Query("radius"); radiusStr != "" {
		if radius, err := strconv.ParseFloat(radiusStr, 64); err == nil {
			filter.Radius = &radius
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			if limit > constants.MaxLimit {
				limit = constants.MaxLimit
			}
			filter.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	if businessDay := c.Query("business_day"); businessDay != "" {
		filter.BusinessDay = strings.TrimSpace(businessDay)
	}

	if businessTime := c.Query("business_time"); businessTime != "" {
		filter.BusinessTime = strings.TrimSpace(businessTime)
	}

	return filter
}

// validateStoreFilter validates store filter parameters
func (h *Handler) validateStoreFilter(filter *repositories.StoreFilter) error {
	// Validate coordinates if provided
	if filter.Latitude != nil && filter.Longitude != nil {
		if err := utils.ValidateCoordinates(*filter.Latitude, *filter.Longitude); err != nil {
			return err
		}
	}

	// Validate business day if provided
	if filter.BusinessDay != "" {
		if err := utils.ValidateBusinessDay(filter.BusinessDay); err != nil {
			return err
		}
	}

	// Validate radius requirements
	if filter.Radius != nil && (filter.Latitude == nil || filter.Longitude == nil) {
		return errors.NewValidationError(
			"Invalid location search",
			"Latitude and longitude are required when radius is specified",
		)
	}

	return nil
}

// sanitizeCSVField sanitizes a field for CSV output by replacing problematic characters
func sanitizeCSVField(field string) string {
	if field == "" {
		return field
	}
	
	// Replace line breaks with space to prevent CSV structure issues
	field = strings.ReplaceAll(field, "\n", " ")
	field = strings.ReplaceAll(field, "\r\n", " ")
	field = strings.ReplaceAll(field, "\r", " ")
	
	// Replace tab characters with space
	field = strings.ReplaceAll(field, "\t", " ")
	
	// Replace multiple spaces with single space
	field = strings.Join(strings.Fields(field), " ")
	
	// Trim leading and trailing whitespace
	field = strings.TrimSpace(field)
	
	return field
}

// formatBusinessHoursForCSV converts BusinessHoursData to CSV-friendly string
func formatBusinessHoursForCSV(businessHours models.BusinessHoursData) string {
	var parts []string
	
	dayNames := map[string]string{
		"monday": "月", "tuesday": "火", "wednesday": "水", "thursday": "木",
		"friday": "金", "saturday": "土", "sunday": "日",
	}
	
	allDays := []struct {
		key      string
		schedule models.DaySchedule
	}{
		{"monday", businessHours.Monday},
		{"tuesday", businessHours.Tuesday},
		{"wednesday", businessHours.Wednesday},
		{"thursday", businessHours.Thursday},
		{"friday", businessHours.Friday},
		{"saturday", businessHours.Saturday},
		{"sunday", businessHours.Sunday},
	}
	
	// 営業時間の抽出
	var openSlots []models.TimeSlot
	var closedDays []string
	
	for _, day := range allDays {
		if day.schedule.IsClosed {
			closedDays = append(closedDays, dayNames[day.key])
		} else if len(day.schedule.TimeSlots) > 0 {
			openSlots = append(openSlots, day.schedule.TimeSlots...)
		}
	}
	
	// 共通の営業時間を見つける
	if len(openSlots) > 0 {
		// 最初のスロットを基準にする
		firstSlot := openSlots[0]
		parts = append(parts, fmt.Sprintf("営業時間: %s-%s", firstSlot.OpenTime, firstSlot.CloseTime))
		
		if firstSlot.LastOrderTime != "" && firstSlot.LastOrderTime != firstSlot.CloseTime {
			parts = append(parts, fmt.Sprintf("ラストオーダー: %s", firstSlot.LastOrderTime))
		}
	}
	
	// 定休日の追加
	if len(closedDays) > 0 {
		parts = append(parts, fmt.Sprintf("定休日: %s", strings.Join(closedDays, "、")))
	} else if len(openSlots) > 0 {
		parts = append(parts, "定休日: 年中無休")
	}
	
	if len(parts) == 0 {
		return "営業時間未設定"
	}
	
	return strings.Join(parts, "; ")
}

