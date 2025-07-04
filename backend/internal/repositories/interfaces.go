package repositories

import (
	"sukimise/internal/models"

	"github.com/google/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock_repositories.go

type StoreRepositoryInterface interface {
	Create(store *models.Store) error
	GetByID(id uuid.UUID) (*models.Store, error)
	GetAll(filter *StoreFilter) ([]*models.Store, error)
	Update(store *models.Store) error
	Delete(id uuid.UUID) error
	GetAllCategories() ([]string, error)
	GetAllTags() ([]string, error)
}

type UserRepositoryInterface interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	GetAll() ([]*models.User, error)
}

type ReviewRepositoryInterface interface {
	Create(review *models.Review) error
	GetByID(id uuid.UUID) (*models.Review, error)
	GetByStoreID(storeID uuid.UUID, limit, offset int) ([]*models.Review, int, error)
	GetByUserID(userID uuid.UUID) ([]*models.Review, error)
	Update(review *models.Review) error
	Delete(id uuid.UUID) error
}

type ViewerAuthRepositoryInterface interface {
	GetViewerSettings() (*models.ViewerSettings, error)
	UpdateViewerSettings(settings *models.ViewerSettings) error
	CreateLoginHistory(history *models.ViewerLoginHistory) error
	GetValidSession(token string) (*models.ViewerLoginHistory, error)
	InvalidateSession(token string) error
	CleanupExpiredSessions() error
	GetLoginHistory(limit, offset int) ([]*models.ViewerLoginHistory, error)
}