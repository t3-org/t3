package config

import (
	"fmt"
	"path"

	"github.com/kamva/hexa"
	hecho "github.com/kamva/hexa-echo"
	huner "github.com/kamva/hexa-tuner"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/hexa/hlog/logdriver"
	"github.com/labstack/echo/v4/middleware"
)

const (
	// AppName is name of the project to load its config...
	AppName = "space"
	// ServiceName is name of the microservice to load its config...
	ServiceName = "space"
	// FileName is name of the config file
	FileName = "config"
	// FileExtension is extension of the config file.
	FileExtension = "yaml"
)

type Config struct {
	hexa.Config

	AppName      string `json:"app_name" mapstructure:"app_name"`
	InstanceName string `json:"instance_name" mapstructure:"instance_name"`
	// Environment is just for log, tracing,... you could set it to
	// something like prod, dev, local,...
	Environment string `json:"environment" mapstructure:"environment"`

	InitializeAdminPhone string `json:"initialize_admin_phone" mapstructure:"initialize_admin_phone"`

	//--------------------------------
	// General Configs
	//--------------------------------

	Debug              bool          `json:"debug" mapstructure:"debug"`
	ListeningIP        string        `json:"listening_ip" mapstructure:"listening_ip"`
	Port               int           `json:"port" mapstructure:"port"`
	DB                 DB            `json:"db" mapstructure:"db"`
	ProbeServerAddress string        `json:"probe_server_address" mapstructure:"probe_server_address"`
	OpenTelemetry      OpenTelemetry `json:"open_telemetry" mapstructure:"open_telemetry"`

	AssetsBasePath    string `json:"assets_base_path" mapstructure:"assets_base_path"`
	ResourcesBasePath string `json:"resources_base_path" mapstructure:"resources_base_path"`

	LogStack          []string `json:"log_stack" mapstructure:"log_stack"`
	SentryDSN         string   `json:"sentry_dsn" mapstructure:"sentry_dsn"`
	SentryEnvironment string   `json:"sentry_environment" mapstructure:"sentry_environment"`
	TranslateFiles    []string `json:"translate_files" mapstructure:"translate_files"`
	FallbackLanguages []string `json:"fallback_languages" mapstructure:"fallback_languages"`
	LogLevel          string   `json:"log_level" mapstructure:"log_level"`
	GRPCLogVerbosity  int      `json:"grpc_log_verbosity" mapstructure:"grpc_log_verbosity"`

	//--------------------------------
	// HTTP server configs
	//--------------------------------
	EchoLogLevel      string   `json:"echo_log_level" mapstructure:"echo_log_level"`
	AllowOrigins      []string `json:"allow_origins" mapstructure:"allow_origins"`
	AllowHeaders      []string `json:"allow_headers" mapstructure:"allow_headers"`
	AllowCredentials  bool     `json:"allow_credentials" mapstructure:"allow_credentials"`
	AllowMethods      []string `json:"allow_methods" mapstructure:"allow_methods"`
	CorsMaxAgeSeconds int      `json:"cors_max_age_seconds" mapstructure:"cors_max_age_seconds"`
	AuthTokenCookie   Cookie   `json:"auth_token_cookie" mapstructure:"auth_token_cookie"`
	CSRFCookie        Cookie   `json:"csrf_cookie" mapstructure:"csrf_cookie"`

	// limit request response
	RequestReadTimeoutMs       int `json:"request_read_timeout_ms" mapstructure:"request_read_timeout_ms"`
	RequestReadHeaderTimeoutMs int `json:"request_read_header_timeout_ms" mapstructure:"request_read_header_timeout_ms"`
	ResponseWriteTimeoutMs     int `json:"response_write_timeout_ms" mapstructure:"response_write_timeout_ms"`
	// IdleTimeout is the maximum amount of time to wait for the next request when keep-alives are enabled.
	ConnectionIdleTimeoutMs int   `json:"connection_idle_timeout_ms" mapstructure:"connection_idle_timeout_ms"`
	MaxHeaderSizeKb         int   `json:"max_header_size_kb" mapstructure:"max_header_size_kb"` // 0 means unlimited.
	MaxBodySizeKb           int64 `json:"max_body_size_kb" mapstructure:"max_body_size_kb"`     // 0 means unlimited.

	AsynqConfig AsynqConfig `json:"asynq_config" mapstructure:"asynq_config"`

	RedisAddress  string `json:"redis_address" mapstructure:"redis_address"`
	RedisPassword string `json:"redis_password" mapstructure:"redis_password"`
	RedisDB       int    `json:"redis_db" mapstructure:"redis_db"`
	Cache         Cache  `json:"cache" mapstructure:"cache"`
}

func (c *Config) validate() error {
	return nil
}

func (c *Config) ServiceName() string {
	return ServiceName
}

func (c Config) I18nPath() string {
	return ResourcePath(c.ResourcesBasePath, "/i18n")
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

func (c *Config) GRPCConfigs() huner.GRPCConfigs {
	return huner.GRPCConfigs{
		Debug:        c.Debug,
		LogVerbosity: c.GRPCLogVerbosity,
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
	return ResourcePath(c.ResourcesBasePath, "doc_route_template.tpl")
}

func (c *Config) ApiDocExportFilePath() string {
	return path.Join(ProjectRootPath(), "internal/router/api/doc/api_docs.go")
}

func (c *Config) ApiDocsDestinationFiles() []string {
	return []string{
		path.Join(ProjectRootPath(), "docs/api/api_docs.json"),
		path.Join(ProjectRootPath(), "docs/api/api_docs.yaml"),
	}
}
