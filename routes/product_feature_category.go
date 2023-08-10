package routes

import (
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/controllers"
	"github.com/esmailemami/eshop/app/middlewares"
	"github.com/esmailemami/eshop/models"
	"github.com/go-chi/chi/v5"
)

func loadAdminProductFeatureCategoryRoutes(r chi.Router) {
	r.Get("/productFeatureCategory",
		app.Handler(controllers.GetProductFeatureCategories,
			middlewares.Permitted(models.ACTION_PRODUCT_FEATURE_CATEGORY_ADMIN_LIST)),
	)
	r.Get("/productFeatureCategory/{id}", app.Handler(controllers.GetProductFeatureCategory,
		middlewares.Permitted(models.ACTION_PRODUCT_FEATURE_CATEGORY_ADMIN_INFO)),
	)
	r.Post("/productFeatureCategory", app.Handler(controllers.CreateProductFeatureCategory,
		middlewares.Permitted(models.ACTION_PRODUCT_FEATURE_CATEGORY_ADMIN_CREATE)),
	)
	r.Post("/productFeatureCategory/edit/{id}", app.Handler(controllers.EditProductFeatureCategory,
		middlewares.Permitted(models.ACTION_PRODUCT_FEATURE_CATEGORY_ADMIN_UPDATE)),
	)
	r.Post("/productFeatureCategory/delete/{id}", app.Handler(controllers.DeleteProductFeatureCategory,
		middlewares.Permitted(models.ACTION_PRODUCT_FEATURE_CATEGORY_ADMIN_DELETE)),
	)
}
