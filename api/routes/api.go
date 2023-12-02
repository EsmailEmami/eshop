package routes

import (
	"net/http"

	"github.com/esmailemami/eshop/docs"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
)

func LoadApiRoutes(root *chi.Mux) {
	// root.Use(func(next http.Handler) http.Handler {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		if r.Method == http.MethodOptions {
	// 			return
	// 		}
	// 		next.ServeHTTP(w, r)
	// 	})
	// })

	root.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Eshop API Server"))
	})

	docs.SwaggerInfo.Title = "Eshop API doc"
	docs.SwaggerInfo.Description = "Eshop API."
	docs.SwaggerInfo.Version = "1.0"
	url := viper.GetString("global.swagger_url")

	docs.SwaggerInfo.Host = url
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"https", "http"}
	root.Mount("/swagger/", httpSwagger.WrapHandler)

	root.Route("/api/v1", func(rootRouter chi.Router) {

		rootRouter.Route("/user", func(r chi.Router) {
			loadUserAnonymousRoutes(r)
			loadUserRoutes(r)
		})

		rootRouter.Route("/admin", func(r chi.Router) {
			loadAdminAnonymousRoutes(r)
			loadAdminRoutes(r)
		})
	})
}
