package handlers

import (
	"sukimise/internal/services"
)

type Handler struct {
	userService   *services.UserService
	storeService  *services.StoreService
	reviewService *services.ReviewService
}

func NewHandler(userService *services.UserService, storeService *services.StoreService, reviewService *services.ReviewService) *Handler {
	return &Handler{
		userService:   userService,
		storeService:  storeService,
		reviewService: reviewService,
	}
}