package routes

import (
	"net/http"

	"github.com/esmailemami/eshop/apphttp"
	"github.com/esmailemami/eshop/apphttp/controllers"
	"github.com/esmailemami/eshop/apphttp/middlewares"
	"github.com/esmailemami/eshop/docs"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
)

func LoadApiRoutes(root *chi.Mux) {
	root.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	root.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Eshop API Server"))
	})

	docs.SwaggerInfo.Title = "Eshop API doc"
	docs.SwaggerInfo.Description = "Eshop API."
	docs.SwaggerInfo.Version = "1.0"
	port := viper.GetString("server.port")
	if port == "" {
		port = "6060"
	}
	docs.SwaggerInfo.Host = "127.0.0.1:" + port
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	root.Mount("/swagger/", httpSwagger.WrapHandler)

	root.Route("/api/v1", func(r chi.Router) {

		r.Post("/auth/login", apphttp.Handler(controllers.Login))
		r.Post("/auth/register", apphttp.Handler(controllers.Register))

		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthenticationHandler)

			// Auth
			r.Get("/auth/is_authenticated", apphttp.Handler(controllers.IsAuthenticated))
			r.Get("/auth/logout", apphttp.Handler(controllers.Logout))

			// category
			r.Get("/category", apphttp.Handler(controllers.GetCategories))
			r.Post("/category", apphttp.Handler(controllers.CreateCategory))
			r.Post("/category/edit/{id}", apphttp.Handler(controllers.EditCategory))
			r.Post("/category/delete/{id}", apphttp.Handler(controllers.DeleteCategory))
		})
	})
}
