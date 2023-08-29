package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/api/middlewares"
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/models"
	"github.com/go-chi/chi/v5"
)

func loadAdminProductRoutes(r chi.Router) {
	r.Get("/product", app.Handler(controllers.GetAdminProducts,
		middlewares.Permitted(models.ACTION_PRODUCT_ADMIN_LIST),
	))
	r.Get("/product/{id}", app.Handler(controllers.GetAdminProduct,
		middlewares.Permitted(models.ACTION_PRODUCT_ADMIN_INFO),
	))
	r.Post("/product", app.Handler(controllers.CreateProduct,
		middlewares.Permitted(models.ACTION_PRODUCT_ADMIN_CREATE),
	))
	r.Post("/product/edit/{id}", app.Handler(controllers.EditProduct,
		middlewares.Permitted(models.ACTION_PRODUCT_ADMIN_UPDATE),
	))
	r.Post("/product/delete/{id}", app.Handler(controllers.DeleteProduct,
		middlewares.Permitted(models.ACTION_PRODUCT_ADMIN_DELETE),
	))
}

func loadAnonymousProductRoutes(r chi.Router) {
	r.Get("/product/suggestions", app.Handler(controllers.GetSuggestionProducts))
	r.Get("/product", app.Handler(controllers.GetUserProducts))
	r.Get("/product/{id}", app.Handler(controllers.GetProduct))
	r.Get("/product/selectList", app.Handler(controllers.GetProductsSelectList))
}
