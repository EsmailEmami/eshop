package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/api/middlewares"
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/models"
	"github.com/go-chi/chi/v5"
)

func loadAdminProductFeatureValueRoutes(r chi.Router) {
	r.Get("/productFeatureValue",
		app.Handler(controllers.GetProductFeatureValues,
			middlewares.Permitted(models.ACTION_PRODUCT_FEATURE_VALUE_ADMIN_LIST)),
	)
	r.Get("/productFeatureValue/{id}", app.Handler(controllers.GetProductFeatureValue,
		middlewares.Permitted(models.ACTION_PRODUCT_FEATURE_VALUE_ADMIN_INFO)),
	)
	r.Post("/productFeatureValue/{productId}", app.Handler(controllers.CreateProductFeatureValue,
		middlewares.Permitted(models.ACTION_PRODUCT_FEATURE_VALUE_ADMIN_CREATE)),
	)
	r.Post("/productFeatureValue/delete/{id}", app.Handler(controllers.DeleteProductFeatureValue,
		middlewares.Permitted(models.ACTION_PRODUCT_FEATURE_VALUE_ADMIN_DELETE)),
	)
}
