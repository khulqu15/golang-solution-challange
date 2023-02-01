package route

import (
	"solution-challange/controller"

	"github.com/gofiber/fiber/v2"
)

func NebulizerRoute(app *fiber.App) {
	app.Get("/nebulizers", controller.GetAllNebulizers)
	app.Post("/nebulizer", controller.CreateNebulizer)
	app.Get("/nebulizer/:nebulizerId", controller.GetANebulizer)
	app.Put("/nebulizer/:nebulizerId", controller.EditNebulizer)
	app.Delete("/nebulizer/:nebulizerId", controller.DeleteNebulizer)

	app.Post("/nebulizer/:nebulizerId/data", controller.CreateNebulizerData)
	app.Put("/nebulizer/:nebulizerId/data/:dataId", controller.EditNebulizerData)
	app.Delete("/nebulizer/:nebulizerId/data/:dataId", controller.DeleteNebulizerData)
}
