package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/app"
	"github.com/go-chi/chi/v5"
)

func loadAdminSettingsRoutes(r chi.Router) {
	r.Get("/settings/{item}",
		app.Handler(controllers.GetSettings),
	)

	r.Post("/settings/{item}",
		app.Handler(controllers.UpdateSetting),
	)
}
