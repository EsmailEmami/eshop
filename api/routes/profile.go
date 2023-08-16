package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/api/middlewares"
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/models"
	"github.com/go-chi/chi/v5"
)

func loadUserProfileRoutes(r chi.Router) {
	r.Get("/profile", app.Handler(controllers.GetUser))
	r.Get("/profile/orders", app.Handler(controllers.GetUserOrders))
	r.Get("/profile/favoriteProducts", app.Handler(controllers.GetUserFavoriteProducts))
}

func loadAdminProfileRoutes(r chi.Router) {
	r.Get("/profile/orders/{userId}", app.Handler(controllers.GetAdminUserOrders,
		middlewares.Permitted(models.ACTION_USER_ADMIN_ORDER_LIST),
	))
	r.Get("/profile/favoriteProducts/{userId}", app.Handler(controllers.GetAdminUserFavoriteProducts,
		middlewares.Permitted(models.ACTION_USER_ADMIN_FAVORITE_PRODUCT_LIST),
	))
}
