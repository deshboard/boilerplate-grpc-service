package main

import promreporter "github.com/uber-go/tally/prometheus"

// newMetricsReporter returns one of tally.StatsReporter and tally.CachedStatsReporter.
func newMetricsReporter(config *configuration) interface{} {
	return promreporter.NewReporter(promreporter.Options{})
}
