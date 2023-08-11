package routes

import (
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/controllers"
	"github.com/esmailemami/eshop/app/middlewares"
	"github.com/esmailemami/eshop/models"
	"github.com/go-chi/chi/v5"
)

func loadAdminBrandRoutes(r chi.Router) {
	r.Get("/brand",
		app.Handler(controllers.GetBrands,
			middlewares.Permitted(models.ACTION_BRAND_ADMIN_LIST)),
	)
	r.Get("/brand/{id}", app.Handler(controllers.GetBrand,
		middlewares.Permitted(models.ACTION_BRAND_ADMIN_INFO)),
	)
	r.Post("/brand", app.Handler(controllers.CreateBrand,
		middlewares.Permitted(models.ACTION_BRAND_ADMIN_CREATE)),
	)
	r.Post("/brand/edit/{id}", app.Handler(controllers.EditBrand,
		middlewares.Permitted(models.ACTION_BRAND_ADMIN_UPDATE)),
	)
	r.Post("/brand/delete/{id}", app.Handler(controllers.DeleteBrand,
		middlewares.Permitted(models.ACTION_BRAND_ADMIN_DELETE)),
	)
}
