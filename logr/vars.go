package logr

import "errors"

var (
	ErrLoggerRequired = errors.New("logger must be initialized")

	ConfigDefault = Config{
		Next: nil,
	}
)
