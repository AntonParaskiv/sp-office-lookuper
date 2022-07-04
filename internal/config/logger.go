package config

import (
	"sp-office-lookuper/internal/app"

	"github.com/kelseyhightower/envconfig"
)

type LoggerConfig struct {
	ServiceName string `envconfig:"SERVICE_NAME"`
	Level       int    `envconfig:"LEVEL" default:"4"` // info level by default
	Destination string `envconfig:"DESTINATION" default:""`
	Host        string `envconfig:"HOST" default:""`
	Port        int    `envconfig:"PORT" default:"12201"`
}

func (c *LoggerConfig) Prepare() error {
	err := envconfig.Process(app.LoggerConfigPrefix, c)
	if err != nil {
		return err
	}

	return nil
}

func PrepareLoggerConfig() (*LoggerConfig, error) {
	config := &LoggerConfig{}
	err := config.Prepare()
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *LoggerConfig) GetServiceName() string {
	return c.ServiceName
}

func (c *LoggerConfig) GetLevel() int {
	return c.Level
}

func (c *LoggerConfig) GetDestination() string {
	return c.Destination
}

func (c *LoggerConfig) GetHost() string {
	return c.Host
}

func (c *LoggerConfig) GetPort() int {
	return c.Port
}
