package controllers

import "github.com/gofiber/fiber/v2"

type AuthController struct{}

func (AuthController) LoginPage(c *fiber.Ctx) error {
	return c.Render("auth/login", fiber.Map{
		"Title": "Masuk",
		"Error": c.Query("error"),
	}, "layouts/base")
}

func (AuthController) LoginSubmit(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return c.Redirect("/login?error=Email+dan+kata+sandi+wajib+diisi")
	}
	if len(password) < 6 {
		return c.Redirect("/login?error=Kata+sandi+minimal+6+karakter")
	}

	return c.Redirect("/")
}

func (AuthController) RegisterPage(c *fiber.Ctx) error {
	return c.Render("auth/register", fiber.Map{
		"Title": "Daftar",
		"Error": c.Query("error"),
	}, "layouts/base")
}

func (AuthController) RegisterSubmit(c *fiber.Ctx) error {
	nama := c.FormValue("nama")
	email := c.FormValue("email")
	password := c.FormValue("password")

	if nama == "" || email == "" || password == "" {
		return c.Redirect("/register?error=Nama,+email,+dan+kata+sandi+wajib+diisi")
	}
	if len(password) < 6 {
		return c.Redirect("/register?error=Kata+sandi+minimal+6+karakter")
	}

	return c.Redirect("/login?registered=1")
}

func (AuthController) Logout(c *fiber.Ctx) error {
	return c.Redirect("/login")
}
