package routes

import (
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/controllers"
	"github.com/go-chi/chi/v5"
)

func loadUserOrderRoutes(r chi.Router) {
	r.Get("/order", app.Handler(controllers.GetOrder))
	r.Post("/order/checkout", app.Handler(controllers.CheckoutOrder))
}
