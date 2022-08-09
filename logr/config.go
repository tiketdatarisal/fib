package logr

import (
	"github.com/go-logr/logr"
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// Logger defines an instance of log framework that implement logr.Logger interface.
	//
	// Required. This will panic if not initiated.
	Logger *logr.Logger
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		panic(ErrLoggerRequired)
	}

	cfg := config[0]
	if cfg.Logger == nil {
		panic(ErrLoggerRequired)
	}

	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}

	return cfg
}
