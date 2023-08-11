package routes

import (
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/controllers"
	"github.com/esmailemami/eshop/app/middlewares"
	"github.com/go-chi/chi/v5"
)

func loadAdminRoutes(root chi.Router) {
	root.Group(func(r chi.Router) {
		r.Use(middlewares.AuthenticationHandler)
		r.Use(middlewares.CanInvokeRouteUnlessUser)

		// ##### Auth #####
		r.Get("/is_authenticated", app.Handler(controllers.IsAuthenticated))
		r.Get("/logout", app.Handler(controllers.Logout))
		// ##### Auth #####

		loadAdminProductFeatureCategoryRoutes(r)
		loadAdminProductFeatureCategoryRoutes(r)
		loadAdminProductFeatureValueRoutes(r)
		loadAdminAppPicRoutes(r)
		loadAdminProductRoutes(r)
		loadAdminAddressRoutes(r)
		loadAdminColorRoutes(r)
		loadAdminCategoryRoutes(r)
		loadAdminBrandRoutes(r)
		loadAdminProfileRoutes(r)
		loadAdminCommentRoutes(r)
		loadAdminProductItemRoutes(r)
		loadAdminRoleRoutes(r)
	})
}
