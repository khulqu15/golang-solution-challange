package route

import (
	"solution-challange/controller"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	app.Get("/users", controller.GetAllUser)
	app.Post("/user", controller.CreateUser)
	app.Get("/user/:userId", controller.GetAUser)
	app.Put("/user/:userId", controller.EditUser)
	app.Delete("/user/:userId", controller.DeleteUser)
}
