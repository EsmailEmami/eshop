package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/api/middlewares"
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/models"
	"github.com/go-chi/chi/v5"
)

func loadUserAddressRoutes(r chi.Router) {
	r.Get("/address", app.Handler(controllers.GetAddresses))
	r.Get("/address/{id}", app.Handler(controllers.GetAddress))
	r.Post("/address", app.Handler(controllers.CreateAddress))
	r.Post("/address/edit/{id}", app.Handler(controllers.EditAddress))
	r.Post("/address/delete/{id}", app.Handler(controllers.DeleteAddress))
}

func loadAdminAddressRoutes(r chi.Router) {
	r.Get("/address/{userId}", app.Handler(controllers.GetAdminUserAddresses,
		middlewares.Permitted(models.ACTION_ADDRESS_ADMIN_LIST)),
	)
}
