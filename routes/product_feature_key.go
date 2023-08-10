package routes

import (
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/controllers"
	"github.com/esmailemami/eshop/app/middlewares"
	"github.com/esmailemami/eshop/models"
	"github.com/go-chi/chi/v5"
)

func loadAdminProductFeatureKeyRoutes(r chi.Router) {
	r.Get("/productFeatureKey",
		app.Handler(controllers.GetProductFeatureKeys,
			middlewares.Permitted(models.ACTION_PRODUCT_FEATURE_KEY_ADMIN_LIST)),
	)
	r.Get("/productFeatureKey/{id}", app.Handler(controllers.GetProductFeatureKey,
		middlewares.Permitted(models.ACTION_PRODUCT_FEATURE_KEY_ADMIN_INFO)),
	)
	r.Post("/productFeatureKey", app.Handler(controllers.CreateProductFeatureKey,
		middlewares.Permitted(models.ACTION_PRODUCT_FEATURE_KEY_ADMIN_CREATE)),
	)
	r.Post("/productFeatureKey/edit/{id}", app.Handler(controllers.EditProductFeatureKey,
		middlewares.Permitted(models.ACTION_PRODUCT_FEATURE_KEY_ADMIN_UPDATE)),
	)
	r.Post("/productFeatureKey/delete/{id}", app.Handler(controllers.DeleteProductFeatureKey,
		middlewares.Permitted(models.ACTION_PRODUCT_FEATURE_KEY_ADMIN_DELETE)),
	)
}
