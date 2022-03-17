package api

import (
	"compress/flate"
	"net/http"

	"github.com/diwise/api-environment/internal/pkg/application"
	"github.com/diwise/api-environment/internal/pkg/presentation/api/ngsi-ld/context"
	ngsi "github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
)

func createContextRegistry(app application.EnvironmentApp, log zerolog.Logger) ngsi.ContextRegistry {
	contextRegistry := ngsi.NewContextRegistry()
	ctxSource := context.CreateSource(app, log)
	contextRegistry.Register(ctxSource)
	return contextRegistry
}

func RegisterHandlers(r chi.Router, app application.EnvironmentApp, log zerolog.Logger) error {
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)

	// Enable gzip compression for ngsi-ld responses
	compressor := middleware.NewCompressor(flate.DefaultCompression, "application/json", "application/ld+json")
	r.Use(compressor.Handler)
	r.Use(middleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ctxReg := createContextRegistry(app, log)

	r.Post("/ngsi-ld/v1/entities", ngsi.NewCreateEntityHandler(ctxReg))
	r.Get("/ngsi-ld/v1/entities", ngsi.NewQueryEntitiesHandler(ctxReg))

	return nil
}
