package graceful

const (
	defaultHost = "0.0.0.0:8080"
)

var (
	ConfigDefault = Config{
		ServerHost: defaultHost,
		Logger:     nil,
	}
)
