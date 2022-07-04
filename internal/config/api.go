package config

import (
	"fmt"

	"sp-office-lookuper/internal/app"

	"github.com/kelseyhightower/envconfig"
)

type APIConfig struct {
	HTTPMaxConnection int    `envconfig:"HTTP_MAX_CONNECTIONS" default:"1000"`
	HTTPListenHost    string `envconfig:"HTTP_HOST" default:"0.0.0.0"`
	HTTPListenPort    int    `envconfig:"HTTP_PORT" default:"8080"`
	HTTPPprofPort     int    `envconfig:"HTTP_PPROF_PORT" default:"8180"`

	GRPCListenHost string `envconfig:"GRPC_HOST" default:"0.0.0.0"`
	GRPCListenPort int    `envconfig:"GRPC_PORT" default:"8090"`

	HTTPListenAddress string
	HTTPPprofAddress  string
	GRPCListenAddress string

	*LoggerConfig
	*JaegerConfig
}

func (c *APIConfig) Prepare() error {
	err := envconfig.Process(app.ServiceName, c)
	if err != nil {
		return err
	}

	err = c.prepareEmbeded()
	if err != nil {
		return err
	}

	c.HTTPListenAddress = fmt.Sprintf("%s:%d", c.HTTPListenHost, c.HTTPListenPort)
	c.HTTPPprofAddress = fmt.Sprintf("%s:%d", c.HTTPListenHost, c.HTTPPprofPort)
	c.GRPCListenAddress = fmt.Sprintf("%s:%d", c.GRPCListenHost, c.GRPCListenPort)

	return nil
}

func PrepareAPIConfig() (*APIConfig, error) {
	config := &APIConfig{}
	err := config.Prepare()
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *APIConfig) prepareEmbeded() error {
	loggerConfig, err := PrepareLoggerConfig()
	if err != nil {
		return err
	}
	c.LoggerConfig = loggerConfig

	jaegerConfig, err := PrepareJaegerConfig()
	if err != nil {
		return err
	}
	c.JaegerConfig = jaegerConfig

	return nil
}
