package services

import (
	"errors"
	"sukimise/internal/models"
	"sukimise/internal/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repositories.UserRepositoryInterface
}

func NewUserService(userRepo repositories.UserRepositoryInterface) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(user *models.User) error {
	existingUser, _ := s.userRepo.GetByUsername(user.Username)
	if existingUser != nil {
		return errors.New("username already exists")
	}

	existingUser, _ = s.userRepo.GetByEmail(user.Email)
	if existingUser != nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.userRepo.Create(user)
}

func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return s.userRepo.GetByUsername(username)
}

func (s *UserService) ValidatePassword(user *models.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (s *UserService) UpdateUser(user *models.User) error {
	existingUser, err := s.userRepo.GetByID(user.ID)
	if err != nil {
		return err
	}

	if existingUser.Username != user.Username {
		existingUserByUsername, _ := s.userRepo.GetByUsername(user.Username)
		if existingUserByUsername != nil && existingUserByUsername.ID != user.ID {
			return errors.New("username already exists")
		}
	}

	if existingUser.Email != user.Email {
		existingUserByEmail, _ := s.userRepo.GetByEmail(user.Email)
		if existingUserByEmail != nil && existingUserByEmail.ID != user.ID {
			return errors.New("email already exists")
		}
	}

	return s.userRepo.Update(user)
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	return s.userRepo.Delete(id)
}