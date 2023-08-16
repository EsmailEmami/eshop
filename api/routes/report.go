package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/api/middlewares"
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/models"
	"github.com/go-chi/chi/v5"
)

func loadAdminReportRoutes(r chi.Router) {
	r.Get("/report/revenueByCategory", app.Handler(controllers.ReportRevenueByCategory,
		middlewares.Permitted(models.ACTION_REPORT_ADMIN_REVENUE_BY_CATEGORY),
	))

	r.Get("/report/sellsChart", app.Handler(controllers.ReportSellsChart))
}
