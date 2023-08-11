package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/api/middlewares"
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/models"
	"github.com/go-chi/chi/v5"
)

func loadAdminCategoryRoutes(r chi.Router) {
	r.Get("/category",
		app.Handler(controllers.GetCategories,
			middlewares.Permitted(models.ACTION_CATEGORY_ADMIN_LIST)),
	)
	r.Get("/category/{id}", app.Handler(controllers.GetCategory,
		middlewares.Permitted(models.ACTION_CATEGORY_ADMIN_INFO)),
	)
	r.Post("/category", app.Handler(controllers.CreateCategory,
		middlewares.Permitted(models.ACTION_CATEGORY_ADMIN_CREATE)),
	)
	r.Post("/category/edit/{id}", app.Handler(controllers.EditCategory,
		middlewares.Permitted(models.ACTION_CATEGORY_ADMIN_UPDATE)),
	)
	r.Post("/category/delete/{id}", app.Handler(controllers.DeleteCategory,
		middlewares.Permitted(models.ACTION_CATEGORY_ADMIN_DELETE)),
	)
}

func loadAnonymousCategoryRoutes(r chi.Router) {
	r.Get("/category/selectList", app.Handler(controllers.GetCategoriesSelectList))
}
