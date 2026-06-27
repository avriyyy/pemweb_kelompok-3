package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"

	"toktik/database"
	"toktik/models"
	"toktik/routes"
)
var Store = session.New()

func main() {
	engine := html.New("./views", ".html")
	engine.Reload(true)

	app := fiber.New(fiber.Config{
		Views:        engine,
		AppName:      "TokTik",
		ServerHeader: "TokTik",
	})

	app.Static("/assets", "./assets")

	database.Connect()

	database.DB.AutoMigrate(
		&models.Film{},
		&models.Studio{},
		&models.Schedule{},
		&models.User{},
		&models.Transaksi{},
		&models.Tiket{},
	)
	routes.Web(app)

	log.Println("TokTik running on port 3334")
	if err := app.Listen("0.0.0.0:3334"); err != nil {
		log.Fatal(err)
	}
}
