package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/api/middlewares"
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/models"
	"github.com/go-chi/chi/v5"
)

func loadAdminDiscountRoutes(r chi.Router) {
	r.Get("/discount",
		app.Handler(controllers.GetDiscounts,
			middlewares.Permitted(models.ACTION_DISCOUNT_ADMIN_LIST)),
	)
	r.Get("/discount/{id}", app.Handler(controllers.GetDiscount,
		middlewares.Permitted(models.ACTION_DISCOUNT_ADMIN_INFO)),
	)
	r.Post("/discount", app.Handler(controllers.CreateDiscount,
		middlewares.Permitted(models.ACTION_DISCOUNT_ADMIN_CREATE)),
	)
	r.Post("/discount/edit/{id}", app.Handler(controllers.EditDiscount,
		middlewares.Permitted(models.ACTION_DISCOUNT_ADMIN_UPDATE)),
	)
	r.Post("/discount/delete/{id}", app.Handler(controllers.DeleteDiscount,
		middlewares.Permitted(models.ACTION_DISCOUNT_ADMIN_DELETE)),
	)
}

func loadUserDiscountRoutes(r chi.Router) {
	r.Get("/discount/validate/code/{code}", app.Handler(controllers.ValidateDiscountByCode))

}
