package graceful

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"os"
	"os/signal"
	"syscall"
)

func Shutdown(app *fiber.App, config ...Config) {
	cfg := configDefault(config...)

	if app == nil {
		if cfg.Logger != nil {
			(*cfg.Logger).Error(errors.New("app is nil; Fiber app is required"),
				"required parameter is not available")
		}
		return
	}

	q := make(chan os.Signal)
	signal.Notify(q, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-q

		if cfg.Logger != nil {
			(*cfg.Logger).Info("shutting down")
		}
		if err := app.Shutdown(); err != nil {
			if cfg.Logger != nil {
				(*cfg.Logger).Error(err, "could not shutdown Fiber gracefully")
			}
		}
	}()

	if cfg.Logger != nil {
		(*cfg.Logger).Info("starting server", "host", cfg.ServerHost)
	}
	if err := app.Listen(cfg.ServerHost); err != nil {
		if cfg.Logger != nil {
			(*cfg.Logger).Error(err, "could not start server")
		}
	}
}
