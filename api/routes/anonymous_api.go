package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/app"
	"github.com/go-chi/chi/v5"
)

func loadUserAnonymousRoutes(r chi.Router) {
	// ##### Auth #####
	r.Post("/login", app.Handler(controllers.LoginUser))
	r.Post("/register", app.Handler(controllers.Register))
	r.Post("/recoveryPasword", app.Handler(controllers.SendRecoveryPasswordRequest))
	r.Post("/recoveryPasword/{key}", app.Handler(controllers.RecoveryPassword))
	r.Get("/logout", app.Handler(controllers.Logout))
	// ##### Auth #####

	loadAnonymousProductRoutes(r)
	loadAnonymousAppPicRoutes(r)
	loadAnonymousCommentRoutes(r)
	loadAnonymousBrandRoutes(r)
	loadAnonymousCategoryRoutes(r)
	loadAnonymousColorRoutes(r)
}

func loadAdminAnonymousRoutes(r chi.Router) {
	// ##### Auth #####
	r.Post("/login", app.Handler(controllers.LoginAdmin))
	r.Get("/logout", app.Handler(controllers.Logout))
	// ##### Auth #####
}
