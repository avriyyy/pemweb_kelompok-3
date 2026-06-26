package middleware

import (
	"fmt"
	"strconv"

	"toktik/database"
	"toktik/models"

	"github.com/gofiber/fiber/v2"
)

func AuthOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Cookies("user_id")
		fmt.Println("AUTH user_id =", userID)

		if userID == "" {
			fmt.Println("Redirect ke login (AuthOnly)")
			return c.Redirect("/login")
		}

		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Cookies("role")
		fmt.Println("ROLE =", role)

		if role != "admin" {
			fmt.Println("Redirect ke login (AdminOnly)")
			return c.Redirect("/login")
		}

		fmt.Println("Admin berhasil masuk")
		return c.Next()
	}
}

func CurrentUser(c *fiber.Ctx) (*models.User, error) {

	userID := c.Cookies("user_id")

	if userID == "" {
		return nil, fiber.ErrUnauthorized
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}

	var user models.User

	if err := database.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}