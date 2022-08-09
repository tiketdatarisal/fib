package ursa

import "github.com/gofiber/fiber/v2"

type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// Host defines an address for URSA GRPC endpoint.
	//
	// Required.
	Host string
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		panic(ErrHostRequired)
	}

	cfg := config[0]
	if cfg.Host == "" {
		panic(ErrHostRequired)
	}

	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}

	return cfg
}
