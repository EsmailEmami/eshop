package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/app"
	"github.com/go-chi/chi/v5"
)

func loadUserFavoriteProductItemRoutes(r chi.Router) {
	r.Post("/favoriteProductItem", app.Handler(controllers.CreateFavoriteProductItem))
	r.Post("/favoriteProductItem/delete/{productItemId}", app.Handler(controllers.DeleteFavoriteProductItem))
}
