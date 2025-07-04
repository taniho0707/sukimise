package services

import (
	"errors"
	"sukimise/internal/models"
	"sukimise/internal/repositories"
	mocks "sukimise/internal/repositories/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestStoreService_CreateStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStoreRepositoryInterface(ctrl)
	service := NewStoreService(mockRepo)

	store := &models.Store{
		Name:    "Test Store",
		Address: "Test Address",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().Create(store).Return(nil)

		err := service.CreateStore(store)
		assert.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.EXPECT().Create(store).Return(errors.New("database error"))

		err := service.CreateStore(store)
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
	})
}

func TestStoreService_GetStoreByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStoreRepositoryInterface(ctrl)
	service := &StoreService{storeRepo: mockRepo}

	storeID := uuid.New()
	expectedStore := &models.Store{
		ID:      storeID,
		Name:    "Test Store",
		Address: "Test Address",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(storeID).Return(expectedStore, nil)

		store, err := service.GetStoreByID(storeID)
		assert.NoError(t, err)
		assert.Equal(t, expectedStore, store)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(storeID).Return(nil, errors.New("store not found"))

		store, err := service.GetStoreByID(storeID)
		assert.Error(t, err)
		assert.Nil(t, store)
	})
}

func TestStoreService_GetStores(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStoreRepositoryInterface(ctrl)
	service := &StoreService{storeRepo: mockRepo}

	expectedStores := []*models.Store{
		{ID: uuid.New(), Name: "Store 1"},
		{ID: uuid.New(), Name: "Store 2"},
	}

	t.Run("success with nil filter", func(t *testing.T) {
		expectedFilter := &repositories.StoreFilter{Limit: 20}
		mockRepo.EXPECT().GetAll(expectedFilter).Return(expectedStores, nil)

		stores, err := service.GetStores(nil)
		assert.NoError(t, err)
		assert.Equal(t, expectedStores, stores)
	})

	t.Run("success with custom filter", func(t *testing.T) {
		filter := &repositories.StoreFilter{
			Name:  "test",
			Limit: 10,
		}
		mockRepo.EXPECT().GetAll(filter).Return(expectedStores, nil)

		stores, err := service.GetStores(filter)
		assert.NoError(t, err)
		assert.Equal(t, expectedStores, stores)
	})

	t.Run("success with zero limit sets default", func(t *testing.T) {
		filter := &repositories.StoreFilter{
			Name:  "test",
			Limit: 0,
		}
		expectedFilter := &repositories.StoreFilter{
			Name:  "test",
			Limit: 20,
		}
		mockRepo.EXPECT().GetAll(expectedFilter).Return(expectedStores, nil)

		stores, err := service.GetStores(filter)
		assert.NoError(t, err)
		assert.Equal(t, expectedStores, stores)
	})
}

func TestStoreService_UpdateStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStoreRepositoryInterface(ctrl)
	service := &StoreService{storeRepo: mockRepo}

	store := &models.Store{
		ID:      uuid.New(),
		Name:    "Updated Store",
		Address: "Updated Address",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().Update(store).Return(nil)

		err := service.UpdateStore(store)
		assert.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.EXPECT().Update(store).Return(errors.New("update failed"))

		err := service.UpdateStore(store)
		assert.Error(t, err)
		assert.Equal(t, "update failed", err.Error())
	})
}

func TestStoreService_DeleteStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStoreRepositoryInterface(ctrl)
	service := &StoreService{storeRepo: mockRepo}

	storeID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().Delete(storeID).Return(nil)

		err := service.DeleteStore(storeID)
		assert.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.EXPECT().Delete(storeID).Return(errors.New("delete failed"))

		err := service.DeleteStore(storeID)
		assert.Error(t, err)
		assert.Equal(t, "delete failed", err.Error())
	})
}

func TestStoreService_SearchStores(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStoreRepositoryInterface(ctrl)
	service := &StoreService{storeRepo: mockRepo}

	expectedStores := []*models.Store{
		{ID: uuid.New(), Name: "Store 1"},
	}

	lat := 35.6762
	lng := 139.6503
	radius := 1000.0

	t.Run("success", func(t *testing.T) {
		expectedFilter := &repositories.StoreFilter{
			Name:       "test",
			Categories: []string{"restaurant"},
			Tags:       []string{"japanese"},
			Latitude:   &lat,
			Longitude:  &lng,
			Radius:     &radius,
			Limit:      10,
			Offset:     0,
		}

		mockRepo.EXPECT().GetAll(expectedFilter).Return(expectedStores, nil)

		stores, err := service.SearchStores("test", []string{"restaurant"}, []string{"japanese"}, &lat, &lng, &radius, 10, 0)
		assert.NoError(t, err)
		assert.Equal(t, expectedStores, stores)
	})

	t.Run("sets default limit when zero", func(t *testing.T) {
		expectedFilter := &repositories.StoreFilter{
			Name:       "test",
			Categories: []string{"restaurant"},
			Tags:       []string{"japanese"},
			Latitude:   &lat,
			Longitude:  &lng,
			Radius:     &radius,
			Limit:      20, // Default limit
			Offset:     0,
		}

		mockRepo.EXPECT().GetAll(expectedFilter).Return(expectedStores, nil)

		stores, err := service.SearchStores("test", []string{"restaurant"}, []string{"japanese"}, &lat, &lng, &radius, 0, 0)
		assert.NoError(t, err)
		assert.Equal(t, expectedStores, stores)
	})
}

func TestStoreService_GetAllCategories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStoreRepositoryInterface(ctrl)
	service := &StoreService{storeRepo: mockRepo}

	expectedCategories := []string{"restaurant", "cafe", "bar"}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().GetAllCategories().Return(expectedCategories, nil)

		categories, err := service.GetAllCategories()
		assert.NoError(t, err)
		assert.Equal(t, expectedCategories, categories)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.EXPECT().GetAllCategories().Return(nil, errors.New("database error"))

		categories, err := service.GetAllCategories()
		assert.Error(t, err)
		assert.Nil(t, categories)
	})
}

func TestStoreService_GetAllTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStoreRepositoryInterface(ctrl)
	service := &StoreService{storeRepo: mockRepo}

	expectedTags := []string{"japanese", "italian", "spicy"}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().GetAllTags().Return(expectedTags, nil)

		tags, err := service.GetAllTags()
		assert.NoError(t, err)
		assert.Equal(t, expectedTags, tags)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.EXPECT().GetAllTags().Return(nil, errors.New("database error"))

		tags, err := service.GetAllTags()
		assert.Error(t, err)
		assert.Nil(t, tags)
	})
}