package handlers

import (
	"log"
	"net/http"
	"sukimise/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateReviewRequest struct {
	StoreID       uuid.UUID  `json:"store_id" binding:"required"`
	Rating        int        `json:"rating" binding:"required,min=1,max=5"`
	Comment       *string    `json:"comment"`
	Photos        []string   `json:"photos"`
	VisitDate     *time.Time `json:"visit_date"`
	IsVisited     bool       `json:"is_visited"`
	PaymentAmount *int       `json:"payment_amount"`
	FoodNotes     *string    `json:"food_notes"`
}

type UpdateReviewRequest struct {
	Rating        int        `json:"rating" binding:"min=1,max=5"`
	Comment       *string    `json:"comment"`
	Photos        []string   `json:"photos"`
	VisitDate     *time.Time `json:"visit_date"`
	IsVisited     bool       `json:"is_visited"`
	PaymentAmount *int       `json:"payment_amount"`
	FoodNotes     *string    `json:"food_notes"`
}

func (h *Handler) CreateReview(c *gin.Context) {
	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	review := &models.Review{
		StoreID:       req.StoreID,
		UserID:        userID.(uuid.UUID),
		Rating:        req.Rating,
		Comment:       req.Comment,
		Photos:        models.StringArray(req.Photos),
		VisitDate:     req.VisitDate,
		IsVisited:     req.IsVisited,
		PaymentAmount: req.PaymentAmount,
		FoodNotes:     req.FoodNotes,
	}

	if err := h.reviewService.CreateReview(review); err != nil {
		if err.Error() == "review already exists for this store and user" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, review)
}

func (h *Handler) GetReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	review, err := h.reviewService.GetReviewByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	c.JSON(http.StatusOK, review)
}

func (h *Handler) GetReviewsByStore(c *gin.Context) {
	storeIDStr := c.Param("id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	reviews, err := h.reviewService.GetReviewsByStoreID(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

func (h *Handler) GetMyReviews(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		log.Println("ERROR: User ID not found in token context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	log.Printf("DEBUG: Getting reviews for user ID: %v", userID)
	
	reviews, err := h.reviewService.GetReviewsByUserID(userID.(uuid.UUID))
	if err != nil {
		log.Printf("ERROR: Failed to get reviews for user %v: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reviews"})
		return
	}

	log.Printf("DEBUG: Successfully retrieved %d reviews for user %v", len(reviews), userID)
	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

func (h *Handler) UpdateReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var req UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	existingReview, err := h.reviewService.GetReviewByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	if req.Rating > 0 {
		existingReview.Rating = req.Rating
	}
	if req.Comment != nil && *req.Comment != "" {
		existingReview.Comment = req.Comment
	}
	if req.Photos != nil {
		existingReview.Photos = models.StringArray(req.Photos)
	}
	if req.VisitDate != nil {
		existingReview.VisitDate = req.VisitDate
	}
	existingReview.IsVisited = req.IsVisited
	existingReview.PaymentAmount = req.PaymentAmount
	if req.FoodNotes != nil {
		existingReview.FoodNotes = req.FoodNotes
	}

	if err := h.reviewService.UpdateReview(existingReview, userID.(uuid.UUID)); err != nil {
		if err.Error() == "unauthorized: you can only update your own reviews" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
		return
	}

	c.JSON(http.StatusOK, existingReview)
}

func (h *Handler) DeleteReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	if err := h.reviewService.DeleteReview(id, userID.(uuid.UUID)); err != nil {
		if err.Error() == "unauthorized: you can only delete your own reviews" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}