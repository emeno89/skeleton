package security

import (
	"go.uber.org/zap"
	"skeleton/shared/bootstrap"
)

type jwtConfig struct {
	Secret string `env:"JWT_SECRET"`
}

func (s *jwtConfig) DisplayFieldMap() map[string]interface{} {
	return bootstrap.CfgDisplayFieldMap(s, "Secret")
}

func makeConfigFromEnv(logger *zap.Logger) (jwtConfig, error) {
	cfg := jwtConfig{}

	err := bootstrap.MakeConfigFromEnv(logger, &cfg)

	return cfg, err
}
