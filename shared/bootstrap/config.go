package bootstrap

import (
	"github.com/caarlos0/env/v8"
	"go.uber.org/zap"
	"reflect"
)

type DebugConfig struct {
	IsDebugMode bool `env:"DEBUG" envDefault:"0"`
	LogLevel    int  `env:"LOG_LEVEL" envDefault:"0"`
}

func MakeDebugConfigFromEnv() DebugConfig {
	cfg := DebugConfig{}

	_ = env.ParseWithOptions(&cfg, env.Options{RequiredIfNoDef: false})

	return cfg
}

const (
	cfgDisplayFieldHidden = "<REDACTED>"
)

type Cfg interface {
	DisplayFieldMap() map[string]interface{}
}

func CfgDisplayFieldMap(cfg Cfg, hiddenFields ...string) map[string]interface{} {
	rElem := reflect.ValueOf(cfg).Elem()

	elems := make(map[string]interface{})

	for i := 0; i < rElem.NumField(); i++ {
		elems[rElem.Type().Field(i).Name] = rElem.Field(i).Interface()
	}

	for _, val := range hiddenFields {
		if _, ok := elems[val]; ok {
			elems[val] = cfgDisplayFieldHidden
		}
	}

	return elems
}

func MakeConfigFromEnv(logger *zap.Logger, obj Cfg) error {
	rElem := reflect.TypeOf(obj).Elem()

	logger = logger.With(zap.String("path", rElem.PkgPath()), zap.String("obj", rElem.Name()))

	if err := env.ParseWithOptions(obj, env.Options{RequiredIfNoDef: true}); err != nil {
		logger.Error("bootstrap cfg read err", zap.Error(err))
		return err
	}

	logger.Info("bootstrap cfg read", zap.Any("data", obj.DisplayFieldMap()))

	return nil
}
