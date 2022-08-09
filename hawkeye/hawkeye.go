package hawkeye

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/metrics"
	"sync"
	"time"
)

func captureLatency(cfg *Config, latency time.Duration, statusCode int, ip, method, path, appName string, isSuccess bool) {
	tags := map[string]interface{}{
		"latency":    latency,
		"statusCode": statusCode,
		"callerIP":   ip,
		"method":     method,
		"path":       path,
	}
	if appName != "" {
		tags["app"] = appName
	}

	if isSuccess {
		_ = cfg.HawkEye.CustomMonitorLatency(apiLatencyMetric, metrics.API_OUT, metrics.Success, statusCode, tags, latency)
	} else {
		_ = cfg.HawkEye.CustomMonitorLatency(apiLatencyMetric, metrics.API_OUT, metrics.Failed, statusCode, tags, latency)
	}
}

// New return a new HawkEye middleware.
func New(config ...Config) fiber.Handler {
	var once sync.Once
	var errHandler fiber.ErrorHandler

	cfg := configDefault(config...)
	return func(c *fiber.Ctx) error {
		if (cfg.Next != nil && cfg.Next(c)) || cfg.HawkEye == nil {
			return c.Next()
		}

		once.Do(func() {
			errHandler = c.App().Config().ErrorHandler
		})

		start := time.Now()
		chainErr := c.Next()
		latency := time.Now().Sub(start)
		statusCode := c.Response().StatusCode()

		if chainErr != nil {
			if err := errHandler(c, chainErr); err != nil {
				go captureLatency(&cfg, latency, statusCode, c.IP(), c.Method(), c.Path(), cfg.AppName, false)
			}
		} else {
			go captureLatency(&cfg, latency, statusCode, c.IP(), c.Method(), c.Path(), cfg.AppName, true)
		}

		return nil
	}
}
