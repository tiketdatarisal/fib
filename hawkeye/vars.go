package hawkeye

import "errors"

const (
	defaultHost                     = "localhost"
	defaultPort                     = 8125
	defaultHawkEyeInstanceName      = "hawkEyeMiddleware"
	apiLatencyMetric                = "ApiLatency"
	errMsgHawkEyeFailedToInitialize = "could not instantiate HawkEye client: %w"
)

var (
	ErrHawkEyeFailedToInitialize = errors.New("could not initialize HawkEye client: Host/Port empty")

	ConfigDefault = Config{
		Next:    nil,
		Host:    defaultHost,
		Port:    defaultPort,
		HawkEye: nil,
	}
)
