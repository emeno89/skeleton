package bootstrap

import (
	"go.elastic.co/apm/module/apmzap/v2"
	"go.elastic.co/ecszap"
	_ "go.elastic.co/ecszap"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func ZapLogger(cfg DebugConfig) *zap.Logger {
	c := &apmzap.Core{}

	encoderConfig := ecszap.EncoderConfig{
		EnableName:       true,
		EnableCaller:     true,
		EnableStackTrace: true,
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeName:       zapcore.FullNameEncoder,
		EncodeLevel:      zapcore.LowercaseLevelEncoder,
		EncodeDuration:   zapcore.MillisDurationEncoder,
		EncodeCaller:     zapcore.FullCallerEncoder,
	}

	level := zapcore.Level(cfg.LogLevel)

	if cfg.IsDebugMode {
		level = zapcore.DebugLevel
	}

	ecs := c.WrapCore(ecszap.NewCore(encoderConfig, os.Stdout, level))

	logger := zap.New(ecs, zap.WrapCore(c.WrapCore))

	return logger.With(
		zap.String("service.environment", os.Getenv("ELASTIC_APM_ENVIRONMENT")),
		zap.String("service.name", os.Getenv("ELASTIC_APM_SERVICE_NAME")),
	)
}

func FxLogger(logger *zap.Logger, debugCfg DebugConfig) fxevent.Logger {
	if !debugCfg.IsDebugMode {
		logger = zap.NewNop()
	}

	return &fxevent.ZapLogger{Logger: logger}
}
