package controllers

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"toktik/database"
	"toktik/models"
)
type AuthController struct{}

func (AuthController) LoginPage(c *fiber.Ctx) error {
	if c.Cookies("user_id") != "" {
		return c.Redirect("/")
	}
	return c.Render("auth/login", fiber.Map{
		"Title": "Masuk",
		"Error": c.Query("error"),
	}, "layouts/base")
}

func (AuthController) LoginSubmit(c *fiber.Ctx) error {

	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return c.Redirect("/login?error=Email+dan+password+wajib+diisi")
	}

	var user models.User

	err := database.DB.Where("email = ?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Redirect("/login?error=Email+tidak+ditemukan")
	}

	if err != nil {
		return c.Redirect("/login?error=Terjadi+kesalahan")
	}

	err = bcrypt.CompareHashAndPassword(
	[]byte(user.Password),
	[]byte(password),
	)

	if err != nil {
		return c.Redirect("/login?error=Password+salah")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "user_id",
		Value:    fmt.Sprintf("%d", user.ID),
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "role",
		Value:    user.Role,
		HTTPOnly: true,
	})

	// Redirect sesuai role
	if user.Role == "admin" {
		return c.Redirect("/admin")
	}

	return c.Redirect("/")
}

func (AuthController) RegisterPage(c *fiber.Ctx) error {
	if c.Cookies("user_id") != "" {
		return c.Redirect("/")
	}
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
		return c.Redirect("/register?error=Semua+field+wajib+diisi")
	}

	if len(password) < 6 {
		return c.Redirect("/register?error=Password+minimal+6+karakter")
	}

	var user models.User

	err := database.DB.Where("email = ?", email).First(&user).Error

	if err == nil {
		return c.Redirect("/register?error=Email+sudah+terdaftar")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Redirect("/register?error=Terjadi+kesalahan")
	}

	hashPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return c.Redirect("/register?error=Gagal+hash+password")
	}

	user = models.User{
		Nama:     nama,
		Email:    email,
		Password: string(hashPassword),
		Role:     "user",
	}

	if err := database.DB.Create(&user).Error; err != nil {
    log.Println("ERROR SIMPAN USER:", err)
    return c.Redirect("/register?error=Gagal+menyimpan+user")
}

	return c.Redirect("/login")
}

func (AuthController) Logout(c *fiber.Ctx) error {

	c.ClearCookie("user_id")
	c.ClearCookie("role")

	return c.Redirect("/login")
}
