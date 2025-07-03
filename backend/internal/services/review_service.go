package services

import (
	"errors"
	"sukimise/internal/models"
	"sukimise/internal/repositories"

	"github.com/google/uuid"
)

type ReviewService struct {
	reviewRepo *repositories.ReviewRepository
}

func NewReviewService(reviewRepo *repositories.ReviewRepository) *ReviewService {
	return &ReviewService{reviewRepo: reviewRepo}
}

func (s *ReviewService) CreateReview(review *models.Review) error {
	return s.reviewRepo.Create(review)
}

func (s *ReviewService) GetReviewByID(id uuid.UUID) (*models.Review, error) {
	return s.reviewRepo.GetByID(id)
}

func (s *ReviewService) GetReviewsByStoreID(storeID uuid.UUID) ([]*models.Review, error) {
	return s.reviewRepo.GetByStoreID(storeID)
}

func (s *ReviewService) GetReviewsByUserID(userID uuid.UUID) ([]*models.Review, error) {
	return s.reviewRepo.GetByUserID(userID)
}

func (s *ReviewService) GetReviewByStoreAndUser(storeID, userID uuid.UUID) (*models.Review, error) {
	return s.reviewRepo.GetByStoreAndUser(storeID, userID)
}

func (s *ReviewService) UpdateReview(review *models.Review, userID uuid.UUID) error {
	existingReview, err := s.reviewRepo.GetByID(review.ID)
	if err != nil {
		return err
	}

	if existingReview.UserID != userID {
		return errors.New("unauthorized: you can only update your own reviews")
	}

	return s.reviewRepo.Update(review)
}

func (s *ReviewService) DeleteReview(id, userID uuid.UUID) error {
	existingReview, err := s.reviewRepo.GetByID(id)
	if err != nil {
		return err
	}

	if existingReview.UserID != userID {
		return errors.New("unauthorized: you can only delete your own reviews")
	}

	return s.reviewRepo.Delete(id)
}

func (s *ReviewService) CreateMenuItem(menuItem *models.MenuItem, userID uuid.UUID) error {
	review, err := s.reviewRepo.GetByID(menuItem.ReviewID)
	if err != nil {
		return err
	}

	if review.UserID != userID {
		return errors.New("unauthorized: you can only add menu items to your own reviews")
	}

	return s.reviewRepo.CreateMenuItem(menuItem)
}

func (s *ReviewService) GetMenuItemsByReviewID(reviewID uuid.UUID) ([]*models.MenuItem, error) {
	return s.reviewRepo.GetMenuItemsByReviewID(reviewID)
}