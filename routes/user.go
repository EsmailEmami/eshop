package routes

import (
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/controllers"
	"github.com/esmailemami/eshop/app/middlewares"
	"github.com/go-chi/chi/v5"
)

func loadUserRoutes(root chi.Router) {
	root.Group(func(r chi.Router) {
		r.Use(middlewares.AuthenticationHandler)

		// ##### Auth #####
		r.Get("/auth/is_authenticated", app.Handler(controllers.IsAuthenticated))
		r.Get("/auth/logout", app.Handler(controllers.Logout))
		// ##### Auth #####

		loadUserOrderItemRoutes(r)
		loadUserOrderRoutes(r)
		loadUserAddressRoutes(r)
		loadUserProfileRoutes(r)
		loadUserFavoriteProductItemRoutes(r)
		loadUserCommentRoutes(r)
		loadUserProductItemRoutes(r)
		loadUserFileRoutes(r)
	})
}
