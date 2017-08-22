package main

import (
	"io"
	"net/http"

	"github.com/uber-go/tally"
	promreporter "github.com/uber-go/tally/prometheus"
)

// newMetrics returns a new tally.Scope used as a root scope.
func newMetrics(config *configuration) interface {
	tally.Scope
	io.Closer
} {
	reporter := promreporter.NewReporter(promreporter.Options{})

	options := tally.ScopeOptions{
		CachedReporter: reporter,
	}

	scope, closer := tally.NewRootScope(options, MetricsReportInterval)

	return struct {
		tally.Scope
		io.Closer
		http.Handler
	}{
		Scope:   scope,
		Closer:  closer,
		Handler: reporter.HTTPHandler(),
	}
}
