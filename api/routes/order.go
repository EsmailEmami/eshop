package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/api/middlewares"
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/models"
	"github.com/go-chi/chi/v5"
)

func loadUserOrderRoutes(r chi.Router) {
	r.Get("/order", app.Handler(controllers.GetOrder))
	r.Post("/order/checkout/{addressId}", app.Handler(controllers.CheckoutOrder))
}

func loadAdminOrderRoutes(r chi.Router) {
	r.Get("/order", app.Handler(controllers.GetAdminOrders,
		middlewares.Permitted(models.ACTION_ORDER_ADMIN_LIST),
	))
}
