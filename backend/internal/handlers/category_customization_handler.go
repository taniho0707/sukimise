package handlers

import (
	"log"
	"net/http"
	"sukimise/internal/errors"
	"sukimise/internal/models"
	"sukimise/internal/services"

	"github.com/gin-gonic/gin"
)

type CategoryCustomizationHandler struct {
	categoryCustomizationService *services.CategoryCustomizationService
	storeService                 *services.StoreService
}

func NewCategoryCustomizationHandler(categoryCustomizationService *services.CategoryCustomizationService, storeService *services.StoreService) *CategoryCustomizationHandler {
	return &CategoryCustomizationHandler{
		categoryCustomizationService: categoryCustomizationService,
		storeService:                 storeService,
	}
}

// CreateCategoryCustomization creates a new category customization
func (h *CategoryCustomizationHandler) CreateCategoryCustomization(c *gin.Context) {
	var req models.CategoryCustomizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.NewValidationError("Invalid request body", err.Error()))
		return
	}

	categoryCustomization, err := h.categoryCustomizationService.CreateCategoryCustomization(&req)
	if err != nil {
		log.Printf("Failed to create category customization: %v", err)
		errors.HandleError(c, err)
		return
	}

	errors.SendSuccess(c, categoryCustomization, nil)
}

// GetCategoryCustomizations returns all category customizations
func (h *CategoryCustomizationHandler) GetCategoryCustomizations(c *gin.Context) {
	customizations, err := h.categoryCustomizationService.GetAllCategoryCustomizations()
	if err != nil {
		log.Printf("Failed to get category customizations: %v", err)
		errors.HandleError(c, errors.NewInternalError("Failed to get category customizations"))
		return
	}

	errors.SendSuccess(c, map[string]interface{}{
		"category_customizations": customizations,
	}, nil)
}

// GetCategoryCustomization returns a specific category customization
func (h *CategoryCustomizationHandler) GetCategoryCustomization(c *gin.Context) {
	categoryName := c.Param("categoryName")
	if categoryName == "" {
		errors.HandleError(c, errors.NewValidationError("Category name is required", ""))
		return
	}

	customization, err := h.categoryCustomizationService.GetCategoryCustomization(categoryName)
	if err != nil {
		log.Printf("Failed to get category customization: %v", err)
		errors.HandleError(c, errors.NewInternalError("Failed to get category customization"))
		return
	}

	if customization == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Category customization not found",
		})
		return
	}

	errors.SendSuccess(c, customization, nil)
}

// UpdateCategoryCustomization updates an existing category customization
func (h *CategoryCustomizationHandler) UpdateCategoryCustomization(c *gin.Context) {
	categoryName := c.Param("categoryName")
	if categoryName == "" {
		errors.HandleError(c, errors.NewValidationError("Category name is required", ""))
		return
	}

	var req models.CategoryCustomizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.NewValidationError("Invalid request body", err.Error()))
		return
	}

	customization, err := h.categoryCustomizationService.UpdateCategoryCustomization(categoryName, &req)
	if err != nil {
		log.Printf("Failed to update category customization: %v", err)
		errors.HandleError(c, err)
		return
	}

	errors.SendSuccess(c, customization, nil)
}

// DeleteCategoryCustomization deletes a category customization
func (h *CategoryCustomizationHandler) DeleteCategoryCustomization(c *gin.Context) {
	categoryName := c.Param("categoryName")
	if categoryName == "" {
		errors.HandleError(c, errors.NewValidationError("Category name is required", ""))
		return
	}

	err := h.categoryCustomizationService.DeleteCategoryCustomization(categoryName)
	if err != nil {
		log.Printf("Failed to delete category customization: %v", err)
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category customization deleted successfully",
	})
}

// SyncCategoriesWithStores synchronizes category customizations with store categories
func (h *CategoryCustomizationHandler) SyncCategoriesWithStores(c *gin.Context) {
	// Get all store categories
	storeCategories, err := h.storeService.GetAllCategories()
	if err != nil {
		log.Printf("Failed to get store categories: %v", err)
		errors.HandleError(c, errors.NewInternalError("Failed to get store categories"))
		return
	}

	// Sync category customizations
	err = h.categoryCustomizationService.SyncWithStoreCategories(storeCategories)
	if err != nil {
		log.Printf("Failed to sync category customizations: %v", err)
		errors.HandleError(c, errors.NewInternalError("Failed to sync category customizations"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category customizations synchronized successfully",
		"data": map[string]interface{}{
			"synchronized_categories": storeCategories,
		},
	})
}