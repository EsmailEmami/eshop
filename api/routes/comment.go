package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/app"
	"github.com/go-chi/chi/v5"
)

func loadUserCommentRoutes(r chi.Router) {
	r.Get("/comment", app.Handler(controllers.GetUserComments))
	r.Get("/comment/{id}", app.Handler(controllers.GetComment))
	r.Post("/comment", app.Handler(controllers.CreateComment))
	r.Post("/comment/edit/{id}", app.Handler(controllers.EditComment))
	r.Post("/comment/delete/{id}", app.Handler(controllers.DeleteComment))
}

func loadAdminCommentRoutes(r chi.Router) {
	r.Get("/comment", app.Handler(controllers.GetAdminUserComments))
}

func loadAnonymousCommentRoutes(r chi.Router) {
	r.Get("/comment/product/{productId}", app.Handler(controllers.GetProductComments))
}
