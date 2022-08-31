package config

import (
	"fmt"

	huner "github.com/kamva/hexa-tuner"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
)

// C is default config to use in project.
var C *Config

func New() (*Config, error) {
	prefix := huner.EnvKeysPrefix()
	files := huner.GetConfigFilePaths(huner.ConfigFilePathsOpts{
		AppName:       AppName,
		ServiceName:   ServiceName,
		HomePath:      ProjectRootPath(),
		FileName:      FileName,
		FileExtension: FileExtension,
		Environment:   huner.Environment(prefix),
	})

	v, err := huner.NewViperConfigDriver(prefix, files)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	cfg := Config{Config: v}
	err = v.Unmarshal(&cfg)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	hlog.Debug("config", hlog.String("conf", fmt.Sprintf("%+v", cfg)))

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SetDefaultConfig set config as global config.
func SetDefaultConfig(cfg *Config) {
	C = cfg
}
