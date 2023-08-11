package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/app"
	"github.com/go-chi/chi/v5"
)

func loadUserOrderItemRoutes(r chi.Router) {
	r.Post("/orderItem", app.Handler(controllers.CreateOrderItem))
	r.Post("/orderItem/delete/{productItemId}", app.Handler(controllers.DeleteOrderItem))
}
