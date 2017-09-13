package main

import (
	"github.com/goph/stdlib/net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// registerDebugRoutes allows to register custom routes in the debug server.
func registerDebugRoutes(a *application, h http.HandlerAcceptor) {
	h.Handle("/metrics", promhttp.Handler())
}
