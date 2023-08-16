package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/app"
	"github.com/go-chi/chi/v5"
)

func loadUserOrderRoutes(r chi.Router) {
	r.Get("/order", app.Handler(controllers.GetOrder))
	r.Post("/order/checkout", app.Handler(controllers.CheckoutOrder))
}

func loadAdminOrderRoutes(r chi.Router) {
	r.Get("/order", app.Handler(controllers.GetAdminOrders))
}
