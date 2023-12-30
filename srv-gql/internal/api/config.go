package api

import (
	"go.uber.org/zap"
	"skeleton/shared/bootstrap"
)

type serverConfig struct {
	Port           string `env:"PORT" envDefault:"80"`
	HealthPort     string `env:"HEALTH_PORT" envDefault:"8080"`
	SwaggerEnabled int    `env:"SWAGGER_ENABLED" envDefault:"0"`
}

func (s *serverConfig) DisplayFieldMap() map[string]interface{} {
	return bootstrap.CfgDisplayFieldMap(s)
}

func makeConfigFromEnv(logger *zap.Logger) (serverConfig, error) {
	cfg := serverConfig{}

	err := bootstrap.MakeConfigFromEnv(logger, &cfg)

	return cfg, err
}
