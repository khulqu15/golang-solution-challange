package main

import (
	"solution-challange/config"
	"solution-challange/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	config.ConnectDB()

	route.UserRoute(app)
	route.NebulizerRoute(app)
	route.WaterRoute(app)

	app.Listen(":6000")
}
