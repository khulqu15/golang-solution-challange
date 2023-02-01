package main

import (
	"solution-challange/config"
	"solution-challange/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	config.ConnectDB()

	route.UserRoute(app)
	route.NebulizerRoute(app)

	app.Listen(":6000")
}
