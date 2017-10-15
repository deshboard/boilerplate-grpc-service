package main

import (
	"flag"
	"time"
)

// defaultTimeout is used as a default for graceful shutdown timeout.
var defaultTimeout = 15 * time.Second

// Config holds any kind of configuration that comes from the outside world and is necessary for running.
type Config struct {
	// Meaningful values are recommended (eg. production, development, staging, release/123, etc)
	//
	// "development" is treated special: address types will use the loopback interface as default when none is defined.
	// This is useful when developing locally and listening on all interfaces requires elevated rights.
	Environment string `default:"production"`

	// Turns on some debug functionality: more verbose logs, exposed pprof, expvar and net trace in the debug server.
	Debug bool `split_words:"true"`

	// Defines the log format.
	// Valid values are: json, logfmt
	LogFormat string `split_words:"true" default:"json"`

	// Address of the debug server (configured by debug.addr flag)
	DebugAddr string `ignored:"true"`

	// Timeout for graceful shutdown (configured by shutdown.timeout flag)
	ShutdownTimeout time.Duration `ignored:"true"`

	// Address of the gRPC server (configured by grpc.addr flag)
	GrpcAddr string `ignored:"true"`

	// Enable the gRPC reflection service.
	GrpcEnableReflection bool `split_words:"true"`
}

// Flags configures a FlagSet.
//
// It still requires resolution (call to FlagSet.Parse) which is out of scope for this method.
func (c *Config) Flags(flags *flag.FlagSet) {
	flags.StringVar(&c.DebugAddr, "debug.addr", ":10000", "Debug and health check address")
	flags.DurationVar(&c.ShutdownTimeout, "shutdown.timeout", defaultTimeout, "Timeout for graceful shutdown")

	flags.StringVar(&c.GrpcAddr, "grpc.addr", ":8000", "gRPC service address")
}
