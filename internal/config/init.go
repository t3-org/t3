package config

import (
	"os"

	huner "github.com/kamva/hexa-tuner"
	"github.com/kamva/tracer"
)

// C is default config to use in project.
var C *Config

func New() (*Config, error) {
	prefix := huner.EnvKeysPrefix()
	env := os.Getenv(huner.EnvironmentKey(prefix))
	files := huner.ConfigFilePaths(huner.ConfigFilePathsOptions{
		EtcPath:       huner.EtcPath(AppName),
		HomePath:      appRootPath(),
		FileName:      FileName,
		FileExtension: FileExtension,
		Environment:   env,
	}, env == huner.EnvironmentTest)

	v, err := huner.NewViperConfig(prefix, files)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	var cfg Config
	err = v.Unmarshal(&cfg, huner.ViperDecoderTagName("json"))
	if err != nil {
		return nil, tracer.Trace(err)
	}

	//prettyCfg, err := yaml.Marshal(cfg)
	//if err != nil {
	//	return nil, tracer.Trace(err)
	//}
	//hlog.Debug("config", hlog.String("values", string(prettyCfg)))

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	if err := cfg.setDefaults(); err != nil {
		return nil, tracer.Trace(err)
	}

	return &cfg, nil
}

// SetDefaultConfig set config as global config.
func SetDefaultConfig(cfg *Config) {
	C = cfg
}
