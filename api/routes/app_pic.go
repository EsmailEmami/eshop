package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/api/middlewares"
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/models"
	"github.com/go-chi/chi/v5"
)

func loadAdminAppPicRoutes(r chi.Router) {
	r.Get("/appPic/{id}", app.Handler(controllers.GetAppPic,
		middlewares.Permitted(models.ACTION_APP_PIC_ADMIN_INFO)),
	)
	r.Post("/appPic", app.Handler(controllers.CreateAppPic,
		middlewares.Permitted(models.ACTION_APP_PIC_ADMIN_CREATE)),
	)
	r.Post("/appPic/edit/{id}", app.Handler(controllers.EditAppPic,
		middlewares.Permitted(models.ACTION_APP_PIC_ADMIN_UPDATE)),
	)
	r.Post("/appPic/delete/{id}", app.Handler(controllers.DeleteAppPic,
		middlewares.Permitted(models.ACTION_APP_PIC_ADMIN_DELETE)),
	)
}

func loadAnonymousAppPicRoutes(r chi.Router) {
	r.Get("/appPic", app.Handler(controllers.GetAppPics))
}
