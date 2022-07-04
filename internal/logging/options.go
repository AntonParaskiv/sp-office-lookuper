package logging

type Config interface {
	GetServiceName() string
	GetLevel() int
	GetDestination() string
	GetHost() string
	GetPort() int
}

type Option func(*Logger)

func WithConfig(conf Config) Option {
	return func(logger *Logger) {
		if conf != nil {
			logger.serviceName = conf.GetServiceName()
			logger.level = conf.GetLevel()
			logger.destination = conf.GetDestination()
			logger.host = conf.GetHost()
			logger.port = conf.GetPort()
		}
	}
}
