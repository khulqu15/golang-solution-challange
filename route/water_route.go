package route

import (
	"solution-challange/controller"

	"github.com/gofiber/fiber/v2"
)

func WaterRoute(app *fiber.App) {
	app.Get("/waters", controller.GetAllWaters)
	app.Post("/water", controller.CreateWater)
	app.Get("/water/:waterId", controller.GetAWater)
	app.Put("/water/:waterId", controller.EditWater)
	app.Delete("/water/:waterId", controller.DeleteWater)

	app.Post("/water/:waterId/data", controller.CreateWaterData)
	app.Put("/water/:waterId/data/:dataId", controller.EditWaterData)
	app.Delete("/water/:waterId/data/:dataId", controller.DeleteWaterData)
}
