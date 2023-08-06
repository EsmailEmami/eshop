package routes

import (
	"net/http"

	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/controllers"
	"github.com/esmailemami/eshop/app/middlewares"
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

		r.Post("/auth/login", app.Handler(controllers.Login))
		r.Post("/auth/register", app.Handler(controllers.Register))
		r.Post("/auth/recoveryPasword", app.Handler(controllers.SendRecoveryPasswordRequest))
		r.Post("/auth/recoveryPasword/{key}", app.Handler(controllers.RecoveryPassword))

		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthenticationHandler)

			// Auth
			r.Get("/auth/is_authenticated", app.Handler(controllers.IsAuthenticated))
			r.Get("/auth/logout", app.Handler(controllers.Logout))

			// category
			r.Get("/category", app.Handler(controllers.GetCategories))
			r.Get("/category/{id}", app.Handler(controllers.GetCategory))
			r.Post("/category", app.Handler(controllers.CreateCategory))
			r.Post("/category/edit/{id}", app.Handler(controllers.EditCategory))
			r.Post("/category/delete/{id}", app.Handler(controllers.DeleteCategory))

			// brand
			r.Get("/brand", app.Handler(controllers.GetBrands))
			r.Get("/brand/{id}", app.Handler(controllers.GetBrand))
			r.Post("/brand", app.Handler(controllers.CreateBrand))
			r.Post("/brand/edit/{id}", app.Handler(controllers.EditBrand))
			r.Post("/brand/delete/{id}", app.Handler(controllers.DeleteBrand))

			// color
			r.Get("/color", app.Handler(controllers.GetColors))
			r.Get("/color/{id}", app.Handler(controllers.GetColor))
			r.Post("/color", app.Handler(controllers.CreateColor))
			r.Post("/color/edit/{id}", app.Handler(controllers.EditColor))
			r.Post("/color/delete/{id}", app.Handler(controllers.DeleteColor))

			// file
			r.Post("/file/uploadImage/{itemId}/{fileType}", app.Handler(controllers.UploadImage))
			r.Post("/file/delete/{fileId}", app.Handler(controllers.DeleteFile))
			r.Get("/file/{fileId}", app.Handler(controllers.GetFile))
			r.Get("/file/stream/{fileId}", app.Handler(controllers.GetStreamingFile))
			r.Get("/file/{itemId}/{fileType}", app.Handler(controllers.GetItemFiles))
			r.Post("/file/changePriority/{fileId}/{itemId}/{priority}",app.Handler(controllers.FileChangePriority))

			// appPic
			r.Get("/appPic", app.Handler(controllers.GetAppPics))
			r.Get("/appPic/{id}", app.Handler(controllers.GetAppPic))
			r.Post("/appPic", app.Handler(controllers.CreateAppPic))
			r.Post("/appPic/edit/{id}", app.Handler(controllers.EditAppPic))
			r.Post("/appPic/delete/{id}", app.Handler(controllers.DeleteAppPic))

			// product
			r.Get("/product/{id}", app.Handler(controllers.GetProduct))
			r.Get("/product", app.Handler(controllers.GetProducts))
			r.Get("/product/list", app.Handler(controllers.GetProductsList))
			r.Post("/product", app.Handler(controllers.CreateProduct))
			r.Post("/product/edit/{id}", app.Handler(controllers.EditProduct))
			r.Post("/product/delete/{id}", app.Handler(controllers.DeleteProduct))
			r.Get("/product/suggestions", app.Handler(controllers.GetSuggestionProducts))

			// productItem
			r.Get("/productItem/{id}", app.Handler(controllers.GetProductItem))
			r.Post("/productItem", app.Handler(controllers.CreateProductItem))
			r.Post("/productItem/edit/{id}", app.Handler(controllers.EditProductItem))
			r.Post("/productItem/delete/{id}", app.Handler(controllers.DeleteProductItem))
			r.Get("/productItem/product/{productId}", app.Handler(controllers.GetProductItems))

			// productFeatureCategory
			r.Get("/productFeatureCategory", app.Handler(controllers.GetProductFeatureCategories))
			r.Get(
				"/productFeatureCategory/{id}",
				app.Handler(controllers.GetProductFeatureCategory),
			)
			r.Post("/productFeatureCategory", app.Handler(controllers.CreateProductFeatureCategory))
			r.Post(
				"/productFeatureCategory/edit/{id}",
				app.Handler(controllers.EditProductFeatureCategory),
			)
			r.Post(
				"/productFeatureCategory/delete/{id}",
				app.Handler(controllers.DeleteProductFeatureCategory),
			)

			// productFeatureKey
			r.Get("/productFeatureKey", app.Handler(controllers.GetProductFeatureKeys))
			r.Get("/productFeatureKey/{id}", app.Handler(controllers.GetProductFeatureKey))
			r.Post("/productFeatureKey", app.Handler(controllers.CreateProductFeatureKey))
			r.Post("/productFeatureKey/edit/{id}", app.Handler(controllers.EditProductFeatureKey))
			r.Post(
				"/productFeatureKey/delete/{id}",
				app.Handler(controllers.DeleteProductFeatureKey),
			)

			// productFeatureValue
			r.Get("/productFeatureValue", app.Handler(controllers.GetProductFeatureValues))
			r.Get("/productFeatureValue/{id}", app.Handler(controllers.GetProductFeatureValue))
			r.Post(
				"/productFeatureValue/{productId}",
				app.Handler(controllers.CreateProductFeatureValue),
			)
			r.Post(
				"/productFeatureValue/delete/{id}",
				app.Handler(controllers.DeleteProductFeatureValue),
			)

			// comment
			r.Get("/comment", app.Handler(controllers.GetComments))
			r.Get("/comment/{id}", app.Handler(controllers.GetComment))
			r.Post("/comment", app.Handler(controllers.CreateComment))
			r.Post("/comment/edit/{id}", app.Handler(controllers.EditComment))
			r.Post("/comment/delete/{id}", app.Handler(controllers.DeleteComment))

			// orderItem
			r.Post("/orderItem", app.Handler(controllers.CreateOrderItem))
			r.Post("/orderItem/delete/{productItemId}", app.Handler(controllers.DeleteOrderItem))

			// order
			r.Get("/order", app.Handler(controllers.GetOrder))
			r.Post("/order/checkout", app.Handler(controllers.CheckoutOrder))

			// user
			r.Get("/user", app.Handler(controllers.GetUser))
			r.Get("/user/orders", app.Handler(controllers.GetUserOrders))
			r.Get("/user/favoriteProducts", app.Handler(controllers.GetUserFavoriteProducts))

			// address
			r.Get("/address", app.Handler(controllers.GetAddresses))
			r.Get("/address/{id}", app.Handler(controllers.GetAddress))
			r.Post("/address", app.Handler(controllers.CreateAddress))
			r.Post("/address/edit/{id}", app.Handler(controllers.EditAddress))
			r.Post("/address/delete/{id}", app.Handler(controllers.DeleteAddress))

			// favoriteProductItem
			r.Post("/favoriteProductItem", app.Handler(controllers.CreateFavoriteProductItem))
			r.Post("/favoriteProductItem/delete/{productItemId}", app.Handler(controllers.DeleteFavoriteProductItem))
		})
	})
}
