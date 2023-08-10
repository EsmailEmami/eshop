package routes

import (
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/controllers"
	"github.com/go-chi/chi/v5"
)

func loadAnonymousRoutes(r chi.Router) {
	// ##### Auth #####
	r.Post("/auth/login", app.Handler(controllers.LoginUser))
	r.Post("/auth/register", app.Handler(controllers.Register))
	r.Post("/auth/recoveryPasword", app.Handler(controllers.SendRecoveryPasswordRequest))
	r.Post("/auth/recoveryPasword/{key}", app.Handler(controllers.RecoveryPassword))
	// ##### Auth #####

	loadAnonymousProductRoutes(r)
	loadAnonymousAppPicRoutes(r)
	loadAnonymousCommentRoutes(r)
}
