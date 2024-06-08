package userRoutes

import (
	"github.com/gofiber/fiber/v2"
	userHandler "sistem-informasi-klinik/internal/handlers/users"
)

func SetupUsersRoutes(router fiber.Router) {
	user := router.Group("/users")
	// Create a user
	user.Post("/", userHandler.Register)
	user.Post("/", userHandler.Login)
	user.Post("/", userHandler.Logout)
	user.Get("/", userHandler.Get)
	// Read all users

	//user.Post("/optimize", belajarHandler.OptimizeReservations)

}
