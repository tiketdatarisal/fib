package logr

import (
	"github.com/gofiber/fiber/v2"
	"sync"
	"time"
)

// New return a new Logger middleware.
func New(config ...Config) fiber.Handler {
	var once sync.Once
	var mutex sync.Mutex
	var errHandler fiber.ErrorHandler

	cfg := configDefault(config...)
	return func(c *fiber.Ctx) error {
		once.Do(func() {
			errHandler = c.App().Config().ErrorHandler
		})

		var start time.Time
		start = time.Now()

		chainErr := c.Next()
		if chainErr != nil {
			if err := errHandler(c, chainErr); err != nil {
				mutex.Lock()
				cfg.Logger.Error(err, "",
					"latency", time.Now().Sub(start).Round(time.Microsecond),
					"responseCode", c.Response().StatusCode(),
					"ip", c.IP(),
					"method", c.Method(),
					"path", c.Path(),
				)
				mutex.Unlock()
				_ = c.SendStatus(fiber.StatusInternalServerError)
				return nil
			}
		}

		mutex.Lock()
		cfg.Logger.Info("",
			"latency", time.Now().Sub(start).Round(time.Microsecond),
			"responseCode", c.Response().StatusCode(),
			"ip", c.IP(),
			"method", c.Method(),
			"path", c.Path(),
		)
		mutex.Unlock()
		return nil
	}
}
