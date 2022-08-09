package ursa

import "github.com/gofiber/fiber/v2"

type Authenticator interface {
	Auth(scopes ...string) fiber.Handler
}

type DefaultAuthenticator struct{}

func (a *DefaultAuthenticator) Auth(...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
