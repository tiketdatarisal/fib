package hawkeye

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/metrics"
)

type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// Host defines Telegraf host.
	//
	// Optional. Default: localhost
	Host string

	// Port defines Telegraf StatsD port.
	//
	// Optional. Default: 8125
	Port int

	// AppName defines application name.
	//
	// Optional.
	AppName string

	// HawkEye defines hawk-eye engines. When defined this variable will be used,
	// otherwise will be instantiated using predefined Host and Port.
	//
	// Optional. Default: nil
	HawkEye metrics.MonitorStatsd
}

func (c *Config) Initialize() error {
	if c.HawkEye != nil {
		return nil
	}

	if c.Host == "" || c.Port <= 0 {
		return ErrHawkEyeFailedToInitialize
	}

	var err error
	c.HawkEye, err = metrics.NewMonitor(c.Host, fmt.Sprintf("%d", c.Port), defaultHawkEyeInstanceName)
	if err != nil {
		return fmt.Errorf(errMsgHawkEyeFailedToInitialize, err)
	}

	return nil
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]
	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}

	err := cfg.Initialize()
	if err != nil {
		panic(fmt.Errorf(errMsgHawkEyeFailedToInitialize, err))
	}

	return cfg
}
