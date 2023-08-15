package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/api/middlewares"
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/models"
	"github.com/go-chi/chi/v5"
)

func loadAdminUserRoutes(r chi.Router) {
	r.Get("/user", app.Handler(controllers.GetUsers,
		middlewares.Permitted(models.ACTION_USER_ADMIN_LIST),
	))
	r.Get("/user/{id}", app.Handler(controllers.GetUser,
		middlewares.Permitted(models.ACTION_USER_ADMIN_INFO),
	))
	r.Post("/user", app.Handler(controllers.CreateUser,
		middlewares.Permitted(models.ACTION_USER_ADMIN_CREATE),
	))
	r.Post("/user/edit/{id}", app.Handler(controllers.EditUser,
		middlewares.Permitted(models.ACTION_USER_ADMIN_UPDATE),
	))
	r.Post("/user/delete/{id}", app.Handler(controllers.DeleteUser,
		middlewares.Permitted(models.ACTION_USER_ADMIN_DELETE),
	))
	r.Post("/user/recoveryPasword/{id}", app.Handler(controllers.AdminUserRecoveryPassword,
		middlewares.Permitted(models.ACTION_USER_ADMIN_RECOVERY_PASSWORD),
	))
}
