package huner

import (
	"errors"
	"fmt"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hconf"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
)

type ConfigFilePathsOpts struct {
	AppName       string // e.g., senna
	ServiceName   string // e.g., order
	HomePath      string // e.g., /home/mehran/senna/order
	FileName      string // e.g., config
	FileExtension string // e.g., json or yaml
	Environment   string // (optional) e.g., staging
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

// GetConfigFilePaths generates config path as follow:
// - /etc/{appName}/{configFile.configExtension}
// - /etc/{appName}/{serviceName.configExtension}
// - {HomePath}/{configFile.configExtension}
// - {HomePath}/.env
// - {HomePath}/.{environment}.env
func GetConfigFilePaths(o ConfigFilePathsOpts) []string {
	configFile := fmt.Sprintf("%s.%s", o.FileName, o.FileExtension)
	msConfigFile := fmt.Sprintf("%s.%s", o.ServiceName, o.FileExtension)

	files := []string{
		path.Join("/etc", o.AppName, configFile),
		path.Join("/etc", o.AppName, msConfigFile),
		path.Join(o.HomePath, configFile),
		path.Join(o.HomePath, ".env"),
	}

	if o.Environment != "" {
		files = append(files, path.Join(o.HomePath, fmt.Sprintf(".%s.env", o.Environment)))
	}

	var existedFiles []string

	for _, f := range files {
		if gutil.FileExists(f) {
			existedFiles = append(existedFiles, f)
		}
	}

	hlog.Debug("generated config file paths",
		hlog.Any("available_paths", files),
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
