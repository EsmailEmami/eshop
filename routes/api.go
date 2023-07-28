package routes

import (
	"net/http"

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
		_, _ = w.Write([]byte("Biling Core API Server"))
	})

	docs.SwaggerInfo.Title = "Biling Core API doc"
	docs.SwaggerInfo.Description = "Biling Core API."
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
		// Auth
		// r.Post("/login", apphttp.Handler(auth.Login))
		// r.Get("/logout", apphttp.Handler(auth.Logout))
		// r.Get("/is_authenticated", apphttp.Handler(auth.IsAuthenticated, middle
	})
}
