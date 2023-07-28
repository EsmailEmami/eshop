package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/esmailemami/eshop/routes"
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

	ip := viper.GetString("server.ip")
	if ip == "" {
		ip = "127.0.0.1"
	}

	srv := &http.Server{
		Handler:      router,
		Addr:         ip + ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("server started at: http://127.0.0.1:" + port)
	defer srv.Close()
	log.Fatalln(srv.ListenAndServe())

}
