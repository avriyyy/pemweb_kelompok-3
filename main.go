package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"toktik/database"
	"toktik/routes"
)

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

	routes.Web(app)

	log.Println("TokTik running on http://localhost:3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
