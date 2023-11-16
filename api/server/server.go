package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/esmailemami/eshop/api/routes"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

func RunServer() {
	router := chi.NewRouter()

	routes.LoadApiRoutes(router)

	port := viper.GetString("server.port")
	if port == "" {
		port = "6060"
	}

	url := viper.GetString("server.url")

	filesDir := viper.GetString("global.file_storage_path")
	fileServer := http.FileServer(http.Dir(filesDir))
	router.Handle("/uploads/*", http.StripPrefix("/uploads/", fileServer))

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("server started at: " + url)
	defer srv.Close()
	log.Fatalln(srv.ListenAndServe())

}
