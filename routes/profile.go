package routes

import (
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/controllers"
	"github.com/go-chi/chi/v5"
)

func loadUserProfileRoutes(r chi.Router) {
	r.Get("/profile", app.Handler(controllers.GetUser))
	r.Get("/profile/orders", app.Handler(controllers.GetUserOrders))
	r.Get("/profile/favoriteProducts", app.Handler(controllers.GetUserFavoriteProducts))
}

func loadAdminProfileRoutes(r chi.Router) {
	r.Get("/profile/orders/{userId}", app.Handler(controllers.GetAdminUserOrders))
	r.Get("/profile/favoriteProducts/{userId}", app.Handler(controllers.GetAdminUserFavoriteProducts))
}
