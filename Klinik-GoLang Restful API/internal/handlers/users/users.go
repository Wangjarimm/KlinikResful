package userHandler

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"sistem-informasi-klinik/database"
	"sistem-informasi-klinik/internal/model"
)

// ...

func Register(c *fiber.Ctx) error {
	var user model.Users
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	// Hash the password before saving to the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}
	user.Password = string(hashedPassword)

	database.DB.Create(&user)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User Berhasil Dibuat!",
		"data":    user,
	})
}

func Login(c *fiber.Ctx) error {
	var user model.Users
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	var users model.Users
	database.DB.Where("username = ?", users.Username).First(&users)

	if users.Id == 0 || bcrypt.CompareHashAndPassword([]byte(users.Username), []byte(users.Password)) != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	// You can generate and return a JWT token here if needed

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
	})
}

func Logout(c *fiber.Ctx) error {
	// You may need to implement session management logic here if using sessions
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout successful",
	})
}

func Get(c *fiber.Ctx) error {
	{
		var user []model.Users
		// Find all users in database
		result := database.DB.Find(&user)
		// Check for errors during query execution
		if result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": result.Error.Error(),
			})
		}
		// Return list of users
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User Berhasil Ditampilkan!",
			"data":    user,
		})
	}
}
