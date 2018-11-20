// This file is safe to edit. Once it exists it will not be overwritten

package api

import (
	"crypto/tls"

	"net/http"
	"strings"

	"github.com/dre1080/recover"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/analyze/asset"
	"github.com/supergiant/analyze/pkg/api/operations"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/rs/cors"

)

//go:generate swagger generate server --target ../pkg --name Analyze --spec ../swagger/api-spec.yml --server-package api --exclude-main

func configureFlags(api *operations.AnalyzeAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.AnalyzeAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.Logger = logrus.Infof

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handlerWithRecovery := recover.New(&recover.Options{
		Log: logrus.Error,
	})(handler)

	//TODO fix CORS till release
	corsHandler := cors.New(cors.Options{
		Debug:          false,
		AllowedHeaders: []string{"*"},
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{},
		MaxAge:         1000,
	}).Handler(handlerWithRecovery)

	handlerWithSwagger := swaggerMiddleware(corsHandler)

	return handlerWithSwagger
}

func swaggerMiddleware(handler http.Handler) http.Handler {
	var staticServer = http.FileServer(asset.SwaggerAssets)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Shortcut helpers for swagger-ui
		if r.URL.Path == "/api/v1/swagger-ui" || r.URL.Path == "/api/v1/help" {
			http.Redirect(w, r, "/api/v1/swagger-ui/", http.StatusFound)
			return
		}
		// Serving ./swagger-ui/
		if strings.Index(r.URL.Path, "/api/v1/swagger-ui/") == 0 {
			http.StripPrefix("/api/v1/swagger-ui/", staticServer).ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
