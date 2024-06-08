package belajarRoutes

import (
	belajarHandler "sistem-informasi-klinik/internal/handlers/belajar"

	"github.com/gofiber/fiber/v2"
)

func SetupBelajarRoutes(router fiber.Router) {
	user := router.Group("/belajar")
	// Create a user
	user.Post("/", belajarHandler.CreateUser)
	// Read all users
	user.Get("/", belajarHandler.GetUsers)
	// Read one user by query parameter
	user.Get("/user", belajarHandler.GetUser) // Menggunakan /user untuk menghindari konflik dengan rute lain
	// Update one user
	user.Put("/", belajarHandler.UpdateUser)
	// Delete one user by URL parameter
	user.Delete("/:id_user", belajarHandler.DeleteUser)
	// Delete one user by body parameter
	user.Delete("/user", belajarHandler.DeleteUser) // Menggunakan DELETE dengan body
	user.Get("/:id_user", belajarHandler.GetJadwalById)

	//user.Post("/optimize", belajarHandler.OptimizeReservations)
}


