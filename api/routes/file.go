package routes

import (
	"github.com/esmailemami/eshop/api/controllers"
	"github.com/esmailemami/eshop/api/middlewares"
	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/models"
	"github.com/go-chi/chi/v5"
)

func loadUserFileRoutes(r chi.Router) {
	r.Post("/file/uploadImage/{itemId}/{fileType}", app.Handler(controllers.UploadImage,
		middlewares.Permitted(models.ACTION_FILE_UPLOAD)),
	)
	r.Post("/file/delete/{fileId}", app.Handler(controllers.DeleteFile,
		middlewares.Permitted(models.ACTION_FILE_DELETE)),
	)
	r.Get("/file/{fileId}", app.Handler(controllers.GetFile,
		middlewares.Permitted(models.ACTION_FILE_DOWNLOAD)),
	)
	r.Get("/file/stream/{fileId}", app.Handler(controllers.GetStreamingFile,
		middlewares.Permitted(models.ACTION_FILE_DOWNLOAD)),
	)
	r.Get("/file/{itemId}/{fileType}", app.Handler(controllers.GetItemFiles,
		middlewares.Permitted(models.ACTION_FILE_LIST)),
	)
	r.Post("/file/changePriority/{fileId}/{itemId}/{priority}", app.Handler(controllers.FileChangePriority,
		middlewares.Permitted(models.ACTION_FILE_CHANGE_PRIORITY)),
	)
}
