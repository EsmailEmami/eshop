package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/api/middlewares"
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/models"
	"github.com/go-chi/chi/v5"
)

func loadAdminRoleRoutes(r chi.Router) {
	r.Get("/role",
		app.Handler(controllers.GetRoles,
			middlewares.Permitted(models.ACTION_ROLE_ADMIN_LIST)),
	)
	r.Get("/role/{id}", app.Handler(controllers.GetRole,
		middlewares.Permitted(models.ACTION_ROLE_ADMIN_INFO)),
	)
	r.Post("/role", app.Handler(controllers.CreateRole,
		middlewares.Permitted(models.ACTION_ROLE_ADMIN_CREATE)),
	)
	r.Post("/role/edit/{id}", app.Handler(controllers.EditRole,
		middlewares.Permitted(models.ACTION_ROLE_ADMIN_UPDATE)),
	)
	r.Post("/role/delete/{id}", app.Handler(controllers.DeleteRole,
		middlewares.Permitted(models.ACTION_ROLE_ADMIN_DELETE)),
	)
	r.Get("/role/permissions",
		app.Handler(controllers.GetPermissions,
			middlewares.Permitted(models.ACTION_ROLE_ADMIN_PERMISSIONS)),
	)
}
