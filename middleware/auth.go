package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// RequireLogin
func RequireLogin(c *fiber.Ctx) error {
	userID := c.Cookies("user_id")
	if userID == "" {
		// Belum login, redirect ke halaman login
		return c.Redirect("/login")
	}
	return c.Next()
}

// RequireAdmin
func RequireAdmin(c *fiber.Ctx) error {
	userRole := c.Cookies("user_role")
	if userRole != "admin" {
		// Bukan admin, tampilkan halaman 403 Forbidden
		return c.Status(403).SendString("Akses ditolak. Halaman ini hanya untuk admin.")
	}
	return c.Next()
}