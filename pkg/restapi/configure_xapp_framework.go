// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi/operations"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi/operations/common"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi/operations/xapp"
)

//go:generate swagger generate server --target ../../pkg --name XappFramework --spec ../../api/xapp_rest_api.yaml --exclude-main

func configureFlags(api *operations.XappFrameworkAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.XappFrameworkAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.XMLProducer = runtime.XMLProducer()

	if api.CommonSubscribeHandler == nil {
		api.CommonSubscribeHandler = common.SubscribeHandlerFunc(func(params common.SubscribeParams) middleware.Responder {
			return middleware.NotImplemented("operation common.Subscribe has not yet been implemented")
		})
	}
	if api.CommonUnsubscribeHandler == nil {
		api.CommonUnsubscribeHandler = common.UnsubscribeHandlerFunc(func(params common.UnsubscribeParams) middleware.Responder {
			return middleware.NotImplemented("operation common.Unsubscribe has not yet been implemented")
		})
	}
	if api.CommonGetAllSubscriptionsHandler == nil {
		api.CommonGetAllSubscriptionsHandler = common.GetAllSubscriptionsHandlerFunc(func(params common.GetAllSubscriptionsParams) middleware.Responder {
			return middleware.NotImplemented("operation common.GetAllSubscriptions has not yet been implemented")
		})
	}
	if api.XappGetXappConfigListHandler == nil {
		api.XappGetXappConfigListHandler = xapp.GetXappConfigListHandlerFunc(func(params xapp.GetXappConfigListParams) middleware.Responder {
			return middleware.NotImplemented("operation xapp.GetXappConfigList has not yet been implemented")
		})
	}

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
	return handler
}
