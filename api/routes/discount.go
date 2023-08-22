package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/app"
	"github.com/go-chi/chi/v5"
)

func loadAdminDiscountRoutes(r chi.Router) {
	r.Get("/discount",
		app.Handler(controllers.GetDiscounts),
	)
	r.Get("/discount/{id}", app.Handler(controllers.GetDiscount))
	r.Post("/discount", app.Handler(controllers.CreateDiscount))
	r.Post("/discount/edit/{id}", app.Handler(controllers.EditDiscount))
	r.Post("/discount/delete/{id}", app.Handler(controllers.DeleteDiscount))
}
