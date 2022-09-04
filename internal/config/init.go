package config

import (
	huner "github.com/kamva/hexa-tuner"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"gopkg.in/yaml.v2"
)

// C is default config to use in project.
var C *Config

func New() (*Config, error) {
	prefix := huner.EnvKeysPrefix()
	env := huner.Environment(prefix)
	files := huner.ConfigFilePaths(huner.ConfigFilePathsOptions{
		EtcPath:       huner.EtcPath(AppName),
		HomePath:      ProjectRootPath(),
		FileName:      FileName,
		FileExtension: FileExtension,
		Environment:   env,
	}, env == huner.EnvironmentTest)

	v, err := huner.NewViperConfigDriver(prefix, files)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	cfg := Config{Config: v}
	err = v.Unmarshal(&cfg)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	prettyCfg, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	hlog.Debug("config", hlog.String("values", string(prettyCfg)))

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SetDefaultConfig set config as global config.
func SetDefaultConfig(cfg *Config) {
	C = cfg
}
