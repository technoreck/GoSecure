package routes

import (
	"GoSecure/server/handlers"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("views/index.html")
	})

	app.Post("/dnsinfo", handlers.Handler)
	app.Get("/getData", handlers.Headerhandler)
	app.Post("/scan", handlers.ScanHandler)
	app.Post("/hsts", handlers.HstsHandler)
	app.Post("/servs", handlers.Servstatushandler)
	app.Post("/dnssec", handlers.Dnssechandler)
	app.Post("/screenshot", handlers.Screenshothandler)
	app.Post("/resolve", handlers.Dnsserverhandler)

}
func Serve() {
	app := fiber.New()
	app.Static("/static", "./static")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("views/index.html")
	})

	SetupRoutes(app)
	log.Fatal(app.Listen(":8080"))
}
