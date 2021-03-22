package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/heyjoakim/devops-21/metrics"
	"github.com/heyjoakim/devops-21/routers"

	_ "github.com/heyjoakim/devops-21/docs" // docs is generated by Swag CLI, you have to import it.
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

func init() {
  // Log as JSON instead of the default ASCII formatter.
  log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	metrics.InitializeMetrics()

	router := mux.NewRouter()
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	router.PathPrefix("/static/").Handler(s)
	routers.AddUIRouter(router)

	apiRouter := router.PathPrefix("/api").Subrouter()
	routers.AddAPIRoutes(apiRouter)

	// Swagger
	apiRouter.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	PORT := ":8000"
	log.Info("Server is running on http://localhost:" + PORT)
	log.Fatal(http.ListenAndServe(PORT, router))
}
