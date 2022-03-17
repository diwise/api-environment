package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/diwise/api-environment/internal/pkg/application"
	"github.com/diwise/api-environment/internal/pkg/infrastructure/repositories/database"
	"github.com/diwise/api-environment/internal/pkg/presentation/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog/log"
)

func main() {
	serviceName := "api-environment"

	logger := log.With().Str("service", strings.ToLower(serviceName)).Logger()
	logger.Info().Msg("starting up ...")

	db, err := database.NewDatabaseConnection(database.NewPostgreSQLConnector(logger))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database, shutting down... ")
	}

	app := application.NewEnvironmentApp(db, logger)

	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(
		httplog.NewLogger(serviceName, httplog.Options{
			JSON: true,
		}),
	))
	api.RegisterHandlers(r, app, logger)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8080"
	}

	log.Info().Str("port", port).Msg("starting to listen for connections")

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen for connections")
	}
}
