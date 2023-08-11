package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/app"
	"github.com/go-chi/chi/v5"
)

func loadUserFileRoutes(r chi.Router) {
	r.Post("/file/uploadImage/{itemId}/{fileType}", app.Handler(controllers.UploadImage))
	r.Post("/file/delete/{fileId}", app.Handler(controllers.DeleteFile))
	r.Get("/file/{fileId}", app.Handler(controllers.GetFile))
	r.Get("/file/stream/{fileId}", app.Handler(controllers.GetStreamingFile))
	r.Get("/file/{itemId}/{fileType}", app.Handler(controllers.GetItemFiles))
	r.Post("/file/changePriority/{fileId}/{itemId}/{priority}", app.Handler(controllers.FileChangePriority))
}
