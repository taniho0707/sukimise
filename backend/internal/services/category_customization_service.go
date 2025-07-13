package services

import (
	"fmt"
	"sukimise/internal/models"
	"sukimise/internal/repositories"
)

type CategoryCustomizationService struct {
	categoryCustomizationRepo repositories.CategoryCustomizationRepositoryInterface
}

func NewCategoryCustomizationService(categoryCustomizationRepo repositories.CategoryCustomizationRepositoryInterface) *CategoryCustomizationService {
	return &CategoryCustomizationService{
		categoryCustomizationRepo: categoryCustomizationRepo,
	}
}

func (s *CategoryCustomizationService) CreateCategoryCustomization(req *models.CategoryCustomizationRequest) (*models.CategoryCustomization, error) {
	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Check if category customization already exists
	existing, err := s.categoryCustomizationRepo.GetByCategoryName(req.CategoryName)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("category customization already exists for this category")
	}

	// Create new category customization
	categoryCustomization := &models.CategoryCustomization{
		CategoryName: req.CategoryName,
		Icon:         req.Icon,
		Color:        req.Color,
	}

	err = s.categoryCustomizationRepo.Create(categoryCustomization)
	if err != nil {
		return nil, err
	}

	return categoryCustomization, nil
}

func (s *CategoryCustomizationService) GetCategoryCustomization(categoryName string) (*models.CategoryCustomization, error) {
	return s.categoryCustomizationRepo.GetByCategoryName(categoryName)
}

func (s *CategoryCustomizationService) GetAllCategoryCustomizations() ([]*models.CategoryCustomization, error) {
	return s.categoryCustomizationRepo.GetAll()
}

func (s *CategoryCustomizationService) UpdateCategoryCustomization(categoryName string, req *models.CategoryCustomizationRequest) (*models.CategoryCustomization, error) {
	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Check if category customization exists
	existing, err := s.categoryCustomizationRepo.GetByCategoryName(categoryName)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("category customization not found")
	}

	// Update the category customization
	existing.Icon = req.Icon
	existing.Color = req.Color

	err = s.categoryCustomizationRepo.Update(existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *CategoryCustomizationService) DeleteCategoryCustomization(categoryName string) error {
	// Check if category customization exists
	existing, err := s.categoryCustomizationRepo.GetByCategoryName(categoryName)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("category customization not found")
	}

	return s.categoryCustomizationRepo.Delete(categoryName)
}

// SyncWithStoreCategories ensures all store categories have corresponding customizations
func (s *CategoryCustomizationService) SyncWithStoreCategories(storeCategories []string) error {
	// Get existing customizations
	existingCustomizations, err := s.categoryCustomizationRepo.GetAll()
	if err != nil {
		return err
	}

	// Create a map for fast lookup
	existingMap := make(map[string]bool)
	for _, customization := range existingCustomizations {
		existingMap[customization.CategoryName] = true
	}

	// Define default icons and colors for common categories
	defaultCustomizations := map[string]struct {
		Icon  string
		Color string
	}{
		"掃除":  {"🧹", "#9E9E9E"},
		"洗濯":  {"👕", "#03A9F4"},
		"牛丼":  {"🍚", "#FF9800"},
		"軽食":  {"🥪", "#FFC107"},
		"その他": {"📍", "#607D8B"},
	}

	// Add missing customizations
	for _, category := range storeCategories {
		if !existingMap[category] {
			// Use default customization if available, otherwise create a basic one
			icon := "📍"
			color := "#607D8B"
			
			if defaults, exists := defaultCustomizations[category]; exists {
				icon = defaults.Icon
				color = defaults.Color
			}

			customization := &models.CategoryCustomization{
				CategoryName: category,
				Icon:         &icon,
				Color:        &color,
			}

			err := s.categoryCustomizationRepo.Create(customization)
			if err != nil {
				return fmt.Errorf("failed to create customization for category '%s': %v", category, err)
			}
		}
	}

	return nil
}