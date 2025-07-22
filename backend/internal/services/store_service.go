package services

import (
	"sukimise/internal/models"
	"sukimise/internal/repositories"

	"github.com/google/uuid"
)

type StoreService struct {
	storeRepo repositories.StoreRepositoryInterface
}

func NewStoreService(storeRepo repositories.StoreRepositoryInterface) *StoreService {
	return &StoreService{storeRepo: storeRepo}
}

func (s *StoreService) CreateStore(store *models.Store) error {
	return s.storeRepo.Create(store)
}

func (s *StoreService) GetStoreByID(id uuid.UUID) (*models.Store, error) {
	return s.storeRepo.GetByID(id)
}

func (s *StoreService) GetStores(filter *repositories.StoreFilter) ([]*models.Store, error) {
	if filter == nil {
		filter = &repositories.StoreFilter{}
	}
	
	if filter.Limit == 0 {
		filter.Limit = 20
	}

	return s.storeRepo.GetAll(filter)
}

func (s *StoreService) GetStoresCount(filter *repositories.StoreFilter) (int, error) {
	if filter == nil {
		filter = &repositories.StoreFilter{}
	}
	
	return s.storeRepo.GetCount(filter)
}

func (s *StoreService) UpdateStore(store *models.Store) error {
	return s.storeRepo.Update(store)
}

func (s *StoreService) DeleteStore(id uuid.UUID) error {
	return s.storeRepo.Delete(id)
}

func (s *StoreService) SearchStores(name string, categories, tags []string, lat, lng, radius *float64, limit, offset int) ([]*models.Store, error) {
	filter := &repositories.StoreFilter{
		Name:       name,
		Categories: categories,
		Tags:       tags,
		Latitude:   lat,
		Longitude:  lng,
		Radius:     radius,
		Limit:      limit,
		Offset:     offset,
	}
	
	if filter.Limit == 0 {
		filter.Limit = 20
	}

	return s.storeRepo.GetAll(filter)
}

func (s *StoreService) GetAllCategories() ([]string, error) {
	return s.storeRepo.GetAllCategories()
}

func (s *StoreService) GetAllTags() ([]string, error) {
	return s.storeRepo.GetAllTags()
}

// CheckForDuplicate checks if a store with the same name and location (within 50m) already exists
func (s *StoreService) CheckForDuplicate(name string, latitude, longitude float64) (*models.Store, error) {
	return s.storeRepo.FindDuplicateByLocationAndName(name, latitude, longitude)
}

