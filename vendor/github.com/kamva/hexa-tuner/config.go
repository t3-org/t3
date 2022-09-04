package huner

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hconf"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"github.com/spf13/viper"
)

const EnvironmentTest = "test" // We can use this reserved env name for test environments. it's completely optional.

type ConfigFilePathsOptions struct {
	EtcPath       string // e.g., senna
	HomePath      string // e.g., /home/mehran/senna/order
	FileName      string // e.g., config
	FileExtension string // e.g., json or yaml
	Environment   string // (optional) e.g., staging
}

// EtcPath returns path to the /etc directory.
// Later we can make its base path os dependant.
func EtcPath(p string) string {
	return path.Join("/etc", p)
}

func EnvKeysPrefix() string {
	return os.Getenv("HEXA_CONF_PREFIX")
}

func Environment(prefix string) string {
	key := "ENV"
	if prefix != "" {
		key = fmt.Sprintf("%s_%s", prefix, key)
	}
	return os.Getenv(key)
}

// ConfigFilePaths generates config path as follow:
// - {EtcPath}/{configFile.configExtension}
// - {HomePath}/{configFile.configExtension}
// - {HomePath}/{configFile.{environment}.configExtension}
// - {HomePath}/.env
// - {HomePath}/.{environment}.env
// If you set onlyEnvDependant to true, it just includes files that include the
// environment name in their config file name.
func ConfigFilePaths(o ConfigFilePathsOptions, onlyEnvDependant bool) []string {
	confFile := fmt.Sprintf("%s.%s", o.FileName, o.FileExtension)

	var files []string

	if !onlyEnvDependant { // including crossEnv paths.
		files = append(files,
			path.Join(o.EtcPath, confFile),  // {EtcPath}/{configFile.configExtension}
			path.Join(o.HomePath, confFile), // {HomePath}/{configFile.configExtension}
			path.Join(o.HomePath, ".env"),   // {HomePath}/.env
		)
	}

	if o.Environment != "" {
		files = append(files,
			// {HomePath}/{configFile}.{env}.{extension}
			path.Join(o.HomePath, fmt.Sprintf("%s.%s.%s", o.FileName, o.Environment, o.FileExtension)),
			// {HomePath}/.{env}.env
			path.Join(o.HomePath, fmt.Sprintf(".%s.env", o.Environment)),
		)
	}

	var existedFiles []string

	for _, f := range files {
		if gutil.FileExists(f) {
			existedFiles = append(existedFiles, f)
		}
	}

	hlog.Debug("generated config file paths",
		hlog.Any("search_paths", files),
		hlog.Any("existed_paths", existedFiles),
		hlog.String("config", fmt.Sprintf("%+v", o)),
	)

	return existedFiles
}

// NewViperConfigDriver returns new instance of the viper driver for hexa config
func NewViperConfigDriver(envPrefix string, files []string) (hexa.Config, error) {
	v := viper.New()

	if len(files) == 0 {
		return nil, tracer.Trace(errors.New("at least one config files should be exists"))
	}

	isFirst := true
	for _, f := range files {
		v.SetConfigFile(f)

		if isFirst {
			isFirst = false
			if err := v.ReadInConfig(); err != nil {
				return nil, tracer.Trace(err)
			}
			continue
		}

		if err := v.MergeInConfig(); err != nil {
			return nil, tracer.Trace(err)
		}
	}

	v.AllowEmptyEnv(true)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()

	return hconf.NewViperDriver(v), nil
}
