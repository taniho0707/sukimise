package services

import (
	"errors"
	"sukimise/internal/models"
	mocks "sukimise/internal/repositories/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "editor",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().GetByUsername(user.Username).Return(nil, errors.New("not found"))
		mockRepo.EXPECT().GetByEmail(user.Email).Return(nil, errors.New("not found"))
		mockRepo.EXPECT().Create(gomock.Any()).Return(nil)

		err := service.CreateUser(user)
		assert.NoError(t, err)
		// Password should be hashed
		assert.NotEqual(t, "password123", user.Password)
	})

	t.Run("username already exists", func(t *testing.T) {
		existingUser := &models.User{Username: user.Username}
		mockRepo.EXPECT().GetByUsername(user.Username).Return(existingUser, nil)

		err := service.CreateUser(user)
		assert.Error(t, err)
		assert.Equal(t, "username already exists", err.Error())
	})

	t.Run("email already exists", func(t *testing.T) {
		existingUser := &models.User{Email: user.Email}
		mockRepo.EXPECT().GetByUsername(user.Username).Return(nil, errors.New("not found"))
		mockRepo.EXPECT().GetByEmail(user.Email).Return(existingUser, nil)

		err := service.CreateUser(user)
		assert.Error(t, err)
		assert.Equal(t, "email already exists", err.Error())
	})
}

func TestUserService_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	expectedUser := &models.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(userID).Return(expectedUser, nil)

		user, err := service.GetUserByID(userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(userID).Return(nil, errors.New("user not found"))

		user, err := service.GetUserByID(userID)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserService_GetUserByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	username := "testuser"
	expectedUser := &models.User{
		ID:       uuid.New(),
		Username: username,
		Email:    "test@example.com",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().GetByUsername(username).Return(expectedUser, nil)

		user, err := service.GetUserByUsername(username)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.EXPECT().GetByUsername(username).Return(nil, errors.New("user not found"))

		user, err := service.GetUserByUsername(username)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserService_ValidatePassword(t *testing.T) {
	service := NewUserService(nil) // No mock needed for this test

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &models.User{
		Password: string(hashedPassword),
	}

	t.Run("correct password", func(t *testing.T) {
		isValid := service.ValidatePassword(user, password)
		assert.True(t, isValid)
	})

	t.Run("incorrect password", func(t *testing.T) {
		isValid := service.ValidatePassword(user, "wrongpassword")
		assert.False(t, isValid)
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	existingUser := &models.User{
		ID:       userID,
		Username: "olduser",
		Email:    "old@example.com",
	}

	updatedUser := &models.User{
		ID:       userID,
		Username: "newuser",
		Email:    "new@example.com",
	}

	t.Run("success - username and email changed", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(userID).Return(existingUser, nil)
		mockRepo.EXPECT().GetByUsername("newuser").Return(nil, errors.New("not found"))
		mockRepo.EXPECT().GetByEmail("new@example.com").Return(nil, errors.New("not found"))
		mockRepo.EXPECT().Update(updatedUser).Return(nil)

		err := service.UpdateUser(updatedUser)
		assert.NoError(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(userID).Return(nil, errors.New("user not found"))

		err := service.UpdateUser(updatedUser)
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("username already exists", func(t *testing.T) {
		anotherUserID := uuid.New()
		anotherUser := &models.User{ID: anotherUserID, Username: "newuser"}
		
		mockRepo.EXPECT().GetByID(userID).Return(existingUser, nil)
		mockRepo.EXPECT().GetByUsername("newuser").Return(anotherUser, nil)

		err := service.UpdateUser(updatedUser)
		assert.Error(t, err)
		assert.Equal(t, "username already exists", err.Error())
	})

	t.Run("email already exists", func(t *testing.T) {
		anotherUserID := uuid.New()
		anotherUser := &models.User{ID: anotherUserID, Email: "new@example.com"}
		
		mockRepo.EXPECT().GetByID(userID).Return(existingUser, nil)
		mockRepo.EXPECT().GetByUsername("newuser").Return(nil, errors.New("not found"))
		mockRepo.EXPECT().GetByEmail("new@example.com").Return(anotherUser, nil)

		err := service.UpdateUser(updatedUser)
		assert.Error(t, err)
		assert.Equal(t, "email already exists", err.Error())
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().Delete(userID).Return(nil)

		err := service.DeleteUser(userID)
		assert.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.EXPECT().Delete(userID).Return(errors.New("delete failed"))

		err := service.DeleteUser(userID)
		assert.Error(t, err)
		assert.Equal(t, "delete failed", err.Error())
	})
}