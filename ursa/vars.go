package ursa

import (
	"context"
	"errors"
	"sync"
)

const (
	userStatusInactive = "inactive"
	userStatusActive   = "active"
	authenticatedUser  = "authenticated-user"
)

var (
	u     *ursa
	mutex sync.Mutex
	bg    context.Context

	ErrHostRequired        = errors.New("host for URSA GRPC is required")
	ErrHostCannotBeReached = errors.New("host cannot be reached")
	ErrUrsaNotInitialized  = errors.New("URSA not initialized")
	ErrEmailNotValid       = errors.New("email address is not valid")
	ErrUserUnauthorized    = errors.New("user unauthorized")
	ErrUserDeactivated     = errors.New("user deactivated")

	ConfigDefault = Config{
		Next: nil,
	}
)
