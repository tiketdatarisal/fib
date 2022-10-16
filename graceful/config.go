package graceful

import "github.com/go-logr/logr"

type Config struct {
	// ServerHost defines server host.
	//
	// Optional. Default: 0.0.0.0:8081
	ServerHost string

	// Logger defines logger engine that will be used.
	//
	// Optional. Default: nil (no log message will be displayed)
	Logger *logr.Logger
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]
	if cfg.ServerHost == "" {
		cfg.ServerHost = defaultHost
	}

	return cfg
}
