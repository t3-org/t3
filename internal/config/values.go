package config

import (
	"fmt"
	"path"

	hecho "github.com/kamva/hexa-echo"
	huner "github.com/kamva/hexa-tuner"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/hexa/hlog/logdriver"
	"github.com/labstack/echo/v4/middleware"
)

const (
	// AppName is the name of the app to load its configurations.
	AppName = "space"
	// FileName is name of the configuration file
	FileName = "config"
	// FileExtension is extension of the configuration file.
	FileExtension = "yaml"
)

type Config struct {
	// Test contains test config which we don't need to provide it in non-test environments.
	Test *Test `json:"test"`

	AppName      string  `json:"app_name"`
	InstanceName string  `json:"instance_name"`
	MachineID    *uint16 `json:"machine_id"` // Optional. use to generate ID values.

	//--------------------------------
	// General Configs
	//--------------------------------

	// Environment is just for log, tracing,... you could set it to
	// something like prod, dev, local,...
	Environment       string   `json:"environment"`
	Debug             bool     `json:"debug"`
	LogStack          []string `json:"log_stack"`
	LogLevel          string   `json:"log_level"`
	SentryDSN         string   `json:"sentry_dsn"`
	SentryEnvironment string   `json:"sentry_environment"`
	TranslateFiles    []string `json:"translate_files"`
	FallbackLanguages []string `json:"fallback_languages"`

	ListeningIP        string        `json:"listening_ip"`
	Port               int           `json:"port"`
	DB                 DB            `json:"db"`
	ProbeServerAddress string        `json:"probe_server_address"`
	OpenTelemetry      OpenTelemetry `json:"open_telemetry"`

	AssetsBasePath    string `json:"assets_base_path"`
	ResourcesBasePath string `json:"resources_base_path"`

	//--------------------------------
	// HTTP server configs
	//--------------------------------
	EchoLogLevel      string   `json:"echo_log_level"`
	AllowOrigins      []string `json:"allow_origins"`
	AllowHeaders      []string `json:"allow_headers"`
	AllowCredentials  bool     `json:"allow_credentials"`
	AllowMethods      []string `json:"allow_methods"`
	CorsMaxAgeSeconds int      `json:"cors_max_age_seconds"`
	AuthTokenCookie   Cookie   `json:"auth_token_cookie"`
	CSRFCookie        Cookie   `json:"csrf_cookie"`

	// limit request response
	RequestReadTimeoutMs       int `json:"request_read_timeout_ms"`
	RequestReadHeaderTimeoutMs int `json:"request_read_header_timeout_ms"`
	ResponseWriteTimeoutMs     int `json:"response_write_timeout_ms"`
	// IdleTimeout is the maximum amount of time to wait for the next request when keep-alives are enabled.
	ConnectionIdleTimeoutMs int   `json:"connection_idle_timeout_ms"`
	MaxHeaderSizeKb         int   `json:"max_header_size_kb"` // 0 means unlimited.
	MaxBodySizeKb           int64 `json:"max_body_size_kb"`   // 0 means unlimited.

	AsynqConfig AsynqConfig `json:"asynq_config"`

	RedisAddress  string       `json:"redis_address"`
	RedisPassword string       `json:"redis_password"`
	RedisDB       int          `json:"redis_db"`
	Matrix        MatrixConfig `json:"matrix"`
}

func (c *Config) validate() error {
	return nil
}

func (c *Config) setDefaults() error {
	if c.AssetsBasePath == "" {
		c.AssetsBasePath = path.Join(appRootPath(), "assets")
	}

	if c.ResourcesBasePath == "" {
		c.ResourcesBasePath = path.Join(appRootPath(), "res")
	}

	return nil
}

func (c *Config) PathToResource(p string) string {
	return path.Join(c.ResourcesBasePath, p)
}

func (c *Config) PathToAssets(p string) string {
	return path.Join(c.ResourcesBasePath, p)
}

func (c *Config) I18nPath() string {
	return c.PathToResource("i18n")
}

func (c *Config) StackLoggerOptions() ([]string, logdriver.StackOptions) {
	return c.LogStack, logdriver.StackOptions{
		Level:      hlog.LevelFromString(c.LogLevel),
		ZapConfig:  logdriver.DefaultZapConfig(c.Debug, hlog.ZapLevel(hlog.LevelFromString(c.LogLevel)), "json"),
		SentryOpts: c.Sentry(),
	}
}

func (c *Config) Sentry() *logdriver.SentryOptions {
	return &logdriver.SentryOptions{
		DSN:         c.SentryDSN,
		Debug:       c.Debug,
		Environment: c.SentryEnvironment,
	}
}

func (c *Config) TranslateOptions() huner.TranslateOpts {
	return huner.TranslateOpts{
		Files:         c.TranslateFiles,
		FallbackLangs: c.FallbackLanguages,
	}
}

func (c *Config) ListeningAddress() string {
	return fmt.Sprintf("%s:%v", c.ListeningIP, c.Port)
}

func (c *Config) CSRFConfig() middleware.CSRFConfig {
	return middleware.CSRFConfig{
		Skipper:        hecho.CSRFSkipperByAuthTokenLocation,
		TokenLength:    32,
		TokenLookup:    middleware.DefaultCSRFConfig.TokenLookup,
		ContextKey:     c.CSRFCookieContextKey(),
		CookieName:     c.CSRFCookie.Name,
		CookieDomain:   c.CSRFCookie.Domain,
		CookiePath:     "/",
		CookieMaxAge:   c.CSRFCookie.MaxAge,
		CookieSecure:   c.CSRFCookie.Secure,
		CookieHTTPOnly: c.CSRFCookie.HttpOnly,
	}
}

func (c *Config) CSRFCookieContextKey() string {
	return middleware.DefaultCSRFConfig.ContextKey
}

func (c *Config) ApiDocsRouteTemplatePath() string {
	return c.PathToResource("doc_route_template.tpl")
}

func (c *Config) ApiDocExportFilePath() string {
	return path.Join(appRootPath(), "internal/router/api/doc/api_docs.go")
}

func (c *Config) ApiDocsDestinationFiles() []string {
	return []string{
		path.Join(appRootPath(), "docs/api/api_docs.json"),
		path.Join(appRootPath(), "docs/api/api_docs.yaml"),
	}
}
