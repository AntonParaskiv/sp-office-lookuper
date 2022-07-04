package config

import (
	"time"

	"sp-office-lookuper/internal/app"

	"github.com/kelseyhightower/envconfig"
)

type JaegerConfig struct {
	ServiceName string `envconfig:"SERVICE_NAME"`
	Host        string `envconfig:"HOST" default:"localhost"`
	Port        int    `envconfig:"PORT" default:"6831"`

	SamplerType         string        `envconfig:"SAMPLER_TYPE" default:"ratelimiting"`
	SamplerParam        float64       `envconfig:"SAMPLER_PARAM" default:"100"` // limit per second
	BufferFlushInterval time.Duration `envconfig:"BUFFER_FLUSH_INTERVAL" default:"1s"`
}

func (c *JaegerConfig) Prepare() error {
	err := envconfig.Process(app.JaegerConfigPrefix, c)
	if err != nil {
		return err
	}

	return nil
}

func PrepareJaegerConfig() (*JaegerConfig, error) {
	config := &JaegerConfig{}
	err := config.Prepare()
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *JaegerConfig) GetServiceName() string {
	return c.ServiceName
}

func (c *JaegerConfig) GetHost() string {
	return c.Host
}

func (c *JaegerConfig) GetPort() int {
	return c.Port
}

func (c *JaegerConfig) GetSamplerType() string {
	return c.SamplerType
}

func (c *JaegerConfig) GetSamplerParam() float64 {
	return c.SamplerParam
}

func (c *JaegerConfig) GetBufferFlushInterval() time.Duration {
	return c.BufferFlushInterval
}
