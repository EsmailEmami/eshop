package routes

import (
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/controllers"
	"github.com/go-chi/chi/v5"
)

func loadAdminProductItemRoutes(r chi.Router) {
	r.Post("/productItem", app.Handler(controllers.CreateProductItem))
	r.Post("/productItem/edit/{id}", app.Handler(controllers.EditProductItem))
	r.Post("/productItem/delete/{id}", app.Handler(controllers.DeleteProductItem))
	r.Get("/productItem/product/{productId}", app.Handler(controllers.GetProductItems))
}

func loadUserProductItemRoutes(r chi.Router) {
	r.Get("/productItem/{id}", app.Handler(controllers.GetProductItem))
}
