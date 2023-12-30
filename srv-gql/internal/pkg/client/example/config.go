package example

import (
	"go.uber.org/zap"
	"skeleton/shared/bootstrap"
)

type clientConfig struct {
	Host string `env:"GRPC_HOST"`
}

func (s *clientConfig) DisplayFieldMap() map[string]interface{} {
	return bootstrap.CfgDisplayFieldMap(s)
}

func makeConfigFromEnv(logger *zap.Logger) (clientConfig, error) {
	cfg := clientConfig{}

	err := bootstrap.MakeConfigFromEnv(logger, &cfg)

	return cfg, err
}
