package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/api/middlewares"
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/models"
	"github.com/go-chi/chi/v5"
)

func loadAdminColorRoutes(r chi.Router) {
	r.Get("/color",
		app.Handler(controllers.GetColors,
			middlewares.Permitted(models.ACTION_COLOR_ADMIN_LIST)),
	)
	r.Get("/color/{id}", app.Handler(controllers.GetColor,
		middlewares.Permitted(models.ACTION_COLOR_ADMIN_INFO)),
	)
	r.Post("/color", app.Handler(controllers.CreateColor,
		middlewares.Permitted(models.ACTION_COLOR_ADMIN_CREATE)),
	)
	r.Post("/color/edit/{id}", app.Handler(controllers.EditColor,
		middlewares.Permitted(models.ACTION_COLOR_ADMIN_UPDATE)),
	)
	r.Post("/color/delete/{id}", app.Handler(controllers.DeleteColor,
		middlewares.Permitted(models.ACTION_COLOR_ADMIN_DELETE)),
	)
}
